package evmvalidator

import (
	"context"
	"encoding/json"
	"strings"

	"cosmossdk.io/core/appmodule"
	"cosmossdk.io/depinject"
	"cosmossdk.io/depinject/appconfig"
	storetypes "cosmossdk.io/store/types"

	abci "github.com/cometbft/cometbft/abci/types"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	grpcruntime "github.com/grpc-ecosystem/grpc-gateway/runtime"
	modulev1 "github.com/mitosis-org/chain/api/mitosis/evmvalidator/module/v1"
	"github.com/mitosis-org/chain/x/evmvalidator/client/cli"
	"github.com/mitosis-org/chain/x/evmvalidator/keeper"
	"github.com/mitosis-org/chain/x/evmvalidator/types"
	"github.com/omni-network/omni/lib/errors"
	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/runtime"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	evmengtypes "github.com/omni-network/omni/octane/evmengine/types"
)

const (
	ConsensusVersion = 1
)

var (
	_ module.AppModuleBasic  = (*AppModule)(nil)
	_ appmodule.AppModule    = (*AppModule)(nil)
	_ module.HasABCIGenesis  = (*AppModule)(nil)
	_ module.HasServices     = (*AppModule)(nil)
	_ module.HasABCIEndBlock = (*AppModule)(nil)
)

// ----------------------------------------------------------------------------
// AppModuleBasic
// ----------------------------------------------------------------------------

// AppModuleBasic implements the AppModuleBasic interface for the evmvalidator module.
type AppModuleBasic struct {
	cdc codec.BinaryCodec
}

func NewAppModuleBasic(cdc codec.BinaryCodec) AppModuleBasic {
	return AppModuleBasic{cdc: cdc}
}

// GetQueryCmd returns the evmvalidator module's root query command.
func (AppModuleBasic) GetQueryCmd() *cobra.Command {
	return cli.GetQueryCmd()
}

// Name returns the evmvalidator module's name.
func (AppModuleBasic) Name() string {
	return types.ModuleName
}

func (AppModuleBasic) ConsensusVersion() uint64 { return ConsensusVersion }

// RegisterLegacyAminoCodec registers the amino codec for the module, which is used
// to marshal and unmarshal structs to/from []byte in order to persist them in the module's KVStore.
func (AppModuleBasic) RegisterLegacyAminoCodec(*codec.LegacyAmino) {}

// RegisterInterfaces registers a module's interface types and their concrete implementations as proto.Message.
func (AppModuleBasic) RegisterInterfaces(reg codectypes.InterfaceRegistry) {
	types.RegisterInterfaces(reg)
}

// RegisterGRPCGatewayRoutes registers the gRPC Gateway routes for the module.
func (AppModuleBasic) RegisterGRPCGatewayRoutes(clientCtx client.Context, mux *grpcruntime.ServeMux) {
	if err := types.RegisterQueryHandlerClient(context.Background(), mux, types.NewQueryClient(clientCtx)); err != nil {
		panic(err)
	}
}

// ----------------------------------------------------------------------------
// AppModule
// ----------------------------------------------------------------------------

// AppModule implements the AppModule interface for the evmvalidator module.
type AppModule struct {
	AppModuleBasic

	keeper *keeper.Keeper
}

func NewAppModule(
	cdc codec.Codec,
	keeper *keeper.Keeper,
) AppModule {
	return AppModule{
		AppModuleBasic: NewAppModuleBasic(cdc),
		keeper:         keeper,
	}
}

// EndBlock executes all ABCI EndBlock logic for the evmvalidator module.
func (am AppModule) EndBlock(ctx context.Context) ([]abci.ValidatorUpdate, error) {
	return am.keeper.EndBlocker(sdk.UnwrapSDKContext(ctx))
}

// InitGenesis performs the evmvalidator module's genesis initialization
func (am AppModule) InitGenesis(ctx sdk.Context, cdc codec.JSONCodec, raw json.RawMessage) []abci.ValidatorUpdate {
	var data types.GenesisState
	cdc.MustUnmarshalJSON(raw, &data)

	vals, err := am.keeper.InitGenesis(ctx, &data)
	if err != nil {
		panic(errors.Wrap(err, "init genesis"))
	}
	return vals
}

