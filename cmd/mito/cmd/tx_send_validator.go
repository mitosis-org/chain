package cmd

import (
	"fmt"
	"math/big"

	"github.com/mitosis-org/chain/bindings"
	"github.com/spf13/cobra"
)

// newTxSendValidatorCreateCmd creates the validator creation command for tx send
func newTxSendValidatorCreateCmd() *cobra.Command {
	var (
		pubKey            string
		operator          string
		rewardManager     string
		commissionRate    string
		metadata          string
		initialCollateral string
	)

	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create, sign and send a validator creation transaction",
		Long: `Create, sign and immediately send a transaction to register a new validator in the ValidatorManager contract.
This command requires signing credentials (private key or keystore).`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runSendValidatorTx(cmd, pubKey, operator, rewardManager, commissionRate, metadata, initialCollateral)
		},
	}

	// Add common flags (signing required)
	AddCommonFlags(cmd, true)

	// Command-specific flags
	cmd.Flags().StringVar(&pubKey, "pubkey", "", "Validator's public key (hex with 0x prefix)")
	cmd.Flags().StringVar(&operator, "operator", "", "Operator address")
	cmd.Flags().StringVar(&rewardManager, "reward-manager", "", "Reward manager address")
	cmd.Flags().StringVar(&commissionRate, "commission-rate", "", "Commission rate in percentage (e.g., \"5%\")")
	cmd.Flags().StringVar(&metadata, "metadata", "", "Validator metadata (JSON string)")
	cmd.Flags().StringVar(&initialCollateral, "initial-collateral", "", "Initial collateral amount in MITO (e.g., \"1.5\")")

	// Mark required flags
	mustMarkFlagRequired(cmd, "pubkey")
	mustMarkFlagRequired(cmd, "operator")
	mustMarkFlagRequired(cmd, "reward-manager")
	mustMarkFlagRequired(cmd, "commission-rate")
	mustMarkFlagRequired(cmd, "metadata")
	mustMarkFlagRequired(cmd, "initial-collateral")

	return cmd
}

func runSendValidatorTx(cmd *cobra.Command, pubKey, operator, rewardManager, commissionRate, metadata, initialCollateral string) error {
	// Get the contract fee
	fee, err := contract.Fee(nil)
	if err != nil {
		return fmt.Errorf("failed to get contract fee: %w", err)
	}

	// Get the config to check the initial deposit requirement
	config, err := contract.GlobalValidatorConfig(nil)
	if err != nil {
		return fmt.Errorf("failed to get global validator config: %w", err)
	}

	// Parse collateral amount as decimal MITO and convert to wei
	collateralAmount, err := ParseValueAsWei(initialCollateral)
	if err != nil {
		return fmt.Errorf("failed to parse initial collateral: %w", err)
	}

	// Ensure collateral is at least the initial deposit requirement
	if collateralAmount.Cmp(config.InitialValidatorDeposit) < 0 {
		return fmt.Errorf("initial collateral must be at least %s MITO",
			FormatWeiToEther(config.InitialValidatorDeposit))
	}

	// Calculate total transaction value (collateral + fee)
	totalValue := new(big.Int).Add(collateralAmount, fee)

	// Validate other parameters
	operatorAddr, err := ValidateAddress(operator)
	if err != nil {
		return fmt.Errorf("invalid operator address: %w", err)
	}

	rewardManagerAddr, err := ValidateAddress(rewardManager)
	if err != nil {
		return fmt.Errorf("invalid reward manager address: %w", err)
	}

	// Parse commission rate
	commissionRateInt, err := ParsePercentageToBasisPoints(commissionRate)
	if err != nil {
		return fmt.Errorf("failed to parse commission rate: %w", err)
	}

	// Validate commission rate
	maxRate, err := contract.MAXCOMMISSIONRATE(nil)
	if err != nil {
		return fmt.Errorf("failed to get max commission rate: %w", err)
	}

	if commissionRateInt.Cmp(big.NewInt(0)) < 0 || commissionRateInt.Cmp(maxRate) > 0 {
		return fmt.Errorf("commission rate must be between 0%% and %s", FormatBasisPointsToPercent(maxRate))
	}

	// Decode public key from hex
	pubKeyBytes, err := DecodeHexWithPrefix(pubKey)
	if err != nil {
		return fmt.Errorf("failed to decode public key: %w", err)
	}

	// Create the request
	request := bindings.IValidatorManagerCreateValidatorRequest{
		Operator:       operatorAddr,
		RewardManager:  rewardManagerAddr,
		CommissionRate: commissionRateInt,
		Metadata:       []byte(metadata),
	}

	// Show summary
	fmt.Println("===== Create Validator Transaction =====")
	fmt.Printf("Public Key                 : %s\n", pubKey)
	fmt.Printf("Operator                   : %s\n", operatorAddr.Hex())
	fmt.Printf("Reward Manager             : %s\n", rewardManagerAddr.Hex())
	fmt.Printf("Commission Rate            : %s\n", FormatBasisPointsToPercent(commissionRateInt))
	fmt.Printf("Metadata                   : %s\n", metadata)
	fmt.Printf("Initial Collateral         : %s MITO\n", FormatWeiToEther(collateralAmount))
	fmt.Printf("Fee                        : %s MITO\n", FormatWeiToEther(fee))
	fmt.Printf("Total Value                : %s MITO\n", FormatWeiToEther(totalValue))
	fmt.Println()

	// Ask for confirmation
	if !ConfirmAction("Do you want to create this validator?") {
		return fmt.Errorf("operation cancelled by user")
	}

	// Execute the transaction
	opts, err := TransactOpts(totalValue)
	if err != nil {
		return fmt.Errorf("failed to create transaction options: %w", err)
	}

	tx, err := contract.CreateValidator(opts, pubKeyBytes, request)
	if err != nil {
		return fmt.Errorf("failed to create validator: %w", err)
	}

	fmt.Printf("Transaction sent: %s\n", tx.Hash().Hex())

	// Wait for transaction confirmation
	err = WaitForTxConfirmation(ethClient, tx.Hash())
	if err != nil {
		return fmt.Errorf("transaction failed: %w", err)
	}

	fmt.Println("✅ Validator created successfully!")
	return nil
}

