package cmd

import (
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/mitosis-org/chain/bindings"
	"github.com/spf13/cobra"
)

// newTxSendCollateralDepositCmd creates the collateral deposit command for tx send
func newTxSendCollateralDepositCmd() *cobra.Command {
	var (
		validator string
		amount    string
	)

	cmd := &cobra.Command{
		Use:   "deposit",
		Short: "Send a collateral deposit transaction",
		Long: `Create, sign, and send a transaction to deposit collateral for a validator.
Only permitted collateral owners can deposit collateral for a validator.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runSendDepositCollateralTx(cmd, validator, amount)
		},
	}

	// Add common flags for contract interaction
	cmd.Flags().StringVar(&rpcURL, "rpc", "", "RPC endpoint URL")
	cmd.Flags().StringVar(&validatorManagerContractAddr, "contract", "", "ValidatorManager contract address")
	mustMarkFlagRequired(cmd, "rpc")
	mustMarkFlagRequired(cmd, "contract")

	// Add signing flags
	AddSigningFlags(cmd)

	// Command-specific flags
	cmd.Flags().StringVar(&validator, "validator", "", "Validator address")
	cmd.Flags().StringVar(&amount, "amount", "", "Amount of collateral to deposit in MITO (e.g., \"1.5\")")

	// Mark required flags
	mustMarkFlagRequired(cmd, "validator")
	mustMarkFlagRequired(cmd, "amount")

	return cmd
}

func runSendDepositCollateralTx(cmd *cobra.Command, validator, amount string) error {
	// Parse collateral amount as decimal MITO and convert to wei
	collateralAmount, err := ParseValueAsWei(amount)
	if err != nil {
		return fmt.Errorf("failed to parse amount: %w", err)
	}

	// Ensure collateral amount is greater than 0
	if collateralAmount.Cmp(big.NewInt(0)) <= 0 {
		return fmt.Errorf("collateral amount must be greater than 0")
	}

	// Validate validator address
	valAddr, err := ValidateAddress(validator)
	if err != nil {
		return fmt.Errorf("invalid validator address: %w", err)
	}

	// Get signing configuration
	signingConfig, err := GetSigningConfig()
	if err != nil {
		return fmt.Errorf("failed to get signing configuration: %w", err)
	}

	// Connect to Ethereum client
	ethClient, err := ConnectToEthereum(rpcURL)
	if err != nil {
		return fmt.Errorf("failed to connect to Ethereum: %w", err)
	}
	defer ethClient.Close()

	// Create contract instance
	contractAddr := common.HexToAddress(validatorManagerContractAddr)
	contract, err := bindings.NewIValidatorManager(contractAddr, ethClient)
	if err != nil {
		return fmt.Errorf("failed to create contract instance: %w", err)
	}

	// Get the contract fee
	fee, err := contract.Fee(nil)
	if err != nil {
		return fmt.Errorf("failed to get contract fee: %w", err)
	}

	// Calculate total value to send (collateral amount + fee)
	totalValue := new(big.Int).Add(collateralAmount, fee)

	// Show summary
	fmt.Println("===== Deposit Collateral Transaction =====")
	fmt.Printf("Validator Address        : %s\n", valAddr.Hex())
	fmt.Printf("Collateral Amount        : %s MITO\n", FormatWeiToEther(collateralAmount))
	fmt.Printf("Fee                      : %s MITO\n", FormatWeiToEther(fee))
	fmt.Printf("Total Value              : %s MITO\n", FormatWeiToEther(totalValue))
	fmt.Println()

	fmt.Println("🚨 IMPORTANT: Only permitted collateral owners can deposit collateral for a validator.")
	fmt.Println("If your address is not a permitted collateral owner for this validator, the transaction will fail.")
	fmt.Println()

	// Ask for confirmation
	if !ConfirmAction("Do you want to deposit this collateral?") {
		return fmt.Errorf("operation cancelled by user")
	}

	// Create transaction options
	transactOpts, err := CreateTransactOpts(ethClient, signingConfig, totalValue)
	if err != nil {
		return fmt.Errorf("failed to create transaction options: %w", err)
	}

	// Execute the transaction
	tx, err := contract.DepositCollateral(transactOpts, valAddr)
	if err != nil {
		return fmt.Errorf("failed to deposit collateral: %w", err)
	}

	fmt.Printf("Transaction sent! Hash: %s\n", tx.Hash().Hex())

	// Wait for transaction confirmation
	err = WaitForTxConfirmation(ethClient, tx.Hash())
	if err != nil {
		return fmt.Errorf("transaction failed: %w", err)
	}

	fmt.Println("✅ Collateral deposited successfully!")
	return nil
}

// newTxSendCollateralWithdrawCmd creates the collateral withdraw command for tx send
func newTxSendCollateralWithdrawCmd() *cobra.Command {
	var (
		validator string
		amount    string
		receiver  string
	)

	cmd := &cobra.Command{
		Use:   "withdraw",
		Short: "Send a collateral withdraw transaction",
		Long: `Create, sign, and send a transaction to withdraw collateral from a validator.
Only the collateral owner can withdraw collateral.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runSendWithdrawCollateralTx(cmd, validator, amount, receiver)
		},
	}

	// Add common flags for contract interaction
	cmd.Flags().StringVar(&rpcURL, "rpc", "", "RPC endpoint URL")
	cmd.Flags().StringVar(&validatorManagerContractAddr, "contract", "", "ValidatorManager contract address")
	mustMarkFlagRequired(cmd, "rpc")
	mustMarkFlagRequired(cmd, "contract")

	// Add signing flags
	AddSigningFlags(cmd)

	// Command-specific flags
	cmd.Flags().StringVar(&validator, "validator", "", "Validator address")
	cmd.Flags().StringVar(&receiver, "receiver", "", "Address to receive the withdrawn collateral")
	cmd.Flags().StringVar(&amount, "amount", "", "Amount of collateral to withdraw in MITO (e.g., \"1.5\")")

	// Mark required flags
	mustMarkFlagRequired(cmd, "validator")
	mustMarkFlagRequired(cmd, "receiver")
	mustMarkFlagRequired(cmd, "amount")

	return cmd
}

