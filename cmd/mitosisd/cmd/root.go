package cmd

import (
	"context"
	"cosmossdk.io/client/v2/autocli"
	"cosmossdk.io/log"
	pvm "github.com/cometbft/cometbft/privval"
	dbm "github.com/cosmos/cosmos-db"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	"github.com/cosmos/cosmos-sdk/version"
	"github.com/mitosis-org/core/app/params"
	"github.com/omni-network/omni/lib/ethclient"
	evmengtypes "github.com/omni-network/omni/octane/evmengine/types"
	"os"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/types/tx/signing"
	"github.com/cosmos/cosmos-sdk/x/auth/tx"
	txmodule "github.com/cosmos/cosmos-sdk/x/auth/tx/config"

	"github.com/mitosis-org/core/app"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/config"
	"github.com/cosmos/cosmos-sdk/server"

	"github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/spf13/cobra"
)

const EnvPrefix = "MITO"

var (
	serverCtx    *server.Context
	addrProvider evmengtypes.AddressProvider
	engineCl     ethclient.EngineClient
)

func NewRootCmd() (*cobra.Command, params.EncodingConfig) {
	app.SetupConfig()

	tmpApp := newTempApp()
	encodingConfig := params.EncodingConfig{
		InterfaceRegistry: tmpApp.InterfaceRegistry(),
		Codec:             tmpApp.AppCodec(),
		TxConfig:          tmpApp.TxConfig(),
		Amino:             tmpApp.LegacyAmino(),
	}

	initClientCtx := client.Context{}.
		WithCodec(encodingConfig.Codec).
		WithInterfaceRegistry(encodingConfig.InterfaceRegistry).
		WithTxConfig(encodingConfig.TxConfig).
		WithLegacyAmino(encodingConfig.Amino).
		WithInput(os.Stdin).
		WithAccountRetriever(types.AccountRetriever{}).
		WithHomeDir(app.DefaultNodeHome).
		WithViper(EnvPrefix)

	rootCmd := &cobra.Command{
		Use:   version.AppName,
		Short: "Mitosis - Consensus Layer",
		PersistentPreRunE: func(cmd *cobra.Command, _ []string) error {
			cmd.SetOut(cmd.OutOrStdout())
			cmd.SetErr(cmd.ErrOrStderr())

			initClientCtx = initClientCtx.WithCmdContext(cmd.Context())
			initClientCtx, err := client.ReadPersistentCommandFlags(initClientCtx, cmd.Flags())
			if err != nil {
				return err
			}

			initClientCtx, err = config.ReadFromClientConfig(initClientCtx)
			if err != nil {
				return err
			}

			// This needs to go after ReadFromClientConfig, as that function
			// sets the RPC client needed for SIGN_MODE_TEXTUAL.
			enabledSignModes := append(tx.DefaultSignModes, signing.SignMode_SIGN_MODE_TEXTUAL)
			txConfigOpts := tx.ConfigOptions{
				EnabledSignModes:           enabledSignModes,
				TextualCoinMetadataQueryFn: txmodule.NewGRPCCoinMetadataQueryFn(initClientCtx),
			}
			txConfigWithTextual, err := tx.NewTxConfigWithOptions(
				codec.NewProtoCodec(encodingConfig.InterfaceRegistry),
				txConfigOpts,
			)
			if err != nil {
				return err
			}
			initClientCtx = initClientCtx.WithTxConfig(txConfigWithTextual)

			if err = client.SetCmdClientContextHandler(initClientCtx, cmd); err != nil {
				return err
			}

			customAppTemplate, customAppConfig := initAppConfig()
			customTMConfig := initTendermintConfig()

			if err = server.InterceptConfigsPreRunHandler(cmd, customAppTemplate, customAppConfig, customTMConfig); err != nil {
				return err
			}

			return initGlobalState(cmd)
		},
	}

	initRootCmd(rootCmd, encodingConfig, tmpApp.BasicModuleManager)

	autoCliOpts, err := enrichAutoCliOpts(tmpApp.AutoCliOpts(), initClientCtx)
	if err != nil {
		panic(err)
	}

	if err := autoCliOpts.EnhanceRootCommand(rootCmd); err != nil {
		panic(err)
	}

	return rootCmd, encodingConfig
}

func newTempApp() *app.MitosisApp {
	return app.NewMitosisApp(log.NewNopLogger(), dbm.NewMemDB(), nil, nil, nil, true, app.EmptyAppOptions{})
}

func enrichAutoCliOpts(autoCliOpts autocli.AppOptions, clientCtx client.Context) (autocli.AppOptions, error) {
	clientCtx, err := config.ReadFromClientConfig(clientCtx)
	if err != nil {
		return autocli.AppOptions{}, err
	}

	autoCliOpts.ClientCtx = clientCtx
	autoCliOpts.Keyring, err = keyring.NewAutoCLIKeyring(clientCtx.Keyring)
	if err != nil {
		return autocli.AppOptions{}, err
	}

	return autoCliOpts, nil
}

func initGlobalState(rootCmd *cobra.Command) error {
	var err error

	// serverCtx
	serverCtx = server.GetServerContextFromCmd(rootCmd)

	conf := DefaultAppConfig()
	if err = serverCtx.Viper.Unmarshal(&conf); err != nil {
		return err
	}

	// addrProvider
	addrProvider = &NoAddressProvider{}
	if conf.Engine.ValidatorMode {
		addrProvider, err = newAddrProvider(rootCmd)
		if err != nil {
			return err
		}
	}

	// engineCl
	engineCl, err = newEngineClient(rootCmd.Context(), conf.Engine)
	if err != nil {
		return err
	}

	return nil
}

func newAddrProvider(rootCmd *cobra.Command) (evmengtypes.AddressProvider, error) {
	serverCtx := server.GetServerContextFromCmd(rootCmd)

	cfg := serverCtx.Config

	privVal := pvm.LoadFilePV(cfg.PrivValidatorKeyFile(), cfg.PrivValidatorStateFile())

	pubKey, err := privVal.GetPubKey()
	if err != nil {
		return nil, err
	}

	return &SimpleAddressProvider{pubKey}, nil
}

func newEngineClient(ctx context.Context, engineCfg *EngineConfig) (ethclient.EngineClient, error) {
	if engineCfg.Mock {
		return ethclient.NewEngineMock()
	}

	jwtSecret, err := ethclient.LoadJWTHexFile(engineCfg.JWTFile)
	if err != nil {
		return nil, err
	}

	engineClient, err := ethclient.NewAuthClient(ctx, engineCfg.Endpoint, jwtSecret)
	if err != nil {
		return nil, err
	}

	return engineClient, nil
}
