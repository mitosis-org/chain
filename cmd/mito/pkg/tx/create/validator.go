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
	"github.com/mitosis-org/chain/cmd/mito/internal/tx"
	"github.com/mitosis-org/chain/cmd/mito/internal/utils"
	"github.com/spf13/cobra"
)

// NewValidatorCmd returns validator commands for tx create
func NewValidatorCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "validator",
		Short: "Create validator transactions",
		Long:  "Create validator transactions (signed or unsigned)",
	}

	cmd.AddCommand(
		newCreateValidatorCreateCmd(),
		newCreateValidatorUpdateMetadataCmd(),
		newCreateValidatorUpdateOperatorCmd(),
		newCreateValidatorUpdateRewardConfigCmd(),
		newCreateValidatorUpdateRewardManagerCmd(),
		newCreateValidatorUnjailCmd(),
	)

	return cmd
}

func newCreateValidatorCreateCmd() *cobra.Command {
	var commonFlags flags.CommonFlags
	var validatorFlags struct {
		pubkey            string
		operator          string
		rewardManager     string
		commissionRate    string
		metadata          string
		initialCollateral string
	}

	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a new validator transaction",
		Long:  "Create a new validator transaction (without sending)",
		RunE: func(cmd *cobra.Command, args []string) error {
			cmd.SilenceUsage = true

			// Resolve configuration
			resolver, err := config.NewResolver()
			if err != nil {
				return err
			}
			resolvedConfig := resolver.ResolveFlags(&commonFlags)

			// Validate required fields
			if err := validateCreateValidatorCreateFields(resolvedConfig, &validatorFlags); err != nil {
				return err
			}

			// Create container (for validation and fee calculation)
			container, err := container.NewContainer(resolvedConfig)
			if err != nil {
				return fmt.Errorf("failed to initialize container: %w", err)
			}
			defer container.Close()

			// Create validator request for validation
			req := &tx.CreateValidatorRequest{
				PubKey:            validatorFlags.pubkey,
				Operator:          validatorFlags.operator,
				RewardManager:     validatorFlags.rewardManager,
				CommissionRate:    validatorFlags.commissionRate,
				Metadata:          validatorFlags.metadata,
				InitialCollateral: validatorFlags.initialCollateral,
			}

			// Validate and get transaction details (without executing)
			txData, err := createValidatorTransaction(container, req, resolvedConfig)
			if err != nil {
				return err
			}

			// Output transaction data
			return outputTransactionData(txData, resolvedConfig)
		},
	}

	// Add flags
	flags.AddCommonFlags(cmd, &commonFlags)
	cmd.Flags().StringVar(&validatorFlags.pubkey, "pubkey", "", "Validator's public key (hex with 0x prefix)")
	cmd.Flags().StringVar(&validatorFlags.operator, "operator", "", "Operator address")
	cmd.Flags().StringVar(&validatorFlags.rewardManager, "reward-manager", "", "Reward manager address")
	cmd.Flags().StringVar(&validatorFlags.commissionRate, "commission-rate", "", "Commission rate in percentage (e.g., \"5%\")")
	cmd.Flags().StringVar(&validatorFlags.metadata, "metadata", "", "Validator metadata (JSON string)")
	cmd.Flags().StringVar(&validatorFlags.initialCollateral, "initial-collateral", "", "Initial collateral amount in MITO (e.g., \"1.5\")")

	// Mark required flags
	cmd.MarkFlagRequired("pubkey")
	cmd.MarkFlagRequired("operator")
	cmd.MarkFlagRequired("reward-manager")
	cmd.MarkFlagRequired("commission-rate")
	cmd.MarkFlagRequired("metadata")
	cmd.MarkFlagRequired("initial-collateral")

	return cmd
}

