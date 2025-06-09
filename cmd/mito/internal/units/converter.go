package units

import (
	"fmt"
	"math/big"
	"strings"
)

const (
	// 1 Mito = 10^18 wei = 10^9 gwei
	// So 1 Mito = 10^9 gwei
	GweiPerMito = 1e9
	WeiPerGwei  = 1e9
	WeiPerMito  = 1e18
)

// ParseGweiToWei converts gwei string to wei (*big.Int)
func ParseGweiToWei(gweiStr string) (*big.Int, error) {
	gweiFloat := new(big.Float)
	_, success := gweiFloat.SetString(gweiStr)
	if !success {
		return nil, fmt.Errorf("invalid gwei format: %s", gweiStr)
	}

	// Multiply by 10^9 to convert gwei to wei
	weiFloat := new(big.Float).Mul(gweiFloat, big.NewFloat(WeiPerGwei))

	// Convert to big.Int
	weiInt, accuracy := weiFloat.Int(nil)
	if accuracy != big.Exact {
		return nil, fmt.Errorf("gwei precision exceeds 9 decimal places: %s", gweiStr)
	}

	return weiInt, nil
}

// ParseMitoToWei converts Mito string to wei (*big.Int)
func ParseMitoToWei(mitoStr string) (*big.Int, error) {
	mitoFloat := new(big.Float)
	_, success := mitoFloat.SetString(mitoStr)
	if !success {
		return nil, fmt.Errorf("invalid Mito format: %s", mitoStr)
	}

	// Multiply by 10^18 to convert Mito to wei
	weiFloat := new(big.Float).Mul(mitoFloat, big.NewFloat(WeiPerMito))

	// Convert to big.Int
	weiInt, accuracy := weiFloat.Int(nil)
	if accuracy != big.Exact {
		return nil, fmt.Errorf("Mito precision exceeds 18 decimal places: %s", mitoStr)
	}

	return weiInt, nil
}

// ParseWeiToWei converts wei string to wei (*big.Int)
func ParseWeiToWei(weiStr string) (*big.Int, error) {
	weiInt := new(big.Int)
	_, success := weiInt.SetString(weiStr, 10)
	if !success {
		return nil, fmt.Errorf("invalid wei format: %s", weiStr)
	}
	return weiInt, nil
}

// FormatWeiToMito formats wei to Mito string for display
func FormatWeiToMito(wei *big.Int) string {
	if wei == nil {
		return "0"
	}

	weiFloat := new(big.Float).SetInt(wei)
	mitoFloat := new(big.Float).Quo(weiFloat, big.NewFloat(WeiPerMito))

	return mitoFloat.Text('f', 18)
}

// FormatWeiToGwei formats wei to gwei string for display
func FormatWeiToGwei(wei *big.Int) string {
	if wei == nil {
		return "0"
	}

	weiFloat := new(big.Float).SetInt(wei)
	gweiFloat := new(big.Float).Quo(weiFloat, big.NewFloat(WeiPerGwei))

	return gweiFloat.Text('f', 9)
}

// FormatWeiToBothUnits formats wei to both MITO and Gwei for display
// Example: "0.000000001000000007 MITO (1.000000007 Gwei)"
func FormatWeiToBothUnits(wei *big.Int) string {
	if wei == nil {
		return "0.000000000000000000 MITO (0.000000000 Gwei)"
	}

	mitoStr := FormatWeiToMito(wei)
	gweiStr := FormatWeiToGwei(wei)

	// Remove trailing zeros for cleaner display
	mitoStr = strings.TrimRight(strings.TrimRight(mitoStr, "0"), ".")
	gweiStr = strings.TrimRight(strings.TrimRight(gweiStr, "0"), ".")

	if mitoStr == "" {
		mitoStr = "0"
	}
	if gweiStr == "" {
		gweiStr = "0"
	}

	return fmt.Sprintf("%s MITO (%s Gwei)", mitoStr, gweiStr)
}

// ParseValueInput parses value input with unit support
// Defaults to gwei if no unit is specified
// Supports: "20" (gwei), "20gwei", "20wei", "0.00000002mito"
func ParseValueInput(input string) (*big.Int, error) {
	if input == "" {
		return big.NewInt(0), nil
	}

	input = strings.ToLower(strings.TrimSpace(input))

	// Check for explicit unit suffix
	if strings.HasSuffix(input, "mito") {
		mitoStr := strings.TrimSuffix(input, "mito")
		return ParseMitoToWei(mitoStr)
	}

	if strings.HasSuffix(input, "gwei") {
		gweiStr := strings.TrimSuffix(input, "gwei")
		return ParseGweiToWei(gweiStr)
	}

	if strings.HasSuffix(input, "wei") {
		weiStr := strings.TrimSuffix(input, "wei")
		return ParseWeiToWei(weiStr)
	}

	// Default to gwei if no unit specified
	return ParseGweiToWei(input)
}

// ParseGasPriceInput parses gas price input (defaults to gwei)
// Examples: "20" (gwei), "20gwei", "20wei", "0.00000002mito"
func ParseGasPriceInput(input string) (*big.Int, error) {
	if input == "" {
		return nil, fmt.Errorf("gas price cannot be empty")
	}
	return ParseValueInput(input)
}

// ParseContractFeeInput parses contract fee input (defaults to gwei)
// Examples: "100" (gwei), "100gwei", "100wei", "0.0000001mito"
func ParseContractFeeInput(input string) (*big.Int, error) {
	return ParseValueInput(input)
}
