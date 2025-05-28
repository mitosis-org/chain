package cmd

import (
	"bufio"
	"context"
	"crypto/ecdsa"
	"fmt"
	"io/ioutil"
	"log"
	"math/big"
	"os"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/common"
	ethcrypto "github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/mitosis-org/chain/bindings"
	"github.com/spf13/cobra"
	"golang.org/x/term"
)

// Common flags used across multiple commands
var (
	rpcURL                       string
	privateKey                   string
	keyfilePath                  string
	keyfilePassword              string
	keyfilePasswordFile          string
	validatorManagerContractAddr string
	yes                          bool
	nonce                        uint64
	nonceSpecified               bool
	outputFile                   string
	signed                       bool
	unsigned                     bool

	// Network information flags
	chainID     string
	gasLimit    uint64
	gasPrice    string
	txNonce     uint64
	contractFee string

	// Shared client and contract instances
	ethClient *ethclient.Client
	contract  *bindings.IValidatorManager

	// Global config
	globalConfig *Config
)

// SigningMethod represents the method used for signing transactions
type SigningMethod int

// Signing method types
const (
	SigningMethodPrivateKey SigningMethod = iota
	SigningMethodKeyfile
)

// SigningConfig holds the signing configuration
type SigningConfig struct {
	Method          SigningMethod
	PrivateKey      *ecdsa.PrivateKey
	KeyfilePath     string
	KeyfilePassword string
}

// MutuallyExclusiveGroup represents a group of mutually exclusive flags
type MutuallyExclusiveGroup struct {
	Name        string
	Description string
	Flags       []string
	Required    bool
}

// FlagValidator manages validation of mutually exclusive flags
type FlagValidator struct {
	groups []MutuallyExclusiveGroup
}

// NewFlagValidator creates a new flag validator
func NewFlagValidator() *FlagValidator {
	return &FlagValidator{
		groups: make([]MutuallyExclusiveGroup, 0),
	}
}

// AddMutuallyExclusiveGroup adds a mutually exclusive group
func (fv *FlagValidator) AddMutuallyExclusiveGroup(group MutuallyExclusiveGroup) {
	fv.groups = append(fv.groups, group)
}

// ValidateFlags validates all mutually exclusive groups for a command
func (fv *FlagValidator) ValidateFlags(cmd *cobra.Command) error {
	for _, group := range fv.groups {
		setFlags := make([]string, 0)

		// Check which flags in the group are set
		for _, flagName := range group.Flags {
			if cmd.Flags().Changed(flagName) {
				setFlags = append(setFlags, flagName)
			}
		}

		// Validate mutual exclusivity
		if len(setFlags) > 1 {
			return fmt.Errorf("flags %v are mutually exclusive (from group: %s)", setFlags, group.Name)
		}

		// Validate required constraint
		if group.Required && len(setFlags) == 0 {
			return fmt.Errorf("one of the following flags is required: %v (group: %s)", group.Flags, group.Name)
		}
	}

	return nil
}

// Common signing method groups
var (
	SigningMethodGroup = MutuallyExclusiveGroup{
		Name:        "signing-method",
		Description: "Method for signing transactions",
		Flags:       []string{"private-key", "keyfile"},
		Required:    false, // Will be set to true for commands that require signing
	}

	// Transaction type group (for future use)
	TransactionTypeGroup = MutuallyExclusiveGroup{
		Name:        "transaction-type",
		Description: "Type of transaction to create",
		Flags:       []string{"signed", "unsigned"},
		Required:    false,
	}

	// Output format group (for future use)
	OutputFormatGroup = MutuallyExclusiveGroup{
		Name:        "output-format",
		Description: "Output format for transaction data",
		Flags:       []string{"json", "raw", "hex"},
		Required:    false,
	}

	// Network information group for offline mode
	NetworkInfoGroup = MutuallyExclusiveGroup{
		Name:        "network-info",
		Description: "Network information for offline transaction creation",
		Flags:       []string{"rpc-url", "chain-id"},
		Required:    false, // Will be validated conditionally
	}
)

// mustMarkFlagRequired marks a flag as required and panics if it fails
func mustMarkFlagRequired(cmd *cobra.Command, flag string) {
	if err := cmd.MarkFlagRequired(flag); err != nil {
		log.Fatalf("Failed to mark flag '%s' as required: %v", flag, err)
	}
}

