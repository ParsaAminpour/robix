use anchor_lang::prelude::*;

declare_id!("AJq49zQoJr8MjU1CHD6XietaT84YiQyRhAYM9FE74ApZ");

#[program]
pub mod coinflip {
    use super::*;

    pub fn initialize(ctx: Context<Initialize>) -> Result<()> {
        msg!("Greetings from: {:?}", ctx.program_id);
        Ok(())
    }
}

#[derive(Accounts)]
pub struct Initialize {}
