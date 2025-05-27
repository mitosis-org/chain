package cmd

import (
	"fmt"
	"math/big"
	"os"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/mitosis-org/chain/bindings"
	"github.com/spf13/cobra"
)

// newTxCreateCollateralDepositCmd creates the collateral deposit command for tx create
func newTxCreateCollateralDepositCmd() *cobra.Command {
	var (
		validator string
		amount    string
	)

	cmd := &cobra.Command{
		Use:   "deposit",
		Short: "Create a collateral deposit transaction",
		Long: `Create a transaction to deposit collateral for a validator.
This command can create either signed or unsigned transactions based on the flags provided.
Only permitted collateral owners can deposit collateral for a validator.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runDepositCollateralTx(cmd, validator, amount)
		},
	}

	// Add common flags for contract interaction
	cmd.Flags().StringVar(&validatorManagerContractAddr, "contract", "", "ValidatorManager contract address")

	// Add transaction creation flags
	cmd.Flags().StringVar(&chainID, "chain-id", "1337", "Chain ID for the transaction")
	cmd.Flags().Uint64Var(&gasLimit, "gas-limit", 500000, "Gas limit for the transaction")
	cmd.Flags().StringVar(&gasPrice, "gas-price", "20000000000", "Gas price in wei")
	cmd.Flags().Uint64Var(&txNonce, "nonce", 0, "Nonce for the transaction (required for signed transactions)")
	cmd.Flags().StringVar(&contractFee, "fee", "0.001", "Contract fee in MITO")

	// Add tx create specific flags
	AddTxCreateFlags(cmd)

	// Command-specific flags
	cmd.Flags().StringVar(&validator, "validator", "", "Validator address")
	cmd.Flags().StringVar(&amount, "amount", "", "Amount of collateral to deposit in MITO (e.g., \"1.5\")")

	// Mark required flags (contract is not required as it can come from config)
	mustMarkFlagRequired(cmd, "validator")
	mustMarkFlagRequired(cmd, "amount")

	// Add PreRun to load config and validate mutual exclusive flags
	existingPreRun := cmd.PreRun
	cmd.PreRun = func(cmd *cobra.Command, args []string) {
		// Call the PreRun from AddTxCreateFlags first (for mutual exclusive validation)
		if existingPreRun != nil {
			existingPreRun(cmd, args)
		}

		// Additional validation specific to this command
		if validatorManagerContractAddr == "" {
			fmt.Println("Error: ValidatorManager contract address is required. Set it with --contract flag or use 'mito config set-contract <address>'")
			os.Exit(1)
		}
	}

	return cmd
}

func runDepositCollateralTx(cmd *cobra.Command, validator, amount string) error {
	// Parse collateral amount as decimal MITO and convert to wei
	collateralAmount, err := ParseValueAsWei(amount)
	if err != nil {
		return fmt.Errorf("failed to parse amount: %w", err)
	}

	// Ensure collateral amount is greater than 0
	if collateralAmount.Cmp(big.NewInt(0)) <= 0 {
		return fmt.Errorf("collateral amount must be greater than 0")
	}

	// Parse fee amount
	feeAmount, err := ParseValueAsWei(contractFee)
	if err != nil {
		return fmt.Errorf("failed to parse fee: %w", err)
	}

	// Calculate total value to send (collateral amount + fee)
	totalValue := new(big.Int).Add(collateralAmount, feeAmount)

	// Validate validator address
	valAddr, err := ValidateAddress(validator)
	if err != nil {
		return fmt.Errorf("invalid validator address: %w", err)
	}

	// Show summary
	fmt.Println("===== Deposit Collateral Transaction =====")
	fmt.Printf("Validator Address        : %s\n", valAddr.Hex())
	fmt.Printf("Collateral Amount        : %s MITO\n", FormatWeiToEther(collateralAmount))
	fmt.Printf("Fee                      : %s MITO\n", FormatWeiToEther(feeAmount))
	fmt.Printf("Total Value              : %s MITO\n", FormatWeiToEther(totalValue))
	fmt.Printf("Chain ID                 : %s\n", chainID)
	fmt.Printf("Gas Limit                : %d\n", gasLimit)
	fmt.Printf("Gas Price                : %s wei\n", gasPrice)
	if signed {
		fmt.Printf("Nonce                    : %d\n", txNonce)
	}
	fmt.Println()

	fmt.Println("🚨 IMPORTANT: Only permitted collateral owners can deposit collateral for a validator.")
	fmt.Println("If your address is not a permitted collateral owner for this validator, the transaction will fail.")
	fmt.Println()

	// Check if we should create signed or unsigned transaction
	if signed {
		return createSignedDepositCollateralTransaction(valAddr, totalValue)
	} else {
		return createUnsignedDepositCollateralTransaction(valAddr, totalValue)
	}
}

func createSignedDepositCollateralTransaction(valAddr common.Address, totalValue *big.Int) error {
	if !signed {
		return fmt.Errorf("signed transaction creation requires --signed flag")
	}

	// Get signing configuration
	signingConfig, err := GetSigningConfig()
	if err != nil {
		return fmt.Errorf("failed to get signing configuration: %w", err)
	}
	privateKey := signingConfig.PrivateKey

	// Get contract address
	contractAddr := common.HexToAddress(validatorManagerContractAddr)

	// Encode the function call data
	abi, err := bindings.IValidatorManagerMetaData.GetAbi()
	if err != nil {
		return fmt.Errorf("failed to get contract ABI: %w", err)
	}

	data, err := abi.Pack("depositCollateral", valAddr)
	if err != nil {
		return fmt.Errorf("failed to pack function call: %w", err)
	}

	// Parse gas price
	gasPriceBig, ok := new(big.Int).SetString(gasPrice, 10)
	if !ok {
		return fmt.Errorf("invalid gas price: %s", gasPrice)
	}

	// Parse chain ID
	chainIDBig, ok := new(big.Int).SetString(chainID, 10)
	if !ok {
		return fmt.Errorf("invalid chain ID: %s", chainID)
	}

	// Create and sign transaction
	tx := types.NewTransaction(txNonce, contractAddr, totalValue, gasLimit, gasPriceBig, data)

	// Sign the transaction
	signer := types.NewEIP155Signer(chainIDBig)
	signedTx, err := types.SignTx(tx, signer, privateKey)
	if err != nil {
		return fmt.Errorf("failed to sign transaction: %w", err)
	}

	// Create signed transaction data
	signedTxData := SignedTransactionData{
		TransactionData: TransactionData{
			To:       signedTx.To().Hex(),
			Value:    signedTx.Value().String(),
			Data:     fmt.Sprintf("0x%x", signedTx.Data()),
			GasLimit: signedTx.Gas(),
			GasPrice: signedTx.GasPrice().String(),
			Nonce:    signedTx.Nonce(),
			ChainID:  chainID,
		},
		Hash: signedTx.Hash().Hex(),
	}

	// Get signature components
	v, r, s := signedTx.RawSignatureValues()
	signedTxData.Signature.V = v.String()
	signedTxData.Signature.R = r.String()
	signedTxData.Signature.S = s.String()

	// Get raw transaction data
	rawTxData, err := signedTx.MarshalBinary()
	if err == nil {
		signedTxData.Raw = fmt.Sprintf("0x%x", rawTxData)
	}

	// Output the transaction
	return outputTransaction(signedTxData, "Signed")
}

func createUnsignedDepositCollateralTransaction(valAddr common.Address, totalValue *big.Int) error {
	// Get contract address
	contractAddr := common.HexToAddress(validatorManagerContractAddr)

	// Encode the function call data
	abi, err := bindings.IValidatorManagerMetaData.GetAbi()
	if err != nil {
		return fmt.Errorf("failed to get contract ABI: %w", err)
	}

	data, err := abi.Pack("depositCollateral", valAddr)
	if err != nil {
		return fmt.Errorf("failed to pack function call: %w", err)
	}

	// Create unsigned transaction data
	unsignedTxData := TransactionData{
		To:       contractAddr.Hex(),
		Value:    totalValue.String(),
		Data:     fmt.Sprintf("0x%x", data),
		ChainID:  chainID,
		GasLimit: gasLimit,
		GasPrice: gasPrice,
	}

	// Output the transaction
	return outputTransaction(unsignedTxData, "Unsigned")
}

// newTxCreateCollateralWithdrawCmd creates the collateral withdraw command for tx create
func newTxCreateCollateralWithdrawCmd() *cobra.Command {
	var (
		validator string
		amount    string
		receiver  string
	)

	cmd := &cobra.Command{
		Use:   "withdraw",
		Short: "Create a collateral withdraw transaction",
		Long: `Create a transaction to withdraw collateral from a validator.
This command can create either signed or unsigned transactions based on the flags provided.
Only the collateral owner can withdraw collateral.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runWithdrawCollateralTx(cmd, validator, amount, receiver)
		},
	}

	// Add common flags for contract interaction
	cmd.Flags().StringVar(&validatorManagerContractAddr, "contract", "", "ValidatorManager contract address")

	// Add transaction creation flags
	cmd.Flags().StringVar(&chainID, "chain-id", "1337", "Chain ID for the transaction")
	cmd.Flags().Uint64Var(&gasLimit, "gas-limit", 500000, "Gas limit for the transaction")
	cmd.Flags().StringVar(&gasPrice, "gas-price", "20000000000", "Gas price in wei")
	cmd.Flags().Uint64Var(&txNonce, "nonce", 0, "Nonce for the transaction (required for signed transactions)")
	cmd.Flags().StringVar(&contractFee, "fee", "0.001", "Contract fee in MITO")

	// Add tx create specific flags
	AddTxCreateFlags(cmd)

	// Command-specific flags
	cmd.Flags().StringVar(&validator, "validator", "", "Validator address")
	cmd.Flags().StringVar(&receiver, "receiver", "", "Address to receive the withdrawn collateral")
	cmd.Flags().StringVar(&amount, "amount", "", "Amount of collateral to withdraw in MITO (e.g., \"1.5\")")

	// Mark required flags (contract is not required as it can come from config)
	mustMarkFlagRequired(cmd, "validator")
	mustMarkFlagRequired(cmd, "receiver")
	mustMarkFlagRequired(cmd, "amount")

	// Add PreRun to load config and validate mutual exclusive flags
	existingPreRun := cmd.PreRun
	cmd.PreRun = func(cmd *cobra.Command, args []string) {
		// Call the PreRun from AddTxCreateFlags first (for mutual exclusive validation)
		if existingPreRun != nil {
			existingPreRun(cmd, args)
		}

		// Additional validation specific to this command
		if validatorManagerContractAddr == "" {
			fmt.Println("Error: ValidatorManager contract address is required. Set it with --contract flag or use 'mito config set-contract <address>'")
			os.Exit(1)
		}
	}

	return cmd
}

