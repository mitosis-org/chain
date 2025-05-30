package tx

import (
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/mitosis-org/chain/bindings"
	"github.com/mitosis-org/chain/cmd/mito/internal/config"
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

// CreateValidator creates an unsigned transaction for creating a validator
func (s *ValidatorService) CreateValidator(req *CreateValidatorRequest) (*types.Transaction, error) {
	// // Get the contract fee
	// fee, err := s.contract.Fee(nil)
	// if err != nil {
	// 	return nil, fmt.Errorf("failed to get contract fee: %w", err)
	// }

	// // Get the config to check the initial deposit requirement
	// config, err := s.contract.GlobalValidatorConfig(nil)
	// if err != nil {
	// 	return nil, fmt.Errorf("failed to get global validator config: %w", err)
	// }

	// // Parse collateral amount as decimal MITO and convert to wei
	// collateralAmount, err := utils.ParseValueAsWei(req.InitialCollateral)
	// if err != nil {
	// 	return nil, fmt.Errorf("failed to parse initial collateral: %w", err)
	// }

	// // Ensure collateral is at least the initial deposit requirement
	// if collateralAmount.Cmp(config.InitialValidatorDeposit) < 0 {
	// 	return nil, fmt.Errorf("initial collateral must be at least %s MITO",
	// 		utils.FormatWeiToEther(config.InitialValidatorDeposit))
	// }

	// // Calculate total transaction value (collateral + fee)
	// totalValue := new(big.Int).Add(collateralAmount, fee)

	// // Validate other parameters
	// operatorAddr, err := utils.ValidateAddress(req.Operator)
	// if err != nil {
	// 	return nil, fmt.Errorf("invalid operator address: %w", err)
	// }

	// rewardManagerAddr, err := utils.ValidateAddress(req.RewardManager)
	// if err != nil {
	// 	return nil, fmt.Errorf("invalid reward manager address: %w", err)
	// }

	// // Parse commission rate
	// commissionRateInt, err := utils.ParsePercentageToBasisPoints(req.CommissionRate)
	// if err != nil {
	// 	return nil, fmt.Errorf("failed to parse commission rate: %w", err)
	// }

	// // Validate commission rate
	// maxRate, err := s.contract.MAXCOMMISSIONRATE(nil)
	// if err != nil {
	// 	return nil, fmt.Errorf("failed to get max commission rate: %w", err)
	// }

	// if commissionRateInt.Cmp(big.NewInt(0)) < 0 || commissionRateInt.Cmp(maxRate) > 0 {
	// 	return nil, fmt.Errorf("commission rate must be between 0%% and %s", utils.FormatBasisPointsToPercent(maxRate))
	// }

	// // Decode public key from hex
	// pubKeyBytes, err := utils.DecodeHexWithPrefix(req.PubKey)
	// if err != nil {
	// 	return nil, fmt.Errorf("failed to decode public key: %w", err)
	// }

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

	// Set default gas limit if not provided
	gasLimit := s.config.GasLimit
	if gasLimit == 0 {
		gasLimit = 500000 // Default gas limit
	}

	collateralAmount, err := utils.ParseValueAsWei(req.InitialCollateral)
	if err != nil {
		return nil, fmt.Errorf("failed to parse initial collateral: %w", err)
	}

	// // Show summary
	// fmt.Println("===== Create Validator Transaction =====")
	// fmt.Printf("Public Key                 : %s\n", req.PubKey)
	// fmt.Printf("Operator                   : %s\n", operatorAddr.Hex())
	// fmt.Printf("Reward Manager             : %s\n", rewardManagerAddr.Hex())
	// fmt.Printf("Commission Rate            : %s\n", utils.FormatBasisPointsToPercent(commissionRateInt))
	// fmt.Printf("Metadata                   : %s\n", req.Metadata)
	// fmt.Printf("Initial Collateral         : %s MITO\n", utils.FormatWeiToEther(collateralAmount))
	// fmt.Println()

	// Create transaction data
	txData := &TransactionData{
		To:       common.HexToAddress(s.config.ValidatorManagerContractAddr),
		Value:    collateralAmount,
		Data:     data,
		GasLimit: gasLimit,
	}

	// Create transaction
	return s.builder.CreateTransactionFromData(txData)
}

// UpdateMetadata creates an unsigned transaction for updating metadata
func (s *ValidatorService) UpdateMetadata(validatorAddr, metadata string) (*types.Transaction, error) {
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

	// Set default gas limit if not provided
	gasLimit := s.config.GasLimit
	if gasLimit == 0 {
		gasLimit = 500000
	}

	// Create transaction data
	txData := &TransactionData{
		To:       common.HexToAddress(s.config.ValidatorManagerContractAddr),
		Value:    big.NewInt(0),
		Data:     data,
		GasLimit: gasLimit,
	}

	// Create transaction
	return s.builder.CreateTransactionFromData(txData)
}

// UpdateOperator creates an unsigned transaction for updating operator
func (s *ValidatorService) UpdateOperator(validatorAddr, newOperator string) (*types.Transaction, error) {
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

	// Set default gas limit if not provided
	gasLimit := s.config.GasLimit
	if gasLimit == 0 {
		gasLimit = 500000
	}

	// Create transaction data
	txData := &TransactionData{
		To:       common.HexToAddress(s.config.ValidatorManagerContractAddr),
		Value:    big.NewInt(0),
		Data:     data,
		GasLimit: gasLimit,
	}

	// Create transaction
	return s.builder.CreateTransactionFromDataWithOptions(txData, true)
}

// UpdateRewardConfig creates an unsigned transaction for updating reward config
func (s *ValidatorService) UpdateRewardConfig(validatorAddr, commissionRate string) (*types.Transaction, error) {
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

	// Set default gas limit if not provided
	gasLimit := s.config.GasLimit
	if gasLimit == 0 {
		gasLimit = 500000
	}

	// Create transaction data
	txData := &TransactionData{
		To:       common.HexToAddress(s.config.ValidatorManagerContractAddr),
		Value:    big.NewInt(0),
		Data:     data,
		GasLimit: gasLimit,
	}

	// Create transaction
	return s.builder.CreateTransactionFromDataWithOptions(txData, true)
}

// UpdateRewardManager creates an unsigned transaction for updating reward manager
func (s *ValidatorService) UpdateRewardManager(validatorAddr, rewardManager string) (*types.Transaction, error) {
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

	// Set default gas limit if not provided
	gasLimit := s.config.GasLimit
	if gasLimit == 0 {
		gasLimit = 500000
	}

	// Create transaction data
	txData := &TransactionData{
		To:       common.HexToAddress(s.config.ValidatorManagerContractAddr),
		Value:    big.NewInt(0),
		Data:     data,
		GasLimit: gasLimit,
	}

	// Create transaction
	return s.builder.CreateTransactionFromDataWithOptions(txData, true)
}

// UnjailValidator creates an unsigned transaction for unjailing a validator
func (s *ValidatorService) UnjailValidator(validatorAddr string) (*types.Transaction, error) {
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

	// Set default gas limit if not provided
	gasLimit := s.config.GasLimit
	if gasLimit == 0 {
		gasLimit = 500000
	}

	// Create transaction data
	txData := &TransactionData{
		To:       common.HexToAddress(s.config.ValidatorManagerContractAddr),
		Value:    big.NewInt(0),
		Data:     data,
		GasLimit: gasLimit,
	}

	// Create transaction
	return s.builder.CreateTransactionFromDataWithOptions(txData, true)
}
