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
		newSetRpcCmd(),
		newSetContractCmd(),
		newShowConfigCmd(),
	)

	return cmd
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

			cfg, err := config.Load()
			if err != nil {
				return fmt.Errorf("failed to load config: %w", err)
			}

			cfg.RpcURL = rpcURL

			if err := config.Save(cfg); err != nil {
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
			if _, err := utils.ValidateAddress(contractAddr); err != nil {
				return fmt.Errorf("invalid contract address: %w", err)
			}

			cfg, err := config.Load()
			if err != nil {
				return fmt.Errorf("failed to load config: %w", err)
			}

			cfg.ValidatorManagerContractAddr = contractAddr

			if err := config.Save(cfg); err != nil {
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
			cfg, err := config.Load()
			if err != nil {
				return fmt.Errorf("failed to load config: %w", err)
			}

			fmt.Println("===== Current Configuration =====")
			if cfg.RpcURL != "" {
				fmt.Printf("RPC URL                      : %s\n", cfg.RpcURL)
			} else {
				fmt.Printf("RPC URL                      : (not set)\n")
			}

			if cfg.ValidatorManagerContractAddr != "" {
				fmt.Printf("ValidatorManager Contract    : %s\n", cfg.ValidatorManagerContractAddr)
			} else {
				fmt.Printf("ValidatorManager Contract    : (not set)\n")
			}

			configPath, _ := config.GetPath()
			fmt.Printf("Config file location         : %s\n", configPath)

			return nil
		},
	}

	return cmd
}
