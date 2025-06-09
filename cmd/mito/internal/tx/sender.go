package tx

import (
	"context"
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/mitosis-org/chain/cmd/mito/internal/client"
)

// Sender handles transaction sending and confirmation
type Sender struct {
	ethClient *client.EthereumClient
}

// NewSender creates a new transaction sender
func NewSender(ethClient *client.EthereumClient) *Sender {
	return &Sender{
		ethClient: ethClient,
	}
}

// SendTransaction sends a signed transaction to the network
func (s *Sender) SendTransaction(signedTx *types.Transaction) (common.Hash, error) {
	ctx := context.Background()

	// Send the transaction
	err := s.ethClient.SendTransaction(ctx, signedTx)
	if err != nil {
		return common.Hash{}, fmt.Errorf("failed to send transaction: %w", err)
	}

	txHash := signedTx.Hash()
	fmt.Printf("Transaction sent with hash: %s\n", txHash.Hex())

	return txHash, nil
}

// SendAndWait sends a transaction and waits for confirmation
func (s *Sender) SendAndWait(signedTx *types.Transaction) (common.Hash, error) {
	// Send the transaction
	txHash, err := s.SendTransaction(signedTx)
	if err != nil {
		return common.Hash{}, err
	}

	// Wait for confirmation
	err = s.ethClient.WaitForTxConfirmation(txHash)
	if err != nil {
		return txHash, fmt.Errorf("transaction sent but confirmation failed: %w", err)
	}

	return txHash, nil
}
