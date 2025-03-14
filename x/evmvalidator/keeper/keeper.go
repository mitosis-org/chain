package keeper

import (
	"cosmossdk.io/core/address"
	"cosmossdk.io/log"
	"cosmossdk.io/store/prefix"
	storetypes "cosmossdk.io/store/types"
	"fmt"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/mitosis-org/chain/bindings"
	mitotypes "github.com/mitosis-org/chain/types"
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

// GetValidator gets a validator by address
func (k Keeper) GetValidator(ctx sdk.Context, valAddr mitotypes.EthAddress) (validator types.Validator, found bool) {
	store := ctx.KVStore(k.storeKey)
	key := types.GetValidatorKey(valAddr)
	bz := store.Get(key)
	if bz == nil {
		return types.Validator{}, false
	}

	k.cdc.MustUnmarshal(bz, &validator)
	return validator, true
}

// HasValidator checks if a validator exists
func (k Keeper) HasValidator(ctx sdk.Context, valAddr mitotypes.EthAddress) bool {
	store := ctx.KVStore(k.storeKey)
	return store.Has(types.GetValidatorKey(valAddr))
}

// SetValidator sets a validator and updates the consensus address mapping
func (k Keeper) SetValidator(ctx sdk.Context, validator types.Validator) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshal(&validator)
	store.Set(types.GetValidatorKey(validator.Addr), bz)
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

	iterator := k.GetValidatorsByPowerIndexIterator(ctx)
	defer iterator.Close()

	for count := uint32(0); iterator.Valid() && count < maxValidators; iterator.Next() {
		// extract the validator address
		valAddr := mitotypes.BytesToEthAddress(iterator.Value())
		validator, found := k.GetValidator(ctx, valAddr)
		if !found {
			panic(fmt.Sprintf("validator not found: %s", valAddr.String()))
		}

		// defensive logic. not possible to have a jailed validator in the power index
		if validator.Jailed {
			k.Logger(ctx).Error(fmt.Sprintf("[BUG] validator %s is jailed", valAddr.String()))
			continue
		}

		validators = append(validators, validator)
		count++
	}

	return validators
}

// GetValidatorByConsAddr returns a validator by consensus address
func (k Keeper) GetValidatorByConsAddr(ctx sdk.Context, consAddr sdk.ConsAddress) (types.Validator, bool) {
	store := ctx.KVStore(k.storeKey)

	// Get the validator's address using the consensus address mapping
	valAddr := store.Get(types.GetValidatorByConsAddrKey(consAddr))
	if valAddr == nil {
		return types.Validator{}, false
	}

	// Get the validator using the EVM address
	return k.GetValidator(ctx, mitotypes.BytesToEthAddress(valAddr))
}

// SetValidatorByConsAddr sets a validator EVM address by consensus address
func (k Keeper) SetValidatorByConsAddr(ctx sdk.Context, consAddr sdk.ConsAddress, valAddr mitotypes.EthAddress) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.GetValidatorByConsAddrKey(consAddr), valAddr.Bytes())
}

// GetValidatorsByPowerIndexIterator returns an iterator for the power index (starting from the most powerful)
func (k Keeper) GetValidatorsByPowerIndexIterator(ctx sdk.Context) storetypes.Iterator {
	store := ctx.KVStore(k.storeKey)
	return storetypes.KVStorePrefixIterator(store, types.ValidatorByPowerIndexKeyPrefix)
}

// SetValidatorByPowerIndex sets a validator in the power index
func (k Keeper) SetValidatorByPowerIndex(ctx sdk.Context, power int64, valAddr mitotypes.EthAddress) {
	store := ctx.KVStore(k.storeKey)
	storeKey := types.GetValidatorByPowerIndexKey(power, valAddr)
	store.Set(storeKey, valAddr.Bytes())
}

// DeleteValidatorByPowerIndex deletes a validator from the power index
func (k Keeper) DeleteValidatorByPowerIndex(ctx sdk.Context, power int64, valAddr mitotypes.EthAddress) {
	store := ctx.KVStore(k.storeKey)
	storeKey := types.GetValidatorByPowerIndexKey(power, valAddr)
	store.Delete(storeKey)
}

// GetLastValidatorPower gets the last validator power for a validator
func (k Keeper) GetLastValidatorPower(ctx sdk.Context, valAddr mitotypes.EthAddress) (power int64, found bool) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.GetLastValidatorPowerKey(valAddr))
	if bz == nil {
		return 0, false
	}

	var lastPower types.LastValidatorPower
	k.cdc.MustUnmarshal(bz, &lastPower)
	return lastPower.Power, true
}

