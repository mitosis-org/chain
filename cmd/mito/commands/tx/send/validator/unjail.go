package validator

import (
	"fmt"

	"github.com/mitosis-org/chain/cmd/mito/internal/config"
	"github.com/mitosis-org/chain/cmd/mito/internal/container"
	"github.com/mitosis-org/chain/cmd/mito/internal/flags"
	"github.com/mitosis-org/chain/cmd/mito/internal/output"
	"github.com/mitosis-org/chain/cmd/mito/internal/utils"
	"github.com/mitosis-org/chain/cmd/mito/internal/validation"
	"github.com/spf13/cobra"
)

// NewUnjailCmd creates the send validator unjail command
func NewUnjailCmd() *cobra.Command {
	var commonFlags flags.CommonFlags
	var validatorAddr string

	cmd := &cobra.Command{
		Use:   "unjail",
		Short: "Send validator unjail transaction",
		Long:  "Create, sign and send a transaction to unjail a validator",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return validation.ValidateSendTxFlagGroups(cmd)
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			cmd.SilenceUsage = true

			// Resolve configuration
			resolver, err := config.NewResolver()
			if err != nil {
				return err
			}
			resolvedConfig := resolver.ResolveFlags(&commonFlags)

			// Validate fields using new validation module
			validateFields := &validation.ValidatorUpdateFields{
				ValidatorAddress: validatorAddr,
				FieldName:        "Jail Status",
				NewValue:         "Unjailed",
			}
			if err := validation.ValidateValidatorUpdateFields(resolvedConfig, validateFields); err != nil {
				return err
			}

			// Validate network and signing requirements for send
			if err := validation.ValidateNetworkFields(resolvedConfig, true); err != nil {
				return err
			}
			if err := validation.ValidateSigningFields(resolvedConfig, true); err != nil {
				return err
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

			// Display transaction information using formatter
			formatter := output.NewTransactionFormatter("")
			info := &output.ValidatorUpdateInfo{
				ValidatorAddress: validatorAddr,
				FieldName:        "Jail Status",
				NewValue:         "Unjailed",
			}

			if err := formatter.FormatValidatorUpdateTransaction(transaction, info); err != nil {
				return fmt.Errorf("failed to format transaction: %w", err)
			}
			fmt.Println()

			// Confirm action
			if !utils.ConfirmAction(resolvedConfig.Yes, "Unjail validator?") {
				return fmt.Errorf("operation cancelled")
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
	cmd.Flags().StringVar(&validatorAddr, "validator", "", "Validator address to unjail")
	cmd.MarkFlagRequired("validator")

	return cmd
}
