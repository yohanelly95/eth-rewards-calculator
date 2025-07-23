# Ethereum Rewards Calculator

A command-line tool for calculating Ethereum validator rewards, penalties, and analyzing network economics with accurate participation rate modeling.

## Features

- **Accurate Participation Model**: Implements Ethereum's actual reward distribution where active validators earn more when participation is low
- **Comprehensive Reward Calculations**: Base rewards, attestation rewards, proposer rewards, and APY
- **Penalty Analysis**: Missed attestations, inactivity leaks, and slashing penalties
- **Comparison Modes**: Compare different validator counts or participation rates
- **Network Health Monitoring**: Warnings for low participation and security risks
- **JSON Output**: Machine-readable output for integration with other tools

## Installation

### From Source

```bash
# Clone the repository
git clone https://github.com/yourusername/eth-rewards-calculator.git
cd eth-rewards-calculator

# Build the binary
make build

# Optional: Install to /usr/local/bin
make install
```

### Quick Start

```bash
# Run the quickstart script to set up the project
./quickstart.sh
```

## Usage

### Basic Usage

Calculate rewards for a specific number of validators:

```bash
# Calculate rewards for 1,000 validators
./bin/eth-rewards -v 1000

# Calculate rewards for 10,000 validators with 90% participation
./bin/eth-rewards -v 10000 -p 0.9
```

### Command-Line Flags

| Flag | Short | Description | Default |
|------|-------|-------------|---------|
| `--validators` | `-v` | Number of validators to simulate | Required* |
| `--participation` | `-p` | Network participation rate (0.0-1.0) | 0.95 |
| `--detailed` | `-d` | Show detailed breakdown of rewards | false |
| `--json` | `-j` | Output results as JSON | false |
| `--compare` | `-c` | Compare multiple validator counts (comma-separated) | - |
| `--compare-participation` | | Compare rewards at different participation rates | false |
| `--penalties` | | Show penalty calculation examples | false |
| `--inactivity` | `-i` | Epochs of inactivity for penalty calculation | 0 |
| `--slashing` | `-s` | Number of validators slashed together | 0 |

*Required unless using `--compare` or `--compare-participation`

### Examples

#### 1. Basic Reward Calculation

```bash
./bin/eth-rewards -v 4096
```

Output shows:
- Network parameters (validator count, total staked ETH)
- Base reward calculations
- Annual rewards breakdown
- APY (Annual Percentage Yield)
- Daily/weekly/monthly projections

#### 2. Detailed Breakdown

```bash
./bin/eth-rewards -v 4096 -d
```

Additional output includes:
- Component rewards (source, target, head votes)
- Proposer statistics and probabilities
- Detailed reward percentages

#### 3. Participation Rate Analysis

```bash
# See how rewards change with 80% participation
./bin/eth-rewards -v 10000 -p 0.8

# Compare different participation rates
./bin/eth-rewards --compare-participation -v 10000
```

The participation comparison shows:
- Reward multiplier for active validators
- Base APY vs Effective APY
- Network health status
- Warnings for low participation

#### 4. Comparing Validator Counts

```bash
./bin/eth-rewards -c 1000,10000,100000,500000,1000000
```

Displays a table comparing:
- Total staked ETH
- Base rewards
- Annual returns
- APY for each validator count

#### 5. Penalty Calculations

```bash
# Show penalty examples
./bin/eth-rewards -v 4096 --penalties

# Include inactivity penalties (10 epochs without finality)
./bin/eth-rewards -v 4096 --penalties -i 10

# Include slashing penalties (100 validators slashed together)
./bin/eth-rewards -v 4096 --penalties -s 100
```

#### 6. JSON Output

```bash
./bin/eth-rewards -v 4096 -j > rewards.json
```

Useful for:
- Scripting and automation
- Data analysis
- Integration with other tools

### Understanding Participation Economics

The calculator implements Ethereum's actual reward distribution model:

- **100% Participation**: Baseline rewards (no multiplier)
- **95% Participation**: 1.05x multiplier for active validators
- **50% Participation**: 2.0x multiplier for active validators
- **<66.67% Participation**: Inactivity leak becomes active
- **<33.33% Participation**: Chain cannot finalize (critical)

Example output at 50% participation:
```
Participation Economics:
- Participation Multiplier: 2.00x
- Base APY (at 100% participation): 11.77%
- Effective APY (with boost): 23.54%
- WARNING: Network participation below 66.67% - inactivity leak active
```

### Network Health Warnings

The calculator provides warnings based on participation rate:

- **80-100%**: Healthy network
- **66.67-80%**: "CAUTION: Network participation below 80% - reduced security"
- **33.33-66.67%**: "WARNING: Network participation below 66.67% - inactivity leak active"
- **<33.33%**: "CRITICAL: Network participation below 33.33% - chain cannot finalize"

## Advanced Features

### Inactivity Leak Simulation

Simulate validator behavior during network non-finality:

```bash
./bin/eth-rewards -v 4096 -i 100 --penalties
```

Shows:
- Inactivity score accumulation
- Daily penalty rates
- Projected losses over time

### Slashing Analysis

Analyze the impact of correlated slashing events:

```bash
./bin/eth-rewards -v 4096 -s 50 --penalties
```

Calculates:
- Initial slashing penalty
- Proportional penalty based on total slashed
- Total ETH lost
- Percentage of stake lost

## Understanding the Output

### Key Metrics Explained

1. **Base Reward**: Fundamental unit of rewards, calculated as:
   ```
   effective_balance × base_reward_factor / sqrt(total_active_balance) / 4
   ```

2. **APY (Annual Percentage Yield)**: Expected annual return as a percentage of staked ETH

3. **Participation Multiplier**: Boost factor for active validators when others are offline

4. **Component Rewards**:
   - Source vote: ~21.875% of base reward
   - Target vote: ~40.625% of base reward
   - Head vote: ~21.875% of base reward
   - Proposer: ~12.5% of base reward

## Build Options

```bash
# Standard build
make build

# Build for multiple platforms
make build-all

# Run tests
make test

# Run with coverage
make test-coverage

# Clean build artifacts
make clean
```

## Development

### Project Structure

```
eth-rewards-calculator/
├── cmd/calculator/      # Main application entry point
├── internal/
│   ├── calculator/      # Core calculation logic
│   ├── config/          # Configuration constants
│   └── types/           # Data structures
├── bin/                 # Compiled binaries
├── Makefile            # Build configuration
└── README.md           # This file
```

### Adding New Features

1. Update types in `internal/types/types.go`
2. Implement calculations in `internal/calculator/`
3. Add command-line flags in `cmd/calculator/main.go`
4. Update this README

## Notes

- All calculations use Gwei (1 ETH = 1e9 Gwei) internally
- Default participation rate is 95% to reflect typical network conditions
- The calculator assumes all validators have 32 ETH effective balance
- Actual rewards may vary based on network conditions and validator performance