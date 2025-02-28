package keeper

import (
	"cosmossdk.io/core/address"
	sdkmath "cosmossdk.io/math"
	"fmt"
	"math"

	"cosmossdk.io/log"
	"cosmossdk.io/store/prefix"
	storetypes "cosmossdk.io/store/types"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/mitosis-org/chain/bindings"
	"github.com/mitosis-org/chain/x/evmvalidator/types"
	"github.com/omni-network/omni/lib/errors"
)

// Keeper of the evmvalidator store
type Keeper struct {
	cdc                            codec.BinaryCodec
	storeKey                       storetypes.StoreKey
	slashingKeeper                 types.SlashingKeeper
	evmValidatorEntrypointAddr     common.Address
	evmValidatorEntrypointContract *bindings.ConsensusValidatorEntrypoint

	// Address codecs for compatibility with other modules
	validatorAddressCodec address.Codec
	consensusAddressCodec address.Codec
}

// NewKeeperWithAddressCodecs creates a new keeper
func NewKeeperWithAddressCodecs(
	cdc codec.BinaryCodec,
	storeKey storetypes.StoreKey,
	slashingKeeper types.SlashingKeeper,
	evmValidatorEntrypointAddr common.Address,
	validatorAddressCodec address.Codec,
	consensusAddressCodec address.Codec,
) *Keeper {
	// Create contract binding to interact with events
	consensusValidatorEntrypointContract, err := bindings.NewConsensusValidatorEntrypoint(evmValidatorEntrypointAddr, nil)
	if err != nil {
		panic(fmt.Sprintf("failed to create consensus validator entrypoint contract: %v", err))
	}

	return &Keeper{
		cdc:                            cdc,
		storeKey:                       storeKey,
		slashingKeeper:                 slashingKeeper,
		evmValidatorEntrypointAddr:     evmValidatorEntrypointAddr,
		evmValidatorEntrypointContract: consensusValidatorEntrypointContract,
		validatorAddressCodec:          validatorAddressCodec,
		consensusAddressCodec:          consensusAddressCodec,
	}
}

// Logger returns a module-specific logger
func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

// GetParams gets the parameters for the x/evmvalidator module
func (k Keeper) GetParams(ctx sdk.Context) (params types.Params) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.ParamsKey)
	if bz == nil {
		return types.DefaultParams()
	}

	k.cdc.MustUnmarshal(bz, &params)
	return params
}

// SetParams sets the parameters for the x/evmvalidator module
func (k Keeper) SetParams(ctx sdk.Context, params types.Params) error {
	if err := params.Validate(); err != nil {
		return err
	}

	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshal(&params)
	store.Set(types.ParamsKey, bz)

	return nil
}

// GetValidator gets a validator by pubkey
func (k Keeper) GetValidator(ctx sdk.Context, pubkey []byte) (validator types.Validator, found bool) {
	store := ctx.KVStore(k.storeKey)
	key := types.GetValidatorKey(pubkey)
	bz := store.Get(key)
	if bz == nil {
		return types.Validator{}, false
	}

	k.cdc.MustUnmarshal(bz, &validator)
	return validator, true
}

// HasValidator checks if a validator exists
func (k Keeper) HasValidator(ctx sdk.Context, pubkey []byte) bool {
	store := ctx.KVStore(k.storeKey)
	return store.Has(types.GetValidatorKey(pubkey))
}

// SetValidator sets a validator
func (k Keeper) SetValidator(ctx sdk.Context, validator types.Validator) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshal(&validator)
	store.Set(types.GetValidatorKey(validator.Pubkey), bz)
}

// RemoveValidator removes a validator
func (k Keeper) RemoveValidator(ctx sdk.Context, pubkey []byte) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(types.GetValidatorKey(pubkey))
}

// IterateValidators iterates through all validators
func (k Keeper) IterateValidators(ctx sdk.Context, cb func(index int64, validator types.Validator) (stop bool)) {
	store := ctx.KVStore(k.storeKey)
	iterator := storetypes.KVStorePrefixIterator(store, types.ValidatorKeyPrefix)
	defer iterator.Close()

	i := int64(0)
	for ; iterator.Valid(); iterator.Next() {
		var validator types.Validator
		k.cdc.MustUnmarshal(iterator.Value(), &validator)

		stop := cb(i, validator)
		if stop {
			break
		}
		i++
	}
}

