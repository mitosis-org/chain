package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

// Config represents the configuration structure
type Config struct {
	RpcURL                       string `json:"rpc_url"`
	ValidatorManagerContractAddr string `json:"validator_manager_contract_addr"`
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

	return filepath.Join(configDir, "config.json"), nil
}

// loadConfig loads the configuration from file
func loadConfig() (*Config, error) {
	configPath, err := getConfigPath()
	if err != nil {
		return nil, err
	}

	// If config file doesn't exist, return empty config
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return &Config{}, nil
	}

	data, err := ioutil.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var config Config
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	return &config, nil
}

// saveConfig saves the configuration to file
func saveConfig(config *Config) error {
	configPath, err := getConfigPath()
	if err != nil {
		return err
	}

	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	if err := ioutil.WriteFile(configPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	return nil
}

// GetConfigValue returns a config value, with command line flag taking precedence
func GetConfigValue(flagValue, configValue string) string {
	if flagValue != "" {
		return flagValue
	}
	return configValue
}

// newSetRpcCmd creates the set-rpc command
func newSetRpcCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "set-rpc <rpc-url>",
		Short: "Set the default RPC URL",
		Long:  "Set the default RPC URL for connecting to the Ethereum network",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			rpcURL := args[0]

			config, err := loadConfig()
			if err != nil {
				return fmt.Errorf("failed to load config: %w", err)
			}

			config.RpcURL = rpcURL

			if err := saveConfig(config); err != nil {
				return fmt.Errorf("failed to save config: %w", err)
			}

			fmt.Printf("✅ RPC URL set to: %s\n", rpcURL)
			return nil
		},
	}

	return cmd
}

// newSetContractCmd creates the set-contract command
func newSetContractCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "set-contract <contract-address>",
		Short: "Set the default ValidatorManager contract address",
		Long:  "Set the default ValidatorManager contract address for validator operations",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			contractAddr := args[0]

			// Validate address format
			if _, err := ValidateAddress(contractAddr); err != nil {
				return fmt.Errorf("invalid contract address: %w", err)
			}

			config, err := loadConfig()
			if err != nil {
				return fmt.Errorf("failed to load config: %w", err)
			}

			config.ValidatorManagerContractAddr = contractAddr

			if err := saveConfig(config); err != nil {
				return fmt.Errorf("failed to save config: %w", err)
			}

			fmt.Printf("✅ ValidatorManager contract address set to: %s\n", contractAddr)
			return nil
		},
	}

	return cmd
}

// newShowConfigCmd creates the show command
func newShowConfigCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "show",
		Short: "Show current configuration",
		Long:  "Display the current configuration settings",
		RunE: func(cmd *cobra.Command, args []string) error {
			config, err := loadConfig()
			if err != nil {
				return fmt.Errorf("failed to load config: %w", err)
			}

			fmt.Println("===== Current Configuration =====")
			if config.RpcURL != "" {
				fmt.Printf("RPC URL                      : %s\n", config.RpcURL)
			} else {
				fmt.Printf("RPC URL                      : (not set)\n")
			}

			if config.ValidatorManagerContractAddr != "" {
				fmt.Printf("ValidatorManager Contract    : %s\n", config.ValidatorManagerContractAddr)
			} else {
				fmt.Printf("ValidatorManager Contract    : (not set)\n")
			}

			configPath, _ := getConfigPath()
			fmt.Printf("Config file location         : %s\n", configPath)

			return nil
		},
	}

	return cmd
}
