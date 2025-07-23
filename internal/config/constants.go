package config

// Reward and penalty constants from Ethereum mainnet
const (
    // Base parameters
    BASE_REWARD_FACTOR             = 64
    BASE_REWARDS_PER_EPOCH         = 4
    PROPOSER_REWARD_QUOTIENT       = 8
    WHISTLEBLOWER_REWARD_QUOTIENT  = 512
    MIN_SLASHING_PENALTY_QUOTIENT  = 128
    PROPORTIONAL_SLASHING_MULTIPLIER = 1
    
    // Altair parameters
    INACTIVITY_PENALTY_QUOTIENT_ALTAIR     = 50331648  // 2**24
    MIN_SLASHING_PENALTY_QUOTIENT_ALTAIR   = 64
    PROPORTIONAL_SLASHING_MULTIPLIER_ALTAIR = 2
    
    // Bellatrix parameters
    INACTIVITY_PENALTY_QUOTIENT_BELLATRIX     = 33554432  // 2**25
    MIN_SLASHING_PENALTY_QUOTIENT_BELLATRIX   = 32
    PROPORTIONAL_SLASHING_MULTIPLIER_BELLATRIX = 3
    
    // Phase 0 parameters (for backwards compatibility)
    INACTIVITY_PENALTY_QUOTIENT    = 67108864  // 2**26
    INACTIVITY_SCORE_BIAS          = 4
    INACTIVITY_SCORE_RECOVERY_RATE = 16
    
    // Participation flag weights
    // TIMELY_SOURCE_WEIGHT = 6
    // TIMELY_TARGET_WEIGHT = 10
    // TIMELY_HEAD_WEIGHT   = 6
    // SYNC_REWARD_WEIGHT   = 1
    // PROPOSER_WEIGHT      = 3
    // WEIGHT_DENOMINATOR   = 26

	TIMELY_SOURCE_WEIGHT = 14
    TIMELY_TARGET_WEIGHT = 26
    TIMELY_HEAD_WEIGHT   = 14
    SYNC_REWARD_WEIGHT   = 2
    PROPOSER_WEIGHT      = 8
    WEIGHT_DENOMINATOR   = 64
    
    // Sync committee
    SYNC_COMMITTEE_SIZE                   = 512
    SYNC_COMMITTEE_SUBNET_COUNT          = 4
    SYNC_REWARD_WEIGHT_DENOMINATOR       = 2
    
    // Balance parameters
    EFFECTIVE_BALANCE_INCREMENT = 1000000000  // 1 ETH in Gwei
    MAX_EFFECTIVE_BALANCE       = 32000000000 // 32 ETH in Gwei
    EJECTION_BALANCE           = 16000000000 // 16 ETH in Gwei
    
    // Time parameters
    SLOTS_PER_EPOCH                  = 32
    EPOCHS_PER_YEAR                  = 82180 // 365.25 * 225
    EPOCHS_PER_DAY                   = 225
    EPOCHS_PER_WEEK                  = 1575
    EPOCHS_PER_MONTH                 = 6848
    SECONDS_PER_SLOT                 = 12
    MIN_ATTESTATION_INCLUSION_DELAY  = 1

	// SLOTS_PER_EPOCH                  = 32
    // EPOCHS_PER_YEAR                  = 98618 // 365.25 * 270
    // EPOCHS_PER_DAY                   = 270
    // EPOCHS_PER_WEEK                  = 1890
    // EPOCHS_PER_MONTH                 = 8219
    // SECONDS_PER_SLOT                 = 10
    // MIN_ATTESTATION_INCLUSION_DELAY  = 1
    
    // Fork versions (for reference)
    PHASE0_FORK_VERSION    = "0x00000000"
    ALTAIR_FORK_VERSION    = "0x01000000"
    BELLATRIX_FORK_VERSION = "0x02000000"
    CAPELLA_FORK_VERSION   = "0x03000000"
    DENEB_FORK_VERSION     = "0x04000000"
    ELECTRA_FORK_VERSION   = "0x05000000"
    
    // Validator set limits
    MIN_GENESIS_ACTIVE_VALIDATOR_COUNT = 16384
    CHURN_LIMIT_QUOTIENT              = 65536
    MIN_PER_EPOCH_CHURN_LIMIT         = 4
    MAX_PER_EPOCH_ACTIVATION_CHURN_LIMIT = 8
    
    // Slashing
    EPOCHS_PER_SLASHINGS_VECTOR = 8192
    WHISTLEBLOWER_REWARD_PROPORTION = 8 // 1/8 of validator effective balance
    
    // Withdrawals
    MAX_VALIDATORS_PER_WITHDRAWALS_SWEEP = 16384
    MAX_WITHDRAWALS_PER_PAYLOAD = 16
)

// Fork configuration
type ForkConfig struct {
    Version                       string
    InactivityPenaltyQuotient    uint64
    MinSlashingPenaltyQuotient   uint64
    ProportionalSlashingMultiplier uint64
}

// GetForkConfig returns configuration for a specific fork
func GetForkConfig(fork string) ForkConfig {
    switch fork {
    case "phase0":
        return ForkConfig{
            Version:                       PHASE0_FORK_VERSION,
            InactivityPenaltyQuotient:    INACTIVITY_PENALTY_QUOTIENT,
            MinSlashingPenaltyQuotient:   MIN_SLASHING_PENALTY_QUOTIENT,
            ProportionalSlashingMultiplier: PROPORTIONAL_SLASHING_MULTIPLIER,
        }
    case "altair":
        return ForkConfig{
            Version:                       ALTAIR_FORK_VERSION,
            InactivityPenaltyQuotient:    INACTIVITY_PENALTY_QUOTIENT_ALTAIR,
            MinSlashingPenaltyQuotient:   MIN_SLASHING_PENALTY_QUOTIENT_ALTAIR,
            ProportionalSlashingMultiplier: PROPORTIONAL_SLASHING_MULTIPLIER_ALTAIR,
        }
    case "bellatrix", "merge":
        return ForkConfig{
            Version:                       BELLATRIX_FORK_VERSION,
            InactivityPenaltyQuotient:    INACTIVITY_PENALTY_QUOTIENT_BELLATRIX,
            MinSlashingPenaltyQuotient:   MIN_SLASHING_PENALTY_QUOTIENT_BELLATRIX,
            ProportionalSlashingMultiplier: PROPORTIONAL_SLASHING_MULTIPLIER_BELLATRIX,
        }
    default:
        // Return latest (Bellatrix) config as default
        return GetForkConfig("bellatrix")
    }
}