package cmd

import (
	"errors"
	"io"
	"time"

	"cosmossdk.io/log"
	confixcmd "cosmossdk.io/tools/confix/cmd"
	cmtcmd "github.com/cometbft/cometbft/cmd/cometbft/commands"
	dbm "github.com/cosmos/cosmos-db"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/debug"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/keys"
	"github.com/cosmos/cosmos-sdk/client/pruning"
	"github.com/cosmos/cosmos-sdk/client/rpc"
	"github.com/cosmos/cosmos-sdk/client/snapshot"
	"github.com/cosmos/cosmos-sdk/server"
	servertypes "github.com/cosmos/cosmos-sdk/server/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/cosmos/cosmos-sdk/version"
	authcmd "github.com/cosmos/cosmos-sdk/x/auth/client/cli"
	"github.com/cosmos/cosmos-sdk/x/crisis"
	genutilcli "github.com/cosmos/cosmos-sdk/x/genutil/client/cli"
	"github.com/mitosis-org/chain/app"
	evmvalcli "github.com/mitosis-org/chain/x/evmvalidator/client/cli"
	"github.com/spf13/cobra"
)

func initRootCmd(rootCmd *cobra.Command, txConfig client.TxConfig, basicManager module.BasicManager) {
	cfg := sdk.GetConfig()
	cfg.Seal()

	rootCmd.AddCommand(
		InitCmd(basicManager, app.DefaultNodeHome),
		evmvalcli.GetGenesisValidatorCmd(app.DefaultNodeHome),
		debug.Cmd(),
		confixcmd.ConfigCommand(),
		pruning.Cmd(newApp, app.DefaultNodeHome),
		snapshot.Cmd(newApp),
	)

	addServerCommands(rootCmd, app.DefaultNodeHome, newApp, appExport, addModuleInitFlags)

	rootCmd.AddCommand(
		server.StatusCommand(),
		genesisCommand(txConfig, app.DefaultNodeHome, basicManager),
		queryCommand(),
		txCommand(),
		keys.Commands(),
	)
}

func addServerCommands(rootCmd *cobra.Command, defaultNodeHome string, appCreator servertypes.AppCreator, appExport servertypes.AppExporter, addStartFlags servertypes.ModuleInitFlags) {
	cometCmd := &cobra.Command{
		Use:     "comet",
		Aliases: []string{"cometbft", "tendermint"},
		Short:   "CometBFT subcommands",
	}

	cometCmd.AddCommand(
		server.ShowNodeIDCmd(),
		server.ShowValidatorCmd(),
		server.ShowAddressCmd(),
		server.VersionCmd(),
		cmtcmd.ResetAllCmd,
		cmtcmd.ResetStateCmd,
		server.BootstrapStateCmd(appCreator),
	)

	startCmd := server.StartCmd(appCreator, defaultNodeHome)
	addStartFlags(startCmd)

	rootCmd.AddCommand(
		startCmd,
		cometCmd,
		server.ExportCmd(appExport, defaultNodeHome),
		version.NewVersionCommand(),
		server.NewRollbackCmd(appCreator, defaultNodeHome),
	)
}

func genesisCommand(txConfig client.TxConfig, defaultNodeHome string, basicManager module.BasicManager, cmds ...*cobra.Command) *cobra.Command {
	cmd := genutilcli.Commands(txConfig, basicManager, defaultNodeHome)

	for _, subCmd := range cmds {
		cmd.AddCommand(subCmd)
	}
	return cmd
}

func queryCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:                        "query",
		Aliases:                    []string{"q"},
		Short:                      "Querying subcommands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	cmd.AddCommand(
		rpc.QueryEventForTxCmd(),
		server.QueryBlockCmd(),
		authcmd.QueryTxsByEventsCmd(),
		server.QueryBlocksCmd(),
		authcmd.QueryTxCmd(),
		server.QueryBlockResultsCmd(),
	)

	return cmd
}

func txCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:                        "tx",
		Short:                      "Transactions subcommands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	cmd.AddCommand(
		authcmd.GetSignCommand(),
		authcmd.GetSignBatchCommand(),
		authcmd.GetMultiSignCommand(),
		authcmd.GetMultiSignBatchCmd(),
		authcmd.GetValidateSignaturesCommand(),
		authcmd.GetBroadcastCommand(),
		authcmd.GetEncodeCommand(),
		authcmd.GetDecodeCommand(),
		authcmd.GetSimulateCmd(),
	)

	return cmd
}

func addModuleInitFlags(startCmd *cobra.Command) {
	crisis.AddModuleInitFlags(startCmd)
}

func newApp(
	logger log.Logger,
	db dbm.DB,
	traceStore io.Writer,
	appOpts servertypes.AppOptions,
) servertypes.Application {
	baseappOptions := server.DefaultBaseappOptions(appOpts)

	appConfig, err := getAppConfig(runningCmd)
	if err != nil {
		panic(err)
	}

	engineCl, err := newEngineClient(runningCmd)
	if err != nil {
		panic(err)
	}

	addrProvider, err := newAddrProvider(runningCmd, appConfig.Engine.FeeRecipient)
	if err != nil {
		panic(err)
	}

	engineBuildDelay, err := time.ParseDuration(appConfig.Engine.BuildDelay)
	if err != nil {
		panic(err)
	}

	govEntrypointContractAddr, err := getGovEntrypointContractAddr(appConfig.EVMGov)
	if err != nil {
		panic(err)
	}

	mitosisApp, err := app.NewMitosisApp(
		logger,
		db,
		traceStore,
		engineCl,
		addrProvider,
		engineBuildDelay,
		appConfig.Engine.BuildOptimistic,
		govEntrypointContractAddr,
		true,
		appOpts,
		baseappOptions...,
	)
	if err != nil {
		panic(err)
	}

	return app.NewABCIWrappedApplication(mitosisApp)
}

func appExport(
	logger log.Logger,
	db dbm.DB,
	traceStore io.Writer,
	height int64,
	forZeroHeight bool,
	jailAllowedAddrs []string,
	appOpts servertypes.AppOptions,
	modulesToExport []string,
) (servertypes.ExportedApp, error) {
	var mitosisApp *app.MitosisApp

	homePath, ok := appOpts.Get(flags.FlagHome).(string)
	if !ok || homePath == "" {
		return servertypes.ExportedApp{}, errors.New("application home not set")
	}

	var loadLatest bool
	if height == -1 {
		loadLatest = true
	}

	appConfig, err := getAppConfig(runningCmd)
	if err != nil {
		return servertypes.ExportedApp{}, err
	}

	engineCl, err := newEngineClient(runningCmd)
	if err != nil {
		return servertypes.ExportedApp{}, err
	}

	addrProvider, err := newAddrProvider(runningCmd, appConfig.Engine.FeeRecipient)
	if err != nil {
		return servertypes.ExportedApp{}, err
	}

	engineBuildDelay, err := time.ParseDuration(appConfig.Engine.BuildDelay)
	if err != nil {
		return servertypes.ExportedApp{}, err
	}

	govEntrypointContractAddr, err := getGovEntrypointContractAddr(appConfig.EVMGov)
	if err != nil {
		return servertypes.ExportedApp{}, err
	}

	mitosisApp, err = app.NewMitosisApp(
		logger,
		db,
		traceStore,
		engineCl,
		addrProvider,
		engineBuildDelay,
		appConfig.Engine.BuildOptimistic,
		govEntrypointContractAddr,
		loadLatest,
		appOpts,
	)
	if err != nil {
		return servertypes.ExportedApp{}, err
	}

	if height != -1 {
		if err := mitosisApp.LoadHeight(height); err != nil {
			return servertypes.ExportedApp{}, err
		}
	}

	return mitosisApp.ExportAppStateAndValidators(forZeroHeight, jailAllowedAddrs, modulesToExport)
}
