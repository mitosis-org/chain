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

// ValidatorService handles validator-related transactions
type ValidatorService struct {
	config  *config.ResolvedConfig
	builder *Builder
}

// NewValidatorService creates a new validator service
func NewValidatorService(config *config.ResolvedConfig, builder *Builder) *ValidatorService {
	return &ValidatorService{
		config:  config,
		builder: builder,
	}
}

// CreateValidatorRequest contains parameters for creating a validator
type CreateValidatorRequest struct {
	PubKey            string
	Operator          string
	RewardManager     string
	CommissionRate    string
	Metadata          string
	InitialCollateral string
}

// TransactionData contains the data needed to build a transaction
type TransactionData struct {
	To       common.Address
	Value    *big.Int
	Data     []byte
	GasLimit uint64
}

// CreateValidator creates a transaction for creating a validator
func (s *ValidatorService) CreateValidator(req *CreateValidatorRequest) (*types.Transaction, error) {
	return s.CreateValidatorWithOptions(req, false)
}

// CreateValidatorWithOptions creates a transaction for creating a validator with options
func (s *ValidatorService) CreateValidatorWithOptions(req *CreateValidatorRequest, unsigned bool) (*types.Transaction, error) {
	// Validate and parse inputs
	pubKeyBytes, err := utils.DecodeHexWithPrefix(req.PubKey)
	if err != nil {
		return nil, fmt.Errorf("failed to decode public key: %w", err)
	}

	operatorAddr, err := utils.ValidateAddress(req.Operator)
	if err != nil {
		return nil, fmt.Errorf("invalid operator address: %w", err)
	}

	rewardManagerAddr, err := utils.ValidateAddress(req.RewardManager)
	if err != nil {
		return nil, fmt.Errorf("invalid reward manager address: %w", err)
	}

	commissionRateInt, err := utils.ParsePercentageToBasisPoints(req.CommissionRate)
	if err != nil {
		return nil, fmt.Errorf("failed to parse commission rate: %w", err)
	}

	// Parse initial collateral using new units module
	collateralAmount, err := units.ParseMitoToWei(req.InitialCollateral)
	if err != nil {
		return nil, fmt.Errorf("failed to parse initial collateral: %w", err)
	}

	// Calculate total value to send (collateral amount + fee)
	totalValue := new(big.Int).Add(collateralAmount, s.config.ContractFee)

	// Create the request
	request := bindings.IValidatorManagerCreateValidatorRequest{
		Operator:       operatorAddr,
		RewardManager:  rewardManagerAddr,
		CommissionRate: commissionRateInt,
		Metadata:       []byte(req.Metadata),
	}

	// Get contract ABI and encode function call data
	abi, err := bindings.IValidatorManagerMetaData.GetAbi()
	if err != nil {
		return nil, fmt.Errorf("failed to get contract ABI: %w", err)
	}

	data, err := abi.Pack("createValidator", pubKeyBytes, request)
	if err != nil {
		return nil, fmt.Errorf("failed to pack function call: %w", err)
	}

	// Create transaction data
	txData := &TransactionData{
		To:       common.HexToAddress(s.config.ValidatorManagerContractAddr),
		Value:    totalValue,
		Data:     data,
		GasLimit: 0, // Let builder handle gas estimation
	}

	// Create transaction
	return s.builder.CreateTransactionFromDataWithOptions(txData, unsigned)
}

// UpdateMetadata creates a transaction for updating metadata
func (s *ValidatorService) UpdateMetadata(validatorAddr, metadata string) (*types.Transaction, error) {
	return s.UpdateMetadataWithOptions(validatorAddr, metadata, false)
}

// UpdateMetadataWithOptions creates a transaction for updating metadata with options
func (s *ValidatorService) UpdateMetadataWithOptions(validatorAddr, metadata string, unsigned bool) (*types.Transaction, error) {
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

	data, err := abi.Pack("updateMetadata", valAddr, []byte(metadata))
	if err != nil {
		return nil, fmt.Errorf("failed to pack function call: %w", err)
	}

	// Create transaction data (no value needed for metadata update)
	txData := &TransactionData{
		To:       common.HexToAddress(s.config.ValidatorManagerContractAddr),
		Value:    big.NewInt(0),
		Data:     data,
		GasLimit: 0, // Let builder handle gas estimation
	}

	// Create transaction
	return s.builder.CreateTransactionFromDataWithOptions(txData, unsigned)
}

// UpdateOperator creates a transaction for updating operator
func (s *ValidatorService) UpdateOperator(validatorAddr, newOperator string) (*types.Transaction, error) {
	return s.UpdateOperatorWithOptions(validatorAddr, newOperator, false)
}

