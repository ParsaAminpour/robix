import * as anchor from "@coral-xyz/anchor";
import { Program } from "@coral-xyz/anchor";
import { Keypair, PublicKey } from "@solana/web3.js";
import { Raffle } from "../target/types/raffle";
import { BN } from "bn.js";
import web3 from "@solana/web3.js";
import { assert } from "chai";
import { Signature } from "ethers";
import { walkUpBindingElementsAndPatterns } from "typescript";
import { bs58 } from "@coral-xyz/anchor/dist/cjs/utils/bytes";
import { deserialize } from "borsh";
import { prototype } from "mocha";

const INIT_TREASURY_FUND = 0.01 * web3.LAMPORTS_PER_SOL;

export class CurrentFeed {
  is_init: number = 0;
  fee: number = 0;
  offset1: number = 0;
  offset2: number = 0;
  offset3: number = 0;
  offset4: number = 0;
  offset5: number = 0;
  offset6: number = 0;
  offset7: number = 0;
  offset8: number = 0;
  account1: number[] = Array.from({ length: 32 }, () => 1);
  account2: number[] = Array.from({ length: 32 }, () => 1);
  account3: number[] = Array.from({ length: 32 }, () => 1);
  fallback_account: number[] = Array.from({ length: 32 }, () => 1);
  bump: number = 0;

  constructor(
    fields:
      | {
          is_init: number;
          fee: number;
          offset1: number;
          offset2: number;
          offset3: number;
          offset4: number;
          offset5: number;
          offset6: number;
          offset7: number;
          offset8: number;
          account1: number[];
          account2: number[];
          account3: number[];
          fallback_account: number[];
          bump: number;
        }
      | undefined = undefined
  ) {
    if (fields) {
      this.is_init = fields.is_init;
      this.fee = fields.fee;
      this.offset1 = fields.offset1;
      this.offset2 = fields.offset2;
      this.offset3 = fields.offset3;
      this.offset4 = fields.offset4;
      this.offset5 = fields.offset5;
      this.offset6 = fields.offset6;
      this.offset7 = fields.offset7;
      this.offset8 = fields.offset8;
      this.account1 = fields.account1;
      this.account2 = fields.account2;
      this.account3 = fields.account3;
      this.fallback_account = fields.fallback_account;
      this.bump = fields.bump;
    }
  }
}

export const CurrentFeedSchema = new Map([
  [
    CurrentFeed,
    {
      kind: "struct",
      fields: [
        ["is_init", "u8"],
        ["fee", "u64"],
        ["offset1", "u8"],
        ["offset2", "u8"],
        ["offset3", "u8"],
        ["offset4", "u8"],
        ["offset5", "u8"],
        ["offset6", "u8"],
        ["offset7", "u8"],
        ["offset8", "u8"],
        ["account1", ["u8", 32]],
        ["account2", ["u8", 32]],
        ["account3", ["u8", 32]],
        ["fallback_account", ["u8", 32]],
        ["bump", "u8"],
      ],
    },
  ],
]);

const confirmTransaction = async (
  connection: web3.Connection,
  signature: web3.TransactionSignature,
  desiredConfirmationStatus: web3.TransactionConfirmationStatus = 'confirmed',
  timeout: number = 30000,
  pollInterval: number = 1000,
  searchTransactionHistory: boolean = false
): Promise<web3.SignatureStatus> => {
  const start = Date.now();
  while (Date.now() - start < timeout) {
      const { value: statuses } = await connection.getSignatureStatuses([signature], { searchTransactionHistory });
      if (!statuses || statuses.length === 0) {
          throw new Error('Failed to get signature status');
      }
      const status = statuses[0];
      if (status === null) {
          await new Promise(resolve => setTimeout(resolve, pollInterval));
          continue;
      }
      if (status.err) throw new Error(`Transaction failed: ${JSON.stringify(status.err)}`);

      if (status.confirmationStatus && status.confirmationStatus === desiredConfirmationStatus) return status;

      if (status.confirmationStatus === 'finalized') return status;

      await new Promise(resolve => setTimeout(resolve, pollInterval));
  }
  throw new Error(`Transaction confirmation timeout after ${timeout}ms`);
}

const raffleNumberBuffer = (value: bigint): Uint8Array => {
  const bytes = new Uint8Array(8);
  for (let i = 0; i < 8; i++) {
      bytes[i] = Number(value & BigInt(0xff));
      value = value >> BigInt(8);
  }
  return bytes;
}

