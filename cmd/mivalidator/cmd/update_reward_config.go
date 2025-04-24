package cmd

import (
	"log"
	"math/big"

	"github.com/mitosis-org/chain/bindings"
	"github.com/mitosis-org/chain/cmd/mivalidator/utils"
	"github.com/spf13/cobra"
)

// NewUpdateRewardConfigCmd creates a new command to update a validator's reward configuration
func NewUpdateRewardConfigCmd() *cobra.Command {
	var (
		validator      string
		commissionRate string
	)

	cmd := &cobra.Command{
		Use:   "update-reward-config",
		Short: "Update a validator's reward configuration",
		Long:  `Update the reward configuration for an existing validator.`,
		Run: func(cmd *cobra.Command, args []string) {
			// Validate validator address
			if validator == "" {
				log.Fatal("Error: validator address is required")
			}
			valAddr, err := utils.ValidateAddress(validator)
			if err != nil {
				log.Fatalf("Error: %v", err)
			}

			// Get validator info to show current values
			validatorInfo, err := contract.ValidatorInfo(nil, valAddr)
			if err != nil {
				log.Fatalf("Error getting validator info: %v", err)
			}

			// Parse commission rate
			commissionRateInt, err := utils.ParsePercentageToBasisPoints(commissionRate)
			if err != nil {
				log.Fatalf("Error parsing commission rate: %v", err)
			}

			// Validate commission rate
			maxRate, err := contract.MAXCOMMISSIONRATE(nil)
			if err != nil {
				log.Fatalf("Error getting max commission rate: %v", err)
			}

			// Get global config for delay information
			config, err := contract.GlobalValidatorConfig(nil)
			if err != nil {
				log.Fatalf("Error getting global validator config: %v", err)
			}

			if commissionRateInt.Cmp(big.NewInt(0)) < 0 || commissionRateInt.Cmp(maxRate) > 0 {
				log.Fatalf("Error: Commission rate must be between 0%% and %s", utils.FormatBasisPointsToPercent(maxRate))
			}

			// Create the request
			request := bindings.IValidatorManagerUpdateRewardConfigRequest{
				CommissionRate: commissionRateInt,
			}

			// Show summary and important information
			log.Println("===== Update Reward Configuration Transaction =====")
			log.Printf("Validator Address           : %s\n", valAddr.Hex())
			log.Printf("Current Commission Rate     : %s\n", utils.FormatBasisPointsToPercent(validatorInfo.CommissionRate))
			log.Printf("New Commission Rate         : %s\n", utils.FormatBasisPointsToPercent(commissionRateInt))
			log.Printf("Update Delay                : %s epochs\n", config.CommissionRateUpdateDelayEpoch.String())

			log.Println("\n🚨 IMPORTANT: Commission rate changes are subject to a delay period.")
			log.Printf("The new commission rate will take effect after %s epochs from now.\n",
				config.CommissionRateUpdateDelayEpoch.String())

			// Ask for confirmation
			if !ConfirmAction("Do you want to update the reward configuration?") {
				log.Fatal("Operation cancelled by user")
			}

			// Execute the transaction
			tx, err := contract.UpdateRewardConfig(TransactOpts(nil), valAddr, request)
			if err != nil {
				log.Fatalf("Error updating reward configuration: %v", err)
			}

			// Handle transaction - either print unsigned or wait for confirmation
			err = HandleTransaction(tx)
			if err != nil {
				log.Fatalf("Transaction failed: %v", err)
			}
		},
	}

	// Add common flags
	AddCommonFlags(cmd, false)

	// Command-specific flags
	cmd.Flags().StringVar(&validator, "validator", "", "Validator address")
	cmd.Flags().StringVar(&commissionRate, "commission-rate", "", "New commission rate in percentage (e.g., \"5%\")")

	// Mark required flags
	cmd.MarkFlagRequired("validator")
	cmd.MarkFlagRequired("commission-rate")

	return cmd
}
