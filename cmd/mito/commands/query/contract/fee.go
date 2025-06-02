package contract

import (
	"fmt"
	"math/big"

	"github.com/mitosis-org/chain/cmd/mito/internal/config"
	"github.com/mitosis-org/chain/cmd/mito/internal/container"
	"github.com/mitosis-org/chain/cmd/mito/internal/flags"
	"github.com/mitosis-org/chain/cmd/mito/internal/units"
	"github.com/spf13/cobra"
)

// NewFeeCmd returns the contract fee query command
func NewFeeCmd() *cobra.Command {
	var commonFlags flags.CommonFlags

	cmd := &cobra.Command{
		Use:   "fee",
		Short: "Query current contract fee",
		Long: `Query current contract fee from the ValidatorManager contract

This command retrieves the current fee required for validator operations
directly from the ValidatorManager contract on the blockchain.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			cmd.SilenceUsage = true

			// Resolve configuration
			resolver, err := config.NewResolver()
			if err != nil {
				return err
			}
			resolvedConfig := resolver.ResolveFlags(&commonFlags)

			// Validate required fields
			if resolvedConfig.RpcURL == "" {
				return fmt.Errorf("RPC URL is required (use --rpc-url or set with 'mito config set-rpc')")
			}
			if resolvedConfig.ValidatorManagerContractAddr == "" {
				return fmt.Errorf("ValidatorManager contract address is required (use --contract or set with 'mito config set-contract')")
			}

			// Create container
			container, err := container.NewContainer(resolvedConfig)
			if err != nil {
				return fmt.Errorf("failed to initialize container: %w", err)
			}
			defer container.Close()

			return runContractFeeQuery(container, resolvedConfig)
		},
	}

	// Add network flags (no signing required for read-only operation)
	flags.AddNetworkFlags(cmd, &commonFlags)

	return cmd
}

func runContractFeeQuery(container *container.Container, config *config.ResolvedConfig) error {
	fmt.Printf("Querying contract fee from ValidatorManager...\n")
	fmt.Printf("Contract Address: %s\n", config.ValidatorManagerContractAddr)
	fmt.Printf("RPC URL: %s\n\n", config.RpcURL)

	// Get contract fee from the ValidatorManager contract
	contractFee, err := container.Contract.Fee(nil)
	if err != nil {
		return fmt.Errorf("failed to get contract fee from blockchain: %w", err)
	}

	// Display contract fee information
	fmt.Println("=== Contract Fee Information (from blockchain) ===")
	fmt.Printf("Fee (wei): %s\n", contractFee.String())
	fmt.Printf("Fee (Mito): %s\n", units.FormatWeiToMito(contractFee))
	fmt.Printf("Fee (gwei): %s\n", units.FormatWeiToGwei(contractFee))
	fmt.Printf("Fee (both units): %s\n", units.FormatWeiToBothUnits(contractFee))

	if contractFee.Cmp(big.NewInt(0)) == 0 {
		fmt.Printf("\nNote: Contract fee is currently set to 0 (no fee required)\n")
	} else {
		fmt.Printf("\nNote: This fee is required for validator operations such as:\n")
		fmt.Printf("- Creating a new validator\n")
		fmt.Printf("- Unjailing a validator\n")
		fmt.Printf("- Transferring collateral ownership\n")
		fmt.Printf("- Withdrawing collateral\n")
	}

	return nil
}
