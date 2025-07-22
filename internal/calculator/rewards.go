package calculator

import (
    "math"
    
    "github.com/eth-rewards-calculator/internal/config"
    "github.com/eth-rewards-calculator/internal/types"
)

// CalculateRewards computes all reward components for the given network state
func CalculateRewards(state *types.NetworkState, participationRate float64) *types.RewardResults {
    validatorCount := len(state.Validators)
    
    // Calculate base reward for a validator with max effective balance
    baseReward := GetBaseReward(state, 0)
    sqrtTotal := IntegerSquareRoot(state.TotalActiveBalance)
    
    // Component rewards
    sourceReward := baseReward * config.TIMELY_SOURCE_WEIGHT / config.WEIGHT_DENOMINATOR
    targetReward := baseReward * config.TIMELY_TARGET_WEIGHT / config.WEIGHT_DENOMINATOR
    headReward := baseReward * config.TIMELY_HEAD_WEIGHT / config.WEIGHT_DENOMINATOR
    attestationReward := sourceReward + targetReward + headReward
    
    // Proposer calculations
    proposerProbability := 1.0 / float64(validatorCount)
    proposalsPerEpoch := proposerProbability
    proposalsPerYear := proposalsPerEpoch * float64(config.EPOCHS_PER_YEAR)
    
    // Average proposer reward per block (simplified)
    avgProposerReward := float64(baseReward) * float64(config.PROPOSER_WEIGHT) / 
                        float64(config.WEIGHT_DENOMINATOR)
    proposerRewardPerEpoch := avgProposerReward * proposerProbability
    
    // Annual calculations with participation rate
    attestationAnnual := float64(attestationReward) * float64(config.EPOCHS_PER_YEAR) * participationRate
    proposerAnnual := proposerRewardPerEpoch * float64(config.EPOCHS_PER_YEAR) * participationRate
    totalAnnual := attestationAnnual + proposerAnnual
    
    // APY calculation
    apy := (totalAnnual / float64(config.MAX_EFFECTIVE_BALANCE)) * 100
    
    return &types.RewardResults{
        // Input parameters
        ValidatorCount:     validatorCount,
        TotalStaked:       state.TotalActiveBalance,
        ParticipationRate: participationRate,
        
        // Base calculations
        SqrtTotalBalance:   sqrtTotal,
        BaseRewardPerEpoch: baseReward,
        
        // Component rewards
        SourceReward:              sourceReward,
        TargetReward:              targetReward,
        HeadReward:                headReward,
        AttestationRewardPerEpoch: attestationReward,
        
        // Proposer calculations
        ProposerProbability:       proposerProbability,
        ExpectedProposalsPerYear:  proposalsPerYear,
        AvgProposerRewardPerBlock: avgProposerReward,
        ProposerRewardPerEpoch:    proposerRewardPerEpoch,
        
        // Annual projections
        AttestationRewardsAnnual: attestationAnnual,
        ProposerRewardsAnnual:    proposerAnnual,
        TotalAnnualRewards:       totalAnnual,
        APY:                      apy,
        
        // Time-based projections
        DailyRewards:   totalAnnual / 365.25,
        WeeklyRewards:  totalAnnual / 52.18,
        MonthlyRewards: totalAnnual / 12,
    }
}

// GetBaseReward calculates the base reward for a validator
func GetBaseReward(state *types.NetworkState, validatorIndex int) uint64 {
    totalBalance := state.TotalActiveBalance
    effectiveBalance := state.Validators[validatorIndex].EffectiveBalance
    
    return effectiveBalance * config.BASE_REWARD_FACTOR / 
           IntegerSquareRoot(totalBalance) / config.BASE_REWARDS_PER_EPOCH
}

// GetBaseRewardPerIncrement calculates base reward per increment
func GetBaseRewardPerIncrement(state *types.NetworkState) uint64 {
    return config.EFFECTIVE_BALANCE_INCREMENT * config.BASE_REWARD_FACTOR / 
           IntegerSquareRoot(state.TotalActiveBalance) / config.BASE_REWARDS_PER_EPOCH
}

