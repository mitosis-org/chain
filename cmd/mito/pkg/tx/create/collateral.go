package create

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/mitosis-org/chain/cmd/mito/internal/config"
	"github.com/mitosis-org/chain/cmd/mito/internal/container"
	"github.com/mitosis-org/chain/cmd/mito/internal/flags"
	"github.com/spf13/cobra"
)

// NewCollateralCmd returns collateral commands for tx create
func NewCollateralCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "collateral",
		Short: "Create collateral transactions",
		Long:  "Create collateral transactions (signed or unsigned)",
	}

	cmd.AddCommand(
		newCreateCollateralDepositCmd(),
		newCreateCollateralWithdrawCmd(),
		newCreateCollateralSetPermittedOwnerCmd(),
		newCreateCollateralTransferOwnershipCmd(),
	)

	return cmd
}

func newCreateCollateralDepositCmd() *cobra.Command {
	var commonFlags flags.CommonFlags
	var collateralFlags struct {
		validator string
		amount    string
	}

	cmd := &cobra.Command{
		Use:   "deposit",
		Short: "Create collateral deposit transaction",
		Long:  "Create a transaction to deposit collateral (without sending)",
		RunE: func(cmd *cobra.Command, args []string) error {
			cmd.SilenceUsage = true

			// Resolve configuration
			resolver, err := config.NewResolver()
			if err != nil {
				return err
			}
			resolvedConfig := resolver.ResolveFlags(&commonFlags)

			// Initialize container
			c, err := container.NewContainer(resolvedConfig)
			if err != nil {
				return fmt.Errorf("failed to initialize container: %w", err)
			}
			defer c.Close()

			tx, err := c.CollateralService.DepositCollateral(collateralFlags.validator, collateralFlags.amount)
			if err != nil {
				return err
			}

			// Sign if not unsigned
			if !commonFlags.Unsigned {
				tx, err = c.TxBuilder.SignTransaction(tx)
				if err != nil {
					return fmt.Errorf("failed to sign transaction: %w", err)
				}
			}

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
	cmd.Flags().StringVar(&collateralFlags.validator, "validator", "", "Validator address")
	cmd.Flags().StringVar(&collateralFlags.amount, "amount", "", "Amount to deposit in MITO (e.g., \"1.5\")")

	// Mark required flags
	cmd.MarkFlagRequired("validator")
	cmd.MarkFlagRequired("amount")

	return cmd
}

func newCreateCollateralWithdrawCmd() *cobra.Command {
	var commonFlags flags.CommonFlags
	var collateralFlags struct {
		validator string
		amount    string
		receiver  string
	}

	cmd := &cobra.Command{
		Use:   "withdraw",
		Short: "Create collateral withdraw transaction",
		Long:  "Create a transaction to withdraw collateral (without sending)",
		RunE: func(cmd *cobra.Command, args []string) error {
			cmd.SilenceUsage = true

			// Resolve configuration
			resolver, err := config.NewResolver()
			if err != nil {
				return err
			}
			resolvedConfig := resolver.ResolveFlags(&commonFlags)

			// Initialize container
			c, err := container.NewContainer(resolvedConfig)
			if err != nil {
				return fmt.Errorf("failed to initialize container: %w", err)
			}
			defer c.Close()

			tx, err := c.CollateralService.WithdrawCollateral(collateralFlags.validator, collateralFlags.amount, collateralFlags.receiver)
			if err != nil {
				return err
			}

			// Sign if not unsigned
			if !commonFlags.Unsigned {
				tx, err = c.TxBuilder.SignTransaction(tx)
				if err != nil {
					return fmt.Errorf("failed to sign transaction: %w", err)
				}
			}

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
	cmd.Flags().StringVar(&collateralFlags.validator, "validator", "", "Validator address")
	cmd.Flags().StringVar(&collateralFlags.amount, "amount", "", "Amount to withdraw in MITO (e.g., \"1.5\")")
	cmd.Flags().StringVar(&collateralFlags.receiver, "receiver", "", "Address to receive the withdrawn collateral")

	// Mark required flags
	cmd.MarkFlagRequired("validator")
	cmd.MarkFlagRequired("amount")
	cmd.MarkFlagRequired("receiver")

	return cmd
}

func newCreateCollateralSetPermittedOwnerCmd() *cobra.Command {
	var commonFlags flags.CommonFlags
	var collateralFlags struct {
		validator       string
		collateralOwner string
		isPermitted     bool
	}

	cmd := &cobra.Command{
		Use:   "set-permitted-owner",
		Short: "Create set permitted owner transaction",
		Long:  "Create a transaction to set permitted collateral owner (without sending)",
		RunE: func(cmd *cobra.Command, args []string) error {
			cmd.SilenceUsage = true

			// Resolve configuration
			resolver, err := config.NewResolver()
			if err != nil {
				return err
			}
			resolvedConfig := resolver.ResolveFlags(&commonFlags)

			// Initialize container
			c, err := container.NewContainer(resolvedConfig)
			if err != nil {
				return fmt.Errorf("failed to initialize container: %w", err)
			}
			defer c.Close()

			tx, err := c.CollateralService.SetPermittedCollateralOwner(collateralFlags.validator, collateralFlags.collateralOwner, collateralFlags.isPermitted)
			if err != nil {
				return err
			}

			// Sign if not unsigned
			if !commonFlags.Unsigned {
				tx, err = c.TxBuilder.SignTransaction(tx)
				if err != nil {
					return fmt.Errorf("failed to sign transaction: %w", err)
				}
			}

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
	cmd.Flags().StringVar(&collateralFlags.validator, "validator", "", "Validator address")
	cmd.Flags().StringVar(&collateralFlags.collateralOwner, "collateral-owner", "", "Collateral owner address to permit or deny")
	cmd.Flags().BoolVar(&collateralFlags.isPermitted, "permitted", false, "Whether to permit (true) or deny (false) the address")

	// Mark required flags
	cmd.MarkFlagRequired("validator")
	cmd.MarkFlagRequired("collateral-owner")

	return cmd
}

func newCreateCollateralTransferOwnershipCmd() *cobra.Command {
	var commonFlags flags.CommonFlags
	var collateralFlags struct {
		validator string
		newOwner  string
	}

	cmd := &cobra.Command{
		Use:   "transfer-ownership",
		Short: "Create transfer ownership transaction",
		Long:  "Create a transaction to transfer collateral ownership (without sending)",
		RunE: func(cmd *cobra.Command, args []string) error {
			cmd.SilenceUsage = true

			// Resolve configuration
			resolver, err := config.NewResolver()
			if err != nil {
				return err
			}
			resolvedConfig := resolver.ResolveFlags(&commonFlags)

			// Initialize container
			c, err := container.NewContainer(resolvedConfig)
			if err != nil {
				return fmt.Errorf("failed to initialize container: %w", err)
			}
			defer c.Close()

			tx, err := c.CollateralService.TransferCollateralOwnership(collateralFlags.validator, collateralFlags.newOwner)
			if err != nil {
				return err
			}

			// Sign if not unsigned
			if !commonFlags.Unsigned {
				tx, err = c.TxBuilder.SignTransaction(tx)
				if err != nil {
					return fmt.Errorf("failed to sign transaction: %w", err)
				}
			}

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
	cmd.Flags().StringVar(&collateralFlags.validator, "validator", "", "Validator address")
	cmd.Flags().StringVar(&collateralFlags.newOwner, "new-owner", "", "New owner address")

	// Mark required flags
	cmd.MarkFlagRequired("validator")
	cmd.MarkFlagRequired("new-owner")

	return cmd
}
