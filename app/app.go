package app

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"cosmossdk.io/x/evidence/exported"

	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/testutil/testdata/testpb"
	"github.com/cosmos/cosmos-sdk/x/auth/ante"
	"github.com/cosmos/cosmos-sdk/x/crisis"
	paramstypes "github.com/cosmos/cosmos-sdk/x/params/types"
	ibctesting "github.com/cosmos/ibc-go/v8/testing"
	ibctestingtypes "github.com/cosmos/ibc-go/v8/testing/types"
	appante "github.com/mitosis-org/chain/app/ante"
	"github.com/mitosis-org/chain/app/params"
	"github.com/omni-network/omni/lib/ethclient"
	"github.com/spf13/cast"

	storetypes "cosmossdk.io/store/types"
	"cosmossdk.io/x/evidence"
	evidencekeeper "cosmossdk.io/x/evidence/keeper"
	evidencetypes "cosmossdk.io/x/evidence/types"
	"cosmossdk.io/x/upgrade"
	upgradetypes "cosmossdk.io/x/upgrade/types"
	tmos "github.com/cometbft/cometbft/libs/os"

	"github.com/cosmos/cosmos-sdk/std"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/bank"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/cosmos/cosmos-sdk/x/consensus"
	consensusparamtypes "github.com/cosmos/cosmos-sdk/x/consensus/types"
	"github.com/cosmos/cosmos-sdk/x/slashing"
	slashingkeeper "github.com/cosmos/cosmos-sdk/x/slashing/keeper"
	slashingtypes "github.com/cosmos/cosmos-sdk/x/slashing/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"

	"cosmossdk.io/client/v2/autocli"
	"cosmossdk.io/core/appmodule"
	"github.com/cosmos/cosmos-sdk/client/grpc/cmtservice"
	"github.com/cosmos/cosmos-sdk/runtime"
	"github.com/cosmos/cosmos-sdk/types/msgservice"
	authcodec "github.com/cosmos/cosmos-sdk/x/auth/codec"
	"github.com/cosmos/cosmos-sdk/x/genutil"
	genutiltypes "github.com/cosmos/cosmos-sdk/x/genutil/types"
	"github.com/cosmos/cosmos-sdk/x/gov"
	govclient "github.com/cosmos/cosmos-sdk/x/gov/client"
	paramsclient "github.com/cosmos/cosmos-sdk/x/params/client"
	"github.com/cosmos/gogoproto/proto"

	autocliv1 "cosmossdk.io/api/cosmos/autocli/v1"
	reflectionv1 "cosmossdk.io/api/cosmos/reflection/v1"
	"cosmossdk.io/log"
	upgradekeeper "cosmossdk.io/x/upgrade/keeper"
	abci "github.com/cometbft/cometbft/abci/types"
	dbm "github.com/cosmos/cosmos-db"
	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/client"
	nodeservice "github.com/cosmos/cosmos-sdk/client/grpc/node"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/codec/types"
	runtimeservices "github.com/cosmos/cosmos-sdk/runtime/services"
	"github.com/cosmos/cosmos-sdk/server"
	"github.com/cosmos/cosmos-sdk/server/api"
	"github.com/cosmos/cosmos-sdk/server/config"
	servertypes "github.com/cosmos/cosmos-sdk/server/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/cosmos/cosmos-sdk/version"
	authkeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	authtx "github.com/cosmos/cosmos-sdk/x/auth/tx"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"

	consensusparamskeeper "github.com/cosmos/cosmos-sdk/x/consensus/keeper"
	crisiskeeper "github.com/cosmos/cosmos-sdk/x/crisis/keeper"
	crisistypes "github.com/cosmos/cosmos-sdk/x/crisis/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	"github.com/cosmos/ibc-go/modules/capability"
	capabilitykeeper "github.com/cosmos/ibc-go/modules/capability/keeper"
	capabilitytypes "github.com/cosmos/ibc-go/modules/capability/types"
	ibctransfer "github.com/cosmos/ibc-go/v8/modules/apps/transfer"
	ibctransferkeeper "github.com/cosmos/ibc-go/v8/modules/apps/transfer/keeper"
	ibctransfertypes "github.com/cosmos/ibc-go/v8/modules/apps/transfer/types"
	ibc "github.com/cosmos/ibc-go/v8/modules/core"
	porttypes "github.com/cosmos/ibc-go/v8/modules/core/05-port/types"
	ibcexported "github.com/cosmos/ibc-go/v8/modules/core/exported"
	ibckeeper "github.com/cosmos/ibc-go/v8/modules/core/keeper"
	ibctm "github.com/cosmos/ibc-go/v8/modules/light-clients/07-tendermint"
	ethostestutilintegration "github.com/ethos-works/ethos/ethos-chain/testutil/integration"
	testutil "github.com/ethos-works/ethos/ethos-chain/testutil/integration"
	ethosconsumermodule "github.com/ethos-works/ethos/ethos-chain/x/ccv/consumer"
	ethosconsumerkeeper "github.com/ethos-works/ethos/ethos-chain/x/ccv/consumer/keeper"
	ethosconsumertypes "github.com/ethos-works/ethos/ethos-chain/x/ccv/consumer/types"

	evmengkeeper "github.com/omni-network/omni/octane/evmengine/keeper"
	evmengmodule "github.com/omni-network/omni/octane/evmengine/module"
	evmengtypes "github.com/omni-network/omni/octane/evmengine/types"
)

