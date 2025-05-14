package cmd

import (
	"fmt"
	"log"

	"github.com/mitosis-org/chain/cmd/mivalidator/utils"
	"github.com/spf13/cobra"
)

// NewTransferCollateralOwnershipCmd creates a new command to transfer collateral ownership
func NewTransferCollateralOwnershipCmd() *cobra.Command {
	var (
		validator string
		newOwner  string
	)

	cmd := &cobra.Command{
		Use:   "transfer-collateral-ownership",
		Short: "Transfer validator collateral ownership",
		Long:  `Transfer ownership of a validator's collateral to a new address.`,
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

			// Validate new owner address
			if newOwner == "" {
				log.Fatal("Error: new owner address is required")
			}
			newOwnerAddr, err := utils.ValidateAddress(newOwner)
			if err != nil {
				log.Fatalf("Error: %v", err)
			}

			// Show summary and important information
			fmt.Println("===== Transfer Collateral Ownership Transaction =====")
			fmt.Printf("Validator Address        : %s\n", valAddr.Hex())
			fmt.Printf("New Collateral Owner     : %s\n", newOwnerAddr.Hex())
			fmt.Printf("Fee                      : %s MITO\n", utils.FormatWeiToEther(fee))

			fmt.Println("\n🚨 IMPORTANT: Only permitted collateral owners can transfer collateral ownership.")
			fmt.Println("If your address is not a permitted collateral owner for this validator, the transaction will fail.")
			fmt.Println("To check all permitted collateral owners, use 'validator-info --validator <validator-address>'")

			fmt.Println("\n🚨 IMPORTANT: This action will transfer ownership of the validator's collateral.")
			fmt.Println("The new owner will have full control over the collateral, including the ability to withdraw it.")
			fmt.Println("Make sure you trust the new owner or it is an address you control.")

			// Ask for confirmation
			if !ConfirmAction("Do you want to transfer collateral ownership?") {
				log.Fatal("Operation cancelled by user")
			}

			// Execute the transaction
			tx, err := contract.TransferCollateralOwnership(TransactOpts(fee), valAddr, newOwnerAddr)
			if err != nil {
				log.Fatalf("Error transferring collateral ownership: %v", err)
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
	cmd.Flags().StringVar(&newOwner, "new-owner", "", "New collateral owner address")

	// Mark required flags
	mustMarkFlagRequired(cmd, "validator")
	mustMarkFlagRequired(cmd, "new-owner")

	return cmd
}
