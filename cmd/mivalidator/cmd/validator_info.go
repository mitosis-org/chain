package cmd

import (
	"encoding/hex"
	"fmt"
	"log"
	"math/big"

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
	fmt.Printf("Withdrawal Recipient     : %s\n", validatorInfo.WithdrawalRecipient.Hex())
	fmt.Printf("Commission Rate          : %s\n", utils.FormatBasisPointsToPercent(validatorInfo.CommissionRate))
	fmt.Printf("Metadata                 : %s\n", string(validatorInfo.Metadata))
}
