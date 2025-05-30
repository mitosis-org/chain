package create

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/mitosis-org/chain/cmd/mito/internal/config"
	"github.com/mitosis-org/chain/cmd/mito/internal/container"
	"github.com/mitosis-org/chain/cmd/mito/internal/flags"
	"github.com/mitosis-org/chain/cmd/mito/internal/tx"
	"github.com/spf13/cobra"
)

// NewValidatorCmd returns validator commands for tx create
func NewValidatorCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "validator",
		Short: "Create validator transactions",
		Long:  "Create validator transactions (signed or unsigned)",
	}

	cmd.AddCommand(
		newCreateValidatorCreateCmd(),
		newCreateValidatorUpdateMetadataCmd(),
		newCreateValidatorUpdateOperatorCmd(),
		newCreateValidatorUpdateRewardConfigCmd(),
		newCreateValidatorUpdateRewardManagerCmd(),
		newCreateValidatorUnjailCmd(),
	)

	return cmd
}

func newCreateValidatorCreateCmd() *cobra.Command {
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
			return flags.ValidateSigningFlags(cmd)
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			cmd.SilenceUsage = true

			// Resolve configuration
			resolver, err := config.NewResolver()
			if err != nil {
				return err
			}
			resolvedConfig := resolver.ResolveFlags(&commonFlags)

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
			tx, err := container.ValidatorService.CreateValidator(req)
			if err != nil {
				return fmt.Errorf("failed to create transaction: %w", err)
			}

			// Sign if not unsigned
			if !commonFlags.Unsigned {
				tx, err = container.TxBuilder.SignTransaction(tx)
				if err != nil {
					return fmt.Errorf("failed to sign transaction: %w", err)
				}
			}

			// Convert to JSON and output
			txJSON, err := json.Marshal(tx)
			if err != nil {
				return fmt.Errorf("failed to convert transaction to JSON: %w", err)
			}

			if commonFlags.OutputFile != "" {
				return os.WriteFile(commonFlags.OutputFile, txJSON, 0644)
			}

			fmt.Println(string(txJSON))
			return nil
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

func newCreateValidatorUpdateMetadataCmd() *cobra.Command {
	var commonFlags flags.CommonFlags
	var validatorAddr, metadata string

	cmd := &cobra.Command{
		Use:   "update-metadata",
		Short: "Create validator metadata update transaction",
		Long:  "Create a transaction to update validator metadata (without sending)",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return flags.ValidateSigningFlags(cmd)
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			cmd.SilenceUsage = true

			// Resolve configuration
			resolver, err := config.NewResolver()
			if err != nil {
				return err
			}
			resolvedConfig := resolver.ResolveFlags(&commonFlags)

			// Create container
			container, err := container.NewContainer(resolvedConfig)
			if err != nil {
				return fmt.Errorf("failed to initialize container: %w", err)
			}
			defer container.Close()

			// Create transaction directly (validation included)
			tx, err := container.ValidatorService.UpdateMetadata(validatorAddr, metadata)
			if err != nil {
				return fmt.Errorf("failed to create transaction: %w", err)
			}

			// Sign if not unsigned
			if !commonFlags.Unsigned {
				tx, err = container.TxBuilder.SignTransaction(tx)
				if err != nil {
					return fmt.Errorf("failed to sign transaction: %w", err)
				}
			}

			// Convert to JSON and output
			txJSON, err := json.Marshal(tx)
			if err != nil {
				return fmt.Errorf("failed to convert transaction to JSON: %w", err)
			}

			if commonFlags.OutputFile != "" {
				return os.WriteFile(commonFlags.OutputFile, txJSON, 0644)
			}

			fmt.Println(string(txJSON))
			return nil
		},
	}

	// Add flags
	flags.AddCommonFlags(cmd, &commonFlags)
	cmd.Flags().StringVar(&validatorAddr, "validator", "", "Validator address")
	cmd.Flags().StringVar(&metadata, "metadata", "", "New validator metadata")
	cmd.MarkFlagRequired("validator")
	cmd.MarkFlagRequired("metadata")

	return cmd
}