describe("Testing Raffle instructions", () => {
  const rngProgram = new anchor.web3.PublicKey('9uSwASSU59XvUS8d1UeU8EwrEzMGFdXZvQ4JSEAfcS7k');
  process.env.ANCHOR_PROVIDER_URL = 'https://api.devnet.solana.com';
  // process.env.ANCHOR_WALLET = './key.json';

  const provider = anchor.AnchorProvider.env();
  anchor.setProvider(provider);
  
  const RAFFLE_SEED = Buffer.from("raffle");
  const TRACKER_SEED = Buffer.from("tracker")
  const TREASURY_SEED = Buffer.from("treasury");

  const program = anchor.workspace.Raffle as Program<Raffle>;
  
  let auth: PublicKey, participant: Keypair;
  let raffle_pda, tracker_pda, treasury_pda;
  let raffle_seeds, tracker_seeds, treasury_seeds;
  
  // Utility function for airdrops
  async function fundWallet(from: PublicKey, to: PublicKey, amount: number) {
      const tx = new anchor.web3.Transaction().add(
      anchor.web3.SystemProgram.transfer({
        fromPubkey: from,
        toPubkey: to,
        lamports: amount
      })
    );
    await provider.sendAndConfirm(tx);
  }
  
  before(async () => {
    auth = provider.wallet.publicKey;
    participant = web3.Keypair.generate();

    const connection = program.provider.connection;
    
    const INIT_RAFFLE_IDX = BigInt(1);
    raffle_seeds = [RAFFLE_SEED, raffleNumberBuffer(INIT_RAFFLE_IDX), auth.toBuffer()];
    raffle_pda = web3.PublicKey.findProgramAddressSync(
      raffle_seeds,
      program.programId
    );

    tracker_seeds = [TRACKER_SEED];
    tracker_pda = web3.PublicKey.findProgramAddressSync(tracker_seeds, program.programId);

    treasury_seeds = [TREASURY_SEED];
    treasury_pda = web3.PublicKey.findProgramAddressSync(tracker_seeds, program.programId);


    ////////////////// FEED PROTOCOL CONFIGURATION //////////////////
    const current_feeds_account = PublicKey.findProgramAddressSync(
      [Buffer.from("c"), Buffer.from([1])],
      rngProgram
    );
  
    const currentFeedsAccountInfo = await connection.getAccountInfo(
      current_feeds_account[0]
    );
    const currentFeedsAccountData = deserialize(
      CurrentFeedSchema,
      CurrentFeed,
      currentFeedsAccountInfo?.data!
    );
  
    const feedAccount1 = new PublicKey(
      bs58.encode(currentFeedsAccountData.account1).toString()
    );
    const feedAccount2 = new PublicKey(
      bs58.encode(currentFeedsAccountData.account2).toString()
    );
    const feedAccount3 = new PublicKey(
      bs58.encode(currentFeedsAccountData.account3).toString()
    );
  
    const fallbackAccount = new PublicKey(
      bs58.encode(currentFeedsAccountData.fallback_account).toString()
    );
  
    const tempKeypair = anchor.web3.Keypair.generate();
  })

    // const tx = await program.methods
    //   .participate(new BN(1))
    //   .accounts({
    //     signer: signer.publicKey,
    //     feedAccount1: feedAccount1,
    //     feedAccount2: feedAccount2,
    //     feedAccount3: feedAccount3,
    //     fallbackAccount: fallbackAccount,
    //     currentFeedsAccount: current_feeds_account[0],
    //     temp: tempKeypair.publicKey,
    //     rngProgram: rngProgram,
    //   })
    //   .signers([player, tempKeypair])
    //   .rpc();
  
    // console.log('Transaction signature:', tx);
  

  it("Is Insitialized", async() => {
    console.log(
      `\x1b[32mRaffle has been initialized\x1b[37m
      raffleInfo pubkey: \x1b[32m${raffle_pda}\x1b[37m
      tracker public key: \x1b[32m${tracker_pda}\x1b[37m
      treasury public key: \x1b[32m${treasury_pda}\x1b[37m
      Signer address: \x1b[32m${auth}\x1b[37m\n`
    );

    // const init_raffle_tx = await program.methods.initializeRaffle().accounts([
    //   { name: "tracker", pda: tracker_pda, signer: false, writable: true },
    //   { name: "raffleInfo", pda: raffle_pda, signer: false, writable: true },
    //   { name: "treasury", pda: treasury_pda, signer: false, writable: true},
    //   { name: "auth", pda: auth, signer: true, writable: true },
    //   { name: "systemProgram", pda: web3.SystemProgram.programId, signer: false, writable: false },
    // ])
    // .signers([])
    // .rpc()

    // await confirmTransaction(provider.connection, init_raffle_tx);
  })

  
  // it("user Participated!", async() => {

  // });


  // it ("multi user participated!", async() => {

  // })

  // it("winner selected!", async() => {

  // });

  // it("reward claimed!", async() => {

  // });
})