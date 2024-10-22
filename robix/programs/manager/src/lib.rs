use core::fmt;

use anchor_lang::prelude::*;
use anchor_lang::{AnchorSerialize, AnchorDeserialize};
// use std::str::FromStr;
#[cfg(not(feature = "no-entrypoint"))]
use solana_security_txt::security_txt;

#[cfg(not(feature = "no-entrypoint"))]
security_txt! {
    name: "Robix",
    project_url: "https://github.com/ParsaAminpour/robix",
    contacts: "email:parsa.aminpour@gmail.com",
    policy: "https://github.com/solana-labs/solana/blob/master/SECURITY.md"
}

declare_id!("88ZpZY4cseY2fU2WY8u7hYyQSzp8t36iJSd4XcqfKh6V");

#[program]
pub mod manager { 
    use super::*;

    pub fn create_account(ctx: Context<CreateAccount>, amount: u64) -> anchor_lang::Result<()> {
        match ctx.accounts.funding(
            ctx.accounts.system_program.to_account_info(), 
            ctx.accounts.signer.to_account_info(), 
            ctx.accounts.treasury.to_account_info(), 
            amount, 
            true,
            &[&[&[1_u8]]]
        )
        {
            Ok(_) => {
                ctx.accounts.account_data.add_balance(amount).unwrap();
                ctx.accounts.account_data.init(ctx.accounts.signer.key(), amount, 0).unwrap();
            },
            Err(_) => return err!(ErrorCode::TransferFailed),
        }

        Ok(())
    }

    pub fn depoit_fund(ctx: Context<DepositFund>, amount: u64) -> anchor_lang::Result<()> {
        require_eq!(ctx.accounts.account_data.owner, ctx.accounts.signer.key(), ErrorCode::SignerIsNotOwner);
        require_gt!(amount, 0, ErrorCode::ZeroAmountNotAllowed);

        match ctx.accounts.funding(
            ctx.accounts.system_program.to_account_info(),
            ctx.accounts.signer.to_account_info(),
            ctx.accounts.treasury.to_account_info(),
            amount,
            true,
            &[&[&[1_u8]]])
        {
            Ok(_) => ctx.accounts.account_data.add_balance(amount).unwrap(),
            Err(_) => return err!(ErrorCode::TransferFailed),
        }
        Ok(())
    }

    pub fn withdraw_fund(ctx: Context<WithdrawFund>, amount: u64) -> anchor_lang::Result<()> {
        require_eq!(ctx.accounts.account_data.owner, ctx.accounts.signer.key(), ErrorCode::SignerIsNotOwner);
        require_gte!(ctx.accounts.account_data.balance - ctx.accounts.account_data.locked_balance, amount, ErrorCode::InsufficientBalance);
        require_gt!(amount, 0, ErrorCode::ZeroAmountNotAllowed);
        
        let bump = &[ctx.bumps.treasury];
        let seeds: &[&[u8]] = &[b"treasury".as_ref(), bump];
        let signer_seeds = &[seeds];

        match ctx.accounts.funding(
            ctx.accounts.system_program.to_account_info(),
            ctx.accounts.treasury.to_account_info(),
            ctx.accounts.signer.to_account_info(),
            amount,
            false,
            signer_seeds)
        {
            Ok(_) => ctx.accounts.account_data.sub_balance(amount).unwrap(),
            Err(_) => return err!(ErrorCode::TransferFailed),
        }
        Ok(())
    }
}


// @audit should consider re-initialization issue.
#[derive(Accounts)]
pub struct CreateAccount<'info> {
    #[account(
        init_if_needed,
        payer = signer,
        owner = signer.key(),
        space = AccountData::ABSTRACTION_ACCOUNT_SPACE,
        seeds = [b"account".as_ref(), signer.key().as_ref()],
        bump
    )]
    pub account_data: Account<'info, AccountData>,

    // @audit adding ownership constraints and owner management for the treasury..
    #[account(mut, seeds=[b"treasury".as_ref()], bump)]
    pub treasury: SystemAccount<'info>,

    #[account(mut)]
    pub signer: Signer<'info>,
    pub system_program: Program<'info, System>
}

impl<'a> Transfer<'a> for CreateAccount<'a> {}


#[derive(Accounts)]
pub struct DepositFund<'info> {
    #[account(
        mut, // @audit init check require for constraint section.
        constraint = account_data.owner == signer.key() @ErrorCode::SignerIsNotOwner,
        seeds = [b"account".as_ref(), signer.key().as_ref()],
        bump
    )]
    pub account_data: Account<'info, AccountData>,

    #[account(mut, seeds=[b"treasury".as_ref()], bump)]
    pub treasury: SystemAccount<'info>,

    #[account(mut)]
    pub signer: Signer<'info>,
    pub system_program: Program<'info, System>
}

