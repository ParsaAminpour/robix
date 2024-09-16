use anchor_lang::prelude::*;

declare_id!("7iKkN561Q2C5w9ooaf5U7LHnVtH3VyyErWiiUr1TJcRk");

#[program]
pub mod raffle {
    use super::*;

    pub fn initialize(ctx: Context<Initialize>) -> Result<()> {
        msg!("Greetings from: {:?}", ctx.program_id);
        Ok(())
    }
}

#[derive(Accounts)]
pub struct Initialize {}
