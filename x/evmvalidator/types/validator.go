package types

import (
	"bytes"
	"fmt"

	"cosmossdk.io/math"
	abci "github.com/cometbft/cometbft/abci/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/omni-network/omni/lib/k1util"
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

// ConsAddr returns the validator's consensus address
func (v Validator) ConsAddr() (sdk.ConsAddress, error) {
	pubkey, err := v.ConsPubKey()
	if err != nil {
		return nil, err
	}
	return sdk.GetConsAddress(pubkey), nil
}

// GetVotingPower returns the validator's voting power
func (v Validator) GetVotingPower() math.Int {
	return v.VotingPower
}

// ComputeVotingPower calculates voting power based on collateral and extra voting power
// with respect to the max leverage ratio
func (v Validator) ComputeVotingPower(maxLeverageRatio math.LegacyDec) math.Int {
	// If collateral is zero, voting power should be zero
	if v.Collateral.IsZero() {
		return math.ZeroInt()
	}

	// Calculate the sum of collateral and extra voting power
	totalPower := v.Collateral.Add(v.ExtraVotingPower)

	// Calculate the maximum allowed by the leverage ratio
	// maxPower = collateral * maxLeverageRatio
	maxCollateralPower := math.LegacyNewDecFromInt(v.Collateral).Mul(maxLeverageRatio).TruncateInt()

	// Return the minimum of the two calculations
	if totalPower.GT(maxCollateralPower) {
		return maxCollateralPower
	}
	return totalPower
}

// NewLastValidatorPower creates a new LastValidatorPower instance
func NewLastValidatorPower(pubkey []byte, power int64) LastValidatorPower {
	return LastValidatorPower{
		Pubkey: pubkey,
		Power:  power,
	}
}

// NewWithdrawal creates a new Withdrawal instance
func NewWithdrawal(pubkey []byte, amount math.Int, receiver string, receivesAt uint64) Withdrawal {
	return Withdrawal{
		Pubkey:     pubkey,
		Amount:     amount,
		Receiver:   receiver,
		ReceivesAt: receivesAt,
	}
}

// ValidatorsEqual checks if two validators are equal based on their public keys
func ValidatorsEqual(v1, v2 Validator) bool {
	return bytes.Equal(v1.Pubkey, v2.Pubkey)
}

// FindValidator finds a validator in a slice by pubkey
func FindValidator(validators []Validator, pubkey []byte) (Validator, bool) {
	for _, v := range validators {
		if bytes.Equal(v.Pubkey, pubkey) {
			return v, true
		}
	}
	return Validator{}, false
}

// ABCIValidatorUpdate creates an ABCI validator update object from a validator
func (v Validator) ABCIValidatorUpdate() (abciVal abci.ValidatorUpdate, err error) {
	tmPubKey, err := k1util.PubKeyBytesToCosmos(v.Pubkey)
	if err != nil {
		return abci.ValidatorUpdate{}, fmt.Errorf("unable to convert validator pubkey: %w", err)
	}

	power := v.VotingPower.Int64()
	if v.Jailed {
		power = 0
	}

	abciVal = abci.ValidatorUpdate{
		PubKey: sdk.ToPubKeyPrototype(tmPubKey),
		Power:  power,
	}

	return abciVal, nil
}
