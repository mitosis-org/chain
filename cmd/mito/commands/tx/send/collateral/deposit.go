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

// NewDepositCmd creates the send collateral deposit command
func NewDepositCmd() *cobra.Command {
	var commonFlags flags.CommonFlags
	var collateralFlags struct {
		validator string
		amount    string
	}

	cmd := &cobra.Command{
		Use:   "deposit",
		Short: "Send collateral deposit transaction",
		Long:  "Create, sign and send a transaction to deposit collateral",
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
			if err := validateSendCollateralDepositFields(resolvedConfig, &collateralFlags); err != nil {
				return err
			}

			// Create container
			container, err := container.NewContainer(resolvedConfig)
			if err != nil {
				return fmt.Errorf("failed to initialize container: %w", err)
			}
			defer container.Close()

			// Create unsigned transaction
			tx, err := container.CollateralService.DepositCollateral(collateralFlags.validator, collateralFlags.amount)
			if err != nil {
				return err
			}

			// Display transaction information using formatter
			formatter := output.NewTransactionFormatter("")
			info := &output.CollateralDepositInfo{
				ValidatorAddress: collateralFlags.validator,
				CollateralAmount: collateralFlags.amount,
				Fee:              resolvedConfig.ContractFee,
				TotalValue:       tx.Value(),
			}

			if err := formatter.FormatCollateralDepositTransaction(tx, info); err != nil {
				return fmt.Errorf("failed to format transaction: %w", err)
			}
			fmt.Println()

			// Confirm action
			if !utils.ConfirmAction(resolvedConfig.Yes, "Deposit collateral?") {
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

			fmt.Println("✅ Collateral deposited successfully!")
			return nil
		},
	}

	// Add flags
	flags.AddCommonFlags(cmd, &commonFlags)
	flags.AddSendFlags(cmd, &commonFlags)
	cmd.Flags().StringVar(&collateralFlags.validator, "validator", "", "Validator address")
	cmd.Flags().StringVar(&collateralFlags.amount, "amount", "", "Amount to deposit in MITO (e.g., \"1.5\")")

	// Mark required flags
	cmd.MarkFlagRequired("validator")
	cmd.MarkFlagRequired("amount")

	return cmd
}

// Helper validation function
func validateSendCollateralDepositFields(config *config.ResolvedConfig, collateralFlags *struct {
	validator string
	amount    string
},
) error {
	if err := validateSendCollateralFields(config); err != nil {
		return err
	}

	if collateralFlags.validator == "" {
		return fmt.Errorf("validator address is required (use --validator)")
	}
	if collateralFlags.amount == "" {
		return fmt.Errorf("amount is required (use --amount)")
	}

	// Validate amount format
	if _, err := utils.ParseValueAsWei(collateralFlags.amount); err != nil {
		return fmt.Errorf("invalid amount format: %w", err)
	}

	// Validate validator address
	if _, err := utils.ValidateAddress(collateralFlags.validator); err != nil {
		return fmt.Errorf("invalid validator address: %w", err)
	}

	return nil
}

func validateSendCollateralFields(config *config.ResolvedConfig) error {
	// Validate network and signing requirements for send
	if err := validation.ValidateNetworkFields(config, true); err != nil {
		return err
	}
	if err := validation.ValidateSigningFields(config, true); err != nil {
		return err
	}
	return nil
}
