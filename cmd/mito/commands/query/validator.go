package query

import (
	"github.com/mitosis-org/chain/cmd/mito/commands/query/validator"
	"github.com/spf13/cobra"
)

// NewValidatorCmd returns the validator command group
func NewValidatorCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "validator",
		Short: "Validator query commands",
		Long: `Commands for querying validator information

This command provides access to read-only operations for:
- validator info: Get detailed validator information`,
	}

	// Add subcommands
	cmd.AddCommand(
		validator.NewInfoCmd(),
	)

	return cmd
}
