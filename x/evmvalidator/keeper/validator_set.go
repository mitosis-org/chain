package keeper

import (
	"context"
	"fmt"

	"github.com/mitosis-org/chain/x/evmvalidator/types"

	abci "github.com/cometbft/cometbft/abci/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	mitotypes "github.com/mitosis-org/chain/types"
	"github.com/omni-network/omni/lib/errors"
)

// ApplyAndReturnValidatorSetUpdates applies and returns accumulated updates to the validator set.
func (k Keeper) ApplyAndReturnValidatorSetUpdates(ctx context.Context) ([]abci.ValidatorUpdate, error) {
	sdkCtx := sdk.UnwrapSDKContext(ctx)

	// Get parameters to determine max validators
	params := k.GetParams(sdkCtx)
	maxValidators := params.MaxValidators

	// Get the current validator set ordered by voting power
	validators := k.GetNotJailedValidatorsByPower(sdkCtx, maxValidators)

	// Collect validator updates by comparing with last validator powers
	var validatorUpdates []abci.ValidatorUpdate

	// Create a map to track validators that are bonded in the active set
	bondedVals := make(map[mitotypes.EthAddress]bool)

	// Process current validators first
	for _, validator := range validators {
		consAddr := validator.MustConsAddr()

		// Get the last power from the store
		lastPower, found := k.GetLastValidatorPower(sdkCtx, validator.Addr)
		if !found {
			lastPower = 0
		}

		// if we get to a zero-power validator (which we don't bond),
		// there are no more possible bonded validators
		currentPower := validator.ConsensusVotingPower()
		if currentPower <= 0 {
			break
		}

		// Record that this validator should be being bonded
		bondedVals[validator.Addr] = true

		// Skip if no change in voting power
		if currentPower == lastPower {
			continue
		}

		// Update the last validator power
		k.SetLastValidatorPower(sdkCtx, validator.Addr, currentPower)

		if !found {
			// Call hook if the validator becomes bonded
			if err := k.slashingKeeper.AfterValidatorBonded(ctx, consAddr); err != nil {
				return nil, errors.Wrap(err, "failed to call AfterValidatorBonded hook")
			}

			// Set the validator as bonded
			validator.Bonded = true
			k.SetValidator(sdkCtx, validator)
		}

		// Append to validator updates
		abciUpdate, err := validator.ABCIValidatorUpdate()
		if err != nil {
			return nil, errors.Wrap(err, "create abci validator update")
		}
		validatorUpdates = append(validatorUpdates, abciUpdate)

		// Log the update
		if found {
			k.Logger(sdkCtx).Debug("😈 Active Validator Set: Power changed",
				"val_addr", validator.Addr.String(),
				"val_pubkey", fmt.Sprintf("%X", validator.Pubkey),
				"cons_addr_hex", fmt.Sprintf("%X", consAddr.Bytes()),
				"previous_power", lastPower,
				"new_power", currentPower,
			)
		} else {
			k.Logger(sdkCtx).Info("😈 Active Validator Set: Bonded",
				"val_addr", validator.Addr.String(),
				"val_pubkey", fmt.Sprintf("%X", validator.Pubkey),
				"cons_addr_hex", fmt.Sprintf("%X", consAddr.Bytes()),
				"new_power", currentPower,
			)
		}
	}

	var err error

	// Process validators that were removed (not in the current set)
	// We need to iterate through all last powers and check if they've been processed
	k.IterateLastValidatorPowers(sdkCtx, func(valAddr mitotypes.EthAddress, power int64) bool {
		if bondedVals[valAddr] {
			// This validator is still bonded in the active set
			return false
		}

		// This validator is no longer bonded in the active set. So we need to unbond it.

		validator, found := k.GetValidator(sdkCtx, valAddr)
		if !found {
			err = errors.Wrap(types.ErrValidatorNotFound, "validator not found for address %s [BUG]", valAddr)
			return true
		}

		// Remove from last validator powers since it's no longer active validator
		k.DeleteLastValidatorPower(sdkCtx, valAddr)

		// Set the validator as not bonded
		validator.Bonded = false
		k.SetValidator(sdkCtx, validator)

		// Append to validator updates
		abciUpdate, err2 := validator.ABCIValidatorUpdateForUnbonding()
		if err2 != nil {
			err = errors.Wrap(err2, "create abci validator update")
			return true
		}
		validatorUpdates = append(validatorUpdates, abciUpdate)

		// Log the removal
		k.Logger(sdkCtx).Info("😈 Active Validator Set: Unbonded",
			"val_addr", valAddr.String(),
			"val_pubkey", fmt.Sprintf("%X", validator.Pubkey),
			"cons_addr_hex", fmt.Sprintf("%X", validator.MustConsAddr().Bytes()),
			"previous_power", power,
		)

		return false // continue iteration
	})

	if err != nil {
		return nil, err
	}

	return validatorUpdates, nil
}
