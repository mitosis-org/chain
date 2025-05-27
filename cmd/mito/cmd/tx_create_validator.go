package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"math/big"
	"os"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/mitosis-org/chain/bindings"
	"github.com/spf13/cobra"
)

// Transaction data structure for JSON output
type TransactionData struct {
	To       string `json:"to"`
	Value    string `json:"value"`
	Data     string `json:"data"`
	GasLimit uint64 `json:"gasLimit,omitempty"`
	GasPrice string `json:"gasPrice,omitempty"`
	Nonce    uint64 `json:"nonce,omitempty"`
	ChainID  string `json:"chainId,omitempty"`
}

// SignedTransactionData structure for signed transactions
type SignedTransactionData struct {
	TransactionData
	Signature struct {
		V string `json:"v"`
		R string `json:"r"`
		S string `json:"s"`
	} `json:"signature"`
	Hash string `json:"hash"`
	Raw  string `json:"raw"`
}

// newTxCreateValidatorCreateCmd creates the validator creation command for tx create
func newTxCreateValidatorCreateCmd() *cobra.Command {
	var (
		pubKey            string
		operator          string
		rewardManager     string
		commissionRate    string
		metadata          string
		initialCollateral string
	)

	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a validator creation transaction",
		Long: `Create a transaction to register a new validator in the ValidatorManager contract.
This command can create either signed or unsigned transactions based on the flags provided.
When creating offline transactions, you need to provide fee, gas parameters manually.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runCreateValidatorTx(cmd, pubKey, operator, rewardManager, commissionRate, metadata, initialCollateral)
		},
	}

	// Add common flags for contract interaction
	cmd.Flags().StringVar(&validatorManagerContractAddr, "contract", "", "ValidatorManager contract address")

	// Add offline mode flags
	cmd.Flags().BoolVar(&offlineMode, "offline", false, "Create transaction offline without RPC connection")
	cmd.Flags().StringVar(&chainID, "chain-id", "1337", "Chain ID for the transaction")
	cmd.Flags().Uint64Var(&gasLimit, "gas-limit", 500000, "Gas limit for the transaction")
	cmd.Flags().StringVar(&gasPrice, "gas-price", "20000000000", "Gas price in wei")
	cmd.Flags().Uint64Var(&txNonce, "nonce", 0, "Nonce for the transaction (required for signed transactions)")
	cmd.Flags().StringVar(&contractFee, "fee", "0.001", "Contract fee in MITO (for offline mode)")

	// Add tx create specific flags
	AddTxCreateFlags(cmd)

	// Command-specific flags
	cmd.Flags().StringVar(&pubKey, "pubkey", "", "Validator's public key (hex with 0x prefix)")
	cmd.Flags().StringVar(&operator, "operator", "", "Operator address")
	cmd.Flags().StringVar(&rewardManager, "reward-manager", "", "Reward manager address")
	cmd.Flags().StringVar(&commissionRate, "commission-rate", "", "Commission rate in percentage (e.g., \"5%\")")
	cmd.Flags().StringVar(&metadata, "metadata", "", "Validator metadata (JSON string)")
	cmd.Flags().StringVar(&initialCollateral, "initial-collateral", "", "Initial collateral amount in MITO (e.g., \"1.5\")")

	// Mark required flags (contract is not required as it can come from config)
	mustMarkFlagRequired(cmd, "pubkey")
	mustMarkFlagRequired(cmd, "operator")
	mustMarkFlagRequired(cmd, "reward-manager")
	mustMarkFlagRequired(cmd, "commission-rate")
	mustMarkFlagRequired(cmd, "metadata")
	mustMarkFlagRequired(cmd, "initial-collateral")

	// Add PreRun to load config
	cmd.PreRun = func(cmd *cobra.Command, args []string) {
		// Load global config
		if err := loadGlobalConfig(); err != nil {
			fmt.Printf("Warning: failed to load config: %v\n", err)
		}

		// Resolve config values
		resolveConfigValues()

		// Validate that required values are set
		if validatorManagerContractAddr == "" {
			fmt.Println("Error: ValidatorManager contract address is required. Set it with --contract flag or use 'mito config set-contract <address>'")
			os.Exit(1)
		}
	}

	return cmd
}

func runCreateValidatorTx(cmd *cobra.Command, pubKey, operator, rewardManager, commissionRate, metadata, initialCollateral string) error {
	// Parse collateral amount as decimal MITO and convert to wei
	collateralAmount, err := ParseValueAsWei(initialCollateral)
	if err != nil {
		return fmt.Errorf("failed to parse initial collateral: %w", err)
	}

	// Parse fee amount
	feeAmount, err := ParseValueAsWei(contractFee)
	if err != nil {
		return fmt.Errorf("failed to parse fee: %w", err)
	}

	// Calculate total transaction value (collateral + fee)
	totalValue := new(big.Int).Add(collateralAmount, feeAmount)

	// Validate other parameters
	operatorAddr, err := ValidateAddress(operator)
	if err != nil {
		return fmt.Errorf("invalid operator address: %w", err)
	}

	rewardManagerAddr, err := ValidateAddress(rewardManager)
	if err != nil {
		return fmt.Errorf("invalid reward manager address: %w", err)
	}

	// Parse commission rate
	commissionRateInt, err := ParsePercentageToBasisPoints(commissionRate)
	if err != nil {
		return fmt.Errorf("failed to parse commission rate: %w", err)
	}

	// Decode public key from hex
	pubKeyBytes, err := DecodeHexWithPrefix(pubKey)
	if err != nil {
		return fmt.Errorf("failed to decode public key: %w", err)
	}

	// Create the request
	request := bindings.IValidatorManagerCreateValidatorRequest{
		Operator:       operatorAddr,
		RewardManager:  rewardManagerAddr,
		CommissionRate: commissionRateInt,
		Metadata:       []byte(metadata),
	}

	// Show summary
	fmt.Println("===== Create Validator Transaction =====")
	fmt.Printf("Public Key                 : %s\n", pubKey)
	fmt.Printf("Operator                   : %s\n", operatorAddr.Hex())
	fmt.Printf("Reward Manager             : %s\n", rewardManagerAddr.Hex())
	fmt.Printf("Commission Rate            : %s\n", FormatBasisPointsToPercent(commissionRateInt))
	fmt.Printf("Metadata                   : %s\n", metadata)
	fmt.Printf("Initial Collateral         : %s MITO\n", FormatWeiToEther(collateralAmount))
	fmt.Printf("Fee                        : %s MITO\n", FormatWeiToEther(feeAmount))
	fmt.Printf("Total Value                : %s MITO\n", FormatWeiToEther(totalValue))
	fmt.Printf("Chain ID                   : %s\n", chainID)
	fmt.Printf("Gas Limit                  : %d\n", gasLimit)
	fmt.Printf("Gas Price                  : %s wei\n", gasPrice)
	if signed {
		fmt.Printf("Nonce                      : %d\n", nonce)
	}
	fmt.Println()

	// Check if we should create signed or unsigned transaction
	if signed {
		return createSignedTransaction(pubKeyBytes, request, totalValue)
	} else {
		return createUnsignedTransaction(pubKeyBytes, request, totalValue)
	}
}

func createSignedTransaction(pubKeyBytes []byte, request bindings.IValidatorManagerCreateValidatorRequest, totalValue *big.Int) error {
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

	data, err := abi.Pack("createValidator", pubKeyBytes, request)
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

func createUnsignedTransaction(pubKeyBytes []byte, request bindings.IValidatorManagerCreateValidatorRequest, totalValue *big.Int) error {
	// Get contract address
	contractAddr := common.HexToAddress(validatorManagerContractAddr)

	// Encode the function call data
	abi, err := bindings.IValidatorManagerMetaData.GetAbi()
	if err != nil {
		return fmt.Errorf("failed to get contract ABI: %w", err)
	}

	data, err := abi.Pack("createValidator", pubKeyBytes, request)
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

// newTxCreateValidatorUpdateOperatorCmd creates the validator operator update command for tx create
func newTxCreateValidatorUpdateOperatorCmd() *cobra.Command {
	var (
		validator string
		operator  string
	)

	cmd := &cobra.Command{
		Use:   "update-operator",
		Short: "Create a validator operator update transaction",
		Long: `Create a transaction to update the operator address for an existing validator.
This command can create either signed or unsigned transactions based on the flags provided.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runUpdateOperatorTx(cmd, validator, operator)
		},
	}

	// Add common flags for contract interaction
	cmd.Flags().StringVar(&validatorManagerContractAddr, "contract", "", "ValidatorManager contract address")
	mustMarkFlagRequired(cmd, "contract")

	// Add offline mode flags
	cmd.Flags().BoolVar(&offlineMode, "offline", false, "Create transaction offline without RPC connection")
	cmd.Flags().StringVar(&chainID, "chain-id", "1337", "Chain ID for the transaction")
	cmd.Flags().Uint64Var(&gasLimit, "gas-limit", 500000, "Gas limit for the transaction")
	cmd.Flags().StringVar(&gasPrice, "gas-price", "20000000000", "Gas price in wei")
	cmd.Flags().Uint64Var(&nonce, "nonce", 0, "Nonce for the transaction (required for signed transactions)")

	// Add tx create specific flags
	AddTxCreateFlags(cmd)

	// Command-specific flags
	cmd.Flags().StringVar(&validator, "validator", "", "Validator address")
	cmd.Flags().StringVar(&operator, "operator", "", "New operator address")

	// Mark required flags
	mustMarkFlagRequired(cmd, "validator")
	mustMarkFlagRequired(cmd, "operator")

	return cmd
}

