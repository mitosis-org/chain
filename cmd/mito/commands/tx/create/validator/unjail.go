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

// NewUnjailCmd creates the create validator unjail command
func NewUnjailCmd() *cobra.Command {
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
	flags.AddCreateFlags(cmd, &commonFlags)
	cmd.Flags().StringVar(&validatorAddr, "validator", "", "Address of the validator to unjail")
	cmd.MarkFlagRequired("validator")

	return cmd
}
