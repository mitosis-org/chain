package query

import (
	"github.com/mitosis-org/chain/cmd/mito/commands/query/contract"
	"github.com/spf13/cobra"
)

// NewContractCmd returns the contract command group
func NewContractCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "contract",
		Short: "Query contract information",
		Long: `Commands for querying contract-related information

This command provides access to read-only operations for:
- contract fee: Query current contract fee settings`,
	}

	// Add subcommands
	cmd.AddCommand(
		contract.NewFeeCmd(),
	)

	return cmd
}
