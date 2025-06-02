package tx

import (
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/mitosis-org/chain/bindings"
	"github.com/mitosis-org/chain/cmd/mito/internal/config"
	"github.com/mitosis-org/chain/cmd/mito/internal/units"
	"github.com/mitosis-org/chain/cmd/mito/internal/utils"
)

// CollateralService handles collateral-related transactions
type CollateralService struct {
	config  *config.ResolvedConfig
	builder *Builder
}

// NewCollateralService creates a new collateral service
func NewCollateralService(config *config.ResolvedConfig, builder *Builder) *CollateralService {
	return &CollateralService{
		config:  config,
		builder: builder,
	}
}

// DepositCollateral creates an unsigned transaction for depositing collateral
func (s *CollateralService) DepositCollateral(validatorAddr, amount string) (*types.Transaction, error) {
	// Parse collateral amount as decimal MITO and convert to wei
	collateralAmount, err := utils.ParseValueAsWei(amount)
	if err != nil {
		return nil, fmt.Errorf("failed to parse amount: %w", err)
	}

	// Ensure collateral amount is greater than 0
	if collateralAmount.Cmp(big.NewInt(0)) <= 0 {
		return nil, fmt.Errorf("collateral amount must be greater than 0")
	}

	// Get contract fee
	fee, err := units.ParseContractFeeInput(s.config.ContractFee)
	if err != nil {
		return nil, fmt.Errorf("failed to get contract fee: %w", err)
	}

	// Calculate total value to send (collateral amount + fee)
	totalValue := new(big.Int).Add(collateralAmount, fee)

	// Validate validator address
	valAddr, err := utils.ValidateAddress(validatorAddr)
	if err != nil {
		return nil, fmt.Errorf("invalid validator address: %w", err)
	}

	// Get contract ABI and encode function call data
	abi, err := bindings.IValidatorManagerMetaData.GetAbi()
	if err != nil {
		return nil, fmt.Errorf("failed to get contract ABI: %w", err)
	}

	data, err := abi.Pack("depositCollateral", valAddr)
	if err != nil {
		return nil, fmt.Errorf("failed to pack function call: %w", err)
	}

	// Set default gas limit if not provided
	gasLimit := s.config.GasLimit
	if gasLimit == 0 {
		gasLimit = 500000
	}

	// Create transaction data
	txData := &TransactionData{
		To:       common.HexToAddress(s.config.ValidatorManagerContractAddr),
		Value:    totalValue,
		Data:     data,
		GasLimit: gasLimit,
	}

	// Create transaction
	return s.builder.CreateTransactionFromDataWithOptions(txData, true)
}

// WithdrawCollateral creates an unsigned transaction for withdrawing collateral
func (s *CollateralService) WithdrawCollateral(validatorAddr, amount, receiver string) (*types.Transaction, error) {
	// Parse collateral amount as decimal MITO and convert to wei
	collateralAmount, err := utils.ParseValueAsWei(amount)
	if err != nil {
		return nil, fmt.Errorf("failed to parse amount: %w", err)
	}

	// Ensure collateral amount is greater than 0
	if collateralAmount.Cmp(big.NewInt(0)) <= 0 {
		return nil, fmt.Errorf("collateral amount must be greater than 0")
	}

	// Get contract fee
	fee, err := units.ParseContractFeeInput(s.config.ContractFee)
	if err != nil {
		return nil, fmt.Errorf("failed to get contract fee: %w", err)
	}

	// Validate addresses
	valAddr, err := utils.ValidateAddress(validatorAddr)
	if err != nil {
		return nil, fmt.Errorf("invalid validator address: %w", err)
	}

	receiverAddr, err := utils.ValidateAddress(receiver)
	if err != nil {
		return nil, fmt.Errorf("invalid receiver address: %w", err)
	}

	// Get contract ABI and encode function call data
	abi, err := bindings.IValidatorManagerMetaData.GetAbi()
	if err != nil {
		return nil, fmt.Errorf("failed to get contract ABI: %w", err)
	}

	data, err := abi.Pack("withdrawCollateral", valAddr, receiverAddr, collateralAmount)
	if err != nil {
		return nil, fmt.Errorf("failed to pack function call: %w", err)
	}

	// Set default gas limit if not provided
	gasLimit := s.config.GasLimit
	if gasLimit == 0 {
		gasLimit = 500000
	}

	// Create transaction data (withdraw only sends fee, not collateral)
	txData := &TransactionData{
		To:       common.HexToAddress(s.config.ValidatorManagerContractAddr),
		Value:    fee,
		Data:     data,
		GasLimit: gasLimit,
	}

	// Create transaction
	return s.builder.CreateTransactionFromDataWithOptions(txData, true)
}