impl <'a> Transfer<'a> for DepositFund<'a> {}

#[derive(Accounts)]
#[instruction(amount: u64)]
pub struct WithdrawFund<'info> {
    #[account(
        mut,
        constraint = account_data.owner == signer.key() @ErrorCode::SignerIsNotOwner,
        seeds = [b"account".as_ref(), signer.key().as_ref()],
        bump
    )]
    pub account_data: Account<'info, AccountData>,

    #[account(mut, seeds=[b"treasury".as_ref()], bump)]
    pub treasury: SystemAccount<'info>,

    #[account(mut)]
    pub signer: Signer<'info>,
    pub system_program: Program<'info, System>
}

impl<'a> Transfer<'a> for WithdrawFund<'a> {}


#[account]
pub struct AccountData {
    pub owner: Pubkey,
    // @audit-info intense security checks should be applied on this arg.
    pub balance: u64,
    // when user participate in any game
    pub locked_balance: u64,
}

impl fmt::Debug for AccountData {
    fn fmt(&self, f: &mut fmt::Formatter<'_>) -> fmt::Result {
        let mut f = f.debug_struct("AccountData");

        f.field("owner", &self.owner)
            .field("balance", &self.balance)
            .field("locked_balance", &self.locked_balance);
        f.finish_non_exhaustive()
    }
}

impl AccountData {
    pub const ABSTRACTION_ACCOUNT_SPACE: usize = 8 + 32 + 4 + 4;
    pub const _SYSTEM_ABSTRACTION_ACCOUNT_SPACE: usize = std::mem::size_of::<AccountData>();

    pub fn init(&mut self, _owner: Pubkey, _balance: u64, _locked_balance: u64) -> Result<()> {
        self.owner = _owner;
        self.balance = _balance;
        self.locked_balance = _locked_balance;
        Ok(())
    }
    pub fn add_balance(&mut self, amount: u64) -> Result<()> {
        require_gt!(amount, 0, ErrorCode::ZeroAmountNotAllowed);
        self.balance.checked_add(amount).unwrap();
        Ok(())
    }
    pub fn sub_balance(&mut self, amount: u64) -> Result<()> {
        require_gte!(self.balance, amount, ErrorCode::BalanceIsLessThanAmountToDecrease);
        require_gte!(self.balance - self.locked_balance, amount, ErrorCode::BalanceIsLessThanAmountToDecrease);
        self.balance.checked_sub(amount).unwrap();
        Ok(())
    }
    pub fn aa_key(&self) -> Pubkey {
        self.owner
    }
    pub fn get_balance(&self) -> u64 {
        self.balance
    }
}

trait Transfer<'a> {
    fn funding(&mut self, system_program: AccountInfo<'a>, sender: AccountInfo<'a>, receiver: AccountInfo<'a>, amount: u64, direction: bool, signer_seeds: &[&[&[u8]]]) -> Result<bool> {
        require_gte!(sender.get_lamports(), amount, ErrorCode::InsufficientBalance);
        require_neq!(sender.key(), receiver.key(), ErrorCode::SameDestination);
        if direction {
            // transfer from user to treasury
            anchor_lang::system_program::transfer(
                CpiContext::new(
                    system_program.to_account_info(), 
                anchor_lang::system_program::Transfer {
                    from: sender.to_account_info(),
                    to: receiver.to_account_info(),
                }), 
                amount).unwrap();
                msg!("SOL transfered from treasury to {}", receiver.key());
                
            } else {
                // vice versa
                anchor_lang::system_program::transfer(
                    CpiContext::new(
                        system_program.to_account_info(),
                        anchor_lang::system_program::Transfer {
                            from: sender.to_account_info(),
                            to: receiver.to_account_info()
                        }
                    ).with_signer(signer_seeds), 
                    amount
                ).unwrap();
                msg!("SOL transfered from {} to treasury", receiver.key());
            }
            Ok(true)
        }
}

#[error_code]
pub enum ErrorCode {
    #[msg("test")]
    TransferFailed,
    SameDestination,
    InsufficientBalance,
    ZeroAmountNotAllowed,
    BalanceIsLessThanAmountToDecrease,
    #[msg("When calling an instruction for an account which you are not the owner")]
    SignerIsNotOwner
}

pub enum AccountType {
    UserAccount,
    TreasuryAccount,
    CommunityAccount
}