// AddCommonFlags adds common flags to a command
func AddCommonFlags(cmd *cobra.Command, requireSigning bool) {
	cmd.Flags().StringVar(&rpcURL, "rpc-url", "", "Ethereum RPC URL (if not set, uses config)")
	cmd.Flags().StringVar(&validatorManagerContractAddr, "contract", "", "ValidatorManager contract address (if not set, uses config)")
	cmd.Flags().BoolVar(&yes, "yes", false, "Skip confirmation prompt")
	cmd.Flags().Uint64Var(&nonce, "nonce", 0, "Manually specify nonce for transaction (optional)")

	if requireSigning {
		// Signing options
		cmd.Flags().StringVar(&privateKey, "private-key", "", "Private key for signing transactions (hex format)")
		cmd.Flags().StringVar(&keyfilePath, "keyfile", "", "Path to geth keyfile")
		cmd.Flags().StringVar(&keyfilePassword, "keyfile-password", "", "Password for keyfile")
		cmd.Flags().StringVar(&keyfilePasswordFile, "keyfile-password-file", "", "Path to file containing keyfile password")
	}

	// Preserve any existing PreRun function
	existingPreRun := cmd.PreRun
	cmd.PreRun = func(cmd *cobra.Command, args []string) {
		var err error

		// Load global config first
		if err := loadGlobalConfig(); err != nil {
			fmt.Printf("Warning: failed to load config: %v\n", err)
		}

		// Resolve config values (command line flags take precedence)
		resolveConfigValues()

		// Validate that required values are set
		if rpcURL == "" {
			fmt.Println("Error: RPC URL is required. Set it with --rpc-url flag or use 'mito config set-rpc <url>'")
			os.Exit(1)
		}
		if validatorManagerContractAddr == "" {
			fmt.Println("Error: ValidatorManager contract address is required. Set it with --contract flag or use 'mito config set-contract <address>'")
			os.Exit(1)
		}

		// Set the nonceSpecified flag
		nonceSpecified = cmd.Flags().Changed("nonce")

		// Setup client
		ethClient, err = GetEthClient(rpcURL)
		if err != nil {
			log.Fatalf("Failed to connect to Ethereum client: %v", err)
		}

		// Get contract instance
		contract, err = GetValidatorManagerContract(ethClient)
		if err != nil {
			log.Fatalf("Failed to initialize contract: %v", err)
		}

		// Call the existing PreRun if it exists
		if existingPreRun != nil {
			existingPreRun(cmd, args)
		}
	}
}

// AddTxCreateFlags adds flags specific to tx create commands with optional signing validation
func AddTxCreateFlags(cmd *cobra.Command) {
	cmd.Flags().BoolVar(&signed, "signed", false, "Create a signed transaction")
	cmd.Flags().BoolVar(&unsigned, "unsigned", true, "Create an unsigned transaction (default)")
	cmd.Flags().StringVar(&outputFile, "output", "", "Output file for the transaction (default: stdout)")

	// Add signing flags when --signed is used
	cmd.Flags().StringVar(&privateKey, "private-key", "", "Private key for signing transactions (hex format) [mutually exclusive with --keyfile]")
	cmd.Flags().StringVar(&keyfilePath, "keyfile", "", "Path to geth keyfile [mutually exclusive with --private-key]")
	cmd.Flags().StringVar(&keyfilePassword, "keyfile-password", "", "Password for keyfile")
	cmd.Flags().StringVar(&keyfilePasswordFile, "keyfile-password-file", "", "Path to file containing keyfile password")

	// Add PreRun to load config and validate flags
	existingPreRun := cmd.PreRun
	cmd.PreRun = func(cmd *cobra.Command, args []string) {
		// Validate mutually exclusive flags
		privateKeySet := cmd.Flags().Changed("private-key")
		keyfileSet := cmd.Flags().Changed("keyfile")

		if privateKeySet && keyfileSet {
			fmt.Println("Error: flags --private-key and --keyfile are mutually exclusive")
			fmt.Println("\nUsage:")
			fmt.Println("  Use either --private-key OR --keyfile, not both")
			fmt.Println("  --private-key: Provide private key directly (hex format)")
			fmt.Println("  --keyfile: Use geth keyfile (more secure)")
			os.Exit(1)
		}

		// Check if signing is required
		if signed {
			if !privateKeySet && !keyfileSet {
				fmt.Println("Error: When using --signed, you must provide either --private-key OR --keyfile")
				fmt.Println("\nUsage:")
				fmt.Println("  --private-key: Provide private key directly (hex format)")
				fmt.Println("  --keyfile: Use geth keyfile (more secure)")
				os.Exit(1)
			}
		}

		// Load global config
		if err := loadGlobalConfig(); err != nil {
			fmt.Printf("Warning: failed to load config: %v\n", err)
		}

		// Resolve config values
		resolveConfigValues()

		// Call existing PreRun if it exists
		if existingPreRun != nil {
			existingPreRun(cmd, args)
		}
	}
}

