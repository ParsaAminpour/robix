use anchor_lang::{AnchorSerialize, AnchorDeserialize};
use anchor_lang::prelude::*;

use crate::constants::RAFFLE_SEED;
use crate::id::ID;

#[account]
pub struct GameTracker {
    pub active_raffle: u64,
    pub active_raffle_owner: Pubkey,
    pub points: Vec<UserPoint>,
    pub bump: u8,
}

impl GameTracker {
    pub fn get_space(users_count: usize) -> usize {
          8     // descriminator
        + 8     // active_raffle index
        + (users_count * UserPoint::INIT_SPACE)
        + 1     // u8 bump
    }

    pub fn initialize(&mut self, gt_bump: u8, active_raffle_owner: Pubkey) {
        self.active_raffle = 0; // first raffle had index 0
        self.active_raffle_owner = active_raffle_owner;
        self.points = vec![];
        self.bump = gt_bump;
    }

    // Calling this function at the end of the Raffle to generate new PDA for new Raffle's round.
    pub fn generate_next_raffle_pda(&self, system_owner: Pubkey) -> (Pubkey, u8) {
        // @audit how to ensure that the exp_pda and exp_bump are 100 percent valid.
        let (exp_pda, exp_bump) = Pubkey::find_program_address(
            &[RAFFLE_SEED.as_ref(), &(self.active_raffle + 1).to_le_bytes(), system_owner.as_ref()], 
            &ID
        );
        msg!("Next Raffle PDA: {} | Next Raffle bump: {}", exp_pda.to_string(), exp_bump);
        (exp_pda, exp_bump)
    }

    pub fn add_raffle_idx(&mut self) {
        self.active_raffle.checked_add(1).unwrap();
    }
}

#[derive(AnchorSerialize, AnchorDeserialize, Clone, Default)]
pub struct UserPoint {
    pub owner: Pubkey,
    pub xp: u64
}

impl Space for UserPoint {
    const INIT_SPACE: usize = 32 + 8;
}