func runUpdateOperatorTx(cmd *cobra.Command, validator, operator string) error {
	// Setup client and contract
	ethClient, err := GetEthClient(rpcURL)
	if err != nil {
		return fmt.Errorf("failed to connect to Ethereum client: %w", err)
	}

	contract, err := GetValidatorManagerContract(ethClient)
	if err != nil {
		return fmt.Errorf("failed to initialize contract: %w", err)
	}

	// Validate validator address
	valAddr, err := ValidateAddress(validator)
	if err != nil {
		return fmt.Errorf("invalid validator address: %w", err)
	}

	// Validate operator address
	operatorAddr, err := ValidateAddress(operator)
	if err != nil {
		return fmt.Errorf("invalid operator address: %w", err)
	}

	// Get validator info to show current values
	validatorInfo, err := contract.ValidatorInfo(nil, valAddr)
	if err != nil {
		return fmt.Errorf("failed to get validator info: %w", err)
	}

	// Show summary
	fmt.Println("===== Update Operator Transaction =====")
	fmt.Printf("Validator Address            : %s\n", valAddr.Hex())
	fmt.Printf("Current Operator             : %s\n", validatorInfo.Operator.Hex())
	fmt.Printf("New Operator                 : %s\n", operatorAddr.Hex())
	fmt.Printf("Current Reward Manager       : %s\n", validatorInfo.RewardManager.Hex())
	fmt.Println()

	// Show important warning
	fmt.Println("🚨 IMPORTANT WARNING 🚨")
	fmt.Println("When changing the operator address, you may also want to update:")
	fmt.Println("1. Reward Manager - To ensure rewards are managed by the correct entity")
	fmt.Println("2. Collateral Ownership - To manage who can deposit or withdraw collateral")
	fmt.Println()

	// Check if we should create signed or unsigned transaction
	if signed {
		return createSignedUpdateOperatorTransaction(ethClient, contract, valAddr, operatorAddr)
	} else {
		return createUnsignedUpdateOperatorTransaction(ethClient, contract, valAddr, operatorAddr)
	}
}

