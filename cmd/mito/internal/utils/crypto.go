package utils

import (
	"crypto/ecdsa"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/mitosis-org/chain/cmd/mito/internal/types"
)

// LoadPrivateKeyFromPrivValidatorKey loads a private key from cosmos priv_validator_key.json file
func LoadPrivateKeyFromPrivValidatorKey(keyfilePath string) (*ecdsa.PrivateKey, error) {
	// Read the priv_validator_key.json file
	keyfileData, err := os.ReadFile(keyfilePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read priv_validator_key.json: %w", err)
	}

	// Parse the JSON structure
	var privValidatorKey types.PrivValidatorKey
	if err := json.Unmarshal(keyfileData, &privValidatorKey); err != nil {
		return nil, fmt.Errorf("failed to parse priv_validator_key.json: %w", err)
	}

	// Validate the private key type
	if privValidatorKey.PrivKey.Type != "tendermint/PrivKeySecp256k1" {
		return nil, fmt.Errorf("unsupported private key type: %s", privValidatorKey.PrivKey.Type)
	}

	// Decode the base64-encoded private key
	privKeyBytes, err := base64.StdEncoding.DecodeString(privValidatorKey.PrivKey.Value)
	if err != nil {
		return nil, fmt.Errorf("failed to decode base64 private key: %w", err)
	}

	// Convert the raw bytes to ECDSA private key
	privKey, err := crypto.ToECDSA(privKeyBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to convert to ECDSA private key: %w", err)
	}

	return privKey, nil
}
