package wallet

import (
	"github.com/spf13/cobra"
)

// NewWalletCmd creates the wallet command with subcommands
func NewWalletCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "wallet",
		Short: "Wallet management commands",
		Long: `Commands for wallet management

This command provides wallet management operations including:
- wallet new: Create a new random keypair
- wallet new-mnemonic: Create a new mnemonic with a set number of words
- wallet import: Import a private key into an encrypted keystore
- wallet list: List all accounts in the keystore directory
- wallet export: Export a private key from an encrypted keystore
- wallet delete: Delete an encrypted keystore file`,
	}

	// Add subcommands
	cmd.AddCommand(
		NewNewCmd(),
		NewNewMnemonicCmd(),
		NewImportCmd(),
		NewListCmd(),
		NewExportCmd(),
		NewDeleteCmd(),
	)

	return cmd
}