var (
	DefaultNodeHome string

	// module account permissions
	maccPerms = map[string][]string{
		authtypes.FeeCollectorName:                      nil,
		ethosconsumertypes.ConsumerRedistributeName:     nil,
		ethosconsumertypes.ConsumerToSendToProviderName: nil,
		ibctransfertypes.ModuleName:                     {authtypes.Minter, authtypes.Burner},
	}
)

var (
	_ runtime.AppI                         = (*MitosisApp)(nil)
	_ servertypes.Application              = (*MitosisApp)(nil)
	_ ibctesting.TestingApp                = (*MitosisApp)(nil)
	_ ethostestutilintegration.ConsumerApp = (*MitosisApp)(nil)
)

type MitosisApp struct {
	*baseapp.BaseApp

	legacyAmino       *codec.LegacyAmino
	appCodec          codec.Codec
	txConfig          client.TxConfig
	interfaceRegistry types.InterfaceRegistry

	// keys to access the substores
	keys    map[string]*storetypes.KVStoreKey
	tkeys   map[string]*storetypes.TransientStoreKey
	memKeys map[string]*storetypes.MemoryStoreKey

	// keepers
	AccountKeeper         authkeeper.AccountKeeper
	BankKeeper            bankkeeper.Keeper
	CapabilityKeeper      *capabilitykeeper.Keeper
	SlashingKeeper        slashingkeeper.Keeper
	CrisisKeeper          *crisiskeeper.Keeper
	UpgradeKeeper         *upgradekeeper.Keeper
	ConsensusParamsKeeper consensusparamskeeper.Keeper
	EvidenceKeeper        evidencekeeper.Keeper

	// ibc keepers
	IBCKeeper      *ibckeeper.Keeper // IBC Keeper must be a pointer in the app, so we can SetRouter on it correctly
	TransferKeeper ibctransferkeeper.Keeper

	// ethos keepers
	ConsumerKeeper ethosconsumerkeeper.Keeper

	// EVM keepers
	EVMEngKeeper *evmengkeeper.Keeper

	// make scoped keepers public for test purposes
	ScopedIBCKeeper      capabilitykeeper.ScopedKeeper
	ScopedTransferKeeper capabilitykeeper.ScopedKeeper
	ScopedConsumerKeeper capabilitykeeper.ScopedKeeper

	// the module manager
	ModuleManager      *module.Manager
	BasicModuleManager module.BasicManager

	// simulation manager
	simulationManager *module.SimulationManager

	// module configurator
	configurator module.Configurator
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
	addrProvider evmengtypes.AddressProvider,
	loadLatest bool,
	appOpts servertypes.AppOptions,
	baseAppOptions ...func(*baseapp.BaseApp),
) *MitosisApp {
	encodingConfig := params.MakeEncodingConfig()
	interfaceRegistry := encodingConfig.InterfaceRegistry
	appCodec := encodingConfig.Codec
	txConfig := encodingConfig.TxConfig
	legacyAmino := encodingConfig.Amino

	std.RegisterLegacyAminoCodec(legacyAmino)
	std.RegisterInterfaces(interfaceRegistry)
	evmengtypes.RegisterInterfaces(interfaceRegistry)

	bApp := baseapp.NewBaseApp(version.AppName, logger, db, txConfig.TxDecoder(), baseAppOptions...)
	bApp.SetCommitMultiStoreTracer(traceStore)
	bApp.SetVersion(version.Version)
	bApp.SetInterfaceRegistry(interfaceRegistry)
	bApp.SetTxEncoder(txConfig.TxEncoder())

	// setup keys
	keys := storetypes.NewKVStoreKeys(
		authtypes.StoreKey, banktypes.StoreKey, capabilitytypes.StoreKey,
		slashingtypes.StoreKey, crisistypes.StoreKey, upgradetypes.StoreKey,
		consensusparamtypes.StoreKey, evidencetypes.StoreKey,
		ibcexported.StoreKey, ibctransfertypes.StoreKey,
		ethosconsumertypes.StoreKey,
		evmengtypes.StoreKey,
	)
	tkeys := storetypes.NewTransientStoreKeys()
	memKeys := storetypes.NewMemoryStoreKeys(capabilitytypes.MemStoreKey, evmengtypes.MemStoreKey)

	// register streaming services
	if err := bApp.RegisterStreamingServices(appOpts, keys); err != nil {
		panic(err)
	}

	app := &MitosisApp{
		BaseApp:           bApp,
		legacyAmino:       legacyAmino,
		appCodec:          appCodec,
		txConfig:          txConfig,
		interfaceRegistry: interfaceRegistry,
		keys:              keys,
		tkeys:             tkeys,
		memKeys:           memKeys,
	}

	accPrefix := sdk.GetConfig().GetBech32AccountAddrPrefix()
	valPrefix := sdk.GetConfig().GetBech32ValidatorAddrPrefix()
	consPrefix := sdk.GetConfig().GetBech32ConsensusAddrPrefix()
	accCodec := authcodec.NewBech32Codec(accPrefix)
	valCodec := authcodec.NewBech32Codec(valPrefix)
	consCodec := authcodec.NewBech32Codec(consPrefix)

	authorityAccAddr := authtypes.NewModuleAddress(govtypes.ModuleName)
	authorityAddr, err := accCodec.BytesToString(authorityAccAddr)
	if err != nil {
		panic(err)
	}

	// set the BaseApp's parameter store
	app.ConsensusParamsKeeper = consensusparamskeeper.NewKeeper(appCodec, runtime.NewKVStoreService(keys[consensusparamtypes.StoreKey]), authorityAddr, runtime.EventService{})
	bApp.SetParamStore(app.ConsensusParamsKeeper.ParamsStore)

	// add capability keeper and ScopeToModule for ibc module
	app.CapabilityKeeper = capabilitykeeper.NewKeeper(appCodec, keys[capabilitytypes.StoreKey], memKeys[capabilitytypes.MemStoreKey])

	// grant capabilities for modules
	app.ScopedIBCKeeper = app.CapabilityKeeper.ScopeToModule(ibcexported.ModuleName)
	app.ScopedTransferKeeper = app.CapabilityKeeper.ScopeToModule(ibctransfertypes.ModuleName)
	app.ScopedConsumerKeeper = app.CapabilityKeeper.ScopeToModule(ethosconsumertypes.ModuleName)

	// seal capability keeper after scoping modules
	// Applications that wish to enforce statically created ScopedKeepers should call `Seal` after creating
	// their scoped modules in `NewMitosisApp` with `ScopeToModule`
	app.CapabilityKeeper.Seal()

	// add keepers
	app.AccountKeeper = authkeeper.NewAccountKeeper(
		appCodec,
		runtime.NewKVStoreService(keys[authtypes.StoreKey]),
		authtypes.ProtoBaseAccount,
		maccPerms,
		accCodec,
		accPrefix,
		authorityAddr,
	)

	app.BankKeeper = bankkeeper.NewBaseKeeper(
		appCodec,
		runtime.NewKVStoreService(keys[banktypes.StoreKey]),
		app.AccountKeeper,
		BlockedAddrs(),
		authorityAddr,
		logger,
	)

	// consumer keeper satisfies the staking keeper interface
	// of the slashing module
	app.SlashingKeeper = slashingkeeper.NewKeeper(
		appCodec,
		legacyAmino,
		runtime.NewKVStoreService(keys[slashingtypes.StoreKey]),
		&app.ConsumerKeeper,
		authorityAddr,
	)

	invCheckPeriod := cast.ToUint(appOpts.Get(server.FlagInvCheckPeriod))
	app.CrisisKeeper = crisiskeeper.NewKeeper(
		appCodec,
		runtime.NewKVStoreService(keys[crisistypes.StoreKey]),
		invCheckPeriod,
		app.BankKeeper,
		authtypes.FeeCollectorName,
		authorityAddr,
		accCodec,
	)

	// create the upgrade keeper
	skipUpgradeHeights := map[int64]bool{}
	for _, h := range cast.ToIntSlice(appOpts.Get(server.FlagUnsafeSkipUpgrades)) {
		skipUpgradeHeights[int64(h)] = true
	}
	homePath := cast.ToString(appOpts.Get(flags.FlagHome))
	app.UpgradeKeeper = upgradekeeper.NewKeeper(
		skipUpgradeHeights,
		runtime.NewKVStoreService(keys[upgradetypes.StoreKey]),
		appCodec,
		homePath,
		app.BaseApp,
		authorityAddr,
	)

	////////////////////////////////////////////
	// IBC / Ethos
	////////////////////////////////////////////

	// pre-initialize ConsumerKeeper to satisfy ibckeeper.NewKeeper
	// which would panic on nil or zero keeper
	// ConsumerKeeper implements StakingKeeper but all function calls result in no-ops so this is safe
	// communication over IBC is not affected by these changes
	app.ConsumerKeeper = ethosconsumerkeeper.NewNonZeroKeeper(
		appCodec,
		runtime.NewKVStoreService(keys[ethosconsumertypes.StoreKey]),
	)

	app.IBCKeeper = ibckeeper.NewKeeper(
		appCodec,
		keys[ibcexported.StoreKey],
		nil,
		app.ConsumerKeeper,
		app.UpgradeKeeper,
		app.ScopedIBCKeeper,
		authorityAddr,
	)

	// initialize the actual consumer keeper
	app.ConsumerKeeper = ethosconsumerkeeper.NewKeeper(
		appCodec,
		runtime.NewKVStoreService(keys[ethosconsumertypes.StoreKey]),
		app.ScopedConsumerKeeper,
		app.IBCKeeper.ChannelKeeper,
		app.IBCKeeper.PortKeeper,
		app.IBCKeeper.ConnectionKeeper,
		app.IBCKeeper.ClientKeeper,
		app.SlashingKeeper,
		app.BankKeeper,
		app.AccountKeeper,
		&app.TransferKeeper,
		app.IBCKeeper,
		authtypes.FeeCollectorName,
		authorityAddr,
		valCodec,
		consCodec,
	)

	// register slashing module Slashing hooks to the consumer keeper
	app.ConsumerKeeper = *app.ConsumerKeeper.SetHooks(app.SlashingKeeper.Hooks())
	consumerModule := ethosconsumermodule.NewAppModule(app.ConsumerKeeper, paramstypes.Subspace{})

	// Create Transfer Keeper and pass IBCFeeKeeper as expected Channel and PortKeeper
	// since fee middleware will wrap the IBCKeeper for underlying application.
	app.TransferKeeper = ibctransferkeeper.NewKeeper(
		appCodec, keys[ibctransfertypes.StoreKey], nil,
		app.IBCKeeper.ChannelKeeper,
		app.IBCKeeper.ChannelKeeper, app.IBCKeeper.PortKeeper,
		app.AccountKeeper, app.BankKeeper, app.ScopedTransferKeeper,
		authorityAddr,
	)
	transferStack := ibctransfer.NewIBCModule(app.TransferKeeper)

	////////////////////////////////////////////
	// IBC router Configuration
	////////////////////////////////////////////

	ibcRouter := porttypes.NewRouter().
		AddRoute(ibctransfertypes.ModuleName, transferStack).
		AddRoute(ethosconsumertypes.ModuleName, consumerModule)

	app.IBCKeeper.SetRouter(ibcRouter)

	////////////////////////////////////////////
	// Evidence
	////////////////////////////////////////////

	// create evidence keeper with router
	app.EvidenceKeeper = *evidencekeeper.NewKeeper(
		appCodec,
		runtime.NewKVStoreService(keys[evidencetypes.StoreKey]),
		&app.ConsumerKeeper,
		app.SlashingKeeper,
		app.AccountKeeper.AddressCodec(),
		runtime.ProvideCometInfoService(),
	)

	router := evidencetypes.NewRouter()
	router = router.AddRoute(evidencetypes.RouteEquivocation, func(ctx context.Context, e exported.Evidence) error {
		slashFractionDoubleSign, err := app.SlashingKeeper.SlashFractionDoubleSign(ctx)
		if err != nil {
			return err
		}

		distributionHeight := e.GetHeight() - sdk.ValidatorUpdateDelay
		_, err = app.ConsumerKeeper.SlashWithInfractionReason(
			ctx,
			e.(*evidencetypes.Equivocation).GetConsensusAddress(app.ConsumerKeeper.ConsensusAddressCodec()),
			distributionHeight,
			e.(*evidencetypes.Equivocation).GetValidatorPower(),
			slashFractionDoubleSign,
			stakingtypes.Infraction_INFRACTION_DOUBLE_SIGN,
		)

		return err
	})
	app.EvidenceKeeper.SetRouter(router)

	///////////////////////////////////////////////////////////////////
	// Integration with Execution Layer through Engine API
	///////////////////////////////////////////////////////////////////

	evmEngKeeper, err := evmengkeeper.NewKeeper(
		appCodec,
		runtime.NewKVStoreService(keys[evmengtypes.StoreKey]),
		engineCl,
		txConfig,
		addrProvider,
		&burnEVMFees{}, // TODO(thai): give fees to block proposer
	)
	if err != nil {
		return nil
	}
	app.EVMEngKeeper = evmEngKeeper

	app.EVMEngKeeper.SetVoteProvider(NoVoteExtensionProvider{})
	// TODO(thai): make it configurable
	app.EVMEngKeeper.SetBuildDelay(time.Millisecond * 600) // 100ms longer than geth's --miner.recommit=500ms.
	app.EVMEngKeeper.SetBuildOptimistic(true)

	////////////////////////////////////////////
	// Module Options
	////////////////////////////////////////////

	// NOTE: we may consider parsing `appOpts` inside module constructors. For the moment
	// we prefer to be more strict in what arguments the modules expect.
	skipGenesisInvariants := cast.ToBool(appOpts.Get(crisis.FlagSkipGenesisInvariants))

	app.ModuleManager = module.NewManager(
		genutil.NewAppModule(app.AccountKeeper, app.ConsumerKeeper, app, txConfig),
		auth.NewAppModule(appCodec, app.AccountKeeper, nil, nil),
		bank.NewAppModule(appCodec, app.BankKeeper, app.AccountKeeper, nil),
		capability.NewAppModule(appCodec, *app.CapabilityKeeper, false),
		crisis.NewAppModule(app.CrisisKeeper, skipGenesisInvariants, nil),
		slashing.NewAppModule(appCodec, app.SlashingKeeper, app.AccountKeeper, app.BankKeeper, app.ConsumerKeeper, nil, app.interfaceRegistry),
		upgrade.NewAppModule(app.UpgradeKeeper, app.AccountKeeper.AddressCodec()),
		evidence.NewAppModule(app.EvidenceKeeper),
		consensus.NewAppModule(appCodec, app.ConsensusParamsKeeper),

		// IBC modules
		ibc.NewAppModule(app.IBCKeeper),
		ibctransfer.NewAppModule(app.TransferKeeper),
		ibctm.NewAppModule(),

		// Ethos modules
		consumerModule,

		// EVM modules
		evmengmodule.NewAppModule(appCodec, app.EVMEngKeeper),
	)

	// BasicModuleManager defines the module BasicManager is in charge of setting up basic,
	// non-dependant module elements, such as codec registration and genesis verification.
	// By default, it is composed of all the module from the module manager.
	// Additionally, app module basics can be overwritten by passing them as argument.
	app.BasicModuleManager = module.NewBasicManagerFromManager(
		app.ModuleManager,
		map[string]module.AppModuleBasic{
			genutiltypes.ModuleName: genutil.NewAppModuleBasic(genutiltypes.DefaultMessageValidator),
			govtypes.ModuleName: gov.NewAppModuleBasic(
				[]govclient.ProposalHandler{
					paramsclient.ProposalHandler,
				},
			),
		})
	app.BasicModuleManager.RegisterLegacyAminoCodec(legacyAmino)
	app.BasicModuleManager.RegisterInterfaces(interfaceRegistry)

	// NOTE: upgrade module is required to be prioritized
	app.ModuleManager.SetOrderPreBlockers(
		upgradetypes.ModuleName,
	)

	app.ModuleManager.SetOrderBeginBlockers(
		// NOTE: capability module's beginblocker must come before any modules using capabilities (e.g. IBC)
		capabilitytypes.ModuleName,
		slashingtypes.ModuleName,
		evidencetypes.ModuleName,
		ibcexported.ModuleName,
		ethosconsumertypes.ModuleName,
	)

	app.ModuleManager.SetOrderEndBlockers(
		crisistypes.ModuleName,
		slashingtypes.ModuleName,
		evidencetypes.ModuleName,
		ethosconsumertypes.ModuleName,
	)

	// NOTE: The genutils module must occur after staking so that pools are
	// properly initialized with tokens from genesis accounts.
	// NOTE: Capability module must occur first so that it can initialize any capabilities
	// so that other modules that want to create or claim capabilities afterwards in InitChain
	// can do so safely.
	genesisModuleOrder := []string{
		capabilitytypes.ModuleName, authtypes.ModuleName, banktypes.ModuleName,
		slashingtypes.ModuleName, crisistypes.ModuleName, genutiltypes.ModuleName, evidencetypes.ModuleName,
		upgradetypes.ModuleName, consensusparamtypes.ModuleName,
		ibcexported.ModuleName, ibctransfertypes.ModuleName,
		ethosconsumertypes.ModuleName,
		evmengtypes.ModuleName,
	}
	app.ModuleManager.SetOrderInitGenesis(genesisModuleOrder...)
	app.ModuleManager.SetOrderExportGenesis(genesisModuleOrder...)

	// register all module invariants
	app.ModuleManager.RegisterInvariants(app.CrisisKeeper)

	app.configurator = module.NewConfigurator(app.appCodec, app.MsgServiceRouter(), app.GRPCQueryRouter())
	if err := app.ModuleManager.RegisterServices(app.configurator); err != nil {
		panic(err)
	}

	// register upgrade handler for later use
	app.RegisterUpgradeHandlers(app.configurator)

	autocliv1.RegisterQueryServer(app.GRPCQueryRouter(), runtimeservices.NewAutoCLIQueryService(app.ModuleManager.Modules))

	reflectionSvc, err := runtimeservices.NewReflectionService()
	if err != nil {
		panic(err)
	}
	reflectionv1.RegisterReflectionServiceServer(app.GRPCQueryRouter(), reflectionSvc)

	// add test gRPC service for testing gRPC queries in isolation
	testpb.RegisterQueryServer(app.GRPCQueryRouter(), testpb.QueryImpl{})

	// initialize stores
	app.MountKVStores(keys)
	app.MountTransientStores(tkeys)
	app.MountMemoryStores(memKeys)

	// initialize BaseApp
	app.SetInitChainer(app.InitChainer)
	app.SetPreBlocker(app.PreBlocker)
	app.SetBeginBlocker(app.BeginBlocker)
	app.SetEndBlocker(app.EndBlocker)

	// set ante handler
	// app.setAnteHandler(txConfig) // TODO(thai): ethos need this but octane should not use this.

	app.SetPrepareProposal(app.EVMEngKeeper.PrepareProposal)
	app.SetProcessProposal(makeProcessProposalHandler(app, txConfig))

	// At startup, after all modules have been registered, check that all prot
	// annotations are correct.
	protoFiles, err := proto.MergedRegistry()
	if err != nil {
		panic(err)
	}
	err = msgservice.ValidateProtoAnnotations(protoFiles)
	if err != nil {
		// Once we switch to using protoreflect-based antehandlers, we might
		// want to panic here instead of logging a warning.
		if _, err = fmt.Fprintln(os.Stderr, err.Error()); err != nil {
			fmt.Println("could not write to stderr")
		}
	}

	if loadLatest {
		if err = app.LoadLatestVersion(); err != nil {
			tmos.Exit(fmt.Sprintf("failed to load latest version: %s", err.Error()))
		}
	}

	return app
}

