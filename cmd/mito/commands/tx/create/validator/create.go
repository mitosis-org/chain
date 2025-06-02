package validator

import (
	"fmt"

	"github.com/mitosis-org/chain/cmd/mito/internal/config"
	"github.com/mitosis-org/chain/cmd/mito/internal/container"
	"github.com/mitosis-org/chain/cmd/mito/internal/flags"
	"github.com/mitosis-org/chain/cmd/mito/internal/output"
	"github.com/mitosis-org/chain/cmd/mito/internal/tx"
	"github.com/mitosis-org/chain/cmd/mito/internal/validation"
	"github.com/spf13/cobra"
)

// NewCreateCmd creates the create validator create command
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
		Short: "Create a new validator transaction",
		Long:  "Create a new validator transaction (without sending)",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return validation.ValidateCreateTxFlagGroups(cmd)
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			cmd.SilenceUsage = true

			// Resolve configuration
			resolver, err := config.NewResolver()
			if err != nil {
				return err
			}
			resolvedConfig := resolver.ResolveFlags(&commonFlags)

			// Validate validator create fields
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

			// Create container (for validation and fee calculation)
			container, err := container.NewContainer(resolvedConfig)
			if err != nil {
				return fmt.Errorf("failed to initialize container: %w", err)
			}
			defer container.Close()

			// Create request
			req := &tx.CreateValidatorRequest{
				PubKey:            validatorFlags.pubkey,
				Operator:          validatorFlags.operator,
				RewardManager:     validatorFlags.rewardManager,
				CommissionRate:    validatorFlags.commissionRate,
				Metadata:          validatorFlags.metadata,
				InitialCollateral: validatorFlags.initialCollateral,
			}

			// Create transaction directly (validation included)
			tx, err := container.ValidatorService.CreateValidatorWithOptions(req, commonFlags.Unsigned)
			if err != nil {
				return fmt.Errorf("failed to create transaction: %w", err)
			}

			// Format and output transaction
			formatter := output.NewTransactionFormatter(commonFlags.OutputFile)
			info := &output.ValidatorCreateInfo{
				PubKey:            validatorFlags.pubkey,
				Operator:          validatorFlags.operator,
				RewardManager:     validatorFlags.rewardManager,
				CommissionRate:    validatorFlags.commissionRate,
				Metadata:          validatorFlags.metadata,
				InitialCollateral: validatorFlags.initialCollateral,
			}

			return formatter.FormatValidatorCreateTransaction(tx, info)
		},
	}

	// Add flags
	flags.AddCommonFlags(cmd, &commonFlags)
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
