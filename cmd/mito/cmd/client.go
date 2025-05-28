package cmd

import (
	"context"
	"fmt"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/mitosis-org/chain/bindings"
)

// GetEthClient creates and returns an Ethereum client
func GetEthClient(rpcURL string) (*ethclient.Client, error) {
	return ethclient.Dial(rpcURL)
}

// ConnectToEthereum creates and returns an Ethereum client (alias for GetEthClient)
func ConnectToEthereum(rpcURL string) (*ethclient.Client, error) {
	return GetEthClient(rpcURL)
}

// GetValidatorManagerContract initializes and returns the ValidatorManager contract
func GetValidatorManagerContract(ethClient *ethclient.Client) (*bindings.IValidatorManager, error) {
	if validatorManagerContractAddr == "" {
		return nil, fmt.Errorf("ValidatorManager contract address is required")
	}

	validatorManagerAddr := common.HexToAddress(validatorManagerContractAddr)
	contract, err := bindings.NewIValidatorManager(validatorManagerAddr, ethClient)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize ValidatorManager contract: %w", err)
	}

	return contract, nil
}

// WaitForTxConfirmation waits for a transaction to be mined and confirmed
func WaitForTxConfirmation(ethClient *ethclient.Client, txHash common.Hash) error {
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
			receipt, err := ethClient.TransactionReceipt(ctx, txHash)
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
