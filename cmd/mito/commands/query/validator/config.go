package validator

import (
	"fmt"

	"github.com/mitosis-org/chain/cmd/mito/internal/config"
	"github.com/mitosis-org/chain/cmd/mito/internal/container"
	"github.com/mitosis-org/chain/cmd/mito/internal/flags"
	"github.com/mitosis-org/chain/cmd/mito/internal/utils"
	"github.com/spf13/cobra"
)

func NewConfigCmd() *cobra.Command {
	var commonFlags flags.CommonFlags

	cmd := &cobra.Command{
		Use:   "config",
		Short: "Query validator config",
		Long:  "Query validator config",
		RunE: func(cmd *cobra.Command, args []string) error {
			resolver, err := config.NewResolver()
			if err != nil {
				return err
			}
			resolvedConfig := resolver.ResolveFlags(&commonFlags)

			if resolvedConfig.RpcURL == "" {
				return fmt.Errorf("RPC URL is required (use --rpc-url or set with 'mito config set-rpc')")
			}
			if resolvedConfig.ValidatorManagerContractAddr == "" {
				return fmt.Errorf("ValidatorManager contract address is required (use --contract or set with 'mito config set-contract')")
			}

			container, err := container.NewContainer(resolvedConfig)
			if err != nil {
				return fmt.Errorf("failed to initialize container: %w", err)
			}
			defer container.Close()
			return runValidatorConfig(container)
		},
	}

	flags.AddNetworkFlags(cmd, &commonFlags)

	return cmd
}

func runValidatorConfig(container *container.Container) error {
	// Get global validator config
	globalConfig, err := container.ValidatorManagerContract.GlobalValidatorConfig(nil)
	if err != nil {
		return fmt.Errorf("failed to get global validator config: %w", err)
	}

	fmt.Printf("%-30s %s\n", "Initial Validator Deposit", utils.FormatWeiToEther(globalConfig.InitialValidatorDeposit))
	fmt.Printf("%-30s %s\n", "Minimum Commission Rate", utils.FormatBasisPointsToPercent(globalConfig.MinimumCommissionRate))
	fmt.Printf("%-30s %s\n", "Commission Rate Update Delay", globalConfig.CommissionRateUpdateDelayEpoch.String())
	fmt.Printf("%-30s %s\n", "Collateral Withdrawal Delay", globalConfig.CollateralWithdrawalDelaySeconds.String())

	return nil
}
