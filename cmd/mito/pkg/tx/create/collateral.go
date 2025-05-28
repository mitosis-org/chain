package create

import (
	"encoding/json"
	"fmt"
	"math/big"
	"os"

	"github.com/mitosis-org/chain/bindings"
	"github.com/mitosis-org/chain/cmd/mito/internal/config"
	"github.com/mitosis-org/chain/cmd/mito/internal/container"
	"github.com/mitosis-org/chain/cmd/mito/internal/flags"
	"github.com/mitosis-org/chain/cmd/mito/internal/utils"
	"github.com/spf13/cobra"
)

// NewCollateralCmd returns collateral commands for tx create
func NewCollateralCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "collateral",
		Short: "Create collateral transactions",
		Long:  "Create collateral transactions (signed or unsigned)",
	}

	cmd.AddCommand(
		newCreateCollateralDepositCmd(),
		newCreateCollateralWithdrawCmd(),
		newCreateCollateralSetPermittedOwnerCmd(),
		newCreateCollateralTransferOwnershipCmd(),
	)

	return cmd
}

func newCreateCollateralDepositCmd() *cobra.Command {
	var commonFlags flags.CommonFlags
	var collateralFlags struct {
		validator string
		amount    string
	}

	cmd := &cobra.Command{
		Use:   "deposit",
		Short: "Create collateral deposit transaction",
		Long:  "Create a transaction to deposit collateral (without sending)",
		RunE: func(cmd *cobra.Command, args []string) error {
			cmd.SilenceUsage = true

			// Resolve configuration
			resolver, err := config.NewResolver()
			if err != nil {
				return err
			}
			resolvedConfig := resolver.ResolveFlags(&commonFlags)

			// Validate required fields
			if err := validateCollateralDepositFields(resolvedConfig, &collateralFlags); err != nil {
				return err
			}

			// Create container
			container, err := container.NewContainer(resolvedConfig)
			if err != nil {
				return fmt.Errorf("failed to initialize container: %w", err)
			}
			defer container.Close()

			// Create transaction data
			txData, err := createCollateralDepositTransaction(container, collateralFlags.validator, collateralFlags.amount, resolvedConfig)
			if err != nil {
				return err
			}

			// Output transaction data
			return outputCollateralTransactionData(txData, resolvedConfig)
		},
	}

	// Add flags
	flags.AddCommonFlags(cmd, &commonFlags)
	cmd.Flags().StringVar(&collateralFlags.validator, "validator", "", "Validator address")
	cmd.Flags().StringVar(&collateralFlags.amount, "amount", "", "Amount to deposit in MITO (e.g., \"1.5\")")

	// Mark required flags
	cmd.MarkFlagRequired("validator")
	cmd.MarkFlagRequired("amount")

	return cmd
}

func newCreateCollateralWithdrawCmd() *cobra.Command {
	var commonFlags flags.CommonFlags
	var collateralFlags struct {
		validator string
		amount    string
		receiver  string
	}

	cmd := &cobra.Command{
		Use:   "withdraw",
		Short: "Create collateral withdraw transaction",
		Long:  "Create a transaction to withdraw collateral (without sending)",
		RunE: func(cmd *cobra.Command, args []string) error {
			cmd.SilenceUsage = true

			// Resolve configuration
			resolver, err := config.NewResolver()
			if err != nil {
				return err
			}
			resolvedConfig := resolver.ResolveFlags(&commonFlags)

			// Validate required fields
			if err := validateCollateralWithdrawFields(resolvedConfig, &collateralFlags); err != nil {
				return err
			}

			// Create container
			container, err := container.NewContainer(resolvedConfig)
			if err != nil {
				return fmt.Errorf("failed to initialize container: %w", err)
			}
			defer container.Close()

			// Create transaction data
			txData, err := createCollateralWithdrawTransaction(container, collateralFlags.validator, collateralFlags.amount, collateralFlags.receiver, resolvedConfig)
			if err != nil {
				return err
			}

			// Output transaction data
			return outputCollateralTransactionData(txData, resolvedConfig)
		},
	}

	// Add flags
	flags.AddCommonFlags(cmd, &commonFlags)
	cmd.Flags().StringVar(&collateralFlags.validator, "validator", "", "Validator address")
	cmd.Flags().StringVar(&collateralFlags.amount, "amount", "", "Amount to withdraw in MITO (e.g., \"1.5\")")
	cmd.Flags().StringVar(&collateralFlags.receiver, "receiver", "", "Address to receive the withdrawn collateral")

	// Mark required flags
	cmd.MarkFlagRequired("validator")
	cmd.MarkFlagRequired("amount")
	cmd.MarkFlagRequired("receiver")

	return cmd
}

