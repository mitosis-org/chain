package types

import (
	cmtprotocrypto "github.com/cometbft/cometbft/proto/tendermint/crypto"
	cryptocodec "github.com/cosmos/cosmos-sdk/crypto/codec"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"

	"cosmossdk.io/math"
	abci "github.com/cometbft/cometbft/abci/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/omni-network/omni/lib/k1util"
)

func (v Validator) ConsPubKey() (cryptotypes.PubKey, error) {
	return k1util.PubKeyBytesToCosmos(v.Pubkey)
}

func (v Validator) MustConsPubKey() cryptotypes.PubKey {
	pk, err := v.ConsPubKey()
	if err != nil {
		panic(err)
	}
	return pk
}

func (v Validator) CmtConsPublicKey() (cmtprotocrypto.PublicKey, error) {
	pk, err := v.ConsPubKey()
	if err != nil {
		return cmtprotocrypto.PublicKey{}, err
	}

	tmPk, err := cryptocodec.ToCmtProtoPublicKey(pk)
	if err != nil {
		return cmtprotocrypto.PublicKey{}, err
	}

	return tmPk, nil
}

// ConsAddr returns the validator's consensus address
func (v Validator) ConsAddr() (sdk.ConsAddress, error) {
	pubkey, err := v.ConsPubKey()
	if err != nil {
		return nil, err
	}
	return sdk.GetConsAddress(pubkey), nil
}

// MustConsAddr returns the validator's consensus address.
// Panics if error.
func (v Validator) MustConsAddr() sdk.ConsAddress {
	consAddr, err := v.ConsAddr()
	if err != nil {
		panic(err)
	}
	return consAddr
}

// ConsensusVotingPower returns the consensus voting power.
func (v Validator) ConsensusVotingPower() int64 {
	if v.Jailed {
		return 0
	}

	return v.VotingPower
}

// ComputeVotingPower calculates voting power based on collateral and extra voting power
// with respect to the max leverage ratio
func (v Validator) ComputeVotingPower(maxLeverageRatio math.LegacyDec) int64 {
	collateralPower := math.LegacyNewDecFromBigInt(v.Collateral.BigInt()).QuoInt(VotingPowerReduction)
	extraPower := math.LegacyNewDecFromBigInt(v.ExtraVotingPower.BigInt()).QuoInt(VotingPowerReduction)
	totalPower := collateralPower.Add(extraPower).TruncateInt64()

	// Calculate the maximum allowed by the leverage ratio
	// maxPower = collateral * maxLeverageRatio
	maxPower := collateralPower.Mul(maxLeverageRatio).TruncateInt64()

	// Return the minimum of the two calculations
	if totalPower > maxPower {
		return maxPower
	} else {
		return totalPower
	}
}

// ABCIValidatorUpdate creates an ABCI validator update object from a validator
func (v Validator) ABCIValidatorUpdate() (abciVal abci.ValidatorUpdate, err error) {
	tmPubKey, err := v.CmtConsPublicKey()
	if err != nil {
		return abci.ValidatorUpdate{}, err
	}

	abciVal = abci.ValidatorUpdate{
		PubKey: tmPubKey,
		Power:  v.ConsensusVotingPower(),
	}

	return abciVal, nil
}

// MustABCIValidatorUpdate creates an ABCI validator update object from a validator.
// Panics if error.
func (v Validator) MustABCIValidatorUpdate() abci.ValidatorUpdate {
	abciVal, err := v.ABCIValidatorUpdate()
	if err != nil {
		panic(err)
	}
	return abciVal
}

// ABCIValidatorUpdateForUnbonding creates an ABCI validator update object from a validator with zero power
func (v Validator) ABCIValidatorUpdateForUnbonding() (abciVal abci.ValidatorUpdate, err error) {
	tmPubKey, err := v.CmtConsPublicKey()
	if err != nil {
		return abci.ValidatorUpdate{}, err
	}

	abciVal = abci.ValidatorUpdate{
		PubKey: tmPubKey,
		Power:  0,
	}

	return abciVal, nil
}

// MustABCIValidatorUpdateForUnbonding creates an ABCI validator update object from a validator with zero power.
// Panics if error.
func (v Validator) MustABCIValidatorUpdateForUnbonding() abci.ValidatorUpdate {
	abciVal, err := v.ABCIValidatorUpdateForUnbonding()
	if err != nil {
		panic(err)
	}
	return abciVal
}
