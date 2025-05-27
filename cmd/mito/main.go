package main

import (
	"fmt"
	"os"

	"github.com/mitosis-org/chain/cmd/mito/cmd"
	"github.com/spf13/cobra"
)

func main() {
	// Root command
	rootCmd := &cobra.Command{
		Use:   "mito",
		Short: "Mitosis chain utilities",
		Long: `Mitosis chain utilities.

This command provides various utilities for interacting with
the Mitosis blockchain, including EVM components, validator management
and transaction handling.`,
	}

	// Add commands
	rootCmd.AddCommand(cmd.NewTxCmd())
	rootCmd.AddCommand(cmd.NewValidatorCmd())
	rootCmd.AddCommand(cmd.NewConfigCmd())

	// Execute the root command
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
