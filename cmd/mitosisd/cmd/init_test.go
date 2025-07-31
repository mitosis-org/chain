package cmd

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	genutiltypes "github.com/cosmos/cosmos-sdk/x/genutil/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewPrintInfoWithEth(t *testing.T) {
	moniker := "test-node"
	chainID := "test-chain-1"
	nodeID := "node123"
	genTxsDir := "/path/to/gentxs"
	appMessage := json.RawMessage(`{"test": "message"}`)
	ethGenesisPath := "/path/to/eth_genesis.json"
	ethChainID := "12345"

	info := newPrintInfoWithEth(moniker, chainID, nodeID, genTxsDir, appMessage, ethGenesisPath, ethChainID)

	assert.Equal(t, moniker, info.Moniker)
	assert.Equal(t, chainID, info.ChainID)
	assert.Equal(t, nodeID, info.NodeID)
	assert.Equal(t, genTxsDir, info.GenTxsDir)
	assert.Equal(t, appMessage, info.AppMessage)
	assert.Equal(t, ethGenesisPath, info.EthGenesis)
	assert.Equal(t, ethChainID, info.EthChainID)
}

func TestDisplayInfo(t *testing.T) {
	info := printInfo{
		Moniker:    "test-node",
		ChainID:    "test-chain",
		NodeID:     "node123",
		GenTxsDir:  "/gentxs",
		AppMessage: json.RawMessage(`{"msg": "test"}`),
		EthGenesis: "/eth_genesis.json",
		EthChainID: "12345",
	}

	// Redirect stderr to capture output
	oldStderr := os.Stderr
	r, w, _ := os.Pipe()
	os.Stderr = w

	err := displayInfo(info)
	require.NoError(t, err)

	// Restore stderr
	w.Close()
	os.Stderr = oldStderr

	// Read captured output
	out := make([]byte, 1024)
	n, _ := r.Read(out)

	// Verify JSON output contains expected fields
	var parsed printInfo
	err = json.Unmarshal(out[:n], &parsed)
	require.NoError(t, err)

	// Compare fields individually to avoid JSON formatting differences
	assert.Equal(t, info.Moniker, parsed.Moniker)
	assert.Equal(t, info.ChainID, parsed.ChainID)
	assert.Equal(t, info.NodeID, parsed.NodeID)
	assert.Equal(t, info.GenTxsDir, parsed.GenTxsDir)
	assert.Equal(t, info.EthGenesis, parsed.EthGenesis)
	assert.Equal(t, info.EthChainID, parsed.EthChainID)

	// For AppMessage, compare the content after parsing
	var expectedMsg, actualMsg map[string]interface{}
	err = json.Unmarshal(info.AppMessage, &expectedMsg)
	require.NoError(t, err)
	err = json.Unmarshal(parsed.AppMessage, &actualMsg)
	require.NoError(t, err)
	assert.Equal(t, expectedMsg, actualMsg)
}

func TestGenerateEthereumGenesis_Integration(t *testing.T) {
	// This tests the integration between init.go and eth_genesis.go
	tempDir := t.TempDir()

	tests := []struct {
		name            string
		chainID         string
		expectedEthID   int64
		expectedBalance string
	}{
		{
			name:            "localnet integration",
			chainID:         "mitosis-localnet-1",
			expectedEthID:   124899,
			expectedBalance: "999000000000000000000000000",
		},
		{
			name:            "devnet integration",
			chainID:         "mitosis-devnet-1",
			expectedEthID:   124864,
			expectedBalance: "999000000000000000000000000",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ethGenesisPath := filepath.Join(tempDir, tt.chainID, "eth_genesis.json")

			// Test the function that would be called from InitCmd
			err := GenerateEthereumGenesis(tt.chainID, ethGenesisPath)
			require.NoError(t, err)

			// Verify the file was created
			data, err := os.ReadFile(ethGenesisPath)
			require.NoError(t, err)

			var genesis EthereumGenesisSpec
			err = json.Unmarshal(data, &genesis)
			require.NoError(t, err)

			// Verify chain ID mapping
			assert.Equal(t, tt.expectedEthID, genesis.Config.ChainID.Int64())

			// Verify balance
			fundedAddr := DefaultFundedAddress
			account, exists := genesis.Alloc[fundedAddr]
			require.True(t, exists)
			assert.Equal(t, tt.expectedBalance, account.Balance)
		})
	}
}

func TestEthChainIDMapping(t *testing.T) {
	// Test the chain ID mapping used in InitCmd
	tests := []struct {
		cosmosChainID string
		expectedEthID int64
	}{
		{"mitosis-localnet-1", 124899},
		{"mitosis-devnet-1", 124864},
		{"custom-chain", 100000},
	}

	for _, tt := range tests {
		t.Run(tt.cosmosChainID, func(t *testing.T) {
			ethChainID := GetEthChainIDFromCosmosChainID(tt.cosmosChainID)
			assert.Equal(t, tt.expectedEthID, ethChainID.Int64())
		})
	}
}

func TestInitCmd_EthGenesisCreation(t *testing.T) {
	// Test that init command would create eth genesis in the right location
	tempDir := t.TempDir()
	configDir := filepath.Join(tempDir, "config")

	// Create config directory (simulating what InitCmd would do)
	err := os.MkdirAll(configDir, 0o755)
	require.NoError(t, err)

	// Test eth genesis generation
	ethGenesisPath := filepath.Join(configDir, "eth_genesis.json")
	err = GenerateEthereumGenesis("test-chain", ethGenesisPath)
	require.NoError(t, err)

	// Verify file exists at expected location
	_, err = os.Stat(ethGenesisPath)
	require.NoError(t, err)

	// Verify we can parse the genesis file
	genesis, err := genutiltypes.AppGenesisFromFile(filepath.Join(configDir, "genesis.json"))
	if err == nil {
		// If cosmos genesis exists, verify it has the right structure
		assert.NotNil(t, genesis)
	}
}
