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

// NewTransferOwnershipCmd creates the send collateral transfer-ownership command
func NewTransferOwnershipCmd() *cobra.Command {
	var commonFlags flags.CommonFlags
	var collateralFlags struct {
		validator string
		newOwner  string
	}

	cmd := &cobra.Command{
		Use:   "transfer-ownership",
		Short: "Send transfer ownership transaction",
		Long:  "Create, sign and send a transaction to transfer collateral ownership",
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
			if err := validateSendCollateralTransferOwnershipFields(resolvedConfig, &collateralFlags); err != nil {
				return err
			}

			// Create container
			container, err := container.NewContainer(resolvedConfig)
			if err != nil {
				return fmt.Errorf("failed to initialize container: %w", err)
			}
			defer container.Close()

			// Create unsigned transaction
			tx, err := container.CollateralService.TransferCollateralOwnership(collateralFlags.validator, collateralFlags.newOwner)
			if err != nil {
				return err
			}

			// Display transaction information using formatter
			formatter := output.NewTransactionFormatter("")
			info := &output.CollateralOwnershipInfo{
				ValidatorAddress: collateralFlags.validator,
				NewOwner:         collateralFlags.newOwner,
				Fee:              resolvedConfig.ContractFee,
			}

			if err := formatter.FormatCollateralOwnershipTransaction(tx, info); err != nil {
				return fmt.Errorf("failed to format transaction: %w", err)
			}
			fmt.Println()

			// Confirm action
			if !utils.ConfirmAction(resolvedConfig.Yes, "Transfer collateral ownership?") {
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

			fmt.Println("✅ Collateral ownership transferred successfully!")
			return nil
		},
	}

	// Add flags
	flags.AddCommonFlags(cmd, &commonFlags)
	cmd.Flags().StringVar(&collateralFlags.validator, "validator", "", "Validator address")
	cmd.Flags().StringVar(&collateralFlags.newOwner, "new-owner", "", "New owner address")

	// Mark required flags
	cmd.MarkFlagRequired("validator")
	cmd.MarkFlagRequired("new-owner")

	return cmd
}

// Helper validation function
func validateSendCollateralTransferOwnershipFields(config *config.ResolvedConfig, collateralFlags *struct {
	validator string
	newOwner  string
}) error {
	if err := validateSendCollateralFields(config); err != nil {
		return err
	}

	if collateralFlags.validator == "" {
		return fmt.Errorf("validator address is required (use --validator)")
	}
	if collateralFlags.newOwner == "" {
		return fmt.Errorf("new owner address is required (use --new-owner)")
	}

	// Validate addresses
	if _, err := utils.ValidateAddress(collateralFlags.validator); err != nil {
		return fmt.Errorf("invalid validator address: %w", err)
	}
	if _, err := utils.ValidateAddress(collateralFlags.newOwner); err != nil {
		return fmt.Errorf("invalid new owner address: %w", err)
	}

	return nil
}
