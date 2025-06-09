package create

import (
	"github.com/mitosis-org/chain/cmd/mito/commands/tx/create/collateral"
	"github.com/spf13/cobra"
)

// NewCollateralCmd returns collateral commands for tx create
func NewCollateralCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "collateral",
		Short: "Create collateral transactions",
		Long: `Create collateral transactions (signed or unsigned)

This command provides access to collateral transaction creation operations:
- collateral deposit: Create collateral deposit transaction
- collateral withdraw: Create collateral withdraw transaction
- collateral set-permitted-owner: Create set permitted owner transaction
- collateral transfer-ownership: Create transfer ownership transaction`,
	}

	// Add subcommands
	cmd.AddCommand(
		collateral.NewDepositCmd(),
		collateral.NewWithdrawCmd(),
		collateral.NewSetPermittedOwnerCmd(),
		collateral.NewTransferOwnershipCmd(),
	)

	return cmd
}
