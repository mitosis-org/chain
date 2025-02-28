package types

import (
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// NewValidator creates a new validator
func NewValidator(pubkey []byte, collateral, extraVotingPower math.Int) Validator {
	return Validator{
		Pubkey:           pubkey,
		Collateral:       collateral,
		ExtraVotingPower: extraVotingPower,
		VotingPower:      math.ZeroInt(), // will be computed on finalizing
		Jailed:           false,
	}
}

func (v Validator) ConsAddr() (sdk.ConsAddress, error) {
	pubkey, err := v.ConsPubKey()
	if err != nil {
		return nil, err
	}
	return sdk.GetConsAddress(pubkey), nil
}
