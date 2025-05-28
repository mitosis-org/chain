package cmd

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
)

// ValidateNetworkInfo validates network information requirements for tx create commands
func ValidateNetworkInfo(cmd *cobra.Command) error {
	// Check if RPC URL is available (from config or flag)
	hasRPC := rpcURL != ""

	// Check if required network info flags are provided
	hasChainID := cmd.Flags().Changed("chain-id")
	hasGasPrice := cmd.Flags().Changed("gas-price")
	hasGasLimit := cmd.Flags().Changed("gas-limit")
	hasFee := cmd.Flags().Changed("fee")

	// If RPC is available, we can fetch network info automatically
	if hasRPC {
		return nil
	}

	// If no RPC, check if all required network info is provided
	missingFlags := []string{}

	if !hasChainID {
		missingFlags = append(missingFlags, "--chain-id")
	}
	if !hasGasPrice {
		missingFlags = append(missingFlags, "--gas-price")
	}
	if !hasGasLimit {
		missingFlags = append(missingFlags, "--gas-limit")
	}

	// For commands that require fee, check if fee is provided
	if cmd.Flags().Lookup("fee") != nil && !hasFee {
		missingFlags = append(missingFlags, "--fee")
	}

	if len(missingFlags) > 0 {
		return fmt.Errorf("no RPC URL configured and missing required network information: %v\n"+
			"Either configure RPC URL with 'mito config set-rpc <url>' or provide the missing flags", missingFlags)
	}

	return nil
}

// SetupNetworkInfo sets up network information either from RPC or from provided flags
func SetupNetworkInfo(cmd *cobra.Command) error {
	// If RPC is available, fetch network info
	if rpcURL != "" {
		ethClient, err := GetEthClient(rpcURL)
		if err != nil {
			return fmt.Errorf("failed to connect to RPC: %w", err)
		}

		// Get chain ID if not provided
		if !cmd.Flags().Changed("chain-id") {
			chainIDBig, err := ethClient.ChainID(context.Background())
			if err != nil {
				return fmt.Errorf("failed to get chain ID from RPC: %w", err)
			}
			chainID = chainIDBig.String()
		}

		// Get gas price if not provided
		if !cmd.Flags().Changed("gas-price") {
			gasPriceBig, err := ethClient.SuggestGasPrice(context.Background())
			if err != nil {
				return fmt.Errorf("failed to get gas price from RPC: %w", err)
			}
			gasPrice = gasPriceBig.String()
		}

		// For commands with contract interaction, get fee if not provided
		if cmd.Flags().Lookup("fee") != nil && !cmd.Flags().Changed("fee") {
			contract, err := GetValidatorManagerContract(ethClient)
			if err != nil {
				return fmt.Errorf("failed to initialize contract: %w", err)
			}

			fee, err := contract.Fee(nil)
			if err != nil {
				return fmt.Errorf("failed to get contract fee from RPC: %w", err)
			}
			contractFee = FormatWeiToEther(fee)
		}

		fmt.Printf("Network info fetched from RPC:\n")
		fmt.Printf("  Chain ID: %s\n", chainID)
		fmt.Printf("  Gas Price: %s wei\n", gasPrice)
		if cmd.Flags().Lookup("fee") != nil {
			fmt.Printf("  Contract Fee: %s MITO\n", contractFee)
		}
		fmt.Println()
	}

	return nil
}

// AddTxCreateNetworkFlags adds network-related flags for tx create commands
func AddTxCreateNetworkFlags(cmd *cobra.Command, requiresFee bool) {
	cmd.Flags().StringVar(&chainID, "chain-id", "", "Chain ID for the transaction (auto-fetched if RPC available)")
	cmd.Flags().Uint64Var(&gasLimit, "gas-limit", 500000, "Gas limit for the transaction")
	cmd.Flags().StringVar(&gasPrice, "gas-price", "", "Gas price in wei (auto-fetched if RPC available)")
	cmd.Flags().Uint64Var(&txNonce, "nonce", 0, "Nonce for the transaction (required for signed transactions)")

	if requiresFee {
		cmd.Flags().StringVar(&contractFee, "fee", "", "Contract fee in MITO (auto-fetched if RPC available)")
	}
}