func (app *MitosisApp) setAnteHandler(txConfig client.TxConfig) {
	anteHandler, err := appante.NewAnteHandler(
		appante.HandlerOptions{
			HandlerOptions: ante.HandlerOptions{
				AccountKeeper:   app.AccountKeeper,
				BankKeeper:      app.BankKeeper,
				FeegrantKeeper:  nil, // TODO(wip):
				SignModeHandler: txConfig.SignModeHandler(),
				SigGasConsumer:  ante.DefaultSigVerificationGasConsumer,
			},
			IBCKeeper:      app.IBCKeeper,
			ConsumerKeeper: app.ConsumerKeeper,
		},
	)
	if err != nil {
		panic(err)
	}

	// Set the AnteHandler for the app
	app.SetAnteHandler(anteHandler)
}

func (app *MitosisApp) Name() string { return app.BaseApp.Name() }

func (app *MitosisApp) PreBlocker(ctx sdk.Context, _ *abci.RequestFinalizeBlock) (*sdk.ResponsePreBlock, error) {
	return app.ModuleManager.PreBlock(ctx)
}

// BeginBlocker application updates every begin block
func (app *MitosisApp) BeginBlocker(ctx sdk.Context) (sdk.BeginBlock, error) {
	return app.ModuleManager.BeginBlock(ctx)
}

