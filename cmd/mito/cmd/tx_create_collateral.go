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
Only permitted collateral owners can deposit collateral for a validator.
Network information (chain-id, gas-price, fee) will be automatically fetched if RPC is configured,
otherwise you must provide them manually.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runDepositCollateralTx(cmd, validator, amount)
		},
	}

	// Add common flags for contract interaction
	cmd.Flags().StringVar(&validatorManagerContractAddr, "contract", "", "ValidatorManager contract address")

	// Add network information flags
	AddTxCreateNetworkFlags(cmd, true) // true because this command requires fee

	// Add tx create specific flags
	AddTxCreateFlags(cmd)

	// Command-specific flags
	cmd.Flags().StringVar(&validator, "validator", "", "Validator address")
	cmd.Flags().StringVar(&amount, "amount", "", "Amount of collateral to deposit in MITO (e.g., \"1.5\")")

	// Mark required flags (contract is not required as it can come from config)
	mustMarkFlagRequired(cmd, "validator")
	mustMarkFlagRequired(cmd, "amount")

	// Add PreRun to load config and validate network info
	existingPreRun := cmd.PreRun
	cmd.PreRun = func(cmd *cobra.Command, args []string) {
		// Call existing PreRun first (from AddTxCreateFlags)
		if existingPreRun != nil {
			existingPreRun(cmd, args)
		}

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

		// Validate network information requirements
		if err := ValidateNetworkInfo(cmd); err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}

		// Setup network information (fetch from RPC if available)
		if err := SetupNetworkInfo(cmd); err != nil {
			fmt.Printf("Error: %v\n", err)
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
Only the collateral owner can withdraw collateral.
Network information (chain-id, gas-price, fee) will be automatically fetched if RPC is configured,
otherwise you must provide them manually.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runWithdrawCollateralTx(cmd, validator, amount, receiver)
		},
	}

	// Add common flags for contract interaction
	cmd.Flags().StringVar(&validatorManagerContractAddr, "contract", "", "ValidatorManager contract address")

	// Add network information flags
	AddTxCreateNetworkFlags(cmd, true) // true because this command requires fee

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

	// Add PreRun to load config and validate network info
	existingPreRun := cmd.PreRun
	cmd.PreRun = func(cmd *cobra.Command, args []string) {
		// Call existing PreRun first (from AddTxCreateFlags)
		if existingPreRun != nil {
			existingPreRun(cmd, args)
		}

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

		// Validate network information requirements
		if err := ValidateNetworkInfo(cmd); err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}

		// Setup network information (fetch from RPC if available)
		if err := SetupNetworkInfo(cmd); err != nil {
			fmt.Printf("Error: %v\n", err)
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

// newTxCreateCollateralSetPermittedOwnerCmd creates the set permitted collateral owner command for tx create
func newTxCreateCollateralSetPermittedOwnerCmd() *cobra.Command {
	var (
		validator       string
		collateralOwner string
		isPermitted     bool
	)

	cmd := &cobra.Command{
		Use:   "set-permitted-owner",
		Short: "Create a set permitted collateral owner transaction",
		Long: `Create a transaction to set a permitted collateral owner for a validator.
This command can create either signed or unsigned transactions based on the flags provided.
Network information (chain-id, gas-price) will be automatically fetched if RPC is configured,
otherwise you must provide them manually.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runSetPermittedCollateralOwnerTx(cmd, validator, collateralOwner, isPermitted)
		},
	}

	// Add common flags for contract interaction
	cmd.Flags().StringVar(&validatorManagerContractAddr, "contract", "", "ValidatorManager contract address")

	// Add network information flags (no fee required for permission update)
	AddTxCreateNetworkFlags(cmd, false)

	// Add tx create specific flags
	AddTxCreateFlags(cmd)

	// Command-specific flags
	cmd.Flags().StringVar(&validator, "validator", "", "Validator address")
	cmd.Flags().StringVar(&collateralOwner, "collateral-owner", "", "Collateral owner address to permit or deny")
	cmd.Flags().BoolVar(&isPermitted, "permitted", false, "Whether to permit (true) or deny (false) the address")

	// Mark required flags
	mustMarkFlagRequired(cmd, "validator")
	mustMarkFlagRequired(cmd, "collateral-owner")

	// Add PreRun to load config and validate network info
	existingPreRun := cmd.PreRun
	cmd.PreRun = func(cmd *cobra.Command, args []string) {
		// Call existing PreRun first (from AddTxCreateFlags)
		if existingPreRun != nil {
			existingPreRun(cmd, args)
		}

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

		// Validate network information requirements
		if err := ValidateNetworkInfo(cmd); err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}

		// Setup network information (fetch from RPC if available)
		if err := SetupNetworkInfo(cmd); err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}
	}

	return cmd
}

func runSetPermittedCollateralOwnerTx(cmd *cobra.Command, validator, collateralOwner string, isPermitted bool) error {
	// Validate validator address
	valAddr, err := ValidateAddress(validator)
	if err != nil {
		return fmt.Errorf("invalid validator address: %w", err)
	}

	// Validate collateral owner address
	collateralOwnerAddr, err := ValidateAddress(collateralOwner)
	if err != nil {
		return fmt.Errorf("invalid collateral owner address: %w", err)
	}

	// Setup client and contract (if RPC is available)
	var currentPermission bool
	var hasPermissionInfo bool

	if rpcURL != "" {
		ethClient, err := GetEthClient(rpcURL)
		if err != nil {
			fmt.Printf("Warning: Could not connect to RPC: %v\n", err)
		} else {
			contract, err := GetValidatorManagerContract(ethClient)
			if err != nil {
				fmt.Printf("Warning: Could not initialize contract: %v\n", err)
			} else {
				// Check current permission status
				currentPermission, err = contract.IsPermittedCollateralOwner(nil, valAddr, collateralOwnerAddr)
				if err != nil {
					fmt.Printf("Warning: Could not get current permission status: %v\n", err)
				} else {
					hasPermissionInfo = true
				}
			}
		}
	}

	// Show summary
	fmt.Println("===== Set Permitted Collateral Owner Transaction =====")
	fmt.Printf("Validator Address          : %s\n", valAddr.Hex())
	fmt.Printf("Collateral Owner Address   : %s\n", collateralOwnerAddr.Hex())
	if hasPermissionInfo {
		fmt.Printf("Current Permission Status  : %t\n", currentPermission)
	}
	fmt.Printf("New Permission Status      : %t\n", isPermitted)
	fmt.Printf("Chain ID                   : %s\n", chainID)
	fmt.Printf("Gas Limit                  : %d\n", gasLimit)
	fmt.Printf("Gas Price                  : %s wei\n", gasPrice)
	if signed {
		fmt.Printf("Nonce                      : %d\n", txNonce)
	}
	fmt.Println()

	// Check if we should create signed or unsigned transaction
	if signed {
		return createSignedSetPermittedCollateralOwnerTransaction(valAddr, collateralOwnerAddr, isPermitted)
	} else {
		return createUnsignedSetPermittedCollateralOwnerTransaction(valAddr, collateralOwnerAddr, isPermitted)
	}
}

func createSignedSetPermittedCollateralOwnerTransaction(valAddr, collateralOwnerAddr common.Address, isPermitted bool) error {
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

	data, err := abi.Pack("setPermittedCollateralOwner", valAddr, collateralOwnerAddr, isPermitted)
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

	// Create and sign transaction (no value needed for permission update)
	tx := types.NewTransaction(txNonce, contractAddr, big.NewInt(0), gasLimit, gasPriceBig, data)

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

func createUnsignedSetPermittedCollateralOwnerTransaction(valAddr, collateralOwnerAddr common.Address, isPermitted bool) error {
	// Get contract address
	contractAddr := common.HexToAddress(validatorManagerContractAddr)

	// Encode the function call data
	abi, err := bindings.IValidatorManagerMetaData.GetAbi()
	if err != nil {
		return fmt.Errorf("failed to get contract ABI: %w", err)
	}

	data, err := abi.Pack("setPermittedCollateralOwner", valAddr, collateralOwnerAddr, isPermitted)
	if err != nil {
		return fmt.Errorf("failed to pack function call: %w", err)
	}

	// Create unsigned transaction data (no value needed for permission update)
	unsignedTxData := TransactionData{
		To:       contractAddr.Hex(),
		Value:    "0",
		Data:     fmt.Sprintf("0x%x", data),
		ChainID:  chainID,
		GasLimit: gasLimit,
		GasPrice: gasPrice,
	}

	// Output the transaction
	return outputTransaction(unsignedTxData, "Unsigned")
}

// newTxCreateCollateralTransferOwnershipCmd creates the transfer collateral ownership command for tx create
func newTxCreateCollateralTransferOwnershipCmd() *cobra.Command {
	var (
		validator string
		newOwner  string
	)

	cmd := &cobra.Command{
		Use:   "transfer-ownership",
		Short: "Create a transfer collateral ownership transaction",
		Long: `Create a transaction to transfer collateral ownership for a validator.
This command can create either signed or unsigned transactions based on the flags provided.
Network information (chain-id, gas-price, fee) will be automatically fetched if RPC is configured,
otherwise you must provide them manually.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runTransferCollateralOwnershipTx(cmd, validator, newOwner)
		},
	}

	// Add common flags for contract interaction
	cmd.Flags().StringVar(&validatorManagerContractAddr, "contract", "", "ValidatorManager contract address")

	// Add network information flags
	AddTxCreateNetworkFlags(cmd, true) // true because this command requires fee

	// Add tx create specific flags
	AddTxCreateFlags(cmd)

	// Command-specific flags
	cmd.Flags().StringVar(&validator, "validator", "", "Validator address")
	cmd.Flags().StringVar(&newOwner, "new-owner", "", "New collateral owner address")

	// Mark required flags
	mustMarkFlagRequired(cmd, "validator")
	mustMarkFlagRequired(cmd, "new-owner")

	// Add PreRun to load config and validate network info
	existingPreRun := cmd.PreRun
	cmd.PreRun = func(cmd *cobra.Command, args []string) {
		// Call existing PreRun first (from AddTxCreateFlags)
		if existingPreRun != nil {
			existingPreRun(cmd, args)
		}

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

		// Validate network information requirements
		if err := ValidateNetworkInfo(cmd); err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}

		// Setup network information (fetch from RPC if available)
		if err := SetupNetworkInfo(cmd); err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}
	}

	return cmd
}

