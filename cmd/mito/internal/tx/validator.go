package tx

import (
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/mitosis-org/chain/bindings"
	"github.com/mitosis-org/chain/cmd/mito/internal/client"
	"github.com/mitosis-org/chain/cmd/mito/internal/config"
	"github.com/mitosis-org/chain/cmd/mito/internal/utils"
)

// ValidatorService handles validator-related transactions
type ValidatorService struct {
	config    *config.ResolvedConfig
	ethClient *client.EthereumClient
	contract  *client.ValidatorManagerContract
	builder   *Builder
}

// NewValidatorService creates a new validator service
func NewValidatorService(config *config.ResolvedConfig, ethClient *client.EthereumClient, contract *client.ValidatorManagerContract, builder *Builder) *ValidatorService {
	return &ValidatorService{
		config:    config,
		ethClient: ethClient,
		contract:  contract,
		builder:   builder,
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

// CreateValidator creates a new validator
func (s *ValidatorService) CreateValidator(req *CreateValidatorRequest) (common.Hash, error) {
	// Get the contract fee
	fee, err := s.contract.Fee(nil)
	if err != nil {
		return common.Hash{}, fmt.Errorf("failed to get contract fee: %w", err)
	}

	// Get the config to check the initial deposit requirement
	config, err := s.contract.GlobalValidatorConfig(nil)
	if err != nil {
		return common.Hash{}, fmt.Errorf("failed to get global validator config: %w", err)
	}

	// Parse collateral amount as decimal MITO and convert to wei
	collateralAmount, err := utils.ParseValueAsWei(req.InitialCollateral)
	if err != nil {
		return common.Hash{}, fmt.Errorf("failed to parse initial collateral: %w", err)
	}

	// Ensure collateral is at least the initial deposit requirement
	if collateralAmount.Cmp(config.InitialValidatorDeposit) < 0 {
		return common.Hash{}, fmt.Errorf("initial collateral must be at least %s MITO",
			utils.FormatWeiToEther(config.InitialValidatorDeposit))
	}

	// Calculate total transaction value (collateral + fee)
	totalValue := new(big.Int).Add(collateralAmount, fee)

	// Validate other parameters
	operatorAddr, err := utils.ValidateAddress(req.Operator)
	if err != nil {
		return common.Hash{}, fmt.Errorf("invalid operator address: %w", err)
	}

	rewardManagerAddr, err := utils.ValidateAddress(req.RewardManager)
	if err != nil {
		return common.Hash{}, fmt.Errorf("invalid reward manager address: %w", err)
	}

	// Parse commission rate
	commissionRateInt, err := utils.ParsePercentageToBasisPoints(req.CommissionRate)
	if err != nil {
		return common.Hash{}, fmt.Errorf("failed to parse commission rate: %w", err)
	}

	// Validate commission rate
	maxRate, err := s.contract.MAXCOMMISSIONRATE(nil)
	if err != nil {
		return common.Hash{}, fmt.Errorf("failed to get max commission rate: %w", err)
	}

	if commissionRateInt.Cmp(big.NewInt(0)) < 0 || commissionRateInt.Cmp(maxRate) > 0 {
		return common.Hash{}, fmt.Errorf("commission rate must be between 0%% and %s", utils.FormatBasisPointsToPercent(maxRate))
	}

	// Decode public key from hex
	pubKeyBytes, err := utils.DecodeHexWithPrefix(req.PubKey)
	if err != nil {
		return common.Hash{}, fmt.Errorf("failed to decode public key: %w", err)
	}

	// Create the request
	request := bindings.IValidatorManagerCreateValidatorRequest{
		Operator:       operatorAddr,
		RewardManager:  rewardManagerAddr,
		CommissionRate: commissionRateInt,
		Metadata:       []byte(req.Metadata),
	}

	// Show summary
	fmt.Println("===== Create Validator Transaction =====")
	fmt.Printf("Public Key                 : %s\n", req.PubKey)
	fmt.Printf("Operator                   : %s\n", operatorAddr.Hex())
	fmt.Printf("Reward Manager             : %s\n", rewardManagerAddr.Hex())
	fmt.Printf("Commission Rate            : %s\n", utils.FormatBasisPointsToPercent(commissionRateInt))
	fmt.Printf("Metadata                   : %s\n", req.Metadata)
	fmt.Printf("Initial Collateral         : %s MITO\n", utils.FormatWeiToEther(collateralAmount))
	fmt.Printf("Fee                        : %s MITO\n", utils.FormatWeiToEther(fee))
	fmt.Printf("Total Value                : %s MITO\n", utils.FormatWeiToEther(totalValue))
	fmt.Println()

	// Execute the transaction
	opts, err := s.builder.CreateTransactOpts(totalValue)
	if err != nil {
		return common.Hash{}, fmt.Errorf("failed to create transaction options: %w", err)
	}

	tx, err := s.contract.CreateValidator(opts, pubKeyBytes, request)
	if err != nil {
		return common.Hash{}, fmt.Errorf("failed to create validator: %w", err)
	}

	return tx.Hash(), nil
}

// UpdateOperator updates the operator address for a validator
func (s *ValidatorService) UpdateOperator(validatorAddr, newOperator string) (common.Hash, error) {
	// Validate addresses
	valAddr, err := utils.ValidateAddress(validatorAddr)
	if err != nil {
		return common.Hash{}, fmt.Errorf("invalid validator address: %w", err)
	}

	operatorAddr, err := utils.ValidateAddress(newOperator)
	if err != nil {
		return common.Hash{}, fmt.Errorf("invalid operator address: %w", err)
	}

	// Get validator info to show current values
	validatorInfo, err := s.contract.ValidatorInfo(nil, valAddr)
	if err != nil {
		return common.Hash{}, fmt.Errorf("failed to get validator info: %w", err)
	}

	// Show summary
	fmt.Println("===== Update Operator Transaction =====")
	fmt.Printf("Validator Address            : %s\n", valAddr.Hex())
	fmt.Printf("Current Operator             : %s\n", validatorInfo.Operator.Hex())
	fmt.Printf("New Operator                 : %s\n", operatorAddr.Hex())
	fmt.Printf("Current Reward Manager       : %s\n", validatorInfo.RewardManager.Hex())
	fmt.Println()

	// Execute the transaction
	opts, err := s.builder.CreateTransactOpts(nil)
	if err != nil {
		return common.Hash{}, fmt.Errorf("failed to create transaction options: %w", err)
	}

	tx, err := s.contract.UpdateOperator(opts, valAddr, operatorAddr)
	if err != nil {
		return common.Hash{}, fmt.Errorf("failed to update operator: %w", err)
	}

	return tx.Hash(), nil
}

// UpdateMetadata updates the metadata for a validator
func (s *ValidatorService) UpdateMetadata(validatorAddr, metadata string) (common.Hash, error) {
	// Validate validator address
	valAddr, err := utils.ValidateAddress(validatorAddr)
	if err != nil {
		return common.Hash{}, fmt.Errorf("invalid validator address: %w", err)
	}

	// Get validator info to show current values
	validatorInfo, err := s.contract.ValidatorInfo(nil, valAddr)
	if err != nil {
		return common.Hash{}, fmt.Errorf("failed to get validator info: %w", err)
	}

	// Show summary
	fmt.Println("===== Update Metadata Transaction =====")
	fmt.Printf("Validator Address        : %s\n", valAddr.Hex())
	fmt.Printf("Current Metadata         : %s\n", string(validatorInfo.Metadata))
	fmt.Printf("New Metadata             : %s\n", metadata)
	fmt.Println()

	// Execute the transaction
	opts, err := s.builder.CreateTransactOpts(nil)
	if err != nil {
		return common.Hash{}, fmt.Errorf("failed to create transaction options: %w", err)
	}

	tx, err := s.contract.UpdateMetadata(opts, valAddr, []byte(metadata))
	if err != nil {
		return common.Hash{}, fmt.Errorf("failed to update metadata: %w", err)
	}

	return tx.Hash(), nil
}

// UpdateRewardConfig updates the reward configuration for a validator
func (s *ValidatorService) UpdateRewardConfig(validatorAddr, commissionRate string) (common.Hash, error) {
	// Validate validator address
	valAddr, err := utils.ValidateAddress(validatorAddr)
	if err != nil {
		return common.Hash{}, fmt.Errorf("invalid validator address: %w", err)
	}

	// Parse commission rate
	commissionRateInt, err := utils.ParsePercentageToBasisPoints(commissionRate)
	if err != nil {
		return common.Hash{}, fmt.Errorf("failed to parse commission rate: %w", err)
	}

	// Validate commission rate
	maxRate, err := s.contract.MAXCOMMISSIONRATE(nil)
	if err != nil {
		return common.Hash{}, fmt.Errorf("failed to get max commission rate: %w", err)
	}

	if commissionRateInt.Cmp(big.NewInt(0)) < 0 || commissionRateInt.Cmp(maxRate) > 0 {
		return common.Hash{}, fmt.Errorf("commission rate must be between 0%% and %s", utils.FormatBasisPointsToPercent(maxRate))
	}

	// Create the request struct
	request := bindings.IValidatorManagerUpdateRewardConfigRequest{
		CommissionRate: commissionRateInt,
	}

	// Show summary
	fmt.Println("===== Update Reward Config Transaction =====")
	fmt.Printf("Validator Address        : %s\n", valAddr.Hex())
	fmt.Printf("New Commission Rate      : %s\n", utils.FormatBasisPointsToPercent(commissionRateInt))
	fmt.Println()

	// Execute the transaction
	opts, err := s.builder.CreateTransactOpts(nil)
	if err != nil {
		return common.Hash{}, fmt.Errorf("failed to create transaction options: %w", err)
	}

	tx, err := s.contract.UpdateRewardConfig(opts, valAddr, request)
	if err != nil {
		return common.Hash{}, fmt.Errorf("failed to update reward config: %w", err)
	}

	return tx.Hash(), nil
}

// UpdateRewardManager updates the reward manager for a validator
func (s *ValidatorService) UpdateRewardManager(validatorAddr, rewardManager string) (common.Hash, error) {
	// Validate addresses
	valAddr, err := utils.ValidateAddress(validatorAddr)
	if err != nil {
		return common.Hash{}, fmt.Errorf("invalid validator address: %w", err)
	}

	rewardManagerAddr, err := utils.ValidateAddress(rewardManager)
	if err != nil {
		return common.Hash{}, fmt.Errorf("invalid reward manager address: %w", err)
	}

	// Get validator info to show current values
	validatorInfo, err := s.contract.ValidatorInfo(nil, valAddr)
	if err != nil {
		return common.Hash{}, fmt.Errorf("failed to get validator info: %w", err)
	}

	// Show summary
	fmt.Println("===== Update Reward Manager Transaction =====")
	fmt.Printf("Validator Address            : %s\n", valAddr.Hex())
	fmt.Printf("Current Reward Manager       : %s\n", validatorInfo.RewardManager.Hex())
	fmt.Printf("New Reward Manager           : %s\n", rewardManagerAddr.Hex())
	fmt.Println()

	// Execute the transaction
	opts, err := s.builder.CreateTransactOpts(nil)
	if err != nil {
		return common.Hash{}, fmt.Errorf("failed to create transaction options: %w", err)
	}

	tx, err := s.contract.UpdateRewardManager(opts, valAddr, rewardManagerAddr)
	if err != nil {
		return common.Hash{}, fmt.Errorf("failed to update reward manager: %w", err)
	}

	return tx.Hash(), nil
}

// UnjailValidator unjails a validator
func (s *ValidatorService) UnjailValidator(validatorAddr string) (common.Hash, error) {
	// Get the contract fee
	fee, err := s.contract.Fee(nil)
	if err != nil {
		return common.Hash{}, fmt.Errorf("failed to get contract fee: %w", err)
	}

	// Validate validator address
	valAddr, err := utils.ValidateAddress(validatorAddr)
	if err != nil {
		return common.Hash{}, fmt.Errorf("invalid validator address: %w", err)
	}

	// Show summary
	fmt.Println("===== Unjail Validator Transaction =====")
	fmt.Printf("Validator Address        : %s\n", valAddr.Hex())
	fmt.Printf("Fee                      : %s MITO\n", utils.FormatWeiToEther(fee))
	fmt.Println()

	// Execute the transaction
	opts, err := s.builder.CreateTransactOpts(fee)
	if err != nil {
		return common.Hash{}, fmt.Errorf("failed to create transaction options: %w", err)
	}

	tx, err := s.contract.UnjailValidator(opts, valAddr)
	if err != nil {
		return common.Hash{}, fmt.Errorf("failed to unjail validator: %w", err)
	}

	return tx.Hash(), nil
}
