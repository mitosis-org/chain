package cmd

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"math/big"
	"os"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	ethcrypto "github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/mitosis-org/chain/bindings"
	"github.com/mitosis-org/chain/cmd/mivalidator/utils"
	"github.com/spf13/cobra"
)

// Common flags used across multiple commands
var (
	rpcURL                       string
	privateKey                   string
	validatorManagerContractAddr string
	yes                          bool
	nonce                        uint64
	nonceSpecified               bool

	// Shared client and contract instances
	client   *ethclient.Client
	contract *bindings.IValidatorManager
)

// mustMarkFlagRequired marks a flag as required and panics if it fails
func mustMarkFlagRequired(cmd *cobra.Command, flag string) {
	if err := cmd.MarkFlagRequired(flag); err != nil {
		log.Fatalf("Failed to mark flag '%s' as required: %v", flag, err)
	}
}

// AddCommonFlags adds common flags to a command
func AddCommonFlags(cmd *cobra.Command, readonly bool) {
	cmd.Flags().StringVar(&rpcURL, "rpc-url", "http://localhost:8545", "Ethereum RPC URL")
	cmd.Flags().StringVar(&privateKey, "private-key", "", "Private key for signing transactions")
	cmd.Flags().StringVar(&validatorManagerContractAddr, "contract", "", "ValidatorManager contract address")
	cmd.Flags().BoolVar(&yes, "yes", false, "Skip confirmation prompt")
	cmd.Flags().Uint64Var(&nonce, "nonce", 0, "Manually specify nonce for transaction (optional)")

	// Mark required flags
	if !readonly {
		mustMarkFlagRequired(cmd, "private-key")
	}
	mustMarkFlagRequired(cmd, "contract")

	// Preserve any existing PreRun function
	existingPreRun := cmd.PreRun
	cmd.PreRun = func(cmd *cobra.Command, args []string) {
		var err error

		// Set the nonceSpecified flag
		nonceSpecified = cmd.Flags().Changed("nonce")

		// Setup client
		client, err = utils.GetEthClient(rpcURL)
		if err != nil {
			log.Fatalf("%v", err)
		}

		// Get contract instance
		contract, err = GetValidatorManagerContract(client)
		if err != nil {
			log.Fatalf("%v", err)
		}

		// Call the existing PreRun if it exists
		if existingPreRun != nil {
			existingPreRun(cmd, args)
		}
	}
}

// GetValidatorManagerContract initializes and returns the ValidatorManager contract
func GetValidatorManagerContract(client *ethclient.Client) (*bindings.IValidatorManager, error) {
	if validatorManagerContractAddr == "" {
		return nil, fmt.Errorf("ValidatorManager contract address is required")
	}

	validatorManagerAddr := common.HexToAddress(validatorManagerContractAddr)
	contract, err := bindings.NewIValidatorManager(validatorManagerAddr, client)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize ValidatorManager contract: %w", err)
	}

	return contract, nil
}

// ConfirmAction prompts the user to confirm an action
// If the yes flag is true, returns true without prompting
func ConfirmAction(message string) bool {
	if yes {
		return true
	}

	fmt.Printf("%s\nType 'yes' to continue: ", message)
	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)
	return strings.ToLower(input) == "yes"
}

// TransactOpts creates transaction options for a contract call
func TransactOpts(value *big.Int) *bind.TransactOpts {
	// Get private key from the flag
	privKey := utils.GetPrivateKey(privateKey)

	// Get address from private key
	addr := common.BytesToAddress(ethcrypto.PubkeyToAddress(privKey.PublicKey).Bytes())

	// Determine nonce - use specified nonce or get from client
	var nVal uint64
	var err error
	if nonceSpecified {
		nVal = nonce
	} else {
		nVal, err = client.PendingNonceAt(context.Background(), addr)
		if err != nil {
			panic(fmt.Errorf("failed to get nonce: %w", err))
		}
	}

	// Get chain ID
	chainID, err := client.ChainID(context.Background())
	if err != nil {
		panic(fmt.Errorf("failed to get chain ID: %w", err))
	}

	// Create transaction options
	opts, err := bind.NewKeyedTransactorWithChainID(privKey, chainID)
	if err != nil {
		panic(fmt.Errorf("failed to create transaction options: %w", err))
	}

	// Set nonce and value
	opts.Nonce = new(big.Int).SetUint64(nVal)
	opts.Value = value

	return opts
}

// WaitForTxConfirmation waits for a transaction to be mined and confirmed
func WaitForTxConfirmation(client *ethclient.Client, txHash common.Hash) error {
	fmt.Printf("Waiting for transaction %s to be confirmed...\n", txHash.Hex())

	ctx := context.Background()

	// Set a timeout for 1 minutes
	timeoutCtx, cancel := context.WithTimeout(ctx, 1*time.Minute)
	defer cancel()

	// Poll for transaction receipt with a 2-second interval
	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-timeoutCtx.Done():
			return fmt.Errorf("timeout waiting for transaction confirmation")
		case <-ticker.C:
			receipt, err := client.TransactionReceipt(ctx, txHash)
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
