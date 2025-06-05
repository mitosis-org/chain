package validator

import (
	"fmt"

	"github.com/mitosis-org/chain/cmd/mito/internal/config"
	"github.com/mitosis-org/chain/cmd/mito/internal/container"
	"github.com/mitosis-org/chain/cmd/mito/internal/flags"
	"github.com/mitosis-org/chain/cmd/mito/internal/output"
	"github.com/mitosis-org/chain/cmd/mito/internal/validation"
	"github.com/spf13/cobra"
)

// NewUpdateRewardConfigCmd creates the create validator update-reward-config command
func NewUpdateRewardConfigCmd() *cobra.Command {
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
	flags.AddCreateFlags(cmd, &commonFlags)
	cmd.Flags().StringVar(&validatorAddr, "validator", "", "Validator address")
	cmd.Flags().StringVar(&commissionRate, "commission-rate", "", "New commission rate in percentage (e.g., \"5%\")")
	cmd.MarkFlagRequired("validator")
	cmd.MarkFlagRequired("commission-rate")

	return cmd
}
