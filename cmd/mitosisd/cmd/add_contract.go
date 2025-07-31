package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/spf13/cobra"
)

// FoundryArtifact represents a Foundry compilation artifact
type FoundryArtifact struct {
	Bytecode struct {
		Object string `json:"object"`
	} `json:"bytecode"`
	DeployedBytecode struct {
		Object string `json:"object"`
	} `json:"deployedBytecode"`
}

// ContractData holds contract information for genesis
type ContractData struct {
	Address string            `json:"address"`
	Code    string            `json:"code"`
	Storage map[string]string `json:"storage,omitempty"`
	Balance string            `json:"balance,omitempty"`
}

func NewAddContractCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "add-contract [contract-address] [artifact-file]",
		Short: "Add a smart contract to the genesis file from Foundry artifact",
		Long: `Add a smart contract to the genesis file from a Foundry compilation artifact.

Usage:
   mitosisd genesis add-contract 0x123... out/MyContract.sol/MyContract.json

The command will add the contract to the alloc section of the genesis file with the deployed bytecode.`,
		Args: cobra.ExactArgs(2),
		RunE: runAddContract,
	}

	cmd.Flags().String("balance", "0", "Initial balance for the contract account")
	cmd.Flags().Bool("use-creation-code", false, "Use creation bytecode instead of deployed bytecode")

	return cmd
}

func runAddContract(cmd *cobra.Command, args []string) error {
	contractAddress := args[0]
	artifactFile := args[1]

	// Get client context to access home directory
	clientCtx, err := client.GetClientQueryContext(cmd)
	if err != nil {
		clientCtx = client.Context{}
	}

	// Get home directory
	homeDir, _ := cmd.Flags().GetString(flags.FlagHome)
	if homeDir == "" {
		homeDir = clientCtx.HomeDir
	}

	// Construct path to Ethereum genesis file
	genesisFile := filepath.Join(homeDir, "config", "eth_genesis.json")

	// Validate contract address format
	if !strings.HasPrefix(contractAddress, "0x") || len(contractAddress) != 42 {
		return fmt.Errorf("invalid contract address format: %s (must be 0x followed by 40 hex characters)", contractAddress)
	}

	// Validate artifact file exists
	if _, err := os.Stat(artifactFile); os.IsNotExist(err) {
		return fmt.Errorf("artifact file does not exist: %s", artifactFile)
	}

	// Get flags
	balance, _ := cmd.Flags().GetString("balance")
	useCreationCode, _ := cmd.Flags().GetBool("use-creation-code")

	// Read bytecode from Foundry artifact file
	bytecode, err := readBytecodeFromArtifact(artifactFile, useCreationCode)  
	if err != nil {
		return fmt.Errorf("failed to read bytecode from artifact: %w", err)
	}

	// Validate bytecode
	if bytecode == "" || bytecode == "0x" {
		return fmt.Errorf("empty bytecode provided")
	}

	// Ensure bytecode has 0x prefix
	if !strings.HasPrefix(bytecode, "0x") {
		bytecode = "0x" + bytecode
	}

	// Read and parse genesis file
	genesis, err := readGenesisFile(genesisFile)
	if err != nil {
		return fmt.Errorf("failed to read genesis file: %w", err)
	}

	// Add contract to alloc
	if genesis.Alloc == nil {
		genesis.Alloc = make(map[string]AllocatedAccount)
	}

	account := AllocatedAccount{
		Balance: balance,
		Code:    bytecode,
	}

	genesis.Alloc[contractAddress] = account

	// Write updated genesis file
	err = writeGenesisFile(genesisFile, genesis)
	if err != nil {
		return fmt.Errorf("failed to write genesis file: %w", err)
	}

	fmt.Printf("Successfully added contract %s to genesis file %s\n", contractAddress, genesisFile)
	fmt.Printf("Bytecode length: %d bytes\n", (len(bytecode)-2)/2) // -2 for 0x prefix, /2for hex encoding

	return nil
}

func readBytecodeFromArtifact(artifactFile string, useCreationCode bool) (string, error) {
	data, err := os.ReadFile(artifactFile)
	if err != nil {
		return "", fmt.Errorf("failed to read artifact file: %w", err)
	}

	var artifact FoundryArtifact
	if err := json.Unmarshal(data, &artifact); err != nil {
		return "", fmt.Errorf("failed to parse artifact JSON: %w", err)
	}

	if useCreationCode {
		return artifact.Bytecode.Object, nil
	}

	return artifact.DeployedBytecode.Object, nil
}


func readGenesisFile(genesisFile string) (*EthereumGenesisSpec, error) {
	data, err := os.ReadFile(genesisFile)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	var genesis EthereumGenesisSpec
	if err := json.Unmarshal(data, &genesis); err != nil {
		return nil, fmt.Errorf("failed to parse genesis JSON: %w", err)
	}

	return &genesis, nil
}

func writeGenesisFile(genesisFile string, genesis *EthereumGenesisSpec) error {
	// Create backup
	backupFile := genesisFile + ".backup"
	if err := copyFile(genesisFile, backupFile); err != nil {
		fmt.Printf("Warning: failed to create backup file: %v\n", err)
	}

	// Marshal to JSON with proper formatting
	data, err := json.MarshalIndent(genesis, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal genesis JSON: %w", err)
	}

	// Write to file
	if err := os.WriteFile(genesisFile, data, 0o644); err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	return nil
}

func copyFile(src, dst string) error {
	sourceFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	destFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destFile.Close()

	_, err = io.Copy(destFile, sourceFile)
	return err
}