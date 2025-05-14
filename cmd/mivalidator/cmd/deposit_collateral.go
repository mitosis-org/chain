package cmd

import (
	"fmt"
	"log"
	"math/big"

	"github.com/mitosis-org/chain/cmd/mivalidator/utils"
	"github.com/spf13/cobra"
)

// NewDepositCollateralCmd creates a new command to deposit collateral
func NewDepositCollateralCmd() *cobra.Command {
	var (
		validator string
		amount    string
	)

	cmd := &cobra.Command{
		Use:   "deposit-collateral",
		Short: "Deposit collateral for a validator",
		Long:  `Add more collateral to an existing validator.`,
		Run: func(cmd *cobra.Command, args []string) {
			// Validate validator address
			if validator == "" {
				log.Fatal("Error: validator address is required")
			}
			valAddr, err := utils.ValidateAddress(validator)
			if err != nil {
				log.Fatalf("Error: %v", err)
			}

			// Get the contract fee
			fee, err := contract.Fee(nil)
			if err != nil {
				log.Fatalf("Error getting contract fee: %v", err)
			}

			// Parse collateral amount as decimal MITO and convert to wei
			collateralAmount, err := utils.ParseValueAsWei(amount)
			if err != nil {
				log.Fatalf("Error parsing amount: %v", err)
			}

			// Ensure collateral amount is greater than 0
			if collateralAmount.Cmp(big.NewInt(0)) <= 0 {
				log.Fatal("Error: Collateral amount must be greater than 0")
			}

			// Calculate total value to send (collateral amount + fee)
			totalValue := new(big.Int).Add(collateralAmount, fee)

			// Show summary
			fmt.Println("===== Deposit Collateral Transaction =====")
			fmt.Printf("Validator Address        : %s\n", valAddr.Hex())
			fmt.Printf("Collateral Amount        : %s MITO\n", utils.FormatWeiToEther(collateralAmount))
			fmt.Printf("Fee                      : %s MITO\n", utils.FormatWeiToEther(fee))

			fmt.Println("\n🚨 IMPORTANT: Only permitted collateral owners can deposit collateral for a validator.")
			fmt.Println("If your address is not a permitted collateral owner for this validator, the transaction will fail.")
			fmt.Println("To check all permitted collateral owners, use 'validator-info --validator <validator-address>'")

			// Ask for confirmation
			if !ConfirmAction("Do you want to deposit this collateral?") {
				log.Fatal("Operation cancelled by user")
			}

			// Execute the transaction
			tx, err := contract.DepositCollateral(TransactOpts(totalValue), valAddr)
			if err != nil {
				log.Fatalf("Error depositing collateral: %v", err)
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
	cmd.Flags().StringVar(&amount, "amount", "", "Amount of collateral to deposit in MITO (e.g., \"1.5\")")

	// Mark required flags
	mustMarkFlagRequired(cmd, "validator")
	mustMarkFlagRequired(cmd, "amount")

	return cmd
}
