package tx

import (
	"context"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/mitosis-org/chain/cmd/mito/internal/client"
	"github.com/mitosis-org/chain/cmd/mito/internal/config"
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
	} else if unsigned {
		// For unsigned transactions, use nonce 0 as default
		nonce = 0
	} else {
		// Try to get signer address for nonce
		signerAddr, err := b.GetSignerAddress()
		if err != nil {
			return nil, fmt.Errorf("failed to get signer address: %w", err)
		}
		nonce, err = b.ethClient.PendingNonceAt(context.Background(), signerAddr)
		if err != nil {
			return nil, fmt.Errorf("failed to get nonce: %w", err)
		}
	}

	// Determine gas limit
	gasLimit := txData.GasLimit
	if gasLimit == 0 && b.config.GasLimit > 0 {
		gasLimit = b.config.GasLimit
	}
	if gasLimit == 0 {
		gasLimit = 500000 // Default
	}

	// Determine gas price - use specified gas price or get from client
	var gasPrice *big.Int
	if b.config.GasPrice != "" {
		gasPrice, _ = new(big.Int).SetString(b.config.GasPrice, 10)
	}
	if gasPrice == nil {
		var err error
		gasPrice, err = b.ethClient.SuggestGasPrice(context.Background())
		if err != nil {
			return nil, fmt.Errorf("failed to get gas price from RPC: %w", err)
		}
	}

	var contractFee *big.Int
	if b.config.ContractFee != "" {
		contractFee, _ = new(big.Int).SetString(b.config.ContractFee, 10)
	}

	totalValue := new(big.Int).Add(txData.Value, contractFee)

	// Create the transaction
	tx := types.NewTransaction(
		nonce,
		txData.To,
		totalValue,
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
		var err error
		chainID, err = b.ethClient.ChainID(context.Background())
		if err != nil {
			return nil, fmt.Errorf("failed to get chain ID from RPC: %w", err)
		}
	}

	return chainID, nil
}