func newCreateValidatorUpdateOperatorCmd() *cobra.Command {
	var commonFlags flags.CommonFlags
	var validatorAddr, operator string

	cmd := &cobra.Command{
		Use:   "update-operator",
		Short: "Create validator operator update transaction",
		Long:  "Create a transaction to update validator operator (without sending)",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return flags.ValidateSigningFlags(cmd)
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			cmd.SilenceUsage = true

			// Resolve configuration
			resolver, err := config.NewResolver()
			if err != nil {
				return err
			}
			resolvedConfig := resolver.ResolveFlags(&commonFlags)

			// Validate required fields
			// Create container
			container, err := container.NewContainer(resolvedConfig)
			if err != nil {
				return fmt.Errorf("failed to initialize container: %w", err)
			}
			defer container.Close()

			// Create transaction directly (validation included)
			tx, err := container.ValidatorService.UpdateOperator(validatorAddr, operator)
			if err != nil {
				return fmt.Errorf("failed to create transaction: %w", err)
			}

			// Sign if not unsigned
			if !commonFlags.Unsigned {
				tx, err = container.TxBuilder.SignTransaction(tx)
				if err != nil {
					return fmt.Errorf("failed to sign transaction: %w", err)
				}
			}

			// Convert to JSON and output
			txJSON, err := json.Marshal(tx)
			if err != nil {
				return fmt.Errorf("failed to convert transaction to JSON: %w", err)
			}

			if commonFlags.OutputFile != "" {
				return os.WriteFile(commonFlags.OutputFile, txJSON, 0644)
			}

			fmt.Println(string(txJSON))
			return nil
		},
	}

	// Add flags
	flags.AddCommonFlags(cmd, &commonFlags)
	cmd.Flags().StringVar(&validatorAddr, "validator", "", "Validator address")
	cmd.Flags().StringVar(&operator, "operator", "", "New operator address")
	cmd.MarkFlagRequired("validator")
	cmd.MarkFlagRequired("operator")

	return cmd
}

func newCreateValidatorUpdateRewardConfigCmd() *cobra.Command {
	var commonFlags flags.CommonFlags
	var validatorAddr, commissionRate string

	cmd := &cobra.Command{
		Use:   "update-reward-config",
		Short: "Create validator reward config update transaction",
		Long:  "Create a transaction to update validator reward configuration (without sending)",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return flags.ValidateSigningFlags(cmd)
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			cmd.SilenceUsage = true

			// Resolve configuration
			resolver, err := config.NewResolver()
			if err != nil {
				return err
			}
			resolvedConfig := resolver.ResolveFlags(&commonFlags)

			// Validate required fields
			// Create container
			container, err := container.NewContainer(resolvedConfig)
			if err != nil {
				return fmt.Errorf("failed to initialize container: %w", err)
			}
			defer container.Close()

			// Create transaction directly (validation included)
			tx, err := container.ValidatorService.UpdateRewardConfig(validatorAddr, commissionRate)
			if err != nil {
				return fmt.Errorf("failed to create transaction: %w", err)
			}

			// Sign if not unsigned
			if !commonFlags.Unsigned {
				tx, err = container.TxBuilder.SignTransaction(tx)
				if err != nil {
					return fmt.Errorf("failed to sign transaction: %w", err)
				}
			}

			// Convert to JSON and output
			txJSON, err := json.Marshal(tx)
			if err != nil {
				return fmt.Errorf("failed to convert transaction to JSON: %w", err)
			}

			if commonFlags.OutputFile != "" {
				return os.WriteFile(commonFlags.OutputFile, txJSON, 0644)
			}

			fmt.Println(string(txJSON))
			return nil
		},
	}

	// Add flags
	flags.AddCommonFlags(cmd, &commonFlags)
	cmd.Flags().StringVar(&validatorAddr, "validator", "", "Validator address")
	cmd.Flags().StringVar(&commissionRate, "commission-rate", "", "New commission rate in percentage (e.g., \"5%\")")
	cmd.MarkFlagRequired("validator")
	cmd.MarkFlagRequired("commission-rate")

	return cmd
}

