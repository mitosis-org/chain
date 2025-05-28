package tx

import (
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/mitosis-org/chain/cmd/mito/internal/client"
	"github.com/mitosis-org/chain/cmd/mito/internal/config"
	"github.com/mitosis-org/chain/cmd/mito/internal/utils"
)

// CollateralService handles collateral-related transactions
type CollateralService struct {
	config    *config.ResolvedConfig
	ethClient *client.EthereumClient
	contract  *client.ValidatorManagerContract
	builder   *Builder
}

// NewCollateralService creates a new collateral service
func NewCollateralService(config *config.ResolvedConfig, ethClient *client.EthereumClient, contract *client.ValidatorManagerContract, builder *Builder) *CollateralService {
	return &CollateralService{
		config:    config,
		ethClient: ethClient,
		contract:  contract,
		builder:   builder,
	}
}

// DepositCollateral deposits collateral for a validator
func (s *CollateralService) DepositCollateral(validatorAddr, amount string) (common.Hash, error) {
	// Parse collateral amount as decimal MITO and convert to wei
	collateralAmount, err := utils.ParseValueAsWei(amount)
	if err != nil {
		return common.Hash{}, fmt.Errorf("failed to parse amount: %w", err)
	}

	// Ensure collateral amount is greater than 0
	if collateralAmount.Cmp(big.NewInt(0)) <= 0 {
		return common.Hash{}, fmt.Errorf("collateral amount must be greater than 0")
	}

	// Get contract fee
	fee, err := s.contract.Fee(nil)
	if err != nil {
		return common.Hash{}, fmt.Errorf("failed to get contract fee: %w", err)
	}

	// Calculate total value to send (collateral amount + fee)
	totalValue := new(big.Int).Add(collateralAmount, fee)

	// Validate validator address
	valAddr, err := utils.ValidateAddress(validatorAddr)
	if err != nil {
		return common.Hash{}, fmt.Errorf("invalid validator address: %w", err)
	}

	// Show summary
	fmt.Println("===== Deposit Collateral Transaction =====")
	fmt.Printf("Validator Address          : %s\n", valAddr.Hex())
	fmt.Printf("Collateral Amount          : %s MITO\n", utils.FormatWeiToEther(collateralAmount))
	fmt.Printf("Fee                        : %s MITO\n", utils.FormatWeiToEther(fee))
	fmt.Printf("Total Value                : %s MITO\n", utils.FormatWeiToEther(totalValue))
	fmt.Println()
	fmt.Println("🚨 IMPORTANT: Only permitted collateral owners can deposit collateral for a validator.")
	fmt.Println("If your address is not a permitted collateral owner for this validator, the transaction will fail.")
	fmt.Println()

	// Execute the transaction
	opts, err := s.builder.CreateTransactOpts(totalValue)
	if err != nil {
		return common.Hash{}, fmt.Errorf("failed to create transaction options: %w", err)
	}

	tx, err := s.contract.DepositCollateral(opts, valAddr)
	if err != nil {
		return common.Hash{}, fmt.Errorf("failed to deposit collateral: %w", err)
	}

	return tx.Hash(), nil
}

// WithdrawCollateral withdraws collateral from a validator
func (s *CollateralService) WithdrawCollateral(validatorAddr, amount, receiver string) (common.Hash, error) {
	// Parse collateral amount as decimal MITO and convert to wei
	collateralAmount, err := utils.ParseValueAsWei(amount)
	if err != nil {
		return common.Hash{}, fmt.Errorf("failed to parse amount: %w", err)
	}

	// Ensure collateral amount is greater than 0
	if collateralAmount.Cmp(big.NewInt(0)) <= 0 {
		return common.Hash{}, fmt.Errorf("collateral amount must be greater than 0")
	}

	// Get contract fee
	fee, err := s.contract.Fee(nil)
	if err != nil {
		return common.Hash{}, fmt.Errorf("failed to get contract fee: %w", err)
	}

	// Validate addresses
	valAddr, err := utils.ValidateAddress(validatorAddr)
	if err != nil {
		return common.Hash{}, fmt.Errorf("invalid validator address: %w", err)
	}

	receiverAddr, err := utils.ValidateAddress(receiver)
	if err != nil {
		return common.Hash{}, fmt.Errorf("invalid receiver address: %w", err)
	}

	// Show summary
	fmt.Println("===== Withdraw Collateral Transaction =====")
	fmt.Printf("Validator Address          : %s\n", valAddr.Hex())
	fmt.Printf("Receiver Address           : %s\n", receiverAddr.Hex())
	fmt.Printf("Collateral Amount          : %s MITO\n", utils.FormatWeiToEther(collateralAmount))
	fmt.Printf("Fee                        : %s MITO\n", utils.FormatWeiToEther(fee))
	fmt.Println()
	fmt.Println("🚨 IMPORTANT: Only the collateral owner can withdraw collateral.")
	fmt.Println()

	// Execute the transaction (withdraw only sends fee, not collateral)
	opts, err := s.builder.CreateTransactOpts(fee)
	if err != nil {
		return common.Hash{}, fmt.Errorf("failed to create transaction options: %w", err)
	}

	tx, err := s.contract.WithdrawCollateral(opts, valAddr, receiverAddr, collateralAmount)
	if err != nil {
		return common.Hash{}, fmt.Errorf("failed to withdraw collateral: %w", err)
	}

	return tx.Hash(), nil
}