// SetLastValidatorPower sets the last validator power for a validator
func (k Keeper) SetLastValidatorPower(ctx sdk.Context, valAddr mitotypes.EthAddress, power int64) {
	store := ctx.KVStore(k.storeKey)
	lastPower := types.LastValidatorPower{ValAddr: valAddr, Power: power}
	bz := k.cdc.MustMarshal(&lastPower)
	store.Set(types.GetLastValidatorPowerKey(valAddr), bz)
}

// DeleteLastValidatorPower deletes the last validator power for a validator
func (k Keeper) DeleteLastValidatorPower(ctx sdk.Context, valAddr mitotypes.EthAddress) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(types.GetLastValidatorPowerKey(valAddr))
}

// IterateLastValidatorPowers iterates through all last validator powers
func (k Keeper) IterateLastValidatorPowers(ctx sdk.Context, cb func(valAddr mitotypes.EthAddress, power int64) (stop bool)) {
	store := ctx.KVStore(k.storeKey)
	iterator := storetypes.KVStorePrefixIterator(store, types.LastValidatorPowerKeyPrefix)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var lastPower types.LastValidatorPower
		k.cdc.MustUnmarshal(iterator.Value(), &lastPower)

		if cb(lastPower.ValAddr, lastPower.Power) {
			break
		}
	}
}

// IterateLastValidators iterates through the active validator set and perform the provided function
func (k Keeper) IterateLastValidators(ctx sdk.Context, cb func(index int64, validator types.Validator) (stop bool)) error {
	var returnErr error

	i := int64(0)
	k.IterateLastValidatorPowers(ctx, func(valAddr mitotypes.EthAddress, power int64) bool {
		validator, found := k.GetValidator(ctx, valAddr)
		if !found {
			// This should never happen
			returnErr = fmt.Errorf("validator not found: %s", valAddr.String())
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

	k.IterateLastValidatorPowers(ctx, func(valAddr mitotypes.EthAddress, power int64) bool {
		powers = append(powers, types.LastValidatorPower{ValAddr: valAddr, Power: power})
		return false
	})

	return powers
}

// AddWithdrawalToQueue adds a withdrawal to the queue
func (k Keeper) AddWithdrawalToQueue(ctx sdk.Context, withdrawal types.Withdrawal) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshal(&withdrawal)
	key := types.GetWithdrawalQueueKey(withdrawal.MaturesAt)
	store.Set(key, bz)

	// Also set in the validator index
	indexKey := types.GetWithdrawalByValidatorKey(withdrawal.ValAddr, withdrawal.MaturesAt)
	store.Set(indexKey, bz)
}

// IterateWithdrawalsQueue iterates through the withdrawals queue
func (k Keeper) IterateWithdrawalsQueue(ctx sdk.Context, cb func(withdrawal types.Withdrawal) (stop bool)) {
	store := ctx.KVStore(k.storeKey)
	prefix := prefix.NewStore(store, types.WithdrawalQueueKeyPrefix)

	iterator := prefix.Iterator(nil, nil)
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
	key := types.GetWithdrawalQueueKey(withdrawal.MaturesAt)
	store.Delete(key)

	// Also remove from the validator index
	indexKey := types.GetWithdrawalByValidatorKey(withdrawal.ValAddr, withdrawal.MaturesAt)
	store.Delete(indexKey)
}

// GetAllWithdrawals gets all withdrawals
func (k Keeper) GetAllWithdrawals(ctx sdk.Context) []types.Withdrawal {
	var withdrawals []types.Withdrawal

	k.IterateWithdrawalsQueue(ctx, func(withdrawal types.Withdrawal) bool {
		withdrawals = append(withdrawals, withdrawal)
		return false
	})

	return withdrawals
}

func (k Keeper) IterateWithdrawalsForValidator(
	ctx sdk.Context,
	valAddr mitotypes.EthAddress,
	cb func(withdrawal types.Withdrawal) (stop bool),
) {
	store := ctx.KVStore(k.storeKey)
	prefix := prefix.NewStore(store, types.GetWithdrawalByValidatorIterationKey(valAddr))

	iterator := prefix.Iterator(nil, nil)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var withdrawal types.Withdrawal
		k.cdc.MustUnmarshal(iterator.Value(), &withdrawal)

		if cb(withdrawal) {
			break
		}
	}
}
