package cmd

import (
	"fmt"
	"log"
	"math/big"

	"github.com/mitosis-org/chain/cmd/mivalidator/utils"
	"github.com/spf13/cobra"
)

// NewWithdrawCollateralCmd creates a new command to withdraw collateral
func NewWithdrawCollateralCmd() *cobra.Command {
	var (
		validator string
		amount    string
	)

	cmd := &cobra.Command{
		Use:   "withdraw-collateral",
		Short: "Withdraw collateral from a validator",
		Long:  `Withdraw collateral from an existing validator.`,
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

			// Parse amount to withdraw
			amountInWei, err := utils.ParseValue(amount)
			if err != nil {
				log.Fatalf("Error parsing amount: %v", err)
			}

			// Ensure amount is greater than 0
			if amountInWei.Cmp(big.NewInt(0)) <= 0 {
				log.Fatal("Error: Amount must be greater than 0")
			}

			// Get the global validator config for collateral withdrawal delay information
			config, err := contract.GlobalValidatorConfig(nil)
			if err != nil {
				log.Fatalf("Error getting global validator config: %v", err)
			}

			// Get validator info to show the withdrawal recipient
			validatorInfo, err := contract.ValidatorInfo(nil, valAddr)
			if err != nil {
				log.Fatalf("Error getting validator info: %v", err)
			}

			// Show summary and important information
			fmt.Println("===== Withdraw Collateral Transaction =====")
			fmt.Printf("Validator Address          : %s\n", valAddr.Hex())
			fmt.Printf("Amount to withdraw         : %s MITO\n", utils.FormatWeiToEther(amountInWei))
			fmt.Printf("Fee                        : %s MITO\n", utils.FormatWeiToEther(fee))
			fmt.Printf("Withdrawal Delay           : %s seconds\n", config.CollateralWithdrawalDelaySeconds.String())
			fmt.Printf("Withdrawal Recipient       : %s\n", validatorInfo.WithdrawalRecipient.Hex())

			fmt.Println("\n🚨 IMPORTANT: The collateral withdrawal is subject to a delay period.")
			fmt.Printf("Your funds will be available after %s seconds from transaction confirmation.\n", config.CollateralWithdrawalDelaySeconds.String())
			fmt.Printf("The withdrawing amount will be sent to your validator's withdrawal recipient address: %s\n", validatorInfo.WithdrawalRecipient.Hex())

			// Ask for confirmation
			if !ConfirmAction("Do you want to withdraw this collateral?") {
				log.Fatal("Operation cancelled by user")
			}

			// Execute the transaction
			tx, err := contract.WithdrawCollateral(TransactOpts(fee), valAddr, amountInWei)
			if err != nil {
				log.Fatalf("Error withdrawing collateral: %v", err)
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
	cmd.Flags().StringVar(&amount, "amount", "", "Amount to withdraw in MITO")

	// Mark required flags
	cmd.MarkFlagRequired("validator")
	cmd.MarkFlagRequired("amount")

	return cmd
}
