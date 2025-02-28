package evmvalidator

import (
	"context"
	"cosmossdk.io/core/address"
	"cosmossdk.io/core/appmodule"
	"cosmossdk.io/depinject"
	"cosmossdk.io/depinject/appconfig"
	storetypes "cosmossdk.io/store/types"
	"encoding/json"
	abci "github.com/cometbft/cometbft/abci/types"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/omni-network/omni/lib/errors"

	"github.com/ethereum/go-ethereum/common"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/spf13/cobra"

	modulev1 "github.com/mitosis-org/chain/api/mitosis/evmvalidator/module/v1"
	"github.com/mitosis-org/chain/x/evmvalidator/keeper"
	"github.com/mitosis-org/chain/x/evmvalidator/types"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"
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
	_ module.HasGenesis      = (*AppModule)(nil)
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

// GetTxCmd returns the evmvalidator module's root tx command.
func (AppModuleBasic) GetTxCmd() *cobra.Command {
	return nil
}

// GetQueryCmd returns the evmvalidator module's root query command.
func (AppModuleBasic) GetQueryCmd() *cobra.Command {
	return nil
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
	// TODO(thai):
}

// RegisterGRPCGatewayRoutes registers the gRPC Gateway routes for the module.
func (AppModuleBasic) RegisterGRPCGatewayRoutes(client.Context, *runtime.ServeMux) {}

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
func (am AppModule) InitGenesis(ctx sdk.Context, cdc codec.JSONCodec, raw json.RawMessage) {
	var data types.GenesisState
	cdc.MustUnmarshalJSON(raw, &data)

	err := am.keeper.InitGenesis(ctx, &data)
	if err != nil {
		panic(errors.Wrap(err, "init genesis"))
	}
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
func (AppModule) RegisterServices(cfg module.Configurator) {
	// TODO(thai):
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
	)
}

type ModuleInputs struct {
	depinject.In

	Config                *modulev1.Module
	Cdc                   codec.Codec
	StoreKey              *storetypes.KVStoreKey
	SlashingKeeper        types.SlashingKeeper
	ValidatorAddressCodec address.Codec
	ConsensusAddressCodec address.Codec
}

type ModuleOutputs struct {
	depinject.Out

	Keeper       *keeper.Keeper
	Module       appmodule.AppModule
	EVMEventProc evmengtypes.InjectedEventProc
}

func ProvideModule(in ModuleInputs) (ModuleOutputs, error) {
	// Parse entrypoint address
	entrypointAddr := common.HexToAddress(in.Config.EvmValidatorEntrypointAddr)

	// Create keeper
	k := keeper.NewKeeperWithAddressCodecs(
		in.Cdc,
		in.StoreKey,
		in.SlashingKeeper,
		entrypointAddr,
		in.ValidatorAddressCodec,
		in.ConsensusAddressCodec,
	)

	// Create module
	m := NewAppModule(
		in.Cdc,
		k,
	)

	return ModuleOutputs{
		Keeper:       k,
		Module:       m,
		EVMEventProc: evmengtypes.InjectEventProc(k),
	}, nil
}
