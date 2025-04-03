package cmd

import (
	"fmt"
	"log"

	"github.com/mitosis-org/chain/cmd/mivalidator/utils"
	"github.com/spf13/cobra"
)

// NewUpdateRewardManagerCmd creates a new command to update a validator's reward manager
func NewUpdateRewardManagerCmd() *cobra.Command {
	var (
		validator     string
		rewardManager string
	)

	cmd := &cobra.Command{
		Use:   "update-reward-manager",
		Short: "Update a validator's reward manager",
		Long:  `Update the reward manager address for an existing validator.`,
		Run: func(cmd *cobra.Command, args []string) {
			// Validate validator address
			if validator == "" {
				log.Fatal("Error: validator address is required")
			}
			valAddr, err := utils.ValidateAddress(validator)
			if err != nil {
				log.Fatalf("Error: %v", err)
			}

			// Validate reward manager address
			if rewardManager == "" {
				log.Fatal("Error: reward manager address is required")
			}
			rewardManagerAddr, err := utils.ValidateAddress(rewardManager)
			if err != nil {
				log.Fatalf("Error: %v", err)
			}

			// Get validator info to show current values
			validatorInfo, err := contract.ValidatorInfo(nil, valAddr)
			if err != nil {
				log.Fatalf("Error getting validator info: %v", err)
			}

			// Show summary and important information
			fmt.Println("===== Update Reward Manager Transaction =====")
			fmt.Printf("Validator Address          : %s\n", valAddr.Hex())
			fmt.Printf("Current Reward Manager     : %s\n", validatorInfo.RewardManager.Hex())
			fmt.Printf("New Reward Manager         : %s\n", rewardManagerAddr.Hex())

			// Show important information
			fmt.Println("\n🚨 IMPORTANT: The reward manager will be responsible for:")
			fmt.Println("1. Managing validator rewards")
			fmt.Println("2. Setting commission rates")
			fmt.Println("Make sure this address is secure and under your control.")

			// Ask for confirmation
			if !ConfirmAction("Do you want to update the reward manager address?") {
				log.Fatal("Operation cancelled by user")
			}

			// Execute the transaction
			tx, err := contract.UpdateRewardManager(TransactOpts(nil), valAddr, rewardManagerAddr)
			if err != nil {
				log.Fatalf("Error updating reward manager: %v", err)
			}

			// Wait for transaction confirmation
			err = WaitForTxConfirmation(client, tx.Hash())
			if err != nil {
				log.Fatalf("Transaction failed: %v", err)
			}
		},
	}

	// Add common flags
	AddCommonFlags(cmd, false)

	// Command-specific flags
	cmd.Flags().StringVar(&validator, "validator", "", "Validator address")
	cmd.Flags().StringVar(&rewardManager, "reward-manager", "", "New reward manager address")

	// Mark required flags
	cmd.MarkFlagRequired("validator")
	cmd.MarkFlagRequired("reward-manager")

	return cmd
}
