package flags

import (
	"github.com/spf13/cobra"
)

// CommonFlags contains all common flags used across commands
type CommonFlags struct {
	// Network flags
	RpcURL                       string
	ChainID                      string
	ValidatorManagerContractAddr string

	// Signing flags
	PrivateKey          string
	KeyfilePath         string
	KeyfilePassword     string
	KeyfilePasswordFile string

	// Transaction flags
	GasLimit    uint64
	GasPrice    string
	Nonce       uint64
	NonceSet    bool
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
	cmd.Flags().StringVar(&flags.RpcURL, "rpc-url", "", "RPC URL for Ethereum node")
	cmd.Flags().StringVar(&flags.ChainID, "chain-id", "", "Chain ID for the network")
	cmd.Flags().StringVar(&flags.ValidatorManagerContractAddr, "contract", "", "ValidatorManager contract address")

	// Signing flags
	cmd.Flags().StringVar(&flags.PrivateKey, "private-key", "", "Private key for signing")
	cmd.Flags().StringVar(&flags.KeyfilePath, "keyfile", "", "Path to keyfile")
	cmd.Flags().StringVar(&flags.KeyfilePassword, "keyfile-password", "", "Password for keyfile")
	cmd.Flags().StringVar(&flags.KeyfilePasswordFile, "keyfile-password-file", "", "File containing keyfile password")

	// Transaction flags
	cmd.Flags().Uint64Var(&flags.GasLimit, "gas-limit", 0, "Gas limit for transaction")
	cmd.Flags().StringVar(&flags.GasPrice, "gas-price", "", "Gas price for transaction")
	cmd.Flags().Uint64Var(&flags.Nonce, "nonce", 0, "Nonce for transaction")
	cmd.Flags().StringVar(&flags.ContractFee, "contract-fee", "", "Contract fee")

	// Output flags
	cmd.Flags().StringVar(&flags.OutputFile, "output", "", "Output file for transaction")
	cmd.Flags().BoolVar(&flags.Signed, "signed", false, "Create signed transaction")
	cmd.Flags().BoolVar(&flags.Unsigned, "unsigned", false, "Create unsigned transaction")
	cmd.Flags().BoolVar(&flags.Yes, "yes", false, "Skip confirmation prompts")
}

// AddSigningFlags adds only signing-related flags
func AddSigningFlags(cmd *cobra.Command, flags *CommonFlags) {
	cmd.Flags().StringVar(&flags.PrivateKey, "private-key", "", "Private key for signing")
	cmd.Flags().StringVar(&flags.KeyfilePath, "keyfile", "", "Path to keyfile")
	cmd.Flags().StringVar(&flags.KeyfilePassword, "keyfile-password", "", "Password for keyfile")
	cmd.Flags().StringVar(&flags.KeyfilePasswordFile, "keyfile-password-file", "", "File containing keyfile password")
}

// AddNetworkFlags adds only network-related flags
func AddNetworkFlags(cmd *cobra.Command, flags *CommonFlags) {
	cmd.Flags().StringVar(&flags.RpcURL, "rpc-url", "", "RPC URL for Ethereum node")
	cmd.Flags().StringVar(&flags.ChainID, "chain-id", "", "Chain ID for the network")
	cmd.Flags().StringVar(&flags.ValidatorManagerContractAddr, "contract", "", "ValidatorManager contract address")
}

// AddTransactionFlags adds transaction-related flags
func AddTransactionFlags(cmd *cobra.Command, flags *CommonFlags) {
	cmd.Flags().Uint64Var(&flags.GasLimit, "gas-limit", 0, "Gas limit for transaction")
	cmd.Flags().StringVar(&flags.GasPrice, "gas-price", "", "Gas price for transaction")
	cmd.Flags().Uint64Var(&flags.Nonce, "nonce", 0, "Nonce for transaction")
	cmd.Flags().StringVar(&flags.ContractFee, "contract-fee", "", "Contract fee")
}

// AddOutputFlags adds output-related flags
func AddOutputFlags(cmd *cobra.Command, flags *CommonFlags) {
	cmd.Flags().StringVar(&flags.OutputFile, "output", "", "Output file for transaction")
	cmd.Flags().BoolVar(&flags.Signed, "signed", false, "Create signed transaction")
	cmd.Flags().BoolVar(&flags.Unsigned, "unsigned", false, "Create unsigned transaction")
	cmd.Flags().BoolVar(&flags.Yes, "yes", false, "Skip confirmation prompts")
}
