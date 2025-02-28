package keeper

import (
	"context"
	"fmt"
	abci "github.com/cometbft/cometbft/abci/types"
	cryptocodec "github.com/cosmos/cosmos-sdk/crypto/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
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
	processedVals := make(map[string]bool)

	// Process current validators first
	for _, validator := range validators {
		// Get the last power from the store
		lastPower, found := k.GetLastValidatorPower(sdkCtx, validator.Pubkey)
		if !found {
			lastPower = 0
		}

		// Skip if no change in voting power
		currentPower := validator.ConsensusVotingPower()
		if currentPower == lastPower {
			continue
		}

		// Create validator update for CometBFT
		abciUpdate, err := validator.ABCIValidatorUpdate()
		if err != nil {
			return nil, errors.Wrap(err, "create validator update")
		}

		validatorUpdates = append(validatorUpdates, abciUpdate)

		// Record that this validator was processed
		processedVals[string(validator.Pubkey)] = true

		// Update the last validator power
		k.SetLastValidatorPower(sdkCtx, validator.Pubkey, currentPower)

		// Log the update
		sdkCtx.Logger().Info("Consensus power changed in validator set",
			"validator", fmt.Sprintf("%X", validator.Pubkey),
			"previous_power", lastPower,
			"new_power", currentPower,
		)
	}

	// Process validators that were removed (not in the current set)
	// We need to iterate through all last powers and check if they've been processed
	k.IterateLastValidatorPowers(sdkCtx, func(pubkey []byte, power int64) bool {
		if processedVals[string(pubkey)] {
			// Already processed this validator
			return false
		}

		// This validator is no longer in the active set or has been jailed
		// Create a validator update with power 0 to remove it
		validator, found := k.GetValidator(sdkCtx, pubkey)
		if !found || validator.Jailed {
			// Create a zero power update
			pk, err := k1util.PubKeyBytesToCosmos(pubkey)
			if err != nil {
				sdkCtx.Logger().Error("Failed to convert pubkey", "err", err)
				return false
			}

			cmtPk, err := cryptocodec.ToCmtProtoPublicKey(pk)
			if err != nil {
				sdkCtx.Logger().Error("Failed to convert to CometBFT pubkey", "err", err)
				return false
			}

			validatorUpdate := abci.ValidatorUpdate{
				PubKey: cmtPk,
				Power:  0,
			}
			validatorUpdates = append(validatorUpdates, validatorUpdate)

			sdkCtx.Logger().Info("Validator excluded from active set",
				"validator", fmt.Sprintf("%X", pubkey),
				"previous_power", power,
			)

			// Remove from last validator powers since it's no longer a validator
			k.DeleteLastValidatorPower(sdkCtx, pubkey)
		}

		return false // continue iteration
	})

	return validatorUpdates, nil
}
