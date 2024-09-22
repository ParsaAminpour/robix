use anchor_lang::{AnchorSerialize, AnchorDeserialize};
use anchor_lang::prelude::*;

use crate::constants::{
    AMOUNT_OF_WINNER, TICKET_PRICE, TICKET_NUMBER_LOWER_RANGE, TICKET_NUMBER_UPPER_RANGE, AMOUNT_OF_WINNER};

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
pub struct RaffleInfo {
    pub raffle_id: u64,
    pub ticket_price: u64,
    pub total_ticket_sold: u32, // it's const for each round
    pub start_time: u64,
    pub ticket_number_bound: (u64, u64), // (upper, lower)
    pub creator: Pubkey,
    pub tickets: Vec<Ticket>,
    pub winner: Option<Vec<Ticket>>,
    pub is_closed: bool,
    pub active: bool,
    pub treasury_bump: u8,
}

#[derive(AnchorSerialize ,AnchorDeserialize, Clone, Default)]
pub struct Ticket {
    owner: Pubkey,
    ticket_number: u64
}

/// NOTE: functionalities realted to the modifying the RaffleInfo struct will define in this implementation.
impl RaffleInfo {
    pub fn get_space(ticket_count: usize) -> usize {
        8 + (
            8 + 8 + 4 + 8 + (8 + 8) + 32 + (4 + Ticket::INIT_SPACE * ticket_count) + (1 + Ticket::INIT_SPACE * (AMOUNT_OF_WINNER as usize)) + 1 + 1 + 1
        )
    }

    pub fn initialize(&mut self, _ticket_price: u64, _creator: Pubkey, t_bump: u8) {
        self.raffle_id = 0; // init raffle id
        // @audit-info the price should be constant.
        self.ticket_price = TICKET_PRICE;
        self.total_ticket_sold = 0;
        self.creator = _creator;
        self.start_time = Clock::get().unwrap().unix_timestamp as u64;
        // @audit-info this bounds should define dynamicly.
        self.ticket_number_bound = (TICKET_NUMBER_LOWER_RANGE, TICKET_NUMBER_UPPER_RANGE); // assume that we have 100 participants.
        self.tickets = Vec::new();
        self.is_closed = false;
        self.active = true;
        self.treasury_bump = t_bump;

        msg!("A Raffle initialized with {} name", self.raffle_id);
    }

    // @audit-info the ticket number should be verified and unique. (it's asumption here)
    pub fn buy_ticket(&mut self, buyer: Pubkey, verified_ticket_number: u64) -> anchor_lang::Result<()> {
        if !self.active {
            return err!(RaffleErrorCode::RaffleAlreadyDeactivated);
        }
        self.tickets.push(Ticket {
            owner: buyer,
            ticket_number: verified_ticket_number
        });
        self.total_ticket_sold.checked_add(1).unwrap();
        Ok(())
    }

    // using defined raffle winner algorithm.
    pub fn select_winner(&mut self, fprng: u64) -> anchor_lang::Result<()>{
        if self.is_closed {
            return err!(RaffleErrorCode::RaffleAlreadyClosed);
        }
        let winning_number = fprng % (self.ticket_number_bound.0 - self.ticket_number_bound.1 + 1) + self.ticket_number_bound.1;
        
        let mut updated_tickets: Vec<Ticket> = self.tickets.clone().into_iter()
            .map(|mut ticket| {
                ticket.ticket_number -= winning_number;
                ticket
        }).collect();
        updated_tickets.sort_by_key(|t| t.ticket_number);

        let winners: Vec<Ticket> = self.tickets[..(AMOUNT_OF_WINNER as usize)].to_vec();
        /* First three tickets that are near to the rnd point are the winners and the
        * prize will be distributed among these guys,
        * 1st person: 50%  |  2nd person: 30%  |  3th person: 20% */
        self.winner = Some(winners);
        self.active = false;
        Ok(())

    }   
    // To avoid same ticket number result.
    pub fn verify_ticket_number_uniqueness(&self, _t_number_to_examine: u64) -> bool {
        // verifying ticket using binary search.
        self.tickets.binary_search_by(|t| t.ticket_number.cmp(&_t_number_to_examine)).is_ok()
    }
}

impl Space for Ticket {
    const INIT_SPACE: usize = 8 + 32 + 8;
}


#[error_code]
pub enum RaffleErrorCode {
    RaffleAlreadyClosed,
    RaffleAlreadyDeactivated,
    WinnerAlreadySelected
}