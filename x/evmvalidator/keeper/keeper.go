package keeper

import (
	"cosmossdk.io/core/address"
	"encoding/hex"
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
)

// Keeper of the evmvalidator store
type Keeper struct {
	cdc                            codec.BinaryCodec
	storeKey                       storetypes.StoreKey
	slashingKeeper                 types.SlashingKeeper // initialized later
	evmEngKeeper                   types.EvmEngineKeeper
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

// SetSlashingKeeper sets the slashing keeper
func (k *Keeper) SetSlashingKeeper(slashingKeeper types.SlashingKeeper) {
	k.slashingKeeper = slashingKeeper
}

// SetEvmEngineKeeper sets the evm engine keeper
func (k *Keeper) SetEvmEngineKeeper(evmEngKeeper types.EvmEngineKeeper) {
	k.evmEngKeeper = evmEngKeeper
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

// SetValidator sets a validator and updates the consensus address mapping
func (k Keeper) SetValidator(ctx sdk.Context, validator types.Validator) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshal(&validator)
	store.Set(types.GetValidatorKey(validator.Pubkey), bz)
}

// IterateValidatorsExec is an internal implementation of IterateValidators that works with the SDK Context
func (k Keeper) IterateValidatorsExec(ctx sdk.Context, fn func(index int64, validator types.Validator) (stop bool)) {
	store := ctx.KVStore(k.storeKey)
	iterator := storetypes.KVStorePrefixIterator(store, types.ValidatorKeyPrefix)
	defer iterator.Close()

	i := int64(0)
	for ; iterator.Valid(); iterator.Next() {
		var validator types.Validator
		k.cdc.MustUnmarshal(iterator.Value(), &validator)

		stop := fn(i, validator)
		if stop {
			break
		}
		i++
	}
}

// GetAllValidators gets all validators
func (k Keeper) GetAllValidators(ctx sdk.Context) (validators []types.Validator) {
	k.IterateValidatorsExec(ctx, func(_ int64, validator types.Validator) bool {
		validators = append(validators, validator)
		return false
	})
	return validators
}

// GetNotJailedValidatorsByPower gets not jailed validators sorted by power
func (k Keeper) GetNotJailedValidatorsByPower(ctx sdk.Context, maxValidators uint32) []types.Validator {
	validators := make([]types.Validator, 0, maxValidators)

	iterator := k.ValidatorsPowerStoreIterator(ctx)
	defer iterator.Close()

	for count := uint32(0); iterator.Valid() && count < maxValidators; iterator.Next() {
		// extract the validator pubkey
		pubkey := iterator.Value()
		validator, found := k.GetValidator(ctx, pubkey)
		if !found {
			panic(fmt.Sprintf("validator with pubkey %s not found", hex.EncodeToString(pubkey)))
		}

		if validator.Jailed {
			continue
		}

		validators = append(validators, validator)
		count++
	}

	return validators
}

// GetValidatorByConsAddr returns a validator pubkey by consensus address
func (k Keeper) GetValidatorByConsAddr(ctx sdk.Context, consAddr sdk.ConsAddress) (types.Validator, bool) {
	store := ctx.KVStore(k.storeKey)

	// Get the validator's pubkey using the consensus address mapping
	pubkey := store.Get(types.GetValidatorByConsAddrKey(consAddr))
	if pubkey == nil {
		return types.Validator{}, false
	}

	// Get the validator using the pubkey
	return k.GetValidator(ctx, pubkey)
}

// SetValidatorByConsAddr sets a validator pubkey by consensus address
func (k Keeper) SetValidatorByConsAddr(ctx sdk.Context, consAddr sdk.ConsAddress, pubkey []byte) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.GetValidatorByConsAddrKey(consAddr), pubkey)
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
func (k Keeper) SetValidatorByPowerIndex(ctx sdk.Context, power int64, pubkey []byte) {
	store := ctx.KVStore(k.storeKey)
	powerRankKey := types.GetValidatorPowerRankKey(power, pubkey)
	store.Set(powerRankKey, pubkey)
}

// DeleteValidatorByPowerIndex deletes a validator from the power index
func (k Keeper) DeleteValidatorByPowerIndex(ctx sdk.Context, power int64, pubkey []byte) {
	store := ctx.KVStore(k.storeKey)
	powerRankKey := types.GetValidatorPowerRankKey(power, pubkey)
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

// IterateLastValidators iterates through the active validator set and perform the provided function
func (k Keeper) IterateLastValidators(ctx sdk.Context, cb func(index int64, validator types.Validator) (stop bool)) error {
	var returnErr error

	i := int64(0)
	k.IterateLastValidatorPowers(ctx, func(pubkey []byte, power int64) bool {
		validator, found := k.GetValidator(ctx, pubkey)
		if !found {
			// This should never happen
			returnErr = fmt.Errorf("validator not found: %s", pubkey)
			return true
		}

		stop := cb(i, validator)
		i++
		return stop
	})

	return returnErr
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
