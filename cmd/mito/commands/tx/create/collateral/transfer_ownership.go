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

// NewTransferOwnershipCmd creates the create collateral transfer-ownership command
func NewTransferOwnershipCmd() *cobra.Command {
	var commonFlags flags.CommonFlags
	var collateralFlags struct {
		validator string
		newOwner  string
	}

	cmd := &cobra.Command{
		Use:   "transfer-ownership",
		Short: "Create transfer ownership transaction",
		Long:  "Create a transaction to transfer collateral ownership (without sending)",
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

			tx, err := c.CollateralService.TransferCollateralOwnership(collateralFlags.validator, collateralFlags.newOwner)
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

			info := &output.CollateralOwnershipInfo{
				ValidatorAddress: collateralFlags.validator,
				NewOwner:         collateralFlags.newOwner,
				Fee:              resolvedConfig.ContractFee,
			}

			return formatter.FormatCollateralOwnershipTransaction(tx, info)
		},
	}

	// Add flags
	flags.AddCommonFlags(cmd, &commonFlags)
	flags.AddCreateFlags(cmd, &commonFlags)
	cmd.Flags().StringVar(&collateralFlags.validator, "validator", "", "Validator address")
	cmd.Flags().StringVar(&collateralFlags.newOwner, "new-owner", "", "New owner address")

	// Mark required flags
	cmd.MarkFlagRequired("validator")
	cmd.MarkFlagRequired("new-owner")

	return cmd
}
