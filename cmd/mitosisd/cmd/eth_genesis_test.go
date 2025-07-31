package cmd

import (
	"encoding/json"
	"math/big"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetEthChainIDFromCosmosChainID(t *testing.T) {
	tests := []struct {
		name           string
		cosmosChainID  string
		expectedChainID *big.Int
	}{
		{
			name:           "localnet chain ID",
			cosmosChainID:  "mitosis-localnet-1",
			expectedChainID: big.NewInt(124899),
		},
		{
			name:           "devnet chain ID",
			cosmosChainID:  "mitosis-devnet-1",
			expectedChainID: big.NewInt(124864),
		},
		{
			name:           "custom chain ID",
			cosmosChainID:  "custom-chain-123",
			expectedChainID: big.NewInt(100000),
		},
		{
			name:           "empty chain ID",
			cosmosChainID:  "",
			expectedChainID: big.NewInt(100000),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GetEthChainIDFromCosmosChainID(tt.cosmosChainID)
			assert.Equal(t, tt.expectedChainID, result)
		})
	}
}

func TestGenerateEthereumGenesis(t *testing.T) {
	tempDir := t.TempDir()

	tests := []struct {
		name               string
		chainID            string
		outputPath         string
		expectedBalance    string
		expectedChainID    *big.Int
		expectError        bool
		validateGenesis    func(t *testing.T, genesis *EthereumGenesisSpec)
	}{
		{
			name:            "generate localnet genesis",
			chainID:         "mitosis-localnet-1",
			outputPath:      filepath.Join(tempDir, "localnet", "genesis.json"),
			expectedBalance: "10000000000000000000000000",
			expectedChainID: big.NewInt(124899),
			expectError:     false,
			validateGenesis: func(t *testing.T, genesis *EthereumGenesisSpec) {
				// Validate chain config
				assert.Equal(t, big.NewInt(124899), genesis.Config.ChainID)
				assert.Equal(t, big.NewInt(0), genesis.Config.HomesteadBlock)
				assert.Equal(t, big.NewInt(0), genesis.Config.ByzantiumBlock)
				assert.Equal(t, big.NewInt(0), genesis.Config.LondonBlock)
				assert.NotNil(t, genesis.Config.ShanghaiTime)
				assert.Equal(t, uint64(0), *genesis.Config.ShanghaiTime)
				assert.NotNil(t, genesis.Config.CancunTime)
				assert.Equal(t, uint64(0), *genesis.Config.CancunTime)
				
				// Validate genesis parameters
				assert.Equal(t, "0x0", genesis.Nonce)
				assert.Equal(t, "0x0", genesis.Timestamp)
				assert.Equal(t, "0x1c9c380", genesis.GasLimit)
				assert.Equal(t, "0x0", genesis.Difficulty)
				
				// Validate alloc
				fundedAddr := "0x2FB9C04d3225b55C964f9ceA934Cc8cD6070a3fF"
				account, exists := genesis.Alloc[fundedAddr]
				assert.True(t, exists)
				assert.Equal(t, "10000000000000000000000000", account.Balance)
			},
		},
		{
			name:            "generate devnet genesis",
			chainID:         "mitosis-devnet-1",
			outputPath:      filepath.Join(tempDir, "devnet", "genesis.json"),
			expectedBalance: "999000000000000000000000000",
			expectedChainID: big.NewInt(124864),
			expectError:     false,
			validateGenesis: func(t *testing.T, genesis *EthereumGenesisSpec) {
				// Validate chain ID
				assert.Equal(t, big.NewInt(124864), genesis.Config.ChainID)
				
				// Validate devnet specific balance
				fundedAddr := "0x2FB9C04d3225b55C964f9ceA934Cc8cD6070a3fF"
				account, exists := genesis.Alloc[fundedAddr]
				assert.True(t, exists)
				assert.Equal(t, "999000000000000000000000000", account.Balance)
			},
		},
		{
			name:            "generate custom chain genesis",
			chainID:         "custom-test-chain",
			outputPath:      filepath.Join(tempDir, "custom", "genesis.json"),
			expectedBalance: "10000000000000000000000000",
			expectedChainID: big.NewInt(100000),
			expectError:     false,
			validateGenesis: func(t *testing.T, genesis *EthereumGenesisSpec) {
				// Validate custom chain ID
				assert.Equal(t, big.NewInt(100000), genesis.Config.ChainID)
			},
		},
		{
			name:        "invalid output path",
			chainID:     "mitosis-localnet-1",
			outputPath:  "/invalid\x00path/genesis.json",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := GenerateEthereumGenesis(tt.chainID, tt.outputPath)

			if tt.expectError {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)

			// Verify file exists
			_, err = os.Stat(tt.outputPath)
			require.NoError(t, err)

			// Read and parse the generated file
			data, err := os.ReadFile(tt.outputPath)
			require.NoError(t, err)

			var genesis EthereumGenesisSpec
			err = json.Unmarshal(data, &genesis)
			require.NoError(t, err)

			// Run custom validation if provided
			if tt.validateGenesis != nil {
				tt.validateGenesis(t, &genesis)
			}
		})
	}
}

