package cli

import (
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/spf13/cobra"
)

// GetTxCmd returns the transaction commands for x/evmvalidator
func GetTxCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:                        "evmvalidator",
		Short:                      "Evmvalidator subcommands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	// Add subcommands here when needed
	// cmd.AddCommand(...)

	return cmd
}
