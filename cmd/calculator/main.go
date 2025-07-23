package main

import (
    "encoding/json"
    "fmt"
    "os"
    "strconv"
    "strings"

    "github.com/eth-rewards-calculator/internal/calculator"
    "github.com/eth-rewards-calculator/internal/config"
    "github.com/eth-rewards-calculator/internal/types"

    "github.com/fatih/color"
    flag "github.com/spf13/pflag"
)

var (
    validatorCount   int
    participation    float64
    detailed         bool
    jsonOutput       bool
    compare          string
    showPenalties    bool
    inactivityEpochs int
    slashingCount    int
    compareParticipation bool
)

func init() {
    flag.IntVarP(&validatorCount, "validators", "v", 0, "Number of validators")
    flag.Float64VarP(&participation, "participation", "p", 0.95, "Network participation rate (0.0-1.0)")
    flag.BoolVarP(&detailed, "detailed", "d", false, "Show detailed breakdown")
    flag.BoolVarP(&jsonOutput, "json", "j", false, "Output results as JSON")
    flag.StringVarP(&compare, "compare", "c", "", "Compare multiple validator counts (comma-separated)")
    flag.BoolVarP(&showPenalties, "penalties", "", false, "Show penalty calculations")
    flag.IntVarP(&inactivityEpochs, "inactivity", "i", 0, "Epochs of inactivity for penalty calculation")
    flag.IntVarP(&slashingCount, "slashing", "s", 0, "Number of validators slashed together")
    flag.BoolVarP(&compareParticipation, "compare-participation", "", false, "Compare rewards at different participation rates")
}

func main() {
    flag.Parse()

    // Validate inputs
    if validatorCount == 0 && compare == "" && !compareParticipation {
        fmt.Println("Error: Please specify validator count with -v, use -c for comparison, or use --compare-participation")
        flag.Usage()
        os.Exit(1)
    }

    if participation < 0 || participation > 1 {
        fmt.Println("Error: Participation rate must be between 0.0 and 1.0")
        os.Exit(1)
    }

    // Handle comparison mode
    if compare != "" {
        handleComparison(compare, participation)
        return
    }
    
    // Handle participation comparison mode
    if compareParticipation {
        if validatorCount == 0 {
            validatorCount = 10000 // Default for participation comparison
        }
        compareParticipationRates(validatorCount)
        return
    }

    // Single validator count calculation
    state := createNetworkState(validatorCount)
    results := calculator.CalculateRewards(state, participation)

    if jsonOutput {
        outputJSON(results)
    } else {
        outputFormatted(results, state, detailed)
    }

    if showPenalties {
        showPenaltyExamples(state)
    }
}

func createNetworkState(validators int) *types.NetworkState {
    state := &types.NetworkState{
        Validators:         make([]types.Validator, validators),
        TotalActiveBalance: uint64(validators) * config.MAX_EFFECTIVE_BALANCE,
        CurrentEpoch:       1000,
        FinalizedEpoch:     998,
    }

    // Initialize validators
    for i := range state.Validators {
        state.Validators[i] = types.Validator{
            EffectiveBalance: config.MAX_EFFECTIVE_BALANCE,
            Slashed:          false,
            InactivityScore:  0,
        }
        
        if inactivityEpochs > 0 {
            state.Validators[i].InactivityScore = uint64(inactivityEpochs * 4)
            state.FinalizedEpoch = state.CurrentEpoch - uint64(inactivityEpochs) - 2
        }
    }

    return state
}

func handleComparison(compareStr string, participation float64) {
    counts := strings.Split(compareStr, ",")
    
    header := color.New(color.FgCyan, color.Bold)
    header.Println("\n=== Ethereum Staking Rewards Comparison ===")
    
    fmt.Printf("\nParticipation Rate: %.1f%%\n\n", participation*100)
    
    // Table header
    fmt.Printf("%-15s %-20s %-20s %-15s %-10s %-15s\n", 
        "Validators", "Total Staked (ETH)", "Base Reward (Gwei)", 
        "Annual ETH", "APY %", "Daily ETH")
    fmt.Println(strings.Repeat("-", 100))

    for _, countStr := range counts {
        count, err := strconv.Atoi(strings.TrimSpace(countStr))
        if err != nil {
            fmt.Printf("Error: Invalid validator count '%s'\n", countStr)
            continue
        }

        state := createNetworkState(count)
        results := calculator.CalculateRewards(state, participation)
        
        fmt.Printf("%-15d %-20s %-20d %-15.6f %-10.2f%% %-15.6f\n",
            count,
            formatNumber(state.TotalActiveBalance/1e9),
            results.BaseRewardPerEpoch,
            results.TotalAnnualRewards/1e9,
            results.APY,
            results.TotalAnnualRewards/1e9/365.25)
    }
    
    fmt.Println()
}

