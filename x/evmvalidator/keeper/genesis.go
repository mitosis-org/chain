package keeper

import (
	abci "github.com/cometbft/cometbft/abci/types"
	cmttypes "github.com/cometbft/cometbft/types"
	cryptocodec "github.com/cosmos/cosmos-sdk/crypto/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	mitotypes "github.com/mitosis-org/chain/types"
	"github.com/mitosis-org/chain/x/evmvalidator/types"
	"github.com/omni-network/omni/lib/errors"
)

// InitGenesis initializes the evmvalidator module's state from a provided genesis state
func (k Keeper) InitGenesis(ctx sdk.Context, data *types.GenesisState) ([]abci.ValidatorUpdate, error) {
	// Set module parameters
	err := k.SetParams(ctx, data.Params)
	if err != nil {
		return []abci.ValidatorUpdate{}, err
	}

	// Set ConsensusValidatorEntrypoint contract address
	k.SetValidatorEntrypointContractAddr(ctx, data.ValidatorEntrypointContractAddr)

	// Validate that each validator has only one collateral owner
	initialCollateralOwnersByValidator := make(map[mitotypes.EthAddress]mitotypes.EthAddress)
	for _, ownership := range data.CollateralOwnerships {
		if _, ok := initialCollateralOwnersByValidator[ownership.ValAddr]; ok {
			return []abci.ValidatorUpdate{}, errors.New("only one collateral owner per validator is allowed in genesis")
		}
		initialCollateralOwnersByValidator[ownership.ValAddr] = ownership.Owner
	}

	// Set validators
	for _, validator := range data.Validators {
		initialCollateralOwner, ok := initialCollateralOwnersByValidator[validator.Addr]
		if !ok {
			return []abci.ValidatorUpdate{}, errors.New("validator has no initial collateral owner")
		}

		// NOTE: validator.CollateralShares is ignored.
		if err = k.RegisterValidator(ctx, validator.Addr, validator.Pubkey, initialCollateralOwner, validator.Collateral, validator.ExtraVotingPower, validator.Jailed); err != nil {
			return nil, err
		}
	}

	// Set withdrawals
	for _, withdrawal := range data.Withdrawals {
		k.AddNewWithdrawalWithNextID(ctx, &withdrawal)
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
		k.GetValidatorEntrypointContractAddr(ctx),
		k.GetAllValidators(ctx),
		k.GetAllWithdrawals(ctx),
		k.GetLastValidatorPowers(ctx),
		k.GetAllCollateralOwnerships(ctx),
	)
}

// WriteValidators returns a slice of bonded genesis validators.
func (k Keeper) WriteValidators(ctx sdk.Context) (vals []cmttypes.GenesisValidator, returnErr error) {
	err := k.IterateLastValidators(ctx, func(_ int64, validator types.Validator) (stop bool) {
		pk := validator.MustConsPubKey()
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