// AddSigningFlags adds signing-related flags to a command with mutual exclusivity validation
func AddSigningFlags(cmd *cobra.Command) {
	cmd.Flags().StringVar(&privateKey, "private-key", "", "Private key for signing transactions (hex format) [mutually exclusive with --keyfile]")
	cmd.Flags().StringVar(&keyfilePath, "keyfile", "", "Path to geth keyfile [mutually exclusive with --private-key]")
	cmd.Flags().StringVar(&keyfilePassword, "keyfile-password", "", "Password for keyfile")
	cmd.Flags().StringVar(&keyfilePasswordFile, "keyfile-password-file", "", "Path to file containing keyfile password")

	// Add PreRun to load config and validate flags
	existingPreRun := cmd.PreRun
	cmd.PreRun = func(cmd *cobra.Command, args []string) {
		// Validate mutually exclusive flags first
		privateKeySet := cmd.Flags().Changed("private-key")
		keyfileSet := cmd.Flags().Changed("keyfile")

		if privateKeySet && keyfileSet {
			fmt.Println("Error: flags --private-key and --keyfile are mutually exclusive")
			fmt.Println("\nUsage:")
			fmt.Println("  Use either --private-key OR --keyfile, not both")
			fmt.Println("  --private-key: Provide private key directly (hex format)")
			fmt.Println("  --keyfile: Use geth keyfile (more secure)")
			os.Exit(1)
		}

		if !privateKeySet && !keyfileSet {
			fmt.Println("Error: one of the following flags is required: --private-key, --keyfile")
			fmt.Println("\nUsage:")
			fmt.Println("  Use either --private-key OR --keyfile")
			fmt.Println("  --private-key: Provide private key directly (hex format)")
			fmt.Println("  --keyfile: Use geth keyfile (more secure)")
			os.Exit(1)
		}

		// Load global config
		if err := loadGlobalConfig(); err != nil {
			fmt.Printf("Warning: failed to load config: %v\n", err)
		}

		// Resolve config values
		resolveConfigValues()

		// Call existing PreRun if it exists
		if existingPreRun != nil {
			existingPreRun(cmd, args)
		}
	}
}

// GetEthClient creates and returns an Ethereum client
func GetEthClient(rpcURL string) (*ethclient.Client, error) {
	return ethclient.Dial(rpcURL)
}

// ConnectToEthereum creates and returns an Ethereum client (alias for GetEthClient)
func ConnectToEthereum(rpcURL string) (*ethclient.Client, error) {
	return GetEthClient(rpcURL)
}

// GetValidatorManagerContract initializes and returns the ValidatorManager contract
func GetValidatorManagerContract(ethClient *ethclient.Client) (*bindings.IValidatorManager, error) {
	if validatorManagerContractAddr == "" {
		return nil, fmt.Errorf("ValidatorManager contract address is required")
	}

	validatorManagerAddr := common.HexToAddress(validatorManagerContractAddr)
	contract, err := bindings.NewIValidatorManager(validatorManagerAddr, ethClient)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize ValidatorManager contract: %w", err)
	}

	return contract, nil
}

// GetSigningConfig determines the signing method and returns the configuration
func GetSigningConfig() (*SigningConfig, error) {
	config := &SigningConfig{}

	// Check if private key is provided
	if privateKey != "" {
		privKey, err := parsePrivateKey(privateKey)
		if err != nil {
			return nil, fmt.Errorf("failed to parse private key: %w", err)
		}
		config.Method = SigningMethodPrivateKey
		config.PrivateKey = privKey
		return config, nil
	}

	// Check if keyfile is provided
	if keyfilePath != "" {
		password, err := getKeyfilePassword()
		if err != nil {
			return nil, fmt.Errorf("failed to get keyfile password: %w", err)
		}

		privKey, err := loadPrivateKeyFromKeyfile(keyfilePath, password)
		if err != nil {
			return nil, fmt.Errorf("failed to load private key from keyfile: %w", err)
		}

		config.Method = SigningMethodKeyfile
		config.PrivateKey = privKey
		config.KeyfilePath = keyfilePath
		config.KeyfilePassword = password
		return config, nil
	}

	return nil, fmt.Errorf("no signing method provided: use --private-key or --keyfile")
}

