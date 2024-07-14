package integration

import (
	"cosmossdk.io/log"
	"encoding/json"
	cometbfttypes "github.com/cometbft/cometbft/abci/types"
	dbm "github.com/cosmos/cosmos-db"
	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	ibctesting "github.com/cosmos/ibc-go/v8/testing"
	providerApp "github.com/ethos-works/ethos/ethos-chain/testapps/provider"
	"github.com/ethos-works/ethos/ethos-chain/tests/integration"
	icstestingutils "github.com/ethos-works/ethos/ethos-chain/testutil/ibc_testing"
	consumertypes "github.com/ethos-works/ethos/ethos-chain/x/ccv/consumer/types"
	ccvtypes "github.com/ethos-works/ethos/ethos-chain/x/ccv/types"
	mitosisApp "github.com/mitosis-org/core/app"
	"github.com/stretchr/testify/suite"
	"testing"
)

var (
	ccvSuite *integration.CCVTestSuite
)

func init() {
	ccvSuite = integration.NewCCVTestSuite[*providerApp.App, *mitosisApp.MitosisApp](
		icstestingutils.ProviderAppIniter, SetupValSetAppIniter, []string{})
}

// TODO(thai): `mitosisApp.SetupConfig()` should be called
func TestCCVTestSuite(t *testing.T) {
	suite.Run(t, ccvSuite)
}

// SetupValSetAppIniter is a wrapper for e2e tests to satisfy test interface
func SetupValSetAppIniter(initValUpdates []cometbfttypes.ValidatorUpdate) icstestingutils.AppIniter {
	return SetupTestingApp(initValUpdates)
}

func SetupTestingApp(initValUpdates []cometbfttypes.ValidatorUpdate) func() (ibctesting.TestingApp, map[string]json.RawMessage) {
	return func() (ibctesting.TestingApp, map[string]json.RawMessage) {
		app := mitosisApp.NewMitosisApp(
			log.NewNopLogger(),
			dbm.NewMemDB(),
			nil,
			false,
			mitosisApp.EmptyAppOptions{},
		)
		encoding := app.AppCodec()

		app.SetInitChainer(app.InitChainer)
		if err := app.LoadLatestVersion(); err != nil {
			panic(err)
		}

		genesisState := app.DefaultGenesis()

		genesisState[stakingtypes.ModuleName] = encoding.MustMarshalJSON(
			&stakingtypes.GenesisState{
				Params: stakingtypes.Params{BondDenom: sdk.DefaultBondDenom},
			},
		)

		var consumerGenesis ccvtypes.ConsumerGenesisState
		encoding.MustUnmarshalJSON(genesisState[consumertypes.ModuleName], &consumerGenesis)
		consumerGenesis.Provider.InitialValSet = initValUpdates
		consumerGenesis.Params.Enabled = true
		genesisState[consumertypes.ModuleName] = encoding.MustMarshalJSON(&consumerGenesis)

		return app, genesisState
	}
}
