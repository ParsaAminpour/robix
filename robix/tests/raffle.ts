import * as anchor from "@coral-xyz/anchor";
import { Program } from "@coral-xyz/anchor";
import { Keypair, PublicKey } from "@solana/web3.js";
import { Raffle } from "../target/types/raffle";
import { BN } from "bn.js";
import web3 from "@solana/web3.js";
import { assert } from "chai";
import { Signature } from "ethers";
import { walkUpBindingElementsAndPatterns } from "typescript";

const INIT_TREASURY_FUND = 0.01 * web3.LAMPORTS_PER_SOL;


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


describe("Testing Raffle instructions", () => {
  const provider = anchor.AnchorProvider.env();
  anchor.setProvider(provider);
  
  const program = anchor.workspace.Raffle as Program<Raffle>;

  const raffle_name: string = "raffle2";
  const ticket_price = new BN(0.1 * web3.LAMPORTS_PER_SOL);
  const max_tickets = 100;
  const end_time = new BN(Math.floor(Date.now() / 1000) + 86400);
  
  let signer: PublicKey, participant: Keypair, treasury: Keypair;  let raffle_seeds, participant_list_seeds, treasury_seed;
  let raffle_info_pda, participant_list_pda, treasury_pda;;


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
    signer = provider.wallet.publicKey;
    participant = web3.Keypair.generate();

    const signer_init_balance = await provider.connection.getBalance(signer);
    const participant_init_balance = await provider.connection.getBalance(participant.publicKey);
    console.log(`Signer Balance: ${signer_init_balance}`);
    console.log(`Participant Balance: ${participant_init_balance}`);

    await fundWallet(signer, participant.publicKey, 0.1 * web3.LAMPORTS_PER_SOL);
    console.log(`participant balance after trasnfer: ${await provider.connection.getBalance(participant.publicKey)}`);

    raffle_seeds = [
      Buffer.from("raffle"), 
      Buffer.from(raffle_name),  // raffle_name should match what is passed to the instruction
      signer.toBuffer()
    ];
    [raffle_info_pda] = web3.PublicKey.findProgramAddressSync(
      raffle_seeds,
      program.programId
    )
  
    participant_list_seeds = [Buffer.from("participant_list"), Buffer.from(raffle_name)];
    [participant_list_pda] = web3.PublicKey.findProgramAddressSync(
      participant_list_seeds,
      program.programId
    );

    treasury_seed = [Buffer.from("treasury")];
    [treasury_pda] = web3.PublicKey.findProgramAddressSync(
      treasury_seed,
      program.programId
    );
  });


  it("Is Insitialized", async() => {
    try {
      const tx = await program.methods.initializeRaffle(
        raffle_name,
        ticket_price,
        max_tickets,
        end_time
      ).accounts([
        { name: "raffleInfo", pda: raffle_info_pda, signer: false, writable: true },
        { name: "participantList", pda: participant_list_pda, signer: false, writable: true },
        { name: "treasury", pda: treasury_pda, signer: false, writable: true},
        { name: "signer", pda: signer, signer: true, writable: true },
        { name: "systemProgram", pda: web3.SystemProgram.programId, signer: false, writable: false },
      ])
      .signers([])
      .rpc();

      await confirmTransaction(provider.connection, tx);
    } catch (err) {
      assert.fail(`Error in transaction: ${err.message}`);
    }

    console.log(
      `\x1b[32mRaffle has been initialized\x1b[37m
      raffleInfo pubkey: \x1b[32m${raffle_info_pda}\x1b[37m
      participants list pubkey: \x1b[32m${participant_list_pda}\x1b[37m
      treasury public key: \x1b[32m${treasury_pda}\x1b[37m
      Signer address: \x1b[32m${signer}\x1b[37m\n`
    );
  
    let raffle_info_fetched: any = await program.account.raffleInfo.fetch(raffle_info_pda.toBase58());    
    const treasury_balance = await provider.connection.getBalance(treasury_pda);
    console.log(`Treasury balance after initializing: ${treasury_balance}`);

    assert.equal(raffle_info_fetched.raffleName, raffle_name);
    assert.equal(raffle_info_fetched.ticketPrice.toString(), ticket_price.toString());
    assert.equal(raffle_info_fetched.maxTickets, max_tickets);
    assert.equal(raffle_info_fetched.creator.toBase58(), signer);
    assert.strictEqual(treasury_balance, INIT_TREASURY_FUND);
  })

  
  it("user Participated!", async() => {
    const t_balance_before_participate = await provider.connection.getBalance(treasury_pda);
    
    const participating_tx1 = await program.methods.participate(raffle_name, signer)
    .accounts([
      { name: "raffleInfo", pda: raffle_info_pda, writable: true },
      { name: "participantList", pda: participant_list_pda, writable: true },
      { name: "sender", pda: signer, signer: true, writable: true },
      { name: "treasury", pda: treasury_pda, signer: false, writable: true },
      { name: "systemProgram", pda: web3.SystemProgram.programId, signer: false, writable: false },
    ])
    .signers([])
    .rpc()
    
    const raffle_fetched = await program.account.raffleInfo.fetch(raffle_info_pda);

    const t_balance_after_participate = await provider.connection.getBalance(treasury_pda);
    const total_ticket_sold_after_participate = raffle_fetched.totalTicketSold;

    const participant_list_fetched: any = await program.account.participantList.fetch(participant_list_pda.toBase58());
    const list_of_participants: PublicKey[] = participant_list_fetched.participants;

    list_of_participants.forEach((addr: PublicKey) => {
      console.log(`address: \x1b[32m${addr.toBase58()}\x1b[37m`);
    });

    assert.equal(t_balance_before_participate + raffle_fetched.ticketPrice.toNumber(), t_balance_after_participate);
    assert.equal(total_ticket_sold_after_participate, 1);
    assert.equal(participant_list_fetched.participants.length, 1);
  });


  it ("multi user participated!", async() => {
    const t_balance_before_mutli_participants = await provider.connection.getBalance(treasury_pda);

    for (let i=0; i<3; i++) {
      await program.methods.participate(raffle_name, signer)
      .accounts([
        { name: "raffleInfo", pda: raffle_info_pda, writable: true },
        { name: "participantList", pda: participant_list_pda, writable: true },
        { name: "sender", pda: signer, signer: true, writable: true },
        { name: "treasury", pda: treasury_pda, signer: false, writable: true },
        { name: "systemProgram", pda: web3.SystemProgram.programId, signer: false, writable: false },
      ])
      .signers([])
      .rpc()
    }

    const participant_list_fetched: any = await program.account.participantList.fetch(participant_list_pda.toBase58());
    const raffle_fetched = await program.account.raffleInfo.fetch(raffle_info_pda);
    const t_balance_after_multi_participants = await provider.connection.getBalance(treasury_pda);

    assert.equal(t_balance_after_multi_participants, t_balance_before_mutli_participants + (3 * raffle_fetched.ticketPrice.toNumber()));
    assert.equal(raffle_fetched.totalTicketSold, 4);
    assert.equal(participant_list_fetched.participants.length, 4);
  })

  it("winner selected!", async() => {
    const t_balance_before_winner_selection = await provider.connection.getBalance(treasury_pda);

    const winner_selection_tx = await program.methods.winnerSelection(raffle_name)
      .accounts([
        { name: "raffleInfo", pda: raffle_info_pda, writable: true },
        { name: "participantList", pda: participant_list_pda, writable: true },
        { name: "creator", pda: signer, signer: true, writable: true },
        { name: "systemProgram", pda: web3.SystemProgram.programId, signer: false, writable: false },
      ])
      .signers([])
      .rpc();
    
    const participant_list_fetched: any = await program.account.participantList.fetch(participant_list_pda.toBase58());
    const raffle_fetched = await program.account.raffleInfo.fetch(raffle_info_pda);
    const t_balance_after_winner_selection = await provider.connection.getBalance(treasury_pda);

    assert.equal(raffle_fetched.winner.toJSON(), signer.toBase58());
  });

  it("reward claimed!", async() => {
    const winner_pda = web3.Keypair.generate();
    const t_balance_before = await provider.connection.getBalance(treasury_pda);
    const raffle_fetched = await program.account.raffleInfo.fetch(raffle_info_pda);
    const raffle_pool_before_claim = raffle_fetched.rafflePool;

    try {
      const claim_reward_tx = await program.methods.claimTheReward(raffle_name, signer)
        .accounts([
          { name: "raffleInfo", pda: raffle_info_pda, writable: true },
          { name: "treasury", pda: treasury_pda, writable: true},
          { name: "signer", pda: signer, signer: true, writable: true },
          { name: "systemProgram", pda: web3.SystemProgram.programId, signer: false, writable: false },
        ])
        .signers([])
        .rpc();

        await confirmTransaction(provider.connection, claim_reward_tx);
    } catch (err) {
      assert.fail(`Error in transaction: ${err.message}`);
    }

      const raffle_fetched_after_clim = await program.account.raffleInfo.fetch(raffle_info_pda);

      const t_balance_after = await provider.connection.getBalance(treasury_pda);
      assert.strictEqual(t_balance_after, t_balance_before - raffle_pool_before_claim.toNumber());
      assert.isTrue(raffle_fetched_after_clim.isClosed);
  });
})