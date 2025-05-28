package cmd

import (
	"context"
	"fmt"
	"log"
	"math/big"
	"os"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	ethcrypto "github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/spf13/cobra"
)

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

	// Add PreRun to load config and validate signing requirements
	existingPreRun := cmd.PreRun
	cmd.PreRun = func(cmd *cobra.Command, args []string) {
		// Check if signing is required
		if signed {
			privateKeySet := cmd.Flags().Changed("private-key")
			keyfileSet := cmd.Flags().Changed("keyfile")

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

	// Add mutually exclusive validation using flags.go (this must be called last to preserve PreRun chain)
	AddMutuallyExclusiveValidation(cmd, TransactionTypeGroup, SigningMethodGroup)
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