// parsePrivateKey converts a hex string to an ECDSA private key
func parsePrivateKey(key string) (*ecdsa.PrivateKey, error) {
	// Remove 0x prefix if present
	key = strings.TrimPrefix(key, "0x")

	privKey, err := ethcrypto.HexToECDSA(key)
	if err != nil {
		return nil, fmt.Errorf("invalid private key format: %w", err)
	}

	return privKey, nil
}

// getKeyfilePassword gets the keyfile password from various sources
func getKeyfilePassword() (string, error) {
	// Check if password is provided directly
	if keyfilePassword != "" {
		return keyfilePassword, nil
	}

	// Check if password file is provided
	if keyfilePasswordFile != "" {
		passwordBytes, err := ioutil.ReadFile(keyfilePasswordFile)
		if err != nil {
			return "", fmt.Errorf("failed to read password file: %w", err)
		}
		return strings.TrimSpace(string(passwordBytes)), nil
	}

	// Prompt for password
	fmt.Print("Enter keyfile password: ")
	passwordBytes, err := term.ReadPassword(int(syscall.Stdin))
	if err != nil {
		return "", fmt.Errorf("failed to read password: %w", err)
	}
	fmt.Println() // Add newline after password input

	return string(passwordBytes), nil
}

// loadPrivateKeyFromKeyfile loads a private key from a geth keyfile
func loadPrivateKeyFromKeyfile(keyfilePath, password string) (*ecdsa.PrivateKey, error) {
	// Read keyfile
	keyfileData, err := ioutil.ReadFile(keyfilePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read keyfile: %w", err)
	}

	// Create a temporary keyfile directory
	tempDir, err := ioutil.TempDir("", "temp_keyfile")
	if err != nil {
		return nil, fmt.Errorf("failed to create temp directory: %w", err)
	}
	defer os.RemoveAll(tempDir)

	// Copy keyfile to temp directory
	tempKeyfilePath := filepath.Join(tempDir, filepath.Base(keyfilePath))
	err = ioutil.WriteFile(tempKeyfilePath, keyfileData, 0600)
	if err != nil {
		return nil, fmt.Errorf("failed to write temp keyfile: %w", err)
	}

	// Create keyfile instance
	ks := keystore.NewKeyStore(tempDir, keystore.StandardScryptN, keystore.StandardScryptP)

	// Get accounts
	accounts := ks.Accounts()
	if len(accounts) == 0 {
		return nil, fmt.Errorf("no accounts found in keyfile")
	}

	// Use the first account
	account := accounts[0]

	// Unlock the account
	err = ks.Unlock(account, password)
	if err != nil {
		return nil, fmt.Errorf("failed to unlock keyfile: %w", err)
	}

	// Export the private key
	key, err := ks.Export(account, password, password)
	if err != nil {
		return nil, fmt.Errorf("failed to export private key: %w", err)
	}

	// Parse the exported key
	parsedKey, err := keystore.DecryptKey(key, password)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt key: %w", err)
	}

	return parsedKey.PrivateKey, nil
}

// ConfirmAction prompts the user to confirm an action
func ConfirmAction(message string) bool {
	if yes {
		return true
	}

	fmt.Printf("%s\nType 'yes' to continue: ", message)
	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)
	return strings.ToLower(input) == "yes"
}

// TransactOpts creates transaction options for a contract call
func TransactOpts(value *big.Int) (*bind.TransactOpts, error) {
	// Get signing configuration
	signingConfig, err := GetSigningConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to get signing configuration: %w", err)
	}

	// Get address from private key
	addr := ethcrypto.PubkeyToAddress(signingConfig.PrivateKey.PublicKey)

	// Determine nonce - use specified nonce or get from client
	var nVal uint64
	if nonceSpecified {
		nVal = nonce
	} else {
		nVal, err = ethClient.PendingNonceAt(context.Background(), addr)
		if err != nil {
			return nil, fmt.Errorf("failed to get nonce: %w", err)
		}
	}

	// Get chain ID
	chainID, err := ethClient.ChainID(context.Background())
	if err != nil {
		return nil, fmt.Errorf("failed to get chain ID: %w", err)
	}

	// Create transaction options
	opts, err := bind.NewKeyedTransactorWithChainID(signingConfig.PrivateKey, chainID)
	if err != nil {
		return nil, fmt.Errorf("failed to create transaction options: %w", err)
	}

	// Set nonce and value
	opts.Nonce = new(big.Int).SetUint64(nVal)
	opts.Value = value

	return opts, nil
}