func newCreateValidatorUpdateRewardManagerCmd() *cobra.Command {
	var commonFlags flags.CommonFlags
	var validatorAddr, rewardManager string

	cmd := &cobra.Command{
		Use:   "update-reward-manager",
		Short: "Create validator reward manager update transaction",
		Long:  "Create a transaction to update validator reward manager (without sending)",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return flags.ValidateSigningFlags(cmd)
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			cmd.SilenceUsage = true

			// Resolve configuration
			resolver, err := config.NewResolver()
			if err != nil {
				return err
			}
			resolvedConfig := resolver.ResolveFlags(&commonFlags)

			// Validate required fields
			// Create container
			container, err := container.NewContainer(resolvedConfig)
			if err != nil {
				return fmt.Errorf("failed to initialize container: %w", err)
			}
			defer container.Close()

			// Create transaction directly (validation included)
			tx, err := container.ValidatorService.UpdateRewardManager(validatorAddr, rewardManager)
			if err != nil {
				return fmt.Errorf("failed to create transaction: %w", err)
			}

			// Sign if not unsigned
			if !commonFlags.Unsigned {
				tx, err = container.TxBuilder.SignTransaction(tx)
				if err != nil {
					return fmt.Errorf("failed to sign transaction: %w", err)
				}
			}

			// Convert to JSON and output
			txJSON, err := json.Marshal(tx)
			if err != nil {
				return fmt.Errorf("failed to convert transaction to JSON: %w", err)
			}

			if commonFlags.OutputFile != "" {
				return os.WriteFile(commonFlags.OutputFile, txJSON, 0644)
			}

			fmt.Println(string(txJSON))
			return nil
		},
	}

	// Add flags
	flags.AddCommonFlags(cmd, &commonFlags)
	cmd.Flags().StringVar(&validatorAddr, "validator", "", "Validator address")
	cmd.Flags().StringVar(&rewardManager, "reward-manager", "", "New reward manager address")
	cmd.MarkFlagRequired("validator")
	cmd.MarkFlagRequired("reward-manager")

	return cmd
}

func newCreateValidatorUnjailCmd() *cobra.Command {
	var commonFlags flags.CommonFlags
	var validatorAddr string

	cmd := &cobra.Command{
		Use:   "unjail",
		Short: "Create validator unjail transaction",
		Long:  "Create a transaction to unjail a validator (without sending)",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return flags.ValidateSigningFlags(cmd)
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			cmd.SilenceUsage = true

			// Resolve configuration
			resolver, err := config.NewResolver()
			if err != nil {
				return err
			}
			resolvedConfig := resolver.ResolveFlags(&commonFlags)

			// Create container
			container, err := container.NewContainer(resolvedConfig)
			if err != nil {
				return fmt.Errorf("failed to initialize container: %w", err)
			}
			defer container.Close()

			// Create transaction directly (validation included)
			tx, err := container.ValidatorService.UnjailValidator(validatorAddr)
			if err != nil {
				return fmt.Errorf("failed to create transaction: %w", err)
			}

			// Sign if not unsigned
			if !commonFlags.Unsigned {
				tx, err = container.TxBuilder.SignTransaction(tx)
				if err != nil {
					return fmt.Errorf("failed to sign transaction: %w", err)
				}
			}

			// Convert to JSON and output
			txJSON, err := json.Marshal(tx)
			if err != nil {
				return fmt.Errorf("failed to convert transaction to JSON: %w", err)
			}

			if commonFlags.OutputFile != "" {
				return os.WriteFile(commonFlags.OutputFile, txJSON, 0644)
			}

			fmt.Println(string(txJSON))
			return nil
		},
	}

	// Add flags
	flags.AddCommonFlags(cmd, &commonFlags)
	cmd.Flags().StringVar(&validatorAddr, "validator", "", "Address of the validator to unjail")
	cmd.MarkFlagRequired("validator")

	return cmd
}

// Transaction data structure for output
type TransactionData struct {
	To       string `json:"to"`
	Value    string `json:"value"`
	Data     string `json:"data"`
	GasLimit uint64 `json:"gasLimit"`
	GasPrice string `json:"gasPrice"`
	ChainID  string `json:"chainId"`
}
