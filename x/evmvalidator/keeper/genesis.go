package keeper

import (
	abci "github.com/cometbft/cometbft/abci/types"
	cmttypes "github.com/cometbft/cometbft/types"
	cryptocodec "github.com/cosmos/cosmos-sdk/crypto/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/mitosis-org/chain/x/evmvalidator/types"
)

// InitGenesis initializes the evmvalidator module's state from a provided genesis state
func (k Keeper) InitGenesis(ctx sdk.Context, data *types.GenesisState) ([]abci.ValidatorUpdate, error) {
	// Set module parameters
	err := k.SetParams(ctx, data.Params)
	if err != nil {
		return []abci.ValidatorUpdate{}, err
	}

	// Set validators
	for _, validator := range data.Validators {
		// voting power will be recalculated
		if err = k.registerValidator(ctx, validator.Addr, validator.Pubkey, validator.Collateral, validator.ExtraVotingPower, validator.Jailed); err != nil {
			return nil, err
		}
	}

	// Set withdrawals
	for _, withdrawal := range data.Withdrawals {
		k.AddWithdrawalToQueue(ctx, withdrawal)
	}

	// Set last validator powers if provided
	for _, lastPower := range data.LastValidatorPowers {
		k.SetLastValidatorPower(ctx, lastPower.ValAddr, lastPower.Power)
	}

	return k.ApplyAndReturnValidatorSetUpdates(ctx)
}

// ExportGenesis returns the evmvalidator module's exported genesis state
func (k Keeper) ExportGenesis(ctx sdk.Context) *types.GenesisState {
	return types.NewGenesisState(
		k.GetParams(ctx),
		k.GetAllValidators(ctx),
		k.GetAllWithdrawals(ctx),
		k.GetLastValidatorPowers(ctx),
	)
}

// WriteValidators returns a slice of bonded genesis validators.
func (k Keeper) WriteValidators(ctx sdk.Context) (vals []cmttypes.GenesisValidator, returnErr error) {
	err := k.IterateLastValidators(ctx, func(_ int64, validator types.Validator) (stop bool) {
		pk, err := validator.ConsPubKey()
		if err != nil {
			returnErr = err
			return true
		}
		cmtPk, err := cryptocodec.ToCmtPubKeyInterface(pk)
		if err != nil {
			returnErr = err
			return true
		}

		vals = append(vals, cmttypes.GenesisValidator{
			Address: sdk.ConsAddress(cmtPk.Address()).Bytes(),
			PubKey:  cmtPk,
			Power:   validator.ConsensusVotingPower(),
			Name:    cmtPk.Address().String(),
		})

		return false
	})
	if err != nil {
		return nil, err
	}

	return vals, returnErr
}