// EndBlocker application updates every end block
func (app *MitosisApp) EndBlocker(ctx sdk.Context) (sdk.EndBlock, error) {
	return app.ModuleManager.EndBlock(ctx)
}

// InitChainer application update at chain initialization
func (app *MitosisApp) InitChainer(ctx sdk.Context, req *abci.RequestInitChain) (*abci.ResponseInitChain, error) {
	var genesisState GenesisState
	if err := json.Unmarshal(req.AppStateBytes, &genesisState); err != nil {
		return nil, err
	}

	if err := app.UpgradeKeeper.SetModuleVersionMap(ctx, app.ModuleManager.GetVersionMap()); err != nil {
		return nil, err
	}
	return app.ModuleManager.InitGenesis(ctx, app.appCodec, genesisState)
}

func (app *MitosisApp) LoadHeight(height int64) error {
	return app.LoadVersion(height)
}

func (app *MitosisApp) LegacyAmino() *codec.LegacyAmino {
	return app.legacyAmino
}

func (app *MitosisApp) AppCodec() codec.Codec {
	return app.appCodec
}

func (app *MitosisApp) InterfaceRegistry() types.InterfaceRegistry {
	return app.interfaceRegistry
}

func (app *MitosisApp) TxConfig() client.TxConfig {
	return app.txConfig
}