func newCreateCollateralSetPermittedOwnerCmd() *cobra.Command {
	var commonFlags flags.CommonFlags
	var collateralFlags struct {
		validator       string
		collateralOwner string
		isPermitted     bool
	}

	cmd := &cobra.Command{
		Use:   "set-permitted-owner",
		Short: "Create set permitted owner transaction",
		Long:  "Create a transaction to set permitted collateral owner (without sending)",
		RunE: func(cmd *cobra.Command, args []string) error {
			cmd.SilenceUsage = true

			// Resolve configuration
			resolver, err := config.NewResolver()
			if err != nil {
				return err
			}
			resolvedConfig := resolver.ResolveFlags(&commonFlags)

			// Validate required fields
			if err := validateCollateralSetPermittedOwnerFields(resolvedConfig, &collateralFlags); err != nil {
				return err
			}

			// Create container
			container, err := container.NewContainer(resolvedConfig)
			if err != nil {
				return fmt.Errorf("failed to initialize container: %w", err)
			}
			defer container.Close()

			// Create transaction data
			txData, err := createCollateralSetPermittedOwnerTransaction(container, collateralFlags.validator, collateralFlags.collateralOwner, collateralFlags.isPermitted, resolvedConfig)
			if err != nil {
				return err
			}

			// Output transaction data
			return outputCollateralTransactionData(txData, resolvedConfig)
		},
	}

	// Add flags
	flags.AddCommonFlags(cmd, &commonFlags)
	cmd.Flags().StringVar(&collateralFlags.validator, "validator", "", "Validator address")
	cmd.Flags().StringVar(&collateralFlags.collateralOwner, "collateral-owner", "", "Collateral owner address to permit or deny")
	cmd.Flags().BoolVar(&collateralFlags.isPermitted, "permitted", false, "Whether to permit (true) or deny (false) the address")

	// Mark required flags
	cmd.MarkFlagRequired("validator")
	cmd.MarkFlagRequired("collateral-owner")

	return cmd
}

func newCreateCollateralTransferOwnershipCmd() *cobra.Command {
	var commonFlags flags.CommonFlags
	var collateralFlags struct {
		validator string
		newOwner  string
	}

	cmd := &cobra.Command{
		Use:   "transfer-ownership",
		Short: "Create transfer ownership transaction",
		Long:  "Create a transaction to transfer collateral ownership (without sending)",
		RunE: func(cmd *cobra.Command, args []string) error {
			cmd.SilenceUsage = true

			// Resolve configuration
			resolver, err := config.NewResolver()
			if err != nil {
				return err
			}
			resolvedConfig := resolver.ResolveFlags(&commonFlags)

			// Validate required fields
			if err := validateCollateralTransferOwnershipFields(resolvedConfig, &collateralFlags); err != nil {
				return err
			}

			// Create container
			container, err := container.NewContainer(resolvedConfig)
			if err != nil {
				return fmt.Errorf("failed to initialize container: %w", err)
			}
			defer container.Close()

			// Create transaction data
			txData, err := createCollateralTransferOwnershipTransaction(container, collateralFlags.validator, collateralFlags.newOwner, resolvedConfig)
			if err != nil {
				return err
			}

			// Output transaction data
			return outputCollateralTransactionData(txData, resolvedConfig)
		},
	}

	// Add flags
	flags.AddCommonFlags(cmd, &commonFlags)
	cmd.Flags().StringVar(&collateralFlags.validator, "validator", "", "Validator address")
	cmd.Flags().StringVar(&collateralFlags.newOwner, "new-owner", "", "New owner address")

	// Mark required flags
	cmd.MarkFlagRequired("validator")
	cmd.MarkFlagRequired("new-owner")

	return cmd
}

