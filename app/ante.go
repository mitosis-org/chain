package app

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth/ante"
	evmengtypes "github.com/omni-network/omni/octane/evmengine/types"
)

func (app *MitosisApp) newAnteHandler() (sdk.AnteHandler, error) {
	anteHandler, err := NewAnteHandler(
		ante.HandlerOptions{
			AccountKeeper:   app.AccountKeeper,
			BankKeeper:      app.BankKeeper,
			SignModeHandler: app.txConfig.SignModeHandler(),
			FeegrantKeeper:  nil,
			SigGasConsumer:  ante.DefaultSigVerificationGasConsumer,
		},
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create ante handler: %w", err)
	}

	return anteHandler, nil
}

func NewAnteHandler(options ante.HandlerOptions) (sdk.AnteHandler, error) {
	defaultAnteHandler, err := ante.NewAnteHandler(options)
	if err != nil {
		return nil, err
	}

	return func(ctx sdk.Context, tx sdk.Tx, simulate bool) (newCtx sdk.Context, err error) {
		if isTxForEVM(tx) {
			// Disable checks for an EVM payload tx because it doesn't have a proper form and
			// is just created by custom prepare proposal handler.
			return ctx, nil
		} else {
			return defaultAnteHandler(ctx, tx, simulate)
		}
	}, nil
}

func isTxForEVM(tx sdk.Tx) bool {
	for _, msg := range tx.GetMsgs() {
		typeURL := sdk.MsgTypeURL(msg)
		if (typeURL == sdk.MsgTypeURL(&evmengtypes.MsgExecutionPayload{})) {
			return true
		}
	}

	return false
}
