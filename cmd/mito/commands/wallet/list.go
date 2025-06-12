package wallet

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

// isValidKeystoreFile checks if a file is a valid Ethereum keystore file
func isValidKeystoreFile(filePath string) bool {
	// Read file content
	content, err := os.ReadFile(filePath)
	if err != nil {
		return false
	}

	// Try to parse as JSON
	var keystore map[string]interface{}
	if err := json.Unmarshal(content, &keystore); err != nil {
		return false
	}

	// Check for required keystore fields
	requiredFields := []string{"address", "crypto", "id", "version"}
	for _, field := range requiredFields {
		if _, exists := keystore[field]; !exists {
			return false
		}
	}

	return true
}

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

			// Filter for valid keystore files only
			for _, file := range files {
				if file.IsDir() {
					continue
				}

				fileName := file.Name()
				filePath := filepath.Join(keystoreDir, fileName)

				// Only show valid keystore files
				if isValidKeystoreFile(filePath) {
					fmt.Printf("%s\n", fileName)
				}
			}

			return nil
		},
	}

	cmd.Flags().StringVar(&dir, "keystore-dir", "", "List all the accounts in the keystore directory. Default keystore directory (~/.mito/keystores) is used if no path provided")

	return cmd
}
