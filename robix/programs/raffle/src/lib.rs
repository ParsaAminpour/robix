use anchor_lang::prelude::*;
use anchor_lang::solana_program::native_token::LAMPORTS_PER_SOL;
use anchor_lang::{AnchorDeserialize, AnchorSerialize};
use solana_program::keccak::hash as keccak_hash;

pub mod raffle_error;

declare_id!("7iKkN561Q2C5w9ooaf5U7LHnVtH3VyyErWiiUr1TJcRk");

const INIT_TREASURY_FUND: f32 = 0.01 * (LAMPORTS_PER_SOL as f32);

#[program]
pub mod raffle {
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


    /// Participate at the raffle needs to calling this function.
    /// # Arguments:
    /// - `ctx` Context for participating to the raffle.
    pub fn participate(ctx: Context<Participate>, _raffle_name: String, _creator: Pubkey) -> anchor_lang::Result<()> {
        require!(ctx.accounts.raffle_info.winner.is_none() || !ctx.accounts.raffle_info.is_closed, raffle_error::ErrorCode::WinnerAlreadySelected);
        require_gte!(ctx.accounts.sender.lamports(), ctx.accounts.raffle_info.ticket_price, raffle_error::ErrorCode::NotSufficientBalance);
        require_gte!(ctx.accounts.raffle_info.max_tickets, ctx.accounts.raffle_info.total_ticket_sold, raffle_error::ErrorCode::TicketAmountThreshold);
        require!(ctx.accounts.raffle_info.start_time < ctx.accounts.raffle_info.end_time, raffle_error::ErrorCode::RaffleExpired);
        
        ctx.accounts.transfer(
            ctx.accounts.system_program.to_account_info(), 
            ctx.accounts.sender.to_account_info(), 
            ctx.accounts.treasury.to_account_info(), 
            ctx.accounts.raffle_info.ticket_price, 
            false, 
            &[&[&[1_u8]]]
        )?;
        
        let associated_raffle = &mut ctx.accounts.raffle_info;
        associated_raffle.total_ticket_sold = associated_raffle.total_ticket_sold.checked_add(1).unwrap();
            
        associated_raffle.raffle_pool += associated_raffle.ticket_price;

        ctx.accounts
            .participant_list
            .participants
            .push(ctx.accounts.sender.key());

        msg!("{} Joined to the list of participants", ctx.accounts.sender.key().to_string());
        Ok(())
    }


    /// Selecting the winner among all participants in the raffle.
    /// # Arguments:
    /// - `ctx` Context for participating to the raffle.
    /// * NOTE: Only the owner of the raffle could call this function.
    pub fn winner_selection(ctx: Context<WinnerSelection>, _raffle_name: String) -> Result<()> {
        require!(ctx.accounts.raffle_info.winner.is_none(), raffle_error::ErrorCode::WinnerAlreadySelected);

        //// Generating randomness - the Feed protocol will add here ////
        let clock: Clock = Clock::get().unwrap();
        let participants_count = ctx.accounts.participant_list.participants.len();
        // @audit-info the feed protocol will be replaced in this section after they resolve their protocol bug.
        let revealed_random_value = ctx
            .accounts
            .rpng(&clock.unix_timestamp, participants_count as u64);

        // @audit-solved type parsing unsupported value to mod between u32&u8
        let winner = ctx.accounts.participant_list.participants[revealed_random_value as usize];
        let winner_selected = ctx.accounts.raffle_info.winner.insert(winner);
        msg!(
            "The Winner selected public key: {}, The relevant raffle name: {}",
            winner_selected.to_string(),
            ctx.accounts.raffle_info.raffle_name
        );
        Ok(())
    }