func createSignedUpdateOperatorTransaction(ethClient *ethclient.Client, contract *bindings.IValidatorManager, valAddr, operatorAddr common.Address) error {
	// Create transaction options
	opts, err := TransactOpts(nil)
	if err != nil {
		return fmt.Errorf("failed to create transaction options: %w", err)
	}

	// Create the transaction
	tx, err := contract.UpdateOperator(opts, valAddr, operatorAddr)
	if err != nil {
		return fmt.Errorf("failed to create update operator transaction: %w", err)
	}

	// Create signed transaction data
	signedTxData := SignedTransactionData{
		TransactionData: TransactionData{
			To:       tx.To().Hex(),
			Value:    tx.Value().String(),
			Data:     fmt.Sprintf("0x%x", tx.Data()),
			GasLimit: tx.Gas(),
			GasPrice: tx.GasPrice().String(),
			Nonce:    tx.Nonce(),
		},
		Hash: tx.Hash().Hex(),
	}

	// Get chain ID
	chainID, err := ethClient.ChainID(context.Background())
	if err == nil {
		signedTxData.ChainID = chainID.String()
	}

	// Get signature components
	v, r, s := tx.RawSignatureValues()
	signedTxData.Signature.V = v.String()
	signedTxData.Signature.R = r.String()
	signedTxData.Signature.S = s.String()

	// Get raw transaction data
	rawTxData, err := tx.MarshalBinary()
	if err == nil {
		signedTxData.Raw = fmt.Sprintf("0x%x", rawTxData)
	}

	// Output the transaction
	return outputTransaction(signedTxData, "Signed")
}