// GetTxConfig satisfies the interface of ibctesting.TestingApp
func (app *MitosisApp) GetTxConfig() client.TxConfig {
	return app.TxConfig()
}

// AutoCliOpts returns the autocli options for the app.
func (app *MitosisApp) AutoCliOpts() autocli.AppOptions {
	modules := make(map[string]appmodule.AppModule, 0)
	for _, m := range app.ModuleManager.Modules {
		if moduleWithName, ok := m.(module.HasName); ok {
			moduleName := moduleWithName.Name()
			if appModule, ok := moduleWithName.(appmodule.AppModule); ok {
				modules[moduleName] = appModule
			}
		}
	}

	return autocli.AppOptions{
		Modules:               modules,
		ModuleOptions:         runtimeservices.ExtractAutoCLIOptions(app.ModuleManager.Modules),
		AddressCodec:          authcodec.NewBech32Codec(sdk.GetConfig().GetBech32AccountAddrPrefix()),
		ValidatorAddressCodec: authcodec.NewBech32Codec(sdk.GetConfig().GetBech32ValidatorAddrPrefix()),
		ConsensusAddressCodec: authcodec.NewBech32Codec(sdk.GetConfig().GetBech32ConsensusAddrPrefix()),
	}
}

// DefaultGenesis returns a default genesis from the registered AppModuleBasic's.
func (app *MitosisApp) DefaultGenesis() map[string]json.RawMessage {
	return app.BasicModuleManager.DefaultGenesis(app.appCodec)
}

