package app

import (
	"cosmossdk.io/depinject"
	"cosmossdk.io/log"
	evidencekeeper "cosmossdk.io/x/evidence/keeper"
	upgradekeeper "cosmossdk.io/x/upgrade/keeper"
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
	"io"
	"os"
	"path/filepath"
	"time"

	consensusparamskeeper "github.com/cosmos/cosmos-sdk/x/consensus/keeper"
	stakingkeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
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
	_ "github.com/cosmos/cosmos-sdk/x/staking"        // import for side-effects
)

var (
	DefaultNodeHome string
)

var (
	_ runtime.AppI            = (*MitosisApp)(nil)
	_ servertypes.Application = (*MitosisApp)(nil)
)

type MitosisApp struct {
	*runtime.App

	appCodec          codec.Codec
	txConfig          client.TxConfig
	interfaceRegistry types.InterfaceRegistry

	// keepers
	AccountKeeper         authkeeper.AccountKeeper
	BankKeeper            bankkeeper.Keeper
	StakingKeeper         *stakingkeeper.Keeper
	SlashingKeeper        slashingkeeper.Keeper
	UpgradeKeeper         *upgradekeeper.Keeper
	ConsensusParamsKeeper consensusparamskeeper.Keeper
	EvidenceKeeper        evidencekeeper.Keeper

	// EVM keepers
	EVMEngKeeper *evmengkeeper.Keeper
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
		&app.StakingKeeper,
		&app.SlashingKeeper,
		&app.EvidenceKeeper,
		&app.ConsensusParamsKeeper,
		&app.UpgradeKeeper,
		&app.EVMEngKeeper,
	); err != nil {
		return nil, errors.Wrap(err, "dep inject")
	}

	app.EVMEngKeeper.SetVoteProvider(NoVoteExtensionProvider{})
	// TODO(thai): make it configurable
	app.EVMEngKeeper.SetBuildDelay(time.Millisecond * 600) // 100ms longer than geth's --miner.recommit=500ms.
	app.EVMEngKeeper.SetBuildOptimistic(true)

	baseAppOpts = append(baseAppOpts, func(bapp *baseapp.BaseApp) {
		// Use evm engine to create block proposals.
		// Note that we do not check MaxTxBytes since all EngineEVM transaction MUST be included since we cannot
		// postpone them to the next block.
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

	app.EVMEngKeeper.SetCometAPI(NewCometAPI(clientCtx.Client))
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
