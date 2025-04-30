package cmd

import (
	"fmt"
	"log"

	"github.com/mitosis-org/chain/cmd/mivalidator/utils"
	"github.com/spf13/cobra"
)

// NewSetPermittedCollateralOwnerCmd creates a new command to set a permitted collateral owner
func NewSetPermittedCollateralOwnerCmd() *cobra.Command {
	var (
		validator       string
		collateralOwner string
		isPermitted     bool
	)

	cmd := &cobra.Command{
		Use:   "set-permitted-collateral-owner",
		Short: "Set a permitted collateral owner for a validator",
		Long:  `Allow or disallow an address to deposit collateral for a validator.`,
		Run: func(cmd *cobra.Command, args []string) {
			// Validate validator address
			if validator == "" {
				log.Fatal("Error: validator address is required")
			}
			valAddr, err := utils.ValidateAddress(validator)
			if err != nil {
				log.Fatalf("Error: %v", err)
			}

			// Validate collateral owner address
			if collateralOwner == "" {
				log.Fatal("Error: collateral owner address is required")
			}
			collateralOwnerAddr, err := utils.ValidateAddress(collateralOwner)
			if err != nil {
				log.Fatalf("Error: %v", err)
			}

			// Check current permission status
			currentPermission, err := contract.IsPermittedCollateralOwner(nil, valAddr, collateralOwnerAddr)
			if err != nil {
				log.Fatalf("Error checking current permission status: %v", err)
			}

			// Show summary and important information
			fmt.Println("===== Set Permitted Collateral Owner Transaction =====")
			fmt.Printf("Validator Address          : %s\n", valAddr.Hex())
			fmt.Printf("Collateral Owner Address   : %s\n", collateralOwnerAddr.Hex())
			fmt.Printf("Current Permission Status  : %t\n", currentPermission)
			fmt.Printf("New Permission Status      : %t\n", isPermitted)

			if isPermitted {
				fmt.Println("\n🚨 IMPORTANT: This will allow the specified address to deposit collateral for your validator.")
				fmt.Println("Make sure you trust this address or it is under your control.")
			} else {
				fmt.Println("\n🚨 IMPORTANT: This will revoke permission for the specified address to deposit collateral for your validator.")
			}

			// Ask for confirmation
			if !ConfirmAction("Do you want to update the permission status?") {
				log.Fatal("Operation cancelled by user")
			}

			// Execute the transaction
			tx, err := contract.SetPermittedCollateralOwner(TransactOpts(nil), valAddr, collateralOwnerAddr, isPermitted)
			if err != nil {
				log.Fatalf("Error setting permitted collateral owner: %v", err)
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
	cmd.Flags().StringVar(&collateralOwner, "collateral-owner", "", "Collateral owner address to permit or deny")
	cmd.Flags().BoolVar(&isPermitted, "permitted", false, "Whether to permit (true) or deny (false) the address")

	// Mark required flags
	mustMarkFlagRequired(cmd, "validator")
	mustMarkFlagRequired(cmd, "collateral-owner")

	return cmd
}
