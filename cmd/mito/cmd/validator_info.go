package cmd

import (
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/spf13/cobra"
)

var (
	validatorAddress string
)

// newValidatorInfoCmd creates the validator info command
func newValidatorInfoCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "info",
		Short: "Get validator information",
		Long:  "Retrieve detailed information about a validator",
		RunE:  runValidatorInfo,
	}

	// Add flags
	cmd.Flags().StringVar(&validatorAddress, "validator-address", "", "Validator address to query")
	mustMarkFlagRequired(cmd, "validator-address")

	// Add common flags (no signing required for read-only operation)
	AddCommonFlags(cmd, false)

	return cmd
}

func runValidatorInfo(cmd *cobra.Command, args []string) error {
	// Validate validator address
	validatorAddr, err := ValidateAddress(validatorAddress)
	if err != nil {
		return fmt.Errorf("invalid validator address: %w", err)
	}

	fmt.Printf("Querying validator information for: %s\n", validatorAddr.Hex())
	fmt.Printf("Using contract: %s\n", validatorManagerContractAddr)
	fmt.Printf("RPC URL: %s\n\n", rpcURL)

	// Check if validator exists
	isValidator, err := contract.IsValidator(nil, validatorAddr)
	if err != nil {
		return fmt.Errorf("failed to check if address is validator: %w", err)
	}

	if !isValidator {
		fmt.Printf("Address %s is not a validator\n", validatorAddr.Hex())
		return nil
	}

	// Get validator info from contract
	validatorInfo, err := contract.ValidatorInfo(nil, validatorAddr)
	if err != nil {
		return fmt.Errorf("failed to get validator info: %w", err)
	}

	// Display validator information
	fmt.Println("=== Validator Information ===")
	fmt.Printf("Address: %s\n", validatorInfo.ValAddr.Hex())
	fmt.Printf("Public Key: %x\n", validatorInfo.PubKey)
	fmt.Printf("Operator: %s\n", validatorInfo.Operator.Hex())
	fmt.Printf("Reward Manager: %s\n", validatorInfo.RewardManager.Hex())
	fmt.Printf("Commission Rate: %s\n", FormatBasisPointsToPercent(validatorInfo.CommissionRate))
	fmt.Printf("Pending Commission Rate: %s\n", FormatBasisPointsToPercent(validatorInfo.PendingCommissionRate))
	fmt.Printf("Pending Commission Rate Update Epoch: %s\n", validatorInfo.PendingCommissionRateUpdateEpoch.String())

	// Display metadata if available
	if len(validatorInfo.Metadata) > 0 {
		fmt.Printf("Metadata: %s\n", string(validatorInfo.Metadata))
	} else {
		fmt.Printf("Metadata: (none)\n")
	}

	// Get additional information
	if err := displayAdditionalInfo(validatorAddr); err != nil {
		fmt.Printf("\nWarning: Could not retrieve additional information: %v\n", err)
	}

	return nil
}

func displayAdditionalInfo(validatorAddr common.Address) error {
	fmt.Println("\n=== Additional Information ===")

	// Get global validator config
	globalConfig, err := contract.GlobalValidatorConfig(nil)
	if err != nil {
		return fmt.Errorf("failed to get global validator config: %w", err)
	}

	fmt.Printf("Initial Validator Deposit: %s MITO (%s wei)\n",
		FormatWeiToEther(globalConfig.InitialValidatorDeposit),
		globalConfig.InitialValidatorDeposit.String())
	fmt.Printf("Minimum Commission Rate: %s\n", FormatBasisPointsToPercent(globalConfig.MinimumCommissionRate))
	fmt.Printf("Commission Rate Update Delay: %s epochs\n", globalConfig.CommissionRateUpdateDelayEpoch.String())
	fmt.Printf("Collateral Withdrawal Delay: %s seconds\n", globalConfig.CollateralWithdrawalDelaySeconds.String())

	// Get permitted collateral owners count
	permittedOwnersCount, err := contract.PermittedCollateralOwnerSize(nil, validatorAddr)
	if err != nil {
		return fmt.Errorf("failed to get permitted collateral owners count: %w", err)
	}

	fmt.Printf("Permitted Collateral Owners Count: %s\n", permittedOwnersCount.String())

	// List permitted collateral owners if any
	if permittedOwnersCount.Uint64() > 0 {
		fmt.Println("\n=== Permitted Collateral Owners ===")
		for i := uint64(0); i < permittedOwnersCount.Uint64() && i < 10; i++ { // Limit to first 10
			owner, err := contract.PermittedCollateralOwnerAt(nil, validatorAddr, big.NewInt(int64(i)))
			if err != nil {
				fmt.Printf("Error getting owner at index %d: %v\n", i, err)
				continue
			}
			fmt.Printf("%d: %s\n", i+1, owner.Hex())
		}
		if permittedOwnersCount.Uint64() > 10 {
			fmt.Printf("... and %d more\n", permittedOwnersCount.Uint64()-10)
		}
	}

	return nil
}