// Transaction data structure for output
type CollateralTransactionData struct {
	To       string `json:"to"`
	Value    string `json:"value"`
	Data     string `json:"data"`
	GasLimit uint64 `json:"gasLimit"`
	GasPrice string `json:"gasPrice"`
	ChainID  string `json:"chainId"`
}

// Helper functions for creating collateral transaction data
func createCollateralDepositTransaction(container *container.Container, validator, amount string, config *config.ResolvedConfig) (*CollateralTransactionData, error) {
	// Parse collateral amount as decimal MITO and convert to wei
	collateralAmount, err := utils.ParseValueAsWei(amount)
	if err != nil {
		return nil, fmt.Errorf("failed to parse amount: %w", err)
	}

	// Ensure collateral amount is greater than 0
	if collateralAmount.Cmp(big.NewInt(0)) <= 0 {
		return nil, fmt.Errorf("collateral amount must be greater than 0")
	}

	// Get contract fee
	var fee *big.Int

	if config.ContractFee != "" {
		// Use provided contract fee
		fee, err = utils.ParseValueAsWei(config.ContractFee)
		if err != nil {
			return nil, fmt.Errorf("failed to parse contract fee: %w", err)
		}
	} else {
		// Get contract fee from RPC
		fee, err = container.Contract.Fee(nil)
		if err != nil {
			return nil, fmt.Errorf("failed to get contract fee: %w", err)
		}
	}

	// Calculate total value to send (collateral amount + fee)
	totalValue := new(big.Int).Add(collateralAmount, fee)

	// Validate validator address
	valAddr, err := utils.ValidateAddress(validator)
	if err != nil {
		return nil, fmt.Errorf("invalid validator address: %w", err)
	}

	// Get contract ABI and encode function call data
	abi, err := bindings.IValidatorManagerMetaData.GetAbi()
	if err != nil {
		return nil, fmt.Errorf("failed to get contract ABI: %w", err)
	}

	data, err := abi.Pack("depositCollateral", valAddr)
	if err != nil {
		return nil, fmt.Errorf("failed to pack function call: %w", err)
	}

	// Set default gas limit if not provided
	gasLimit := config.GasLimit
	if gasLimit == 0 {
		gasLimit = 500000
	}

	// Show summary
	fmt.Println("===== Deposit Collateral Transaction =====")
	fmt.Printf("Validator Address          : %s\n", valAddr.Hex())
	fmt.Printf("Collateral Amount          : %s MITO\n", utils.FormatWeiToEther(collateralAmount))
	fmt.Printf("Fee                        : %s MITO\n", utils.FormatWeiToEther(fee))
	fmt.Printf("Total Value                : %s MITO\n", utils.FormatWeiToEther(totalValue))
	fmt.Printf("Chain ID                   : %s\n", config.ChainID)
	fmt.Printf("Gas Limit                  : %d\n", gasLimit)
	fmt.Printf("Gas Price                  : %s wei\n", config.GasPrice)
	fmt.Println()
	fmt.Println("🚨 IMPORTANT: Only permitted collateral owners can deposit collateral for a validator.")
	fmt.Println("If your address is not a permitted collateral owner for this validator, the transaction will fail.")
	fmt.Println()

	return &CollateralTransactionData{
		To:       config.ValidatorManagerContractAddr,
		Value:    totalValue.String(),
		Data:     fmt.Sprintf("0x%x", data),
		GasLimit: gasLimit,
		GasPrice: config.GasPrice,
		ChainID:  config.ChainID,
	}, nil
}

