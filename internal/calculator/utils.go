package calculator

import (
    "fmt"
    "math"
    
    "github.com/eth-rewards-calculator/internal/config"
    "github.com/eth-rewards-calculator/internal/types"
)

// ValidatorSetComparison compares rewards across different validator set sizes
func ValidatorSetComparison(participation float64, validatorCounts ...int) []types.ComparisonResult {
    results := make([]types.ComparisonResult, len(validatorCounts))
    
    for i, count := range validatorCounts {
        state := &types.NetworkState{
            Validators:         make([]types.Validator, count),
            TotalActiveBalance: uint64(count) * config.MAX_EFFECTIVE_BALANCE,
            CurrentEpoch:       1000,
            FinalizedEpoch:     998,
        }
        
        // Initialize validators
        for j := range state.Validators {
            state.Validators[j] = types.Validator{
                EffectiveBalance: config.MAX_EFFECTIVE_BALANCE,
            }
        }
        
        rewards := CalculateRewards(state, participation)
        
        results[i] = types.ComparisonResult{
            ValidatorCount: count,
            TotalStaked:    state.TotalActiveBalance / 1e9,
            BaseReward:     rewards.BaseRewardPerEpoch,
            AnnualRewards:  rewards.TotalAnnualRewards / 1e9,
            APY:            rewards.APY,
            DailyRewards:   rewards.DailyRewards / 1e9,
        }
    }
    
    return results
}

// CalculateBreakEvenTime calculates how long until rewards cover initial stake
func CalculateBreakEvenTime(apy float64) (years, months, days float64) {
    if apy <= 0 {
        return math.Inf(1), math.Inf(1), math.Inf(1)
    }
    
    // Time to double investment (100% return)
    yearsToDouble := 100.0 / apy
    years = yearsToDouble
    months = yearsToDouble * 12
    days = yearsToDouble * 365.25
    
    return
}

// EstimateValidatorQueue estimates activation queue time based on churn limit
func EstimateValidatorQueue(currentValidators, pendingValidators int) (epochs, days float64) {
    // Calculate churn limit
    churnLimit := max(config.MIN_PER_EPOCH_CHURN_LIMIT, 
                     uint64(currentValidators)/config.CHURN_LIMIT_QUOTIENT)
    
    // Cap at max churn limit
    churnLimit = min(churnLimit, config.MAX_PER_EPOCH_ACTIVATION_CHURN_LIMIT)
    
    epochs = float64(pendingValidators) / float64(churnLimit)
    days = epochs / float64(config.EPOCHS_PER_DAY)
    
    return
}

// CalculateCompoundingReturns calculates returns with reinvestment
func CalculateCompoundingReturns(initialStake float64, apy float64, years int) map[string]float64 {
    results := make(map[string]float64)
    
    // Convert APY to decimal
    rate := apy / 100.0
    
    // Calculate for each year
    for year := 1; year <= years; year++ {
        value := initialStake * math.Pow(1+rate, float64(year))
        results[fmt.Sprintf("year_%d", year)] = value
    }
    
    // Calculate total return
    finalValue := initialStake * math.Pow(1+rate, float64(years))
    results["total_return"] = finalValue - initialStake
    results["total_return_percentage"] = ((finalValue - initialStake) / initialStake) * 100
    
    return results
}

// OptimalValidatorDistribution suggests optimal validator distribution for a given ETH amount
func OptimalValidatorDistribution(totalETH float64) map[string]interface{} {
    validatorCount := int(totalETH / 32.0)
    remainingETH := math.Mod(totalETH, 32.0)
    
    distribution := map[string]interface{}{
        "total_eth":           totalETH,
        "full_validators":     validatorCount,
        "staked_eth":         float64(validatorCount) * 32.0,
        "remaining_eth":      remainingETH,
        "efficiency":         (float64(validatorCount) * 32.0 / totalETH) * 100,
    }
    
    // Add recommendation
    if remainingETH >= 16 {
        distribution["recommendation"] = "Consider waiting to accumulate 32 ETH for another validator"
    } else if remainingETH > 0 {
        distribution["recommendation"] = fmt.Sprintf("Keep %.2f ETH liquid or in DeFi", remainingETH)
    } else {
        distribution["recommendation"] = "Optimal distribution achieved"
    }
    
    return distribution
}

// CalculateNetReturns calculates returns after considering various factors
func CalculateNetReturns(grossAPY, inflationRate, taxRate float64) map[string]float64 {
    realReturn := grossAPY - inflationRate
    afterTaxReturn := grossAPY * (1 - taxRate/100)
    realAfterTaxReturn := afterTaxReturn - inflationRate
    
    return map[string]float64{
        "gross_apy":            grossAPY,
        "inflation_adjusted":   realReturn,
        "after_tax":           afterTaxReturn,
        "real_after_tax":      realAfterTaxReturn,
        "effective_rate":      realAfterTaxReturn,
    }
}

// Helper functions

func max(a, b uint64) uint64 {
    if a > b {
        return a
    }
    return b
}

func maxFloat(a, b float64) float64 {
    if a > b {
        return a
    }
    return b
}

func minFloat(a, b float64) float64 {
    if a < b {
        return a
    }
    return b
}

// FormatGwei formats Gwei values for display
func FormatGwei(gwei uint64) string {
    if gwei >= 1e9 {
        return fmt.Sprintf("%.6f ETH", float64(gwei)/1e9)
    } else if gwei >= 1e6 {
        return fmt.Sprintf("%.3f mETH", float64(gwei)/1e6)
    }
    return fmt.Sprintf("%d Gwei", gwei)
}

// FormatPercentage formats percentage with appropriate precision
func FormatPercentage(value float64) string {
    if value < 0.01 {
        return fmt.Sprintf("%.4f%%", value)
    } else if value < 1 {
        return fmt.Sprintf("%.3f%%", value)
    } else if value < 10 {
        return fmt.Sprintf("%.2f%%", value)
    }
    return fmt.Sprintf("%.1f%%", value)
}