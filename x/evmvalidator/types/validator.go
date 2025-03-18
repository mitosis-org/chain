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
	collateralPower := math.LegacyNewDecFromInt(v.Collateral).QuoInt(VotingPowerReductionForGwei)
	totalPower := collateralPower.Add(v.ExtraVotingPower).TruncateInt64()

	// Calculate the maximum allowed by the leverage ratio
	// maxPower = collateral * maxLeverageRatio
	maxCollateralPower := collateralPower.Mul(maxLeverageRatio).TruncateInt64()

	// Return the minimum of the two calculations
	if totalPower > maxCollateralPower {
		return maxCollateralPower
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
