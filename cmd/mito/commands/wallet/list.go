package wallet

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

// NewListCmd creates the wallet list command
func NewListCmd() *cobra.Command {
	var dir string

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List all the accounts in the keystore default directory",
		Long:  `List all the accounts in the keystore default directory ~/.mito/keystores.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			// Get keystore directory
			keystoreDir := dir
			if keystoreDir == "" {
				homeDir, err := os.UserHomeDir()
				if err != nil {
					return fmt.Errorf("failed to get home directory: %w", err)
				}
				keystoreDir = filepath.Join(homeDir, ".mito", "keystores")
			}

			// Check if keystore directory exists
			if _, err := os.Stat(keystoreDir); os.IsNotExist(err) {
				// Cast doesn't print anything if directory doesn't exist
				return nil
			}

			// Read directory contents
			files, err := os.ReadDir(keystoreDir)
			if err != nil {
				return fmt.Errorf("failed to read keystore directory: %w", err)
			}

			// Filter for keystore files and extract account names
			for _, file := range files {
				if file.IsDir() {
					continue
				}

				fileName := file.Name()
				// Don't check if it's a "valid" keystore file - just include all files
				// This matches cast's behavior more closely
				fmt.Printf("%s (Local)\n", fileName)
			}

			return nil
		},
	}

	cmd.Flags().StringVar(&dir, "dir", "", "List all the accounts in the keystore directory. Default keystore directory (~/.mito/keystores) is used if no path provided")

	return cmd
}