// newTxSendValidatorUpdateOperatorCmd creates the validator operator update command for tx send
func newTxSendValidatorUpdateOperatorCmd() *cobra.Command {
	var (
		validator string
		operator  string
	)

	cmd := &cobra.Command{
		Use:   "update-operator",
		Short: "Create, sign and send a validator operator update transaction",
		Long: `Create, sign and immediately send a transaction to update the operator address for an existing validator.
This command requires signing credentials (private key or keystore).`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runSendUpdateOperatorTx(cmd, validator, operator)
		},
	}

	// Add common flags (signing required)
	AddCommonFlags(cmd, true)

	// Command-specific flags
	cmd.Flags().StringVar(&validator, "validator", "", "Validator address")
	cmd.Flags().StringVar(&operator, "operator", "", "New operator address")

	// Mark required flags
	mustMarkFlagRequired(cmd, "validator")
	mustMarkFlagRequired(cmd, "operator")

	return cmd
}

func runSendUpdateOperatorTx(cmd *cobra.Command, validator, operator string) error {
	// Validate validator address
	valAddr, err := ValidateAddress(validator)
	if err != nil {
		return fmt.Errorf("invalid validator address: %w", err)
	}

	// Validate operator address
	operatorAddr, err := ValidateAddress(operator)
	if err != nil {
		return fmt.Errorf("invalid operator address: %w", err)
	}

	// Get validator info to show current values
	validatorInfo, err := contract.ValidatorInfo(nil, valAddr)
	if err != nil {
		return fmt.Errorf("failed to get validator info: %w", err)
	}

	// Show summary
	fmt.Println("===== Update Operator Transaction =====")
	fmt.Printf("Validator Address            : %s\n", valAddr.Hex())
	fmt.Printf("Current Operator             : %s\n", validatorInfo.Operator.Hex())
	fmt.Printf("New Operator                 : %s\n", operatorAddr.Hex())
	fmt.Printf("Current Reward Manager       : %s\n", validatorInfo.RewardManager.Hex())
	fmt.Println()

	// Show important warning
	fmt.Println("🚨 IMPORTANT WARNING 🚨")
	fmt.Println("When changing the operator address, you may also want to update:")
	fmt.Println("1. Reward Manager - To ensure rewards are managed by the correct entity")
	fmt.Println("2. Collateral Ownership - To manage who can deposit or withdraw collateral")
	fmt.Println()
	fmt.Println("After this operation completes, consider using:")
	fmt.Printf("  - evm tx send validator update-reward-manager --validator %s --reward-manager <new-address>\n", valAddr.Hex())
	fmt.Printf("  - evm tx send collateral set-permitted-owner --validator %s --collateral-owner <old-address> --permitted false\n", valAddr.Hex())
	fmt.Printf("  - evm tx send collateral set-permitted-owner --validator %s --collateral-owner <new-address> --permitted true\n", valAddr.Hex())
	fmt.Printf("  - evm tx send collateral transfer-ownership --validator %s --new-owner <new-address>\n", valAddr.Hex())
	fmt.Println()

	// Ask for confirmation
	if !ConfirmAction("Do you want to update the operator address?") {
		return fmt.Errorf("operation cancelled by user")
	}

	// Execute the transaction
	opts, err := TransactOpts(nil)
	if err != nil {
		return fmt.Errorf("failed to create transaction options: %w", err)
	}

	tx, err := contract.UpdateOperator(opts, valAddr, operatorAddr)
	if err != nil {
		return fmt.Errorf("failed to update operator: %w", err)
	}

	fmt.Printf("Transaction sent: %s\n", tx.Hash().Hex())

	// Wait for transaction confirmation
	err = WaitForTxConfirmation(ethClient, tx.Hash())
	if err != nil {
		return fmt.Errorf("transaction failed: %w", err)
	}

	fmt.Println("✅ Operator updated successfully!")
	return nil
}

