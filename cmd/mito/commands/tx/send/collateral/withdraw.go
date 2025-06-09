package collateral

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

// NewWithdrawCmd creates the send collateral withdraw command
func NewWithdrawCmd() *cobra.Command {
	var commonFlags flags.CommonFlags
	var collateralFlags struct {
		validator string
		amount    string
		receiver  string
	}

	cmd := &cobra.Command{
		Use:   "withdraw",
		Short: "Send collateral withdraw transaction",
		Long:  "Create, sign and send a transaction to withdraw collateral",
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

			// Validate required fields
			if err := validateSendCollateralWithdrawFields(resolvedConfig, &collateralFlags); err != nil {
				return err
			}

			// Create container
			container, err := container.NewContainer(resolvedConfig)
			if err != nil {
				return fmt.Errorf("failed to initialize container: %w", err)
			}
			defer container.Close()

			// Create unsigned transaction
			tx, err := container.CollateralService.WithdrawCollateral(collateralFlags.validator, collateralFlags.amount, collateralFlags.receiver)
			if err != nil {
				return err
			}

			// Display transaction information using formatter
			formatter := output.NewTransactionFormatter("")
			info := &output.CollateralWithdrawInfo{
				ValidatorAddress: collateralFlags.validator,
				ReceiverAddress:  collateralFlags.receiver,
				CollateralAmount: collateralFlags.amount,
				Fee:              resolvedConfig.ContractFee,
			}

			if err := formatter.FormatCollateralWithdrawTransaction(tx, info); err != nil {
				return fmt.Errorf("failed to format transaction: %w", err)
			}
			fmt.Println()

			// Confirm action
			if !utils.ConfirmAction(resolvedConfig.Yes, "Withdraw collateral?") {
				return fmt.Errorf("operation cancelled")
			}

			// Sign transaction
			signedTx, err := container.TxBuilder.SignTransaction(tx)
			if err != nil {
				return fmt.Errorf("failed to sign transaction: %w", err)
			}

			// Send transaction and wait for confirmation
			_, err = container.TxSender.SendAndWait(signedTx)
			if err != nil {
				return err
			}

			fmt.Println("✅ Collateral withdrawn successfully!")
			return nil
		},
	}

	// Add flags
	flags.AddCommonFlags(cmd, &commonFlags)
	flags.AddSendFlags(cmd, &commonFlags)
	cmd.Flags().StringVar(&collateralFlags.validator, "validator", "", "Validator address")
	cmd.Flags().StringVar(&collateralFlags.amount, "amount", "", "Amount to withdraw in MITO (e.g., \"1.5\")")
	cmd.Flags().StringVar(&collateralFlags.receiver, "receiver", "", "Address to receive the withdrawn collateral")

	// Mark required flags
	cmd.MarkFlagRequired("validator")
	cmd.MarkFlagRequired("amount")
	cmd.MarkFlagRequired("receiver")

	return cmd
}

// Helper validation function
func validateSendCollateralWithdrawFields(config *config.ResolvedConfig, collateralFlags *struct {
	validator string
	amount    string
	receiver  string
}) error {
	if err := validateSendCollateralFields(config); err != nil {
		return err
	}

	if collateralFlags.validator == "" {
		return fmt.Errorf("validator address is required (use --validator)")
	}
	if collateralFlags.amount == "" {
		return fmt.Errorf("amount is required (use --amount)")
	}
	if collateralFlags.receiver == "" {
		return fmt.Errorf("receiver address is required (use --receiver)")
	}

	// Validate amount format
	if _, err := utils.ParseValueAsWei(collateralFlags.amount); err != nil {
		return fmt.Errorf("invalid amount format: %w", err)
	}

	// Validate addresses
	if _, err := utils.ValidateAddress(collateralFlags.validator); err != nil {
		return fmt.Errorf("invalid validator address: %w", err)
	}
	if _, err := utils.ValidateAddress(collateralFlags.receiver); err != nil {
		return fmt.Errorf("invalid receiver address: %w", err)
	}

	return nil
}