func compareParticipationRates(validatorCount int) {
    header := color.New(color.FgCyan, color.Bold)
    header.Println("\n=== Participation Rate Impact Analysis ===")
    
    fmt.Printf("\nValidator Count: %s\n\n", formatNumber(uint64(validatorCount)))
    
    // Create network state once
    state := createNetworkState(validatorCount)
    
    // Table header
    fmt.Printf("%-20s %-15s %-15s %-20s %-15s %-25s\n", 
        "Participation Rate", "Multiplier", "Base APY %", "Effective APY %", 
        "Annual ETH", "Network Status")
    fmt.Println(strings.Repeat("-", 110))
    
    // Compare different participation rates
    participationRates := []float64{1.0, 0.95, 0.9, 0.8, 0.7, 0.6667, 0.6, 0.5, 0.4, 0.3333}
    
    for _, rate := range participationRates {
        results := calculator.CalculateRewards(state, rate)
        
        statusColor := color.New(color.FgGreen)
        status := "Healthy"
        
        if rate < 0.3333 {
            statusColor = color.New(color.FgRed, color.Bold)
            status = "CRITICAL - No finality"
        } else if rate < 0.6667 {
            statusColor = color.New(color.FgRed)
            status = "Inactivity leak active"
        } else if rate < 0.8 {
            statusColor = color.New(color.FgYellow)
            status = "Reduced security"
        }
        
        fmt.Printf("%-20s %-15s %-15.2f%% %-20.2f%% %-15.6f ",
            fmt.Sprintf("%.1f%%", rate*100),
            fmt.Sprintf("%.2fx", results.ParticipationMultiplier),
            results.BaseAPY,
            results.EffectiveAPY,
            results.TotalAnnualRewards/1e9)
        
        statusColor.Printf("%-25s\n", status)
    }
    
    fmt.Println("\nNOTE: This model shows how active validators benefit from others being offline.")
    fmt.Println("      At low participation rates, inactivity penalties and network instability become significant factors.")
}

func outputFormatted(results *types.RewardResults, state *types.NetworkState, detailed bool) {
    header := color.New(color.FgCyan, color.Bold)
    subheader := color.New(color.FgYellow, color.Bold)
    highlight := color.New(color.FgGreen, color.Bold)
    
    header.Println("\n=== Ethereum Staking Rewards Calculator ===")
    
    // Network Parameters
    subheader.Println("\nNetwork Parameters:")
    fmt.Printf("- Validator Count: %s\n", formatNumber(uint64(len(state.Validators))))
    fmt.Printf("- Total Staked: %s ETH\n", formatNumber(state.TotalActiveBalance/1e9))
    fmt.Printf("- Participation Rate: %.1f%%\n", results.ParticipationRate*100)
    fmt.Printf("- Effective Balance: %.0f ETH\n", float64(config.MAX_EFFECTIVE_BALANCE)/1e9)
    
    // Base Reward Calculation
    subheader.Println("\nBase Reward Calculation:")
    fmt.Printf("- Base Reward Factor: %d\n", config.BASE_REWARD_FACTOR)
    fmt.Printf("- Square Root of Total Balance: %s\n", formatNumber(results.SqrtTotalBalance))
    fmt.Printf("- Base Reward per Epoch: %s Gwei (%.9f ETH)\n", 
        formatNumber(results.BaseRewardPerEpoch), float64(results.BaseRewardPerEpoch)/1e9)
    
    if detailed {
        // Detailed Reward Breakdown
        subheader.Println("\nDetailed Reward Breakdown (per epoch):")
        fmt.Printf("- Source Vote Reward: %s Gwei (%.2f%%)\n", 
            formatNumber(results.SourceReward), 
            float64(config.TIMELY_SOURCE_WEIGHT)/float64(config.WEIGHT_DENOMINATOR)*100)
        fmt.Printf("- Target Vote Reward: %s Gwei (%.2f%%)\n", 
            formatNumber(results.TargetReward),
            float64(config.TIMELY_TARGET_WEIGHT)/float64(config.WEIGHT_DENOMINATOR)*100)
        fmt.Printf("- Head Vote Reward: %s Gwei (%.2f%%)\n", 
            formatNumber(results.HeadReward),
            float64(config.TIMELY_HEAD_WEIGHT)/float64(config.WEIGHT_DENOMINATOR)*100)
        fmt.Printf("- Total Attestation Reward: %s Gwei\n", 
            formatNumber(results.AttestationRewardPerEpoch))
        
        subheader.Println("\nProposer Statistics:")
        fmt.Printf("- Probability per Epoch: %.4f%%\n", results.ProposerProbability*100)
        fmt.Printf("- Expected Proposals per Year: %.2f\n", results.ExpectedProposalsPerYear)
        fmt.Printf("- Average Proposer Reward per Block: %s Gwei\n", 
            formatNumber(uint64(results.AvgProposerRewardPerBlock)))
    }
    
    // Participation Economics
    if results.ParticipationRate < 1.0 {
        subheader.Println("\nParticipation Economics:")
        fmt.Printf("- Participation Multiplier: %.2fx\n", results.ParticipationMultiplier)
        fmt.Printf("- Base APY (at 100%% participation): %.2f%%\n", results.BaseAPY)
        fmt.Printf("- Effective APY (with boost): %.2f%%\n", results.EffectiveAPY)
        if results.NetworkHealthWarning != "" {
            warningColor := color.New(color.FgRed, color.Bold)
            warningColor.Printf("- %s\n", results.NetworkHealthWarning)
        }
    }
    
    // Annual Rewards
    subheader.Println("\nAnnual Rewards:")
    fmt.Printf("- Attestation Rewards: %.6f ETH\n", results.AttestationRewardsAnnual/1e9)
    fmt.Printf("- Proposer Rewards: %.6f ETH\n", results.ProposerRewardsAnnual/1e9)
    fmt.Printf("- Total Annual Rewards: %.6f ETH\n", results.TotalAnnualRewards/1e9)
    
    highlight.Printf("- Annual Percentage Yield (APY): %.2f%%\n", results.APY)
    
    // Daily/Monthly projections
    subheader.Println("\nProjected Earnings:")
    fmt.Printf("- Daily: %.6f ETH\n", results.TotalAnnualRewards/1e9/365.25)
    fmt.Printf("- Weekly: %.6f ETH\n", results.TotalAnnualRewards/1e9/52.18)
    fmt.Printf("- Monthly: %.6f ETH\n", results.TotalAnnualRewards/1e9/12)
}

