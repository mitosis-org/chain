package keeper

import (
	"context"
	"fmt"

	abci "github.com/cometbft/cometbft/abci/types"
	cryptocodec "github.com/cosmos/cosmos-sdk/crypto/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	mitotypes "github.com/mitosis-org/chain/types"
	"github.com/omni-network/omni/lib/errors"
	"github.com/omni-network/omni/lib/k1util"
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

	// Create a map to track validators that have been updated
	processedVals := make(map[mitotypes.EthAddress]bool)

	// Process current validators first
	for _, validator := range validators {
		// Get the last power from the store
		lastPower, found := k.GetLastValidatorPower(sdkCtx, validator.Addr)
		if !found {
			lastPower = 0
		}

		// Skip if no change in voting power
		currentPower := validator.ConsensusVotingPower()
		if currentPower == lastPower {
			continue
		}

		// if we get to a zero-power validator (which we don't bond),
		// there are no more possible bonded validators
		if currentPower <= 0 {
			break
		}

		// Update the last validator power
		k.SetLastValidatorPower(sdkCtx, validator.Addr, currentPower)

		if !found {
			// Call hook if the validator becomes bonded
			consAddr, err := validator.ConsAddr()
			if err != nil {
				return nil, errors.Wrap(err, "failed to get consensus address")
			}
			if err = k.slashingKeeper.AfterValidatorBonded(ctx, consAddr); err != nil {
				return nil, errors.Wrap(err, "failed to call AfterValidatorBonded hook")
			}

			// Set the validator as bonded
			validator.Bonded = true
			k.SetValidator(sdkCtx, validator)
		}

		// Record that this validator was processed
		processedVals[validator.Addr] = true

		// Append to validator updates
		abciUpdate, err := validator.ABCIValidatorUpdate()
		if err != nil {
			return nil, errors.Wrap(err, "create validator update")
		}
		validatorUpdates = append(validatorUpdates, abciUpdate)

		consAddr, err := validator.ConsAddr()
		if err != nil {
			return nil, errors.Wrap(err, "failed to get consensus address")
		}

		// Log the update
		if found {
			k.Logger(sdkCtx).Info("⚡️Active Validator Set: Power changed",
				"val_addr", validator.Addr.String(),
				"val_pubkey", fmt.Sprintf("%X", validator.Pubkey),
				"cons_addr_hex", fmt.Sprintf("%X", consAddr.Bytes()),
				"previous_power", lastPower,
				"new_power", currentPower,
			)
		} else {
			k.Logger(sdkCtx).Info("⚡️Active Validator Set: Bonded",
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
		if processedVals[valAddr] {
			// Already processed this validator
			return false
		}

		// This validator is no longer in the active set or has been jailed
		// Create a validator update with power 0 to remove it
		validator, found := k.GetValidator(sdkCtx, valAddr)
		if !found || validator.Jailed {
			// Create a zero power update
			pk, err2 := k1util.PubKeyBytesToCosmos(validator.Pubkey)
			if err2 != nil {
				err = errors.Wrap(err2, "failed to convert pubkey")
				return true
			}

			// Remove from last validator powers since it's no longer active validator
			k.DeleteLastValidatorPower(sdkCtx, valAddr)

			// Set the validator as not bonded
			if validator.Bonded {
				validator.Bonded = false
				k.SetValidator(sdkCtx, validator)
			}

			// Append to validator updates
			cmtPk, err2 := cryptocodec.ToCmtProtoPublicKey(pk)
			if err2 != nil {
				err = errors.Wrap(err2, "failed to convert to CometBFT pubkey")
				return true
			}
			validatorUpdate := abci.ValidatorUpdate{
				PubKey: cmtPk,
				Power:  0,
			}
			validatorUpdates = append(validatorUpdates, validatorUpdate)

			// Log the removal
			consAddr, err2 := validator.ConsAddr()
			if err2 != nil {
				err = errors.Wrap(err2, "failed to get consensus address")
				return true
			}
			k.Logger(sdkCtx).Info("⚡️Active Validator Set: Unbonded",
				"val_addr", valAddr.String(),
				"val_pubkey", fmt.Sprintf("%X", validator.Pubkey),
				"cons_addr_hex", fmt.Sprintf("%X", consAddr.Bytes()),
				"previous_power", power,
			)
		}

		return false // continue iteration
	})

	if err != nil {
		return nil, err
	}

	return validatorUpdates, nil
}
