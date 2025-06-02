package send

import (
	"github.com/mitosis-org/chain/cmd/mito/commands/tx/send/collateral"
	"github.com/spf13/cobra"
)

// NewCollateralCmd returns collateral commands for tx send
func NewCollateralCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "collateral",
		Short: "Send collateral transactions",
		Long: `Create, sign and send collateral transactions to the network

This command provides access to collateral transaction operations:
- collateral deposit: Send collateral deposit transaction
- collateral withdraw: Send collateral withdraw transaction
- collateral set-permitted-owner: Send set permitted owner transaction
- collateral transfer-ownership: Send transfer ownership transaction`,
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
