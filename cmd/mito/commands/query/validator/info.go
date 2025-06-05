package validator

import (
	"fmt"

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

	cmd := &cobra.Command{
		Use:   "info",
		Short: "Get validator information",
		Long:  "Retrieve detailed information about a validator",
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

			return runValidatorInfo(container, validatorAddr)
		},
	}

	// Add flags
	cmd.Flags().StringVar(&validatorAddress, "address", "", "Validator address to query")
	cmd.MarkFlagRequired("address")

	// Add network flags (no signing required for read-only operation)
	flags.AddNetworkFlags(cmd, &commonFlags)

	return cmd
}

func runValidatorInfo(container *container.Container, validatorAddr common.Address) error {
	// Check if validator exists
	isValidator, err := container.Contract.IsValidator(nil, validatorAddr)
	if err != nil {
		return fmt.Errorf("failed to check if address is validator: %w", err)
	}

	if !isValidator {
		fmt.Printf("Address %s is not a validator\n", validatorAddr.Hex())
		return nil
	}

	// Get validator info from contract
	validatorInfo, err := container.Contract.ValidatorInfo(nil, validatorAddr)
	if err != nil {
		return fmt.Errorf("failed to get validator info: %w", err)
	}

	// Display validator information
	fmt.Println("=== Validator Information ===")
	fmt.Printf("Address: %s\n", validatorInfo.ValAddr.Hex())
	fmt.Printf("Public Key: %x\n", validatorInfo.PubKey)
	fmt.Printf("Operator: %s\n", validatorInfo.Operator.Hex())
	fmt.Printf("Reward Manager: %s\n", validatorInfo.RewardManager.Hex())
	fmt.Printf("Commission Rate: %s\n", utils.FormatBasisPointsToPercent(validatorInfo.CommissionRate))
	fmt.Printf("Pending Commission Rate: %s\n", utils.FormatBasisPointsToPercent(validatorInfo.PendingCommissionRate))
	fmt.Printf("Pending Commission Rate Update Epoch: %s\n", validatorInfo.PendingCommissionRateUpdateEpoch.String())

	// Display metadata if available
	if len(validatorInfo.Metadata) > 0 {
		fmt.Printf("Metadata: %s\n", string(validatorInfo.Metadata))
	} else {
		fmt.Printf("Metadata: (none)\n")
	}

	return nil
}
