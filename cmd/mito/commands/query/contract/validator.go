package contract

import (
	"fmt"

	"github.com/mitosis-org/chain/cmd/mito/internal/config"
	"github.com/mitosis-org/chain/cmd/mito/internal/container"
	"github.com/mitosis-org/chain/cmd/mito/internal/flags"
	"github.com/mitosis-org/chain/cmd/mito/internal/units"
	"github.com/spf13/cobra"
)

// NewValidatorCmd returns the contract fee query command
func NewValidatorCmd() *cobra.Command {
	var commonFlags flags.CommonFlags

	cmd := &cobra.Command{
		Use:   "validator",
		Short: "Query current validator contracts",
		Long:  `Query current validator contracts from the ValidatorManager contract.`,
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

			return runContractQuery(container)
		},
	}

	// Add network flags (no signing required for read-only operation)
	flags.AddNetworkFlags(cmd, &commonFlags)
	return cmd
}

func runContractQuery(container *container.Container) error {
	return runValidatorManagerQuery(container)
}

func runValidatorManagerQuery(container *container.Container) error {
	// Get contract fee from the ValidatorManager contract
	contractFee, err := container.ValidatorManagerContract.Fee(nil)
	if err != nil {
		return fmt.Errorf("failed to get contract fee from blockchain: %w", err)
	}

	// Display contract fee information with aligned formatting
	fmt.Printf("%-20s %s\n", "address", container.ValidatorManagerContract.GetAddress())
	fmt.Printf("%-20s %s\n", "fee", units.FormatWeiToMito(contractFee))

	return nil
}
