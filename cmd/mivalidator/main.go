package main

import (
	"fmt"
	"os"

	"github.com/mitosis-org/chain/cmd/mivalidator/cmd"
	"github.com/spf13/cobra"
)

func main() {
	// Root command
	rootCmd := &cobra.Command{
		Use:   "mivalidator",
		Short: "Command-line tool for Mitosis Validator operations",
		Long: `A command-line interface for managing validators in the Mitosis network.
This tool allows operators to create validators, manage collateral, and update validator settings.`,
	}

	// Add commands
	rootCmd.AddCommand(cmd.NewValidatorInfoCmd())
	rootCmd.AddCommand(cmd.NewCreateValidatorCmd())
	rootCmd.AddCommand(cmd.NewDepositCollateralCmd())
	rootCmd.AddCommand(cmd.NewWithdrawCollateralCmd())
	rootCmd.AddCommand(cmd.NewUnjailValidatorCmd())
	rootCmd.AddCommand(cmd.NewUpdateOperatorCmd())
	rootCmd.AddCommand(cmd.NewUpdateWithdrawalRecipientCmd())
	rootCmd.AddCommand(cmd.NewUpdateRewardManagerCmd())
	rootCmd.AddCommand(cmd.NewUpdateRewardConfigCmd())
	rootCmd.AddCommand(cmd.NewUpdateMetadataCmd())

	// Execute the root command
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