func newCreateValidatorUpdateMetadataCmd() *cobra.Command {
	var commonFlags flags.CommonFlags
	var validatorAddr, metadata string

	cmd := &cobra.Command{
		Use:   "update-metadata",
		Short: "Create validator metadata update transaction",
		Long:  "Create a transaction to update validator metadata (without sending)",
		RunE: func(cmd *cobra.Command, args []string) error {
			cmd.SilenceUsage = true

			// Resolve configuration
			resolver, err := config.NewResolver()
			if err != nil {
				return err
			}
			resolvedConfig := resolver.ResolveFlags(&commonFlags)

			// Validate required fields
			if err := validateCreateValidatorBasicFields(resolvedConfig); err != nil {
				return err
			}

			// Create container
			container, err := container.NewContainer(resolvedConfig)
			if err != nil {
				return fmt.Errorf("failed to initialize container: %w", err)
			}
			defer container.Close()

			// Create transaction data
			txData, err := createUpdateMetadataTransaction(container, validatorAddr, metadata, resolvedConfig)
			if err != nil {
				return err
			}

			// Output transaction data
			return outputTransactionData(txData, resolvedConfig)
		},
	}

	// Add flags
	flags.AddCommonFlags(cmd, &commonFlags)
	cmd.Flags().StringVar(&validatorAddr, "validator", "", "Validator address")
	cmd.Flags().StringVar(&metadata, "metadata", "", "New validator metadata")
	cmd.MarkFlagRequired("validator")
	cmd.MarkFlagRequired("metadata")

	return cmd
}

func newCreateValidatorUpdateOperatorCmd() *cobra.Command {
	var commonFlags flags.CommonFlags
	var validatorAddr, operator string

	cmd := &cobra.Command{
		Use:   "update-operator",
		Short: "Create validator operator update transaction",
		Long:  "Create a transaction to update validator operator (without sending)",
		RunE: func(cmd *cobra.Command, args []string) error {
			cmd.SilenceUsage = true

			// Resolve configuration
			resolver, err := config.NewResolver()
			if err != nil {
				return err
			}
			resolvedConfig := resolver.ResolveFlags(&commonFlags)

			// Validate required fields
			if err := validateCreateValidatorBasicFields(resolvedConfig); err != nil {
				return err
			}

			// Create container
			container, err := container.NewContainer(resolvedConfig)
			if err != nil {
				return fmt.Errorf("failed to initialize container: %w", err)
			}
			defer container.Close()

			// Create transaction data
			txData, err := createUpdateOperatorTransaction(container, validatorAddr, operator, resolvedConfig)
			if err != nil {
				return err
			}

			// Output transaction data
			return outputTransactionData(txData, resolvedConfig)
		},
	}

	// Add flags
	flags.AddCommonFlags(cmd, &commonFlags)
	cmd.Flags().StringVar(&validatorAddr, "validator", "", "Validator address")
	cmd.Flags().StringVar(&operator, "operator", "", "New operator address")
	cmd.MarkFlagRequired("validator")
	cmd.MarkFlagRequired("operator")

	return cmd
}

func newCreateValidatorUpdateRewardConfigCmd() *cobra.Command {
	var commonFlags flags.CommonFlags
	var validatorAddr, commissionRate string

	cmd := &cobra.Command{
		Use:   "update-reward-config",
		Short: "Create validator reward config update transaction",
		Long:  "Create a transaction to update validator reward configuration (without sending)",
		RunE: func(cmd *cobra.Command, args []string) error {
			cmd.SilenceUsage = true

			// Resolve configuration
			resolver, err := config.NewResolver()
			if err != nil {
				return err
			}
			resolvedConfig := resolver.ResolveFlags(&commonFlags)

			// Validate required fields
			if err := validateCreateValidatorBasicFields(resolvedConfig); err != nil {
				return err
			}

			// Create container
			container, err := container.NewContainer(resolvedConfig)
			if err != nil {
				return fmt.Errorf("failed to initialize container: %w", err)
			}
			defer container.Close()

			// Create transaction data
			txData, err := createUpdateRewardConfigTransaction(container, validatorAddr, commissionRate, resolvedConfig)
			if err != nil {
				return err
			}

			// Output transaction data
			return outputTransactionData(txData, resolvedConfig)
		},
	}

	// Add flags
	flags.AddCommonFlags(cmd, &commonFlags)
	cmd.Flags().StringVar(&validatorAddr, "validator", "", "Validator address")
	cmd.Flags().StringVar(&commissionRate, "commission-rate", "", "New commission rate in percentage (e.g., \"5%\")")
	cmd.MarkFlagRequired("validator")
	cmd.MarkFlagRequired("commission-rate")

	return cmd
}

