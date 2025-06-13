package wallet

import (
	"crypto/ecdsa"
	"crypto/rand"
	"fmt"

	"github.com/cosmos/go-bip39"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/spf13/cobra"
)

// NewNewMnemonicCmd creates the wallet new-mnemonic command
func NewNewMnemonicCmd() *cobra.Command {
	var (
		words    int
		accounts int
	)

	cmd := &cobra.Command{
		Use:   "new-mnemonic",
		Short: "Creates a new mnemonic with a set number of words",
		Long:  `Generates a random BIP39 mnemonic phrase.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			// Validate word count and get entropy size (in bytes)
			var entropySize int
			switch words {
			case 12:
				entropySize = 16 // 128 bits = 16 bytes
			case 15:
				entropySize = 20 // 160 bits = 20 bytes
			case 18:
				entropySize = 24 // 192 bits = 24 bytes
			case 21:
				entropySize = 28 // 224 bits = 28 bytes
			case 24:
				entropySize = 32 // 256 bits = 32 bytes
			default:
				return fmt.Errorf("invalid word count: %d (must be 12, 15, 18, 21, or 24)", words)
			}

			// Generate entropy using crypto/rand
			entropy := make([]byte, entropySize)
			_, err := rand.Read(entropy)
			if err != nil {
				return fmt.Errorf("failed to generate entropy: %w", err)
			}

			// Generate mnemonic from entropy
			mnemonic, err := bip39.NewMnemonic(entropy)
			if err != nil {
				return fmt.Errorf("failed to generate mnemonic: %w", err)
			}

			fmt.Printf("Successfully generated a new mnemonic.\n")
			fmt.Printf("Phrase:\n%s\n\n", mnemonic)

			// Generate and display accounts
			fmt.Printf("Accounts:\n")
			for i := 0; i < accounts; i++ {
				// For simplicity, derive keys from mnemonic seed using a deterministic method
				privateKey, err := deriveKeyFromMnemonic(mnemonic, i)
				if err != nil {
					return fmt.Errorf("failed to derive private key for account %d: %w", i, err)
				}

				address := crypto.PubkeyToAddress(privateKey.PublicKey)

				fmt.Printf("- Account %d:\n", i)
				fmt.Printf("Address:     %s\n", address.Hex())
				fmt.Printf("Private key: %s\n", formatPrivateKey(privateKey))
				if i < accounts-1 {
					fmt.Println()
				}
			}

			return nil
		},
	}

	cmd.Flags().IntVarP(&words, "words", "w", 12, "The amount of words for the mnemonic")
	cmd.Flags().IntVarP(&accounts, "accounts", "a", 1, "The number of accounts to display")

	return cmd
}

// deriveKeyFromMnemonic derives a private key from mnemonic using a simple deterministic method
func deriveKeyFromMnemonic(mnemonic string, index int) (*ecdsa.PrivateKey, error) {
	// Generate seed from mnemonic
	seed := bip39.NewSeed(mnemonic, "")

	// For simplicity, we'll use a hash-based key derivation
	// This is not a full BIP32/BIP44 implementation but serves our purpose
	derivationInput := fmt.Sprintf("%x%d", seed, index)
	keyBytes := crypto.Keccak256([]byte(derivationInput))

	// Create ECDSA private key from hash
	privateKey, err := crypto.ToECDSA(keyBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to create ECDSA key: %w", err)
	}

	return privateKey, nil
}
