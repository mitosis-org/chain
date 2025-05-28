package cmd

import (
	"fmt"
	"math/big"
	"os"

	"github.com/ethereum/go-ethereum/common"
	"github.com/mitosis-org/chain/bindings"
	"github.com/spf13/cobra"
)

// newTxSendCollateralDepositCmd creates the collateral deposit command for tx send
func newTxSendCollateralDepositCmd() *cobra.Command {
	var (
		validator string
		amount    string
	)

	cmd := &cobra.Command{
		Use:   "deposit",
		Short: "Send a collateral deposit transaction",
		Long: `Create, sign, and send a transaction to deposit collateral for a validator.
Only permitted collateral owners can deposit collateral for a validator.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runSendDepositCollateralTx(cmd, validator, amount)
		},
	}

	// Add common flags for contract interaction
	cmd.Flags().StringVar(&rpcURL, "rpc", "", "RPC endpoint URL")
	cmd.Flags().StringVar(&validatorManagerContractAddr, "contract", "", "ValidatorManager contract address")

	// Add signing flags
	AddSigningFlags(cmd)

	// Command-specific flags
	cmd.Flags().StringVar(&validator, "validator", "", "Validator address")
	cmd.Flags().StringVar(&amount, "amount", "", "Amount of collateral to deposit in MITO (e.g., \"1.5\")")

	// Mark required flags
	mustMarkFlagRequired(cmd, "validator")
	mustMarkFlagRequired(cmd, "amount")

	// Add PreRun to load config and validate
	existingPreRun := cmd.PreRun
	cmd.PreRun = func(cmd *cobra.Command, args []string) {
		// Call existing PreRun first (from AddSigningFlags)
		if existingPreRun != nil {
			existingPreRun(cmd, args)
		}

		// Load global config
		if err := loadGlobalConfig(); err != nil {
			fmt.Printf("Warning: failed to load config: %v\n", err)
		}

		// Resolve config values
		resolveConfigValues()

		// Validate that required values are set
		if rpcURL == "" {
			fmt.Println("Error: RPC URL is required. Set it with --rpc flag or use 'mito config set-rpc <url>'")
			os.Exit(1)
		}
		if validatorManagerContractAddr == "" {
			fmt.Println("Error: ValidatorManager contract address is required. Set it with --contract flag or use 'mito config set-contract <address>'")
			os.Exit(1)
		}
	}

	return cmd
}

func runSendDepositCollateralTx(cmd *cobra.Command, validator, amount string) error {
	// Parse collateral amount as decimal MITO and convert to wei
	collateralAmount, err := ParseValueAsWei(amount)
	if err != nil {
		return fmt.Errorf("failed to parse amount: %w", err)
	}

	// Ensure collateral amount is greater than 0
	if collateralAmount.Cmp(big.NewInt(0)) <= 0 {
		return fmt.Errorf("collateral amount must be greater than 0")
	}

	// Validate validator address
	valAddr, err := ValidateAddress(validator)
	if err != nil {
		return fmt.Errorf("invalid validator address: %w", err)
	}

	// Get signing configuration
	signingConfig, err := GetSigningConfig()
	if err != nil {
		return fmt.Errorf("failed to get signing configuration: %w", err)
	}

	// Connect to Ethereum client
	ethClient, err := ConnectToEthereum(rpcURL)
	if err != nil {
		return fmt.Errorf("failed to connect to Ethereum: %w", err)
	}
	defer ethClient.Close()

	// Create contract instance
	contractAddr := common.HexToAddress(validatorManagerContractAddr)
	contract, err := bindings.NewIValidatorManager(contractAddr, ethClient)
	if err != nil {
		return fmt.Errorf("failed to create contract instance: %w", err)
	}

	// Get the contract fee
	fee, err := contract.Fee(nil)
	if err != nil {
		return fmt.Errorf("failed to get contract fee: %w", err)
	}

	// Calculate total value to send (collateral amount + fee)
	totalValue := new(big.Int).Add(collateralAmount, fee)

	// Show summary
	fmt.Println("===== Deposit Collateral Transaction =====")
	fmt.Printf("Validator Address        : %s\n", valAddr.Hex())
	fmt.Printf("Collateral Amount        : %s MITO\n", FormatWeiToEther(collateralAmount))
	fmt.Printf("Fee                      : %s MITO\n", FormatWeiToEther(fee))
	fmt.Printf("Total Value              : %s MITO\n", FormatWeiToEther(totalValue))
	fmt.Println()

	fmt.Println("🚨 IMPORTANT: Only permitted collateral owners can deposit collateral for a validator.")
	fmt.Println("If your address is not a permitted collateral owner for this validator, the transaction will fail.")
	fmt.Println()

	// Ask for confirmation
	if !ConfirmAction("Do you want to deposit this collateral?") {
		return fmt.Errorf("operation cancelled by user")
	}

	// Create transaction options
	transactOpts, err := CreateTransactOpts(ethClient, signingConfig, totalValue)
	if err != nil {
		return fmt.Errorf("failed to create transaction options: %w", err)
	}

	// Execute the transaction
	tx, err := contract.DepositCollateral(transactOpts, valAddr)
	if err != nil {
		return fmt.Errorf("failed to deposit collateral: %w", err)
	}

	fmt.Printf("Transaction sent! Hash: %s\n", tx.Hash().Hex())

	// Wait for transaction confirmation
	err = WaitForTxConfirmation(ethClient, tx.Hash())
	if err != nil {
		return fmt.Errorf("transaction failed: %w", err)
	}

	fmt.Println("✅ Collateral deposited successfully!")
	return nil
}

// newTxSendCollateralWithdrawCmd creates the collateral withdraw command for tx send
func newTxSendCollateralWithdrawCmd() *cobra.Command {
	var (
		validator string
		amount    string
		receiver  string
	)

	cmd := &cobra.Command{
		Use:   "withdraw",
		Short: "Send a collateral withdraw transaction",
		Long: `Create, sign, and send a transaction to withdraw collateral from a validator.
Only the collateral owner can withdraw collateral.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runSendWithdrawCollateralTx(cmd, validator, amount, receiver)
		},
	}

	// Add common flags for contract interaction
	cmd.Flags().StringVar(&rpcURL, "rpc", "", "RPC endpoint URL")
	cmd.Flags().StringVar(&validatorManagerContractAddr, "contract", "", "ValidatorManager contract address")

	// Add signing flags
	AddSigningFlags(cmd)

	// Command-specific flags
	cmd.Flags().StringVar(&validator, "validator", "", "Validator address")
	cmd.Flags().StringVar(&receiver, "receiver", "", "Address to receive the withdrawn collateral")
	cmd.Flags().StringVar(&amount, "amount", "", "Amount of collateral to withdraw in MITO (e.g., \"1.5\")")

	// Mark required flags
	mustMarkFlagRequired(cmd, "validator")
	mustMarkFlagRequired(cmd, "receiver")
	mustMarkFlagRequired(cmd, "amount")

	// Add PreRun to load config and validate
	existingPreRun := cmd.PreRun
	cmd.PreRun = func(cmd *cobra.Command, args []string) {
		// Call existing PreRun first (from AddSigningFlags)
		if existingPreRun != nil {
			existingPreRun(cmd, args)
		}

		// Load global config
		if err := loadGlobalConfig(); err != nil {
			fmt.Printf("Warning: failed to load config: %v\n", err)
		}

		// Resolve config values
		resolveConfigValues()

		// Validate that required values are set
		if rpcURL == "" {
			fmt.Println("Error: RPC URL is required. Set it with --rpc flag or use 'mito config set-rpc <url>'")
			os.Exit(1)
		}
		if validatorManagerContractAddr == "" {
			fmt.Println("Error: ValidatorManager contract address is required. Set it with --contract flag or use 'mito config set-contract <address>'")
			os.Exit(1)
		}
	}

	return cmd
}

