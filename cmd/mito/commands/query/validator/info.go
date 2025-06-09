package validator

import (
	"encoding/json"
	"fmt"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"github.com/mitosis-org/chain/cmd/mito/internal/config"
	"github.com/mitosis-org/chain/cmd/mito/internal/container"
	"github.com/mitosis-org/chain/cmd/mito/internal/flags"
	"github.com/mitosis-org/chain/cmd/mito/internal/utils"
	"github.com/spf13/cobra"
)

// NewInfoCmd creates the validator info command
func NewInfoCmd() *cobra.Command {
	var commonFlags flags.CommonFlags
	var validatorAddress string
	var head uint64
	var tail uint64

	cmd := &cobra.Command{
		Use:   "info",
		Short: "Get validator information",
		Long:  "Retrieve detailed information about a validator including collateral information",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			if head > 0 && tail > 0 {
				return fmt.Errorf("either --head or --tail must be specified, not both")
			}
			return nil
		},
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

			// Validate validator address
			validatorAddr, err := utils.ValidateAddress(validatorAddress)
			if err != nil {
				return fmt.Errorf("invalid validator address: %w", err)
			}

			// Create container
			container, err := container.NewContainer(resolvedConfig)
			if err != nil {
				return fmt.Errorf("failed to initialize container: %w", err)
			}
			defer container.Close()

			return runValidatorInfo(container, validatorAddr, head, tail)
		},
	}

	// Add flags
	cmd.Flags().StringVar(&validatorAddress, "address", "", "Validator address to query")
	cmd.Flags().Uint64Var(&head, "head", 0, "Top N permitted collateral owners to query")
	cmd.Flags().Uint64Var(&tail, "tail", 0, "Bottom N permitted collateral owners to query")
	cmd.MarkFlagRequired("address")

	// Add network flags (no signing required for read-only operation)
	flags.AddNetworkFlags(cmd, &commonFlags)

	return cmd
}

func runValidatorInfo(container *container.Container, validatorAddr common.Address, head uint64, tail uint64) error {
	// Check if validator exists
	isValidator, err := container.ValidatorManagerContract.IsValidator(nil, validatorAddr)
	if err != nil {
		return fmt.Errorf("failed to check if address is validator: %w", err)
	}

	if !isValidator {
		fmt.Printf("Address %s is not a validator\n", validatorAddr.Hex())
		return nil
	}

	// Get validator info from contract
	validatorInfo, err := container.ValidatorManagerContract.ValidatorInfo(nil, validatorAddr)
	if err != nil {
		return fmt.Errorf("failed to get validator info: %w", err)
	}

	// Display validator information
	fmt.Printf("%-30s %s\n", "Address", validatorInfo.ValAddr.Hex())
	fmt.Printf("%-30s %x\n", "Public Key", validatorInfo.PubKey)
	fmt.Printf("%-30s %s\n", "Operator", validatorInfo.Operator.Hex())
	fmt.Printf("%-30s %s\n", "Reward Manager", validatorInfo.RewardManager.Hex())
	fmt.Printf("%-30s %s\n", "Commission Rate", utils.FormatBasisPointsToPercent(validatorInfo.CommissionRate))
	fmt.Printf("%-30s %s\n", "Pending Commission Rate", utils.FormatBasisPointsToPercent(validatorInfo.PendingCommissionRate))
	fmt.Printf("%-30s %s\n", "Update Epoch", validatorInfo.PendingCommissionRateUpdateEpoch.String())

	// Display metadata
	displayMetadata(validatorInfo.Metadata)

	// Display collateral information
	return runValidatorCollateral(container, validatorAddr, head, tail)
}

func displayMetadata(metadata []byte) {
	if len(metadata) == 0 {
		fmt.Printf("%-30s (none)\n", "Metadata")
		return
	}

	// Try to parse and format as JSON
	var jsonData interface{}
	if err := json.Unmarshal(metadata, &jsonData); err == nil {
		// Successfully parsed as JSON, format it nicely
		formatted, err := json.MarshalIndent(jsonData, "", "  ")
		if err == nil {
			lines := strings.Split(string(formatted), "\n")
			if len(lines) > 0 {
				// First line with label
				fmt.Printf("%-30s %s\n", "Metadata", lines[0])
				// Remaining lines with proper indentation
				for i := 1; i < len(lines); i++ {
					fmt.Printf("%-30s %s\n", "", lines[i])
				}
			}
			return
		}
	}

	// Fallback to raw string if JSON parsing/formatting fails
	fmt.Printf("%-30s %s\n", "Metadata", string(metadata))
}

func runValidatorCollateral(container *container.Container, validatorAddr common.Address, head uint64, tail uint64) error {
	// Add collateral information
	fmt.Println() // Empty line for separation

	// Get permitted collateral owners count
	permittedOwnersCount, err := container.ValidatorManagerContract.PermittedCollateralOwnerSize(nil, validatorAddr)
	if err != nil {
		return fmt.Errorf("failed to get permitted collateral owners count: %w", err)
	}

	fmt.Printf("%-30s %s\n", "Collateral Owners Count", permittedOwnersCount.String())

	if permittedOwnersCount.Uint64() > 0 {
		if head > 0 {
			fmt.Println()
			fmt.Printf("Top %d Collateral Owners:\n", head)
			fmt.Println("========================")
			return runValidatorCollateralOwnerList(container, validatorAddr, permittedOwnersCount.Uint64(), head, false)
		} else if tail > 0 {
			fmt.Println()
			fmt.Printf("Bottom %d Collateral Owners:\n", tail)
			fmt.Println("===========================")
			return runValidatorCollateralOwnerList(container, validatorAddr, permittedOwnersCount.Uint64(), tail, true)
		}
	}

	return nil
}

func runValidatorCollateralOwnerList(container *container.Container, validatorAddr common.Address, count uint64, target uint64, desc bool) error {
	if !desc {
		end := target
		if end > count {
			end = count
		}

		for i := uint64(0); i < end; i++ {
			owner, err := container.ValidatorManagerContract.PermittedCollateralOwnerAt(nil, validatorAddr, big.NewInt(int64(i)))
			if err != nil {
				fmt.Printf("Error getting owner at index %d: %v\n", i, err)
				continue
			}
			fmt.Printf("%-30s %s\n", fmt.Sprintf("%d:", i+1), owner.Hex())
		}
	} else {
		printed := uint64(0)
		for i := int64(count - 1); i >= 0 && printed < target; i-- {
			owner, err := container.ValidatorManagerContract.PermittedCollateralOwnerAt(nil, validatorAddr, big.NewInt(i))
			if err != nil {
				fmt.Printf("Error getting owner at index %d: %v\n", i, err)
				continue
			}
			fmt.Printf("%-20s %s\n", fmt.Sprintf("%d:", i+1), owner.Hex())
			printed++
		}
	}
	return nil
}
