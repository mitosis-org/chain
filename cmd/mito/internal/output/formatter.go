package output

import (
	"encoding/json"
	"fmt"
	"math/big"
	"os"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/mitosis-org/chain/cmd/mito/internal/units"
)

// TransactionFormatter handles transaction output formatting
type TransactionFormatter struct {
	outputFile string
}

// NewTransactionFormatter creates a new transaction formatter
func NewTransactionFormatter(outputFile string) *TransactionFormatter {
	return &TransactionFormatter{
		outputFile: outputFile,
	}
}

// ValidatorCreateInfo contains information for validator creation transaction display
type ValidatorCreateInfo struct {
	PubKey            string
	Operator          string
	RewardManager     string
	CommissionRate    string
	Metadata          string
	InitialCollateral string
}

// ValidatorUpdateInfo contains information for validator update transaction display
type ValidatorUpdateInfo struct {
	ValidatorAddress string
	FieldName        string
	NewValue         string
}

// CollateralInfo contains information for collateral transaction display
type CollateralInfo struct {
	ValidatorAddress string
	Amount           string
	Operation        string // "deposit", "withdraw", etc.
}

// CollateralDepositInfo contains information for collateral deposit transaction display
type CollateralDepositInfo struct {
	ValidatorAddress string
	CollateralAmount string   // Keep as string since it's MITO amount
	Fee              *big.Int // Wei amount
	TotalValue       *big.Int // Wei amount
}

// CollateralWithdrawInfo contains information for collateral withdraw transaction display
type CollateralWithdrawInfo struct {
	ValidatorAddress string
	ReceiverAddress  string
	CollateralAmount string   // Keep as string since it's MITO amount
	Fee              *big.Int // Wei amount
}

// CollateralPermissionInfo contains information for collateral permission transaction display
type CollateralPermissionInfo struct {
	ValidatorAddress string
	CollateralOwner  string
	IsPermitted      bool
}

// CollateralOwnershipInfo contains information for collateral ownership transfer transaction display
type CollateralOwnershipInfo struct {
	ValidatorAddress string
	NewOwner         string
	Fee              *big.Int // Wei amount
}

// FormatValidatorCreateTransaction formats and outputs a validator creation transaction
func (tf *TransactionFormatter) FormatValidatorCreateTransaction(tx *types.Transaction, info *ValidatorCreateInfo) error {
	// Print header information
	fmt.Println("===== Create Validator Transaction =====")
	fmt.Printf("Public Key                 : %s\n", info.PubKey)
	fmt.Printf("Operator                   : %s\n", info.Operator)
	fmt.Printf("Reward Manager             : %s\n", info.RewardManager)
	fmt.Printf("Commission Rate            : %s\n", info.CommissionRate)
	fmt.Printf("Metadata                   : %s\n", info.Metadata)
	fmt.Printf("Initial Collateral         : %s MITO\n", info.InitialCollateral)
	fmt.Println()

	return nil
}

// FormatValidatorUpdateTransaction formats and outputs a validator update transaction
func (tf *TransactionFormatter) FormatValidatorUpdateTransaction(tx *types.Transaction, info *ValidatorUpdateInfo) error {
	// Print header information
	fmt.Printf("===== Update %s Transaction =====\n", info.FieldName)
	fmt.Printf("Validator Address        : %s\n", info.ValidatorAddress)
	fmt.Printf("New %s             : %s\n", info.FieldName, info.NewValue)
	fmt.Println()

	return nil
}

// FormatCollateralTransaction formats and outputs a collateral transaction
func (tf *TransactionFormatter) FormatCollateralTransaction(tx *types.Transaction, info *CollateralInfo) error {
	// Print header information
	operationTitle := fmt.Sprintf("%s%s",
		string(rune(info.Operation[0]-32)), info.Operation[1:]) // Capitalize first letter

	fmt.Printf("===== %s Collateral Transaction =====\n", operationTitle)
	fmt.Printf("Validator Address        : %s\n", info.ValidatorAddress)
	fmt.Printf("Amount                   : %s MITO\n", info.Amount)
	fmt.Println()

	return nil
}

// FormatCollateralDepositTransaction formats and outputs a collateral deposit transaction
func (tf *TransactionFormatter) FormatCollateralDepositTransaction(tx *types.Transaction, info *CollateralDepositInfo) error {
	// Print header information
	fmt.Println("===== Deposit Collateral Transaction =====")
	fmt.Printf("Validator Address          : %s\n", info.ValidatorAddress)
	fmt.Printf("Collateral Amount          : %s MITO\n", info.CollateralAmount)
	fmt.Printf("Fee                        : %s\n", units.FormatWeiToBothUnits(info.Fee))
	fmt.Printf("Total Value                : %s\n", units.FormatWeiToBothUnits(info.TotalValue))
	fmt.Println()
	fmt.Println("🚨 IMPORTANT: Only permitted collateral owners can deposit collateral for a validator.")
	fmt.Println("If your address is not a permitted collateral owner for this validator, the transaction will fail.")
	fmt.Println()

	return nil
}

