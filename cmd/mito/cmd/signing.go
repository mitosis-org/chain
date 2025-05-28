package cmd

import (
	"crypto/ecdsa"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"syscall"

	"github.com/ethereum/go-ethereum/accounts/keystore"
	ethcrypto "github.com/ethereum/go-ethereum/crypto"
	"golang.org/x/term"
)

// SigningMethod represents the method used for signing transactions
type SigningMethod int

// Signing method types
const (
	SigningMethodPrivateKey SigningMethod = iota
	SigningMethodKeyfile
)

// SigningConfig holds the signing configuration
type SigningConfig struct {
	Method          SigningMethod
	PrivateKey      *ecdsa.PrivateKey
	KeyfilePath     string
	KeyfilePassword string
}

// GetSigningConfig determines the signing method and returns the configuration
func GetSigningConfig() (*SigningConfig, error) {
	config := &SigningConfig{}

	// Check if private key is provided
	if privateKey != "" {
		privKey, err := parsePrivateKey(privateKey)
		if err != nil {
			return nil, fmt.Errorf("failed to parse private key: %w", err)
		}
		config.Method = SigningMethodPrivateKey
		config.PrivateKey = privKey
		return config, nil
	}

	// Check if keyfile is provided
	if keyfilePath != "" {
		password, err := getKeyfilePassword()
		if err != nil {
			return nil, fmt.Errorf("failed to get keyfile password: %w", err)
		}

		privKey, err := loadPrivateKeyFromKeyfile(keyfilePath, password)
		if err != nil {
			return nil, fmt.Errorf("failed to load private key from keyfile: %w", err)
		}

		config.Method = SigningMethodKeyfile
		config.PrivateKey = privKey
		config.KeyfilePath = keyfilePath
		config.KeyfilePassword = password
		return config, nil
	}

	return nil, fmt.Errorf("no signing method provided: use --private-key or --keyfile")
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
func getKeyfilePassword() (string, error) {
	// Check if password is provided directly
	if keyfilePassword != "" {
		return keyfilePassword, nil
	}

	// Check if password file is provided
	if keyfilePasswordFile != "" {
		passwordBytes, err := ioutil.ReadFile(keyfilePasswordFile)
		if err != nil {
			return "", fmt.Errorf("failed to read password file: %w", err)
		}
		return strings.TrimSpace(string(passwordBytes)), nil
	}

	// Prompt for password
	fmt.Print("Enter keyfile password: ")
	passwordBytes, err := term.ReadPassword(int(syscall.Stdin))
	if err != nil {
		return "", fmt.Errorf("failed to read password: %w", err)
	}
	fmt.Println() // Add newline after password input

	return string(passwordBytes), nil
}

// loadPrivateKeyFromKeyfile loads a private key from a geth keyfile
func loadPrivateKeyFromKeyfile(keyfilePath, password string) (*ecdsa.PrivateKey, error) {
	// Read keyfile
	keyfileData, err := ioutil.ReadFile(keyfilePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read keyfile: %w", err)
	}

	// Create a temporary keyfile directory
	tempDir, err := ioutil.TempDir("", "temp_keyfile")
	if err != nil {
		return nil, fmt.Errorf("failed to create temp directory: %w", err)
	}
	defer os.RemoveAll(tempDir)

	// Copy keyfile to temp directory
	tempKeyfilePath := filepath.Join(tempDir, filepath.Base(keyfilePath))
	err = ioutil.WriteFile(tempKeyfilePath, keyfileData, 0600)
	if err != nil {
		return nil, fmt.Errorf("failed to write temp keyfile: %w", err)
	}

	// Create keyfile instance
	ks := keystore.NewKeyStore(tempDir, keystore.StandardScryptN, keystore.StandardScryptP)

	// Get accounts
	accounts := ks.Accounts()
	if len(accounts) == 0 {
		return nil, fmt.Errorf("no accounts found in keyfile")
	}

	// Use the first account
	account := accounts[0]

	// Unlock the account
	err = ks.Unlock(account, password)
	if err != nil {
		return nil, fmt.Errorf("failed to unlock keyfile: %w", err)
	}

	// Export the private key
	key, err := ks.Export(account, password, password)
	if err != nil {
		return nil, fmt.Errorf("failed to export private key: %w", err)
	}

	// Parse the exported key
	parsedKey, err := keystore.DecryptKey(key, password)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt key: %w", err)
	}

	return parsedKey.PrivateKey, nil
}