func newCreateValidatorUpdateRewardManagerCmd() *cobra.Command {
	var commonFlags flags.CommonFlags
	var validatorAddr, rewardManager string

	cmd := &cobra.Command{
		Use:   "update-reward-manager",
		Short: "Create validator reward manager update transaction",
		Long:  "Create a transaction to update validator reward manager (without sending)",
		RunE: func(cmd *cobra.Command, args []string) error {
			cmd.SilenceUsage = true

			// Resolve configuration
			resolver, err := config.NewResolver()
			if err != nil {
				return err
			}
			resolvedConfig := resolver.ResolveFlags(&commonFlags)

			// Validate required fields
			if err := validateCreateValidatorBasicFields(resolvedConfig); err != nil {
				return err
			}

			// Create container
			container, err := container.NewContainer(resolvedConfig)
			if err != nil {
				return fmt.Errorf("failed to initialize container: %w", err)
			}
			defer container.Close()

			// Create transaction data
			txData, err := createUpdateRewardManagerTransaction(container, validatorAddr, rewardManager, resolvedConfig)
			if err != nil {
				return err
			}

			// Output transaction data
			return outputTransactionData(txData, resolvedConfig)
		},
	}

	// Add flags
	flags.AddCommonFlags(cmd, &commonFlags)
	cmd.Flags().StringVar(&validatorAddr, "validator", "", "Validator address")
	cmd.Flags().StringVar(&rewardManager, "reward-manager", "", "New reward manager address")
	cmd.MarkFlagRequired("validator")
	cmd.MarkFlagRequired("reward-manager")

	return cmd
}

func newCreateValidatorUnjailCmd() *cobra.Command {
	var commonFlags flags.CommonFlags
	var validatorAddr string

	cmd := &cobra.Command{
		Use:   "unjail",
		Short: "Create validator unjail transaction",
		Long:  "Create a transaction to unjail a validator (without sending)",
		RunE: func(cmd *cobra.Command, args []string) error {
			cmd.SilenceUsage = true

			// Resolve configuration
			resolver, err := config.NewResolver()
			if err != nil {
				return err
			}
			resolvedConfig := resolver.ResolveFlags(&commonFlags)

			// Validate required fields
			if err := validateCreateValidatorBasicFields(resolvedConfig); err != nil {
				return err
			}

			// Create container
			container, err := container.NewContainer(resolvedConfig)
			if err != nil {
				return fmt.Errorf("failed to initialize container: %w", err)
			}
			defer container.Close()

			// Create transaction data
			txData, err := createUnjailValidatorTransaction(container, validatorAddr, resolvedConfig)
			if err != nil {
				return err
			}

			// Output transaction data
			return outputTransactionData(txData, resolvedConfig)
		},
	}

	// Add flags
	flags.AddCommonFlags(cmd, &commonFlags)
	cmd.Flags().StringVar(&validatorAddr, "validator", "", "Address of the validator to unjail")
	cmd.MarkFlagRequired("validator")

	return cmd
}

// Transaction data structure for output
type TransactionData struct {
	To       string `json:"to"`
	Value    string `json:"value"`
	Data     string `json:"data"`
	GasLimit uint64 `json:"gasLimit"`
	GasPrice string `json:"gasPrice"`
	ChainID  string `json:"chainId"`
}