// SetPermittedCollateralOwner sets a permitted collateral owner for a validator
func (s *CollateralService) SetPermittedCollateralOwner(validatorAddr, collateralOwner string, isPermitted bool) (common.Hash, error) {
	// Validate addresses
	valAddr, err := utils.ValidateAddress(validatorAddr)
	if err != nil {
		return common.Hash{}, fmt.Errorf("invalid validator address: %w", err)
	}

	collateralOwnerAddr, err := utils.ValidateAddress(collateralOwner)
	if err != nil {
		return common.Hash{}, fmt.Errorf("invalid collateral owner address: %w", err)
	}

	// Show summary
	permissionText := "DENY"
	if isPermitted {
		permissionText = "PERMIT"
	}

	fmt.Println("===== Set Permitted Collateral Owner Transaction =====")
	fmt.Printf("Validator Address          : %s\n", valAddr.Hex())
	fmt.Printf("Collateral Owner           : %s\n", collateralOwnerAddr.Hex())
	fmt.Printf("Permission                 : %s\n", permissionText)
	fmt.Println()

	// Execute the transaction (no value needed for permission update)
	opts, err := s.builder.CreateTransactOpts(nil)
	if err != nil {
		return common.Hash{}, fmt.Errorf("failed to create transaction options: %w", err)
	}

	tx, err := s.contract.SetPermittedCollateralOwner(opts, valAddr, collateralOwnerAddr, isPermitted)
	if err != nil {
		return common.Hash{}, fmt.Errorf("failed to set permitted collateral owner: %w", err)
	}

	return tx.Hash(), nil
}

// TransferCollateralOwnership transfers collateral ownership for a validator
func (s *CollateralService) TransferCollateralOwnership(validatorAddr, newOwner string) (common.Hash, error) {
	// Get contract fee
	fee, err := s.contract.Fee(nil)
	if err != nil {
		return common.Hash{}, fmt.Errorf("failed to get contract fee: %w", err)
	}

	// Validate addresses
	valAddr, err := utils.ValidateAddress(validatorAddr)
	if err != nil {
		return common.Hash{}, fmt.Errorf("invalid validator address: %w", err)
	}

	newOwnerAddr, err := utils.ValidateAddress(newOwner)
	if err != nil {
		return common.Hash{}, fmt.Errorf("invalid new owner address: %w", err)
	}

	// Show summary
	fmt.Println("===== Transfer Collateral Ownership Transaction =====")
	fmt.Printf("Validator Address          : %s\n", valAddr.Hex())
	fmt.Printf("New Owner                  : %s\n", newOwnerAddr.Hex())
	fmt.Printf("Fee                        : %s MITO\n", utils.FormatWeiToEther(fee))
	fmt.Println()

	// Execute the transaction
	opts, err := s.builder.CreateTransactOpts(fee)
	if err != nil {
		return common.Hash{}, fmt.Errorf("failed to create transaction options: %w", err)
	}

	tx, err := s.contract.TransferCollateralOwnership(opts, valAddr, newOwnerAddr)
	if err != nil {
		return common.Hash{}, fmt.Errorf("failed to transfer collateral ownership: %w", err)
	}

	return tx.Hash(), nil
}
