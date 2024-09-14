use anchor_lang::prelude::*;

declare_id!("88ZpZY4cseY2fU2WY8u7hYyQSzp8t36iJSd4XcqfKh6V");

#[program]
pub mod manager {
    use super::*;

    pub fn initialize(ctx: Context<Initialize>) -> Result<()> {
        msg!("Greetings from: {:?}", ctx.program_id);
        Ok(())
    }
}

#[derive(Accounts)]
pub struct Initialize {}
