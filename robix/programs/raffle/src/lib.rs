#![allow(warnings)]
use anchor_lang::prelude::*;
// use anchor_lang::solana_program::native_token::LAMPORTS_PER_SOL;
use anchor_lang::{AnchorDeserialize, AnchorSerialize};
use instructions::*;
use state::RaffleInfo;

pub mod raffle_error;
pub mod state;
pub mod instructions;
pub mod id;
pub mod constants;
pub use id::ID;

#[program]
pub mod raffle {
    use super::*;

    /// Initializes a new raffle for the users.
    /// # Arguments
    /// - `ctx`: Context for the instruction.
    /// * NOTE: Adding initial lamports to paying the Switchboard randomness fees.
    // @audit-info an EndRaffle should be an initialization either.
    pub fn initialize_raffle(ctx: Context<InitializeRaffle>) -> Result<()> {
        instructions::initialize_raffle(ctx).unwrap();
        Ok(())
    }

    /// Participate at the raffle needs to calling this function.
    /// # Arguments:
    /// - `ctx` Context for participating to the raffle.
    /// - `active_raffle` is the raffle index that is currently active (it'll check based on the game tracker actice raffle argument)
    pub fn participate(ctx: Context<Participate>,active_raffle: u64) -> anchor_lang::Result<()> {
        instructions::participate(ctx, active_raffle).unwrap();
        Ok(())
    }

    /// Claim reward to the winner from Treasury.
    /// # Arguments:
    /// - `ctx` Context for participating to the raffle.
    /// * NOTE: This instruction is provided for treasury to send its funds to the winner if he's already selected.
    // pub fn claim_the_reward(
    //     ctx: Context<ClaimReward>,
    //     _raffle_name: String,
    //     _creator: Pubkey,
    // ) -> Result<()> {
    //     require_eq!(ctx.accounts.raffle_info.winner.unwrap().key(), ctx.accounts.winner.key(), raffle_error::ErrorCode::CallerIsNotWinner);
    //     require!(ctx.accounts.raffle_info.winner.is_some(), raffle_error::ErrorCode::WinnerIsNotSelectedYet);
    //     let dest_balance = ctx.accounts.raffle_info.raffle_pool;
    //     let bump = &[ctx.bumps.treasury];
    //     let seeds: &[&[u8]] = &[b"treasury".as_ref(), bump];
    //     let signer_seeds = &[seeds];
    //     ctx.accounts.transfer(
    //         ctx.accounts.system_program.to_account_info(), 
    //         ctx.accounts.treasury.to_account_info(), 
    //         ctx.accounts.winner.to_account_info(), 
    //         ctx.accounts.raffle_info.raffle_pool, 
    //         true, 
    //         signer_seeds
    //     )?;
    //     ctx.accounts.raffle_info.is_closed = true;
    //     msg!("Reward: {} has been calimed by {}", dest_balance, ctx.accounts.raffle_info.winner.unwrap().key().to_string());
    //     Ok(())
    // }

    pub fn close(ctx: Context<Close>, _raffle_name: String, _creator: Pubkey) -> Result<()> {
        require!(
            ctx.accounts.account_to_close.is_closed,
            raffle_error::ErrorCode::RaffleIsNotClosed
        );
        Ok(())
    }
}

#[derive(Accounts)]
#[instruction(raffle_name: String, creator: Pubkey)]
pub struct Close<'info> {
    #[account(
        mut,
        has_one = creator,
        seeds = ["raffle".as_bytes(), raffle_name.as_bytes(), creator.key().as_ref()], 
        close = creator,
        bump
    )]
    pub account_to_close: Account<'info, RaffleInfo>,

    #[account(mut)]
    pub creator: Signer<'info>,
}
