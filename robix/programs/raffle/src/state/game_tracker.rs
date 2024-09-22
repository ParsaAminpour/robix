use anchor_lang::{AnchorSerialize, AnchorDeserialize};
use anchor_lang::prelude::*;
// use crate::constants::*;

#[account]
pub struct GameTracker {
    pub active_raffle: u64,
    pub points: Vec<UserPoint>,
    pub bump: u8,
}

#[derive(AnchorSerialize, AnchorDeserialize, Clone, Default)]
pub struct UserPoint {
    pub owner: Pubkey,
    pub xp: u64
}

impl GameTracker {
    pub fn get_size(users_count: usize) -> usize {
        return 8 // descriminator
        + 8 // active_raffle index
        + (users_count * UserPoint::INIT_SPACE)
        + 1
    }

    // pub fn generate_next_raffle(&mut self) -> (Pubkey, u8) {
    //     let (next_raffle_pda, next_raffle_bump) = Pubkey::find_program::address(
    //         &[
    //             constants::RAFFLE_INIT_SEED.as_ref(),
    //             &(self.active_raffle + 1).to_le_bytes(),
    //         ],
    //     );
    // }
}

impl Space for UserPoint {
    const INIT_SPACE: usize = 32 + 8;
}