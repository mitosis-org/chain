package cmd

import (
	"bufio"
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"math/big"
	"os"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
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
	generateUnsigned             bool
	fromAddress                  string

	// Shared client and contract instances
	client   *ethclient.Client
	contract *bindings.IValidatorManager
)

// AddCommonFlags adds common flags to a command
func AddCommonFlags(cmd *cobra.Command, readonly bool) {
	cmd.Flags().StringVar(&rpcURL, "rpc-url", "http://localhost:8545", "Ethereum RPC URL")
	cmd.Flags().StringVar(&privateKey, "private-key", "", "Private key for signing transactions")
	cmd.Flags().StringVar(&validatorManagerContractAddr, "contract", "", "ValidatorManager contract address")
	cmd.Flags().BoolVar(&yes, "yes", false, "Skip confirmation prompt")
	cmd.Flags().Uint64Var(&nonce, "nonce", 0, "Manually specify nonce for transaction (optional)")
	cmd.Flags().BoolVar(&generateUnsigned, "unsigned", false, "Generate unsigned transaction data instead of sending transaction")
	cmd.Flags().StringVar(&fromAddress, "from", "", "From address for unsigned transaction (required when --unsigned is used)")

	// Mark required flags
	if !readonly {
		if cmd.Annotations == nil {
			cmd.Annotations = make(map[string]string)
		}
		cmd.Annotations["private-key-required"] = "true"
	}
	cmd.MarkFlagRequired("contract")

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

		// Check private key or unsigned mode
		if !readonly && !generateUnsigned && privateKey == "" {
			log.Fatalf("Either --private-key or --unsigned with --from must be provided")
		}

		// Check from address is provided when unsigned is used
		if generateUnsigned && fromAddress == "" {
			log.Fatalf("--from address is required when using --unsigned")
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
	// If generating unsigned transaction
	if generateUnsigned {
		// Parse from address
		fromAddr := common.HexToAddress(fromAddress)

		// Determine nonce - use specified nonce or get from client
		var nVal uint64
		var err error
		if nonceSpecified {
			nVal = nonce
		} else {
			nVal, err = client.PendingNonceAt(context.Background(), fromAddr)
			if err != nil {
				panic(fmt.Errorf("failed to get nonce: %w", err))
			}
		}

		// Create dummy signer that doesn't actually sign
		opts := &bind.TransactOpts{
			From:   fromAddr,
			Nonce:  big.NewInt(int64(nVal)),
			Value:  value,
			NoSend: true,
			Signer: func(address common.Address, tx *types.Transaction) (*types.Transaction, error) {
				return tx, nil
			},
		}

		return opts
	} else {
		// Regular transaction with private key
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
}

// PrintUnsignedTransaction formats and prints the unsigned transaction data in a JSON format
func PrintUnsignedTransaction(tx *types.Transaction) {
	// Extract chain ID for the transaction
	chainID, err := client.ChainID(context.Background())
	if err != nil {
		log.Fatalf("Failed to get chain ID: %v", err)
	}

	// Convert binary data to hex string, keep type in hex format
	typeHex := fmt.Sprintf("0x%x", tx.Type())
	dataHex := "0x" + hex.EncodeToString(tx.Data())

	var txJSON []byte

	// Handle EIP-1559 transaction (type 2)
	if tx.Type() == 2 {
		// Define struct for type 2 transaction with specific field order
		txData := struct {
			Type                 string `json:"type"`
			ChainID              string `json:"chainId"`
			From                 string `json:"from"`
			To                   string `json:"to"`
			Value                string `json:"value"`
			Data                 string `json:"data"`
			Nonce                string `json:"nonce"`
			Gas                  string `json:"gas"`
			MaxFeePerGas         string `json:"maxFeePerGas"`
			MaxPriorityFeePerGas string `json:"maxPriorityFeePerGas"`
		}{
			Type:                 typeHex,
			ChainID:              chainID.String(),
			From:                 fromAddress,
			To:                   tx.To().Hex(),
			Value:                tx.Value().String(),
			Data:                 dataHex,
			Nonce:                fmt.Sprintf("%d", tx.Nonce()),
			Gas:                  fmt.Sprintf("%d", tx.Gas()),
			MaxFeePerGas:         tx.GasFeeCap().String(),
			MaxPriorityFeePerGas: tx.GasTipCap().String(),
		}

		// Marshal to JSON
		txJSON, err = json.MarshalIndent(txData, "", "  ")
		if err != nil {
			log.Fatalf("Failed to marshal transaction to JSON: %v", err)
		}
	} else if tx.Type() == 0 { // Legacy transaction
		// Define struct for legacy transaction with specific field order
		txData := struct {
			ChainID  string `json:"chainId"`
			From     string `json:"from"`
			To       string `json:"to"`
			Value    string `json:"value"`
			Data     string `json:"data"`
			Nonce    string `json:"nonce"`
			Gas      string `json:"gas"`
			GasPrice string `json:"gasPrice"`
		}{
			ChainID:  chainID.String(),
			From:     fromAddress,
			To:       tx.To().Hex(),
			Value:    tx.Value().String(),
			Data:     dataHex,
			Nonce:    fmt.Sprintf("%d", tx.Nonce()),
			Gas:      fmt.Sprintf("%d", tx.Gas()),
			GasPrice: tx.GasPrice().String(),
		}

		// Marshal to JSON
		txJSON, err = json.MarshalIndent(txData, "", "  ")
		if err != nil {
			log.Fatalf("Failed to marshal transaction to JSON: %v", err)
		}
	} else {
		log.Fatalf("Unsupported transaction type: %d. Only legacy (type 0) and EIP-1559 (type 2) transactions are supported.", tx.Type())
	}

	fmt.Println("\n===== Unsigned Transaction Data =====")
	fmt.Println(string(txJSON))
	fmt.Println("\nThis transaction data can be signed offline with hardware wallets or other signing tools.")
}

// HandleTransaction processes a transaction - either printing unsigned data or sending and waiting for confirmation
func HandleTransaction(tx *types.Transaction) error {
	if generateUnsigned {
		PrintUnsignedTransaction(tx)
		return nil
	} else {
		return WaitForTxConfirmation(client, tx.Hash())
	}
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