func runTransferCollateralOwnershipTx(cmd *cobra.Command, validator, newOwner string) error {
	// Validate validator address
	valAddr, err := ValidateAddress(validator)
	if err != nil {
		return fmt.Errorf("invalid validator address: %w", err)
	}

	// Validate new owner address
	newOwnerAddr, err := ValidateAddress(newOwner)
	if err != nil {
		return fmt.Errorf("invalid new owner address: %w", err)
	}

	// Parse fee amount
	feeAmount, err := ParseValueAsWei(contractFee)
	if err != nil {
		return fmt.Errorf("failed to parse fee: %w", err)
	}

	// Setup client and contract (if RPC is available)
	if rpcURL != "" {
		ethClient, err := GetEthClient(rpcURL)
		if err != nil {
			fmt.Printf("Warning: Could not connect to RPC: %v\n", err)
		} else {
			contract, err := GetValidatorManagerContract(ethClient)
			if err != nil {
				fmt.Printf("Warning: Could not initialize contract: %v\n", err)
			} else {
				// Get the actual contract fee
				actualFee, err := contract.Fee(nil)
				if err != nil {
					fmt.Printf("Warning: Could not get contract fee: %v\n", err)
				} else {
					feeAmount = actualFee
				}
			}
		}
	}

	// Show summary
	fmt.Println("===== Transfer Collateral Ownership Transaction =====")
	fmt.Printf("Validator Address        : %s\n", valAddr.Hex())
	fmt.Printf("New Collateral Owner     : %s\n", newOwnerAddr.Hex())
	fmt.Printf("Fee                      : %s MITO\n", FormatWeiToEther(feeAmount))
	fmt.Printf("Chain ID                 : %s\n", chainID)
	fmt.Printf("Gas Limit                : %d\n", gasLimit)
	fmt.Printf("Gas Price                : %s wei\n", gasPrice)
	if signed {
		fmt.Printf("Nonce                    : %d\n", txNonce)
	}
	fmt.Println()

	// Check if we should create signed or unsigned transaction
	if signed {
		return createSignedTransferCollateralOwnershipTransaction(valAddr, newOwnerAddr, feeAmount)
	} else {
		return createUnsignedTransferCollateralOwnershipTransaction(valAddr, newOwnerAddr, feeAmount)
	}
}

func createSignedTransferCollateralOwnershipTransaction(valAddr, newOwnerAddr common.Address, feeAmount *big.Int) error {
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

	data, err := abi.Pack("transferCollateralOwnership", valAddr, newOwnerAddr)
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

	// Create and sign transaction (fee is required for ownership transfer)
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

func createUnsignedTransferCollateralOwnershipTransaction(valAddr, newOwnerAddr common.Address, feeAmount *big.Int) error {
	// Get contract address
	contractAddr := common.HexToAddress(validatorManagerContractAddr)

	// Encode the function call data
	abi, err := bindings.IValidatorManagerMetaData.GetAbi()
	if err != nil {
		return fmt.Errorf("failed to get contract ABI: %w", err)
	}

	data, err := abi.Pack("transferCollateralOwnership", valAddr, newOwnerAddr)
	if err != nil {
		return fmt.Errorf("failed to pack function call: %w", err)
	}

	// Create unsigned transaction data (fee is required for ownership transfer)
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
