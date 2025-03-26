package keeper

import (
	"context"

	"github.com/omni-network/omni/lib/errors"

	"cosmossdk.io/core/address"
	"cosmossdk.io/math"
	evidencetypes "cosmossdk.io/x/evidence/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	genutiltypes "github.com/cosmos/cosmos-sdk/x/genutil/types"
	slashingtypes "github.com/cosmos/cosmos-sdk/x/slashing/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/mitosis-org/chain/x/evmvalidator/types"
)

var (
	_ slashingtypes.StakingKeeper = (*Keeper)(nil)
	_ genutiltypes.StakingKeeper  = (*Keeper)(nil)
	_ evidencetypes.StakingKeeper = (*KeeperWrapperForEvidence)(nil)
)

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
func (k Keeper) IterateValidators(ctx context.Context, fn func(index int64, validator slashingtypes.ValidatorI) (stop bool)) error {
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	var index int64
	k.IterateValidators_(sdkCtx, func(_ int64, validator types.Validator) (stop bool) {
		return fn(index, validator)
	})
	return nil
}

func (k *Keeper) ValidatorByConsAddr(ctx context.Context, consAddress sdk.ConsAddress) (slashingtypes.ValidatorI, error) {
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
	sdkCtx := sdk.UnwrapSDKContext(ctx)

	// Find the validator by consensus address
	validator, found := k.GetValidatorByConsAddr(sdkCtx, consAddress)
	if !found {
		return math.ZeroInt(), errors.Wrap(types.ErrValidatorNotFound, consAddress.String())
	}

	slashedAmount, err := k.Slash_(sdkCtx, &validator, infractionHeight, power, slashFraction)
	if err != nil {
		return math.Int{}, err
	}

	return math.NewIntFromBigInt(slashedAmount.BigInt()), nil
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
	k.Jail_(sdkCtx, &validator, "triggered from slashing module")
	return nil
}

func (k *Keeper) Unjail(ctx context.Context, consAddress sdk.ConsAddress) error {
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	validator, found := k.GetValidatorByConsAddr(sdkCtx, consAddress)
	if !found {
		return types.ErrValidatorNotFound
	}
	return k.Unjail_(sdkCtx, &validator)
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

type KeeperWrapperForEvidence struct {
	K *Keeper
}

func (k KeeperWrapperForEvidence) ConsensusAddressCodec() address.Codec {
	return k.K.ConsensusAddressCodec()
}

func (k KeeperWrapperForEvidence) ValidatorByConsAddr(ctx context.Context, consAddress sdk.ConsAddress) (evidencetypes.ValidatorI, error) {
	return k.K.ValidatorByConsAddr(ctx, consAddress)
}
