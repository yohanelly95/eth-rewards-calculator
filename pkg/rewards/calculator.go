package rewards

type Calculator struct {
	BaseRewardFactor    uint64
	TotalActiveBalance  uint64
}

func NewCalculator(baseRewardFactor, totalActiveBalance uint64) *Calculator {
	if baseRewardFactor == 0 {
		baseRewardFactor = DEFAULT_BASE_REWARD_FACTOR
	}
	if totalActiveBalance == 0 {
		totalActiveBalance = DEFAULT_TOTAL_ACTIVE_BALANCE
	}
	
	return &Calculator{
		BaseRewardFactor:   baseRewardFactor,
		TotalActiveBalance: totalActiveBalance,
	}
}

// CalculateBaseReward returns the base reward for a validator with given effective balance
func (c *Calculator) CalculateBaseReward(effectiveBalance uint64) uint64 {
	if effectiveBalance > MAX_EFFECTIVE_BALANCE {
		effectiveBalance = MAX_EFFECTIVE_BALANCE
	}
	
	return effectiveBalance * c.BaseRewardFactor / IntegerSquareRoot(c.TotalActiveBalance)
}

// CalculateMaxAttestationReward returns the maximum reward for perfect attestations
func (c *Calculator) CalculateMaxAttestationReward(effectiveBalance uint64) uint64 {
	baseReward := c.CalculateBaseReward(effectiveBalance)
	
	sourceReward := baseReward * TIMELY_SOURCE_WEIGHT / WEIGHT_DENOMINATOR
	targetReward := baseReward * TIMELY_TARGET_WEIGHT / WEIGHT_DENOMINATOR
	headReward := baseReward * TIMELY_HEAD_WEIGHT / WEIGHT_DENOMINATOR
	
	return sourceReward + targetReward + headReward
}

// CalculateProposerReward returns the reward for proposing a block
func (c *Calculator) CalculateProposerReward(effectiveBalance uint64) uint64 {
	baseReward := c.CalculateBaseReward(effectiveBalance)
	return baseReward * PROPOSER_WEIGHT / WEIGHT_DENOMINATOR
}

// CalculateSyncCommitteeReward returns the reward for sync committee participation
func (c *Calculator) CalculateSyncCommitteeReward(effectiveBalance uint64) uint64 {
	baseReward := c.CalculateBaseReward(effectiveBalance)
	return baseReward * SYNC_REWARD_WEIGHT / WEIGHT_DENOMINATOR
}

// CalculateAnnualReward estimates annual rewards for a validator with perfect performance
func (c *Calculator) CalculateAnnualReward(effectiveBalance uint64, proposerProbability float64) uint64 {
	// Attestation rewards per epoch
	attestationRewardPerEpoch := c.CalculateMaxAttestationReward(effectiveBalance)
	annualAttestationReward := attestationRewardPerEpoch * EPOCHS_PER_YEAR
	
	// Expected proposer rewards
	proposerRewardPerBlock := c.CalculateProposerReward(effectiveBalance)
	expectedBlocksPerYear := uint64(float64(EPOCHS_PER_YEAR*SLOTS_PER_EPOCH) * proposerProbability)
	annualProposerReward := proposerRewardPerBlock * expectedBlocksPerYear
	
	return annualAttestationReward + annualProposerReward
}

// CalculateAPR calculates the Annual Percentage Rate
func (c *Calculator) CalculateAPR(effectiveBalance uint64, proposerProbability float64) float64 {
	annualReward := c.CalculateAnnualReward(effectiveBalance, proposerProbability)
	return float64(annualReward) / float64(effectiveBalance) * 100
}