// ExportGenesis returns the evmvalidator module's exported genesis state as raw JSON bytes.
func (am AppModule) ExportGenesis(ctx sdk.Context, cdc codec.JSONCodec) json.RawMessage {
	genState := am.keeper.ExportGenesis(ctx)
	return cdc.MustMarshalJSON(genState)
}

// DefaultGenesis returns the evmvalidator module's default genesis state.
func (am AppModule) DefaultGenesis(cdc codec.JSONCodec) json.RawMessage {
	return cdc.MustMarshalJSON(types.DefaultGenesisState())
}

// ValidateGenesis performs validation of the evmvalidator module's genesis state.
func (am AppModule) ValidateGenesis(cdc codec.JSONCodec, _ client.TxEncodingConfig, bz json.RawMessage) error {
	var data types.GenesisState
	if err := cdc.UnmarshalJSON(bz, &data); err != nil {
		return errors.Wrap(err, "unmarshal genesis state")
	}

	err := data.Validate()
	if err != nil {
		return errors.Wrap(err, "validate genesis state")
	}

	return nil
}

// RegisterServices registers a gRPC query service to respond to the module-specific gRPC queries.
func (am AppModule) RegisterServices(cfg module.Configurator) {
	types.RegisterMsgServer(cfg.MsgServer(), keeper.NewMsgServerImpl(am.keeper))
	types.RegisterQueryServer(cfg.QueryServer(), keeper.NewQueryServer(am.keeper))
}

// IsOnePerModuleType implements the depinject.OnePerModuleType interface.
func (AppModule) IsOnePerModuleType() {}

// IsAppModule implements the appmodule.AppModule interface.
func (AppModule) IsAppModule() {}

// ----------------------------------------------------------------------------
// App Wiring Setup
// ----------------------------------------------------------------------------

//nolint:gochecknoinits
func init() {
	appconfig.RegisterModule(
		&modulev1.Module{},
		appconfig.Provide(ProvideModule),
		appmodule.Invoke(InvokeProhibitStakingHooks),
	)
}

type ModuleInputs struct {
	depinject.In

	Config                *modulev1.Module
	Cdc                   codec.Codec
	StoreKey              *storetypes.KVStoreKey
	ValidatorAddressCodec runtime.ValidatorAddressCodec
	ConsensusAddressCodec runtime.ConsensusAddressCodec
}

type ModuleOutputs struct {
	depinject.Out

	Module            appmodule.AppModule
	Keeper            *keeper.Keeper
	KeeperForEvidence *keeper.KeeperWrapperForEvidence
	EVMEventProc      evmengtypes.InjectedEventProc
}

func ProvideModule(in ModuleInputs) (ModuleOutputs, error) {
	// Parse authority address
	if in.Config.Authority == "" {
		return ModuleOutputs{}, errors.New("authority address is empty")
	}
	authority := authtypes.NewModuleAddressOrBech32Address(in.Config.Authority)

	// Create keeper
	k := keeper.NewKeeper(
		in.Cdc,
		in.StoreKey,
		in.ValidatorAddressCodec,
		in.ConsensusAddressCodec,
		authority.String(),
	)

	// Create module
	m := NewAppModule(
		in.Cdc,
		k,
	)

	return ModuleOutputs{
		Module:            m,
		Keeper:            k,
		KeeperForEvidence: &keeper.KeeperWrapperForEvidence{K: k},
		EVMEventProc:      evmengtypes.InjectEventProc(k),
	}, nil
}

// InvokeProhibitStakingHooks is an invoker that prohibits the use of staking hooks.
// x/evmvalidator is compatible with x/staking partially, but it does not support staking hooks.
// So, other modules should not use staking hooks with x/evmvalidator.
func InvokeProhibitStakingHooks(
	stakingHooks map[string]stakingtypes.StakingHooksWrapper,
) error {
	if len(stakingHooks) == 0 {
		return nil
	}

	modules := make([]string, 0, len(stakingHooks))
	for k := range stakingHooks {
		modules = append(modules, k)
	}
	modulesStr := strings.Join(modules, ",")

	return errors.New("staking hooks are not supported", "modules", modulesStr)
}
