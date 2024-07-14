package app

import (
	"encoding/json"
	"fmt"
	tmtypes "github.com/cometbft/cometbft/types"
	servertypes "github.com/cosmos/cosmos-sdk/server/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	slashingtypes "github.com/cosmos/cosmos-sdk/x/slashing/types"
)

// ExportAppStateAndValidators exports the state of the application for a genesis
// file.
func (app *MitosisApp) ExportAppStateAndValidators(
	forZeroHeight bool,
	jailAllowedAddrs []string,
	modulesToExport []string,
) (servertypes.ExportedApp, error) {
	// as if they could withdraw from the start of the next block
	ctx := app.NewContext(true)

	// We export at last height + 1, because that's the height at which
	// Tendermint will start InitChain.
	height := app.LastBlockHeight() + 1
	if forZeroHeight {
		height = 0
		if err := app.prepForZeroHeightGenesis(ctx, jailAllowedAddrs); err != nil {
			return servertypes.ExportedApp{}, err
		}
	}

	genState, err := app.ModuleManager.ExportGenesisForModules(ctx, app.appCodec, modulesToExport)
	if err != nil {
		return servertypes.ExportedApp{}, err
	}

	appState, err := json.MarshalIndent(genState, "", "  ")
	if err != nil {
		return servertypes.ExportedApp{}, err
	}

	validators, err := app.GetValidatorSet(ctx)
	if err != nil {
		return servertypes.ExportedApp{}, err
	}

	return servertypes.ExportedApp{
		AppState:        appState,
		Validators:      validators,
		Height:          height,
		ConsensusParams: app.BaseApp.GetConsensusParams(ctx),
	}, err
}

// prepare for fresh start at zero height
// NOTE zero height genesis is a temporary feature which will be deprecated
// in favour of export at a block height
func (app *MitosisApp) prepForZeroHeightGenesis(ctx sdk.Context, jailAllowedAddrs []string) error {
	/* Just to be safe, assert the invariants on current state. */
	app.CrisisKeeper.AssertInvariants(ctx)

	// set context height to zero
	height := ctx.BlockHeight()
	ctx = ctx.WithBlockHeight(0)

	// reset context height
	ctx = ctx.WithBlockHeight(height)

	/* Handle slashing state. */

	// reset start height on signing infos
	var iteratorErr error
	if err := app.SlashingKeeper.IterateValidatorSigningInfos(
		ctx,
		func(addr sdk.ConsAddress, info slashingtypes.ValidatorSigningInfo) (stop bool) {
			info.StartHeight = 0
			if err := app.SlashingKeeper.SetValidatorSigningInfo(ctx, addr, info); err != nil {
				iteratorErr = err
				return true
			}
			return false
		},
	); err != nil {
		return err
	}
	if iteratorErr != nil {
		return iteratorErr
	}

	return nil
}

// GetValidatorSet returns a slice of bonded validators.
func (app *MitosisApp) GetValidatorSet(ctx sdk.Context) ([]tmtypes.GenesisValidator, error) {
	cVals := app.ConsumerKeeper.GetAllCCValidator(ctx)
	if len(cVals) == 0 {
		return nil, fmt.Errorf("empty validator set")
	}

	vals := []tmtypes.GenesisValidator{}
	for _, v := range cVals {
		vals = append(vals, tmtypes.GenesisValidator{Address: v.Address, Power: v.Power})
	}
	return vals, nil
}
