package calculator

import (
    "github.com/eth-rewards-calculator/internal/config"
    "github.com/eth-rewards-calculator/internal/types"
)

// CalculatePenalties computes attestation penalties for missed duties
func CalculatePenalties(state *types.NetworkState, validatorIndex int,
    correctSource, correctTarget, correctHead bool) *types.PenaltyResults {
    
    baseReward := GetBaseReward(state, validatorIndex)
    
    results := &types.PenaltyResults{
        InactivityScore: state.Validators[validatorIndex].InactivityScore,
    }
    
    // Calculate penalties for missed attestation components
    if !correctSource {
        results.SourcePenalty = baseReward * config.TIMELY_SOURCE_WEIGHT / config.WEIGHT_DENOMINATOR
    }
    if !correctTarget {
        results.TargetPenalty = baseReward * config.TIMELY_TARGET_WEIGHT / config.WEIGHT_DENOMINATOR
    }
    if !correctHead {
        results.HeadPenalty = baseReward * config.TIMELY_HEAD_WEIGHT / config.WEIGHT_DENOMINATOR
    }
    
    results.TotalAttestationPenalty = results.SourcePenalty + results.TargetPenalty + results.HeadPenalty
    
    // Calculate inactivity penalty if applicable
    if state.CurrentEpoch > state.FinalizedEpoch+config.MIN_ATTESTATION_INCLUSION_DELAY {
        results.InactivityPenalty = GetInactivityPenalty(state, validatorIndex)
    }
    
    // Daily projections
    results.DailyAttestationPenalty = float64(results.TotalAttestationPenalty*config.EPOCHS_PER_DAY) / 1e9
    results.DailyInactivityPenalty = float64(results.InactivityPenalty*config.EPOCHS_PER_DAY) / 1e9
    
    return results
}

// GetInactivityPenalty calculates the inactivity leak penalty
func GetInactivityPenalty(state *types.NetworkState, validatorIndex int) uint64 {
    validator := &state.Validators[validatorIndex]
    
    // Only applies during non-finality
    if state.CurrentEpoch <= state.FinalizedEpoch+config.MIN_ATTESTATION_INCLUSION_DELAY {
        return 0
    }
    
    // Get appropriate penalty quotient based on fork
    forkConfig := config.GetForkConfig(state.CurrentFork)
    
    penaltyNumerator := validator.EffectiveBalance * validator.InactivityScore
    penaltyDenominator := config.INACTIVITY_SCORE_BIAS * forkConfig.InactivityPenaltyQuotient
    
    return penaltyNumerator / penaltyDenominator
}

// CalculateInactivityScore computes the inactivity score for a validator
func CalculateInactivityScore(previousScore uint64, isActive bool, isFinalized bool) uint64 {
    if isFinalized {
        if previousScore > 0 {
            // Decrease score during finality
            return previousScore - min(1, previousScore)
        }
        return 0
    }
    
    // Increase score during non-finality
    if !isActive {
        return previousScore + config.INACTIVITY_SCORE_BIAS
    }
    
    // Active but not finalizing
    return previousScore + 1
}

// CalculateSlashingPenalties computes all slashing-related penalties
func CalculateSlashingPenalties(state *types.NetworkState, validatorIndex int, 
    totalSlashedBalance uint64) *types.SlashingResults {
    
    validator := &state.Validators[validatorIndex]
    forkConfig := config.GetForkConfig(state.CurrentFork)
    
    // Initial penalty
    initialPenalty := validator.EffectiveBalance / forkConfig.MinSlashingPenaltyQuotient
    
    // Proportional penalty (correlation penalty)
    proportionalPenalty := validator.EffectiveBalance * 
                          min(totalSlashedBalance*forkConfig.ProportionalSlashingMultiplier, 
                              state.TotalActiveBalance) / 
                          state.TotalActiveBalance
    
    totalPenalty := initialPenalty + proportionalPenalty
    
    // Whistleblower rewards
    whistleblowerReward := validator.EffectiveBalance / config.WHISTLEBLOWER_REWARD_QUOTIENT
    proposerReward := whistleblowerReward / config.PROPOSER_REWARD_QUOTIENT
    
    return &types.SlashingResults{
        InitialPenalty:      initialPenalty,
        ProportionalPenalty: proportionalPenalty,
        TotalPenalty:        totalPenalty,
        PercentageOfStake:   float64(totalPenalty) / float64(validator.EffectiveBalance) * 100,
        WhistleblowerReward: whistleblowerReward,
        ProposerReward:      proposerReward,
    }
}

// EstimateSlashingImpact estimates the impact of a slashing event on the network
func EstimateSlashingImpact(state *types.NetworkState, slashedValidatorCount int) map[string]interface{} {
    slashedBalance := uint64(slashedValidatorCount) * config.MAX_EFFECTIVE_BALANCE
    slashingPercentage := float64(slashedBalance) / float64(state.TotalActiveBalance) * 100
    
    // Calculate penalties for different scenarios
    singleSlashing := CalculateSlashingPenalties(state, 0, config.MAX_EFFECTIVE_BALANCE)
    correlatedSlashing := CalculateSlashingPenalties(state, 0, slashedBalance)
    
    return map[string]interface{}{
        "slashed_validator_count": slashedValidatorCount,
        "slashed_balance_eth":     float64(slashedBalance) / 1e9,
        "network_percentage":      slashingPercentage,
        "single_validator_penalty": map[string]interface{}{
            "initial_eth":      float64(singleSlashing.InitialPenalty) / 1e9,
            "proportional_eth": float64(singleSlashing.ProportionalPenalty) / 1e9,
            "total_eth":        float64(singleSlashing.TotalPenalty) / 1e9,
            "percentage":       singleSlashing.PercentageOfStake,
        },
        "correlated_penalty": map[string]interface{}{
            "initial_eth":      float64(correlatedSlashing.InitialPenalty) / 1e9,
            "proportional_eth": float64(correlatedSlashing.ProportionalPenalty) / 1e9,
            "total_eth":        float64(correlatedSlashing.TotalPenalty) / 1e9,
            "percentage":       correlatedSlashing.PercentageOfStake,
        },
        "network_impact": map[string]interface{}{
            "total_penalties_eth":  float64(correlatedSlashing.TotalPenalty*uint64(slashedValidatorCount)) / 1e9,
            "reduced_staking_eth":  float64(slashedBalance) / 1e9,
            "security_impact":      getSecurityImpactLevel(slashingPercentage),
        },
    }
}

// Helper function to determine security impact level
func getSecurityImpactLevel(slashingPercentage float64) string {
    switch {
    case slashingPercentage < 0.1:
        return "Minimal"
    case slashingPercentage < 1.0:
        return "Low"
    case slashingPercentage < 5.0:
        return "Moderate"
    case slashingPercentage < 10.0:
        return "High"
    case slashingPercentage < 33.3:
        return "Critical"
    default:
        return "Catastrophic"
    }
}

// min returns the minimum of two uint64 values
func min(a, b uint64) uint64 {
    if a < b {
        return a
    }
    return b
}