    /// Claim reward to the winner from Treasury.
    /// # Arguments:
    /// - `ctx` Context for participating to the raffle.
    /// * NOTE: This instruction is provided for treasury to send its funds to the winner if he's already selected.
    pub fn claim_the_reward(
        ctx: Context<ClaimReward>,
        _raffle_name: String,
        _creator: Pubkey,
    ) -> Result<()> {
        require_eq!(ctx.accounts.raffle_info.winner.unwrap().key(), ctx.accounts.winner.key(), raffle_error::ErrorCode::CallerIsNotWinner);
        require!(ctx.accounts.raffle_info.winner.is_some(), raffle_error::ErrorCode::WinnerIsNotSelectedYet);

        let dest_balance = ctx.accounts.raffle_info.raffle_pool;

        let bump = &[ctx.bumps.treasury];
        let seeds: &[&[u8]] = &[b"treasury".as_ref(), bump];
        let signer_seeds = &[seeds];

        ctx.accounts.transfer(
            ctx.accounts.system_program.to_account_info(), 
            ctx.accounts.treasury.to_account_info(), 
            ctx.accounts.winner.to_account_info(), 
            ctx.accounts.raffle_info.raffle_pool, 
            true, 
            signer_seeds
        )?;
        ctx.accounts.raffle_info.is_closed = true;

        msg!("Reward: {} has been calimed by {}", dest_balance, ctx.accounts.raffle_info.winner.unwrap().key().to_string());
        Ok(())
    }

    pub fn close(ctx: Context<Close>, _raffle_name: String, _creator: Pubkey) -> Result<()> {
        require!(
            ctx.accounts.account_to_close.is_closed,
            raffle_error::ErrorCode::RaffleIsNotClosed
        );
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
#[instruction(raffle_name: String)]
pub struct InitializeRaffle<'info> {
    #[account(
        init,
        payer = signer,
        seeds = ["raffle".as_ref(), raffle_name.as_ref(), signer.key().as_ref()],
        space = RaffleInfo::INIT_SPACE,
        bump,
    )]
    pub raffle_info: Account<'info, RaffleInfo>,

    #[account(
        init,
        payer = signer,
        space = ParticipantList::INIT_SPACE,
        seeds = ["participant_list".as_ref(), raffle_name.as_ref()],
        bump
    )]
    pub participant_list: Account<'info, ParticipantList>,

    #[account(
        mut,
        seeds = [b"treasury".as_ref()],
        bump
    )]
    /// CHECK: This is okay - it's a PDA to store SOL and doesn't need a data layout
    pub treasury: SystemAccount<'info>,

    #[account(mut)]
    pub signer: Signer<'info>, // aka creator.
    pub system_program: Program<'info, System>,
}

impl<'a> Transfer<'a> for InitializeRaffle<'a> {}


/// Instruction to participate in the raffle.
/// # Arguments
/// - `ctx`: Context of the instruction, providing accounts and program state.
/// - `raffle_info`: The information provided for the raffle.
/// - `sender`: The user participating in the raffle.
/// - `treasury`: The program's account to receive the ticket cost paid by the user like a treasury.
/// - `system_program`: The system program account.
/// *  NOTE: the treasury pubkey of this account should be as same as the ClaimReward's treasury pubkey.
#[derive(Accounts)]
#[instruction(raffle_name: String, creator: Pubkey)]
pub struct Participate<'info> {
    #[account(
        mut,
        seeds = ["raffle".as_ref(), raffle_name.as_ref(), creator.key().as_ref()],
        bump
    )]
    pub raffle_info: Account<'info, RaffleInfo>,

    #[account(mut, seeds=["participant_list".as_ref(), raffle_name.as_ref()], bump)]
    pub participant_list: Account<'info, ParticipantList>,

    #[account(mut)]
    pub sender: Signer<'info>,

    /// CHECK: This is not dangerous because we are transferring SOL to the program's account
    #[account(mut, seeds = [b"treasury".as_ref()], bump)]
    pub treasury: SystemAccount<'info>,
    pub system_program: Program<'info, System>,
}

impl<'a> Transfer<'a> for Participate<'a> {}


