use anchor_lang::prelude::*;

#[error_code]
pub enum ErrorCode {
    #[msg("Invalid amount to send, zero anount is not allowed in transfer")]
    InvalidAmount,
    #[msg("Empty Raffle name is not allowed")]
    EmptyStringNotAllowed,
    #[msg("Invalid time range specified for end time")]
    InvalidEndTime,
    #[msg("The tickets amount reached to the max tickets amount threshod")]
    TicketAmountThreshold,
    #[msg("The winner related to this raffle has already selected")]
    WinnerAlreadySelected,
    #[msg("Insufficient amount to transfer")]
    InsufficientBalance,
    WinnerIsNotSelectedYet,
    RaffleExpired,
    SameDestinationAddressNotAllowed,
    RaffleIsNotClosed,
    CallerIsNotWinner,
    NotSufficientBalance
}
