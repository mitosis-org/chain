package send

import (
	"github.com/spf13/cobra"
)

// NewSendCmd returns the tx send command group
func NewSendCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "send",
		Short: "Create, sign and send transactions",
		Long:  "Create, sign and immediately send transactions to the network",
	}

	cmd.AddCommand(
		NewValidatorCmd(),
		NewCollateralCmd(),
	)

	return cmd
}