/// Instruction to participate in the raffle.
/// # Arguments
/// - `ctx`: Context of the instruction, providing accounts and program state.
/// - `raffle_info`: The information provided for the raffle.
/// - `sender`: The user participating in the raffle.
/// - `receiver`: The program's account to receive the ticket cost paid by the user.
/// - `system_program`: The system program account.
/// *  NOTE: the treasury pubkey of this account should be as same as the ClaimReward's treasury pubkey.
#[derive(Accounts)]
#[instruction(raffle_name: String)]
pub struct WinnerSelection<'info> {
    #[account(
        mut,
        has_one = creator,
        seeds = ["raffle".as_ref(), raffle_name.as_ref(), creator.key().as_ref()],
        bump
    )]
    pub raffle_info: Account<'info, RaffleInfo>,
    // @audit-info This `randomness_account_data` will be handled after the Feed protocol VRF is fixed.
    /// CHECK: The account's data is validated manually within the handler.
    // pub randomness_account_data: AccountInfo<'info>,

    #[account(mut, seeds=["participant_list".as_ref(), raffle_name.as_ref()], bump)]
    pub participant_list: Account<'info, ParticipantList>,

    #[account(mut)]
    pub creator: Signer<'info>,

    pub system_program: Program<'info, System>,
}

// TODO: will replace by Feed protocol VRF.
impl<'a> GenerateRandomness<'a> for WinnerSelection<'a> {
    fn rpng(&mut self, seed: &i64, rng: u64) -> u64 {
        let hash = keccak_hash(&seed.to_le_bytes()).to_bytes();
        let random_value = u64::from_le_bytes([
            hash[0], hash[1], hash[2], hash[3], hash[4], hash[5], hash[6], hash[7],
        ]);
        random_value % rng
    }
}


/// Instruction to send the colllected funds to the selected winner by `treasury`.
/// # Arguments
/// - `ctx`: Context of the instruction, providing accounts and program state.
/// - `raffle_info`: The information provided for the raffle.
/// - `treasury`: The treasury as the signer to send the reward to the selected winner.
/// - `winner`: The selected winner that the treasury's fund will send to him.
/// - `system_program`: The system program account.
/// *  NOTE: the treasury pubkey of this account should be as same as the WinnerSelection's treasury pubkey.
#[derive(Accounts)]
#[instruction(raffle_name: String, creator: Pubkey)]
pub struct ClaimReward<'info> {
    #[account(
        mut,
        seeds = ["raffle".as_bytes(), raffle_name.as_bytes(), creator.key().as_ref()],
        bump,
    )]
    pub raffle_info: Account<'info, RaffleInfo>,

    /// CHECK: This is not dangerous because we are transferring SOL to the program's account
    #[account(
        mut, 
        seeds = [b"treasury".as_ref()],
        bump
    )]
    pub treasury: SystemAccount<'info>,

    #[account(mut)]
    pub winner: Signer<'info>,
    pub system_program: Program<'info, System>,
}

impl<'a> Transfer<'a> for ClaimReward<'a> {}



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



/// Struct representing the raffle state.
/// # Fields
/// - `raffle_name`: The name of the raffle when it's initialized.
/// - `ticket_price`: The fixed price users must pay for each ticket.
/// - `raffle_pool`: The SOL treasury collected for the raffle by participants.
/// - `max_tickets`: The maximum number of tickets available in this raffle.
/// - `total_tickets_sold`: The number of tickets that have been sold.
/// - `start_time`: The raffle start time, set to the account initialization time.
/// - `end_time`: The deadline for participating in the raffle.
/// - `creator`: The creator of the raffle.
/// - `winner`: The winner's address, selected at the end of the raffle.
/// - `is_closed`: Indicates if the raffle is closed; this field will be true when the raffle is closed.
#[account]
#[derive(InitSpace)]
pub struct RaffleInfo {
    #[max_len(20)]
    pub raffle_name: String,
    pub ticket_price: u64,
    pub raffle_pool: u64,
    pub max_tickets: u32,
    pub total_ticket_sold: u32,
    pub start_time: u64,
    pub end_time: u64,
    pub creator: Pubkey,
    pub winner: Option<Pubkey>,
    pub is_closed: bool,
    pub treasury_bump: u8,
}

#[account]
#[derive(InitSpace)]
pub struct ParticipantList {
    #[max_len(20)]
    pub participants: Vec<Pubkey>,
    pub bump: u8,
}

pub struct RandomNumber {
    pub random_number: u64,
}
