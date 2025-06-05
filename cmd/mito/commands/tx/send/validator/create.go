package validator

import (
	"fmt"

	"github.com/mitosis-org/chain/cmd/mito/internal/config"
	"github.com/mitosis-org/chain/cmd/mito/internal/container"
	"github.com/mitosis-org/chain/cmd/mito/internal/flags"
	"github.com/mitosis-org/chain/cmd/mito/internal/output"
	"github.com/mitosis-org/chain/cmd/mito/internal/tx"
	"github.com/mitosis-org/chain/cmd/mito/internal/utils"
	"github.com/mitosis-org/chain/cmd/mito/internal/validation"
	"github.com/spf13/cobra"
)

// NewCreateCmd creates the send validator create command
func NewCreateCmd() *cobra.Command {
	var commonFlags flags.CommonFlags
	var validatorFlags struct {
		pubkey            string
		operator          string
		rewardManager     string
		commissionRate    string
		metadata          string
		initialCollateral string
	}

	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create and send a new validator transaction",
		Long:  "Create, sign and send a new validator transaction to the network",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			// Validate mutually exclusive flags
			return validation.ValidateSendTxFlagGroups(cmd)
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			cmd.SilenceUsage = true

			// Resolve configuration
			resolver, err := config.NewResolver()
			if err != nil {
				return err
			}
			resolvedConfig := resolver.ResolveFlags(&commonFlags)

			// Validate fields using new validation module
			validateFields := &validation.ValidatorCreateFields{
				PubKey:            validatorFlags.pubkey,
				Operator:          validatorFlags.operator,
				RewardManager:     validatorFlags.rewardManager,
				CommissionRate:    validatorFlags.commissionRate,
				Metadata:          validatorFlags.metadata,
				InitialCollateral: validatorFlags.initialCollateral,
			}
			if err := validation.ValidateValidatorCreateFields(resolvedConfig, validateFields); err != nil {
				return err
			}

			// Validate network and signing requirements for send
			if err := validation.ValidateNetworkFields(resolvedConfig, true); err != nil {
				return err
			}
			if err := validation.ValidateSigningFields(resolvedConfig, true); err != nil {
				return err
			}

			// Create container
			container, err := container.NewContainer(resolvedConfig)
			if err != nil {
				return fmt.Errorf("failed to initialize container: %w", err)
			}
			defer container.Close()

			// Create validator request
			req := &tx.CreateValidatorRequest{
				PubKey:            validatorFlags.pubkey,
				Operator:          validatorFlags.operator,
				RewardManager:     validatorFlags.rewardManager,
				CommissionRate:    validatorFlags.commissionRate,
				Metadata:          validatorFlags.metadata,
				InitialCollateral: validatorFlags.initialCollateral,
			}

			// Create unsigned transaction
			transaction, err := container.ValidatorService.CreateValidator(req)
			if err != nil {
				return err
			}

			// Display transaction information using formatter
			formatter := output.NewTransactionFormatter("")
			info := &output.ValidatorCreateInfo{
				PubKey:            validatorFlags.pubkey,
				Operator:          validatorFlags.operator,
				RewardManager:     validatorFlags.rewardManager,
				CommissionRate:    validatorFlags.commissionRate,
				Metadata:          validatorFlags.metadata,
				InitialCollateral: validatorFlags.initialCollateral,
			}

			if err := formatter.FormatValidatorCreateTransaction(transaction, info); err != nil {
				return fmt.Errorf("failed to format transaction: %w", err)
			}
			fmt.Println()

			// Confirm action
			if !utils.ConfirmAction(resolvedConfig.Yes, "Create and send validator transaction?") {
				return fmt.Errorf("operation cancelled")
			}

			// Sign transaction
			signedTx, err := container.TxBuilder.SignTransaction(transaction)
			if err != nil {
				return fmt.Errorf("failed to sign transaction: %w", err)
			}

			// Send transaction and wait for confirmation
			_, err = container.TxSender.SendAndWait(signedTx)
			if err != nil {
				return err
			}

			fmt.Println("✅ Validator created successfully!")
			return nil
		},
	}

	// Add flags
	flags.AddCommonFlags(cmd, &commonFlags)
	flags.AddSendFlags(cmd, &commonFlags)
	cmd.Flags().StringVar(&validatorFlags.pubkey, "pubkey", "", "Validator's public key (hex with 0x prefix)")
	cmd.Flags().StringVar(&validatorFlags.operator, "operator", "", "Operator address")
	cmd.Flags().StringVar(&validatorFlags.rewardManager, "reward-manager", "", "Reward manager address")
	cmd.Flags().StringVar(&validatorFlags.commissionRate, "commission-rate", "", "Commission rate in percentage (e.g., \"5%\")")
	cmd.Flags().StringVar(&validatorFlags.metadata, "metadata", "", "Validator metadata (JSON string)")
	cmd.Flags().StringVar(&validatorFlags.initialCollateral, "initial-collateral", "", "Initial collateral amount in MITO (e.g., \"1.5\")")

	// Mark required flags
	cmd.MarkFlagRequired("pubkey")
	cmd.MarkFlagRequired("operator")
	cmd.MarkFlagRequired("reward-manager")
	cmd.MarkFlagRequired("commission-rate")
	cmd.MarkFlagRequired("metadata")
	cmd.MarkFlagRequired("initial-collateral")

	return cmd
}