func runSendWithdrawCollateralTx(cmd *cobra.Command, validator, amount, receiver string) error {
	// Parse collateral amount as decimal MITO and convert to wei
	collateralAmount, err := ParseValueAsWei(amount)
	if err != nil {
		return fmt.Errorf("failed to parse amount: %w", err)
	}

	// Ensure collateral amount is greater than 0
	if collateralAmount.Cmp(big.NewInt(0)) <= 0 {
		return fmt.Errorf("collateral amount must be greater than 0")
	}

	// Validate validator address
	valAddr, err := ValidateAddress(validator)
	if err != nil {
		return fmt.Errorf("invalid validator address: %w", err)
	}

	// Validate receiver address
	receiverAddr, err := ValidateAddress(receiver)
	if err != nil {
		return fmt.Errorf("invalid receiver address: %w", err)
	}

	// Get signing configuration
	signingConfig, err := GetSigningConfig()
	if err != nil {
		return fmt.Errorf("failed to get signing configuration: %w", err)
	}

	// Connect to Ethereum client
	ethClient, err := ConnectToEthereum(rpcURL)
	if err != nil {
		return fmt.Errorf("failed to connect to Ethereum: %w", err)
	}
	defer ethClient.Close()

	// Create contract instance
	contractAddr := common.HexToAddress(validatorManagerContractAddr)
	contract, err := bindings.NewIValidatorManager(contractAddr, ethClient)
	if err != nil {
		return fmt.Errorf("failed to create contract instance: %w", err)
	}

	// Get the contract fee
	fee, err := contract.Fee(nil)
	if err != nil {
		return fmt.Errorf("failed to get contract fee: %w", err)
	}

	// Show summary
	fmt.Println("===== Withdraw Collateral Transaction =====")
	fmt.Printf("Validator Address        : %s\n", valAddr.Hex())
	fmt.Printf("Collateral Amount        : %s MITO\n", FormatWeiToEther(collateralAmount))
	fmt.Printf("Fee                      : %s MITO\n", FormatWeiToEther(fee))
	fmt.Printf("Receiver Address         : %s\n", receiverAddr.Hex())
	fmt.Println()

	fmt.Println("🚨 IMPORTANT: Only the collateral owner can withdraw collateral.")
	fmt.Println("Make sure you are the owner of the collateral for this validator.")
	fmt.Println()

	// Ask for confirmation
	if !ConfirmAction("Do you want to withdraw this collateral?") {
		return fmt.Errorf("operation cancelled by user")
	}

	// Create transaction options (withdraw only sends fee, not collateral)
	transactOpts, err := CreateTransactOpts(ethClient, signingConfig, fee)
	if err != nil {
		return fmt.Errorf("failed to create transaction options: %w", err)
	}

	// Execute the transaction
	tx, err := contract.WithdrawCollateral(transactOpts, valAddr, receiverAddr, collateralAmount)
	if err != nil {
		return fmt.Errorf("failed to withdraw collateral: %w", err)
	}

	fmt.Printf("Transaction sent! Hash: %s\n", tx.Hash().Hex())

	// Wait for transaction confirmation
	err = WaitForTxConfirmation(ethClient, tx.Hash())
	if err != nil {
		return fmt.Errorf("transaction failed: %w", err)
	}

	fmt.Println("✅ Collateral withdrawn successfully!")
	return nil
}

