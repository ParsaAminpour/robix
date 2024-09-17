use anchor_lang::prelude::*;

declare_id!("A7tUzPVsRRmPLMAR16Qu8KCS6CS8RpiG3QanqGF83oDN");

// @audit-info transfer these constants to a seperate module.
const MINIMUM_WHEEL_NAME_SIZE: usize = 3;
const MINIMUM_PARTICIPANTS_COUNT: usize = 2;
const MINIMUM_WHEEL_LIFE_CYCLE: u64 = 1;

#[program]
pub mod wheel {
    use super::*;

    pub fn InitializeSpin(ctx: Context<InitializeSpin>, _wheel_name: String, _min_participants: u8, _end_time: u64, _creator: Pubkey) -> Result<()> {
        require_gte!(_wheel_name.len(), MINIMUM_WHEEL_NAME_SIZE, ErrorCode::WheelNameIsTooShort);
        require_gte!(_min_participants, MINIMUM_PARTICIPANTS_COUNT as u8, ErrorCode::ParticipantsNumberExceeded);
        // require
        Ok(())
    }
}

#[derive(Accounts)]
#[instruction(wheel_name: String, start_time: u64)]
pub struct InitializeSpin<'info> {
    #[account(
        init,
        space = WheelInfo::WHEEL_INFO_SPACE,
        payer = signer,
        owner = signer.key(),
        seeds = [b"wheel".as_ref(), wheel_name.as_bytes(), start_time.to_string().as_bytes()],
        bump
    )]
    pub wheel_info: Account<'info, WheelInfo>,

    #[account(
        mut, seeds = [b"treasury".as_ref()], bump
    )]
    pub treasury: SystemAccount<'info>,

    #[account(mut)]
    pub signer: Signer<'info>,
    pub system_program: Program<'info, System>
}


#[account]
pub struct WheelInfo {
    pub wheel_name: String,
    pub wheel_pool: u64,
    pub min_participants: u8,
    pub start_time: u64,
    pub end_time: u64,
    pub creator: Pubkey,
    pub is_closed: bool,
    pub participants: Vec<Pubkey>,
    pub treasury_bump: u8, 
}

impl WheelInfo {
    pub const WHEEL_INFO_SPACE: usize = std::mem::size_of::<WheelInfo>();
    
    pub fn init(&mut self, _wheel_name: String, _min_participants: u8, _start_time: u64, _end_time: u64, _creator: Pubkey) -> Result<()> {
        self.wheel_name = _wheel_name;
        self.wheel_pool = 0;
        self.min_participants = _min_participants;
        self.start_time = Clock::get().unwrap().unix_timestamp as u64;
        self.end_time = _end_time;
        self.creator = _creator;
        self.is_closed = false;
        self.participants = Vec::new();

        Ok(())
    }
}

#[error_code]
pub enum ErrorCode {
    WheelNameIsTooShort,
    ParticipantsNumberExceeded
}