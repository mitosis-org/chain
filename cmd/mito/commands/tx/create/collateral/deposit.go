package collateral

import (
	"fmt"

	"github.com/mitosis-org/chain/cmd/mito/internal/config"
	"github.com/mitosis-org/chain/cmd/mito/internal/container"
	"github.com/mitosis-org/chain/cmd/mito/internal/flags"
	"github.com/mitosis-org/chain/cmd/mito/internal/output"
	"github.com/mitosis-org/chain/cmd/mito/internal/validation"
	"github.com/spf13/cobra"
)

// NewDepositCmd creates the create collateral deposit command
func NewDepositCmd() *cobra.Command {
	var commonFlags flags.CommonFlags
	var collateralFlags struct {
		validator string
		amount    string
	}

	cmd := &cobra.Command{
		Use:   "deposit",
		Short: "Create collateral deposit transaction",
		Long:  "Create a transaction to deposit collateral (without sending)",
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

			// Initialize container
			c, err := container.NewContainer(resolvedConfig)
			if err != nil {
				return fmt.Errorf("failed to initialize container: %w", err)
			}
			defer c.Close()

			tx, err := c.CollateralService.DepositCollateral(collateralFlags.validator, collateralFlags.amount)
			if err != nil {
				return err
			}

			// Sign if not unsigned
			if !commonFlags.Unsigned {
				tx, err = c.TxBuilder.SignTransaction(tx)
				if err != nil {
					return fmt.Errorf("failed to sign transaction: %w", err)
				}
			}

			// Create formatter and format transaction
			formatter := output.NewTransactionFormatter(commonFlags.OutputFile)

			return formatter.OutputTransaction(tx)
		},
	}

	// Add flags
	flags.AddCommonFlags(cmd, &commonFlags)
	flags.AddCreateFlags(cmd, &commonFlags)
	cmd.Flags().StringVar(&collateralFlags.validator, "validator", "", "Validator address")
	cmd.Flags().StringVar(&collateralFlags.amount, "amount", "", "Amount to deposit in MITO (e.g., \"1.5\")")

	// Mark required flags
	cmd.MarkFlagRequired("validator")
	cmd.MarkFlagRequired("amount")

	return cmd
}
