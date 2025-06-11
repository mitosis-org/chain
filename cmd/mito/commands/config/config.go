package config

import (
	"fmt"

	"github.com/mitosis-org/chain/cmd/mito/internal/config"
	"github.com/mitosis-org/chain/cmd/mito/internal/utils"
	"github.com/spf13/cobra"
)

// NewConfigCmd returns the config command group
func NewConfigCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "config",
		Short: "Configuration commands",
		Long:  "Commands for managing configuration settings",
	}

	cmd.AddCommand(
		newSetRPCCmd(),
		newSetContractCmd(),
		newShowConfigCmd(),
	)

	return cmd
}

// newSetRPCCmd creates the set-rpc command
func newSetRPCCmd() *cobra.Command {
	var network string

	cmd := &cobra.Command{
		Use:   "set-rpc <rpc-url>",
		Short: "Set the default RPC URL",
		Long:  "Set the default RPC URL for connecting to the Ethereum network",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			rpcURL := args[0]

			cfg, err := config.Load()
			if err != nil {
				return fmt.Errorf("failed to load config: %w", err)
			}

			// Get current network config
			networkConfig := cfg.GetNetworkConfig(network)
			networkConfig.RPCURL = rpcURL

			// Set the updated network config
			cfg.SetNetworkConfig(network, networkConfig)

			if err := config.Save(cfg); err != nil {
				return fmt.Errorf("failed to save config: %w", err)
			}

			networkDisplay := network
			if networkDisplay == "" || networkDisplay == config.DefaultNetworkName {
				networkDisplay = config.DefaultNetworkName
			}
			fmt.Printf("✅ RPC URL set to: %s (network: %s)\n", rpcURL, networkDisplay)
			return nil
		},
	}

	cmd.Flags().StringVar(&network, "network", "", "Network name (defaults to '"+config.DefaultNetworkName+"')")

	return cmd
}

// newSetContractCmd creates the set-contract command
func newSetContractCmd() *cobra.Command {
	var network string

	cmd := &cobra.Command{
		Use:   "set-contract --validator-manager <contract-address>",
		Short: "Set the ValidatorManager contract address",
		Long:  "Set the ValidatorManager contract address for validator operations",
		RunE: func(cmd *cobra.Command, args []string) error {
			contractAddr, err := cmd.Flags().GetString("validator-manager")
			if err != nil {
				return fmt.Errorf("failed to get validator-manager flag: %w", err)
			}

			if contractAddr == "" {
				return fmt.Errorf("validator-manager address is required (use --validator-manager)")
			}

			// Validate address format
			if _, err := utils.ValidateAddress(contractAddr); err != nil {
				return fmt.Errorf("invalid contract address: %w", err)
			}

			cfg, err := config.Load()
			if err != nil {
				return fmt.Errorf("failed to load config: %w", err)
			}

			// Get current network config
			networkConfig := cfg.GetNetworkConfig(network)
			networkConfig.ValidatorManagerContractAddr = contractAddr

			// Set the updated network config
			cfg.SetNetworkConfig(network, networkConfig)

			if err := config.Save(cfg); err != nil {
				return fmt.Errorf("failed to save config: %w", err)
			}

			networkDisplay := network
			if networkDisplay == "" || networkDisplay == config.DefaultNetworkName {
				networkDisplay = config.DefaultNetworkName
			}
			fmt.Printf("✅ ValidatorManager contract address set to: %s (network: %s)\n", contractAddr, networkDisplay)
			return nil
		},
	}

	cmd.Flags().StringVar(&network, "network", "", "Network name (defaults to '"+config.DefaultNetworkName+"')")
	cmd.Flags().String("validator-manager", "", "ValidatorManager contract address")
	cmd.MarkFlagRequired("validator-manager")

	return cmd
}

// newShowConfigCmd creates the show command
func newShowConfigCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "show",
		Short: "Show current configuration",
		Long:  "Display the current configuration settings",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := config.Load()
			if err != nil {
				return fmt.Errorf("failed to load config: %w", err)
			}

			fmt.Println("===== Current Configuration =====")

			// Show default network
			fmt.Printf("\n[%s]\n", config.DefaultNetworkName)
			if cfg.Default.RPCURL != "" {
				fmt.Printf("rpc-url                                = %s\n", cfg.Default.RPCURL)
			} else {
				fmt.Printf("rpc-url                                = (not set)\n")
			}

			if cfg.Default.ValidatorManagerContractAddr != "" {
				fmt.Printf("validator-manager-contract-address     = %s\n", cfg.Default.ValidatorManagerContractAddr)
			} else {
				fmt.Printf("validator-manager-contract-address     = (not set)\n")
			}

			// Show other networks
			for _, networkName := range cfg.GetNetworkNames() {
				networkConfig := cfg.GetNetworkConfig(networkName)
				fmt.Printf("\n[%s]\n", networkName)
				if networkConfig.RPCURL != "" {
					fmt.Printf("rpc-url                                = %s\n", networkConfig.RPCURL)
				} else {
					fmt.Printf("rpc-url                                = (not set)\n")
				}

				if networkConfig.ValidatorManagerContractAddr != "" {
					fmt.Printf("validator-manager-contract-address     = %s\n", networkConfig.ValidatorManagerContractAddr)
				} else {
					fmt.Printf("validator-manager-contract-address     = (not set)\n")
				}
			}

			configPath, _ := config.GetPath()
			fmt.Printf("\nConfig file location: %s\n", configPath)

			return nil
		},
	}

	return cmd
}
