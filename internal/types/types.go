package types

// Validator represents a single validator in the network
type Validator struct {
    // Core fields
    Pubkey                     [48]byte `json:"pubkey,omitempty"`
    WithdrawalCredentials      [32]byte `json:"withdrawal_credentials,omitempty"`
    EffectiveBalance          uint64   `json:"effective_balance"`
    Slashed                   bool     `json:"slashed"`
    
    // Activation and exit epochs
    ActivationEligibilityEpoch uint64   `json:"activation_eligibility_epoch"`
    ActivationEpoch           uint64   `json:"activation_epoch"`
    ExitEpoch                 uint64   `json:"exit_epoch"`
    WithdrawableEpoch         uint64   `json:"withdrawable_epoch"`
    
    // For penalty calculations
    InactivityScore           uint64   `json:"inactivity_score"`
}

// NetworkState represents the current state of the network
type NetworkState struct {
    // Validators
    Validators         []Validator `json:"validators"`
    TotalActiveBalance uint64      `json:"total_active_balance"`
    
    // Epoch information
    CurrentEpoch       uint64      `json:"current_epoch"`
    FinalizedEpoch     uint64      `json:"finalized_epoch"`
    JustifiedEpoch     uint64      `json:"justified_epoch"`
    
    // Fork information
    CurrentFork        string      `json:"current_fork"`
    
    // Slashing tracking
    SlashingsPerEpoch  []uint64    `json:"slashings_per_epoch,omitempty"`
}

// RewardResults contains all calculated reward information
type RewardResults struct {
    // Input parameters
    ValidatorCount     int         `json:"validator_count"`
    TotalStaked       uint64      `json:"total_staked_gwei"`
    ParticipationRate float64     `json:"participation_rate"`
    
    // Base calculations
    SqrtTotalBalance   uint64      `json:"sqrt_total_balance"`
    BaseRewardPerEpoch uint64      `json:"base_reward_per_epoch"`
    
    // Component rewards (per epoch)
    SourceReward       uint64      `json:"source_reward"`
    TargetReward       uint64      `json:"target_reward"`
    HeadReward         uint64      `json:"head_reward"`
    AttestationRewardPerEpoch uint64 `json:"attestation_reward_per_epoch"`
    
    // Proposer calculations
    ProposerProbability       float64 `json:"proposer_probability"`
    ExpectedProposalsPerYear  float64 `json:"expected_proposals_per_year"`
    AvgProposerRewardPerBlock float64 `json:"avg_proposer_reward_per_block"`
    ProposerRewardPerEpoch    float64 `json:"proposer_reward_per_epoch"`
    
    // Annual projections
    AttestationRewardsAnnual  float64 `json:"attestation_rewards_annual"`
    ProposerRewardsAnnual     float64 `json:"proposer_rewards_annual"`
    TotalAnnualRewards        float64 `json:"total_annual_rewards"`
    APY                       float64 `json:"apy_percentage"`
    
    // Time-based projections
    DailyRewards   float64 `json:"daily_rewards"`
    WeeklyRewards  float64 `json:"weekly_rewards"`
    MonthlyRewards float64 `json:"monthly_rewards"`
}

// PenaltyResults contains penalty calculations
type PenaltyResults struct {
    // Attestation penalties
    SourcePenalty           uint64 `json:"source_penalty"`
    TargetPenalty           uint64 `json:"target_penalty"`
    HeadPenalty             uint64 `json:"head_penalty"`
    TotalAttestationPenalty uint64 `json:"total_attestation_penalty"`
    
    // Inactivity penalties
    InactivityScore   uint64 `json:"inactivity_score"`
    InactivityPenalty uint64 `json:"inactivity_penalty"`
    
    // Daily projections
    DailyAttestationPenalty float64 `json:"daily_attestation_penalty_eth"`
    DailyInactivityPenalty  float64 `json:"daily_inactivity_penalty_eth"`
}

// SlashingResults contains slashing penalty calculations
type SlashingResults struct {
    InitialPenalty       uint64  `json:"initial_penalty"`
    ProportionalPenalty  uint64  `json:"proportional_penalty"`
    TotalPenalty         uint64  `json:"total_penalty"`
    PercentageOfStake    float64 `json:"percentage_of_stake"`
    WhistleblowerReward  uint64  `json:"whistleblower_reward"`
    ProposerReward       uint64  `json:"proposer_reward"`
}

// ComparisonResult for comparing different validator counts
type ComparisonResult struct {
    ValidatorCount int     `json:"validator_count"`
    TotalStaked    uint64  `json:"total_staked_eth"`
    BaseReward     uint64  `json:"base_reward_gwei"`
    AnnualRewards  float64 `json:"annual_rewards_eth"`
    APY            float64 `json:"apy_percentage"`
    DailyRewards   float64 `json:"daily_rewards_eth"`
}

// DetailedBreakdown provides comprehensive reward breakdown
type DetailedBreakdown struct {
    RewardResults    *RewardResults    `json:"reward_results"`
    PenaltyResults   *PenaltyResults   `json:"penalty_results,omitempty"`
    SlashingResults  *SlashingResults  `json:"slashing_results,omitempty"`
    NetworkMetrics   *NetworkMetrics   `json:"network_metrics"`
}

// NetworkMetrics contains additional network statistics
type NetworkMetrics struct {
    // Issuance metrics
    NewIssuancePerEpoch  uint64  `json:"new_issuance_per_epoch"`
    NewIssuancePerYear   float64 `json:"new_issuance_per_year_eth"`
    InflationRate        float64 `json:"inflation_rate_percentage"`
    
    // Network participation
    ActiveValidators     int     `json:"active_validators"`
    TotalValidators      int     `json:"total_validators"`
    NetworkParticipation float64 `json:"network_participation_rate"`
    
    // Economic metrics
    TotalSupply          uint64  `json:"total_supply_eth"`
    StakedPercentage     float64 `json:"staked_percentage"`
    YieldPerValidator    float64 `json:"yield_per_validator_eth"`
}

// ValidatorPerformance tracks individual validator metrics
type ValidatorPerformance struct {
    ValidatorIndex       int     `json:"validator_index"`
    EffectiveBalance     uint64  `json:"effective_balance"`
    AttestationAccuracy  float64 `json:"attestation_accuracy"`
    ProposerDuties       int     `json:"proposer_duties"`
    TotalRewards         uint64  `json:"total_rewards"`
    TotalPenalties       uint64  `json:"total_penalties"`
    NetEarnings          int64   `json:"net_earnings"`
}