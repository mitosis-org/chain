package tx

import (
	"crypto/ecdsa"
	"fmt"
	"math/big"
	"os"
	"strings"
	"syscall"

	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	ethcrypto "github.com/ethereum/go-ethereum/crypto"
	"github.com/mitosis-org/chain/cmd/mito/internal/config"
	"golang.org/x/term"
)

// SigningMethod represents the method used for signing transactions
type SigningMethod int

// Signing method types
const (
	SigningMethodPrivateKey SigningMethod = iota
	SigningMethodKeyfile
)

// Signer handles transaction signing with different methods
type Signer struct {
	config *config.ResolvedConfig
}

// SigningConfig holds the signing configuration
type SigningConfig struct {
	Method          SigningMethod
	PrivateKey      *ecdsa.PrivateKey
	KeyfilePath     string
	KeyfilePassword string
	Address         common.Address
}

// NewSigner creates a new transaction signer
func NewSigner(config *config.ResolvedConfig) *Signer {
	return &Signer{
		config: config,
	}
}

// GetSigningConfig determines the signing method and returns the configuration
func (s *Signer) GetSigningConfig() (*SigningConfig, error) {
	config := &SigningConfig{}

	// Check if private key is provided
	if s.config.PrivateKey != "" {
		privKey, err := parsePrivateKey(s.config.PrivateKey)
		if err != nil {
			return nil, fmt.Errorf("failed to parse private key: %w", err)
		}
		config.Method = SigningMethodPrivateKey
		config.PrivateKey = privKey
		config.Address = ethcrypto.PubkeyToAddress(privKey.PublicKey)
		return config, nil
	}

	// Check if keyfile is provided
	if s.config.KeyfilePath != "" {
		password, err := s.getKeyfilePassword()
		if err != nil {
			return nil, fmt.Errorf("failed to get keyfile password: %w", err)
		}

		privKey, err := loadPrivateKeyFromKeyfile(s.config.KeyfilePath, password)
		if err != nil {
			return nil, fmt.Errorf("failed to load private key from keyfile: %w", err)
		}

		config.Method = SigningMethodKeyfile
		config.PrivateKey = privKey
		config.KeyfilePath = s.config.KeyfilePath
		config.KeyfilePassword = password
		config.Address = ethcrypto.PubkeyToAddress(privKey.PublicKey)
		return config, nil
	}

	return nil, fmt.Errorf("no signing method provided: use --private-key or --keyfile")
}

// SignTransaction signs a transaction with the configured signing method
func (s *Signer) SignTransaction(tx *types.Transaction, chainID int64) (*types.Transaction, error) {
	signingConfig, err := s.GetSigningConfig()
	if err != nil {
		return nil, err
	}

	// Create signer for the chain
	signer := types.NewEIP155Signer(big.NewInt(chainID))

	// Sign the transaction
	signedTx, err := types.SignTx(tx, signer, signingConfig.PrivateKey)
	if err != nil {
		return nil, fmt.Errorf("failed to sign transaction: %w", err)
	}

	return signedTx, nil
}

// GetSignerAddress returns the address of the signer
func (s *Signer) GetSignerAddress() (common.Address, error) {
	signingConfig, err := s.GetSigningConfig()
	if err != nil {
		return common.Address{}, err
	}
	return signingConfig.Address, nil
}

// parsePrivateKey converts a hex string to an ECDSA private key
func parsePrivateKey(key string) (*ecdsa.PrivateKey, error) {
	// Remove 0x prefix if present
	key = strings.TrimPrefix(key, "0x")

	privKey, err := ethcrypto.HexToECDSA(key)
	if err != nil {
		return nil, fmt.Errorf("invalid private key format: %w", err)
	}

	return privKey, nil
}

// getKeyfilePassword gets the keyfile password from various sources
func (s *Signer) getKeyfilePassword() (string, error) {
	// Check if password is provided directly
	if s.config.KeyfilePassword != "" {
		return s.config.KeyfilePassword, nil
	}

	// Check if password file is provided
	if s.config.KeyfilePasswordFile != "" {
		passwordBytes, err := os.ReadFile(s.config.KeyfilePasswordFile)
		if err != nil {
			return "", fmt.Errorf("failed to read password file: %w", err)
		}
		return strings.TrimSpace(string(passwordBytes)), nil
	}

	// Prompt for password
	fmt.Print("Enter keyfile password: ")
	passwordBytes, err := term.ReadPassword(syscall.Stdin)
	if err != nil {
		return "", fmt.Errorf("failed to read password: %w", err)
	}
	fmt.Println() // Add newline after password input

	return string(passwordBytes), nil
}

// loadPrivateKeyFromKeyfile loads a private key from a keyfile (supports geth, cast, and other wallet tools)
func loadPrivateKeyFromKeyfile(keyfilePath, password string) (*ecdsa.PrivateKey, error) {
	// Read keyfile
	keyfileData, err := os.ReadFile(keyfilePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read keyfile: %w", err)
	}

	// Directly decrypt the keyfile using go-ethereum's DecryptKey
	// This works for both geth keyfiles (with address field) and cast keyfiles (without address field)
	parsedKey, err := keystore.DecryptKey(keyfileData, password)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt keyfile: %w", err)
	}

	return parsedKey.PrivateKey, nil
}
