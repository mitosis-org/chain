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

// NewUpdateOperatorCmd creates the send validator update-operator command
func NewUpdateOperatorCmd() *cobra.Command {
	var commonFlags flags.CommonFlags
	var validatorAddr, operator string

	cmd := &cobra.Command{
		Use:   "update-operator",
		Short: "Send validator operator update transaction",
		Long:  "Create, sign and send a transaction to update validator operator",
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
				FieldName:        "Operator",
				NewValue:         operator,
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
			transaction, err := container.ValidatorService.UpdateOperator(validatorAddr, operator)
			if err != nil {
				return err
			}

			// Display transaction information using formatter
			formatter := output.NewTransactionFormatter("")
			info := &output.ValidatorUpdateInfo{
				ValidatorAddress: validatorAddr,
				FieldName:        "Operator",
				NewValue:         operator,
			}

			if err := formatter.FormatValidatorUpdateTransaction(transaction, info); err != nil {
				return fmt.Errorf("failed to format transaction: %w", err)
			}
			fmt.Println()

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
	flags.AddSendFlags(cmd, &commonFlags)
	cmd.Flags().StringVar(&validatorAddr, "validator", "", "Validator address")
	cmd.Flags().StringVar(&operator, "operator", "", "New operator address")
	cmd.MarkFlagRequired("validator")
	cmd.MarkFlagRequired("operator")

	return cmd
}
