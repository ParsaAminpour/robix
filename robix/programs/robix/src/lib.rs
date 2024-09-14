use anchor_lang::prelude::*;

declare_id!("9GaG1Kh8mdVMY3UZWbE3sYikwhiM8qyuup9GzgHUFxCh");

#[program]
pub mod robix {
    use super::*;

    pub fn initialize(ctx: Context<Initialize>) -> Result<()> {
        msg!("Greetings from: {:?}", ctx.program_id);
        Ok(())
    }
}

#[derive(Accounts)]
pub struct Initialize {}