// newTxSendCollateralSetPermittedOwnerCmd creates the set permitted collateral owner command for tx send
func newTxSendCollateralSetPermittedOwnerCmd() *cobra.Command {
	var (
		validator       string
		collateralOwner string
		isPermitted     bool
	)

	cmd := &cobra.Command{
		Use:   "set-permitted-owner",
		Short: "Send a set permitted collateral owner transaction",
		Long: `Create, sign, and send a transaction to set a permitted collateral owner for a validator.
Only the validator operator can set permitted collateral owners.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runSendSetPermittedCollateralOwnerTx(cmd, validator, collateralOwner, isPermitted)
		},
	}

	// Add common flags for contract interaction
	cmd.Flags().StringVar(&rpcURL, "rpc", "", "RPC endpoint URL")
	cmd.Flags().StringVar(&validatorManagerContractAddr, "contract", "", "ValidatorManager contract address")

	// Add signing flags
	AddSigningFlags(cmd)

	// Command-specific flags
	cmd.Flags().StringVar(&validator, "validator", "", "Validator address")
	cmd.Flags().StringVar(&collateralOwner, "collateral-owner", "", "Collateral owner address to permit or deny")
	cmd.Flags().BoolVar(&isPermitted, "permitted", false, "Whether to permit (true) or deny (false) the address")

	// Mark required flags
	mustMarkFlagRequired(cmd, "validator")
	mustMarkFlagRequired(cmd, "collateral-owner")

	// Add PreRun to load config and validate
	existingPreRun := cmd.PreRun
	cmd.PreRun = func(cmd *cobra.Command, args []string) {
		// Call existing PreRun first (from AddSigningFlags)
		if existingPreRun != nil {
			existingPreRun(cmd, args)
		}

		// Load global config
		if err := loadGlobalConfig(); err != nil {
			fmt.Printf("Warning: failed to load config: %v\n", err)
		}

		// Resolve config values
		resolveConfigValues()

		// Validate that required values are set
		if rpcURL == "" {
			fmt.Println("Error: RPC URL is required. Set it with --rpc flag or use 'mito config set-rpc <url>'")
			os.Exit(1)
		}
		if validatorManagerContractAddr == "" {
			fmt.Println("Error: ValidatorManager contract address is required. Set it with --contract flag or use 'mito config set-contract <address>'")
			os.Exit(1)
		}
	}

	return cmd
}

func runSendSetPermittedCollateralOwnerTx(cmd *cobra.Command, validator, collateralOwner string, isPermitted bool) error {
	// Validate validator address
	valAddr, err := ValidateAddress(validator)
	if err != nil {
		return fmt.Errorf("invalid validator address: %w", err)
	}

	// Validate collateral owner address
	collateralOwnerAddr, err := ValidateAddress(collateralOwner)
	if err != nil {
		return fmt.Errorf("invalid collateral owner address: %w", err)
	}

	// Get signing configuration
	signingConfig, err := GetSigningConfig()
	if err != nil {
		return fmt.Errorf("failed to get signing configuration: %w", err)
	}

	// Connect to Ethereum client
	ethClient, err := ConnectToEthereum(rpcURL)
	if err != nil {
		return fmt.Errorf("failed to connect to Ethereum: %w", err)
	}
	defer ethClient.Close()

	// Create contract instance
	contractAddr := common.HexToAddress(validatorManagerContractAddr)
	contract, err := bindings.NewIValidatorManager(contractAddr, ethClient)
	if err != nil {
		return fmt.Errorf("failed to create contract instance: %w", err)
	}

	// Check current permission status
	currentPermission, err := contract.IsPermittedCollateralOwner(nil, valAddr, collateralOwnerAddr)
	if err != nil {
		return fmt.Errorf("failed to check current permission status: %w", err)
	}

	// Show summary
	fmt.Println("===== Set Permitted Collateral Owner Transaction =====")
	fmt.Printf("Validator Address          : %s\n", valAddr.Hex())
	fmt.Printf("Collateral Owner Address   : %s\n", collateralOwnerAddr.Hex())
	fmt.Printf("Current Permission Status  : %t\n", currentPermission)
	fmt.Printf("New Permission Status      : %t\n", isPermitted)
	fmt.Println()

	if isPermitted {
		fmt.Println("🚨 IMPORTANT: This will allow the specified address to deposit collateral for your validator.")
		fmt.Println("Make sure you trust this address or it is under your control.")
	} else {
		fmt.Println("🚨 IMPORTANT: This will revoke permission for the specified address to deposit collateral for your validator.")
	}
	fmt.Println()

	// Ask for confirmation
	if !ConfirmAction("Do you want to update the permission status?") {
		return fmt.Errorf("operation cancelled by user")
	}

	// Create transaction options (no value needed for permission update)
	transactOpts, err := CreateTransactOpts(ethClient, signingConfig, nil)
	if err != nil {
		return fmt.Errorf("failed to create transaction options: %w", err)
	}

	// Execute the transaction
	tx, err := contract.SetPermittedCollateralOwner(transactOpts, valAddr, collateralOwnerAddr, isPermitted)
	if err != nil {
		return fmt.Errorf("failed to set permitted collateral owner: %w", err)
	}

	fmt.Printf("Transaction sent! Hash: %s\n", tx.Hash().Hex())

	// Wait for transaction confirmation
	err = WaitForTxConfirmation(ethClient, tx.Hash())
	if err != nil {
		return fmt.Errorf("transaction failed: %w", err)
	}

	fmt.Println("✅ Permitted collateral owner status updated successfully!")
	return nil
}

// newTxSendCollateralTransferOwnershipCmd creates the transfer collateral ownership command for tx send
func newTxSendCollateralTransferOwnershipCmd() *cobra.Command {
	var (
		validator string
		newOwner  string
	)

	cmd := &cobra.Command{
		Use:   "transfer-ownership",
		Short: "Send a transfer collateral ownership transaction",
		Long: `Create, sign, and send a transaction to transfer collateral ownership for a validator.
Only permitted collateral owners can transfer collateral ownership.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runSendTransferCollateralOwnershipTx(cmd, validator, newOwner)
		},
	}

	// Add common flags for contract interaction
	cmd.Flags().StringVar(&rpcURL, "rpc", "", "RPC endpoint URL")
	cmd.Flags().StringVar(&validatorManagerContractAddr, "contract", "", "ValidatorManager contract address")

	// Add signing flags
	AddSigningFlags(cmd)

	// Command-specific flags
	cmd.Flags().StringVar(&validator, "validator", "", "Validator address")
	cmd.Flags().StringVar(&newOwner, "new-owner", "", "New collateral owner address")

	// Mark required flags
	mustMarkFlagRequired(cmd, "validator")
	mustMarkFlagRequired(cmd, "new-owner")

	// Add PreRun to load config and validate
	existingPreRun := cmd.PreRun
	cmd.PreRun = func(cmd *cobra.Command, args []string) {
		// Call existing PreRun first (from AddSigningFlags)
		if existingPreRun != nil {
			existingPreRun(cmd, args)
		}

		// Load global config
		if err := loadGlobalConfig(); err != nil {
			fmt.Printf("Warning: failed to load config: %v\n", err)
		}

		// Resolve config values
		resolveConfigValues()

		// Validate that required values are set
		if rpcURL == "" {
			fmt.Println("Error: RPC URL is required. Set it with --rpc flag or use 'mito config set-rpc <url>'")
			os.Exit(1)
		}
		if validatorManagerContractAddr == "" {
			fmt.Println("Error: ValidatorManager contract address is required. Set it with --contract flag or use 'mito config set-contract <address>'")
			os.Exit(1)
		}
	}

	return cmd
}

