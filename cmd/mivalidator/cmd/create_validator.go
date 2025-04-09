package cmd

import (
	"fmt"
	"log"
	"math/big"

	"github.com/mitosis-org/chain/bindings"
	"github.com/mitosis-org/chain/cmd/mivalidator/utils"
	"github.com/spf13/cobra"
)

// NewCreateValidatorCmd creates a new command to create a validator
func NewCreateValidatorCmd() *cobra.Command {
	var (
		pubKey              string
		operator            string
		rewardManager       string
		withdrawalRecipient string
		commissionRate      string
		metadata            string
		initialCollateral   string
	)

	cmd := &cobra.Command{
		Use:   "create-validator",
		Short: "Create a new validator",
		Long:  `Register a new validator in the ValidatorManager contract.`,
		Run: func(cmd *cobra.Command, args []string) {
			// Get the contract fee
			fee, err := contract.Fee(nil)
			if err != nil {
				log.Fatalf("Error getting contract fee: %v", err)
			}

			// Get the config to check the initial deposit requirement
			config, err := contract.GlobalValidatorConfig(nil)
			if err != nil {
				log.Fatalf("Error getting global validator config: %v", err)
			}

			// Parse collateral amount as decimal MITO and convert to wei
			collateralAmount, err := utils.ParseValueAsWei(initialCollateral)
			if err != nil {
				log.Fatalf("Error parsing initial collateral: %v", err)
			}

			// Ensure collateral is at least the initial deposit requirement
			if collateralAmount.Cmp(config.InitialValidatorDeposit) < 0 {
				log.Fatalf("Error: Initial collateral must be at least %s MITO",
					utils.FormatWeiToEther(config.InitialValidatorDeposit))
			}

			// Calculate total transaction value (collateral + fee)
			totalValue := new(big.Int).Add(collateralAmount, fee)

			// Validate other parameters
			operatorAddr, err := utils.ValidateAddress(operator)
			if err != nil {
				log.Fatalf("Error in operator address: %v", err)
			}

			rewardManagerAddr, err := utils.ValidateAddress(rewardManager)
			if err != nil {
				log.Fatalf("Error in reward manager address: %v", err)
			}

			withdrawalRecipientAddr, err := utils.ValidateAddress(withdrawalRecipient)
			if err != nil {
				log.Fatalf("Error in withdrawal recipient address: %v", err)
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

			if commissionRateInt.Cmp(big.NewInt(0)) < 0 || commissionRateInt.Cmp(maxRate) > 0 {
				log.Fatalf("Error: Commission rate must be between 0%% and %s", utils.FormatBasisPointsToPercent(maxRate))
			}

			// Decode public key from hex
			pubKeyBytes, err := utils.DecodeHexWithPrefix(pubKey)
			if err != nil {
				log.Fatalf("Error decoding public key: %v", err)
			}

			// Create the request
			request := bindings.IValidatorManagerCreateValidatorRequest{
				Operator:            operatorAddr,
				RewardManager:       rewardManagerAddr,
				WithdrawalRecipient: withdrawalRecipientAddr,
				CommissionRate:      commissionRateInt,
				Metadata:            []byte(metadata),
			}

			// Show summary
			fmt.Println("===== Create Validator Transaction =====")
			fmt.Printf("Public Key                 : %s\n", pubKey)
			fmt.Printf("Operator                   : %s\n", operatorAddr.Hex())
			fmt.Printf("Reward Manager             : %s\n", rewardManagerAddr.Hex())
			fmt.Printf("Withdrawal Recipient       : %s\n", withdrawalRecipientAddr.Hex())
			fmt.Printf("Commission Rate            : %s\n", utils.FormatBasisPointsToPercent(commissionRateInt))
			fmt.Printf("Metadata                   : %s\n", metadata)
			fmt.Printf("Initial Collateral         : %s MITO\n", utils.FormatWeiToEther(collateralAmount))
			fmt.Printf("Fee                        : %s MITO\n", utils.FormatWeiToEther(fee))

			// Ask for confirmation
			if !ConfirmAction("Do you want to create this validator?") {
				log.Fatal("Operation cancelled by user")
			}

			// Execute the transaction
			tx, err := contract.CreateValidator(TransactOpts(totalValue), pubKeyBytes, request)
			if err != nil {
				log.Fatalf("Error creating validator: %v", err)
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
	cmd.Flags().StringVar(&pubKey, "pubkey", "", "Validator's public key (hex with 0x prefix)")
	cmd.Flags().StringVar(&operator, "operator", "", "Operator address")
	cmd.Flags().StringVar(&rewardManager, "reward-manager", "", "Reward manager address")
	cmd.Flags().StringVar(&withdrawalRecipient, "withdrawal-recipient", "", "Withdrawal recipient address")
	cmd.Flags().StringVar(&commissionRate, "commission-rate", "", "Commission rate in percentage (e.g., \"5%\")")
	cmd.Flags().StringVar(&metadata, "metadata", "", "Validator metadata (JSON string)")
	cmd.Flags().StringVar(&initialCollateral, "initial-collateral", "", "Initial collateral amount in MITO (e.g., \"1.5\")")

	// Mark required flags
	cmd.MarkFlagRequired("pubkey")
	cmd.MarkFlagRequired("operator")
	cmd.MarkFlagRequired("reward-manager")
	cmd.MarkFlagRequired("withdrawal-recipient")
	cmd.MarkFlagRequired("commission-rate")
	cmd.MarkFlagRequired("metadata")
	cmd.MarkFlagRequired("initial-collateral")

	return cmd
}
