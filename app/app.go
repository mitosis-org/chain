package app

import (
	"io"
	"os"
	"path/filepath"
	"time"

	mitotypes "github.com/mitosis-org/chain/types"

	"cosmossdk.io/depinject"
	"cosmossdk.io/log"
	evidencekeeper "cosmossdk.io/x/evidence/keeper"
	upgradekeeper "cosmossdk.io/x/upgrade/keeper"
	rpcclient "github.com/cometbft/cometbft/rpc/client"
	dbm "github.com/cosmos/cosmos-db"
	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/runtime"
	servertypes "github.com/cosmos/cosmos-sdk/server/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	authkeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	slashingkeeper "github.com/cosmos/cosmos-sdk/x/slashing/keeper"
	"github.com/omni-network/omni/lib/errors"
	"github.com/omni-network/omni/lib/ethclient"

	consensusparamskeeper "github.com/cosmos/cosmos-sdk/x/consensus/keeper"
	evmgovkeeper "github.com/mitosis-org/chain/x/evmgov/keeper"
	evmvalkeeper "github.com/mitosis-org/chain/x/evmvalidator/keeper"
	evmengkeeper "github.com/omni-network/omni/octane/evmengine/keeper"

	_ "cosmossdk.io/api/cosmos/tx/config/v1"          // import for side-effects
	_ "cosmossdk.io/x/evidence"                       // import for side-effects
	_ "cosmossdk.io/x/upgrade"                        // import for side-effects
	_ "github.com/cosmos/cosmos-sdk/x/auth"           // import for side-effects
	_ "github.com/cosmos/cosmos-sdk/x/auth/tx/config" // import for side-effects
	_ "github.com/cosmos/cosmos-sdk/x/bank"           // import for side-effects
	_ "github.com/cosmos/cosmos-sdk/x/consensus"      // import for side-effects
	_ "github.com/cosmos/cosmos-sdk/x/genutil"        // import for side-effects
	_ "github.com/cosmos/cosmos-sdk/x/slashing"       // import for side-effects
	_ "github.com/mitosis-org/chain/x/evmgov"         // import for side-effects
	_ "github.com/mitosis-org/chain/x/evmvalidator"   // import for side-effects
)

var DefaultNodeHome string

var (
	_ runtime.AppI            = (*MitosisApp)(nil)
	_ servertypes.Application = (*MitosisApp)(nil)
)

type MitosisApp struct {
	*runtime.App

	appCodec          codec.Codec
	txConfig          client.TxConfig
	interfaceRegistry types.InterfaceRegistry

	// Cosmos SDK keepers
	AccountKeeper         authkeeper.AccountKeeper
	BankKeeper            bankkeeper.Keeper
	SlashingKeeper        slashingkeeper.Keeper
	UpgradeKeeper         *upgradekeeper.Keeper
	ConsensusParamsKeeper consensusparamskeeper.Keeper
	EvidenceKeeper        evidencekeeper.Keeper

	// Octane keepers
	EVMEngKeeper *evmengkeeper.Keeper

	// Mitosis keepers
	EVMValKeeper *evmvalkeeper.Keeper
	EVMGovKeeper *evmgovkeeper.Keeper
}

func init() {
	userHomeDir, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}

	DefaultNodeHome = filepath.Join(userHomeDir, ".mitosisd")
}

func NewMitosisApp(
	logger log.Logger,
	db dbm.DB,
	traceStore io.Writer,
	engineCl ethclient.EngineClient,
	addrProvider ValidatorAddressProvider,
	engineBuildDelay time.Duration,
	engineBuildOptimistic bool,
	govEntrypointContractAddr mitotypes.EthAddress,
	loadLatest bool,
	appOpts servertypes.AppOptions,
	baseAppOpts ...func(*baseapp.BaseApp),
) (*MitosisApp, error) {
	var (
		app        = new(MitosisApp)
		appBuilder = new(runtime.AppBuilder)
	)
	if err := depinject.Inject(
		depinject.Configs(
			AppConfig(),
			depinject.Supply(
				logger,
				engineCl,
				&addrProvider,
				appOpts,
			),
		),
		&appBuilder,
		&app.appCodec,
		&app.txConfig,
		&app.interfaceRegistry,
		&app.AccountKeeper,
		&app.BankKeeper,
		&app.SlashingKeeper,
		&app.EvidenceKeeper,
		&app.ConsensusParamsKeeper,
		&app.UpgradeKeeper,
		&app.EVMEngKeeper,
		&app.EVMValKeeper,
		&app.EVMGovKeeper,
	); err != nil {
		return nil, errors.Wrap(err, "dep inject")
	}

	app.EVMEngKeeper.SetBuildDelay(engineBuildDelay)
	app.EVMEngKeeper.SetBuildOptimistic(engineBuildOptimistic)
	app.EVMEngKeeper.SetVoteProvider(NoVoteExtensionProvider{})

	app.EVMValKeeper.SetSlashingKeeper(app.SlashingKeeper)
	app.EVMValKeeper.SetEvmEngineKeeper(app.EVMEngKeeper)

	if err := app.EVMGovKeeper.SetGovEntrypointContractAddr(govEntrypointContractAddr); err != nil {
		return nil, errors.Wrap(err, "failed to set governance entrypoint contract address")
	}

	baseAppOpts = append(baseAppOpts, func(bapp *baseapp.BaseApp) {
		bapp.SetPrepareProposal(app.EVMEngKeeper.PrepareProposal)

		// Route proposed messages to keepers for verification and external state updates.
		bapp.SetProcessProposal(makeProcessProposalHandler(makeProcessProposalRouter(app), app.txConfig))
	})

	app.App = appBuilder.Build(db, traceStore, baseAppOpts...)

	if err := app.Load(loadLatest); err != nil {
		return nil, errors.Wrap(err, "failed to load latest version")
	}

	return app, nil
}

// RegisterTendermintService implements the Application.RegisterTendermintService method.
func (app *MitosisApp) RegisterTendermintService(clientCtx client.Context) {
	app.App.RegisterTendermintService(clientCtx)

	rpcClient, ok := clientCtx.Client.(rpcclient.Client)
	if !ok {
		panic("invalid rpc client")
	}

	app.EVMEngKeeper.SetCometAPI(NewCometAPI(rpcClient, app.ChainID()))
}

func (app *MitosisApp) LegacyAmino() *codec.LegacyAmino {
	return nil
}

func (app *MitosisApp) SimulationManager() *module.SimulationManager {
	return nil
}

type EmptyAppOptions struct{}

func (ao EmptyAppOptions) Get(_ string) interface{} {
	return nil
}
