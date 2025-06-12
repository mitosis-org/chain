package wallet

import (
	"crypto/ecdsa"
	"fmt"
	"os"
	"path/filepath"
	"syscall"

	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/google/uuid"
	"github.com/spf13/cobra"
	"golang.org/x/term"
)

// NewNewCmd creates the wallet new command
func NewNewCmd() *cobra.Command {
	var (
		unsafePassword string
	)

	cmd := &cobra.Command{
		Use:   "new [path] [account_name]",
		Short: "Create a new random keypair",
		Long: `Create a new random keypair.

If path is specified, then keypair will be written to an encrypted JSON keystore.
If account_name is provided, the keystore file will be named using this account name.`,
		Args: cobra.MaximumNArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			// Generate a new private key
			privateKey, err := crypto.GenerateKey()
			if err != nil {
				return fmt.Errorf("failed to generate private key: %w", err)
			}

			address := crypto.PubkeyToAddress(privateKey.PublicKey)

			// If no path specified, just print the keypair
			if len(args) == 0 {
				fmt.Printf("Successfully created new keypair.\n")
				fmt.Printf("Address:     %s\n", address.Hex())
				fmt.Printf("Private key: %s\n", formatPrivateKey(privateKey))
				return nil
			}

			// If path specified, save to keystore
			keystorePath := args[0]

			// Check if directory exists
			if _, err := os.Stat(keystorePath); os.IsNotExist(err) {
				return fmt.Errorf("directory does not exist: %s", keystorePath)
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

			// Determine final file name
			var finalPath string
			if len(args) >= 2 {
				// Account name provided - rename to account name (like cast)
				accountName := args[1]
				finalPath = filepath.Join(keystorePath, accountName)
				if err := os.Rename(account.URL.Path, finalPath); err != nil {
					return fmt.Errorf("failed to rename keystore file: %w", err)
				}
			} else {
				// No account name - use UUID (like cast)
				oldPath := account.URL.Path
				uuidName := uuid.New().String()
				finalPath = filepath.Join(keystorePath, uuidName)
				if err := os.Rename(oldPath, finalPath); err != nil {
					return fmt.Errorf("failed to rename keystore file: %w", err)
				}
			}

			fmt.Printf("Created new encrypted keystore file: %s\n", finalPath)
			fmt.Printf("Address: %s\n", address.Hex())

			return nil
		},
	}

	cmd.Flags().StringVar(&unsafePassword, "unsafe-password", "", "Password for the JSON keystore in cleartext (unsafe)")

	return cmd
}

// formatPrivateKey formats the private key as hex string
func formatPrivateKey(key *ecdsa.PrivateKey) string {
	return fmt.Sprintf("0x%x", crypto.FromECDSA(key))
}

// promptPassword prompts user for password securely
func promptPassword() (string, error) {
	fmt.Print("Enter password: ")
	password, err := term.ReadPassword(int(syscall.Stdin))
	fmt.Println() // Add newline after password input
	if err != nil {
		return "", err
	}
	return string(password), nil
}
