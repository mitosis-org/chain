package wallet

import (
	"bufio"
	"crypto/ecdsa"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/cosmos/go-bip39"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/mitosis-org/chain/cmd/mito/internal/utils"
	"github.com/spf13/cobra"
)

// NewImportCmd creates the wallet import command
func NewImportCmd() *cobra.Command {
	var (
		keystoreDir      string
		interactive      bool
		unsafePassword   string
		mnemonicIndex    int
		mnemonic         string
		privateKey       string
		privValidatorKey string
	)

	cmd := &cobra.Command{
		Use:   "import [account_name]",
		Short: "Import a private key into an encrypted keystore",
		Long: `Import a private key into an encrypted keystore.

If no keystore-dir is specified, it will be saved in the default ~/.mito/keystores,
so it can be accessed through the --account option in methods like forge script,
cast send or any other that requires a private key.`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			accountName := args[0]

			// Set default keystore directory
			if keystoreDir == "" {
				homeDir, err := os.UserHomeDir()
				if err != nil {
					return fmt.Errorf("failed to get home directory: %w", err)
				}
				keystoreDir = filepath.Join(homeDir, ".mito", "keystores")
			}

			// Create keystore directory if it doesn't exist
			if err := os.MkdirAll(keystoreDir, 0o755); err != nil {
				return fmt.Errorf("failed to create keystore directory: %w", err)
			}

			var key *ecdsa.PrivateKey
			var err error

			// Handle different input methods with proper precedence (like cast)
			// Priority: interactive > private-key > priv-validator-key > mnemonic
			switch {
			case interactive:
				key, err = getPrivateKeyInteractive()
				if err != nil {
					return fmt.Errorf("failed to get private key interactively: %w", err)
				}

			case privateKey != "":
				key, err = crypto.HexToECDSA(strings.TrimPrefix(privateKey, "0x"))
				if err != nil {
					return fmt.Errorf("failed to parse private key: %w", err)
				}

			case privValidatorKey != "":
				key, err = utils.LoadPrivateKeyFromPrivValidatorKey(privValidatorKey)
				if err != nil {
					return fmt.Errorf("failed to load private key from priv_validator_key.json: %w", err)
				}

			case mnemonic != "":
				// Handle mnemonic import
				if !bip39.IsMnemonicValid(mnemonic) {
					return fmt.Errorf("invalid mnemonic phrase")
				}

				key, err = deriveKeyFromMnemonic(mnemonic, mnemonicIndex)
				if err != nil {
					return fmt.Errorf("failed to derive key from mnemonic: %w", err)
				}

			default:
				return fmt.Errorf("no private key source specified (use --interactive, --private-key, --priv-validator-key, or --mnemonic)")
			}

			// Get password for keystore encryption
			var password string
			if unsafePassword != "" {
				password = unsafePassword
			} else {
				password, err = promptPassword()
				if err != nil {
					return fmt.Errorf("failed to get password: %w", err)
				}
			}

			// Create keystore and import key
			ks := keystore.NewKeyStore(keystoreDir, keystore.StandardScryptN, keystore.StandardScryptP)
			account, err := ks.ImportECDSA(key, password)
			if err != nil {
				return fmt.Errorf("failed to import key to keystore: %w", err)
			}

			// Rename the keystore file to just the account name (like cast)
			oldPath := account.URL.Path
			newPath := filepath.Join(keystoreDir, accountName)
			if err := os.Rename(oldPath, newPath); err != nil {
				return fmt.Errorf("failed to rename keystore file: %w", err)
			}

			fmt.Printf("`%s` keystore was saved successfully. Address: %s\n", accountName, account.Address.Hex())

			return nil
		},
	}

	// Flags (matching cast exactly)
	cmd.Flags().StringVarP(&keystoreDir, "keystore-dir", "k", "", "If provided, keystore will be saved here instead of the default keystores directory (~/.mito/keystores)")
	cmd.Flags().StringVar(&unsafePassword, "unsafe-password", "", "Password for the JSON keystore in cleartext This is unsafe, we recommend using the default hidden password prompt")

	// Wallet options - raw (matching cast)
	cmd.Flags().BoolVarP(&interactive, "interactive", "i", false, "Open an interactive prompt to enter your private key")
	cmd.Flags().StringVar(&privateKey, "private-key", "", "Use the provided private key")
	cmd.Flags().StringVar(&privValidatorKey, "priv-validator-key", "", "Use the private key from cosmos priv_validator_key.json file")
	cmd.Flags().StringVar(&mnemonic, "mnemonic", "", "Use the mnemonic phrase of mnemonic file at the specified path")
	cmd.Flags().IntVar(&mnemonicIndex, "mnemonic-index", 0, "Use the private key from the given mnemonic index. Used with --mnemonic.")

	return cmd
}

// getPrivateKeyInteractive prompts user to enter private key interactively
func getPrivateKeyInteractive() (*ecdsa.PrivateKey, error) {
	fmt.Print("Enter private key: ")

	// Read from stdin
	reader := bufio.NewReader(os.Stdin)
	privateKeyStr, err := reader.ReadString('\n')
	if err != nil {
		return nil, fmt.Errorf("failed to read private key: %w", err)
	}

	privateKeyStr = strings.TrimSpace(privateKeyStr)
	privateKeyStr = strings.TrimPrefix(privateKeyStr, "0x")

	privateKey, err := crypto.HexToECDSA(privateKeyStr)
	if err != nil {
		return nil, fmt.Errorf("failed to parse private key: %w", err)
	}

	return privateKey, nil
}
