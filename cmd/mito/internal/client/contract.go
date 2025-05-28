package client

import (
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/mitosis-org/chain/bindings"
)

// ValidatorManagerContract wraps the ValidatorManager contract with additional functionality
type ValidatorManagerContract struct {
	*bindings.IValidatorManager
	address common.Address
}

// NewValidatorManagerContract creates and returns a new ValidatorManager contract instance
func NewValidatorManagerContract(contractAddr string, ethClient *EthereumClient) (*ValidatorManagerContract, error) {
	if contractAddr == "" {
		return nil, fmt.Errorf("ValidatorManager contract address is required")
	}

	address := common.HexToAddress(contractAddr)
	contract, err := bindings.NewIValidatorManager(address, ethClient.Client)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize ValidatorManager contract: %w", err)
	}

	return &ValidatorManagerContract{
		IValidatorManager: contract,
		address:           address,
	}, nil
}

// GetAddress returns the contract address
func (c *ValidatorManagerContract) GetAddress() common.Address {
	return c.address
}
