package keeper

import (
	abci "github.com/cometbft/cometbft/abci/types"
	"github.com/cosmos/cosmos-sdk/telemetry"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/mitosis-org/chain/x/evmvalidator/types"
)

// EndBlocker is called at the end of every block
func (k Keeper) EndBlocker(ctx sdk.Context) ([]abci.ValidatorUpdate, error) {
	defer telemetry.ModuleMeasureSince(types.ModuleName, telemetry.Now(), telemetry.MetricKeyEndBlocker)

	// Process matured withdrawals
	if err := k.ProcessMaturedWithdrawals(ctx); err != nil {
		return nil, err
	}

	// Update active validator set
	return k.ApplyAndReturnValidatorSetUpdates(ctx)
}
