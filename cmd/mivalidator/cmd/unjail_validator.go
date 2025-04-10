package cmd

import (
	"fmt"
	"log"

	"github.com/mitosis-org/chain/cmd/mivalidator/utils"
	"github.com/spf13/cobra"
)

// NewUnjailValidatorCmd creates a new command to unjail a validator
func NewUnjailValidatorCmd() *cobra.Command {
	var (
		validator string
	)

	cmd := &cobra.Command{
		Use:   "unjail-validator",
		Short: "Unjail a validator",
		Long:  `Unjail a validator that has been jailed.`,
		Run: func(cmd *cobra.Command, args []string) {
			// Get the contract fee
			fee, err := contract.Fee(nil)
			if err != nil {
				log.Fatalf("Error getting contract fee: %v", err)
			}

			// Validate validator address
			if validator == "" {
				log.Fatal("Error: validator address is required")
			}
			valAddr, err := utils.ValidateAddress(validator)
			if err != nil {
				log.Fatalf("Error: %v", err)
			}

			// Show summary
			fmt.Println("===== Unjail Validator Transaction =====")
			fmt.Printf("Validator Address        : %s\n", valAddr.Hex())
			fmt.Printf("Fee                      : %s MITO\n", utils.FormatWeiToEther(fee))

			// Ask for confirmation
			if !ConfirmAction("Do you want to unjail this validator?") {
				log.Fatal("Operation cancelled by user")
			}

			// Execute the transaction
			tx, err := contract.UnjailValidator(TransactOpts(fee), valAddr)
			if err != nil {
				log.Fatalf("Error unjailing validator: %v", err)
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
	cmd.Flags().StringVar(&validator, "validator", "", "Address of the validator to unjail")

	// Mark required flags
	cmd.MarkFlagRequired("validator")

	return cmd
}
