package cmd

import (
	"fmt"
	"log"

	"github.com/mitosis-org/chain/cmd/mivalidator/utils"
	"github.com/spf13/cobra"
)

// NewUpdateOperatorCmd creates a new command to update a validator's operator
func NewUpdateOperatorCmd() *cobra.Command {
	var (
		validator string
		operator  string
	)

	cmd := &cobra.Command{
		Use:   "update-operator",
		Short: "Update a validator's operator",
		Long:  `Update the operator address for an existing validator.`,
		Run: func(cmd *cobra.Command, args []string) {
			// Validate validator address
			if validator == "" {
				log.Fatal("Error: validator address is required")
			}
			valAddr, err := utils.ValidateAddress(validator)
			if err != nil {
				log.Fatalf("Error: %v", err)
			}

			// Validate operator address
			if operator == "" {
				log.Fatal("Error: operator address is required")
			}
			operatorAddr, err := utils.ValidateAddress(operator)
			if err != nil {
				log.Fatalf("Error: %v", err)
			}

			// Get validator info to show current values and warn about changes
			validatorInfo, err := contract.ValidatorInfo(nil, valAddr)
			if err != nil {
				log.Fatalf("Error getting validator info: %v", err)
			}

			// Show summary and important information
			fmt.Println("===== Update Operator Transaction =====")
			fmt.Printf("Validator Address            : %s\n", valAddr.Hex())
			fmt.Printf("Current Operator             : %s\n", validatorInfo.Operator.Hex())
			fmt.Printf("New Operator                 : %s\n", operatorAddr.Hex())
			fmt.Printf("Current Withdrawal Recipient : %s\n", validatorInfo.WithdrawalRecipient.Hex())
			fmt.Printf("Current Reward Manager       : %s\n", validatorInfo.RewardManager.Hex())

			// Show important warning about other updates that might be needed
			fmt.Println("\n🚨 IMPORTANT WARNING 🚨")
			fmt.Println("When changing the operator address, you may also want to update:")
			fmt.Println("1. Withdrawal Recipient - To ensure funds are sent to the correct address")
			fmt.Println("2. Reward Manager - To ensure rewards are managed by the correct entity")
			fmt.Println("\nAfter this operation completes, consider using:")
			fmt.Printf("  - mivalidator update-withdrawal-recipient --validator %s --recipient <new-address>\n", valAddr.Hex())
			fmt.Printf("  - mivalidator update-reward-manager --validator %s --reward-manager <new-address>\n", valAddr.Hex())

			// Ask for confirmation
			if !ConfirmAction("Do you want to update the operator address?") {
				log.Fatal("Operation cancelled by user")
			}

			// Execute the transaction
			tx, err := contract.UpdateOperator(TransactOpts(nil), valAddr, operatorAddr)
			if err != nil {
				log.Fatalf("Error updating operator: %v", err)
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
	cmd.Flags().StringVar(&operator, "operator", "", "New operator address")

	// Mark required flags
	mustMarkFlagRequired(cmd, "validator")
	mustMarkFlagRequired(cmd, "operator")

	return cmd
}