// Helper functions for creating transaction data without sending
func createValidatorTransaction(container *container.Container, req *tx.CreateValidatorRequest, config *config.ResolvedConfig) (*TransactionData, error) {
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

	// Parse collateral amount as decimal MITO and convert to wei
	collateralAmount, err := utils.ParseValueAsWei(req.InitialCollateral)
	if err != nil {
		return nil, fmt.Errorf("failed to parse initial collateral: %w", err)
	}

	// Calculate total transaction value (collateral + fee)
	totalValue := new(big.Int).Add(collateralAmount, fee)

	// Validate addresses
	operatorAddr, err := utils.ValidateAddress(req.Operator)
	if err != nil {
		return nil, fmt.Errorf("invalid operator address: %w", err)
	}

	rewardManagerAddr, err := utils.ValidateAddress(req.RewardManager)
	if err != nil {
		return nil, fmt.Errorf("invalid reward manager address: %w", err)
	}

	// Parse commission rate
	commissionRateInt, err := utils.ParsePercentageToBasisPoints(req.CommissionRate)
	if err != nil {
		return nil, fmt.Errorf("failed to parse commission rate: %w", err)
	}

	// Decode public key from hex
	pubKeyBytes, err := utils.DecodeHexWithPrefix(req.PubKey)
	if err != nil {
		return nil, fmt.Errorf("failed to decode public key: %w", err)
	}

	// Create the request struct
	request := bindings.IValidatorManagerCreateValidatorRequest{
		Operator:       operatorAddr,
		RewardManager:  rewardManagerAddr,
		CommissionRate: commissionRateInt,
		Metadata:       []byte(req.Metadata),
	}

	// Get contract ABI and encode function call data
	abi, err := bindings.IValidatorManagerMetaData.GetAbi()
	if err != nil {
		return nil, fmt.Errorf("failed to get contract ABI: %w", err)
	}

	data, err := abi.Pack("createValidator", pubKeyBytes, request)
	if err != nil {
		return nil, fmt.Errorf("failed to pack function call: %w", err)
	}

	// Set default gas limit if not provided
	gasLimit := config.GasLimit
	if gasLimit == 0 {
		gasLimit = 500000 // Default gas limit like the original version
	}

	// Show summary (like the original version)
	fmt.Println("===== Create Validator Transaction =====")
	fmt.Printf("Public Key                 : %s\n", req.PubKey)
	fmt.Printf("Operator                   : %s\n", operatorAddr.Hex())
	fmt.Printf("Reward Manager             : %s\n", rewardManagerAddr.Hex())
	fmt.Printf("Commission Rate            : %s\n", utils.FormatBasisPointsToPercent(commissionRateInt))
	fmt.Printf("Metadata                   : %s\n", req.Metadata)
	fmt.Printf("Initial Collateral         : %s MITO\n", utils.FormatWeiToEther(collateralAmount))
	fmt.Printf("Fee                        : %s MITO\n", utils.FormatWeiToEther(fee))
	fmt.Printf("Total Value                : %s MITO\n", utils.FormatWeiToEther(totalValue))
	fmt.Printf("Chain ID                   : %s\n", config.ChainID)
	fmt.Printf("Gas Limit                  : %d\n", gasLimit)
	fmt.Printf("Gas Price                  : %s wei\n", config.GasPrice)
	fmt.Println()

	return &TransactionData{
		To:       config.ValidatorManagerContractAddr,
		Value:    totalValue.String(),
		Data:     fmt.Sprintf("0x%x", data),
		GasLimit: gasLimit,
		GasPrice: config.GasPrice,
		ChainID:  config.ChainID,
	}, nil
}

func createUpdateMetadataTransaction(container *container.Container, validatorAddr, metadata string, config *config.ResolvedConfig) (*TransactionData, error) {
	// Validate validator address
	valAddr, err := utils.ValidateAddress(validatorAddr)
	if err != nil {
		return nil, fmt.Errorf("invalid validator address: %w", err)
	}

	// Get contract ABI and encode function call data
	abi, err := bindings.IValidatorManagerMetaData.GetAbi()
	if err != nil {
		return nil, fmt.Errorf("failed to get contract ABI: %w", err)
	}

	data, err := abi.Pack("updateMetadata", valAddr, []byte(metadata))
	if err != nil {
		return nil, fmt.Errorf("failed to pack function call: %w", err)
	}

	// Set default gas limit if not provided
	gasLimit := config.GasLimit
	if gasLimit == 0 {
		gasLimit = 500000
	}

	// Show summary
	fmt.Println("===== Update Metadata Transaction =====")
	fmt.Printf("Validator                  : %s\n", valAddr.Hex())
	fmt.Printf("New Metadata               : %s\n", metadata)
	fmt.Printf("Chain ID                   : %s\n", config.ChainID)
	fmt.Printf("Gas Limit                  : %d\n", gasLimit)
	fmt.Printf("Gas Price                  : %s wei\n", config.GasPrice)
	fmt.Println()

	return &TransactionData{
		To:       config.ValidatorManagerContractAddr,
		Value:    "0",
		Data:     fmt.Sprintf("0x%x", data),
		GasLimit: gasLimit,
		GasPrice: config.GasPrice,
		ChainID:  config.ChainID,
	}, nil
}