// CreateTransactOpts creates transaction options for a contract call (alias for TransactOpts)
func CreateTransactOpts(ethClient *ethclient.Client, signingConfig *SigningConfig, value *big.Int) (*bind.TransactOpts, error) {
	// Get address from private key
	addr := ethcrypto.PubkeyToAddress(signingConfig.PrivateKey.PublicKey)

	// Get nonce
	nVal, err := ethClient.PendingNonceAt(context.Background(), addr)
	if err != nil {
		return nil, fmt.Errorf("failed to get nonce: %w", err)
	}

	// Get chain ID
	chainID, err := ethClient.ChainID(context.Background())
	if err != nil {
		return nil, fmt.Errorf("failed to get chain ID: %w", err)
	}

	// Create transaction options
	opts, err := bind.NewKeyedTransactorWithChainID(signingConfig.PrivateKey, chainID)
	if err != nil {
		return nil, fmt.Errorf("failed to create transaction options: %w", err)
	}

	// Set nonce and value
	opts.Nonce = new(big.Int).SetUint64(nVal)
	opts.Value = value

	return opts, nil
}

// WaitForTxConfirmation waits for a transaction to be mined and confirmed
func WaitForTxConfirmation(ethClient *ethclient.Client, txHash common.Hash) error {
	fmt.Printf("Waiting for transaction %s to be confirmed...\n", txHash.Hex())

	ctx := context.Background()

	// Set a timeout for 2 minutes
	timeoutCtx, cancel := context.WithTimeout(ctx, 2*time.Minute)
	defer cancel()

	// Poll for transaction receipt with a 2-second interval
	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-timeoutCtx.Done():
			return fmt.Errorf("timeout waiting for transaction confirmation")
		case <-ticker.C:
			receipt, err := ethClient.TransactionReceipt(ctx, txHash)
			if err != nil {
				// If error, likely the tx is not yet mined
				fmt.Print(".")
				continue
			}

			// Once we have a receipt, check its status
			if receipt.Status == 1 {
				blockNumber := receipt.BlockNumber
				fmt.Printf("\nTransaction confirmed in block %d\n", blockNumber.Uint64())
				return nil
			} else {
				return fmt.Errorf("transaction failed with status: %d", receipt.Status)
			}
		}
	}
}

// loadGlobalConfig loads the global configuration
func loadGlobalConfig() error {
	var err error
	globalConfig, err = loadConfig()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}
	return nil
}

// resolveConfigValues resolves configuration values with command line flags taking precedence
func resolveConfigValues() {
	if rpcURL == "" && globalConfig.RpcURL != "" {
		rpcURL = globalConfig.RpcURL
	}
	if validatorManagerContractAddr == "" && globalConfig.ValidatorManagerContractAddr != "" {
		validatorManagerContractAddr = globalConfig.ValidatorManagerContractAddr
	}
}

// AddMutuallyExclusiveValidation is a helper function to add validation to any command
func AddMutuallyExclusiveValidation(cmd *cobra.Command, groups ...MutuallyExclusiveGroup) {
	validator := NewFlagValidator()

	for _, group := range groups {
		validator.AddMutuallyExclusiveGroup(group)
	}

	// Store the original PreRun
	existingPreRun := cmd.PreRun

	cmd.PreRun = func(cmd *cobra.Command, args []string) {
		// Validate mutually exclusive flags
		if err := validator.ValidateFlags(cmd); err != nil {
			fmt.Printf("Error: %v\n", err)

			// Provide helpful usage information
			for _, group := range groups {
				if containsAnyFlag(err.Error(), group.Flags) {
					fmt.Printf("\nGroup '%s': %s\n", group.Name, group.Description)
					fmt.Printf("  Use only one of: %v\n", group.Flags)
				}
			}
			os.Exit(1)
		}

		// Call existing PreRun if it exists
		if existingPreRun != nil {
			existingPreRun(cmd, args)
		}
	}
}