func createUnsignedUpdateOperatorTransaction(ethClient *ethclient.Client, contract *bindings.IValidatorManager, valAddr, operatorAddr common.Address) error {
	// Get contract address
	contractAddr := common.HexToAddress(validatorManagerContractAddr)

	// Encode the function call data
	abi, err := bindings.IValidatorManagerMetaData.GetAbi()
	if err != nil {
		return fmt.Errorf("failed to get contract ABI: %w", err)
	}

	data, err := abi.Pack("updateOperator", valAddr, operatorAddr)
	if err != nil {
		return fmt.Errorf("failed to pack function call: %w", err)
	}

	// Get chain ID
	chainID, err := ethClient.ChainID(context.Background())
	if err != nil {
		return fmt.Errorf("failed to get chain ID: %w", err)
	}

	// Create unsigned transaction data
	unsignedTxData := TransactionData{
		To:      contractAddr.Hex(),
		Value:   "0",
		Data:    fmt.Sprintf("0x%x", data),
		ChainID: chainID.String(),
	}

	// Estimate gas if possible
	if gasLimit, err := ethClient.EstimateGas(context.Background(), ethereum.CallMsg{
		To:   &contractAddr,
		Data: data,
	}); err == nil {
		unsignedTxData.GasLimit = gasLimit
	}

	// Get suggested gas price
	if gasPrice, err := ethClient.SuggestGasPrice(context.Background()); err == nil {
		unsignedTxData.GasPrice = gasPrice.String()
	}

	// Output the transaction
	return outputTransaction(unsignedTxData, "Unsigned")
}

