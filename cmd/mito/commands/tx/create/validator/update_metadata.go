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

// NewUpdateMetadataCmd creates the create validator update-metadata command
func NewUpdateMetadataCmd() *cobra.Command {
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