func runSendTransferCollateralOwnershipTx(cmd *cobra.Command, validator, newOwner string) error {
	// Validate validator address
	valAddr, err := ValidateAddress(validator)
	if err != nil {
		return fmt.Errorf("invalid validator address: %w", err)
	}

	// Validate new owner address
	newOwnerAddr, err := ValidateAddress(newOwner)
	if err != nil {
		return fmt.Errorf("invalid new owner address: %w", err)
	}

	// Get signing configuration
	signingConfig, err := GetSigningConfig()
	if err != nil {
		return fmt.Errorf("failed to get signing configuration: %w", err)
	}

	// Connect to Ethereum client
	ethClient, err := ConnectToEthereum(rpcURL)
	if err != nil {
		return fmt.Errorf("failed to connect to Ethereum: %w", err)
	}
	defer ethClient.Close()

	// Create contract instance
	contractAddr := common.HexToAddress(validatorManagerContractAddr)
	contract, err := bindings.NewIValidatorManager(contractAddr, ethClient)
	if err != nil {
		return fmt.Errorf("failed to create contract instance: %w", err)
	}

	// Get the contract fee
	fee, err := contract.Fee(nil)
	if err != nil {
		return fmt.Errorf("failed to get contract fee: %w", err)
	}

	// Show summary
	fmt.Println("===== Transfer Collateral Ownership Transaction =====")
	fmt.Printf("Validator Address        : %s\n", valAddr.Hex())
	fmt.Printf("New Collateral Owner     : %s\n", newOwnerAddr.Hex())
	fmt.Printf("Fee                      : %s MITO\n", FormatWeiToEther(fee))
	fmt.Println()

	fmt.Println("🚨 IMPORTANT: Only permitted collateral owners can transfer collateral ownership.")
	fmt.Println("If your address is not a permitted collateral owner for this validator, the transaction will fail.")
	fmt.Println("To check all permitted collateral owners, use 'mito validator info --validator-address <validator-address>'")
	fmt.Println()

	fmt.Println("🚨 IMPORTANT: This action will transfer ownership of the validator's collateral.")
	fmt.Println("The new owner will have full control over the collateral, including the ability to withdraw it.")
	fmt.Println("Make sure you trust the new owner or it is an address you control.")
	fmt.Println()

	// Ask for confirmation
	if !ConfirmAction("Do you want to transfer collateral ownership?") {
		return fmt.Errorf("operation cancelled by user")
	}

	// Create transaction options (fee is required for ownership transfer)
	transactOpts, err := CreateTransactOpts(ethClient, signingConfig, fee)
	if err != nil {
		return fmt.Errorf("failed to create transaction options: %w", err)
	}

	// Execute the transaction
	tx, err := contract.TransferCollateralOwnership(transactOpts, valAddr, newOwnerAddr)
	if err != nil {
		return fmt.Errorf("failed to transfer collateral ownership: %w", err)
	}

	fmt.Printf("Transaction sent! Hash: %s\n", tx.Hash().Hex())

	// Wait for transaction confirmation
	err = WaitForTxConfirmation(ethClient, tx.Hash())
	if err != nil {
		return fmt.Errorf("transaction failed: %w", err)
	}

	fmt.Println("✅ Collateral ownership transferred successfully!")
	return nil
}