func createCollateralWithdrawTransaction(container *container.Container, validator, amount, receiver string, config *config.ResolvedConfig) (*CollateralTransactionData, error) {
	// Parse collateral amount as decimal MITO and convert to wei
	collateralAmount, err := utils.ParseValueAsWei(amount)
	if err != nil {
		return nil, fmt.Errorf("failed to parse amount: %w", err)
	}

	// Ensure collateral amount is greater than 0
	if collateralAmount.Cmp(big.NewInt(0)) <= 0 {
		return nil, fmt.Errorf("collateral amount must be greater than 0")
	}

	// Get contract fee
	var fee *big.Int

	if config.ContractFee != "" {
		// Use provided contract fee
		fee, err = utils.ParseValueAsWei(config.ContractFee)
		if err != nil {
			return nil, fmt.Errorf("failed to parse contract fee: %w", err)
		}
	} else {
		// Get contract fee from RPC
		fee, err = container.Contract.Fee(nil)
		if err != nil {
			return nil, fmt.Errorf("failed to get contract fee: %w", err)
		}
	}

	// Validate addresses
	valAddr, err := utils.ValidateAddress(validator)
	if err != nil {
		return nil, fmt.Errorf("invalid validator address: %w", err)
	}

	receiverAddr, err := utils.ValidateAddress(receiver)
	if err != nil {
		return nil, fmt.Errorf("invalid receiver address: %w", err)
	}

	// Get contract ABI and encode function call data
	abi, err := bindings.IValidatorManagerMetaData.GetAbi()
	if err != nil {
		return nil, fmt.Errorf("failed to get contract ABI: %w", err)
	}

	data, err := abi.Pack("withdrawCollateral", valAddr, receiverAddr, collateralAmount)
	if err != nil {
		return nil, fmt.Errorf("failed to pack function call: %w", err)
	}

	// Set default gas limit if not provided
	gasLimit := config.GasLimit
	if gasLimit == 0 {
		gasLimit = 500000
	}

	// Show summary
	fmt.Println("===== Withdraw Collateral Transaction =====")
	fmt.Printf("Validator Address          : %s\n", valAddr.Hex())
	fmt.Printf("Receiver Address           : %s\n", receiverAddr.Hex())
	fmt.Printf("Collateral Amount          : %s MITO\n", utils.FormatWeiToEther(collateralAmount))
	fmt.Printf("Fee                        : %s MITO\n", utils.FormatWeiToEther(fee))
	fmt.Printf("Chain ID                   : %s\n", config.ChainID)
	fmt.Printf("Gas Limit                  : %d\n", gasLimit)
	fmt.Printf("Gas Price                  : %s wei\n", config.GasPrice)
	fmt.Println()
	fmt.Println("🚨 IMPORTANT: Only the collateral owner can withdraw collateral.")
	fmt.Println()

	return &CollateralTransactionData{
		To:       config.ValidatorManagerContractAddr,
		Value:    fee.String(), // withdraw only sends fee, not collateral
		Data:     fmt.Sprintf("0x%x", data),
		GasLimit: gasLimit,
		GasPrice: config.GasPrice,
		ChainID:  config.ChainID,
	}, nil
}

func createCollateralSetPermittedOwnerTransaction(container *container.Container, validator, collateralOwner string, isPermitted bool, config *config.ResolvedConfig) (*CollateralTransactionData, error) {
	// Validate addresses
	valAddr, err := utils.ValidateAddress(validator)
	if err != nil {
		return nil, fmt.Errorf("invalid validator address: %w", err)
	}

	collateralOwnerAddr, err := utils.ValidateAddress(collateralOwner)
	if err != nil {
		return nil, fmt.Errorf("invalid collateral owner address: %w", err)
	}

	// Get contract ABI and encode function call data
	abi, err := bindings.IValidatorManagerMetaData.GetAbi()
	if err != nil {
		return nil, fmt.Errorf("failed to get contract ABI: %w", err)
	}

	data, err := abi.Pack("setPermittedCollateralOwner", valAddr, collateralOwnerAddr, isPermitted)
	if err != nil {
		return nil, fmt.Errorf("failed to pack function call: %w", err)
	}

	// Set default gas limit if not provided
	gasLimit := config.GasLimit
	if gasLimit == 0 {
		gasLimit = 500000
	}

	// Show summary
	permissionText := "DENY"
	if isPermitted {
		permissionText = "PERMIT"
	}

	fmt.Println("===== Set Permitted Collateral Owner Transaction =====")
	fmt.Printf("Validator Address          : %s\n", valAddr.Hex())
	fmt.Printf("Collateral Owner           : %s\n", collateralOwnerAddr.Hex())
	fmt.Printf("Permission                 : %s\n", permissionText)
	fmt.Printf("Chain ID                   : %s\n", config.ChainID)
	fmt.Printf("Gas Limit                  : %d\n", gasLimit)
	fmt.Printf("Gas Price                  : %s wei\n", config.GasPrice)
	fmt.Println()

	return &CollateralTransactionData{
		To:       config.ValidatorManagerContractAddr,
		Value:    "0", // no value needed for permission update
		Data:     fmt.Sprintf("0x%x", data),
		GasLimit: gasLimit,
		GasPrice: config.GasPrice,
		ChainID:  config.ChainID,
	}, nil
}

