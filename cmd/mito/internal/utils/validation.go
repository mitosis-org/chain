package utils

import (
	"encoding/hex"
	"fmt"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	ethcrypto "github.com/ethereum/go-ethereum/crypto"
)

// ParseValue parses a string into a big.Int value
func ParseValue(value string) (*big.Int, error) {
	valueInt := big.NewInt(0)
	if value != "" {
		var success bool
		valueInt, success = big.NewInt(0).SetString(value, 10)
		if !success {
			return nil, fmt.Errorf("invalid format")
		}
	}
	return valueInt, nil
}

// ParseValueAsWei parses a decimal string (e.g. "1.5") into wei units (multiplied by 10^18)
func ParseValueAsWei(value string) (*big.Int, error) {
	if value == "" {
		return big.NewInt(0), nil
	}

	// Create a new Float
	valueFloat := new(big.Float)
	_, success := valueFloat.SetString(value)
	if !success {
		return nil, fmt.Errorf("invalid decimal format: %s", value)
	}

	// Multiply by 10^18
	multiplier := new(big.Float).SetInt(new(big.Int).Exp(big.NewInt(10), big.NewInt(18), nil))
	valueFloat.Mul(valueFloat, multiplier)

	// Convert to big.Int and check accuracy
	valueInt, accuracy := valueFloat.Int(nil)

	// Check if there was loss of precision
	if accuracy != big.Exact {
		return nil, fmt.Errorf("decimal precision exceeds 18 places: %s", value)
	}

	return valueInt, nil
}

// FormatWeiToEther formats Wei to MITO for display
func FormatWeiToEther(wei *big.Int) string {
	if wei == nil {
		return "0"
	}

	// Create a float representation of wei / 10^18
	fValue := new(big.Float).SetInt(wei)
	ethValue := new(big.Float).Quo(fValue, big.NewFloat(1e18))

	// Convert to string with precision
	result := ethValue.Text('f', 18)

	// Remove trailing zeros and decimal point if not needed
	result = strings.TrimRight(result, "0")
	result = strings.TrimRight(result, ".")

	return result
}

// ValidateAddress validates that the provided string is a valid Ethereum address
func ValidateAddress(addr string) (common.Address, error) {
	if !common.IsHexAddress(addr) {
		return common.Address{}, fmt.Errorf("invalid Ethereum address format")
	}
	return common.HexToAddress(addr), nil
}

// IncreaseFee adds the fee amount to the transaction value
func IncreaseFee(value *big.Int, fee *big.Int) *big.Int {
	if value == nil {
		value = big.NewInt(0)
	}
	if fee == nil {
		return value
	}
	return new(big.Int).Add(value, fee)
}

// DecodeHexWithPrefix decodes a hex string with 0x prefix
func DecodeHexWithPrefix(hexStr string) ([]byte, error) {
	// Remove 0x prefix if present
	if len(hexStr) >= 2 && hexStr[0:2] == "0x" {
		hexStr = hexStr[2:]
	}

	// Decode the hex string
	bytes := make([]byte, len(hexStr)/2)
	_, err := hex.Decode(bytes, []byte(hexStr))
	if err != nil {
		return nil, fmt.Errorf("failed to decode hex string: %w", err)
	}

	return bytes, nil
}

// GetAddressFromPrivateKey returns the address corresponding to a private key
func GetAddressFromPrivateKey(privKey string) (string, error) {
	// Remove 0x prefix if present
	privKey = strings.TrimPrefix(privKey, "0x")

	// Parse private key
	privateKey, err := ethcrypto.HexToECDSA(privKey)
	if err != nil {
		return "", fmt.Errorf("invalid private key: %w", err)
	}

	return ethcrypto.PubkeyToAddress(privateKey.PublicKey).Hex(), nil
}

// ParsePercentageToBasisPoints parses a percentage string (with % sign) into basis points (1% = 100 basis points)
func ParsePercentageToBasisPoints(value string) (*big.Int, error) {
	// Require the % sign
	if len(value) == 0 || value[len(value)-1] != '%' {
		return nil, fmt.Errorf("percentage value must end with %% symbol (e.g., \"5%%\")")
	}

	// Remove the % sign
	percentValue := value[:len(value)-1]

	// Parse the percentage value
	percentFloat, ok := new(big.Float).SetString(percentValue)
	if !ok {
		return nil, fmt.Errorf("error parsing percentage value: %s", value)
	}

	// Multiply by 100 to convert percentage to basis points
	basisPoints := new(big.Float).Mul(percentFloat, big.NewFloat(100))

	// Convert to int
	var accuracy big.Accuracy
	basisPointsInt, accuracy := basisPoints.Int(nil)
	if accuracy == big.Exact {
		return basisPointsInt, nil
	}

	return nil, fmt.Errorf("error converting percentage to basis points: %s", value)
}

// FormatBasisPointsToPercent formats basis points as a percentage string with % symbol
func FormatBasisPointsToPercent(basisPoints *big.Int) string {
	if basisPoints == nil {
		return "0.00%"
	}

	// Create a float representation of basisPoints / 100
	fValue := new(big.Float).SetInt(basisPoints)
	percentage := new(big.Float).Quo(fValue, big.NewFloat(100))

	// Convert to string with precision
	return fmt.Sprintf("%.2f%%", percentage)
}