func runWithdrawCollateralTx(cmd *cobra.Command, validator, amount, receiver string) error {
	// Parse collateral amount as decimal MITO and convert to wei
	collateralAmount, err := ParseValueAsWei(amount)
	if err != nil {
		return fmt.Errorf("failed to parse amount: %w", err)
	}

	// Ensure collateral amount is greater than 0
	if collateralAmount.Cmp(big.NewInt(0)) <= 0 {
		return fmt.Errorf("collateral amount must be greater than 0")
	}

	// Parse fee amount
	feeAmount, err := ParseValueAsWei(contractFee)
	if err != nil {
		return fmt.Errorf("failed to parse fee: %w", err)
	}

	// Validate validator address
	valAddr, err := ValidateAddress(validator)
	if err != nil {
		return fmt.Errorf("invalid validator address: %w", err)
	}

	// Validate receiver address
	receiverAddr, err := ValidateAddress(receiver)
	if err != nil {
		return fmt.Errorf("invalid receiver address: %w", err)
	}

	// Show summary
	fmt.Println("===== Withdraw Collateral Transaction =====")
	fmt.Printf("Validator Address        : %s\n", valAddr.Hex())
	fmt.Printf("Collateral Amount        : %s MITO\n", FormatWeiToEther(collateralAmount))
	fmt.Printf("Fee                      : %s MITO\n", FormatWeiToEther(feeAmount))
	fmt.Printf("Receiver Address         : %s\n", receiverAddr.Hex())
	fmt.Printf("Chain ID                 : %s\n", chainID)
	fmt.Printf("Gas Limit                : %d\n", gasLimit)
	fmt.Printf("Gas Price                : %s wei\n", gasPrice)
	if signed {
		fmt.Printf("Nonce                    : %d\n", txNonce)
	}
	fmt.Println()

	fmt.Println("🚨 IMPORTANT: Only the collateral owner can withdraw collateral.")
	fmt.Println("Make sure you are the owner of the collateral for this validator.")
	fmt.Println()

	// Check if we should create signed or unsigned transaction
	if signed {
		return createSignedWithdrawCollateralTransaction(valAddr, receiverAddr, collateralAmount, feeAmount)
	} else {
		return createUnsignedWithdrawCollateralTransaction(valAddr, receiverAddr, collateralAmount, feeAmount)
	}
}