// UpdateOperatorWithOptions creates a transaction for updating operator with options
func (s *ValidatorService) UpdateOperatorWithOptions(validatorAddr, newOperator string, unsigned bool) (*types.Transaction, error) {
	// Validate addresses
	valAddr, err := utils.ValidateAddress(validatorAddr)
	if err != nil {
		return nil, fmt.Errorf("invalid validator address: %w", err)
	}

	operatorAddr, err := utils.ValidateAddress(newOperator)
	if err != nil {
		return nil, fmt.Errorf("invalid operator address: %w", err)
	}

	// Get contract ABI and encode function call data
	abi, err := bindings.IValidatorManagerMetaData.GetAbi()
	if err != nil {
		return nil, fmt.Errorf("failed to get contract ABI: %w", err)
	}

	data, err := abi.Pack("updateOperator", valAddr, operatorAddr)
	if err != nil {
		return nil, fmt.Errorf("failed to pack function call: %w", err)
	}

	// Create transaction data (no value needed for operator update)
	txData := &TransactionData{
		To:       common.HexToAddress(s.config.ValidatorManagerContractAddr),
		Value:    big.NewInt(0),
		Data:     data,
		GasLimit: 0, // Let builder handle gas estimation
	}

	// Create transaction
	return s.builder.CreateTransactionFromDataWithOptions(txData, unsigned)
}

// UpdateRewardConfig creates a transaction for updating reward config
func (s *ValidatorService) UpdateRewardConfig(validatorAddr, commissionRate string) (*types.Transaction, error) {
	return s.UpdateRewardConfigWithOptions(validatorAddr, commissionRate, false)
}

// UpdateRewardConfigWithOptions creates a transaction for updating reward config with options
func (s *ValidatorService) UpdateRewardConfigWithOptions(validatorAddr, commissionRate string, unsigned bool) (*types.Transaction, error) {
	// Validate validator address
	valAddr, err := utils.ValidateAddress(validatorAddr)
	if err != nil {
		return nil, fmt.Errorf("invalid validator address: %w", err)
	}

	// Parse commission rate
	commissionRateInt, err := utils.ParsePercentageToBasisPoints(commissionRate)
	if err != nil {
		return nil, fmt.Errorf("failed to parse commission rate: %w", err)
	}

	// Create the request struct
	request := bindings.IValidatorManagerUpdateRewardConfigRequest{
		CommissionRate: commissionRateInt,
	}

	// Get contract ABI and encode function call data
	abi, err := bindings.IValidatorManagerMetaData.GetAbi()
	if err != nil {
		return nil, fmt.Errorf("failed to get contract ABI: %w", err)
	}

	data, err := abi.Pack("updateRewardConfig", valAddr, request)
	if err != nil {
		return nil, fmt.Errorf("failed to pack function call: %w", err)
	}

	// Create transaction data (no value needed for reward config update)
	txData := &TransactionData{
		To:       common.HexToAddress(s.config.ValidatorManagerContractAddr),
		Value:    big.NewInt(0),
		Data:     data,
		GasLimit: 0, // Let builder handle gas estimation
	}

	// Create transaction
	return s.builder.CreateTransactionFromDataWithOptions(txData, unsigned)
}

// UpdateRewardManager creates a transaction for updating reward manager
func (s *ValidatorService) UpdateRewardManager(validatorAddr, rewardManager string) (*types.Transaction, error) {
	return s.UpdateRewardManagerWithOptions(validatorAddr, rewardManager, false)
}

// UpdateRewardManagerWithOptions creates a transaction for updating reward manager with options
func (s *ValidatorService) UpdateRewardManagerWithOptions(validatorAddr, rewardManager string, unsigned bool) (*types.Transaction, error) {
	// Validate addresses
	valAddr, err := utils.ValidateAddress(validatorAddr)
	if err != nil {
		return nil, fmt.Errorf("invalid validator address: %w", err)
	}

	rewardManagerAddr, err := utils.ValidateAddress(rewardManager)
	if err != nil {
		return nil, fmt.Errorf("invalid reward manager address: %w", err)
	}

	// Get contract ABI and encode function call data
	abi, err := bindings.IValidatorManagerMetaData.GetAbi()
	if err != nil {
		return nil, fmt.Errorf("failed to get contract ABI: %w", err)
	}

	data, err := abi.Pack("updateRewardManager", valAddr, rewardManagerAddr)
	if err != nil {
		return nil, fmt.Errorf("failed to pack function call: %w", err)
	}

	// Create transaction data (no value needed for reward manager update)
	txData := &TransactionData{
		To:       common.HexToAddress(s.config.ValidatorManagerContractAddr),
		Value:    big.NewInt(0),
		Data:     data,
		GasLimit: 0, // Let builder handle gas estimation
	}

	// Create transaction
	return s.builder.CreateTransactionFromDataWithOptions(txData, unsigned)
}

// UnjailValidator creates a transaction for unjailing a validator
func (s *ValidatorService) UnjailValidator(validatorAddr string) (*types.Transaction, error) {
	return s.UnjailValidatorWithOptions(validatorAddr, false)
}

// UnjailValidatorWithOptions creates a transaction for unjailing a validator with options
func (s *ValidatorService) UnjailValidatorWithOptions(validatorAddr string, unsigned bool) (*types.Transaction, error) {
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

	data, err := abi.Pack("unjailValidator", valAddr)
	if err != nil {
		return nil, fmt.Errorf("failed to pack function call: %w", err)
	}

	// Create transaction data
	txData := &TransactionData{
		To:       common.HexToAddress(s.config.ValidatorManagerContractAddr),
		Value:    s.config.ContractFee,
		Data:     data,
		GasLimit: 0, // Let builder handle gas estimation
	}

	// Create transaction
	return s.builder.CreateTransactionFromDataWithOptions(txData, unsigned)
}