func createCollateralTransferOwnershipTransaction(container *container.Container, validator, newOwner string, config *config.ResolvedConfig) (*CollateralTransactionData, error) {
	// Get contract fee
	var fee *big.Int
	var err error

	if config.ContractFee != "" {
		// Use provided contract fee
		fee, err = utils.ParseValueAsWei(config.ContractFee)
		if err != nil {
			return nil, fmt.Errorf("failed to parse contract fee: %w", err)
		}
	} else {
		// Get contract fee from RPC
		fee, err = container.Contract.Fee(nil)
		if err != nil {
			return nil, fmt.Errorf("failed to get contract fee: %w", err)
		}
	}

	// Validate addresses
	valAddr, err := utils.ValidateAddress(validator)
	if err != nil {
		return nil, fmt.Errorf("invalid validator address: %w", err)
	}

	newOwnerAddr, err := utils.ValidateAddress(newOwner)
	if err != nil {
		return nil, fmt.Errorf("invalid new owner address: %w", err)
	}

	// Get contract ABI and encode function call data
	abi, err := bindings.IValidatorManagerMetaData.GetAbi()
	if err != nil {
		return nil, fmt.Errorf("failed to get contract ABI: %w", err)
	}

	data, err := abi.Pack("transferCollateralOwnership", valAddr, newOwnerAddr)
	if err != nil {
		return nil, fmt.Errorf("failed to pack function call: %w", err)
	}

	// Set default gas limit if not provided
	gasLimit := config.GasLimit
	if gasLimit == 0 {
		gasLimit = 500000
	}

	// Show summary
	fmt.Println("===== Transfer Collateral Ownership Transaction =====")
	fmt.Printf("Validator Address          : %s\n", valAddr.Hex())
	fmt.Printf("New Owner                  : %s\n", newOwnerAddr.Hex())
	fmt.Printf("Fee                        : %s MITO\n", utils.FormatWeiToEther(fee))
	fmt.Printf("Chain ID                   : %s\n", config.ChainID)
	fmt.Printf("Gas Limit                  : %d\n", gasLimit)
	fmt.Printf("Gas Price                  : %s wei\n", config.GasPrice)
	fmt.Println()

	return &CollateralTransactionData{
		To:       config.ValidatorManagerContractAddr,
		Value:    fee.String(),
		Data:     fmt.Sprintf("0x%x", data),
		GasLimit: gasLimit,
		GasPrice: config.GasPrice,
		ChainID:  config.ChainID,
	}, nil
}

func outputCollateralTransactionData(txData *CollateralTransactionData, config *config.ResolvedConfig) error {
	if config.OutputFile != "" {
		// Write to file
		data, err := json.MarshalIndent(txData, "", "  ")
		if err != nil {
			return fmt.Errorf("failed to marshal transaction data: %w", err)
		}

		err = os.WriteFile(config.OutputFile, data, 0644)
		if err != nil {
			return fmt.Errorf("failed to write transaction data to file: %w", err)
		}

		fmt.Printf("Transaction data written to: %s\n", config.OutputFile)
	} else {
		// Print to stdout
		data, err := json.MarshalIndent(txData, "", "  ")
		if err != nil {
			return fmt.Errorf("failed to marshal transaction data: %w", err)
		}

		fmt.Println("===== Transaction Data =====")
		fmt.Println(string(data))
	}

	return nil
}

