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

// NewSetPermittedOwnerCmd creates the send collateral set-permitted-owner command
func NewSetPermittedOwnerCmd() *cobra.Command {
	var commonFlags flags.CommonFlags
	var collateralFlags struct {
		validator       string
		collateralOwner string
		isPermitted     bool
	}

	cmd := &cobra.Command{
		Use:   "set-permitted-owner",
		Short: "Send set permitted owner transaction",
		Long:  "Create, sign and send a transaction to set permitted collateral owner",
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
			if err := validateSendCollateralSetPermittedOwnerFields(resolvedConfig, &collateralFlags); err != nil {
				return err
			}

			// Create container
			container, err := container.NewContainer(resolvedConfig)
			if err != nil {
				return fmt.Errorf("failed to initialize container: %w", err)
			}
			defer container.Close()

			// Create unsigned transaction
			tx, err := container.CollateralService.SetPermittedCollateralOwner(collateralFlags.validator, collateralFlags.collateralOwner, collateralFlags.isPermitted)
			if err != nil {
				return err
			}

			// Display transaction information using formatter
			formatter := output.NewTransactionFormatter("")
			info := &output.CollateralPermissionInfo{
				ValidatorAddress: collateralFlags.validator,
				CollateralOwner:  collateralFlags.collateralOwner,
				IsPermitted:      collateralFlags.isPermitted,
			}

			if err := formatter.FormatCollateralPermissionTransaction(tx, info); err != nil {
				return fmt.Errorf("failed to format transaction: %w", err)
			}
			fmt.Println()

			// Confirm action
			permissionText := "deny"
			if collateralFlags.isPermitted {
				permissionText = "permit"
			}
			confirmMsg := fmt.Sprintf("Set collateral owner permission to %s?", permissionText)
			if !utils.ConfirmAction(resolvedConfig.Yes, confirmMsg) {
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

			fmt.Println("✅ Collateral owner permission updated successfully!")
			return nil
		},
	}

	// Add flags
	flags.AddCommonFlags(cmd, &commonFlags)
	flags.AddSendFlags(cmd, &commonFlags)
	cmd.Flags().StringVar(&collateralFlags.validator, "validator", "", "Validator address")
	cmd.Flags().StringVar(&collateralFlags.collateralOwner, "collateral-owner", "", "Collateral owner address to permit or deny")
	cmd.Flags().BoolVar(&collateralFlags.isPermitted, "permitted", false, "Whether to permit (true) or deny (false) the address")

	// Mark required flags
	cmd.MarkFlagRequired("validator")
	cmd.MarkFlagRequired("collateral-owner")

	return cmd
}

// Helper validation function
func validateSendCollateralSetPermittedOwnerFields(config *config.ResolvedConfig, collateralFlags *struct {
	validator       string
	collateralOwner string
	isPermitted     bool
},
) error {
	if err := validateSendCollateralFields(config); err != nil {
		return err
	}

	if collateralFlags.validator == "" {
		return fmt.Errorf("validator address is required (use --validator)")
	}
	if collateralFlags.collateralOwner == "" {
		return fmt.Errorf("collateral owner address is required (use --collateral-owner)")
	}

	// Validate addresses
	if _, err := utils.ValidateAddress(collateralFlags.validator); err != nil {
		return fmt.Errorf("invalid validator address: %w", err)
	}
	if _, err := utils.ValidateAddress(collateralFlags.collateralOwner); err != nil {
		return fmt.Errorf("invalid collateral owner address: %w", err)
	}

	return nil
}
