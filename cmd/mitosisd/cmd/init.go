package cmd

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"os"
	"path/filepath"

	"github.com/cosmos/cosmos-sdk/x/genutil"

	cfg "github.com/cometbft/cometbft/config"
	"github.com/cosmos/go-bip39"
	"github.com/spf13/cobra"

	errorsmod "cosmossdk.io/errors"
	"cosmossdk.io/math/unsafe"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/input"
	"github.com/cosmos/cosmos-sdk/server"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/cosmos/cosmos-sdk/version"
	"github.com/cosmos/cosmos-sdk/x/genutil/types"
)

const (
	// FlagOverwrite defines a flag to overwrite an existing genesis JSON file.
	FlagOverwrite = "overwrite"

	// FlagRecover defines a flag to recover the validator key from a BIP39 mnemonic.
	FlagRecover = "recover"

	// FlagDefaultBondDenom defines the default denom to use in the genesis file.
	FlagDefaultBondDenom = "default-denom"

	// Ethereum genesis flags
	FlagEthChainID     = "eth-chain-id"
	FlagEthGasLimit    = "eth-gas-limit"
	FlagEthFundedAddr  = "eth-funded-address"
	FlagEthInitBalance = "eth-initial-balance"
)

type printInfo struct {
	Moniker    string          `json:"moniker" yaml:"moniker"`
	ChainID    string          `json:"chain_id" yaml:"chain_id"`
	NodeID     string          `json:"node_id" yaml:"node_id"`
	GenTxsDir  string          `json:"gentxs_dir" yaml:"gentxs_dir"`
	AppMessage json.RawMessage `json:"app_message" yaml:"app_message"`
	EthGenesis string          `json:"eth_genesis_path" yaml:"eth_genesis_path"`
	EthChainID string          `json:"eth_chain_id" yaml:"eth_chain_id"`
}

func newPrintInfoWithEth(moniker, chainID, nodeID, genTxsDir string, appMessage json.RawMessage, ethGenesisPath string, ethChainID string) printInfo {
	return printInfo{
		Moniker:    moniker,
		ChainID:    chainID,
		NodeID:     nodeID,
		GenTxsDir:  genTxsDir,
		AppMessage: appMessage,
		EthGenesis: ethGenesisPath,
		EthChainID: ethChainID,
	}
}

func displayInfo(info printInfo) error {
	out, err := json.MarshalIndent(info, "", " ")
	if err != nil {
		return err
	}

	_, err = fmt.Fprintf(os.Stderr, "%s\n", out)

	return err
}

