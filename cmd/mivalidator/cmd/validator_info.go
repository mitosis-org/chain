package cmd

import (
	"encoding/hex"
	"fmt"
	"log"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/mitosis-org/chain/bindings"
	"github.com/mitosis-org/chain/cmd/mivalidator/utils"
	"github.com/spf13/cobra"
)

// NewValidatorInfoCmd creates a new command to get validator information
func NewValidatorInfoCmd() *cobra.Command {
	var (
		validator string
		epoch     uint64
	)

	cmd := &cobra.Command{
		Use:   "validator-info",
		Short: "Get information about a validator",
		Long:  `Retrieve detailed information about a validator from the ValidatorManager contract.`,
		Run: func(cmd *cobra.Command, args []string) {
			// Validate validator address
			if validator == "" {
				log.Fatal("Error: validator address is required")
			}
			valAddr, err := utils.ValidateAddress(validator)
			if err != nil {
				log.Fatalf("Error: %v", err)
			}

			// Check if this is a validator
			isValidator, err := contract.IsValidator(nil, valAddr)
			if err != nil {
				log.Fatalf("Error checking if address is a validator: %v", err)
			}

			if !isValidator {
				fmt.Println("NOTE: This address is not registered as a validator.")
			}

			var validatorInfo bindings.IValidatorManagerValidatorInfoResponse
			if epoch > 0 {
				// Get validator info at a specific epoch
				validatorInfo, err = contract.ValidatorInfoAt(nil, new(big.Int).SetUint64(epoch), valAddr)
			} else {
				// Get current validator info
				validatorInfo, err = contract.ValidatorInfo(nil, valAddr)
			}

			if err != nil {
				log.Fatalf("Error getting validator info: %v", err)
			}

			// Display validator info
			displayValidatorInfo(validatorInfo)

			// Display permitted collateral owners
			displayPermittedCollateralOwners(valAddr)
		},
	}

	// Add common flags, specifying this is a read-only command
	AddCommonFlags(cmd, true)

	// Command-specific flags
	cmd.Flags().StringVar(&validator, "validator", "", "Validator address")
	cmd.Flags().Uint64Var(&epoch, "epoch", 0, "Epoch number (0 for current epoch)")
	mustMarkFlagRequired(cmd, "validator")

	return cmd
}

// displayValidatorInfo formats and displays the validator information
func displayValidatorInfo(validatorInfo bindings.IValidatorManagerValidatorInfoResponse) {
	fmt.Println("===== Validator Information =====")
	fmt.Printf("Validator Address        : %s\n", validatorInfo.ValAddr.Hex())
	fmt.Printf("Public Key               : 0x%s\n", hex.EncodeToString(validatorInfo.PubKey))
	fmt.Printf("Operator                 : %s\n", validatorInfo.Operator.Hex())
	fmt.Printf("Reward Manager           : %s\n", validatorInfo.RewardManager.Hex())
	fmt.Printf("Commission Rate          : %s\n", utils.FormatBasisPointsToPercent(validatorInfo.CommissionRate))
	fmt.Printf("Metadata                 : %s\n", string(validatorInfo.Metadata))
}

// displayPermittedCollateralOwners retrieves and displays all permitted collateral owners for a validator
func displayPermittedCollateralOwners(valAddr common.Address) {
	// Get the total count of permitted collateral owners
	size, err := contract.PermittedCollateralOwnerSize(nil, valAddr)
	if err != nil {
		log.Printf("Error getting permitted collateral owner size: %v", err)
		return
	}

	fmt.Println("\n===== Permitted Collateral Owners =====")
	fmt.Println("Permitted collateral owners are addresses that are allowed to:")
	fmt.Println("1. Deposit collateral for this validator")
	fmt.Println("2. Transfer collateral ownership to another address")
	fmt.Println("Only these addresses can perform these operations for this validator.")

	// If there are no permitted collateral owners, display a message
	if size.Cmp(big.NewInt(0)) == 0 {
		fmt.Println("\nNo permitted collateral owners found")
		return
	}

	fmt.Println("\nList of permitted collateral owners:")

	// Iterate through all permitted collateral owners
	for i := int64(0); i < size.Int64(); i++ {
		index := big.NewInt(i)
		collateralOwner, err := contract.PermittedCollateralOwnerAt(nil, valAddr, index)
		if err != nil {
			log.Printf("Error getting permitted collateral owner at index %d: %v", i, err)
			continue
		}

		fmt.Printf("%d. %s\n", i+1, collateralOwner.Hex())
	}
}
