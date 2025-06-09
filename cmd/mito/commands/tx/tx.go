package tx

import (
	"github.com/mitosis-org/chain/cmd/mito/commands/tx/create"
	"github.com/mitosis-org/chain/cmd/mito/commands/tx/send"
	"github.com/spf13/cobra"
)

// NewTxCmd returns the transaction command group
func NewTxCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "tx",
		Short: "Transaction commands",
		Long:  "Commands for creating and sending transactions",
	}

	cmd.AddCommand(
		create.NewCreateCmd(),
		send.NewSendCmd(),
	)

	return cmd
}
