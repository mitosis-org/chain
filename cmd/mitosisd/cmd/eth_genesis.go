package cmd

import (
	"encoding/json"
	"fmt"
	"math/big"
	"os"
	"path/filepath"
)

// EthereumChainConfig represents the chain configuration for Ethereum genesis
type EthereumChainConfig struct {
	ChainID                 *big.Int `json:"chainId"`
	HomesteadBlock          *big.Int `json:"homesteadBlock"`
	DAOForkBlock            *big.Int `json:"daoForkBlock,omitempty"`
	DAOForkSupport          bool     `json:"daoForkSupport,omitempty"`
	EIP150Block             *big.Int `json:"eip150Block"`
	EIP155Block             *big.Int `json:"eip155Block"`
	EIP158Block             *big.Int `json:"eip158Block"`
	ByzantiumBlock          *big.Int `json:"byzantiumBlock"`
	ConstantinopleBlock     *big.Int `json:"constantinopleBlock"`
	PetersburgBlock         *big.Int `json:"petersburgBlock"`
	IstanbulBlock           *big.Int `json:"istanbulBlock"`
	BerlinBlock             *big.Int `json:"berlinBlock"`
	LondonBlock             *big.Int `json:"londonBlock"`
	ArrowGlacierBlock       *big.Int `json:"arrowGlacierBlock,omitempty"`
	GrayGlacierBlock        *big.Int `json:"grayGlacierBlock,omitempty"`
	MergeNetsplitBlock      *big.Int `json:"mergeNetsplitBlock,omitempty"`
	TerminalTotalDifficulty *big.Int `json:"terminalTotalDifficulty"`
	ShanghaiTime            *uint64  `json:"shanghaiTime"`
	CancunTime              *uint64  `json:"cancunTime"`
	PragueTime              *uint64  `json:"pragueTime,omitempty"`
}

// EthereumGenesisSpec represents the Ethereum genesis specification
type EthereumGenesisSpec struct {
	Config     *EthereumChainConfig        `json:"config"`
	Nonce      string                      `json:"nonce"`
	Timestamp  string                      `json:"timestamp"`
	ExtraData  string                      `json:"extraData"`
	GasLimit   string                      `json:"gasLimit"`
	Difficulty string                      `json:"difficulty"`
	MixHash    string                      `json:"mixHash"`
	Coinbase   string                      `json:"coinbase"`
	Alloc      map[string]AllocatedAccount `json:"alloc"`
}

// AllocatedAccount represents a pre-funded account in the genesis
type AllocatedAccount struct {
	Balance string `json:"balance"`
}

// GetEthChainIDFromCosmosChainID derives the Ethereum chain ID from Cosmos chain ID
func GetEthChainIDFromCosmosChainID(cosmosChainID string) *big.Int {
	// Default mappings based on existing configuration
	switch cosmosChainID {
	case "mitosis-localnet-1":
		return big.NewInt(124899)
	case "mitosis-devnet-1":
		return big.NewInt(124864)
	default:
		// For custom chains, use a deterministic derivation
		// This is a simple hash-based approach - you may want to customize this
		return big.NewInt(100000)
	}
}

// GenerateEthereumGenesis creates an Ethereum genesis file
func GenerateEthereumGenesis(chainID string, outputPath string) error {
	ethChainID := GetEthChainIDFromCosmosChainID(chainID)

	// Create chain config with all EIPs enabled from genesis (similar to existing configs)
	zero := big.NewInt(0)
	zeroTime := uint64(0)

	chainConfig := &EthereumChainConfig{
		ChainID:                 ethChainID,
		HomesteadBlock:          zero,
		EIP150Block:             zero,
		EIP155Block:             zero,
		EIP158Block:             zero,
		ByzantiumBlock:          zero,
		ConstantinopleBlock:     zero,
		PetersburgBlock:         zero,
		IstanbulBlock:           zero,
		BerlinBlock:             zero,
		LondonBlock:             zero,
		TerminalTotalDifficulty: zero,
		ShanghaiTime:            &zeroTime,
		CancunTime:              &zeroTime,
		PragueTime:              &zeroTime,
	}

	// Default funded address (same as in existing genesis files)
	defaultFundedAddress := "0x2FB9C04d3225b55C964f9ceA934Cc8cD6070a3fF"

	// Initial balance: 10,000,000 ETH for localnet
	initialBalance := "10000000000000000000000000"
	if chainID == "mitosis-devnet-1" {
		// 999,000,000 ETH for devnet
		initialBalance = "999000000000000000000000000"
	}

	genesis := &EthereumGenesisSpec{
		Config:     chainConfig,
		Nonce:      "0x0",
		Timestamp:  "0x0",
		ExtraData:  "0x",
		GasLimit:   "0x1c9c380", // 30,000,000 gas limit
		Difficulty: "0x0",       // PoS from start
		MixHash:    "0x0000000000000000000000000000000000000000000000000000000000000000",
		Coinbase:   "0x0000000000000000000000000000000000000000",
		Alloc: map[string]AllocatedAccount{
			defaultFundedAddress: {
				Balance: initialBalance,
			},
		},
	}

	// Marshal to JSON with proper formatting
	jsonData, err := json.MarshalIndent(genesis, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal ethereum genesis: %w", err)
	}

	// Ensure the directory exists
	dir := filepath.Dir(outputPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create directory %s: %w", dir, err)
	}

	// Write the genesis file
	if err := os.WriteFile(outputPath, jsonData, 0600); err != nil {
		return fmt.Errorf("failed to write ethereum genesis file: %w", err)
	}

	return nil
}