// containsAnyFlag checks if the error message contains any of the specified flags
func containsAnyFlag(errorMsg string, flags []string) bool {
	for _, flag := range flags {
		if strings.Contains(errorMsg, flag) {
			return true
		}
	}
	return false
}

// ValidateNetworkInfo validates network information requirements for tx create commands
func ValidateNetworkInfo(cmd *cobra.Command) error {
	// Check if RPC URL is available (from config or flag)
	hasRPC := rpcURL != ""

	// Check if required network info flags are provided
	hasChainID := cmd.Flags().Changed("chain-id")
	hasGasPrice := cmd.Flags().Changed("gas-price")
	hasGasLimit := cmd.Flags().Changed("gas-limit")
	hasFee := cmd.Flags().Changed("fee")

	// If RPC is available, we can fetch network info automatically
	if hasRPC {
		return nil
	}

	// If no RPC, check if all required network info is provided
	missingFlags := []string{}

	if !hasChainID {
		missingFlags = append(missingFlags, "--chain-id")
	}
	if !hasGasPrice {
		missingFlags = append(missingFlags, "--gas-price")
	}
	if !hasGasLimit {
		missingFlags = append(missingFlags, "--gas-limit")
	}

	// For commands that require fee, check if fee is provided
	if cmd.Flags().Lookup("fee") != nil && !hasFee {
		missingFlags = append(missingFlags, "--fee")
	}

	if len(missingFlags) > 0 {
		return fmt.Errorf("no RPC URL configured and missing required network information: %v\n"+
			"Either configure RPC URL with 'mito config set-rpc <url>' or provide the missing flags", missingFlags)
	}

	return nil
}

// SetupNetworkInfo sets up network information either from RPC or from provided flags
func SetupNetworkInfo(cmd *cobra.Command) error {
	// If RPC is available, fetch network info
	if rpcURL != "" {
		ethClient, err := GetEthClient(rpcURL)
		if err != nil {
			return fmt.Errorf("failed to connect to RPC: %w", err)
		}

		// Get chain ID if not provided
		if !cmd.Flags().Changed("chain-id") {
			chainIDBig, err := ethClient.ChainID(context.Background())
			if err != nil {
				return fmt.Errorf("failed to get chain ID from RPC: %w", err)
			}
			chainID = chainIDBig.String()
		}

		// Get gas price if not provided
		if !cmd.Flags().Changed("gas-price") {
			gasPriceBig, err := ethClient.SuggestGasPrice(context.Background())
			if err != nil {
				return fmt.Errorf("failed to get gas price from RPC: %w", err)
			}
			gasPrice = gasPriceBig.String()
		}

		// For commands with contract interaction, get fee if not provided
		if cmd.Flags().Lookup("fee") != nil && !cmd.Flags().Changed("fee") {
			contract, err := GetValidatorManagerContract(ethClient)
			if err != nil {
				return fmt.Errorf("failed to initialize contract: %w", err)
			}

			fee, err := contract.Fee(nil)
			if err != nil {
				return fmt.Errorf("failed to get contract fee from RPC: %w", err)
			}
			contractFee = FormatWeiToEther(fee)
		}

		fmt.Printf("Network info fetched from RPC:\n")
		fmt.Printf("  Chain ID: %s\n", chainID)
		fmt.Printf("  Gas Price: %s wei\n", gasPrice)
		if cmd.Flags().Lookup("fee") != nil {
			fmt.Printf("  Contract Fee: %s MITO\n", contractFee)
		}
		fmt.Println()
	}

	return nil
}

// AddTxCreateNetworkFlags adds network-related flags for tx create commands
func AddTxCreateNetworkFlags(cmd *cobra.Command, requiresFee bool) {
	cmd.Flags().StringVar(&chainID, "chain-id", "", "Chain ID for the transaction (auto-fetched if RPC available)")
	cmd.Flags().Uint64Var(&gasLimit, "gas-limit", 500000, "Gas limit for the transaction")
	cmd.Flags().StringVar(&gasPrice, "gas-price", "", "Gas price in wei (auto-fetched if RPC available)")
	cmd.Flags().Uint64Var(&txNonce, "nonce", 0, "Nonce for the transaction (required for signed transactions)")

	if requiresFee {
		cmd.Flags().StringVar(&contractFee, "fee", "", "Contract fee in MITO (auto-fetched if RPC available)")
	}
}
