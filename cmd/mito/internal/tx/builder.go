package tx

import (
	"context"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/mitosis-org/chain/cmd/mito/internal/client"
	"github.com/mitosis-org/chain/cmd/mito/internal/config"
	"github.com/mitosis-org/chain/cmd/mito/internal/units"
)

// Builder handles transaction building
type Builder struct {
	config    *config.ResolvedConfig
	ethClient *client.EthereumClient
	signer    *Signer
}

// NewBuilder creates a new transaction builder
func NewBuilder(config *config.ResolvedConfig, ethClient *client.EthereumClient) *Builder {
	return &Builder{
		config:    config,
		ethClient: ethClient,
		signer:    NewSigner(config),
	}
}

// GetSignerAddress returns the address of the transaction signer
func (b *Builder) GetSignerAddress() (common.Address, error) {
	return b.signer.GetSignerAddress()
}

// SignTransaction signs a transaction
func (b *Builder) SignTransaction(tx *types.Transaction) (*types.Transaction, error) {
	// Get chain ID
	chainID, err := b.GetChainID()
	if err != nil {
		return nil, fmt.Errorf("failed to get chain ID: %w", err)
	}

	return b.signer.SignTransaction(tx, chainID.Int64())
}

// CreateTransactionFromData creates a types.Transaction from TransactionData
func (b *Builder) CreateTransactionFromData(txData *TransactionData) (*types.Transaction, error) {
	return b.CreateTransactionFromDataWithOptions(txData, false)
}

// CreateTransactionFromDataWithOptions creates a types.Transaction from TransactionData with unsigned option
func (b *Builder) CreateTransactionFromDataWithOptions(txData *TransactionData, unsigned bool) (*types.Transaction, error) {
	// Determine nonce - use specified nonce or get from client (if possible)
	var nonce uint64
	if b.config.Nonce != "" {
		nonceInt, ok := new(big.Int).SetString(b.config.Nonce, 10)
		if ok {
			nonce = nonceInt.Uint64()
		}

		if !ok {
			return nil, fmt.Errorf("invalid nonce: %s", b.config.Nonce)
		}
	} else {
		// Try to get signer address for nonce
		signerAddr, err := b.GetSignerAddress()
		if err != nil {
			return nil, fmt.Errorf("failed to get signer address: %w", err)
		}

		// Check if network is available for nonce lookup
		if b.ethClient == nil {
			return nil, fmt.Errorf("RPC connection required to get nonce automatically. Please provide --nonce manually or set RPC URL with --rpc-url or 'mito config set-rpc'")
		}

		nonce, err = b.ethClient.PendingNonceAt(context.Background(), signerAddr)
		if err != nil {
			return nil, fmt.Errorf("failed to get nonce from network: %w", err)
		}
	}

	// Determine gas limit
	gasLimit := txData.GasLimit
	if gasLimit == 0 && b.config.GasLimit > 0 {
		gasLimit = b.config.GasLimit
	}
	if gasLimit == 0 {
		if b.ethClient != nil && !unsigned {
			// Only estimate gas for signed transactions since we need a valid from address
			fromAddr, err := b.GetSignerAddress()
			if err != nil {
				// Fallback to default if we can't get signer address
				gasLimit = 200000
			} else {
				msg := ethereum.CallMsg{
					From:  fromAddr,
					To:    &txData.To,
					Value: txData.Value,
					Data:  txData.Data,
				}

				estimatedGas, err := b.ethClient.EstimateGas(context.Background(), msg)
				if err != nil {
					return nil, fmt.Errorf("gas estimation failed: %w", err)
				} else {
					// Add 20% buffer to estimated gas
					gasLimit = estimatedGas + (estimatedGas / 5)
				}
			}
		} else {
			gasLimit = 200000 // Conservative default when no RPC available or unsigned transaction
		}
	}

	// Determine gas price - use specified gas price or get from client
	var gasPrice *big.Int
	if b.config.GasPrice != "" {
		// Parse gas price using new units module (defaults to gwei)
		var err error
		gasPrice, err = units.ParseGasPriceInput(b.config.GasPrice)
		if err != nil {
			return nil, fmt.Errorf("failed to parse gas price: %w", err)
		}
	}
	if gasPrice == nil {
		// Check if network is available for gas price lookup
		if b.ethClient == nil {
			return nil, fmt.Errorf("RPC connection required to get gas price automatically. Please provide --gas-price manually or set RPC URL with --rpc-url or 'mito config set-rpc'")
		}

		var err error
		gasPrice, err = b.ethClient.SuggestGasPrice(context.Background())
		if err != nil {
			return nil, fmt.Errorf("failed to get gas price from network: %w", err)
		}
	}

	// Create the transaction
	tx := types.NewTransaction(
		nonce,
		txData.To,
		txData.Value,
		gasLimit,
		gasPrice,
		txData.Data,
	)

	var err error
	if !unsigned {
		tx, err = b.SignTransaction(tx)
		if err != nil {
			return nil, fmt.Errorf("failed to sign transaction: %w", err)
		}
	}

	return tx, nil
}

func (b *Builder) GetChainID() (*big.Int, error) {
	// Get chain ID if not specified in config
	var chainID *big.Int
	if b.config.ChainID != "" {
		chainID, _ = new(big.Int).SetString(b.config.ChainID, 10)
	}
	if chainID == nil {
		// Check if network is available for chain ID lookup
		if b.ethClient == nil {
			return nil, fmt.Errorf("RPC connection required to get chain ID automatically. Please provide --chain-id manually or set RPC URL with --rpc-url or 'mito config set-rpc'")
		}

		var err error
		chainID, err = b.ethClient.ChainID(context.Background())
		if err != nil {
			return nil, fmt.Errorf("failed to get chain ID from network: %w", err)
		}
	}

	return chainID, nil
}
