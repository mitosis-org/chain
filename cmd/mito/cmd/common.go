package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/mitosis-org/chain/bindings"
)

// Common flags used across multiple commands
var (
	rpcURL                       string
	privateKey                   string
	keyfilePath                  string
	keyfilePassword              string
	keyfilePasswordFile          string
	validatorManagerContractAddr string
	yes                          bool
	nonce                        uint64
	nonceSpecified               bool
	outputFile                   string
	signed                       bool
	unsigned                     bool

	// Network information flags
	chainID     string
	gasLimit    uint64
	gasPrice    string
	txNonce     uint64
	contractFee string

	// Shared client and contract instances
	ethClient *ethclient.Client
	contract  *bindings.IValidatorManager

	// Global config
	globalConfig *Config
)

// ConfirmAction prompts the user to confirm an action
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

// loadGlobalConfig loads the global configuration
func loadGlobalConfig() error {
	var err error
	globalConfig, err = loadConfig()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}
	return nil
}

// resolveConfigValues resolves configuration values with command line flags taking precedence
func resolveConfigValues() {
	if rpcURL == "" && globalConfig.RpcURL != "" {
		rpcURL = globalConfig.RpcURL
	}
	if validatorManagerContractAddr == "" && globalConfig.ValidatorManagerContractAddr != "" {
		validatorManagerContractAddr = globalConfig.ValidatorManagerContractAddr
	}
}