// GetAllValidators gets all validators
func (k Keeper) GetAllValidators(ctx sdk.Context) (validators []types.Validator) {
	k.IterateValidators(ctx, func(_ int64, validator types.Validator) bool {
		validators = append(validators, validator)
		return false
	})
	return validators
}

// GetBondedValidatorsByPower gets validators sorted by power
func (k Keeper) GetBondedValidatorsByPower(ctx sdk.Context) []types.Validator {
	maxValidators := k.GetParams(ctx).MaxValidators
	validators := make([]types.Validator, 0, maxValidators)

	iterator := k.ValidatorsPowerStoreIterator(ctx)
	defer iterator.Close()

	for count := uint32(0); iterator.Valid() && count < maxValidators; iterator.Next() {
		// extract the validator pubkey
		pubkey := iterator.Value()
		validator, found := k.GetValidator(ctx, pubkey)
		if !found {
			panic(fmt.Sprintf("validator with pubkey %s not found", sdk.ValAddress(pubkey)))
		}

		if validator.Jailed {
			continue
		}

		validators = append(validators, validator)
		count++
	}

	return validators
}

// ValidatorsPowerStoreIterator returns an iterator for the validators power store
func (k Keeper) ValidatorsPowerStoreIterator(ctx sdk.Context) storetypes.Iterator {
	store := ctx.KVStore(k.storeKey)
	return storetypes.KVStorePrefixIterator(store, types.ValidatorPowerRankStoreKeyPrefix)
}

// GetValidatorsByPowerIndexIterator returns an iterator for the validators power store with a starting rank
func (k Keeper) GetValidatorsByPowerIndexIterator(ctx sdk.Context) storetypes.Iterator {
	store := ctx.KVStore(k.storeKey)
	return storetypes.KVStorePrefixIterator(store, types.ValidatorPowerRankStoreKeyPrefix)
}

// SetValidatorByPowerIndex sets a validator in the power index
func (k Keeper) SetValidatorByPowerIndex(ctx sdk.Context, validator types.Validator) {
	// If jailed, delete from the power ranking
	if validator.Jailed {
		k.DeleteValidatorByPowerIndex(ctx, validator.VotingPower.Int64(), validator.Pubkey)
		return
	}

	// Set in the power ranking
	store := ctx.KVStore(k.storeKey)
	powerRankKey := types.GetValidatorPowerRankKey(validator)
	store.Set(powerRankKey, validator.Pubkey)
}

// DeleteValidatorByPowerIndex deletes a validator from the power index
func (k Keeper) DeleteValidatorByPowerIndex(ctx sdk.Context, power int64, pubkey []byte) {
	store := ctx.KVStore(k.storeKey)
	validator := types.Validator{
		VotingPower: sdkmath.NewInt(power),
		Pubkey:      pubkey,
	}
	powerRankKey := types.GetValidatorPowerRankKey(validator)
	store.Delete(powerRankKey)
}

// GetLastValidatorPower gets the last validator power for a validator
func (k Keeper) GetLastValidatorPower(ctx sdk.Context, pubkey []byte) (power int64, found bool) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.GetLastValidatorPowerKey(pubkey))
	if bz == nil {
		return 0, false
	}

	var lastPower types.LastValidatorPower
	k.cdc.MustUnmarshal(bz, &lastPower)
	return lastPower.Power, true
}

// SetLastValidatorPower sets the last validator power for a validator
func (k Keeper) SetLastValidatorPower(ctx sdk.Context, pubkey []byte, power int64) {
	store := ctx.KVStore(k.storeKey)
	lastPower := types.NewLastValidatorPower(pubkey, power)
	bz := k.cdc.MustMarshal(&lastPower)
	store.Set(types.GetLastValidatorPowerKey(pubkey), bz)
}

// DeleteLastValidatorPower deletes the last validator power for a validator
func (k Keeper) DeleteLastValidatorPower(ctx sdk.Context, pubkey []byte) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(types.GetLastValidatorPowerKey(pubkey))
}

