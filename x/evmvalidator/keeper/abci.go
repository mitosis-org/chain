package keeper

import (
	abci "github.com/cometbft/cometbft/abci/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// EndBlocker is called at the end of every block
func (k Keeper) EndBlocker(ctx sdk.Context) ([]abci.ValidatorUpdate, error) {
	// Process matured withdrawals
	if err := k.ProcessMaturedWithdrawals(ctx); err != nil {
		return nil, err
	}

	// Update active validator set
	return k.ApplyAndReturnValidatorSetUpdates(ctx)
}
