package commands

import (
	"github.com/mitosis-org/chain/cmd/mito/commands/config"
	"github.com/mitosis-org/chain/cmd/mito/commands/query"
	"github.com/mitosis-org/chain/cmd/mito/commands/tx/create"
	"github.com/mitosis-org/chain/cmd/mito/commands/tx/send"
	"github.com/mitosis-org/chain/cmd/mito/commands/version"
	"github.com/mitosis-org/chain/cmd/mito/commands/wallet"
	"github.com/spf13/cobra"
)

// NewRootCmd creates the root command with new architecture
func NewRootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "mito",
		Short: "Mitosis chain utilities",
		Long: `Mitosis chain utilities.

This command provides various utilities for interacting with
the Mitosis blockchain, including EVM components, validator management
and transaction handling.

Features:
- Online/Offline transaction support
- Comprehensive validator management
- Extensible query operations
- Wallet management utilities
- Clean and maintainable architecture`,
	}

	// Add subcommands
	cmd.AddCommand(
		// Transaction commands
		newTxCmd(),

		// Query commands
		newQueryCmd(),

		// Wallet commands
		wallet.NewWalletCmd(),

		// Existing commands
		config.NewConfigCmd(),
		version.NewVersionCmd(),
	)

	return cmd
}

// newTxCmd creates the transaction command with subcommands
func newTxCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "tx",
		Short: "Transaction commands",
		Long: `Commands for creating and sending transactions

This command provides access to transaction operations with support for:
- tx create: Create transactions (signed or unsigned) without sending
- tx send: Create, sign and immediately send transactions to the network`,
	}

	// Add subcommands
	cmd.AddCommand(
		newTxCreateCmd(),
		newTxSendCmd(),
	)

	return cmd
}

// newTxCreateCmd creates the tx create command
func newTxCreateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create transactions (signed or unsigned)",
		Long: `Create transactions that can be signed immediately or later

This command creates transactions without sending them to the network.
Supports both signed and unsigned transaction creation:
- Signed: Transaction is ready to be broadcast
- Unsigned: Transaction can be signed later (requires nonce)`,
	}

	// Add subcommands
	cmd.AddCommand(
		create.NewValidatorCmd(),
		create.NewCollateralCmd(),
	)

	return cmd
}

// newTxSendCmd creates the tx send command
func newTxSendCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "send",
		Short: "Create, sign and send transactions",
		Long: `Create, sign and immediately send transactions to the network

This command creates, signs and broadcasts transactions in one step.
Requires network connection and signing method. All transactions
created through this command are automatically signed.`,
	}

	// Add subcommands
	cmd.AddCommand(
		send.NewValidatorCmd(),
		send.NewCollateralCmd(),
	)

	return cmd
}

// newQueryCmd creates the query command with subcommands
func newQueryCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "query",
		Short: "Query blockchain data",
		Long: `Commands for querying blockchain data

This command provides access to read-only operations for:
- query validator: Query validator information
- query contract: Query contract-related information`,
	}

	// Add subcommands
	cmd.AddCommand(
		query.NewValidatorCmd(),
		query.NewContractCmd(),
	)

	return cmd
}