// CalculateAttestationReward computes reward for a single attestation
func CalculateAttestationReward(state *types.NetworkState, validatorIndex int,
    correctSource, correctTarget, correctHead bool, inclusionDelay uint64) uint64 {
    
    baseReward := GetBaseReward(state, validatorIndex)
    reward := uint64(0)
    
    if correctSource {
        reward += baseReward * config.TIMELY_SOURCE_WEIGHT / config.WEIGHT_DENOMINATOR
    }
    if correctTarget {
        reward += baseReward * config.TIMELY_TARGET_WEIGHT / config.WEIGHT_DENOMINATOR
    }
    if correctHead {
        reward += baseReward * config.TIMELY_HEAD_WEIGHT / config.WEIGHT_DENOMINATOR
    }
    
    // Apply inclusion delay penalty (for late attestations)
    if inclusionDelay > config.MIN_ATTESTATION_INCLUSION_DELAY && reward > 0 {
        reward = reward * config.PROPOSER_REWARD_QUOTIENT / 
                (config.PROPOSER_REWARD_QUOTIENT + inclusionDelay - config.MIN_ATTESTATION_INCLUSION_DELAY)
    }
    
    return reward
}

// CalculateProposerReward computes reward for block proposer
func CalculateProposerReward(state *types.NetworkState, attestingBalance uint64) uint64 {
    baseRewardPerIncrement := GetBaseRewardPerIncrement(state)
    proposerRewardPerIncrement := baseRewardPerIncrement / config.PROPOSER_REWARD_QUOTIENT
    
    return proposerRewardPerIncrement * attestingBalance / config.EFFECTIVE_BALANCE_INCREMENT
}

// CalculateSyncCommitteeReward computes sync committee participation reward
func CalculateSyncCommitteeReward(state *types.NetworkState, participantCount int) uint64 {
    baseReward := GetBaseReward(state, 0) // Assume max effective balance
    totalActiveIncrements := state.TotalActiveBalance / config.EFFECTIVE_BALANCE_INCREMENT
    totalBaseRewards := baseReward * totalActiveIncrements
    
    maxParticipantRewards := totalBaseRewards * config.SYNC_REWARD_WEIGHT / 
                            config.WEIGHT_DENOMINATOR / config.SLOTS_PER_EPOCH
    participantReward := maxParticipantRewards / config.SYNC_COMMITTEE_SIZE
    
    return participantReward * uint64(participantCount)
}

// CalculateWhistleblowerReward computes reward for reporting slashable offense
func CalculateWhistleblowerReward(slashedValidatorBalance uint64) (whistleblowerReward, proposerReward uint64) {
    whistleblowerReward = slashedValidatorBalance / config.WHISTLEBLOWER_REWARD_QUOTIENT
    proposerReward = whistleblowerReward / config.PROPOSER_REWARD_QUOTIENT
    return
}

// EstimateNetworkIssuance calculates total new issuance for the network
func EstimateNetworkIssuance(state *types.NetworkState, participationRate float64) *types.NetworkMetrics {
    validatorCount := len(state.Validators)
    
    // Calculate per-validator rewards
    results := CalculateRewards(state, participationRate)
    
    // Network-wide issuance
    totalIssuancePerEpoch := results.BaseRewardPerEpoch * 4 * uint64(validatorCount) * 
                            uint64(participationRate * float64(config.WEIGHT_DENOMINATOR)) / 
                            config.WEIGHT_DENOMINATOR
    
    totalIssuancePerYear := float64(totalIssuancePerEpoch) * float64(config.EPOCHS_PER_YEAR) / 1e9
    
    // Assume total ETH supply (this would need to be tracked properly)
    totalSupply := uint64(120_000_000) // Approximate ETH supply
    inflationRate := (totalIssuancePerYear / float64(totalSupply)) * 100
    
    return &types.NetworkMetrics{
        NewIssuancePerEpoch:  totalIssuancePerEpoch,
        NewIssuancePerYear:   totalIssuancePerYear,
        InflationRate:        inflationRate,
        ActiveValidators:     int(float64(validatorCount) * participationRate),
        TotalValidators:      validatorCount,
        NetworkParticipation: participationRate,
        TotalSupply:          totalSupply,
        StakedPercentage:     float64(state.TotalActiveBalance/1e9) / float64(totalSupply) * 100,
        YieldPerValidator:    results.TotalAnnualRewards / 1e9,
    }
}

// IntegerSquareRoot computes integer square root
func IntegerSquareRoot(n uint64) uint64 {
    if n == 0 {
        return 0
    }
    
    // Use floating point for initial guess
    x := uint64(math.Sqrt(float64(n)))
    
    // Newton's method for refinement
    for {
        x1 := (x + n/x) / 2
        if x1 >= x {
            return x
        }
        x = x1
    }
}