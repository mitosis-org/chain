package keeper

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/mitosis-org/chain/x/evmvalidator/types"
	"time"
)

// ProcessMaturedWithdrawals processes withdrawals that have matured
func (k Keeper) ProcessMaturedWithdrawals(ctx sdk.Context) error {
	var err error

	params := k.GetParams(ctx)
	currentTime := ctx.BlockTime().Unix()
	withdrawalLimit := params.WithdrawalLimit
	processedCount := uint32(0)

	k.IterateWithdrawalsByMaturesAt(ctx, func(withdrawal types.Withdrawal) bool {
		// Check if we've processed enough withdrawals for this block
		if processedCount >= withdrawalLimit {
			return true
		}

		// Check if the withdrawal has not matured yet
		if currentTime < withdrawal.MaturesAt {
			return true // afterward, all withdrawals are not matured
		}

		// Insert the withdrawal into the EVM engine
		if err = k.evmEngKeeper.InsertWithdrawal(ctx, withdrawal.Receiver.Address(), withdrawal.Amount); err != nil {
			// err will be returned in the end of the function
			return true
		}

		// Emit an event and remove the withdrawal
		ctx.EventManager().EmitEvent(
			sdk.NewEvent(
				types.EventTypeWithdrawalMatured,
				sdk.NewAttribute(types.AttributeKeyWithdrawalID, fmt.Sprintf("%d", withdrawal.ID)),
				sdk.NewAttribute(types.AttributeKeyValAddr, withdrawal.ValAddr.String()),
				sdk.NewAttribute(types.AttributeKeyAmount, fmt.Sprintf("%d", withdrawal.Amount)),
				sdk.NewAttribute(types.AttributeKeyReceiver, withdrawal.Receiver.String()),
				sdk.NewAttribute(types.AttributeKeyMaturesAt, time.Unix(withdrawal.MaturesAt, 0).String()),
			),
		)

		k.DeleteWithdrawal(ctx, withdrawal)
		processedCount++

		return false
	})

	return err
}
