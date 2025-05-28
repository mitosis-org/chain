package tx

import (
	"context"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	ethcrypto "github.com/ethereum/go-ethereum/crypto"
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

// CreateTransactOpts creates transaction options for contract calls
func (b *Builder) CreateTransactOpts(value *big.Int) (*bind.TransactOpts, error) {
	// Get signing configuration
	signingConfig, err := b.signer.GetSigningConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to get signing configuration: %w", err)
	}

	// Get address from private key
	addr := ethcrypto.PubkeyToAddress(signingConfig.PrivateKey.PublicKey)

	// Determine nonce - use specified nonce or get from client
	var nVal uint64
	if b.config.NonceSet {
		nVal = b.config.Nonce
	} else {
		nVal, err = b.ethClient.PendingNonceAt(context.Background(), addr)
		if err != nil {
			return nil, fmt.Errorf("failed to get nonce: %w", err)
		}
	}

	// Get chain ID
	chainID, err := b.ethClient.ChainID(context.Background())
	if err != nil {
		return nil, fmt.Errorf("failed to get chain ID: %w", err)
	}

	// Create transaction options
	opts, err := bind.NewKeyedTransactorWithChainID(signingConfig.PrivateKey, chainID)
	if err != nil {
		return nil, fmt.Errorf("failed to create transaction options: %w", err)
	}

	// Set nonce and value
	opts.Nonce = new(big.Int).SetUint64(nVal)
	opts.Value = value

	// Set gas limit if specified
	if b.config.GasLimit > 0 {
		opts.GasLimit = b.config.GasLimit
	}

	// Set gas price if specified
	if b.config.GasPrice != "" {
		gasPrice, ok := new(big.Int).SetString(b.config.GasPrice, 10)
		if !ok {
			return nil, fmt.Errorf("invalid gas price format: %s", b.config.GasPrice)
		}
		opts.GasPrice = gasPrice
	}

	return opts, nil
}

// GetSignerAddress returns the address of the transaction signer
func (b *Builder) GetSignerAddress() (common.Address, error) {
	return b.signer.GetSignerAddress()
}

// SignTransaction signs a transaction
func (b *Builder) SignTransaction(tx *types.Transaction) (*types.Transaction, error) {
	// Get chain ID
	chainID, err := b.ethClient.ChainID(context.Background())
	if err != nil {
		return nil, fmt.Errorf("failed to get chain ID: %w", err)
	}

	return b.signer.SignTransaction(tx, chainID.Int64())
}