func TestEthereumChainConfig_AllFields(t *testing.T) {
	// Test that all chain config fields are properly set
	tempDir := t.TempDir()
	outputPath := filepath.Join(tempDir, "test_genesis.json")

	err := GenerateEthereumGenesis("test-chain", outputPath)
	require.NoError(t, err)

	data, err := os.ReadFile(outputPath)
	require.NoError(t, err)

	var genesis EthereumGenesisSpec
	err = json.Unmarshal(data, &genesis)
	require.NoError(t, err)

	config := genesis.Config

	// Verify all required blocks are set to 0
	assert.Equal(t, big.NewInt(0), config.HomesteadBlock)
	assert.Equal(t, big.NewInt(0), config.EIP150Block)
	assert.Equal(t, big.NewInt(0), config.EIP155Block)
	assert.Equal(t, big.NewInt(0), config.EIP158Block)
	assert.Equal(t, big.NewInt(0), config.ByzantiumBlock)
	assert.Equal(t, big.NewInt(0), config.ConstantinopleBlock)
	assert.Equal(t, big.NewInt(0), config.PetersburgBlock)
	assert.Equal(t, big.NewInt(0), config.IstanbulBlock)
	assert.Equal(t, big.NewInt(0), config.BerlinBlock)
	assert.Equal(t, big.NewInt(0), config.LondonBlock)
	assert.Equal(t, big.NewInt(0), config.TerminalTotalDifficulty)

	// Verify time-based forks
	assert.NotNil(t, config.ShanghaiTime)
	assert.Equal(t, uint64(0), *config.ShanghaiTime)
	assert.NotNil(t, config.CancunTime)
	assert.Equal(t, uint64(0), *config.CancunTime)
	assert.NotNil(t, config.PragueTime)
	assert.Equal(t, uint64(0), *config.PragueTime)

	// Verify optional fields are nil (not set)
	assert.Nil(t, config.DAOForkBlock)
	assert.False(t, config.DAOForkSupport)
	assert.Nil(t, config.ArrowGlacierBlock)
	assert.Nil(t, config.GrayGlacierBlock)
	assert.Nil(t, config.MergeNetsplitBlock)
}

func TestGenerateEthereumGenesis_DirectoryCreation(t *testing.T) {
	tempDir := t.TempDir()
	// Test nested directory creation
	nestedPath := filepath.Join(tempDir, "a", "b", "c", "genesis.json")

	err := GenerateEthereumGenesis("test-chain", nestedPath)
	require.NoError(t, err)

	// Verify file exists
	_, err = os.Stat(nestedPath)
	require.NoError(t, err)
}

func TestGenerateEthereumGenesis_JSONFormat(t *testing.T) {
	tempDir := t.TempDir()
	outputPath := filepath.Join(tempDir, "genesis.json")

	err := GenerateEthereumGenesis("mitosis-localnet-1", outputPath)
	require.NoError(t, err)

	// Read the file
	data, err := os.ReadFile(outputPath)
	require.NoError(t, err)

	// Verify it's valid JSON
	var rawJSON map[string]interface{}
	err = json.Unmarshal(data, &rawJSON)
	require.NoError(t, err)

	// Verify proper indentation (should contain newlines and spaces)
	assert.Contains(t, string(data), "\n")
	assert.Contains(t, string(data), "  ")
}