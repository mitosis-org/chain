package validator

import (
	"encoding/json"
	"fmt"
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

	// Display metadata if available
	if len(validatorInfo.Metadata) > 0 {
		// Try to parse and format as JSON
		var jsonData interface{}
		if err := json.Unmarshal(validatorInfo.Metadata, &jsonData); err == nil {
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
			} else {
				// Fallback to raw string if formatting fails
				fmt.Printf("%-30s %s\n", "Metadata", string(validatorInfo.Metadata))
			}
		} else {
			// Not valid JSON, display as raw string
			fmt.Printf("%-30s %s\n", "Metadata", string(validatorInfo.Metadata))
		}
	} else {
		fmt.Printf("%-30s (none)\n", "Metadata")
	}

	return nil
}
