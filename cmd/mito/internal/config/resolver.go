package config

import (
	"math/big"

	"github.com/mitosis-org/chain/cmd/mito/internal/flags"
)

// Resolver resolves configuration values from multiple sources
type Resolver struct {
	config *Config
}

// NewResolver creates a new config resolver
func NewResolver() (*Resolver, error) {
	config, err := Load()
	if err != nil {
		return nil, err
	}
	return &Resolver{config: config}, nil
}

// ResolveFlags resolves flag values with config file fallback
func (r *Resolver) ResolveFlags(commonFlags *flags.CommonFlags) *ResolvedConfig {
	// Get network-specific configuration
	networkConfig := r.config.GetNetworkConfig(commonFlags.Network)

	return &ResolvedConfig{
		Network:                      commonFlags.Network,
		RpcURL:                       resolveValue(commonFlags.RpcURL, networkConfig.RpcURL),
		ValidatorManagerContractAddr: resolveValue(commonFlags.ValidatorManagerContractAddr, networkConfig.ValidatorManagerContractAddr),
		ChainID:                      commonFlags.ChainID,
		PrivateKey:                   commonFlags.PrivateKey,
		KeyfilePath:                  commonFlags.KeyfilePath,
		KeyfilePassword:              commonFlags.KeyfilePassword,
		KeyfilePasswordFile:          commonFlags.KeyfilePasswordFile,
		GasLimit:                     commonFlags.GasLimit,
		GasPrice:                     commonFlags.GasPrice,
		Nonce:                        commonFlags.Nonce,
		ContractFee:                  big.NewInt(0),
		OutputFile:                   commonFlags.OutputFile,
		Signed:                       commonFlags.Signed,
		Unsigned:                     commonFlags.Unsigned,
		Yes:                          commonFlags.Yes,
	}
}

// ResolvedConfig contains resolved configuration values
type ResolvedConfig struct {
	Network                      string
	RpcURL                       string
	ValidatorManagerContractAddr string
	ChainID                      string
	PrivateKey                   string
	KeyfilePath                  string
	KeyfilePassword              string
	KeyfilePasswordFile          string
	GasLimit                     uint64
	GasPrice                     string
	Nonce                        string
	ContractFee                  *big.Int
	OutputFile                   string
	Signed                       bool
	Unsigned                     bool
	Yes                          bool
}

// resolveValue returns flag value if set, otherwise returns config value
func resolveValue(flagValue, configValue string) string {
	if flagValue != "" {
		return flagValue
	}
	return configValue
}

// HasSigningMethod checks if any signing method is configured
func (r *ResolvedConfig) HasSigningMethod() bool {
	return r.PrivateKey != "" || r.KeyfilePath != ""
}

// GetSigningMethod returns the configured signing method
func (r *ResolvedConfig) GetSigningMethod() string {
	if r.PrivateKey != "" {
		return "private-key"
	}
	if r.KeyfilePath != "" {
		return "keyfile"
	}
	return ""
}
