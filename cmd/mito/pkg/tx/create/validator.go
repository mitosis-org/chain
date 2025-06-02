package create

import (
	"fmt"

	"github.com/mitosis-org/chain/cmd/mito/internal/config"
	"github.com/mitosis-org/chain/cmd/mito/internal/container"
	"github.com/mitosis-org/chain/cmd/mito/internal/flags"
	"github.com/mitosis-org/chain/cmd/mito/internal/output"
	"github.com/mitosis-org/chain/cmd/mito/internal/tx"
	"github.com/mitosis-org/chain/cmd/mito/internal/validation"
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
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return validation.ValidateCreateTxFlagGroups(cmd)
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			cmd.SilenceUsage = true

			// Resolve configuration
			resolver, err := config.NewResolver()
			if err != nil {
				return err
			}
			resolvedConfig := resolver.ResolveFlags(&commonFlags)

			// Validate validator create fields
			validateFields := &validation.ValidatorCreateFields{
				PubKey:            validatorFlags.pubkey,
				Operator:          validatorFlags.operator,
				RewardManager:     validatorFlags.rewardManager,
				CommissionRate:    validatorFlags.commissionRate,
				Metadata:          validatorFlags.metadata,
				InitialCollateral: validatorFlags.initialCollateral,
			}
			if err := validation.ValidateValidatorCreateFields(resolvedConfig, validateFields); err != nil {
				return err
			}

			// Create container (for validation and fee calculation)
			container, err := container.NewContainer(resolvedConfig)
			if err != nil {
				return fmt.Errorf("failed to initialize container: %w", err)
			}
			defer container.Close()

			// Create request
			req := &tx.CreateValidatorRequest{
				PubKey:            validatorFlags.pubkey,
				Operator:          validatorFlags.operator,
				RewardManager:     validatorFlags.rewardManager,
				CommissionRate:    validatorFlags.commissionRate,
				Metadata:          validatorFlags.metadata,
				InitialCollateral: validatorFlags.initialCollateral,
			}

			// Create transaction directly (validation included)
			tx, err := container.ValidatorService.CreateValidatorWithOptions(req, commonFlags.Unsigned)
			if err != nil {
				return fmt.Errorf("failed to create transaction: %w", err)
			}

			// Format and output transaction
			formatter := output.NewTransactionFormatter(commonFlags.OutputFile)
			info := &output.ValidatorCreateInfo{
				PubKey:            validatorFlags.pubkey,
				Operator:          validatorFlags.operator,
				RewardManager:     validatorFlags.rewardManager,
				CommissionRate:    validatorFlags.commissionRate,
				Metadata:          validatorFlags.metadata,
				InitialCollateral: validatorFlags.initialCollateral,
			}

			return formatter.FormatValidatorCreateTransaction(tx, info)
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
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return validation.ValidateCreateTxFlagGroups(cmd)
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			cmd.SilenceUsage = true

			// Resolve configuration
			resolver, err := config.NewResolver()
			if err != nil {
				return err
			}
			resolvedConfig := resolver.ResolveFlags(&commonFlags)

			// Validate validator update fields
			validateFields := &validation.ValidatorUpdateFields{
				ValidatorAddress: validatorAddr,
				FieldName:        "Metadata",
				NewValue:         metadata,
			}
			if err := validation.ValidateValidatorUpdateFields(resolvedConfig, validateFields); err != nil {
				return err
			}

			// Create container
			container, err := container.NewContainer(resolvedConfig)
			if err != nil {
				return fmt.Errorf("failed to initialize container: %w", err)
			}
			defer container.Close()

			// Create transaction directly (validation included)
			tx, err := container.ValidatorService.UpdateMetadataWithOptions(validatorAddr, metadata, commonFlags.Unsigned)
			if err != nil {
				return fmt.Errorf("failed to create transaction: %w", err)
			}

			// Format and output transaction
			formatter := output.NewTransactionFormatter(commonFlags.OutputFile)
			info := &output.ValidatorUpdateInfo{
				ValidatorAddress: validatorAddr,
				FieldName:        "Metadata",
				NewValue:         metadata,
			}

			return formatter.FormatValidatorUpdateTransaction(tx, info)
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
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return validation.ValidateCreateTxFlagGroups(cmd)
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			cmd.SilenceUsage = true

			// Resolve configuration
			resolver, err := config.NewResolver()
			if err != nil {
				return err
			}
			resolvedConfig := resolver.ResolveFlags(&commonFlags)

			// Validate validator update fields
			validateFields := &validation.ValidatorUpdateFields{
				ValidatorAddress: validatorAddr,
				FieldName:        "Operator",
				NewValue:         operator,
			}
			if err := validation.ValidateValidatorUpdateFields(resolvedConfig, validateFields); err != nil {
				return err
			}

			// Create container
			container, err := container.NewContainer(resolvedConfig)
			if err != nil {
				return fmt.Errorf("failed to initialize container: %w", err)
			}
			defer container.Close()

			// Create transaction directly (validation included)
			tx, err := container.ValidatorService.UpdateOperatorWithOptions(validatorAddr, operator, commonFlags.Unsigned)
			if err != nil {
				return fmt.Errorf("failed to create transaction: %w", err)
			}

			// Format and output transaction
			formatter := output.NewTransactionFormatter(commonFlags.OutputFile)
			info := &output.ValidatorUpdateInfo{
				ValidatorAddress: validatorAddr,
				FieldName:        "Operator",
				NewValue:         operator,
			}

			return formatter.FormatValidatorUpdateTransaction(tx, info)
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
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return validation.ValidateCreateTxFlagGroups(cmd)
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			cmd.SilenceUsage = true

			// Resolve configuration
			resolver, err := config.NewResolver()
			if err != nil {
				return err
			}
			resolvedConfig := resolver.ResolveFlags(&commonFlags)

			// Validate validator update fields
			validateFields := &validation.ValidatorUpdateFields{
				ValidatorAddress: validatorAddr,
				FieldName:        "Reward Config",
				NewValue:         commissionRate,
			}
			if err := validation.ValidateValidatorUpdateFields(resolvedConfig, validateFields); err != nil {
				return err
			}

			// Create container
			container, err := container.NewContainer(resolvedConfig)
			if err != nil {
				return fmt.Errorf("failed to initialize container: %w", err)
			}
			defer container.Close()

			// Create transaction directly (validation included)
			tx, err := container.ValidatorService.UpdateRewardConfigWithOptions(validatorAddr, commissionRate, commonFlags.Unsigned)
			if err != nil {
				return fmt.Errorf("failed to create transaction: %w", err)
			}

			// Format and output transaction
			formatter := output.NewTransactionFormatter(commonFlags.OutputFile)
			info := &output.ValidatorUpdateInfo{
				ValidatorAddress: validatorAddr,
				FieldName:        "Reward Config",
				NewValue:         commissionRate,
			}

			return formatter.FormatValidatorUpdateTransaction(tx, info)
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
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return validation.ValidateCreateTxFlagGroups(cmd)
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			cmd.SilenceUsage = true

			// Resolve configuration
			resolver, err := config.NewResolver()
			if err != nil {
				return err
			}
			resolvedConfig := resolver.ResolveFlags(&commonFlags)

			// Validate validator update fields
			validateFields := &validation.ValidatorUpdateFields{
				ValidatorAddress: validatorAddr,
				FieldName:        "Reward Manager",
				NewValue:         rewardManager,
			}
			if err := validation.ValidateValidatorUpdateFields(resolvedConfig, validateFields); err != nil {
				return err
			}

			// Create container
			container, err := container.NewContainer(resolvedConfig)
			if err != nil {
				return fmt.Errorf("failed to initialize container: %w", err)
			}
			defer container.Close()

			// Create transaction directly (validation included)
			tx, err := container.ValidatorService.UpdateRewardManagerWithOptions(validatorAddr, rewardManager, commonFlags.Unsigned)
			if err != nil {
				return fmt.Errorf("failed to create transaction: %w", err)
			}

			// Format and output transaction
			formatter := output.NewTransactionFormatter(commonFlags.OutputFile)
			info := &output.ValidatorUpdateInfo{
				ValidatorAddress: validatorAddr,
				FieldName:        "Reward Manager",
				NewValue:         rewardManager,
			}

			return formatter.FormatValidatorUpdateTransaction(tx, info)
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
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return validation.ValidateCreateTxFlagGroups(cmd)
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			cmd.SilenceUsage = true

			// Resolve configuration
			resolver, err := config.NewResolver()
			if err != nil {
				return err
			}
			resolvedConfig := resolver.ResolveFlags(&commonFlags)

			// Validate validator update fields
			validateFields := &validation.ValidatorUpdateFields{
				ValidatorAddress: validatorAddr,
				FieldName:        "Unjail",
				NewValue:         "unjail",
			}
			if err := validation.ValidateValidatorUpdateFields(resolvedConfig, validateFields); err != nil {
				return err
			}

			// Create container
			container, err := container.NewContainer(resolvedConfig)
			if err != nil {
				return fmt.Errorf("failed to initialize container: %w", err)
			}
			defer container.Close()

			// Create transaction directly (validation included)
			tx, err := container.ValidatorService.UnjailValidatorWithOptions(validatorAddr, commonFlags.Unsigned)
			if err != nil {
				return fmt.Errorf("failed to create transaction: %w", err)
			}

			// Format and output transaction
			formatter := output.NewTransactionFormatter(commonFlags.OutputFile)
			info := &output.ValidatorUpdateInfo{
				ValidatorAddress: validatorAddr,
				FieldName:        "Unjail",
				NewValue:         "unjail validator",
			}

			return formatter.FormatValidatorUpdateTransaction(tx, info)
		},
	}

	// Add flags
	flags.AddCommonFlags(cmd, &commonFlags)
	cmd.Flags().StringVar(&validatorAddr, "validator", "", "Address of the validator to unjail")
	cmd.MarkFlagRequired("validator")

	return cmd
}
