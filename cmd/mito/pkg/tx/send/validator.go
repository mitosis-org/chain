package send

import (
	"fmt"

	"github.com/mitosis-org/chain/cmd/mito/internal/config"
	"github.com/mitosis-org/chain/cmd/mito/internal/container"
	"github.com/mitosis-org/chain/cmd/mito/internal/flags"
	"github.com/mitosis-org/chain/cmd/mito/internal/tx"
	"github.com/mitosis-org/chain/cmd/mito/internal/utils"
	"github.com/spf13/cobra"
)

// NewValidatorCmd returns validator commands for tx send
func NewValidatorCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "validator",
		Short: "Send validator transactions",
		Long:  "Create, sign and send validator transactions to the network",
	}

	cmd.AddCommand(
		newSendValidatorCreateCmd(),
		newSendValidatorUpdateMetadataCmd(),
		newSendValidatorUpdateOperatorCmd(),
		newSendValidatorUpdateRewardConfigCmd(),
		newSendValidatorUpdateRewardManagerCmd(),
		newSendValidatorUnjailCmd(),
	)

	return cmd
}

func newSendValidatorCreateCmd() *cobra.Command {
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
		Short: "Create and send a new validator transaction",
		Long:  "Create, sign and send a new validator transaction to the network",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			// Validate mutually exclusive flags
			return flags.ValidateSigningFlags(cmd)
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			cmd.SilenceUsage = true

			// Resolve configuration
			resolver, err := config.NewResolver()
			if err != nil {
				return err
			}
			resolvedConfig := resolver.ResolveFlags(&commonFlags)

			// Validate required fields
			if err := validateSendValidatorCreateFields(resolvedConfig, &validatorFlags); err != nil {
				return err
			}

			// Confirm action
			if !utils.ConfirmAction(resolvedConfig.Yes, "Create and send validator transaction?") {
				return fmt.Errorf("operation cancelled")
			}

			// Create container
			container, err := container.NewContainer(resolvedConfig)
			if err != nil {
				return fmt.Errorf("failed to initialize container: %w", err)
			}
			defer container.Close()

			// Create validator request
			req := &tx.CreateValidatorRequest{
				PubKey:            validatorFlags.pubkey,
				Operator:          validatorFlags.operator,
				RewardManager:     validatorFlags.rewardManager,
				CommissionRate:    validatorFlags.commissionRate,
				Metadata:          validatorFlags.metadata,
				InitialCollateral: validatorFlags.initialCollateral,
			}

			// Create unsigned transaction
			transaction, err := container.ValidatorService.CreateValidator(req)
			if err != nil {
				return err
			}

			// Sign transaction
			signedTx, err := container.TxBuilder.SignTransaction(transaction)
			if err != nil {
				return fmt.Errorf("failed to sign transaction: %w", err)
			}

			// Send transaction and wait for confirmation
			_, err = container.TxSender.SendAndWait(signedTx)
			if err != nil {
				return err
			}

			fmt.Println("✅ Validator created successfully!")
			return nil
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

func newSendValidatorUpdateMetadataCmd() *cobra.Command {
	var commonFlags flags.CommonFlags
	var validatorAddr, metadata string

	cmd := &cobra.Command{
		Use:   "update-metadata",
		Short: "Send validator metadata update transaction",
		Long:  "Create, sign and send a transaction to update validator metadata",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return flags.ValidateSigningFlags(cmd)
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			cmd.SilenceUsage = true

			// Resolve configuration
			resolver, err := config.NewResolver()
			if err != nil {
				return err
			}
			resolvedConfig := resolver.ResolveFlags(&commonFlags)

			// Validate required fields
			if err := validateSendValidatorFields(resolvedConfig); err != nil {
				return err
			}

			// Confirm action
			if !utils.ConfirmAction(resolvedConfig.Yes, "Update validator metadata?") {
				return fmt.Errorf("operation cancelled")
			}

			// Create container
			container, err := container.NewContainer(resolvedConfig)
			if err != nil {
				return fmt.Errorf("failed to initialize container: %w", err)
			}
			defer container.Close()

			// Create unsigned transaction
			transaction, err := container.ValidatorService.UpdateMetadata(validatorAddr, metadata)
			if err != nil {
				return err
			}

			// Sign transaction
			signedTx, err := container.TxBuilder.SignTransaction(transaction)
			if err != nil {
				return fmt.Errorf("failed to sign transaction: %w", err)
			}

			// Send transaction and wait for confirmation
			_, err = container.TxSender.SendAndWait(signedTx)
			if err != nil {
				return err
			}

			fmt.Println("✅ Metadata updated successfully!")
			return nil
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

func newSendValidatorUpdateOperatorCmd() *cobra.Command {
	var commonFlags flags.CommonFlags
	var validatorAddr, operator string

	cmd := &cobra.Command{
		Use:   "update-operator",
		Short: "Send validator operator update transaction",
		Long:  "Create, sign and send a transaction to update validator operator",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return flags.ValidateSigningFlags(cmd)
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			cmd.SilenceUsage = true

			// Resolve configuration
			resolver, err := config.NewResolver()
			if err != nil {
				return err
			}
			resolvedConfig := resolver.ResolveFlags(&commonFlags)

			// Validate required fields
			if err := validateSendValidatorFields(resolvedConfig); err != nil {
				return err
			}

			// Show important warning
			fmt.Println("🚨 IMPORTANT WARNING 🚨")
			fmt.Println("When changing the operator address, you may also want to update:")
			fmt.Println("1. Reward Manager - To ensure rewards are managed by the correct entity")
			fmt.Println("2. Collateral Ownership - To manage who can deposit or withdraw collateral")
			fmt.Println()

			// Confirm action
			if !utils.ConfirmAction(resolvedConfig.Yes, "Update validator operator?") {
				return fmt.Errorf("operation cancelled")
			}

			// Create container
			container, err := container.NewContainer(resolvedConfig)
			if err != nil {
				return fmt.Errorf("failed to initialize container: %w", err)
			}
			defer container.Close()

			// Create unsigned transaction
			transaction, err := container.ValidatorService.UpdateOperator(validatorAddr, operator)
			if err != nil {
				return err
			}

			// Sign transaction
			signedTx, err := container.TxBuilder.SignTransaction(transaction)
			if err != nil {
				return fmt.Errorf("failed to sign transaction: %w", err)
			}

			// Send transaction and wait for confirmation
			_, err = container.TxSender.SendAndWait(signedTx)
			if err != nil {
				return err
			}

			fmt.Println("✅ Operator updated successfully!")
			return nil
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

func newSendValidatorUpdateRewardConfigCmd() *cobra.Command {
	var commonFlags flags.CommonFlags
	var validatorAddr, commissionRate string

	cmd := &cobra.Command{
		Use:   "update-reward-config",
		Short: "Send validator reward config update transaction",
		Long:  "Create, sign and send a transaction to update validator reward configuration",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return flags.ValidateSigningFlags(cmd)
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			cmd.SilenceUsage = true

			// Resolve configuration
			resolver, err := config.NewResolver()
			if err != nil {
				return err
			}
			resolvedConfig := resolver.ResolveFlags(&commonFlags)

			// Validate required fields
			if err := validateSendValidatorFields(resolvedConfig); err != nil {
				return err
			}

			// Confirm action
			if !utils.ConfirmAction(resolvedConfig.Yes, "Update validator reward configuration?") {
				return fmt.Errorf("operation cancelled")
			}

			// Create container
			container, err := container.NewContainer(resolvedConfig)
			if err != nil {
				return fmt.Errorf("failed to initialize container: %w", err)
			}
			defer container.Close()

			// Create unsigned transaction
			transaction, err := container.ValidatorService.UpdateRewardConfig(validatorAddr, commissionRate)
			if err != nil {
				return err
			}

			// Sign transaction
			signedTx, err := container.TxBuilder.SignTransaction(transaction)
			if err != nil {
				return fmt.Errorf("failed to sign transaction: %w", err)
			}

			// Send transaction and wait for confirmation
			_, err = container.TxSender.SendAndWait(signedTx)
			if err != nil {
				return err
			}

			fmt.Println("✅ Reward configuration updated successfully!")
			return nil
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

func newSendValidatorUpdateRewardManagerCmd() *cobra.Command {
	var commonFlags flags.CommonFlags
	var validatorAddr, rewardManager string

	cmd := &cobra.Command{
		Use:   "update-reward-manager",
		Short: "Send validator reward manager update transaction",
		Long:  "Create, sign and send a transaction to update validator reward manager",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return flags.ValidateSigningFlags(cmd)
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			cmd.SilenceUsage = true

			// Resolve configuration
			resolver, err := config.NewResolver()
			if err != nil {
				return err
			}
			resolvedConfig := resolver.ResolveFlags(&commonFlags)

			// Validate required fields
			if err := validateSendValidatorFields(resolvedConfig); err != nil {
				return err
			}

			// Confirm action
			if !utils.ConfirmAction(resolvedConfig.Yes, "Update validator reward manager?") {
				return fmt.Errorf("operation cancelled")
			}

			// Create container
			container, err := container.NewContainer(resolvedConfig)
			if err != nil {
				return fmt.Errorf("failed to initialize container: %w", err)
			}
			defer container.Close()

			// Create unsigned transaction
			transaction, err := container.ValidatorService.UpdateRewardManager(validatorAddr, rewardManager)
			if err != nil {
				return err
			}

			// Sign transaction
			signedTx, err := container.TxBuilder.SignTransaction(transaction)
			if err != nil {
				return fmt.Errorf("failed to sign transaction: %w", err)
			}

			// Send transaction and wait for confirmation
			_, err = container.TxSender.SendAndWait(signedTx)
			if err != nil {
				return err
			}

			fmt.Println("✅ Reward manager updated successfully!")
			return nil
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

func newSendValidatorUnjailCmd() *cobra.Command {
	var commonFlags flags.CommonFlags
	var validatorAddr string

	cmd := &cobra.Command{
		Use:   "unjail",
		Short: "Send validator unjail transaction",
		Long:  "Create, sign and send a transaction to unjail a validator",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return flags.ValidateSigningFlags(cmd)
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			cmd.SilenceUsage = true

			// Resolve configuration
			resolver, err := config.NewResolver()
			if err != nil {
				return err
			}
			resolvedConfig := resolver.ResolveFlags(&commonFlags)

			// Validate required fields
			if err := validateSendValidatorFields(resolvedConfig); err != nil {
				return err
			}

			// Confirm action
			if !utils.ConfirmAction(resolvedConfig.Yes, "Unjail validator?") {
				return fmt.Errorf("operation cancelled")
			}

			// Create container
			container, err := container.NewContainer(resolvedConfig)
			if err != nil {
				return fmt.Errorf("failed to initialize container: %w", err)
			}
			defer container.Close()

			// Create unsigned transaction
			transaction, err := container.ValidatorService.UnjailValidator(validatorAddr)
			if err != nil {
				return err
			}

			// Sign transaction
			signedTx, err := container.TxBuilder.SignTransaction(transaction)
			if err != nil {
				return fmt.Errorf("failed to sign transaction: %w", err)
			}

			// Send transaction and wait for confirmation
			_, err = container.TxSender.SendAndWait(signedTx)
			if err != nil {
				return err
			}

			fmt.Println("✅ Validator unjailed successfully!")
			return nil
		},
	}

	// Add flags
	flags.AddCommonFlags(cmd, &commonFlags)
	cmd.Flags().StringVar(&validatorAddr, "validator", "", "Address of the validator to unjail")
	cmd.MarkFlagRequired("validator")

	return cmd
}

func validateSendValidatorCreateFields(config *config.ResolvedConfig, validatorFlags *struct {
	pubkey            string
	operator          string
	rewardManager     string
	commissionRate    string
	metadata          string
	initialCollateral string
}) error {
	if err := validateSendValidatorFields(config); err != nil {
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

func validateSendValidatorFields(config *config.ResolvedConfig) error {
	if config.RpcURL == "" {
		return fmt.Errorf("RPC URL is required (use --rpc-url or set with 'mito config set-rpc')")
	}
	if config.ValidatorManagerContractAddr == "" {
		return fmt.Errorf("ValidatorManager contract address is required (use --contract or set with 'mito config set-contract')")
	}

	// Validate signing method is provided
	if !config.HasSigningMethod() {
		return fmt.Errorf("signing method is required (use --private-key or --keyfile)")
	}

	return nil
}