// newTxSendValidatorUpdateMetadataCmd creates the validator metadata update command for tx send
func newTxSendValidatorUpdateMetadataCmd() *cobra.Command {
	var (
		validator string
		metadata  string
	)

	cmd := &cobra.Command{
		Use:   "update-metadata",
		Short: "Create, sign and send a validator metadata update transaction",
		Long: `Create, sign and immediately send a transaction to update the metadata for an existing validator.
This command requires signing credentials (private key or keystore).`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runSendUpdateMetadataTx(cmd, validator, metadata)
		},
	}

	// Add common flags (signing required)
	AddCommonFlags(cmd, true)

	// Command-specific flags
	cmd.Flags().StringVar(&validator, "validator", "", "Validator address")
	cmd.Flags().StringVar(&metadata, "metadata", "", "New metadata (JSON string)")

	// Mark required flags
	mustMarkFlagRequired(cmd, "validator")
	mustMarkFlagRequired(cmd, "metadata")

	return cmd
}

func runSendUpdateMetadataTx(cmd *cobra.Command, validator, metadata string) error {
	// Validate validator address
	valAddr, err := ValidateAddress(validator)
	if err != nil {
		return fmt.Errorf("invalid validator address: %w", err)
	}

	// Get validator info to show current values
	validatorInfo, err := contract.ValidatorInfo(nil, valAddr)
	if err != nil {
		return fmt.Errorf("failed to get validator info: %w", err)
	}

	// Show summary
	fmt.Println("===== Update Metadata Transaction =====")
	fmt.Printf("Validator Address        : %s\n", valAddr.Hex())
	fmt.Printf("Current Metadata         : %s\n", string(validatorInfo.Metadata))
	fmt.Printf("New Metadata             : %s\n", metadata)
	fmt.Println()

	// Ask for confirmation
	if !ConfirmAction("Do you want to update the metadata?") {
		return fmt.Errorf("operation cancelled by user")
	}

	// Execute the transaction
	opts, err := TransactOpts(nil)
	if err != nil {
		return fmt.Errorf("failed to create transaction options: %w", err)
	}

	tx, err := contract.UpdateMetadata(opts, valAddr, []byte(metadata))
	if err != nil {
		return fmt.Errorf("failed to update metadata: %w", err)
	}

	fmt.Printf("Transaction sent: %s\n", tx.Hash().Hex())

	// Wait for transaction confirmation
	err = WaitForTxConfirmation(ethClient, tx.Hash())
	if err != nil {
		return fmt.Errorf("transaction failed: %w", err)
	}

	fmt.Println("✅ Metadata updated successfully!")
	return nil
}

// newTxSendValidatorUnjailCmd creates the validator unjail command for tx send
func newTxSendValidatorUnjailCmd() *cobra.Command {
	var validator string

	cmd := &cobra.Command{
		Use:   "unjail",
		Short: "Create, sign and send a validator unjail transaction",
		Long: `Create, sign and immediately send a transaction to unjail a validator that has been jailed.
This command requires signing credentials (private key or keystore).`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runSendUnjailValidatorTx(cmd, validator)
		},
	}

	// Add common flags (signing required)
	AddCommonFlags(cmd, true)

	// Command-specific flags
	cmd.Flags().StringVar(&validator, "validator", "", "Address of the validator to unjail")

	// Mark required flags
	mustMarkFlagRequired(cmd, "validator")

	return cmd
}

func runSendUnjailValidatorTx(cmd *cobra.Command, validator string) error {
	// Get the contract fee
	fee, err := contract.Fee(nil)
	if err != nil {
		return fmt.Errorf("failed to get contract fee: %w", err)
	}

	// Validate validator address
	valAddr, err := ValidateAddress(validator)
	if err != nil {
		return fmt.Errorf("invalid validator address: %w", err)
	}

	// Show summary
	fmt.Println("===== Unjail Validator Transaction =====")
	fmt.Printf("Validator Address        : %s\n", valAddr.Hex())
	fmt.Printf("Fee                      : %s MITO\n", FormatWeiToEther(fee))
	fmt.Println()

	// Ask for confirmation
	if !ConfirmAction("Do you want to unjail this validator?") {
		return fmt.Errorf("operation cancelled by user")
	}

	// Execute the transaction
	opts, err := TransactOpts(fee)
	if err != nil {
		return fmt.Errorf("failed to create transaction options: %w", err)
	}

	tx, err := contract.UnjailValidator(opts, valAddr)
	if err != nil {
		return fmt.Errorf("failed to unjail validator: %w", err)
	}

	fmt.Printf("Transaction sent: %s\n", tx.Hash().Hex())

	// Wait for transaction confirmation
	err = WaitForTxConfirmation(ethClient, tx.Hash())
	if err != nil {
		return fmt.Errorf("transaction failed: %w", err)
	}

	fmt.Println("✅ Validator unjailed successfully!")
	return nil
}
