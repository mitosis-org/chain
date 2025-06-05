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

// NewUpdateOperatorCmd creates the create validator update-operator command
func NewUpdateOperatorCmd() *cobra.Command {
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

			return formatter.OutputTransaction(tx)
		},
	}

	// Add flags
	flags.AddCommonFlags(cmd, &commonFlags)
	flags.AddCreateFlags(cmd, &commonFlags)
	cmd.Flags().StringVar(&validatorAddr, "validator", "", "Validator address")
	cmd.Flags().StringVar(&operator, "operator", "", "New operator address")
	cmd.MarkFlagRequired("validator")
	cmd.MarkFlagRequired("operator")

	return cmd
}
