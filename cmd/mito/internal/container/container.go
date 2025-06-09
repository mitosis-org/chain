package container

import (
	"fmt"
	"math/big"

	"github.com/mitosis-org/chain/cmd/mito/internal/client"
	"github.com/mitosis-org/chain/cmd/mito/internal/config"
	"github.com/mitosis-org/chain/cmd/mito/internal/tx"
)

// Container holds all dependencies for transaction operations
type Container struct {
	Config                   *config.ResolvedConfig
	EthClient                *client.EthereumClient
	ValidatorManagerContract *client.ValidatorManagerContract
	TxBuilder                *tx.Builder
	TxSender                 *tx.Sender
	ValidatorService         *tx.ValidatorService
	CollateralService        *tx.CollateralService
}

// NewContainer creates a new dependency container
func NewContainer(resolvedConfig *config.ResolvedConfig) (*Container, error) {
	// Create Ethereum client - allow nil for offline mode
	var ethClient *client.EthereumClient
	var err error
	if resolvedConfig.RpcURL != "" {
		ethClient, err = client.NewEthereumClient(resolvedConfig.RpcURL)
		if err != nil {
			return nil, fmt.Errorf("failed to create ethereum client: %w", err)
		}
	}

	// Create contract instance - allow nil for offline mode
	var contract *client.ValidatorManagerContract
	if resolvedConfig.ValidatorManagerContractAddr != "" && ethClient != nil {
		contract, err = client.NewValidatorManagerContract(resolvedConfig.ValidatorManagerContractAddr, ethClient)
		if err != nil {
			return nil, fmt.Errorf("failed to create contract instance: %w", err)
		}

		// Try to get contract fee if not specified
		if resolvedConfig.ContractFee.Cmp(big.NewInt(0)) == 0 {
			contractFee, err := contract.Fee(nil)
			if err != nil {
				// Don't fail if we can't get contract fee - use default
				resolvedConfig.ContractFee = big.NewInt(0)
			} else {
				resolvedConfig.ContractFee = contractFee
			}
		}
	}

	// Create transaction components - allow nil ethClient for offline mode
	txBuilder := tx.NewBuilder(resolvedConfig, ethClient)
	txSender := tx.NewSender(ethClient)

	// Create services
	validatorService := tx.NewValidatorService(resolvedConfig, txBuilder)
	collateralService := tx.NewCollateralService(resolvedConfig, txBuilder)

	return &Container{
		Config:                   resolvedConfig,
		EthClient:                ethClient,
		ValidatorManagerContract: contract,
		TxBuilder:                txBuilder,
		TxSender:                 txSender,
		ValidatorService:         validatorService,
		CollateralService:        collateralService,
	}, nil
}

// Close closes the Ethereum client connection
func (c *Container) Close() {
	if c.EthClient != nil {
		c.EthClient.Close()
	}
}