func showPenaltyExamples(state *types.NetworkState) {
    header := color.New(color.FgRed, color.Bold)
    subheader := color.New(color.FgYellow, color.Bold)
    
    header.Println("\n=== Penalty Examples ===")
    
    validatorIndex := 0
    
    // Missed attestation
    penalties := calculator.CalculatePenalties(state, validatorIndex, false, false, false)
    subheader.Println("\nMissed Attestation Penalties:")
    fmt.Printf("- Source Penalty: %s Gwei\n", formatNumber(penalties.SourcePenalty))
    fmt.Printf("- Target Penalty: %s Gwei\n", formatNumber(penalties.TargetPenalty))
    fmt.Printf("- Head Penalty: %s Gwei\n", formatNumber(penalties.HeadPenalty))
    fmt.Printf("- Total per Epoch: %s Gwei\n", formatNumber(penalties.TotalAttestationPenalty))
    fmt.Printf("- Daily Cost: %.6f ETH\n", float64(penalties.TotalAttestationPenalty*225)/1e9)
    
    // Inactivity leak
    if inactivityEpochs > 0 {
        inactivityPenalty := calculator.GetInactivityPenalty(state, validatorIndex)
        subheader.Printf("\nInactivity Leak (%d epochs without finality):\n", inactivityEpochs)
        fmt.Printf("- Inactivity Score: %d\n", state.Validators[validatorIndex].InactivityScore)
        fmt.Printf("- Penalty per Epoch: %s Gwei (%.6f ETH)\n", 
            formatNumber(inactivityPenalty), float64(inactivityPenalty)/1e9)
        fmt.Printf("- Daily Penalty: %.6f ETH\n", float64(inactivityPenalty*225)/1e9)
        fmt.Printf("- Projected Loss in 30 days: %.6f ETH\n", float64(inactivityPenalty*225*30)/1e9)
    }
    
    // Slashing
    if slashingCount > 0 {
        subheader.Printf("\nSlashing Penalties (%d validators slashed together):\n", slashingCount)
        slashingResults := calculator.CalculateSlashingPenalties(
            state, validatorIndex, uint64(slashingCount)*config.MAX_EFFECTIVE_BALANCE)
        
        fmt.Printf("- Initial Penalty: %.6f ETH\n", float64(slashingResults.InitialPenalty)/1e9)
        fmt.Printf("- Proportional Penalty: %.6f ETH\n", float64(slashingResults.ProportionalPenalty)/1e9)
        fmt.Printf("- Total Penalty: %.6f ETH (%.2f%% of stake)\n", 
            float64(slashingResults.TotalPenalty)/1e9,
            float64(slashingResults.TotalPenalty)/float64(config.MAX_EFFECTIVE_BALANCE)*100)
    }
}

func outputJSON(results *types.RewardResults) {
    output, err := json.MarshalIndent(results, "", "  ")
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error marshaling JSON: %v\n", err)
        os.Exit(1)
    }
    fmt.Println(string(output))
}

func formatNumber(n uint64) string {
    str := strconv.FormatUint(n, 10)
    var result []string
    
    for i, digit := range str {
        if i > 0 && (len(str)-i)%3 == 0 {
            result = append(result, ",")
        }
        result = append(result, string(digit))
    }
    
    return strings.Join(result, "")
}