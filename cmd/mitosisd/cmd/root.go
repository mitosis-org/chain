package cmd

import (
	"cosmossdk.io/client/v2/autocli"
	clientv2keyring "cosmossdk.io/client/v2/autocli/keyring"
	"cosmossdk.io/core/address"
	"cosmossdk.io/depinject"
	"cosmossdk.io/log"
	pvm "github.com/cometbft/cometbft/privval"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/cosmos/cosmos-sdk/version"
	"github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/omni-network/omni/lib/ethclient"
	"github.com/omni-network/omni/lib/k1util"
	"os"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/types/tx/signing"
	"github.com/cosmos/cosmos-sdk/x/auth/tx"
	txmodule "github.com/cosmos/cosmos-sdk/x/auth/tx/config"

	"github.com/mitosis-org/chain/app"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/config"
	"github.com/cosmos/cosmos-sdk/server"

	"github.com/spf13/cobra"
)

const EnvPrefix = "MITO"

var runningCmd *cobra.Command

func NewRootCmd() *cobra.Command {
	app.SetupConfig()

	var (
		txConfigOpts tx.ConfigOptions
		autoCliOpts  autocli.AppOptions
		basicManager module.BasicManager
		clientCtx    client.Context
	)

	mockEngineClient, err := ethclient.NewEngineMock()
	if err != nil {
		panic(err)
	}
	if err := depinject.Inject(
		depinject.Configs(
			app.AppConfig(),
			depinject.Supply(
				log.NewNopLogger(),
				mockEngineClient,
				&app.ValidatorAddressProvider{Addr: common.Address{}},
			),
			depinject.Provide(
				ProvideClientContext,
				ProvideKeyring,
			),
		),
		&txConfigOpts,
		&autoCliOpts,
		&basicManager,
		&clientCtx,
	); err != nil {
		panic(err)
	}

	rootCmd := &cobra.Command{
		Use:   version.AppName,
		Short: "Mitosis - Consensus Client",
		PersistentPreRunE: func(cmd *cobra.Command, _ []string) error {
			cmd.SetOut(cmd.OutOrStdout())
			cmd.SetErr(cmd.ErrOrStderr())

			clientCtx = clientCtx.WithCmdContext(cmd.Context())
			clientCtx, err := client.ReadPersistentCommandFlags(clientCtx, cmd.Flags())
			if err != nil {
				return err
			}

			clientCtx, err = config.ReadFromClientConfig(clientCtx)
			if err != nil {
				return err
			}

			// sign mode textual is only available in online mode
			if !clientCtx.Offline {
				// This needs to go after ReadFromClientConfig, as that function
				// sets the RPC client needed for SIGN_MODE_TEXTUAL.
				enabledSignModes := append(tx.DefaultSignModes, signing.SignMode_SIGN_MODE_TEXTUAL)
				txConfigOpts := tx.ConfigOptions{
					EnabledSignModes:           enabledSignModes,
					TextualCoinMetadataQueryFn: txmodule.NewGRPCCoinMetadataQueryFn(clientCtx),
				}
				txConfigWithTextual, err := tx.NewTxConfigWithOptions(
					codec.NewProtoCodec(clientCtx.InterfaceRegistry),
					txConfigOpts,
				)
				if err != nil {
					return err
				}
				clientCtx = clientCtx.WithTxConfig(txConfigWithTextual)
			}

			if err = client.SetCmdClientContextHandler(clientCtx, cmd); err != nil {
				return err
			}

			customAppTemplate, customAppConfig := initAppConfig()
			customTMConfig := initTendermintConfig()

			if err = server.InterceptConfigsPreRunHandler(cmd, customAppTemplate, customAppConfig, customTMConfig); err != nil {
				return err
			}

			runningCmd = cmd

			return nil
		},
	}

	initRootCmd(rootCmd, clientCtx.TxConfig, basicManager)

	if err := autoCliOpts.EnhanceRootCommand(rootCmd); err != nil {
		panic(err)
	}

	return rootCmd
}

func ProvideClientContext(
	appCodec codec.Codec,
	interfaceRegistry codectypes.InterfaceRegistry,
	txConfig client.TxConfig,
	legacyAmino *codec.LegacyAmino,
) client.Context {
	clientCtx := client.Context{}.
		WithCodec(appCodec).
		WithInterfaceRegistry(interfaceRegistry).
		WithTxConfig(txConfig).
		WithLegacyAmino(legacyAmino).
		WithInput(os.Stdin).
		WithAccountRetriever(types.AccountRetriever{}).
		WithHomeDir(app.DefaultNodeHome).
		WithViper(EnvPrefix) // env variable prefix

	// Read the config again to overwrite the default values with the values from the config file
	clientCtx, _ = config.ReadFromClientConfig(clientCtx)

	return clientCtx
}

func ProvideKeyring(clientCtx client.Context, addressCodec address.Codec) (clientv2keyring.Keyring, error) {
	kb, err := client.NewKeyringFromBackend(clientCtx, clientCtx.Keyring.Backend())
	if err != nil {
		return nil, err
	}

	return keyring.NewAutoCLIKeyring(kb)
}

func newAddrProvider(rootCmd *cobra.Command) (app.ValidatorAddressProvider, error) {
	serverCtx := server.GetServerContextFromCmd(rootCmd)

	cfg := serverCtx.Config

	privVal := pvm.LoadFilePV(cfg.PrivValidatorKeyFile(), cfg.PrivValidatorStateFile())

	addr, err := k1util.PubKeyToAddress(privVal.Key.PrivKey.PubKey())
	if err != nil {
		return app.ValidatorAddressProvider{}, err
	}

	return app.ValidatorAddressProvider{Addr: addr}, nil
}

func newEngineClient(rootCmd *cobra.Command) (ethclient.EngineClient, error) {
	serverCtx := server.GetServerContextFromCmd(rootCmd)

	conf := DefaultAppConfig()
	if err := serverCtx.Viper.Unmarshal(&conf); err != nil {
		return nil, err
	}

	if conf.Engine.Mock {
		return ethclient.NewEngineMock()
	}

	jwtSecret, err := ethclient.LoadJWTHexFile(conf.Engine.JWTFile)
	if err != nil {
		return nil, err
	}

	engineClient, err := ethclient.NewAuthClient(rootCmd.Context(), conf.Engine.Endpoint, jwtSecret)
	if err != nil {
		return nil, err
	}

	return engineClient, nil
}
