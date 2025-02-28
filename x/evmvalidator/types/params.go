package types

import (
	"fmt"

	"cosmossdk.io/math"
)

// DefaultMaxValidators is the default maximum number of validators.
const DefaultMaxValidators uint32 = 100

// DefaultMaxLeverageRatio is the default maximum leverage ratio.
var DefaultMaxLeverageRatio = math.LegacyNewDec(2) // 2.0x leverage

// DefaultMinVotingPower is the default minimum voting power required.
var DefaultMinVotingPower = math.NewInt(1000000) // 1M minimum voting power

// DefaultWithdrawalLimit is the default withdrawal limit per block.
const DefaultWithdrawalLimit uint32 = 10

// DefaultParams returns a default set of parameters.
func DefaultParams() Params {
	return Params{
		MaxValidators:    DefaultMaxValidators,
		MaxLeverageRatio: DefaultMaxLeverageRatio,
		MinVotingPower:   DefaultMinVotingPower,
		WithdrawalLimit:  DefaultWithdrawalLimit,
	}
}

// Validate validates the set of params.
func (p Params) Validate() error {
	if p.MaxValidators == 0 {
		return fmt.Errorf("max validators must be positive: %d", p.MaxValidators)
	}
	if p.MaxLeverageRatio.IsNil() || p.MaxLeverageRatio.LT(math.LegacyNewDec(1)) {
		return fmt.Errorf("max leverage ratio must be at least 1: %s", p.MaxLeverageRatio)
	}
	if p.MinVotingPower.IsNil() || p.MinVotingPower.IsNegative() {
		return fmt.Errorf("min voting power must be non-negative: %s", p.MinVotingPower)
	}
	if p.WithdrawalLimit == 0 {
		return fmt.Errorf("withdrawal limit must be positive: %d", p.WithdrawalLimit)
	}
	return nil
}
