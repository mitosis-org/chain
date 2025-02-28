package keeper

import (
	"context"
	"cosmossdk.io/core/address"
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/mitosis-org/chain/x/evmvalidator/types"
)

var _ types.StakingKeeper = (*Keeper)(nil)

// ValidatorAddressCodec returns the address codec for validators
func (k Keeper) ValidatorAddressCodec() address.Codec {
	return k.validatorAddressCodec
}

// ConsensusAddressCodec returns the address codec for consensus nodes
func (k Keeper) ConsensusAddressCodec() address.Codec {
	return k.consensusAddressCodec
}

// IterateValidators implements the StakingKeeper interface
// It iterates through validators and executes the provided function for each validator
func (k Keeper) IterateValidators(ctx context.Context, fn func(index int64, validator types.ValidatorI) (stop bool)) error {
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	var index int64
	k.IterateValidatorsExec(sdkCtx, func(_ int64, validator types.Validator) (stop bool) {
		return fn(index, validator)
	})
	return nil
}

// Validator returns a validator by validator address
func (k Keeper) Validator(ctx context.Context, valAddr sdk.ValAddress) (types.ValidatorI, error) {
	sdkCtx := sdk.UnwrapSDKContext(ctx)

	var foundValidator types.Validator
	var found bool

	k.IterateValidatorsExec(sdkCtx, func(_ int64, validator types.Validator) bool {
		consAddr, err := validator.ConsAddr()
		if err != nil {
			return false
		}

		// Derive validator address from consensus address
		// This is a bit of a workaround since we're storing by pubkey not by address
		valAddrFromConsAddr := sdk.ValAddress(consAddr)
		if valAddrFromConsAddr.Equals(valAddr) {
			foundValidator = validator
			found = true
			return true
		}
		return false
	})

	if !found {
		return nil, types.ErrValidatorNotFound
	}

	return foundValidator, nil
}

func (k *Keeper) ValidatorByConsAddr(ctx context.Context, consAddress sdk.ConsAddress) (types.ValidatorI, error) {
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	validator, found := k.GetValidatorByConsAddr(sdkCtx, consAddress)
	if !found {
		return nil, types.ErrValidatorNotFound
	}

	return validator, nil
}

func (k *Keeper) Slash(
	ctx context.Context,
	consAddress sdk.ConsAddress,
	infractionHeight int64,
	power int64,
	slashFraction math.LegacyDec,
) (math.Int, error) {
	return k.slash(sdk.UnwrapSDKContext(ctx), consAddress, infractionHeight, power, slashFraction)
}

// SlashWithInfractionReason implements the StakingKeeper interface
// It slashes a validator for an infraction committed at a specific height with a specific reason
func (k Keeper) SlashWithInfractionReason(
	ctx context.Context,
	consAddr sdk.ConsAddress,
	infractionHeight int64,
	power int64,
	slashFraction math.LegacyDec,
	_ stakingtypes.Infraction,
) (math.Int, error) {
	return k.Slash(ctx, consAddr, infractionHeight, power, slashFraction)
}

func (k *Keeper) Jail(ctx context.Context, consAddress sdk.ConsAddress) error {
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	validator, found := k.GetValidatorByConsAddr(sdkCtx, consAddress)
	if !found {
		return types.ErrValidatorNotFound
	}
	return k.jail(sdkCtx, &validator, "triggered from slashing module")
}

func (k *Keeper) Unjail(ctx context.Context, consAddress sdk.ConsAddress) error {
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	validator, found := k.GetValidatorByConsAddr(sdkCtx, consAddress)
	if !found {
		return types.ErrValidatorNotFound
	}
	return k.unjail(sdkCtx, &validator)
}

// MaxValidators implements the StakingKeeper interface
func (k Keeper) MaxValidators(ctx context.Context) (uint32, error) {
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	params := k.GetParams(sdkCtx)
	return params.MaxValidators, nil
}

// IsValidatorJailed implements the StakingKeeper interface
func (k Keeper) IsValidatorJailed(ctx context.Context, addr sdk.ConsAddress) (bool, error) {
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	validator, found := k.GetValidatorByConsAddr(sdkCtx, addr)
	if !found {
		return false, types.ErrValidatorNotFound
	}
	return validator.Jailed, nil
}
