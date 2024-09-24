use std::str::FromStr;
use std::sync::atomic;

use anchor_lang::prelude::*;
use crate::state::raffle_account::{RaffleInfo, RaffleFeeVault};
use crate::constants::{
    RAFFLE_SEED, RAFFLE_OWNER, TREASURY_SEED, TRACKER_SEED, FEE_VAULT, FPRNG_PROGRAM_ADDRESS};
use crate::raffle_error;
use crate::state::game_tracker::GameTracker;

use anchor_lang::solana_program::{instruction::Instruction, program::{get_return_data, invoke}};

/// Instruction to initialize a raffle.
/// # Arguments
/// - `ctx`: Context of the instruction, providing accounts and program state.
/// - `raffle_info`: The information provided for the raffle.'
/// - `participant_list`: An account for storing the participants.
/// - `signer`: The signer of the account.
/// - `system_program`: The system program account.
#[derive(Accounts)]
pub struct InitializeRaffle<'info> {
    #[account(
        init,
        payer = auth,
        seeds = [TRACKER_SEED.as_ref()],
        space = GameTracker::get_space(0),
        bump
    )]
    pub tracker: Account<'info, GameTracker>,

    #[account(
        init,
        payer = auth,
        seeds = [RAFFLE_SEED.as_ref(), &(1_u64).to_le_bytes(), auth.key().as_ref()],
        space = RaffleInfo::get_space(0),
        bump,
    )]
    pub raffle_info: Account<'info, RaffleInfo>,

    #[account(
        init,
        payer = auth,
        seeds = [TREASURY_SEED.as_ref()],
        space = RaffleFeeVault::INIT_SPACE,
        bump
    )]
    pub treasury: Account<'info, RaffleFeeVault>,

    #[account(mut, address = Pubkey::from_str(RAFFLE_OWNER).unwrap() @ raffle_error::ErrorCode::SignerIsNotValid)]
    pub auth: Signer<'info>, // aka creator.

    #[account(address = anchor_lang::system_program::ID)]
    pub system_program: Program<'info, System>,
}

impl<'a> Transfer<'a> for InitializeRaffle<'a> {}

/// Initializes a new raffle for the users.
///
/// # Arguments
/// - `ctx`: Context for the instruction.
/// * NOTE: Adding initial lamports to paying the Switchboard randomness fees.
pub fn initialize_raffle(ctx: Context<InitializeRaffle>) -> Result<()> {
    let raffle = &mut ctx.accounts.raffle_info;
    let game_tracker = &mut ctx.accounts.tracker;
    let auth = &mut ctx.accounts.auth;

    raffle.initialize(
        auth.key(), ctx.bumps.raffle_info
    );
    game_tracker.initialize(ctx.bumps.tracker, auth.key());

    msg!("Raffle initialized at {}", raffle.key());
    msg!("Game Tracked initialized at {}", game_tracker.key());
    Ok(())
}



#[derive(Accounts)]
pub struct Participate<'info> {
    #[account(
        mut,
        address = Pubkey::from_str(FEE_VAULT).unwrap() @ raffle_error::ErrorCode::NotValidTreasuryAddress
    )]
    pub treasury: Account<'info, RaffleFeeVault>,

    #[account(
        mut,
        seeds = [TRACKER_SEED.as_ref()],
        bump = game_tracker.bump
    )]
    pub game_tracker: Account<'info, GameTracker>,

    #[account(
        mut,
        seeds = [RAFFLE_SEED.as_ref(), &(game_tracker.active_raffle).to_le_bytes(), game_tracker.active_raffle_owner.as_ref()],
        bump = raffle_info.raffle_bump,
        realloc = RaffleInfo::get_space(raffle_info.tickets.len() + 1),
        realloc::payer = participant,
        realloc::zero = false  
    )]
    pub raffle_info: Account<'info, RaffleInfo>,

    #[account(mut)]
    pub participant: Signer<'info>,

    /// CHECK: Feed account
    pub feed_account_1: AccountInfo<'info>,
    /// CHECK: Feed account
    pub feed_account_2: AccountInfo<'info>,
    /// CHECK: Feed account
    pub feed_account_3: AccountInfo<'info>,
    /// CHECK: Feed fallback accoubt
    pub fallback_account: AccountInfo<'info>,
    #[account(mut)]
    /// CHECK: Current Feed Account
    pub current_feeds_account: AccountInfo<'info>,
    #[account(mut)]
    /// CHECK: idk what the hell it this
    pub temp: Signer<'info>,
    #[account(address = Pubkey::from_str(FPRNG_PROGRAM_ADDRESS).unwrap() @ raffle_error::ErrorCode::InvalidFeedRNGAddress)]
    pub rng_program: AccountInfo<'info>,

    #[account(address = anchor_lang::system_program::ID)]
    pub system_program: Program<'info, System>
}