func createUpdateOperatorTransaction(container *container.Container, validatorAddr, operator string, config *config.ResolvedConfig) (*TransactionData, error) {
	// Validate addresses
	valAddr, err := utils.ValidateAddress(validatorAddr)
	if err != nil {
		return nil, fmt.Errorf("invalid validator address: %w", err)
	}

	operatorAddr, err := utils.ValidateAddress(operator)
	if err != nil {
		return nil, fmt.Errorf("invalid operator address: %w", err)
	}

	// Get contract ABI and encode function call data
	abi, err := bindings.IValidatorManagerMetaData.GetAbi()
	if err != nil {
		return nil, fmt.Errorf("failed to get contract ABI: %w", err)
	}

	data, err := abi.Pack("updateOperator", valAddr, operatorAddr)
	if err != nil {
		return nil, fmt.Errorf("failed to pack function call: %w", err)
	}

	// Set default gas limit if not provided
	gasLimit := config.GasLimit
	if gasLimit == 0 {
		gasLimit = 500000
	}

	// Show summary
	fmt.Println("===== Update Operator Transaction =====")
	fmt.Printf("Validator                  : %s\n", valAddr.Hex())
	fmt.Printf("New Operator               : %s\n", operatorAddr.Hex())
	fmt.Printf("Chain ID                   : %s\n", config.ChainID)
	fmt.Printf("Gas Limit                  : %d\n", gasLimit)
	fmt.Printf("Gas Price                  : %s wei\n", config.GasPrice)
	fmt.Println()

	return &TransactionData{
		To:       config.ValidatorManagerContractAddr,
		Value:    "0",
		Data:     fmt.Sprintf("0x%x", data),
		GasLimit: gasLimit,
		GasPrice: config.GasPrice,
		ChainID:  config.ChainID,
	}, nil
}

func createUpdateRewardConfigTransaction(container *container.Container, validatorAddr, commissionRate string, config *config.ResolvedConfig) (*TransactionData, error) {
	// Validate validator address
	valAddr, err := utils.ValidateAddress(validatorAddr)
	if err != nil {
		return nil, fmt.Errorf("invalid validator address: %w", err)
	}

	// Parse commission rate
	commissionRateInt, err := utils.ParsePercentageToBasisPoints(commissionRate)
	if err != nil {
		return nil, fmt.Errorf("failed to parse commission rate: %w", err)
	}

	// Create the request struct
	request := bindings.IValidatorManagerUpdateRewardConfigRequest{
		CommissionRate: commissionRateInt,
	}

	// Get contract ABI and encode function call data
	abi, err := bindings.IValidatorManagerMetaData.GetAbi()
	if err != nil {
		return nil, fmt.Errorf("failed to get contract ABI: %w", err)
	}

	data, err := abi.Pack("updateRewardConfig", valAddr, request)
	if err != nil {
		return nil, fmt.Errorf("failed to pack function call: %w", err)
	}

	// Set default gas limit if not provided
	gasLimit := config.GasLimit
	if gasLimit == 0 {
		gasLimit = 500000
	}

	// Show summary
	fmt.Println("===== Update Reward Config Transaction =====")
	fmt.Printf("Validator                  : %s\n", valAddr.Hex())
	fmt.Printf("New Commission Rate        : %s\n", utils.FormatBasisPointsToPercent(commissionRateInt))
	fmt.Printf("Chain ID                   : %s\n", config.ChainID)
	fmt.Printf("Gas Limit                  : %d\n", gasLimit)
	fmt.Printf("Gas Price                  : %s wei\n", config.GasPrice)
	fmt.Println()

	return &TransactionData{
		To:       config.ValidatorManagerContractAddr,
		Value:    "0",
		Data:     fmt.Sprintf("0x%x", data),
		GasLimit: gasLimit,
		GasPrice: config.GasPrice,
		ChainID:  config.ChainID,
	}, nil
}

func createUpdateRewardManagerTransaction(container *container.Container, validatorAddr, rewardManager string, config *config.ResolvedConfig) (*TransactionData, error) {
	// Validate addresses
	valAddr, err := utils.ValidateAddress(validatorAddr)
	if err != nil {
		return nil, fmt.Errorf("invalid validator address: %w", err)
	}

	rewardManagerAddr, err := utils.ValidateAddress(rewardManager)
	if err != nil {
		return nil, fmt.Errorf("invalid reward manager address: %w", err)
	}

	// Get contract ABI and encode function call data
	abi, err := bindings.IValidatorManagerMetaData.GetAbi()
	if err != nil {
		return nil, fmt.Errorf("failed to get contract ABI: %w", err)
	}

	data, err := abi.Pack("updateRewardManager", valAddr, rewardManagerAddr)
	if err != nil {
		return nil, fmt.Errorf("failed to pack function call: %w", err)
	}

	// Set default gas limit if not provided
	gasLimit := config.GasLimit
	if gasLimit == 0 {
		gasLimit = 500000
	}

	// Show summary
	fmt.Println("===== Update Reward Manager Transaction =====")
	fmt.Printf("Validator                  : %s\n", valAddr.Hex())
	fmt.Printf("New Reward Manager         : %s\n", rewardManagerAddr.Hex())
	fmt.Printf("Chain ID                   : %s\n", config.ChainID)
	fmt.Printf("Gas Limit                  : %d\n", gasLimit)
	fmt.Printf("Gas Price                  : %s wei\n", config.GasPrice)
	fmt.Println()

	return &TransactionData{
		To:       config.ValidatorManagerContractAddr,
		Value:    "0",
		Data:     fmt.Sprintf("0x%x", data),
		GasLimit: gasLimit,
		GasPrice: config.GasPrice,
		ChainID:  config.ChainID,
	}, nil
}

