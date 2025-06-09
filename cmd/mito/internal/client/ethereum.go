package client

import (
	"context"
	"fmt"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

// EthereumClient wraps the go-ethereum client with additional functionality
type EthereumClient struct {
	*ethclient.Client
	rpcURL string
}

// NewEthereumClient creates and returns a new Ethereum client
func NewEthereumClient(rpcURL string) (*EthereumClient, error) {
	if rpcURL == "" {
		return nil, fmt.Errorf("RPC URL is required")
	}

	client, err := ethclient.Dial(rpcURL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to Ethereum node at %s: %w", rpcURL, err)
	}

	return &EthereumClient{
		Client: client,
		rpcURL: rpcURL,
	}, nil
}

// GetRPCURL returns the RPC URL used by this client
func (c *EthereumClient) GetRPCURL() string {
	return c.rpcURL
}

// WaitForTxConfirmation waits for a transaction to be mined and confirmed
func (c *EthereumClient) WaitForTxConfirmation(txHash common.Hash) error {
	fmt.Printf("Waiting for transaction %s to be confirmed...\n", txHash.Hex())

	ctx := context.Background()

	// Set a timeout for 2 minutes
	timeoutCtx, cancel := context.WithTimeout(ctx, 2*time.Minute)
	defer cancel()

	// Poll for transaction receipt with a 2-second interval
	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-timeoutCtx.Done():
			return fmt.Errorf("timeout waiting for transaction confirmation")
		case <-ticker.C:
			receipt, err := c.TransactionReceipt(ctx, txHash)
			if err != nil {
				// If error, likely the tx is not yet mined
				fmt.Print(".")
				continue
			}

			// Once we have a receipt, check its status
			if receipt.Status == 1 {
				blockNumber := receipt.BlockNumber
				fmt.Printf("\nTransaction confirmed in block %d\n", blockNumber.Uint64())
				return nil
			} else {
				return fmt.Errorf("transaction failed with status: %d", receipt.Status)
			}
		}
	}
}
