use anchor_lang::{AnchorSerialize, AnchorDeserialize};
use anchor_lang::prelude::*;
use inline_colorization::*;

use crate::constants::{
    AMOUNT_OF_WINNER, TICKET_PRICE, TICKET_NUMBER_LOWER_RANGE, TICKET_NUMBER_UPPER_RANGE};

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
    pub raffle_bump: u8,
}

#[derive(AnchorSerialize ,AnchorDeserialize, Clone, Default, Debug)]
pub struct Ticket {
    owner: Pubkey,
    ticket_number: u64
}


/// NOTE: functionalities realted to the modifying the RaffleInfo struct will define in this implementation.
// @audit-info each method doesn't need raffle args check like activity, these constraints performed in instructions.
impl RaffleInfo {
    pub fn get_space(ticket_count: usize) -> usize {
        8 + (
            8 + 8 + 4 + 8 + (8 + 8) + 32 + (4 + Ticket::INIT_SPACE * ticket_count) + (1 + Ticket::INIT_SPACE * (AMOUNT_OF_WINNER as usize)) + 1 + 1 + 1
        )
    }

    pub fn initialize(&mut self, creator: Pubkey, r_bump: u8) {
        self.raffle_id = 0; // init raffle id
        // @audit-info the price should be constant.
        self.ticket_price = TICKET_PRICE;
        self.total_ticket_sold = 0;
        self.creator = creator;
        self.start_time = Clock::get().unwrap().unix_timestamp as u64;
        // @audit-info this bounds should define dynamicly.
        self.ticket_number_bound = (TICKET_NUMBER_UPPER_RANGE, TICKET_NUMBER_LOWER_RANGE); // assume that we have 100 participants.
        self.tickets = Vec::new();
        self.is_closed = false;
        self.active = true;
        self.raffle_bump = r_bump;

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
        self.total_ticket_sold += 1;
        Ok(())
    }

    // using defined raffle winner algorithm.
    pub fn select_winner(&mut self, fprng: u64) -> anchor_lang::Result<()>{
        if self.is_closed {
            return err!(RaffleErrorCode::RaffleAlreadyClosed);
        }
        // @audit overflow vulnerable.
        let winning_number = fprng % (self.ticket_number_bound.0 - self.ticket_number_bound.1 + 1) + self.ticket_number_bound.1;
        // println!("\n{color_yellow}self.tickets (before): {:?}{color_reset}\n", self.tickets.clone());

        let mut updated_tickets: Vec<Ticket> = self.tickets.clone().into_iter()
            .map(|ticket| {
                let _ = u64::abs_diff(ticket.ticket_number, winning_number);
                ticket
        }).collect();
        updated_tickets.sort_by_key(|t| t.ticket_number);
        updated_tickets.reverse();

        // println!("{color_green}self.tickets (after): {:?}{color_reset}\n", self.tickets.clone());
        // println!("{color_green}updated_tickets (after): {:?}{color_reset}\n", updated_tickets.clone());

        let winners: Vec<Ticket> = updated_tickets[..(AMOUNT_OF_WINNER as usize)].to_vec();
        /* First three tickets that are near to the rnd point are the winners and the
        * prize will be distributed among these guys,
        * 1st person: 50%  |  2nd person: 30%  |  3th person: 20% */
        self.winner = Some(winners);
        self.active = false;
        Ok(())

    }   
    // To avoid same ticket number result.
    pub fn verify_ticket_number_uniqueness(&self, _t_number_to_examine: u64) -> bool {
        // verifying ticket using binary search. if true: the number is useable.
        self.tickets.binary_search_by(|t| t.ticket_number.cmp(&_t_number_to_examine)).is_err()
    }
}

impl Space for Ticket {
    const INIT_SPACE: usize = 8 + 32 + 8;
}

#[account]
pub struct RaffleFeeVault {
    pub bump: u8
}
impl Space for RaffleFeeVault {
    const INIT_SPACE: usize= 8 + 1;
}

#[error_code]
pub enum RaffleErrorCode {
    RaffleAlreadyClosed,
    RaffleAlreadyDeactivated,
    WinnerAlreadySelected
}


#[cfg(test)]
mod tests {
    use std::borrow::BorrowMut;
    use inline_colorization::*;
    use super::*;

    pub const RAFFLE_OWNER_PUBKEY: Pubkey = pubkey!("So11111111111111111111111111111111111111112");

    pub const RAFFLE_PARTICIPANT1: Pubkey = pubkey!("So11111111111111111111111111111111111111113");
    pub const RAFFLE_PARTICIPANT2: Pubkey = pubkey!("So11111111111111111111111111111111111111114");
    pub const RAFFLE_PARTICIPANT3: Pubkey = pubkey!("So11111111111111111111111111111111111111115");
    pub const RAFFLE_PARTICIPANT4: Pubkey = pubkey!("So11111111111111111111111111111111111111116");
    pub const RAFFLE_PARTICIPANT5: Pubkey = pubkey!("So11111111111111111111111111111111111111117");