// GetKey returns the KVStoreKey for the provided store key.
//
// NOTE: This is solely to be used for testing purposes.
func (app *MitosisApp) GetKey(storeKey string) *storetypes.KVStoreKey {
	return app.keys[storeKey]
}

// GetStoreKeys returns all the stored store keys.
func (app *MitosisApp) GetStoreKeys() []storetypes.StoreKey {
	keys := make([]storetypes.StoreKey, len(app.keys))
	for _, key := range app.keys {
		keys = append(keys, key)
	}

	return keys
}

// SimulationManager implements the SimulationApp interface
func (app *MitosisApp) SimulationManager() *module.SimulationManager {
	return app.simulationManager
}

func (app *MitosisApp) RegisterAPIRoutes(apiSvr *api.Server, apiConfig config.APIConfig) {
	clientCtx := apiSvr.ClientCtx

	// Register new tx routes from grpc-gateway.
	authtx.RegisterGRPCGatewayRoutes(clientCtx, apiSvr.GRPCGatewayRouter)

	// Register new CometBFT queries routes from grpc-gateway.
	cmtservice.RegisterGRPCGatewayRoutes(clientCtx, apiSvr.GRPCGatewayRouter)

	// Register node gRPC service for grpc-gateway.
	nodeservice.RegisterGRPCGatewayRoutes(clientCtx, apiSvr.GRPCGatewayRouter)

	// Register grpc-gateway routes for all modules.
	app.BasicModuleManager.RegisterGRPCGatewayRoutes(clientCtx, apiSvr.GRPCGatewayRouter)

	// register swagger API from root so that other applications can override easily
	if err := server.RegisterSwaggerAPI(apiSvr.ClientCtx, apiSvr.Router, apiConfig.Swagger); err != nil {
		panic(err)
	}
}

