package create

import (
	"github.com/spf13/cobra"
)

// NewCreateCmd returns the tx create command group
func NewCreateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create transactions (signed or unsigned)",
		Long:  "Create transactions that can be signed immediately or later",
	}

	cmd.AddCommand(
		NewValidatorCmd(),
		NewCollateralCmd(),
	)

	return cmd
}
