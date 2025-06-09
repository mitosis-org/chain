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

// NewSetPermittedOwnerCmd creates the create collateral set-permitted-owner command
func NewSetPermittedOwnerCmd() *cobra.Command {
	var commonFlags flags.CommonFlags
	var collateralFlags struct {
		validator       string
		collateralOwner string
		isPermitted     bool
	}

	cmd := &cobra.Command{
		Use:   "set-permitted-owner",
		Short: "Create set permitted owner transaction",
		Long:  "Create a transaction to set permitted collateral owner (without sending)",
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

			tx, err := c.CollateralService.SetPermittedCollateralOwner(collateralFlags.validator, collateralFlags.collateralOwner, collateralFlags.isPermitted)
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
	cmd.Flags().StringVar(&collateralFlags.collateralOwner, "collateral-owner", "", "Collateral owner address to permit or deny")
	cmd.Flags().BoolVar(&collateralFlags.isPermitted, "permitted", false, "Whether to permit (true) or deny (false) the address")

	// Mark required flags
	cmd.MarkFlagRequired("validator")
	cmd.MarkFlagRequired("collateral-owner")

	return cmd
}
