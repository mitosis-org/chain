package cmd

import (
	"encoding/json"
	"fmt"
	"math/big"
	"os"
	"path/filepath"
)

// DefaultFundedAddress is the default Ethereum address that receives initial funding
const DefaultFundedAddress = "0xF530AC32044B7bCF6B6ac9E2B65d8eDb7794d64f"

// BlobConfig represents blob configuration for specific network upgrades
type BlobConfig struct {
	Target                int `json:"target"`
	Max                   int `json:"max"`
	BaseFeeUpdateFraction int `json:"baseFeeUpdateFraction"`
}

// BlobSchedule represents the blob schedule configuration
type BlobSchedule struct {
	Cancun *BlobConfig `json:"cancun,omitempty"`
	Prague *BlobConfig `json:"prague,omitempty"`
}

// EthereumChainConfig represents the chain configuration for Ethereum genesis
type EthereumChainConfig struct {
	ChainID                 *big.Int      `json:"chainId"`
	HomesteadBlock          *big.Int      `json:"homesteadBlock"`
	DAOForkBlock            *big.Int      `json:"daoForkBlock,omitempty"`
	DAOForkSupport          bool          `json:"daoForkSupport,omitempty"`
	EIP150Block             *big.Int      `json:"eip150Block"`
	EIP155Block             *big.Int      `json:"eip155Block"`
	EIP158Block             *big.Int      `json:"eip158Block"`
	ByzantiumBlock          *big.Int      `json:"byzantiumBlock"`
	ConstantinopleBlock     *big.Int      `json:"constantinopleBlock"`
	PetersburgBlock         *big.Int      `json:"petersburgBlock"`
	IstanbulBlock           *big.Int      `json:"istanbulBlock"`
	BerlinBlock             *big.Int      `json:"berlinBlock"`
	LondonBlock             *big.Int      `json:"londonBlock"`
	ArrowGlacierBlock       *big.Int      `json:"arrowGlacierBlock,omitempty"`
	GrayGlacierBlock        *big.Int      `json:"grayGlacierBlock,omitempty"`
	MergeNetsplitBlock      *big.Int      `json:"mergeNetsplitBlock,omitempty"`
	TerminalTotalDifficulty *big.Int      `json:"terminalTotalDifficulty"`
	ShanghaiTime            *uint64       `json:"shanghaiTime"`
	CancunTime              *uint64       `json:"cancunTime"`
	PragueTime              *uint64       `json:"pragueTime,omitempty"`
	BlobSchedule            *BlobSchedule `json:"blobSchedule,omitempty"`
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

// EthGenesisOptions contains options for generating Ethereum genesis
type EthGenesisOptions struct {
	ChainID        string
	OutputPath     string
	EthChainID     *big.Int
	GasLimit       uint64
	FundedAddress  string
	InitialBalance string
}

// GenerateEthereumGenesis creates an Ethereum genesis file
func GenerateEthereumGenesis(chainID string, outputPath string) error {
	return GenerateEthereumGenesisWithOptions(EthGenesisOptions{
		ChainID:    chainID,
		OutputPath: outputPath,
	})
}

// GenerateEthereumGenesisWithOptions creates an Ethereum genesis file with custom options
func GenerateEthereumGenesisWithOptions(opts EthGenesisOptions) error {
	// Determine Ethereum chain ID
	var ethChainID *big.Int
	if opts.EthChainID != nil {
		ethChainID = opts.EthChainID
	} else {
		ethChainID = GetEthChainIDFromCosmosChainID(opts.ChainID)
	}

	// Create chain config with all EIPs enabled from genesis
	zero := big.NewInt(0)
	zeroTime := uint64(0)

	// Create blob schedule configuration
	blobSchedule := &BlobSchedule{
		Cancun: &BlobConfig{
			Target:                3,
			Max:                   6,
			BaseFeeUpdateFraction: 3338477,
		},
		Prague: &BlobConfig{
			Target:                6,
			Max:                   9,
			BaseFeeUpdateFraction: 5007716,
		},
	}

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
		BlobSchedule:            blobSchedule,
	}

	// Determine funded address
	fundedAddress := opts.FundedAddress
	if fundedAddress == "" {
		fundedAddress = DefaultFundedAddress
	}

	// Determine initial balance
	initialBalance := opts.InitialBalance
	if initialBalance == "" {
		// Default balance: 999,000,000 ETH
		initialBalance = "999000000000000000000000000"
	}

	// Determine gas limit
	gasLimit := opts.GasLimit
	if gasLimit == 0 {
		gasLimit = 30000000 // Default: 30,000,000
	}
	gasLimitStr := fmt.Sprintf("%d", gasLimit)

	genesis := &EthereumGenesisSpec{
		Config:     chainConfig,
		Nonce:      "0",
		Timestamp:  "0",
		ExtraData:  "0x",
		GasLimit:   gasLimitStr,
		Difficulty: "0",
		Alloc: map[string]AllocatedAccount{
			fundedAddress: {
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
	dir := filepath.Dir(opts.OutputPath)
	if err := os.MkdirAll(dir, 0o750); err != nil {
		return fmt.Errorf("failed to create directory %s: %w", dir, err)
	}

	// Write the genesis file
	if err := os.WriteFile(opts.OutputPath, jsonData, 0o600); err != nil {
		return fmt.Errorf("failed to write ethereum genesis file: %w", err)
	}

	return nil
}
