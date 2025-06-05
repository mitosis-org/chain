package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/pelletier/go-toml/v2"
)

// NetworkConfig represents configuration for a specific network
type NetworkConfig struct {
	RpcURL                       string `toml:"rpc-url"`
	ValidatorManagerContractAddr string `toml:"validator-manager-contract-address"`
}

// Config represents the overall configuration structure
type Config struct {
	Default  NetworkConfig
	networks map[string]NetworkConfig // internal storage
}

// getConfigPath returns the path to the config file
func getConfigPath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get home directory: %w", err)
	}

	configDir := filepath.Join(homeDir, ".mito")
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create config directory: %w", err)
	}

	return filepath.Join(configDir, "config.toml"), nil
}

// Load loads the configuration from file
func Load() (*Config, error) {
	configPath, err := getConfigPath()
	if err != nil {
		return nil, err
	}

	// If config file doesn't exist, return empty config
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return &Config{
			networks: make(map[string]NetworkConfig),
		}, nil
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	// Parse as raw interface{} first
	var rawData map[string]interface{}
	if err := toml.Unmarshal(data, &rawData); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	config := &Config{
		networks: make(map[string]NetworkConfig),
	}

	// Process each section
	for sectionName, sectionData := range rawData {
		if sectionMap, ok := sectionData.(map[string]interface{}); ok {
			networkConfig := NetworkConfig{}

			// Parse rpc-url
			if rpcURL, exists := sectionMap["rpc-url"]; exists {
				if rpcStr, ok := rpcURL.(string); ok {
					networkConfig.RpcURL = rpcStr
				}
			}

			// Parse validator-manager-contract-address
			if contractAddr, exists := sectionMap["validator-manager-contract-address"]; exists {
				if contractStr, ok := contractAddr.(string); ok {
					networkConfig.ValidatorManagerContractAddr = contractStr
				}
			}

			// Store in appropriate location
			if sectionName == "default" {
				config.Default = networkConfig
			} else {
				config.networks[sectionName] = networkConfig
			}
		}
	}

	return config, nil
}

// Save saves the configuration to file
func Save(config *Config) error {
	configPath, err := getConfigPath()
	if err != nil {
		return err
	}

	// Build raw data structure for flat TOML output
	rawData := make(map[string]interface{})

	// Add default section
	rawData["default"] = map[string]interface{}{
		"rpc-url":                            config.Default.RpcURL,
		"validator-manager-contract-address": config.Default.ValidatorManagerContractAddr,
	}

	// Add other networks as top-level sections
	for name, networkConfig := range config.networks {
		rawData[name] = map[string]interface{}{
			"rpc-url":                            networkConfig.RpcURL,
			"validator-manager-contract-address": networkConfig.ValidatorManagerContractAddr,
		}
	}

	data, err := toml.Marshal(rawData)
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	if err := os.WriteFile(configPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	return nil
}

// GetPath returns the configuration file path
func GetPath() (string, error) {
	return getConfigPath()
}

// GetNetworkConfig returns the configuration for a specific network
func (c *Config) GetNetworkConfig(networkName string) NetworkConfig {
	if networkName == "" || networkName == "default" {
		return c.Default
	}

	if network, exists := c.networks[networkName]; exists {
		return network
	}

	// Return empty config if network doesn't exist
	return NetworkConfig{}
}

// SetNetworkConfig sets the configuration for a specific network
func (c *Config) SetNetworkConfig(networkName string, config NetworkConfig) {
	if networkName == "" || networkName == "default" {
		c.Default = config
		return
	}

	if c.networks == nil {
		c.networks = make(map[string]NetworkConfig)
	}
	c.networks[networkName] = config
}

// GetNetworkNames returns all configured network names (excluding default)
func (c *Config) GetNetworkNames() []string {
	var names []string
	for name := range c.networks {
		names = append(names, name)
	}
	return names
}
