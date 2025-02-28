package keeper

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/mitosis-org/chain/x/evmvalidator/types"
	"time"
)

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

		// Insert the withdrawal into the EVM engine
		if err := k.evmEngKeeper.InsertWithdrawal(ctx, common.BytesToAddress(withdrawal.Receiver), withdrawal.Amount); err != nil {
			return false
		}

		// Emit an event and remove from queue
		ctx.EventManager().EmitEvent(
			sdk.NewEvent(
				types.EventTypeWithdrawalMatured,
				sdk.NewAttribute(types.AttributeKeyPubkey, fmt.Sprintf("%X", withdrawal.Pubkey)),
				sdk.NewAttribute(types.AttributeKeyAmount, fmt.Sprintf("%d", withdrawal.Amount)),
				sdk.NewAttribute(types.AttributeKeyReceiver, common.Bytes2Hex(withdrawal.Receiver)),
				sdk.NewAttribute(types.AttributeKeyReceivesAt, time.Unix(int64(withdrawal.ReceivesAt), 0).String()),
			),
		)

		k.DeleteWithdrawalFromQueue(ctx, withdrawal)
		processedCount++

		return false // continue iteration
	})

	return nil
}
