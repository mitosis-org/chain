package cmd

import (
	"fmt"
	"log"

	"github.com/mitosis-org/chain/cmd/mivalidator/utils"
	"github.com/spf13/cobra"
)

// NewUpdateWithdrawalRecipientCmd creates a new command to update a validator's withdrawal recipient
func NewUpdateWithdrawalRecipientCmd() *cobra.Command {
	var (
		validator string
		recipient string
	)

	cmd := &cobra.Command{
		Use:   "update-withdrawal-recipient",
		Short: "Update a validator's withdrawal recipient",
		Long: `Update the withdrawal recipient for an existing validator. The withdrawal recipient
is the address that receives the withdrawal of any ETH and validator rewards when a validator
exits the system.`,
		Run: func(cmd *cobra.Command, args []string) {
			// Validate validator address
			if validator == "" {
				log.Fatal("Error: validator address is required")
			}
			valAddr, err := utils.ValidateAddress(validator)
			if err != nil {
				log.Fatalf("Error: %v", err)
			}

			// Validate recipient address
			if recipient == "" {
				log.Fatal("Error: withdrawal recipient address is required")
			}
			recipientAddr, err := utils.ValidateAddress(recipient)
			if err != nil {
				log.Fatalf("Error: %v", err)
			}

			// Get validator info to show current values
			validatorInfo, err := contract.ValidatorInfo(nil, valAddr)
			if err != nil {
				log.Fatalf("Error getting validator info: %v", err)
			}

			// Show summary and important information
			fmt.Println("===== Update Withdrawal Recipient Transaction =====")
			fmt.Printf("Validator Address            : %s\n", valAddr.Hex())
			fmt.Printf("Current Withdrawal Recipient : %s\n", validatorInfo.WithdrawalRecipient.Hex())
			fmt.Printf("New Withdrawal Recipient     : %s\n", recipientAddr.Hex())

			fmt.Println("\n🚨 IMPORTANT: The withdrawal recipient is the address that will receive funds when a validator exits.")
			fmt.Println("Make sure you have full control over this address or it belongs to your intended recipient.")

			// Ask for confirmation
			if !ConfirmAction("Do you want to update the withdrawal recipient?") {
				log.Fatal("Operation cancelled by user")
			}

			// Execute the transaction
			tx, err := contract.UpdateWithdrawalRecipient(TransactOpts(nil), valAddr, recipientAddr)
			if err != nil {
				log.Fatalf("Error updating withdrawal recipient: %v", err)
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
	cmd.Flags().StringVar(&recipient, "recipient", "", "New withdrawal recipient address")

	// Mark required flags
	cmd.MarkFlagRequired("validator")
	cmd.MarkFlagRequired("recipient")

	return cmd
}
