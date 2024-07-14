package cmd

import (
	"cosmossdk.io/client/v2/autocli"
	"cosmossdk.io/log"
	dbm "github.com/cosmos/cosmos-db"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	"github.com/cosmos/cosmos-sdk/version"
	"github.com/mitosis-org/core/app/params"
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

			return server.InterceptConfigsPreRunHandler(cmd, customAppTemplate, customAppConfig, customTMConfig)
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
	return app.NewMitosisApp(log.NewNopLogger(), dbm.NewMemDB(), nil, true, app.EmptyAppOptions{})
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
