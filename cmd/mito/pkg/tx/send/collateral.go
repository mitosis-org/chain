package send

import (
	"fmt"

	"github.com/mitosis-org/chain/cmd/mito/internal/config"
	"github.com/mitosis-org/chain/cmd/mito/internal/container"
	"github.com/mitosis-org/chain/cmd/mito/internal/flags"
	"github.com/mitosis-org/chain/cmd/mito/internal/utils"
	"github.com/spf13/cobra"
)

// NewCollateralCmd returns collateral commands for tx send
func NewCollateralCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "collateral",
		Short: "Send collateral transactions",
		Long:  "Create, sign and send collateral transactions to the network",
	}

	cmd.AddCommand(
		newSendCollateralDepositCmd(),
		newSendCollateralWithdrawCmd(),
		newSendCollateralSetPermittedOwnerCmd(),
		newSendCollateralTransferOwnershipCmd(),
	)

	return cmd
}

func newSendCollateralDepositCmd() *cobra.Command {
	var commonFlags flags.CommonFlags
	var collateralFlags struct {
		validator string
		amount    string
	}

	cmd := &cobra.Command{
		Use:   "deposit",
		Short: "Send collateral deposit transaction",
		Long:  "Create, sign and send a transaction to deposit collateral",
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
			if err := validateSendCollateralDepositFields(resolvedConfig, &collateralFlags); err != nil {
				return err
			}

			// Confirm action
			if !utils.ConfirmAction(resolvedConfig.Yes, "Deposit collateral?") {
				return fmt.Errorf("operation cancelled")
			}

			// Create container
			container, err := container.NewContainer(resolvedConfig)
			if err != nil {
				return fmt.Errorf("failed to initialize container: %w", err)
			}
			defer container.Close()

			// Create unsigned transaction
			tx, err := container.CollateralService.DepositCollateral(collateralFlags.validator, collateralFlags.amount)
			if err != nil {
				return err
			}

			// Sign transaction
			signedTx, err := container.TxBuilder.SignTransaction(tx)
			if err != nil {
				return fmt.Errorf("failed to sign transaction: %w", err)
			}

			// Send transaction and wait for confirmation
			_, err = container.TxSender.SendAndWait(signedTx)
			if err != nil {
				return err
			}

			fmt.Println("✅ Collateral deposited successfully!")
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

func newSendCollateralWithdrawCmd() *cobra.Command {
	var commonFlags flags.CommonFlags
	var collateralFlags struct {
		validator string
		amount    string
		receiver  string
	}

	cmd := &cobra.Command{
		Use:   "withdraw",
		Short: "Send collateral withdraw transaction",
		Long:  "Create, sign and send a transaction to withdraw collateral",
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
			if err := validateSendCollateralWithdrawFields(resolvedConfig, &collateralFlags); err != nil {
				return err
			}

			// Confirm action
			if !utils.ConfirmAction(resolvedConfig.Yes, "Withdraw collateral?") {
				return fmt.Errorf("operation cancelled")
			}

			// Create container
			container, err := container.NewContainer(resolvedConfig)
			if err != nil {
				return fmt.Errorf("failed to initialize container: %w", err)
			}
			defer container.Close()

			// Create unsigned transaction
			tx, err := container.CollateralService.WithdrawCollateral(collateralFlags.validator, collateralFlags.amount, collateralFlags.receiver)
			if err != nil {
				return err
			}

			// Sign transaction
			signedTx, err := container.TxBuilder.SignTransaction(tx)
			if err != nil {
				return fmt.Errorf("failed to sign transaction: %w", err)
			}

			// Send transaction and wait for confirmation
			_, err = container.TxSender.SendAndWait(signedTx)
			if err != nil {
				return err
			}

			fmt.Println("✅ Collateral withdrawn successfully!")
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

func newSendCollateralSetPermittedOwnerCmd() *cobra.Command {
	var commonFlags flags.CommonFlags
	var collateralFlags struct {
		validator       string
		collateralOwner string
		isPermitted     bool
	}

	cmd := &cobra.Command{
		Use:   "set-permitted-owner",
		Short: "Send set permitted owner transaction",
		Long:  "Create, sign and send a transaction to set permitted collateral owner",
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
			if err := validateSendCollateralSetPermittedOwnerFields(resolvedConfig, &collateralFlags); err != nil {
				return err
			}

			// Confirm action
			permissionText := "deny"
			if collateralFlags.isPermitted {
				permissionText = "permit"
			}
			confirmMsg := fmt.Sprintf("Set collateral owner permission to %s?", permissionText)
			if !utils.ConfirmAction(resolvedConfig.Yes, confirmMsg) {
				return fmt.Errorf("operation cancelled")
			}

			// Create container
			container, err := container.NewContainer(resolvedConfig)
			if err != nil {
				return fmt.Errorf("failed to initialize container: %w", err)
			}
			defer container.Close()

			// Create unsigned transaction
			tx, err := container.CollateralService.SetPermittedCollateralOwner(collateralFlags.validator, collateralFlags.collateralOwner, collateralFlags.isPermitted)
			if err != nil {
				return err
			}

			// Sign transaction
			signedTx, err := container.TxBuilder.SignTransaction(tx)
			if err != nil {
				return fmt.Errorf("failed to sign transaction: %w", err)
			}

			// Send transaction and wait for confirmation
			_, err = container.TxSender.SendAndWait(signedTx)
			if err != nil {
				return err
			}

			fmt.Println("✅ Collateral owner permission updated successfully!")
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

func newSendCollateralTransferOwnershipCmd() *cobra.Command {
	var commonFlags flags.CommonFlags
	var collateralFlags struct {
		validator string
		newOwner  string
	}

	cmd := &cobra.Command{
		Use:   "transfer-ownership",
		Short: "Send transfer ownership transaction",
		Long:  "Create, sign and send a transaction to transfer collateral ownership",
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
			if err := validateSendCollateralTransferOwnershipFields(resolvedConfig, &collateralFlags); err != nil {
				return err
			}

			// Confirm action
			if !utils.ConfirmAction(resolvedConfig.Yes, "Transfer collateral ownership?") {
				return fmt.Errorf("operation cancelled")
			}

			// Create container
			container, err := container.NewContainer(resolvedConfig)
			if err != nil {
				return fmt.Errorf("failed to initialize container: %w", err)
			}
			defer container.Close()

			// Create unsigned transaction
			tx, err := container.CollateralService.TransferCollateralOwnership(collateralFlags.validator, collateralFlags.newOwner)
			if err != nil {
				return err
			}

			// Sign transaction
			signedTx, err := container.TxBuilder.SignTransaction(tx)
			if err != nil {
				return fmt.Errorf("failed to sign transaction: %w", err)
			}

			// Send transaction and wait for confirmation
			_, err = container.TxSender.SendAndWait(signedTx)
			if err != nil {
				return err
			}

			fmt.Println("✅ Collateral ownership transferred successfully!")
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

// Validation functions
func validateSendCollateralDepositFields(config *config.ResolvedConfig, collateralFlags *struct {
	validator string
	amount    string
}) error {
	if err := validateSendCollateralFields(config); err != nil {
		return err
	}

	if collateralFlags.validator == "" {
		return fmt.Errorf("validator address is required (use --validator)")
	}
	if collateralFlags.amount == "" {
		return fmt.Errorf("amount is required (use --amount)")
	}

	// Validate amount format
	if _, err := utils.ParseValueAsWei(collateralFlags.amount); err != nil {
		return fmt.Errorf("invalid amount format: %w", err)
	}

	// Validate validator address
	if _, err := utils.ValidateAddress(collateralFlags.validator); err != nil {
		return fmt.Errorf("invalid validator address: %w", err)
	}

	return nil
}

func validateSendCollateralWithdrawFields(config *config.ResolvedConfig, collateralFlags *struct {
	validator string
	amount    string
	receiver  string
}) error {
	if err := validateSendCollateralFields(config); err != nil {
		return err
	}

	if collateralFlags.validator == "" {
		return fmt.Errorf("validator address is required (use --validator)")
	}
	if collateralFlags.amount == "" {
		return fmt.Errorf("amount is required (use --amount)")
	}
	if collateralFlags.receiver == "" {
		return fmt.Errorf("receiver address is required (use --receiver)")
	}

	// Validate amount format
	if _, err := utils.ParseValueAsWei(collateralFlags.amount); err != nil {
		return fmt.Errorf("invalid amount format: %w", err)
	}

	// Validate addresses
	if _, err := utils.ValidateAddress(collateralFlags.validator); err != nil {
		return fmt.Errorf("invalid validator address: %w", err)
	}
	if _, err := utils.ValidateAddress(collateralFlags.receiver); err != nil {
		return fmt.Errorf("invalid receiver address: %w", err)
	}

	return nil
}

func validateSendCollateralSetPermittedOwnerFields(config *config.ResolvedConfig, collateralFlags *struct {
	validator       string
	collateralOwner string
	isPermitted     bool
}) error {
	if err := validateSendCollateralFields(config); err != nil {
		return err
	}

	if collateralFlags.validator == "" {
		return fmt.Errorf("validator address is required (use --validator)")
	}
	if collateralFlags.collateralOwner == "" {
		return fmt.Errorf("collateral owner address is required (use --collateral-owner)")
	}

	// Validate addresses
	if _, err := utils.ValidateAddress(collateralFlags.validator); err != nil {
		return fmt.Errorf("invalid validator address: %w", err)
	}
	if _, err := utils.ValidateAddress(collateralFlags.collateralOwner); err != nil {
		return fmt.Errorf("invalid collateral owner address: %w", err)
	}

	return nil
}

func validateSendCollateralTransferOwnershipFields(config *config.ResolvedConfig, collateralFlags *struct {
	validator string
	newOwner  string
}) error {
	if err := validateSendCollateralFields(config); err != nil {
		return err
	}

	if collateralFlags.validator == "" {
		return fmt.Errorf("validator address is required (use --validator)")
	}
	if collateralFlags.newOwner == "" {
		return fmt.Errorf("new owner address is required (use --new-owner)")
	}

	// Validate addresses
	if _, err := utils.ValidateAddress(collateralFlags.validator); err != nil {
		return fmt.Errorf("invalid validator address: %w", err)
	}
	if _, err := utils.ValidateAddress(collateralFlags.newOwner); err != nil {
		return fmt.Errorf("invalid new owner address: %w", err)
	}

	return nil
}

func validateSendCollateralFields(config *config.ResolvedConfig) error {
	if config.RpcURL == "" {
		return fmt.Errorf("RPC URL is required (use --rpc-url or set with 'mito config set-rpc')")
	}
	if config.ValidatorManagerContractAddr == "" {
		return fmt.Errorf("ValidatorManager contract address is required (use --contract or set with 'mito config set-contract')")
	}
	if !config.HasSigningMethod() {
		return fmt.Errorf("signing method is required (use --private-key or --keyfile)")
	}

	return nil
}