// InitCmd returns a command that initializes all files needed for Tendermint
// and the respective application.
func InitCmd(mbm module.BasicManager, defaultNodeHome string) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "init [moniker]",
		Short: "Initialize private validator, p2p, genesis, and application configuration files",
		Long:  `Initialize validators's and node's configuration files.`,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)
			cdc := clientCtx.Codec

			serverCtx := server.GetServerContextFromCmd(cmd)
			config := serverCtx.Config
			config.SetRoot(clientCtx.HomeDir)

			chainID, _ := cmd.Flags().GetString(flags.FlagChainID)
			switch {
			case chainID != "":
			case clientCtx.ChainID != "":
				chainID = clientCtx.ChainID
			default:
				chainID = fmt.Sprintf("test-chain-%v", unsafe.Str(6))
			}

			// Get bip39 mnemonic
			var mnemonic string
			recover, _ := cmd.Flags().GetBool(FlagRecover)
			if recover {
				inBuf := bufio.NewReader(cmd.InOrStdin())
				value, err := input.GetString("Enter your bip39 mnemonic", inBuf)
				if err != nil {
					return err
				}

				mnemonic = value
				if !bip39.IsMnemonicValid(mnemonic) {
					return errors.New("invalid mnemonic")
				}
			}

			// Get initial height
			initHeight, _ := cmd.Flags().GetInt64(flags.FlagInitHeight)
			if initHeight < 1 {
				initHeight = 1
			}

			nodeID, _, err := InitializeNodeValidatorFilesFromMnemonic(config, mnemonic)
			if err != nil {
				return err
			}

			config.Moniker = args[0]

			genFile := config.GenesisFile()
			overwrite, _ := cmd.Flags().GetBool(FlagOverwrite)
			defaultDenom, _ := cmd.Flags().GetString(FlagDefaultBondDenom)

			// use os.Stat to check if the file exists
			_, err = os.Stat(genFile)
			if !overwrite && !os.IsNotExist(err) {
				return fmt.Errorf("genesis.json file already exists: %v", genFile)
			}

			// Overwrites the SDK default denom for side-effects
			if defaultDenom != "" {
				sdk.DefaultBondDenom = defaultDenom
			}
			appGenState := mbm.DefaultGenesis(cdc)

			appState, err := json.MarshalIndent(appGenState, "", " ")
			if err != nil {
				return errorsmod.Wrap(err, "Failed to marshal default genesis state")
			}

			appGenesis := &types.AppGenesis{}
			if _, err := os.Stat(genFile); err != nil {
				if !os.IsNotExist(err) {
					return err
				}
			} else {
				appGenesis, err = types.AppGenesisFromFile(genFile)
				if err != nil {
					return errorsmod.Wrap(err, "Failed to read genesis doc from file")
				}
			}

			appGenesis.AppName = version.AppName
			appGenesis.AppVersion = version.Version
			appGenesis.ChainID = chainID
			appGenesis.AppState = appState
			appGenesis.InitialHeight = initHeight
			appGenesis.Consensus = &types.ConsensusGenesis{
				Validators: nil,
			}

			if err = genutil.ExportGenesisFile(appGenesis, genFile); err != nil {
				return errorsmod.Wrap(err, "Failed to export genesis file")
			}

			// Generate Ethereum genesis file with custom options
			ethGenesisPath := filepath.Join(config.RootDir, "config", "eth_genesis.json")

			// Get CLI flag values for Ethereum genesis customization
			ethChainIDFlag, _ := cmd.Flags().GetInt64(FlagEthChainID)
			ethGasLimit, _ := cmd.Flags().GetUint64(FlagEthGasLimit)
			ethFundedAddr, _ := cmd.Flags().GetString(FlagEthFundedAddr)
			ethInitBalance, _ := cmd.Flags().GetString(FlagEthInitBalance)

			opts := EthGenesisOptions{
				ChainID:        chainID,
				OutputPath:     ethGenesisPath,
				GasLimit:       ethGasLimit,
				FundedAddress:  ethFundedAddr,
				InitialBalance: ethInitBalance,
			}

			// Set custom Ethereum chain ID if provided
			if ethChainIDFlag > 0 {
				opts.EthChainID = big.NewInt(ethChainIDFlag)
			}

			if err = GenerateEthereumGenesisWithOptions(opts); err != nil {
				return errorsmod.Wrap(err, "Failed to generate Ethereum genesis file")
			}

			// Get Ethereum chain ID for display
			var ethChainID *big.Int
			if opts.EthChainID != nil {
				ethChainID = opts.EthChainID
			} else {
				ethChainID = GetEthChainIDFromCosmosChainID(chainID)
			}

			toPrint := newPrintInfoWithEth(config.Moniker, chainID, nodeID, "", appState, ethGenesisPath, ethChainID.String())

			cfg.WriteConfigFile(filepath.Join(config.RootDir, "config", "config.toml"), config)
			return displayInfo(toPrint)
		},
	}

	cmd.Flags().String(flags.FlagHome, defaultNodeHome, "node's home directory")
	cmd.Flags().BoolP(FlagOverwrite, "o", false, "overwrite the genesis.json file")
	cmd.Flags().Bool(FlagRecover, false, "provide seed phrase to recover existing key instead of creating")
	cmd.Flags().String(flags.FlagChainID, "", "genesis file chain-id, if left blank will be randomly created")
	cmd.Flags().String(FlagDefaultBondDenom, "", "genesis file default denomination, if left blank default value is 'stake'")
	cmd.Flags().Int64(flags.FlagInitHeight, 1, "specify the initial block height at genesis")

	// Ethereum genesis flags
	cmd.Flags().Int64(FlagEthChainID, 0, "ethereum chain ID (overrides default mapping)")
	cmd.Flags().Uint64(FlagEthGasLimit, 0, "ethereum genesis gas limit (default: 30000000)")
	cmd.Flags().String(FlagEthFundedAddr, "", "ethereum funded address (default: "+DefaultFundedAddress+")")
	cmd.Flags().String(FlagEthInitBalance, "", "ethereum initial balance in wei (default: 999000000000000000000000000)")

	return cmd
}