    pub const TEST_RAFFLE_TICKET_NUMBER: u64 = 5244;
    pub const TEST_VRF_NUMBER_IN_RANGE: u64 = 342349698505;

    pub fn setup() -> RaffleInfo {
        RaffleInfo {
            ticket_number_bound: (TICKET_NUMBER_UPPER_RANGE, TICKET_NUMBER_LOWER_RANGE),
            raffle_id: 0,                   ticket_price: TICKET_PRICE,
            total_ticket_sold: 0,           start_time: 1727163188,
            creator: RAFFLE_OWNER_PUBKEY,   tickets: Vec::new(),
            winner: Option::None,           is_closed: false,
            active: true,                   raffle_bump: 1_u8   
        }
    }
    #[test]
    fn test_buy_ticket() {
        let raffle: RaffleInfo = setup();
        let mut binding = raffle.clone();
        let  m_raffle = binding.borrow_mut();
        
        m_raffle.buy_ticket(RAFFLE_PARTICIPANT1, TEST_RAFFLE_TICKET_NUMBER).unwrap();

        assert_eq!(m_raffle.total_ticket_sold, 1);
        assert_eq!(m_raffle.tickets.len(), 1);
        assert_eq!(m_raffle.tickets[0].owner, RAFFLE_PARTICIPANT1);
        assert_eq!(m_raffle.tickets[0].ticket_number, TEST_RAFFLE_TICKET_NUMBER);
    }

    #[test]
    fn test_select_winner() {
        let raffle: RaffleInfo = setup();
        let mut binding = raffle.clone();
        let  m_raffle = binding.borrow_mut();
        // ticket numbers are not in this order obvs.
        m_raffle.buy_ticket(RAFFLE_PARTICIPANT1, TEST_RAFFLE_TICKET_NUMBER).unwrap();
        m_raffle.buy_ticket(RAFFLE_PARTICIPANT2, TEST_RAFFLE_TICKET_NUMBER + 5).unwrap();
        m_raffle.buy_ticket(RAFFLE_PARTICIPANT4, TEST_RAFFLE_TICKET_NUMBER - 30).unwrap();
        m_raffle.buy_ticket(RAFFLE_PARTICIPANT3, TEST_RAFFLE_TICKET_NUMBER + 25).unwrap();
        m_raffle.buy_ticket(RAFFLE_PARTICIPANT5, TEST_RAFFLE_TICKET_NUMBER - 15).unwrap();

        assert_eq!(m_raffle.total_ticket_sold, 5);
        assert_eq!(m_raffle.tickets.len(), 5);

        // the result of vrf based on the TEST_VRF_NUMBER_IN_RANGE IS 5815 
        m_raffle.select_winner(TEST_VRF_NUMBER_IN_RANGE).unwrap();
        println!("{color_green}Winners: {:?}{color_reset}", m_raffle.winner);
        assert_eq!(m_raffle.clone().winner.unwrap()[0].ticket_number, 5269);
        assert_eq!(m_raffle.clone().winner.unwrap()[1].ticket_number, 5249);
        assert_eq!(m_raffle.clone().winner.unwrap()[2].ticket_number, 5244);

        assert!(!m_raffle.active);
    }
    
    #[test]
    fn test_verify_ticket_number_uniqueness() {
        let raffle: RaffleInfo = setup();
        let mut binding = raffle.clone();
        let  m_raffle = binding.borrow_mut();

        m_raffle.buy_ticket(RAFFLE_PARTICIPANT1, TEST_RAFFLE_TICKET_NUMBER).unwrap();
        m_raffle.buy_ticket(RAFFLE_PARTICIPANT2, TEST_RAFFLE_TICKET_NUMBER + 5).unwrap();
        m_raffle.buy_ticket(RAFFLE_PARTICIPANT3, TEST_RAFFLE_TICKET_NUMBER + 25).unwrap();
        m_raffle.buy_ticket(RAFFLE_PARTICIPANT4, TEST_RAFFLE_TICKET_NUMBER - 30).unwrap();
        m_raffle.buy_ticket(RAFFLE_PARTICIPANT5, TEST_RAFFLE_TICKET_NUMBER - 15).unwrap();
        m_raffle.buy_ticket(RAFFLE_PARTICIPANT5, 5815).unwrap();

        m_raffle.select_winner(TEST_RAFFLE_TICKET_NUMBER).unwrap();

        let not_verify: bool = m_raffle.verify_ticket_number_uniqueness(5815);
        let verify: bool = m_raffle.verify_ticket_number_uniqueness(5815 + 1);
        assert!(!not_verify);
        assert!(verify);
    }
}