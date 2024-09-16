use anchor_lang::prelude::*;

declare_id!("A7tUzPVsRRmPLMAR16Qu8KCS6CS8RpiG3QanqGF83oDN");

#[program]
pub mod wheel {
    use super::*;

    pub fn initialize(ctx: Context<Initialize>) -> Result<()> {
        msg!("Greetings from: {:?}", ctx.program_id);
        Ok(())
    }
}

#[derive(Accounts)]
pub struct Initialize {}
