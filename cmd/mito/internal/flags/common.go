package flags

import (
	"github.com/spf13/cobra"
)

// CommonFlags contains all common flags used across commands
type CommonFlags struct {
	// Network flags
	Network                      string
	RPCURL                       string
	ChainID                      string
	ValidatorManagerContractAddr string

	// Signing flags
	PrivateKey          string
	KeyfilePath         string
	KeyfilePassword     string
	KeyfilePasswordFile string

	// Account-based signing flags
	Account      string
	KeystorePath string

	// Transaction flags
	GasLimit    uint64
	GasPrice    string
	Nonce       string
	ContractFee string

	// Output flags
	OutputFile string
	Signed     bool
	Unsigned   bool
	Yes        bool
}

// AddCommonFlags adds common flags to a command
func AddCommonFlags(cmd *cobra.Command, flags *CommonFlags) {
	// Network flags
	AddNetworkFlags(cmd, flags)

	// Signing flags
	AddSigningFlags(cmd, flags)

	// Transaction flags
	AddTransactionFlags(cmd, flags)

	// Output flags
	AddOutputFlags(cmd, flags)
}

// AddSigningFlags adds only signing-related flags
func AddSigningFlags(cmd *cobra.Command, flags *CommonFlags) {
	cmd.Flags().StringVar(&flags.PrivateKey, "private-key", "", "Private key for signing")
	cmd.Flags().StringVar(&flags.KeyfilePath, "keyfile", "", "Path to keyfile")
	cmd.Flags().StringVar(&flags.KeyfilePassword, "keyfile-password", "", "Password for keyfile")
	cmd.Flags().StringVar(&flags.KeyfilePasswordFile, "keyfile-password-file", "", "File containing keyfile password")

	// Account-based signing flags
	cmd.Flags().StringVar(&flags.Account, "account", "", "Account name from keystore (use with wallet new)")
	cmd.Flags().StringVar(&flags.KeystorePath, "keystore-dir", "", "Custom keystore directory path (default: ~/.mito/keystores)")
}

// AddNetworkFlags adds only network-related flags
func AddNetworkFlags(cmd *cobra.Command, flags *CommonFlags) {
	cmd.Flags().StringVar(&flags.Network, "network", "", "Network name for configuration")
	cmd.Flags().StringVar(&flags.RPCURL, "rpc-url", "", "RPC URL for Ethereum node")
	cmd.Flags().StringVar(&flags.ChainID, "chain-id", "", "Chain ID for the network")
	cmd.Flags().StringVar(&flags.ValidatorManagerContractAddr, "contract", "", "ValidatorManager contract address")
}

// AddTransactionFlags adds transaction-related flags
func AddTransactionFlags(cmd *cobra.Command, flags *CommonFlags) {
	cmd.Flags().Uint64Var(&flags.GasLimit, "gas-limit", 0, "Gas limit for transaction")
	cmd.Flags().StringVar(&flags.GasPrice, "gas-price", "", "Gas price for transaction")
	cmd.Flags().StringVar(&flags.Nonce, "nonce", "", "Nonce for transaction")
	cmd.Flags().StringVar(&flags.ContractFee, "contract-fee", "", "Contract fee")
}

func AddSendFlags(cmd *cobra.Command, flags *CommonFlags) {
	cmd.Flags().BoolVar(&flags.Yes, "yes", false, "Skip confirmation prompts")
}

func AddCreateFlags(cmd *cobra.Command, flags *CommonFlags) {
	cmd.Flags().BoolVar(&flags.Signed, "signed", false, "Create signed transaction")
	cmd.Flags().BoolVar(&flags.Unsigned, "unsigned", false, "Create unsigned transaction")
}

// AddOutputFlags adds output-related flags
func AddOutputFlags(cmd *cobra.Command, flags *CommonFlags) {
	cmd.Flags().StringVar(&flags.OutputFile, "output", "", "Output file for transaction")
}