func createUnjailValidatorTransaction(container *container.Container, validatorAddr string, config *config.ResolvedConfig) (*TransactionData, error) {
	// Validate validator address
	valAddr, err := utils.ValidateAddress(validatorAddr)
	if err != nil {
		return nil, fmt.Errorf("invalid validator address: %w", err)
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

	// Get contract ABI and encode function call data
	abi, err := bindings.IValidatorManagerMetaData.GetAbi()
	if err != nil {
		return nil, fmt.Errorf("failed to get contract ABI: %w", err)
	}

	data, err := abi.Pack("unjailValidator", valAddr)
	if err != nil {
		return nil, fmt.Errorf("failed to pack function call: %w", err)
	}

	// Set default gas limit if not provided
	gasLimit := config.GasLimit
	if gasLimit == 0 {
		gasLimit = 500000
	}

	// Show summary
	fmt.Println("===== Unjail Validator Transaction =====")
	fmt.Printf("Validator                  : %s\n", valAddr.Hex())
	fmt.Printf("Fee                        : %s MITO\n", utils.FormatWeiToEther(fee))
	fmt.Printf("Chain ID                   : %s\n", config.ChainID)
	fmt.Printf("Gas Limit                  : %d\n", gasLimit)
	fmt.Printf("Gas Price                  : %s wei\n", config.GasPrice)
	fmt.Println()

	return &TransactionData{
		To:       config.ValidatorManagerContractAddr,
		Value:    fee.String(),
		Data:     fmt.Sprintf("0x%x", data),
		GasLimit: gasLimit,
		GasPrice: config.GasPrice,
		ChainID:  config.ChainID,
	}, nil
}

func outputTransactionData(txData *TransactionData, config *config.ResolvedConfig) error {
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

func validateCreateValidatorCreateFields(config *config.ResolvedConfig, validatorFlags *struct {
	pubkey            string
	operator          string
	rewardManager     string
	commissionRate    string
	metadata          string
	initialCollateral string
}) error {
	if err := validateCreateValidatorBasicFields(config); err != nil {
		return err
	}

	if validatorFlags.pubkey == "" {
		return fmt.Errorf("public key is required (use --pubkey)")
	}
	if validatorFlags.operator == "" {
		return fmt.Errorf("operator address is required (use --operator)")
	}
	if validatorFlags.rewardManager == "" {
		return fmt.Errorf("reward manager address is required (use --reward-manager)")
	}
	if validatorFlags.commissionRate == "" {
		return fmt.Errorf("commission rate is required (use --commission-rate)")
	}
	if validatorFlags.metadata == "" {
		return fmt.Errorf("metadata is required (use --metadata)")
	}
	if validatorFlags.initialCollateral == "" {
		return fmt.Errorf("initial collateral is required (use --initial-collateral)")
	}

	// Validate address formats
	if _, err := utils.ValidateAddress(validatorFlags.operator); err != nil {
		return fmt.Errorf("invalid operator address: %w", err)
	}
	if _, err := utils.ValidateAddress(validatorFlags.rewardManager); err != nil {
		return fmt.Errorf("invalid reward manager address: %w", err)
	}

	return nil
}

func validateCreateValidatorBasicFields(config *config.ResolvedConfig) error {
	if config.RpcURL == "" {
		return fmt.Errorf("RPC URL is required (use --rpc-url or set with 'mito config set-rpc')")
	}
	if config.ValidatorManagerContractAddr == "" {
		return fmt.Errorf("ValidatorManager contract address is required (use --contract or set with 'mito config set-contract')")
	}

	return nil
}