func (app *MitosisApp) RegisterTxService(clientCtx client.Context) {
	authtx.RegisterTxService(app.BaseApp.GRPCQueryRouter(), clientCtx, app.BaseApp.Simulate, app.interfaceRegistry)
}

// RegisterTendermintService implements the Application.RegisterTendermintService method.
func (app *MitosisApp) RegisterTendermintService(clientCtx client.Context) {
	cmtApp := server.NewCometABCIWrapper(app)
	cmtservice.RegisterTendermintService(
		clientCtx,
		app.BaseApp.GRPCQueryRouter(),
		app.interfaceRegistry,
		cmtApp.Query,
	)

	app.EVMEngKeeper.SetCometAPI(NewCometAPI(clientCtx.Client))
}

func (app *MitosisApp) RegisterNodeService(clientCtx client.Context, cfg config.Config) {
	nodeservice.RegisterNodeService(clientCtx, app.GRPCQueryRouter(), cfg)
}

func (app *MitosisApp) OnTxSucceeded(_ sdk.Context, _, _ string, _ []byte, _ []byte) {
}

func (app *MitosisApp) OnTxFailed(_ sdk.Context, _, _ string, _ []byte, _ []byte) {
}

func (app *MitosisApp) GetBaseApp() *baseapp.BaseApp {
	return app.BaseApp
}