// Validation functions
func validateCollateralDepositFields(config *config.ResolvedConfig, collateralFlags *struct {
	validator string
	amount    string
}) error {
	if err := validateCollateralBasicFields(config); err != nil {
		return err
	}

	if collateralFlags.validator == "" {
		return fmt.Errorf("validator address is required (use --validator)")
	}
	if collateralFlags.amount == "" {
		return fmt.Errorf("amount is required (use --amount)")
	}

	// Validate amount format
	if _, err := utils.ParseValueAsWei(collateralFlags.amount); err != nil {
		return fmt.Errorf("invalid amount format: %w", err)
	}

	// Validate validator address
	if _, err := utils.ValidateAddress(collateralFlags.validator); err != nil {
		return fmt.Errorf("invalid validator address: %w", err)
	}

	return nil
}

func validateCollateralWithdrawFields(config *config.ResolvedConfig, collateralFlags *struct {
	validator string
	amount    string
	receiver  string
}) error {
	if err := validateCollateralBasicFields(config); err != nil {
		return err
	}

	if collateralFlags.validator == "" {
		return fmt.Errorf("validator address is required (use --validator)")
	}
	if collateralFlags.amount == "" {
		return fmt.Errorf("amount is required (use --amount)")
	}
	if collateralFlags.receiver == "" {
		return fmt.Errorf("receiver address is required (use --receiver)")
	}

	// Validate amount format
	if _, err := utils.ParseValueAsWei(collateralFlags.amount); err != nil {
		return fmt.Errorf("invalid amount format: %w", err)
	}

	// Validate addresses
	if _, err := utils.ValidateAddress(collateralFlags.validator); err != nil {
		return fmt.Errorf("invalid validator address: %w", err)
	}
	if _, err := utils.ValidateAddress(collateralFlags.receiver); err != nil {
		return fmt.Errorf("invalid receiver address: %w", err)
	}

	return nil
}

func validateCollateralSetPermittedOwnerFields(config *config.ResolvedConfig, collateralFlags *struct {
	validator       string
	collateralOwner string
	isPermitted     bool
}) error {
	if err := validateCollateralBasicFields(config); err != nil {
		return err
	}

	if collateralFlags.validator == "" {
		return fmt.Errorf("validator address is required (use --validator)")
	}
	if collateralFlags.collateralOwner == "" {
		return fmt.Errorf("collateral owner address is required (use --collateral-owner)")
	}

	// Validate addresses
	if _, err := utils.ValidateAddress(collateralFlags.validator); err != nil {
		return fmt.Errorf("invalid validator address: %w", err)
	}
	if _, err := utils.ValidateAddress(collateralFlags.collateralOwner); err != nil {
		return fmt.Errorf("invalid collateral owner address: %w", err)
	}

	return nil
}

func validateCollateralTransferOwnershipFields(config *config.ResolvedConfig, collateralFlags *struct {
	validator string
	newOwner  string
}) error {
	if err := validateCollateralBasicFields(config); err != nil {
		return err
	}

	if collateralFlags.validator == "" {
		return fmt.Errorf("validator address is required (use --validator)")
	}
	if collateralFlags.newOwner == "" {
		return fmt.Errorf("new owner address is required (use --new-owner)")
	}

	// Validate addresses
	if _, err := utils.ValidateAddress(collateralFlags.validator); err != nil {
		return fmt.Errorf("invalid validator address: %w", err)
	}
	if _, err := utils.ValidateAddress(collateralFlags.newOwner); err != nil {
		return fmt.Errorf("invalid new owner address: %w", err)
	}

	return nil
}

func validateCollateralBasicFields(config *config.ResolvedConfig) error {
	if config.ValidatorManagerContractAddr == "" {
		return fmt.Errorf("ValidatorManager contract address is required (use --contract or set with 'mito config set-contract')")
	}

	return nil
}
