package validator

import (
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/mitosis-org/chain/cmd/mito/internal/config"
	"github.com/mitosis-org/chain/cmd/mito/internal/container"
	"github.com/mitosis-org/chain/cmd/mito/internal/flags"
	"github.com/mitosis-org/chain/cmd/mito/internal/utils"
	"github.com/spf13/cobra"
)

func NewCollateralCmd() *cobra.Command {
	var commonFlags flags.CommonFlags
	var validatorAddress string
	var head uint64
	var tail uint64

	cmd := &cobra.Command{
		Use:   "collateral",
		Short: "Query validator collateral",
		Long:  "Query validator collateral",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			if head > 0 && tail > 0 {
				return fmt.Errorf("either --head or --tail must be specified, not both")
			}

			return nil
		},
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

			validatorAddr, err := utils.ValidateAddress(validatorAddress)
			if err != nil {
				return fmt.Errorf("invalid validator address: %w", err)
			}

			container, err := container.NewContainer(resolvedConfig)
			if err != nil {
				return fmt.Errorf("failed to initialize container: %w", err)
			}
			defer container.Close()

			if runValidatorCollateral(container, validatorAddr, head, tail) != nil {
				return fmt.Errorf("failed to run validator collateral: %w", err)
			}

			return nil
		},
	}

	cmd.Flags().StringVar(&validatorAddress, "address", "", "Validator address to query")
	cmd.Flags().Uint64Var(&head, "head", 0, "Top N permitted collateral owners to query")
	cmd.Flags().Uint64Var(&tail, "tail", 0, "Bottom N permitted collateral owners to query")
	cmd.MarkFlagRequired("address")

	return cmd
}

func runValidatorCollateral(container *container.Container, validatorAddr common.Address, head uint64, tail uint64) error {
	fmt.Println("=== Collateral Information ===")

	// Get permitted collateral owners count
	permittedOwnersCount, err := container.Contract.PermittedCollateralOwnerSize(nil, validatorAddr)
	if err != nil {
		return fmt.Errorf("failed to get permitted collateral owners count: %w", err)
	}

	fmt.Printf("Permitted Collateral Owners Count: %s\n", permittedOwnersCount.String())

	if permittedOwnersCount.Uint64() == 0 {
		return nil
	}

	if head > 0 {
		return runValidatorCollateralOwnerList(container, validatorAddr, permittedOwnersCount.Uint64(), head, false)
	} else if tail > 0 {
		return runValidatorCollateralOwnerList(container, validatorAddr, permittedOwnersCount.Uint64(), tail, true)
	}

	return nil
}

func runValidatorCollateralOwnerList(container *container.Container, validatorAddr common.Address, count uint64, target uint64, desc bool) error {
	if !desc {
		// head: 처음부터 target개
		fmt.Println("\n=== Top Permitted Collateral Owners ===")

		end := target
		if end > count {
			end = count
		}

		for i := uint64(0); i < end; i++ {
			owner, err := container.Contract.PermittedCollateralOwnerAt(nil, validatorAddr, big.NewInt(int64(i)))
			if err != nil {
				fmt.Printf("Error getting owner at index %d: %v\n", i, err)
				continue
			}
			fmt.Printf("%d: %s\n", i+1, owner.Hex())
		}
	} else {
		// tail: 마지막부터 target개 (역순)
		fmt.Println("\n=== Bottom Permitted Collateral Owners ===")

		printed := uint64(0)
		for i := int64(count - 1); i >= 0 && printed < target; i-- {
			owner, err := container.Contract.PermittedCollateralOwnerAt(nil, validatorAddr, big.NewInt(i))
			if err != nil {
				fmt.Printf("Error getting owner at index %d: %v\n", i, err)
				continue
			}
			fmt.Printf("%d: %s\n", i+1, owner.Hex())
			printed++
		}
	}
	return nil
}