func runSendWithdrawCollateralTx(cmd *cobra.Command, validator, amount, receiver string) error {
	// Parse collateral amount as decimal MITO and convert to wei
	collateralAmount, err := ParseValueAsWei(amount)
	if err != nil {
		return fmt.Errorf("failed to parse amount: %w", err)
	}

	// Ensure collateral amount is greater than 0
	if collateralAmount.Cmp(big.NewInt(0)) <= 0 {
		return fmt.Errorf("collateral amount must be greater than 0")
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

	// Get signing configuration
	signingConfig, err := GetSigningConfig()
	if err != nil {
		return fmt.Errorf("failed to get signing configuration: %w", err)
	}

	// Connect to Ethereum client
	ethClient, err := ConnectToEthereum(rpcURL)
	if err != nil {
		return fmt.Errorf("failed to connect to Ethereum: %w", err)
	}
	defer ethClient.Close()

	// Create contract instance
	contractAddr := common.HexToAddress(validatorManagerContractAddr)
	contract, err := bindings.NewIValidatorManager(contractAddr, ethClient)
	if err != nil {
		return fmt.Errorf("failed to create contract instance: %w", err)
	}

	// Get the contract fee
	fee, err := contract.Fee(nil)
	if err != nil {
		return fmt.Errorf("failed to get contract fee: %w", err)
	}

	// Show summary
	fmt.Println("===== Withdraw Collateral Transaction =====")
	fmt.Printf("Validator Address        : %s\n", valAddr.Hex())
	fmt.Printf("Collateral Amount        : %s MITO\n", FormatWeiToEther(collateralAmount))
	fmt.Printf("Fee                      : %s MITO\n", FormatWeiToEther(fee))
	fmt.Printf("Receiver Address         : %s\n", receiverAddr.Hex())
	fmt.Println()

	fmt.Println("🚨 IMPORTANT: Only the collateral owner can withdraw collateral.")
	fmt.Println("Make sure you are the owner of the collateral for this validator.")
	fmt.Println()

	// Ask for confirmation
	if !ConfirmAction("Do you want to withdraw this collateral?") {
		return fmt.Errorf("operation cancelled by user")
	}

	// Create transaction options (withdraw only sends fee, not collateral)
	transactOpts, err := CreateTransactOpts(ethClient, signingConfig, fee)
	if err != nil {
		return fmt.Errorf("failed to create transaction options: %w", err)
	}

	// Execute the transaction
	tx, err := contract.WithdrawCollateral(transactOpts, valAddr, receiverAddr, collateralAmount)
	if err != nil {
		return fmt.Errorf("failed to withdraw collateral: %w", err)
	}

	fmt.Printf("Transaction sent! Hash: %s\n", tx.Hash().Hex())

	// Wait for transaction confirmation
	err = WaitForTxConfirmation(ethClient, tx.Hash())
	if err != nil {
		return fmt.Errorf("transaction failed: %w", err)
	}

	fmt.Println("✅ Collateral withdrawn successfully!")
	return nil
}