// newTxCreateValidatorUpdateMetadataCmd creates the validator metadata update command for tx create
func newTxCreateValidatorUpdateMetadataCmd() *cobra.Command {
	var (
		validator string
		metadata  string
	)

	cmd := &cobra.Command{
		Use:   "update-metadata",
		Short: "Create a validator metadata update transaction",
		Long: `Create a transaction to update the metadata for an existing validator.
This command can create either signed or unsigned transactions based on the flags provided.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runUpdateMetadataTx(cmd, validator, metadata)
		},
	}

	// Add common flags for contract interaction
	cmd.Flags().StringVar(&validatorManagerContractAddr, "contract", "", "ValidatorManager contract address")
	mustMarkFlagRequired(cmd, "contract")

	// Add offline mode flags
	cmd.Flags().BoolVar(&offlineMode, "offline", false, "Create transaction offline without RPC connection")
	cmd.Flags().StringVar(&chainID, "chain-id", "1337", "Chain ID for the transaction")
	cmd.Flags().Uint64Var(&gasLimit, "gas-limit", 500000, "Gas limit for the transaction")
	cmd.Flags().StringVar(&gasPrice, "gas-price", "20000000000", "Gas price in wei")
	cmd.Flags().Uint64Var(&nonce, "nonce", 0, "Nonce for the transaction (required for signed transactions)")

	// Add tx create specific flags
	AddTxCreateFlags(cmd)

	// Command-specific flags
	cmd.Flags().StringVar(&validator, "validator", "", "Validator address")
	cmd.Flags().StringVar(&metadata, "metadata", "", "New metadata (JSON string)")

	// Mark required flags
	mustMarkFlagRequired(cmd, "validator")
	mustMarkFlagRequired(cmd, "metadata")

	return cmd
}

func runUpdateMetadataTx(cmd *cobra.Command, validator, metadata string) error {
	// Setup client and contract
	ethClient, err := GetEthClient(rpcURL)
	if err != nil {
		return fmt.Errorf("failed to connect to Ethereum client: %w", err)
	}

	contract, err := GetValidatorManagerContract(ethClient)
	if err != nil {
		return fmt.Errorf("failed to initialize contract: %w", err)
	}

	// Validate validator address
	valAddr, err := ValidateAddress(validator)
	if err != nil {
		return fmt.Errorf("invalid validator address: %w", err)
	}

	// Get validator info to show current values
	validatorInfo, err := contract.ValidatorInfo(nil, valAddr)
	if err != nil {
		return fmt.Errorf("failed to get validator info: %w", err)
	}

	// Show summary
	fmt.Println("===== Update Metadata Transaction =====")
	fmt.Printf("Validator Address        : %s\n", valAddr.Hex())
	fmt.Printf("Current Metadata         : %s\n", string(validatorInfo.Metadata))
	fmt.Printf("New Metadata             : %s\n", metadata)
	fmt.Println()

	// Check if we should create signed or unsigned transaction
	if signed {
		return createSignedUpdateMetadataTransaction(ethClient, contract, valAddr, metadata)
	} else {
		return createUnsignedUpdateMetadataTransaction(ethClient, contract, valAddr, metadata)
	}
}

func createSignedUpdateMetadataTransaction(ethClient *ethclient.Client, contract *bindings.IValidatorManager, valAddr common.Address, metadata string) error {
	// Create transaction options
	opts, err := TransactOpts(nil)
	if err != nil {
		return fmt.Errorf("failed to create transaction options: %w", err)
	}

	// Create the transaction
	tx, err := contract.UpdateMetadata(opts, valAddr, []byte(metadata))
	if err != nil {
		return fmt.Errorf("failed to create update metadata transaction: %w", err)
	}

	// Create signed transaction data
	signedTxData := SignedTransactionData{
		TransactionData: TransactionData{
			To:       tx.To().Hex(),
			Value:    tx.Value().String(),
			Data:     fmt.Sprintf("0x%x", tx.Data()),
			GasLimit: tx.Gas(),
			GasPrice: tx.GasPrice().String(),
			Nonce:    tx.Nonce(),
		},
		Hash: tx.Hash().Hex(),
	}

	// Get chain ID
	chainID, err := ethClient.ChainID(context.Background())
	if err == nil {
		signedTxData.ChainID = chainID.String()
	}

	// Get signature components
	v, r, s := tx.RawSignatureValues()
	signedTxData.Signature.V = v.String()
	signedTxData.Signature.R = r.String()
	signedTxData.Signature.S = s.String()

	// Get raw transaction data
	rawTxData, err := tx.MarshalBinary()
	if err == nil {
		signedTxData.Raw = fmt.Sprintf("0x%x", rawTxData)
	}

	// Output the transaction
	return outputTransaction(signedTxData, "Signed")
}

func createUnsignedUpdateMetadataTransaction(ethClient *ethclient.Client, contract *bindings.IValidatorManager, valAddr common.Address, metadata string) error {
	// Get contract address
	contractAddr := common.HexToAddress(validatorManagerContractAddr)

	// Encode the function call data
	abi, err := bindings.IValidatorManagerMetaData.GetAbi()
	if err != nil {
		return fmt.Errorf("failed to get contract ABI: %w", err)
	}

	data, err := abi.Pack("updateMetadata", valAddr, []byte(metadata))
	if err != nil {
		return fmt.Errorf("failed to pack function call: %w", err)
	}

	// Get chain ID
	chainID, err := ethClient.ChainID(context.Background())
	if err != nil {
		return fmt.Errorf("failed to get chain ID: %w", err)
	}

	// Create unsigned transaction data
	unsignedTxData := TransactionData{
		To:      contractAddr.Hex(),
		Value:   "0",
		Data:    fmt.Sprintf("0x%x", data),
		ChainID: chainID.String(),
	}

	// Estimate gas if possible
	if gasLimit, err := ethClient.EstimateGas(context.Background(), ethereum.CallMsg{
		To:   &contractAddr,
		Data: data,
	}); err == nil {
		unsignedTxData.GasLimit = gasLimit
	}

	// Get suggested gas price
	if gasPrice, err := ethClient.SuggestGasPrice(context.Background()); err == nil {
		unsignedTxData.GasPrice = gasPrice.String()
	}

	// Output the transaction
	return outputTransaction(unsignedTxData, "Unsigned")
}

// newTxCreateValidatorUnjailCmd creates the validator unjail command for tx create
func newTxCreateValidatorUnjailCmd() *cobra.Command {
	var validator string

	cmd := &cobra.Command{
		Use:   "unjail",
		Short: "Create a validator unjail transaction",
		Long: `Create a transaction to unjail a validator that has been jailed.
This command can create either signed or unsigned transactions based on the flags provided.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runUnjailValidatorTx(cmd, validator)
		},
	}

	// Add common flags for contract interaction
	cmd.Flags().StringVar(&validatorManagerContractAddr, "contract", "", "ValidatorManager contract address")
	mustMarkFlagRequired(cmd, "contract")

	// Add offline mode flags
	cmd.Flags().BoolVar(&offlineMode, "offline", false, "Create transaction offline without RPC connection")
	cmd.Flags().StringVar(&chainID, "chain-id", "1337", "Chain ID for the transaction")
	cmd.Flags().Uint64Var(&gasLimit, "gas-limit", 500000, "Gas limit for the transaction")
	cmd.Flags().StringVar(&gasPrice, "gas-price", "20000000000", "Gas price in wei")
	cmd.Flags().Uint64Var(&nonce, "nonce", 0, "Nonce for the transaction (required for signed transactions)")

	// Add tx create specific flags
	AddTxCreateFlags(cmd)

	// Command-specific flags
	cmd.Flags().StringVar(&validator, "validator", "", "Address of the validator to unjail")

	// Mark required flags
	mustMarkFlagRequired(cmd, "validator")

	return cmd
}