// IterateLastValidatorPowers iterates through all last validator powers
func (k Keeper) IterateLastValidatorPowers(ctx sdk.Context, cb func(pubkey []byte, power int64) (stop bool)) {
	store := ctx.KVStore(k.storeKey)
	iterator := storetypes.KVStorePrefixIterator(store, types.LastValidatorPowerKeyPrefix)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var lastPower types.LastValidatorPower
		k.cdc.MustUnmarshal(iterator.Value(), &lastPower)

		if cb(lastPower.Pubkey, lastPower.Power) {
			break
		}
	}
}

// GetLastValidatorPowers gets all last validator powers
func (k Keeper) GetLastValidatorPowers(ctx sdk.Context) []types.LastValidatorPower {
	var powers []types.LastValidatorPower

	k.IterateLastValidatorPowers(ctx, func(pubkey []byte, power int64) bool {
		powers = append(powers, types.NewLastValidatorPower(pubkey, power))
		return false
	})

	return powers
}

// AddWithdrawalToQueue adds a withdrawal to the queue
func (k Keeper) AddWithdrawalToQueue(ctx sdk.Context, withdrawal types.Withdrawal) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshal(&withdrawal)
	key := types.GetWithdrawalQueueKey(withdrawal.ReceivesAt)
	store.Set(key, bz)
}

// IterateWithdrawalsQueue iterates through the withdrawals queue
func (k Keeper) IterateWithdrawalsQueue(ctx sdk.Context, endTime uint64, cb func(withdrawal types.Withdrawal) (stop bool)) {
	store := ctx.KVStore(k.storeKey)
	prefix := prefix.NewStore(store, types.WithdrawalQueueKeyPrefix)

	iterator := prefix.Iterator(nil, sdk.Uint64ToBigEndian(endTime+1))
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var withdrawal types.Withdrawal
		k.cdc.MustUnmarshal(iterator.Value(), &withdrawal)

		if cb(withdrawal) {
			break
		}
	}
}

// DeleteWithdrawalFromQueue deletes a withdrawal from the queue
func (k Keeper) DeleteWithdrawalFromQueue(ctx sdk.Context, withdrawal types.Withdrawal) {
	store := ctx.KVStore(k.storeKey)
	key := types.GetWithdrawalQueueKey(withdrawal.ReceivesAt)
	store.Delete(key)
}

// GetAllWithdrawals gets all withdrawals
func (k Keeper) GetAllWithdrawals(ctx sdk.Context) []types.Withdrawal {
	var withdrawals []types.Withdrawal

	k.IterateWithdrawalsQueue(ctx, math.MaxUint64, func(withdrawal types.Withdrawal) bool {
		withdrawals = append(withdrawals, withdrawal)
		return false
	})

	return withdrawals
}

// Slash slashes a validator's collateral by a fraction
// It calls the StakingKeeper's Slash method to maintain compatibility with x/slashing module
func (k Keeper) Slash(ctx sdk.Context, consAddr sdk.ConsAddress, infractionHeight int64, power int64, slashFraction sdkmath.LegacyDec) (sdkmath.Int, error) {
	// Find the validator by consensus address
	var validator types.Validator
	var found bool

	k.IterateValidators(ctx, func(_ int64, val types.Validator) bool {
		valConsAddr, err := val.ConsAddr()
		if err != nil {
			ctx.Logger().Error("failed to get consensus address", "err", err)
			return false
		}

		if valConsAddr.Equals(consAddr) {
			validator = val
			found = true
			return true
		}
		return false
	})

	if !found {
		return sdkmath.ZeroInt(), errors.Wrap(types.ErrValidatorNotFound, consAddr.String())
	}

	// Calculate the amount to slash
	// Note that we're slashing collateral, not voting power
	slashAmount := sdkmath.LegacyNewDecFromInt(validator.Collateral).Mul(slashFraction).TruncateInt()
	if slashAmount.GT(validator.Collateral) {
		slashAmount = validator.Collateral
	}

	// Update validator's collateral
	validator.Collateral = validator.Collateral.Sub(slashAmount)

	// Recompute voting power
	params := k.GetParams(ctx)
	oldVotingPower := validator.VotingPower
	validator.VotingPower = validator.ComputeVotingPower(params.MaxLeverageRatio)

	// Check if validator should be jailed due to insufficient voting power
	if !validator.Jailed && validator.VotingPower.LT(params.MinVotingPower) {
		validator.Jailed = true
		ctx.EventManager().EmitEvent(
			sdk.NewEvent(
				types.EventTypeJailValidator,
				sdk.NewAttribute(types.AttributeKeyPubkey, fmt.Sprintf("%X", validator.Pubkey)),
				sdk.NewAttribute(types.AttributeKeyReason, "insufficient voting power after slash"),
			),
		)
	}

	// Update the validator in state
	k.SetValidator(ctx, validator)

	// Update the validator in power index
	k.DeleteValidatorByPowerIndex(ctx, oldVotingPower.Int64(), validator.Pubkey)
	k.SetValidatorByPowerIndex(ctx, validator)

	// Record last validator power
	k.SetLastValidatorPower(ctx, validator.Pubkey, validator.VotingPower.Int64())

	// Emit events
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeSlashValidator,
			sdk.NewAttribute(types.AttributeKeyPubkey, fmt.Sprintf("%X", validator.Pubkey)),
			sdk.NewAttribute(types.AttributeKeyAmount, slashAmount.String()),
			sdk.NewAttribute(types.AttributeKeySlashFraction, slashFraction.String()),
			sdk.NewAttribute(types.AttributeKeyInfractionHeight, fmt.Sprintf("%d", infractionHeight)),
			sdk.NewAttribute(types.AttributeKeyInfractionPower, fmt.Sprintf("%d", power)),
		),
	)

	return slashAmount, nil
}

