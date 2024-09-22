use anchor_lang::prelude::*;
use crate::state::raffle::{RaffleInfo};


#[program]
pub mod raffle_instructions {
    use super::*;

    /// Initializes a new raffle for the users.
    ///
    /// # Arguments
    /// - `ctx`: Context for the instruction.
    /// - `raffle_name`: The specific name for the new raffle, also used as a discriminator.
    /// - `ticket_price`: The fixed price users must pay for each ticket.
    /// - `max_tickets`: The maximum number of tickets available in this raffle.
    /// - `end_time`: The deadline for the raffle.
    /// * NOTE: Adding initial lamports to paying the Switchboard randomness fees.
    // @audit-info an EndRaffle should be an initialization either.
    pub fn initialize_raffle(ctx: Context<InitializeRaffle>, raffle_name: String, ticket_price: u64, max_tickets: u32, end_time: u64) -> Result<()> {
        let time = Clock::get().unwrap();
        require!(raffle_name.len() > 1, raffle_error::ErrorCode::EmptyStringNotAllowed);
        require!(end_time > time.unix_timestamp as u64, raffle_error::ErrorCode::InvalidEndTime);
        require!(ticket_price * (max_tickets as u64) != 0, raffle_error::ErrorCode::InvalidAmount);

        ctx.accounts.raffle_info.raffle_name = raffle_name;
        ctx.accounts.raffle_info.ticket_price = ticket_price;
        ctx.accounts.raffle_info.raffle_pool = 0;
        ctx.accounts.raffle_info.max_tickets = max_tickets;
        ctx.accounts.raffle_info.total_ticket_sold = 0;
        ctx.accounts.raffle_info.start_time = time.unix_timestamp as u64;
        ctx.accounts.raffle_info.end_time = end_time;
        ctx.accounts.raffle_info.creator = ctx.accounts.signer.key();
        ctx.accounts.raffle_info.is_closed = false;
        ctx.accounts.raffle_info.active = true;
        
        ctx.accounts.transfer(
            ctx.accounts.system_program.to_account_info(), 
            ctx.accounts.signer.to_account_info(), 
            ctx.accounts.treasury.to_account_info(), 
            INIT_TREASURY_FUND as u64, 
            false, 
            &[&[&[1_u8]]]
        )?;
        msg!("Raffle {} has been initialized!", ctx.accounts.raffle_info.raffle_name);
        Ok(())
    }

}

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



/// Instruction to initialize a raffle.
/// # Arguments
/// - `ctx`: Context of the instruction, providing accounts and program state.
/// - `raffle_info`: The information provided for the raffle.'
/// - `participant_list`: An account for storing the participants.
/// - `signer`: The signer of the account.
/// - `system_program`: The system program account.
#[derive(Accounts)]
#[instruction(raffle_id: u64)]
pub struct InitializeRaffle<'info> {
    #[account(
        init,
        payer = signer,
        // @audit-info adding raffle number to the seeds for generate.
        seeds = ["raffle".as_ref(), raffle_id.to_strig().as_bytes(), signer.key().as_ref()],
        space = RaffleInfo::INIT_SPACE,
        bump,
    )]
    pub raffle_info: Account<'info, RaffleInfo>,

    #[account(
        mut,
        seeds = [b"treasury".as_ref()],
        bump
    )]

    // @audit-info this should be the fee vault.
    /// CHECK: This is okay - it's a PDA to store SOL and doesn't need a data layout
    pub treasury: SystemAccount<'info>,

    // @audit should determine a constant address as the signer, not an arbitary signer.
    #[account(mut)]
    pub signer: Signer<'info>, // aka creator.
    #[account(address = anchor_lang::system_program::ID)]
    pub system_program: Program<'info, System>,
}

impl<'a> Transfer<'a> for InitializeRaffle<'a> {}
