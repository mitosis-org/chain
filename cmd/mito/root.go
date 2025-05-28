package main

import (
	"github.com/mitosis-org/chain/cmd/mito/pkg/config"
	"github.com/mitosis-org/chain/cmd/mito/pkg/tx"
	"github.com/mitosis-org/chain/cmd/mito/pkg/validator"
	"github.com/mitosis-org/chain/cmd/mito/pkg/version"
	"github.com/spf13/cobra"
)

// NewRootCmd creates the root command
func NewRootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "mito",
		Short: "Mitosis chain utilities",
		Long: `Mitosis chain utilities.

This command provides various utilities for interacting with
the Mitosis blockchain, including EVM components, validator management
and transaction handling.`,
	}

	// Add subcommands
	cmd.AddCommand(
		config.NewConfigCmd(),
		tx.NewTxCmd(),
		validator.NewValidatorCmd(),
		version.NewVersionCmd(),
	)

	return cmd
}