// GetValidatorByConsAddr returns a validator by consensus address
func (k Keeper) GetValidatorByConsAddr(ctx sdk.Context, consAddr sdk.ConsAddress) (types.Validator, bool) {
	var validator types.Validator
	var found bool

	k.IterateValidators(ctx, func(_ int64, val types.Validator) bool {
		valConsAddr, err := val.ConsAddr()
		if err != nil {
			ctx.Logger().Error("failed to get consensus address", "err", err)
			return false
		}

		if valConsAddr.Equals(consAddr) {
			validator = val
			found = true
			return true
		}
		return false
	})

	return validator, found
}

// Jail jails a validator
func (k Keeper) Jail(ctx sdk.Context, consAddr sdk.ConsAddress) error {
	validator, found := k.GetValidatorByConsAddr(ctx, consAddr)
	if !found {
		return types.ErrValidatorNotFound
	}

	if validator.Jailed {
		return nil // already jailed
	}

	validator.Jailed = true
	oldVotingPower := validator.VotingPower

	// Update the validator in state
	k.SetValidator(ctx, validator)

	// Update the validator in power index
	k.DeleteValidatorByPowerIndex(ctx, oldVotingPower.Int64(), validator.Pubkey)
	k.SetValidatorByPowerIndex(ctx, validator)

	// Record last validator power
	k.SetLastValidatorPower(ctx, validator.Pubkey, 0) // Zero power when jailed

	// Emit event
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeJailValidator,
			sdk.NewAttribute(types.AttributeKeyPubkey, fmt.Sprintf("%X", validator.Pubkey)),
			sdk.NewAttribute(types.AttributeKeyReason, "jailed by slashing module"),
		),
	)

	return nil
}

// Unjail unjails a validator
func (k Keeper) Unjail(ctx sdk.Context, consAddr sdk.ConsAddress) error {
	validator, found := k.GetValidatorByConsAddr(ctx, consAddr)
	if !found {
		return types.ErrValidatorNotFound
	}

	if !validator.Jailed {
		return nil // already unjailed
	}

	// Check if voting power meets minimum requirement
	params := k.GetParams(ctx)
	if validator.VotingPower.LT(params.MinVotingPower) {
		return errors.Wrap(types.ErrInvalidVotingPower, "voting power below minimum requirement")
	}

	validator.Jailed = false

	// Update the validator in state
	k.SetValidator(ctx, validator)

	// Update the validator in power index
	k.DeleteValidatorByPowerIndex(ctx, 0, validator.Pubkey) // Delete with power 0 (jailed)
	k.SetValidatorByPowerIndex(ctx, validator)

	// Record last validator power
	k.SetLastValidatorPower(ctx, validator.Pubkey, validator.VotingPower.Int64())

	// Emit event
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeUnjailValidator,
			sdk.NewAttribute(types.AttributeKeyPubkey, fmt.Sprintf("%X", validator.Pubkey)),
		),
	)

	return nil
}
