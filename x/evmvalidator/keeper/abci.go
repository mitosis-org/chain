package keeper

import (
	"bytes"
	"fmt"
	abci "github.com/cometbft/cometbft/abci/types"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/mitosis-org/chain/x/evmvalidator/types"
)

// EndBlocker is called at the end of every block
func (k Keeper) EndBlocker(ctx sdk.Context) ([]abci.ValidatorUpdate, error) {
	// Process matured withdrawals
	if err := k.ProcessMaturedWithdrawals(ctx); err != nil {
		return nil, err
	}

	// Update validator set
	return k.GetValidatorUpdates(ctx)
}

// ProcessMaturedWithdrawals processes withdrawals that have matured
func (k Keeper) ProcessMaturedWithdrawals(ctx sdk.Context) error {
	params := k.GetParams(ctx)
	currentTime := uint64(ctx.BlockTime().Unix())
	withdrawalLimit := params.WithdrawalLimit
	processedCount := uint32(0)

	k.IterateWithdrawalsQueue(ctx, currentTime, func(withdrawal types.Withdrawal) bool {
		// Check if we've processed enough withdrawals for this block
		if processedCount >= withdrawalLimit {
			return true // stop iteration
		}

		// Process the withdrawal
		// TODO: Send tokens to the receiver (will be implemented by the chain owner)
		// For now, just emit an event and remove from queue

		ctx.EventManager().EmitEvent(
			sdk.NewEvent(
				types.EventTypeWithdrawalMatured,
				sdk.NewAttribute(types.AttributeKeyPubkey, fmt.Sprintf("%X", withdrawal.Pubkey)),
				sdk.NewAttribute(types.AttributeKeyAmount, withdrawal.Amount.String()),
				sdk.NewAttribute(types.AttributeKeyReceiver, withdrawal.Receiver),
				sdk.NewAttribute(types.AttributeKeyReceivesAt, sdk.Uint64ToBigEndian(withdrawal.ReceivesAt)),
			),
		)

		k.DeleteWithdrawalFromQueue(ctx, withdrawal)
		processedCount++

		return false // continue iteration
	})

	return nil
}

// GetValidatorUpdates gets the validator set updates
func (k Keeper) GetValidatorUpdates(ctx sdk.Context) ([]abci.ValidatorUpdate, error) {
	// Get the last validator set
	lastValidatorPowers := make(map[string]int64)
	k.IterateLastValidatorPowers(ctx, func(pubkey []byte, power int64) bool {
		lastValidatorPowers[string(pubkey)] = power
		return false
	})

	// Get the current validators
	currentValidators := k.GetAllValidators(ctx)
	currentValidatorsByPubkey := make(map[string]types.Validator)
	currentPowersByPubkey := make(map[string]int64)

	for _, validator := range currentValidators {
		pubkeyStr := string(validator.Pubkey)
		currentValidatorsByPubkey[pubkeyStr] = validator

		// If jailed, power is 0
		power := validator.VotingPower.Int64()
		if validator.Jailed {
			power = 0
		}
		currentPowersByPubkey[pubkeyStr] = power
	}

	// Calculate the validator updates by comparing last powers with current powers
	var updates []abci.ValidatorUpdate

	// Check for validators that have been removed or changed power
	for pubkeyStr, lastPower := range lastValidatorPowers {
		// Get current power
		currentPower, exists := currentPowersByPubkey[pubkeyStr]

		// If validator no longer exists or power changed
		if !exists || lastPower != currentPower {
			pubkey := []byte(pubkeyStr)

			// If it still exists, use the current validator
			if exists {
				validator := currentValidatorsByPubkey[pubkeyStr]
				update, err := validator.ABCIValidatorUpdate()
				if err != nil {
					ctx.Logger().Error("failed to get validator update", "pubkey", fmt.Sprintf("%X", pubkey), "error", err)
					continue
				}
				updates = append(updates, update)
			} else {
				// If the validator no longer exists, remove it with power 0
				tmPubKey, err := types.GetAddressFromConsensusPublicKey(pubkey)
				if err != nil {
					ctx.Logger().Error("failed to convert pubkey to consensus address", "pubkey", fmt.Sprintf("%X", pubkey), "error", err)
					continue
				}

				// Create an update with power 0 to remove the validator
				update := abci.ValidatorUpdate{
					PubKey: abci.PubKey{
						Type: "secp256k1",
						Data: pubkey,
					},
					Power: 0,
				}
				updates = append(updates, update)

				// Delete the last power record
				k.DeleteLastValidatorPower(ctx, pubkey)
			}
		}
	}

	// Check for new validators
	for pubkeyStr, currentPower := range currentPowersByPubkey {
		pubkey := []byte(pubkeyStr)
		_, exists := lastValidatorPowers[pubkeyStr]

		// If validator is new
		if !exists {
			validator := currentValidatorsByPubkey[pubkeyStr]
			update, err := validator.ABCIValidatorUpdate()
			if err != nil {
				ctx.Logger().Error("failed to get validator update", "pubkey", fmt.Sprintf("%X", pubkey), "error", err)
				continue
			}
			updates = append(updates, update)

			// Set the last power record
			k.SetLastValidatorPower(ctx, pubkey, currentPower)
		}
	}

	// Ensure validator updates don't exceed max validators
	params := k.GetParams(ctx)
	bondedValidators := k.GetBondedValidatorsByPower(ctx)

	// Sort validators by power (done in GetBondedValidatorsByPower)
	// Get top N validators where N is MaxValidators
	var bondedPubkeys [][]byte
	for i, validator := range bondedValidators {
		if uint32(i) >= params.MaxValidators {
			break
		}
		bondedPubkeys = append(bondedPubkeys, validator.Pubkey)
	}

	// Ensure all updates are either removing validators or updating to the top N
	var filteredUpdates []abci.ValidatorUpdate
	for _, update := range updates {
		pubkey := update.PubKey.Data

		// Always include updates that remove validators (power 0)
		if update.Power == 0 {
			filteredUpdates = append(filteredUpdates, update)
			continue
		}

		// Check if this validator is in the top N
		found := false
		for _, bondedPubkey := range bondedPubkeys {
			if bytes.Equal(pubkey, bondedPubkey) {
				found = true
				break
			}
		}

		if found {
			filteredUpdates = append(filteredUpdates, update)
		} else {
			// If not in top N, set power to 0 to remove
			update.Power = 0
			filteredUpdates = append(filteredUpdates, update)
		}
	}

	return filteredUpdates, nil
}