// GetScopedIBCKeeper implements the TestingApp interface.
func (app *MitosisApp) GetScopedIBCKeeper() capabilitykeeper.ScopedKeeper {
	return app.ScopedIBCKeeper
}

// GetIBCKeeper implements the TestingApp interface.
func (app *MitosisApp) GetIBCKeeper() *ibckeeper.Keeper {
	return app.IBCKeeper
}

func (app *MitosisApp) GetMemKey(storeKey string) *storetypes.MemoryStoreKey {
	return app.memKeys[storeKey]
}

// ConsumerApp interface implementations for integration tests

// GetConsumerKeeper implements the ConsumerApp interface.
func (app *MitosisApp) GetConsumerKeeper() ethosconsumerkeeper.Keeper {
	return app.ConsumerKeeper
}

// GetTestBankKeeper implements the ConsumerApp interface.
func (app *MitosisApp) GetTestBankKeeper() testutil.TestBankKeeper {
	return app.BankKeeper
}

// GetTestAccountKeeper implements the ConsumerApp interface.
func (app *MitosisApp) GetTestAccountKeeper() testutil.TestAccountKeeper {
	return app.AccountKeeper
}

// GetTestSlashingKeeper implements the ConsumerApp interface.
func (app *MitosisApp) GetTestSlashingKeeper() testutil.TestSlashingKeeper {
	return app.SlashingKeeper
}

// GetTestEvidenceKeeper implements the ConsumerApp interface.
func (app *MitosisApp) GetTestEvidenceKeeper() evidencekeeper.Keeper {
	return app.EvidenceKeeper
}

func (app *MitosisApp) GetSubspace(moduleName string) paramstypes.Subspace {
	// TODO(thai): this function should be removed later.
	//   it is just for satisfying the interface. this function might not be called actually.
	panic("this function should be not called.")
}

// TestingApp functions

// GetStakingKeeper implements the TestingApp interface.
func (app *MitosisApp) GetStakingKeeper() ibctestingtypes.StakingKeeper {
	return app.ConsumerKeeper
}

type EmptyAppOptions struct{}

func (ao EmptyAppOptions) Get(_ string) interface{} {
	return nil
}

func BlockedAddrs() map[string]bool {
	modAccAddrs := make(map[string]bool)
	for acc := range maccPerms {
		modAccAddrs[authtypes.NewModuleAddress(acc).String()] = true
	}

	// this is required for the consumer chain to be able to send tokens to
	// the provider chain
	delete(modAccAddrs, authtypes.NewModuleAddress(ethosconsumertypes.ConsumerToSendToProviderName).String())

	return modAccAddrs
}
