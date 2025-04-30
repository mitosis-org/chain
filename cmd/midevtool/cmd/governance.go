package cmd

import (
	"encoding/json"
	"fmt"
	"log"
	"math/big"
	"os"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"github.com/mitosis-org/chain/bindings"
	"github.com/mitosis-org/chain/cmd/midevtool/utils"
	"github.com/spf13/cobra"
)

// NewGovernanceCmd creates a new governance command
func NewGovernanceCmd() *cobra.Command {
	governanceCmd := &cobra.Command{
		Use:   "governance",
		Short: "Interact with Consensus Governance contract",
		Long:  `Execute transactions on the Consensus Governance Entrypoint contract.`,
	}

	// Add subcommands
	governanceCmd.AddCommand(newGovernanceExecuteCmd())

	return governanceCmd
}

// newGovernanceExecuteCmd creates a new execute subcommand
func newGovernanceExecuteCmd() *cobra.Command {
	var (
		rpcURL                    string
		privateKey                string
		msgFile                   string
		msgString                 string
		govEntrypointContractAddr string
	)

	executeCmd := &cobra.Command{
		Use:   "execute",
		Short: "Execute messages through the Governance Entrypoint",
		Long:  `Submit Cosmos SDK messages for execution through the EVM Governance Entrypoint.`,
		Run: func(cmd *cobra.Command, args []string) {
			// Check required flags
			if privateKey == "" {
				log.Fatal("Error: private key is required")
			}

			if msgFile == "" && msgString == "" {
				log.Fatal("Error: either --msg-file or --msg must be provided")
			}

			if govEntrypointContractAddr == "" {
				log.Fatal("Error: ConsensusGovernanceEntrypoint contract address is required")
			}

			// Parse messages
			messages, err := parseMessages(msgFile, msgString)
			if err != nil {
				log.Fatalf("Error parsing messages: %v", err)
			}
			log.Printf("messages: %v", messages)

			// Setup client and contract
			privKey := utils.GetPrivateKey(privateKey)
			client, err := utils.GetEthClient(rpcURL)
			if err != nil {
				log.Fatalf("Error connecting to Ethereum client: %v", err)
			}

			// Get contract instance
			contract, err := bindings.NewConsensusGovernanceEntrypoint(
				common.HexToAddress(govEntrypointContractAddr),
				client,
			)
			if err != nil {
				log.Fatalf("Error initializing contract: %v", err)
			}

			// Create transaction options with zero value
			opts := utils.CreateTransactOpts(client, privKey, big.NewInt(0))

			// Execute the transaction
			tx, err := contract.Execute(opts, messages)
			if err != nil {
				log.Fatalf("Error executing transaction: %v", err)
			}

			log.Printf("Transaction sent: %s", tx.Hash().Hex())
		},
	}

	// Command-specific flags
	executeCmd.Flags().StringVar(&rpcURL, "rpc-url", "http://localhost:8545", "Ethereum RPC URL")
	executeCmd.Flags().StringVar(&privateKey, "private-key", "", "Private key for signing transactions")
	executeCmd.Flags().StringVar(&msgFile, "msg-file", "", "A file containing JSON messages")
	executeCmd.Flags().StringVar(&msgString, "msg", "", "JSON message string")
	executeCmd.Flags().StringVar(&govEntrypointContractAddr, "entrypoint", "", "ConsensusGovernanceEntrypoint contract address")

	return executeCmd
}

// parseMessages parses messages from a file or string, validating that it's a JSON array.
func parseMessages(filePath, msgString string) ([]string, error) {
	var content string

	if filePath != "" {
		// Read from file
		data, err := os.ReadFile(filePath)
		if err != nil {
			return nil, err
		}
		content = string(data)
	} else if msgString != "" {
		content = msgString
	} else {
		return nil, fmt.Errorf("no message provided")
	}

	// Validate that the content is a JSON array
	content = strings.TrimSpace(content)

	// Check if it starts with "[" and ends with "]"
	if !strings.HasPrefix(content, "[") || !strings.HasSuffix(content, "]") {
		return nil, fmt.Errorf("message must be a JSON array (starting with '[' and ending with ']')")
	}

	// Validate that it's a proper JSON array
	var jsonArray []json.RawMessage
	if err := json.Unmarshal([]byte(content), &jsonArray); err != nil {
		return nil, fmt.Errorf("invalid JSON array: %w", err)
	}

	// Convert each JSON object to string and return
	messages := make([]string, len(jsonArray))
	for i, msgRaw := range jsonArray {
		messages[i] = string(msgRaw)
	}

	return messages, nil
}
