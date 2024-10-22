/// Fee Vault @audit should be changed these constants.
pub const FEE_VAULT: &str = "264r45MGssHfx4gnof4wu7uoSk44Rk5WFgPkRb5JrzA9,255";

/// This is the authority to initiatialize stuff.
pub const RAFFLE_OWNER: &str = "BGAdbX9mWxxZ1PKPy9J5bvvEhL5vbwJo9wodYafefEfK";  
// (DEV: BGAdbX9mWxxZ1PKPy9J5bvvEhL5vbwJo9wodYafefEfK)

/// The current version of the Raffle account.
const _RENT_ADDITION: u64 = 1_120_560; // amount for + rent of 33 bytes

/// The current ticket price.
pub const TICKET_PRICE: u64 = 500_000_000; // 0.5 SOL in lamports

/// The current fee collected per ticket.
pub const TICKET_FEE: u64 = 13_100_000; // 0.0131 SOL in lamports

/// The current fee collected per ticket.
pub const SUPER_RAFFLE_FEE: u64 = 6_900_000; // 0.0069 SOL in lamports

/// The  cost per new raffle (rounded up)
pub const NEW_RAFFLE_COST: u64 = 1_500_000; // 0.0015 SOL in lamports

/// The maximum number of tickets that can be purchased per user.
pub const MAX_TICKETS_PER_USER: u8 = 50;

/// The number of points per ticket.
pub const POINTS_PER_TICKET: u32 = 1;

/// The number of points for selling.
pub const POINTS_FOR_SELLING: u32 = 10;

/// Price Feeds
pub const SOL_PRICE_FEED: &str = "H6ARHf6YXhGYeQfUzQNGk6rDNnLBQKrenN712K4AQJEG";
pub const STALENESS_THRESHOLD: u64 = 1; // staleness threshold in seconds


pub const RAFFLE_INIT_SEED: &str = "raffle";

/// Amount of tickets for selecting winners.
pub const AMOUNT_OF_WINNER: u32 = 3;

/// The ticket number boundaries.
pub const TICKET_NUMBER_LOWER_RANGE: u64 = 5000;
pub const TICKET_NUMBER_UPPER_RANGE: u64 = 6000;

/// Anchor seeds for raffle pda generating using `create_program_address`
pub const RAFFLE_SEED: [u8; 6] = *b"raffle";

pub const TREASURY_SEED: [u8; 8] = *b"treasury";

pub const TRACKER_SEED: [u8; 7] = *b"tracker";

pub const FPRNG_PROGRAM_ADDRESS: &str = "9uSwASSU59XvUS8d1UeU8EwrEzMGFdXZvQ4JSEAfcS7k";