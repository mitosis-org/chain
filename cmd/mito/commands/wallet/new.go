package wallet

import (
	"crypto/ecdsa"
	"fmt"
	"os"
	"path/filepath"
	"syscall"

	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/spf13/cobra"
	"golang.org/x/term"
)

// NewNewCmd creates the wallet new command
func NewNewCmd() *cobra.Command {
	var unsafePassword string
	var keystoreDir string

	cmd := &cobra.Command{
		Use:   "new <account_name>",
		Short: "Create a new random keypair",
		Long: `Create a new random keypair.

If account_name is provided, the keypair will be saved to an encrypted JSON keystore.
Use --keystore-dir to specify the keystore directory (default: ~/.mito/keystores).`,
		Args: cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			// Generate a new private key
			privateKey, err := crypto.GenerateKey()
			if err != nil {
				return fmt.Errorf("failed to generate private key: %w", err)
			}

			address := crypto.PubkeyToAddress(privateKey.PublicKey)

			// If no account name specified, just print the keypair
			if len(args) == 0 {
				fmt.Printf("Successfully created new keypair.\n")
				fmt.Printf("Address:     %s\n", address.Hex())
				fmt.Printf("Private key: %s\n", formatPrivateKey(privateKey))
				return nil
			}

			// Account name provided, save to keystore
			accountName := args[0]

			// Determine keystore directory
			var keystorePath string
			if keystoreDir == "" {
				// Use default keystore directory
				homeDir, err := os.UserHomeDir()
				if err != nil {
					return fmt.Errorf("failed to get user home directory: %w", err)
				}
				keystorePath = filepath.Join(homeDir, ".mito", "keystores")
			} else {
				keystorePath = keystoreDir
			}

			// Create keystore directory if it doesn't exist
			if err := os.MkdirAll(keystorePath, 0o755); err != nil {
				return fmt.Errorf("failed to create keystore directory: %w", err)
			}

			// Get password
			var password string
			if unsafePassword != "" {
				password = unsafePassword
			} else {
				password, err = promptPassword()
				if err != nil {
					return fmt.Errorf("failed to get password: %w", err)
				}
			}

			// Create keystore
			ks := keystore.NewKeyStore(keystorePath, keystore.StandardScryptN, keystore.StandardScryptP)
			account, err := ks.ImportECDSA(privateKey, password)
			if err != nil {
				return fmt.Errorf("failed to import key to keystore: %w", err)
			}

			// Rename to account name
			finalPath := filepath.Join(keystorePath, accountName)
			if err := os.Rename(account.URL.Path, finalPath); err != nil {
				return fmt.Errorf("failed to rename keystore file: %w", err)
			}

			fmt.Printf("Created new encrypted keystore file: %s\n", finalPath)
			fmt.Printf("Address: %s\n", address.Hex())

			return nil
		},
	}

	cmd.Flags().StringVar(&unsafePassword, "unsafe-password", "", "Password for the JSON keystore in cleartext (unsafe)")
	cmd.Flags().StringVar(&keystoreDir, "keystore-dir", "", "Directory to store the keystore file (default: ~/.mito/keystores)")

	return cmd
}

// formatPrivateKey formats the private key as hex string
func formatPrivateKey(key *ecdsa.PrivateKey) string {
	return fmt.Sprintf("0x%x", crypto.FromECDSA(key))
}

// promptPassword prompts user for password securely
func promptPassword() (string, error) {
	fmt.Print("Enter password: ")
	password, err := term.ReadPassword(syscall.Stdin)
	fmt.Println() // Add newline after password input
	if err != nil {
		return "", err
	}
	return string(password), nil
}
