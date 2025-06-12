package wallet

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

// NewDeleteCmd creates the wallet delete command
func NewDeleteCmd() *cobra.Command {
	var keystoreDir string
	var yes bool

	cmd := &cobra.Command{
		Use:   "delete [account_name]",
		Short: "Delete an encrypted keystore file",
		Long: `Delete an encrypted keystore file.

This will permanently remove the keystore file from disk.
WARNING: This action cannot be undone. Make sure you have a backup of your private key.`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			accountName := args[0]

			// Determine keystore directory
			var keystorePath string
			if keystoreDir == "" {
				homeDir, err := os.UserHomeDir()
				if err != nil {
					return fmt.Errorf("failed to get user home directory: %w", err)
				}
				keystorePath = filepath.Join(homeDir, ".mito", "keystores")
			} else {
				keystorePath = keystoreDir
			}

			// Check if keystore file exists
			keystoreFile := filepath.Join(keystorePath, accountName)
			if _, err := os.Stat(keystoreFile); os.IsNotExist(err) {
				return fmt.Errorf("keystore file not found: %s", keystoreFile)
			}

			// Show confirmation unless --yes flag is used
			if !yes {
				fmt.Printf("Are you sure you want to delete the keystore file '%s'? [y/N]: ", accountName)

				reader := bufio.NewReader(os.Stdin)
				response, err := reader.ReadString('\n')
				if err != nil {
					return fmt.Errorf("failed to read confirmation: %w", err)
				}

				response = strings.ToLower(strings.TrimSpace(response))
				if response != "y" && response != "yes" {
					fmt.Println("Deletion cancelled.")
					return nil
				}
			}

			// Delete the keystore file
			if err := os.Remove(keystoreFile); err != nil {
				return fmt.Errorf("failed to delete keystore file: %w", err)
			}

			fmt.Printf("Successfully deleted keystore file: %s\n", keystoreFile)

			return nil
		},
	}

	cmd.Flags().StringVar(&keystoreDir, "keystore-dir", "", "Directory containing the keystore file (default: ~/.mito/keystores)")
	cmd.Flags().BoolVarP(&yes, "yes", "y", false, "Skip confirmation prompt")

	return cmd
}
