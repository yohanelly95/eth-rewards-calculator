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
    
    // Calculate realistic proposer reward including attestation inclusion
    attestationInclusionReward := CalculateAttestationInclusionReward(state, participationRate)
    estimatedAttestationsPerBlock := EstimateAttestationsPerBlock(state)
    inclusionEffectivenessRate := CalculateInclusionEffectivenessRate(participationRate)
    
    // Average proposer reward per block (with attestation inclusion)
    avgProposerReward := float64(attestationInclusionReward)
    proposerRewardPerEpoch := avgProposerReward * proposerProbability
    
    // Calculate base annual rewards (at 100% participation)
    baseAttestationAnnual := float64(attestationReward) * float64(config.EPOCHS_PER_YEAR)
    baseProposerAnnual := proposerRewardPerEpoch * float64(config.EPOCHS_PER_YEAR)
    baseTotalAnnual := baseAttestationAnnual + baseProposerAnnual
    baseAPY := (baseTotalAnnual / float64(config.MAX_EFFECTIVE_BALANCE)) * 100
    
    // Apply participation economics - active validators get higher rewards when participation is low
    participationMultiplier := 1.0 / participationRate
    
    // Effective rewards for active validators
    attestationAnnual := baseAttestationAnnual * participationMultiplier
    proposerAnnual := baseProposerAnnual * participationMultiplier
    totalAnnual := attestationAnnual + proposerAnnual
    
    // Effective APY with participation boost
    effectiveAPY := (totalAnnual / float64(config.MAX_EFFECTIVE_BALANCE)) * 100
    
    // Check for inactivity leak conditions
    inactivityLeakActive := participationRate < 0.6667
    networkHealthWarning := ""
    if participationRate < 0.3333 {
        networkHealthWarning = "CRITICAL: Network participation below 33.33% - chain cannot finalize"
    } else if participationRate < 0.6667 {
        networkHealthWarning = "WARNING: Network participation below 66.67% - inactivity leak active"
    } else if participationRate < 0.8 {
        networkHealthWarning = "CAUTION: Network participation below 80% - reduced security"
    }
    
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
        
        // Attestation inclusion details
        EstimatedAttestationsPerBlock: estimatedAttestationsPerBlock,
        AttestationInclusionReward:    attestationInclusionReward,
        InclusionEffectivenessRate:    inclusionEffectivenessRate,
        
        // Annual projections
        AttestationRewardsAnnual: attestationAnnual,
        ProposerRewardsAnnual:    proposerAnnual,
        TotalAnnualRewards:       totalAnnual,
        APY:                      effectiveAPY,
        
        // Time-based projections
        DailyRewards:   totalAnnual / 365.25,
        WeeklyRewards:  totalAnnual / 52.18,
        MonthlyRewards: totalAnnual / 12,
        
        // Participation economics
        ParticipationMultiplier: participationMultiplier,
        BaseAPY:                baseAPY,
        EffectiveAPY:           effectiveAPY,
        InactivityLeakActive:   inactivityLeakActive,
        NetworkHealthWarning:   networkHealthWarning,
    }
}

// GetBaseReward calculates the base reward for a validator using Electra formula (Altair+)
func GetBaseReward(state *types.NetworkState, validatorIndex int) uint64 {
    totalBalance := state.TotalActiveBalance
    effectiveBalance := state.Validators[validatorIndex].EffectiveBalance
    
    // Electra formula: removes division by BASE_REWARDS_PER_EPOCH (used in Phase 0)
    return effectiveBalance * config.BASE_REWARD_FACTOR / 
           IntegerSquareRoot(totalBalance)
}

// GetBaseRewardPerIncrement calculates base reward per increment using Electra formula (Altair+)
func GetBaseRewardPerIncrement(state *types.NetworkState) uint64 {
    return config.EFFECTIVE_BALANCE_INCREMENT * config.BASE_REWARD_FACTOR / 
           IntegerSquareRoot(state.TotalActiveBalance)
}

// EstimateAttestationsPerBlock estimates how many attestations can fit in a block
func EstimateAttestationsPerBlock(state *types.NetworkState) float64 {
    validatorCount := float64(len(state.Validators))
    
    // Attestations come from validators in previous epochs
    // Each epoch has 32 slots, so we get attestations from ~32 slots worth of validators
    // But blocks have size limits, so we can't include all attestations
    
    // Conservative estimate: ~60% of validator attestations can be included per block
    // This accounts for:
    // - Block size limits
    // - Some attestations being too old
    // - Network propagation delays
    maxIncludableRate := 0.6
    
    // Attestations per slot = validators / slots_per_epoch
    attestationsPerSlot := validatorCount / float64(config.SLOTS_PER_EPOCH)
    
    // Estimate attestations from multiple previous slots that can be included
    slotsToInclude := 8.0 // Conservative estimate of slots we can include from
    
    estimatedAttestations := attestationsPerSlot * slotsToInclude * maxIncludableRate
    
    return estimatedAttestations
}

// CalculateAttestationInclusionReward calculates rewards for including attestations in a block
func CalculateAttestationInclusionReward(state *types.NetworkState, participationRate float64) uint64 {
    baseRewardIncrement := GetBaseRewardPerIncrement(state)
    estimatedAttestations := EstimateAttestationsPerBlock(state)
    
    // Each attestation has 3 components: source, target, head
    // Proposer gets reward for each component included
    avgComponentsPerAttestation := 2.8 // Account for some missed/late votes
    
    // Apply participation rate - not all validators are active
    effectiveAttestations := estimatedAttestations * participationRate
    
    // Apply inclusion effectiveness - some attestations are late or missed
    inclusionEffectiveness := 0.9 // 90% effectiveness rate
    finalAttestations := effectiveAttestations * inclusionEffectiveness
    
    // Calculate total proposer reward
    // Proposer gets 1/PROPOSER_REWARD_QUOTIENT of the attestation reward
    proposerRewardPerComponent := baseRewardIncrement / config.PROPOSER_REWARD_QUOTIENT
    totalInclusionReward := uint64(finalAttestations * avgComponentsPerAttestation) * proposerRewardPerComponent
    
    return totalInclusionReward
}

// CalculateInclusionEffectivenessRate calculates the effective inclusion rate
func CalculateInclusionEffectivenessRate(participationRate float64) float64 {
    // Base effectiveness of 90% (some attestations are late or missed)
    baseEffectiveness := 0.9
    
    // Lower participation means less competition for inclusion, slightly higher effectiveness
    // But also means more empty slots, so balance these effects
    participationAdjustment := 0.95 + (participationRate-0.95)*0.5
    
    return baseEffectiveness * participationAdjustment
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