pub fn participate(ctx: Context<Participate>) -> anchor_lang::Result<()> {
    // get vrf from Feed protocol using CPI
    let rng_program: &Pubkey = ctx.accounts.rng_program.key;
   
    //Creating instruction for CPI to RNG_PROGRAM of Feed protocol
    let instruction: Instruction = Instruction {
        program_id: *rng_program,
        accounts: vec![
            ctx.accounts.participant.to_account_metas(Some(true))[0].clone(),
            ctx.accounts.feed_account_1.to_account_metas(Some(false))[0].clone(),
            ctx.accounts.feed_account_2.to_account_metas(Some(false))[0].clone(),
            ctx.accounts.feed_account_3.to_account_metas(Some(false))[0].clone(),
            ctx.accounts.fallback_account.to_account_metas(Some(false))[0].clone(),
            ctx.accounts.current_feeds_account.to_account_metas(Some(false))[0].clone(),
            ctx.accounts.temp.to_account_metas(Some(true))[0].clone(),
            ctx.accounts.system_program.to_account_metas(Some(false))[0].clone(),
        ],
        data: vec![0],
    };

    //Creating account infos for CPI to RNG_PROGRAM
    let account_infos: &[AccountInfo; 8] = &[
        ctx.accounts.participant.to_account_info().clone(),
        ctx.accounts.feed_account_1.to_account_info().clone(),
        ctx.accounts.feed_account_2.to_account_info().clone(),
        ctx.accounts.feed_account_3.to_account_info().clone(),
        ctx.accounts.fallback_account.to_account_info().clone(),
        ctx.accounts.current_feeds_account.to_account_info().clone(),
        ctx.accounts.temp.to_account_info().clone(),
        ctx.accounts.system_program.to_account_info().clone(),
    ];

    invoke(&instruction, account_infos)?;
    let return_data: (Pubkey, Vec<u8>) = get_return_data().unwrap();

    // let rnd_number: u64;
    // if &return_data.0 == rng_program {
    //     rnd_number = return_data.1;
    // }

    // verify that vrf number based on tickets number

    // other obvs instructions.
    Ok(())
}

impl<'a> Transfer<'a> for Participate<'a> {}



trait Transfer<'a> {
    fn transfer(&mut self, system_program: AccountInfo<'a>, sender: AccountInfo<'a>, receiver: AccountInfo<'a>, amount: u64, from_vault: bool, _signer_seeds: &[&[&[u8]]]) -> Result<()> {
        require_gte!(sender.get_lamports(), amount, raffle_error::ErrorCode::InsufficientBalance);
        require_neq!(sender.key(), receiver.key(), raffle_error::ErrorCode::SameDestinationAddressNotAllowed);
        require_neq!(amount, 0, raffle_error::ErrorCode::InvalidAmount);

        if from_vault {
            anchor_lang::system_program::transfer(
                CpiContext::new(
                    system_program.to_account_info(),
                    anchor_lang::system_program::Transfer {
                        from: sender.to_account_info(),
                        to: receiver.to_account_info()
                    }
                ).with_signer(_signer_seeds), 
                amount
            ).unwrap();
            msg!("SOL transfered from {} to {}", sender.key(), receiver.key());

        } else {
            anchor_lang::system_program::transfer(
                CpiContext::new(
                    system_program.to_account_info(),
                    anchor_lang::system_program::Transfer {
                        from: sender.to_account_info(),
                        to: receiver.to_account_info()
                    }),
                    amount
            ).unwrap();
            msg!("SOL transfered from treasury to {}", receiver.key());
        }
        Ok(())
    }
}
trait GenerateRandomness<'a> {
    fn rpng(&mut self, seed: &i64, rng: u64) -> u64;
}