func createSignedWithdrawCollateralTransaction(valAddr, receiverAddr common.Address, collateralAmount, feeAmount *big.Int) error {
	if !signed {
		return fmt.Errorf("signed transaction creation requires --signed flag")
	}

	// Get signing configuration
	signingConfig, err := GetSigningConfig()
	if err != nil {
		return fmt.Errorf("failed to get signing configuration: %w", err)
	}
	privateKey := signingConfig.PrivateKey

	// Get contract address
	contractAddr := common.HexToAddress(validatorManagerContractAddr)

	// Encode the function call data
	abi, err := bindings.IValidatorManagerMetaData.GetAbi()
	if err != nil {
		return fmt.Errorf("failed to get contract ABI: %w", err)
	}

	data, err := abi.Pack("withdrawCollateral", valAddr, receiverAddr, collateralAmount)
	if err != nil {
		return fmt.Errorf("failed to pack function call: %w", err)
	}

	// Parse gas price
	gasPriceBig, ok := new(big.Int).SetString(gasPrice, 10)
	if !ok {
		return fmt.Errorf("invalid gas price: %s", gasPrice)
	}

	// Parse chain ID
	chainIDBig, ok := new(big.Int).SetString(chainID, 10)
	if !ok {
		return fmt.Errorf("invalid chain ID: %s", chainID)
	}

	// Create and sign transaction (withdraw only sends fee, not collateral)
	tx := types.NewTransaction(txNonce, contractAddr, feeAmount, gasLimit, gasPriceBig, data)

	// Sign the transaction
	signer := types.NewEIP155Signer(chainIDBig)
	signedTx, err := types.SignTx(tx, signer, privateKey)
	if err != nil {
		return fmt.Errorf("failed to sign transaction: %w", err)
	}

	// Create signed transaction data
	signedTxData := SignedTransactionData{
		TransactionData: TransactionData{
			To:       signedTx.To().Hex(),
			Value:    signedTx.Value().String(),
			Data:     fmt.Sprintf("0x%x", signedTx.Data()),
			GasLimit: signedTx.Gas(),
			GasPrice: signedTx.GasPrice().String(),
			Nonce:    signedTx.Nonce(),
			ChainID:  chainID,
		},
		Hash: signedTx.Hash().Hex(),
	}

	// Get signature components
	v, r, s := signedTx.RawSignatureValues()
	signedTxData.Signature.V = v.String()
	signedTxData.Signature.R = r.String()
	signedTxData.Signature.S = s.String()

	// Get raw transaction data
	rawTxData, err := signedTx.MarshalBinary()
	if err == nil {
		signedTxData.Raw = fmt.Sprintf("0x%x", rawTxData)
	}

	// Output the transaction
	return outputTransaction(signedTxData, "Signed")
}

func createUnsignedWithdrawCollateralTransaction(valAddr, receiverAddr common.Address, collateralAmount, feeAmount *big.Int) error {
	// Get contract address
	contractAddr := common.HexToAddress(validatorManagerContractAddr)

	// Encode the function call data
	abi, err := bindings.IValidatorManagerMetaData.GetAbi()
	if err != nil {
		return fmt.Errorf("failed to get contract ABI: %w", err)
	}

	data, err := abi.Pack("withdrawCollateral", valAddr, receiverAddr, collateralAmount)
	if err != nil {
		return fmt.Errorf("failed to pack function call: %w", err)
	}

	// Create unsigned transaction data (withdraw only sends fee, not collateral)
	unsignedTxData := TransactionData{
		To:       contractAddr.Hex(),
		Value:    feeAmount.String(),
		Data:     fmt.Sprintf("0x%x", data),
		ChainID:  chainID,
		GasLimit: gasLimit,
		GasPrice: gasPrice,
	}

	// Output the transaction
	return outputTransaction(unsignedTxData, "Unsigned")
}
