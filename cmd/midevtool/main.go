package main

import (
	"fmt"
	"os"

	"github.com/mitosis-org/chain/cmd/midevtool/cmd"
	"github.com/spf13/cobra"
)

func main() {
	// Root command
	rootCmd := &cobra.Command{
		Use:   "midevtool",
		Short: "Development and testing tools for Mitosis Chain (Not production-purpose)",
	}

	// Add commands
	rootCmd.AddCommand(cmd.NewGovernanceCmd())

	// Execute the root command
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
