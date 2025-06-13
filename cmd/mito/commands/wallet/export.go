package wallet

import (
	"crypto/ecdsa"
	"fmt"
	"os"
	"path/filepath"

	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/spf13/cobra"
)

// formatPublicKey formats the public key as hex string
func formatPublicKey(key *ecdsa.PublicKey) string {
	return fmt.Sprintf("0x%x", crypto.FromECDSAPub(key))
}

// NewExportCmd creates the wallet export command
func NewExportCmd() *cobra.Command {
	var keystoreDir string
	var unsafePassword string

	cmd := &cobra.Command{
		Use:   "export [account_name]",
		Short: "Export a private key from an encrypted keystore",
		Long: `Export a private key from an encrypted keystore file.

This will decrypt the keystore file and display the private key in hex format.
WARNING: This will display your private key in plain text. Use with caution.`,
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

			// Get password
			var password string
			if unsafePassword != "" {
				password = unsafePassword
			} else {
				var err error
				password, err = promptPassword()
				if err != nil {
					return fmt.Errorf("failed to get password: %w", err)
				}
			}

			// Read keystore file
			keystoreContent, err := os.ReadFile(keystoreFile)
			if err != nil {
				return fmt.Errorf("failed to read keystore file: %w", err)
			}

			// Decrypt the keystore
			key, err := keystore.DecryptKey(keystoreContent, password)
			if err != nil {
				return fmt.Errorf("failed to decrypt keystore (wrong password?): %w", err)
			}

			// Display the key information
			fmt.Printf("Private key: %s\n", formatPrivateKey(key.PrivateKey))
			fmt.Printf("Public key:  %s\n", formatPublicKey(&key.PrivateKey.PublicKey))
			fmt.Printf("Address:     %s\n", key.Address.Hex())

			return nil
		},
	}

	cmd.Flags().StringVar(&keystoreDir, "keystore-dir", "", "Directory containing the keystore file (default: ~/.mito/keystores)")
	cmd.Flags().StringVar(&unsafePassword, "unsafe-password", "", "Password for the JSON keystore in cleartext (unsafe)")

	return cmd
}
