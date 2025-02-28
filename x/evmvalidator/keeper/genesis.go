package keeper

import (
	"encoding/hex"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/mitosis-org/chain/x/evmvalidator/types"
)

// InitGenesis initializes the evmvalidator module's state from a provided genesis state
func (k Keeper) InitGenesis(ctx sdk.Context, data *types.GenesisState) error {
	// Set module parameters
	err := k.SetParams(ctx, data.Params)
	if err != nil {
		return err
	}

	// Set validators
	for _, validator := range data.Validators {
		// Ensure voting power is computed
		validator.VotingPower = validator.ComputeVotingPower(data.Params.MaxLeverageRatio)

		// Set validator
		k.SetValidator(ctx, validator)

		// Set validator in power index
		k.SetValidatorByPowerIndex(ctx, validator)

		// Set last validator power
		power := validator.VotingPower.Int64()
		if validator.Jailed {
			power = 0
		}
		k.SetLastValidatorPower(ctx, validator.Pubkey, power)

		ctx.EventManager().EmitEvent(
			sdk.NewEvent(
				types.EventTypeRegisterValidator,
				sdk.NewAttribute(types.AttributeKeyPubkey, hex.EncodeToString(validator.Pubkey)),
				sdk.NewAttribute(types.AttributeKeyCollateral, validator.Collateral.String()),
				sdk.NewAttribute(types.AttributeKeyExtraVotingPower, validator.ExtraVotingPower.String()),
				sdk.NewAttribute(types.AttributeKeyVotingPower, validator.VotingPower.String()),
			),
		)
	}

	// Set withdrawals
	for _, withdrawal := range data.Withdrawals {
		k.AddWithdrawalToQueue(ctx, withdrawal)
	}

	// Set last validator powers if provided
	for _, lastPower := range data.LastValidatorPowers {
		k.SetLastValidatorPower(ctx, lastPower.Pubkey, lastPower.Power)
	}

	return nil
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
