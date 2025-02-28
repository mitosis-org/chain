package keeper

import (
	"fmt"
	abci "github.com/cometbft/cometbft/abci/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"time"

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
		// TODO(thai): using evmengine to process withdrawal

		// Emit an event and remove from queue
		ctx.EventManager().EmitEvent(
			sdk.NewEvent(
				types.EventTypeWithdrawalMatured,
				sdk.NewAttribute(types.AttributeKeyPubkey, fmt.Sprintf("%X", withdrawal.Pubkey)),
				sdk.NewAttribute(types.AttributeKeyAmount, withdrawal.Amount.String()),
				sdk.NewAttribute(types.AttributeKeyReceiver, withdrawal.Receiver),
				sdk.NewAttribute(types.AttributeKeyReceivesAt, time.Unix(int64(withdrawal.ReceivesAt), 0).String()),
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
	// TODO(thai): implement this
	return nil, nil
}