// SetPermittedCollateralOwner creates an unsigned transaction for setting a permitted collateral owner
func (s *CollateralService) SetPermittedCollateralOwner(validatorAddr, collateralOwner string, isPermitted bool) (*types.Transaction, error) {
	// Validate addresses
	valAddr, err := utils.ValidateAddress(validatorAddr)
	if err != nil {
		return nil, fmt.Errorf("invalid validator address: %w", err)
	}

	collateralOwnerAddr, err := utils.ValidateAddress(collateralOwner)
	if err != nil {
		return nil, fmt.Errorf("invalid collateral owner address: %w", err)
	}

	// Get contract ABI and encode function call data
	abi, err := bindings.IValidatorManagerMetaData.GetAbi()
	if err != nil {
		return nil, fmt.Errorf("failed to get contract ABI: %w", err)
	}

	data, err := abi.Pack("setPermittedCollateralOwner", valAddr, collateralOwnerAddr, isPermitted)
	if err != nil {
		return nil, fmt.Errorf("failed to pack function call: %w", err)
	}

	// Set default gas limit if not provided
	gasLimit := s.config.GasLimit
	if gasLimit == 0 {
		gasLimit = 500000
	}

	// Create transaction data (no value needed for permission update)
	txData := &TransactionData{
		To:       common.HexToAddress(s.config.ValidatorManagerContractAddr),
		Value:    big.NewInt(0),
		Data:     data,
		GasLimit: gasLimit,
	}

	// Create transaction
	return s.builder.CreateTransactionFromDataWithOptions(txData, true)
}

// TransferCollateralOwnership creates an unsigned transaction for transferring collateral ownership
func (s *CollateralService) TransferCollateralOwnership(validatorAddr, newOwner string) (*types.Transaction, error) {
	// Get contract fee
	fee, err := units.ParseContractFeeInput(s.config.ContractFee)
	if err != nil {
		return nil, fmt.Errorf("failed to get contract fee: %w", err)
	}

	// Validate addresses
	valAddr, err := utils.ValidateAddress(validatorAddr)
	if err != nil {
		return nil, fmt.Errorf("invalid validator address: %w", err)
	}

	newOwnerAddr, err := utils.ValidateAddress(newOwner)
	if err != nil {
		return nil, fmt.Errorf("invalid new owner address: %w", err)
	}

	// Get contract ABI and encode function call data
	abi, err := bindings.IValidatorManagerMetaData.GetAbi()
	if err != nil {
		return nil, fmt.Errorf("failed to get contract ABI: %w", err)
	}

	data, err := abi.Pack("transferCollateralOwnership", valAddr, newOwnerAddr)
	if err != nil {
		return nil, fmt.Errorf("failed to pack function call: %w", err)
	}

	// Set default gas limit if not provided
	gasLimit := s.config.GasLimit
	if gasLimit == 0 {
		gasLimit = 500000
	}

	// Create transaction data
	txData := &TransactionData{
		To:       common.HexToAddress(s.config.ValidatorManagerContractAddr),
		Value:    fee,
		Data:     data,
		GasLimit: gasLimit,
	}

	// Create transaction
	return s.builder.CreateTransactionFromDataWithOptions(txData, true)
}
