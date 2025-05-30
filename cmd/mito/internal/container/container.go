package container

import (
	"github.com/mitosis-org/chain/cmd/mito/internal/client"
	"github.com/mitosis-org/chain/cmd/mito/internal/config"
	"github.com/mitosis-org/chain/cmd/mito/internal/tx"
)

// Container holds all dependencies for transaction operations
type Container struct {
	Config            *config.ResolvedConfig
	EthClient         *client.EthereumClient
	Contract          *client.ValidatorManagerContract
	TxBuilder         *tx.Builder
	TxSender          *tx.Sender
	ValidatorService  *tx.ValidatorService
	CollateralService *tx.CollateralService
}

// NewContainer creates a new dependency container
func NewContainer(resolvedConfig *config.ResolvedConfig) (*Container, error) {
	// Create Ethereum client
	ethClient, err := client.NewEthereumClient(resolvedConfig.RpcURL)
	if err != nil && resolvedConfig.RpcURL != "" {
		return nil, err
	}

	// Create contract instance
	contract, err := client.NewValidatorManagerContract(resolvedConfig.ValidatorManagerContractAddr, ethClient)
	if err != nil && resolvedConfig.ValidatorManagerContractAddr != "" {
		return nil, err
	}

	if contract != nil && resolvedConfig.ContractFee == "" {
		contractFee, err := contract.Fee(nil)
		if err != nil {
			return nil, err
		}
		resolvedConfig.ContractFee = contractFee.String()
	}

	// Create transaction components
	txBuilder := tx.NewBuilder(resolvedConfig, ethClient)
	txSender := tx.NewSender(ethClient)

	// Create services
	validatorService := tx.NewValidatorService(resolvedConfig, txBuilder)
	collateralService := tx.NewCollateralService(resolvedConfig, txBuilder)

	return &Container{
		Config:            resolvedConfig,
		EthClient:         ethClient,
		Contract:          contract,
		TxBuilder:         txBuilder,
		TxSender:          txSender,
		ValidatorService:  validatorService,
		CollateralService: collateralService,
	}, nil
}

// Close closes the Ethereum client connection
func (c *Container) Close() {
	if c.EthClient != nil {
		c.EthClient.Close()
	}
}