func runUnjailValidatorTx(cmd *cobra.Command, validator string) error {
	// Setup client and contract
	ethClient, err := GetEthClient(rpcURL)
	if err != nil {
		return fmt.Errorf("failed to connect to Ethereum client: %w", err)
	}

	contract, err := GetValidatorManagerContract(ethClient)
	if err != nil {
		return fmt.Errorf("failed to initialize contract: %w", err)
	}

	// Get the contract fee
	fee, err := contract.Fee(nil)
	if err != nil {
		return fmt.Errorf("failed to get contract fee: %w", err)
	}

	// Validate validator address
	valAddr, err := ValidateAddress(validator)
	if err != nil {
		return fmt.Errorf("invalid validator address: %w", err)
	}

	// Show summary
	fmt.Println("===== Unjail Validator Transaction =====")
	fmt.Printf("Validator Address        : %s\n", valAddr.Hex())
	fmt.Printf("Fee                      : %s MITO\n", FormatWeiToEther(fee))
	fmt.Println()

	// Check if we should create signed or unsigned transaction
	if signed {
		return createSignedUnjailValidatorTransaction(ethClient, contract, valAddr, fee)
	} else {
		return createUnsignedUnjailValidatorTransaction(ethClient, contract, valAddr, fee)
	}
}

func createSignedUnjailValidatorTransaction(ethClient *ethclient.Client, contract *bindings.IValidatorManager, valAddr common.Address, fee *big.Int) error {
	// Create transaction options
	opts, err := TransactOpts(fee)
	if err != nil {
		return fmt.Errorf("failed to create transaction options: %w", err)
	}

	// Create the transaction
	tx, err := contract.UnjailValidator(opts, valAddr)
	if err != nil {
		return fmt.Errorf("failed to create unjail validator transaction: %w", err)
	}

	// Create signed transaction data
	signedTxData := SignedTransactionData{
		TransactionData: TransactionData{
			To:       tx.To().Hex(),
			Value:    tx.Value().String(),
			Data:     fmt.Sprintf("0x%x", tx.Data()),
			GasLimit: tx.Gas(),
			GasPrice: tx.GasPrice().String(),
			Nonce:    tx.Nonce(),
		},
		Hash: tx.Hash().Hex(),
	}

	// Get chain ID
	chainID, err := ethClient.ChainID(context.Background())
	if err == nil {
		signedTxData.ChainID = chainID.String()
	}

	// Get signature components
	v, r, s := tx.RawSignatureValues()
	signedTxData.Signature.V = v.String()
	signedTxData.Signature.R = r.String()
	signedTxData.Signature.S = s.String()

	// Get raw transaction data
	rawTxData, err := tx.MarshalBinary()
	if err == nil {
		signedTxData.Raw = fmt.Sprintf("0x%x", rawTxData)
	}

	// Output the transaction
	return outputTransaction(signedTxData, "Signed")
}

