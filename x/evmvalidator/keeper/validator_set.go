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

		// Call hook if the validator was not found in the last validator set
		if !found {
			consAddr, err := validator.ConsAddr()
			if err != nil {
				return nil, errors.Wrap(err, "failed to get consensus address")
			}
			if err = k.slashingKeeper.AfterValidatorBonded(ctx, consAddr); err != nil {
				return nil, errors.Wrap(err, "failed to call AfterValidatorBonded hook")
			}
		}

		// Create validator update for CometBFT
		abciUpdate, err := validator.ABCIValidatorUpdate()
		if err != nil {
			return nil, errors.Wrap(err, "create validator update")
		}

		validatorUpdates = append(validatorUpdates, abciUpdate)

		// Record that this validator was processed
		processedVals[validator.Addr] = true

		// Update the last validator power
		k.SetLastValidatorPower(sdkCtx, validator.Addr, currentPower)

		// Log the update
		if found {
			k.Logger(sdkCtx).Info("[Active Validator Set] Power changed",
				"val_addr", validator.Addr.String(),
				"val_pubkey", fmt.Sprintf("%X", validator.Pubkey),
				"previous_power", lastPower,
				"new_power", currentPower,
			)
		} else {
			k.Logger(sdkCtx).Info("[Active Validator Set] Added",
				"val_addr", validator.Addr.String(),
				"val_pubkey", fmt.Sprintf("%X", validator.Pubkey),
				"new_power", currentPower,
			)
		}
	}

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
			pk, err := k1util.PubKeyBytesToCosmos(validator.Pubkey)
			if err != nil {
				k.Logger(sdkCtx).Error("Failed to convert pubkey", "err", err)
				return false
			}

			cmtPk, err := cryptocodec.ToCmtProtoPublicKey(pk)
			if err != nil {
				k.Logger(sdkCtx).Error("Failed to convert to CometBFT pubkey", "err", err)
				return false
			}

			validatorUpdate := abci.ValidatorUpdate{
				PubKey: cmtPk,
				Power:  0,
			}
			validatorUpdates = append(validatorUpdates, validatorUpdate)

			k.Logger(sdkCtx).Info("[Active Validator Set] Removed",
				"val_addr", valAddr.String(),
				"val_pubkey", fmt.Sprintf("%X", validator.Pubkey),
				"previous_power", power,
			)

			// Remove from last validator powers since it's no longer a validator
			k.DeleteLastValidatorPower(sdkCtx, valAddr)
		}

		return false // continue iteration
	})

	return validatorUpdates, nil
}