// FormatCollateralWithdrawTransaction formats and outputs a collateral withdraw transaction
func (tf *TransactionFormatter) FormatCollateralWithdrawTransaction(tx *types.Transaction, info *CollateralWithdrawInfo) error {
	// Print header information
	fmt.Println("===== Withdraw Collateral Transaction =====")
	fmt.Printf("Validator Address          : %s\n", info.ValidatorAddress)
	fmt.Printf("Receiver Address           : %s\n", info.ReceiverAddress)
	fmt.Printf("Collateral Amount          : %s MITO\n", info.CollateralAmount)
	fmt.Printf("Fee                        : %s\n", units.FormatWeiToBothUnits(info.Fee))
	fmt.Println()
	fmt.Println("🚨 IMPORTANT: Only the collateral owner can withdraw collateral.")
	fmt.Println()

	return nil
}

// FormatCollateralPermissionTransaction formats and outputs a collateral permission transaction
func (tf *TransactionFormatter) FormatCollateralPermissionTransaction(tx *types.Transaction, info *CollateralPermissionInfo) error {
	// Print header information
	permissionText := "DENY"
	if info.IsPermitted {
		permissionText = "PERMIT"
	}

	fmt.Println("===== Set Permitted Collateral Owner Transaction =====")
	fmt.Printf("Validator Address          : %s\n", info.ValidatorAddress)
	fmt.Printf("Collateral Owner           : %s\n", info.CollateralOwner)
	fmt.Printf("Permission                 : %s\n", permissionText)
	fmt.Println()

	return nil
}

// FormatCollateralOwnershipTransaction formats and outputs a collateral ownership transfer transaction
func (tf *TransactionFormatter) FormatCollateralOwnershipTransaction(tx *types.Transaction, info *CollateralOwnershipInfo) error {
	// Print header information
	fmt.Println("===== Transfer Collateral Ownership Transaction =====")
	fmt.Printf("Validator Address          : %s\n", info.ValidatorAddress)
	fmt.Printf("New Owner                  : %s\n", info.NewOwner)
	fmt.Printf("Fee                        : %s\n", units.FormatWeiToBothUnits(info.Fee))
	fmt.Println()

	return nil
}

// OutputTransaction outputs the transaction in the appropriate format
func (tf *TransactionFormatter) OutputTransaction(tx *types.Transaction) error {
	// Create a formatted transaction object for JSON output
	formattedTx := map[string]interface{}{
		"type":                 fmt.Sprintf("0x%x", tx.Type()),
		"chainId":              fmt.Sprintf("0x%x", tx.ChainId()),
		"nonce":                fmt.Sprintf("0x%x", tx.Nonce()),
		"to":                   tx.To().Hex(),
		"gas":                  fmt.Sprintf("0x%x", tx.Gas()),
		"gasPrice":             fmt.Sprintf("0x%x", tx.GasPrice()),
		"maxPriorityFeePerGas": nil,
		"maxFeePerGas":         nil,
		"value":                fmt.Sprintf("0x%x", tx.Value()),
		"input":                fmt.Sprintf("0x%x", tx.Data()),
	}

	// Add signature fields if transaction is signed
	v, r, s := tx.RawSignatureValues()
	if v != nil && v.Cmp(big.NewInt(0)) != 0 {
		formattedTx["v"] = fmt.Sprintf("0x%x", v)
		formattedTx["r"] = fmt.Sprintf("0x%x", r)
		formattedTx["s"] = fmt.Sprintf("0x%x", s)
	}

	// Convert to JSON
	txJSON, err := json.Marshal(formattedTx)
	if err != nil {
		return fmt.Errorf("failed to convert transaction to JSON: %w", err)
	}

	// Output to file or stdout
	if tf.outputFile != "" {
		return os.WriteFile(tf.outputFile, txJSON, 0o600)
	}

	fmt.Println(string(txJSON))
	return nil
}

// FormatTransactionSummary prints a summary of transaction details
func FormatTransactionSummary(tx *types.Transaction) {
	fmt.Println("Transaction Summary:")
	fmt.Printf("  Chain ID: %s\n", tx.ChainId().String())
	fmt.Printf("  Nonce: %d\n", tx.Nonce())
	fmt.Printf("  To: %s\n", tx.To().Hex())
	fmt.Printf("  Value: %s\n", units.FormatWeiToBothUnits(tx.Value()))
	fmt.Printf("  Gas Limit: %d\n", tx.Gas())
	fmt.Printf("  Gas Price: %s\n", units.FormatWeiToBothUnits(tx.GasPrice()))
	fmt.Printf("  Data Size: %d bytes\n", len(tx.Data()))

	// Check if signed
	v, _, _ := tx.RawSignatureValues()
	if v != nil && v.Cmp(big.NewInt(0)) != 0 {
		fmt.Println("  Status: Signed")
	} else {
		fmt.Println("  Status: Unsigned")
	}
	fmt.Println()
}