func createUnsignedUnjailValidatorTransaction(ethClient *ethclient.Client, contract *bindings.IValidatorManager, valAddr common.Address, fee *big.Int) error {
	// Get contract address
	contractAddr := common.HexToAddress(validatorManagerContractAddr)

	// Encode the function call data
	abi, err := bindings.IValidatorManagerMetaData.GetAbi()
	if err != nil {
		return fmt.Errorf("failed to get contract ABI: %w", err)
	}

	data, err := abi.Pack("unjailValidator", valAddr)
	if err != nil {
		return fmt.Errorf("failed to pack function call: %w", err)
	}

	// Get chain ID
	chainID, err := ethClient.ChainID(context.Background())
	if err != nil {
		return fmt.Errorf("failed to get chain ID: %w", err)
	}

	// Create unsigned transaction data
	unsignedTxData := TransactionData{
		To:      contractAddr.Hex(),
		Value:   fee.String(),
		Data:    fmt.Sprintf("0x%x", data),
		ChainID: chainID.String(),
	}

	// Estimate gas if possible
	if gasLimit, err := ethClient.EstimateGas(context.Background(), ethereum.CallMsg{
		To:    &contractAddr,
		Value: fee,
		Data:  data,
	}); err == nil {
		unsignedTxData.GasLimit = gasLimit
	}

	// Get suggested gas price
	if gasPrice, err := ethClient.SuggestGasPrice(context.Background()); err == nil {
		unsignedTxData.GasPrice = gasPrice.String()
	}

	// Output the transaction
	return outputTransaction(unsignedTxData, "Unsigned")
}

func outputTransaction(txData interface{}, txType string) error {
	// Convert to JSON
	jsonData, err := json.MarshalIndent(txData, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal transaction data: %w", err)
	}

	fmt.Printf("===== %s Transaction Data =====\n", txType)

	// Output to file or stdout
	if outputFile != "" {
		err = os.WriteFile(outputFile, jsonData, 0644)
		if err != nil {
			return fmt.Errorf("failed to write to file %s: %w", outputFile, err)
		}
		fmt.Printf("Transaction data written to: %s\n", outputFile)
	} else {
		fmt.Println(string(jsonData))
	}

	return nil
}
