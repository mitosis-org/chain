package cmd

import (
	"fmt"
	"log"

	"github.com/mitosis-org/chain/cmd/mivalidator/utils"
	"github.com/spf13/cobra"
)

// NewUpdateMetadataCmd creates a new command to update a validator's metadata
func NewUpdateMetadataCmd() *cobra.Command {
	var (
		validator string
		metadata  string
	)

	cmd := &cobra.Command{
		Use:   "update-metadata",
		Short: "Update a validator's metadata",
		Long:  `Update the metadata for an existing validator.`,
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

			// Show summary and important information
			fmt.Println("===== Update Metadata Transaction =====")
			fmt.Printf("Validator Address        : %s\n", valAddr.Hex())
			fmt.Printf("Current Metadata         : %s\n", string(validatorInfo.Metadata))
			fmt.Printf("New Metadata             : %s\n", metadata)

			// Ask for confirmation
			if !ConfirmAction("Do you want to update the metadata?") {
				log.Fatal("Operation cancelled by user")
			}

			// Execute the transaction
			tx, err := contract.UpdateMetadata(TransactOpts(nil), valAddr, []byte(metadata))
			if err != nil {
				log.Fatalf("Error updating metadata: %v", err)
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
	cmd.Flags().StringVar(&metadata, "metadata", "", "New metadata (JSON string)")

	// Mark required flags
	cmd.MarkFlagRequired("validator")
	cmd.MarkFlagRequired("metadata")

	return cmd
}
