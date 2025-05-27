package cmd

import (
	"github.com/spf13/cobra"
)

// NewTxCmd returns the transaction command group
func NewTxCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "tx",
		Short: "Transaction commands",
		Long:  "Commands for creating and sending transactions",
	}

	cmd.AddCommand(
		newTxSendCmd(),
		newTxCreateCmd(),
	)

	return cmd
}

// NewValidatorCmd returns the validator command group (read-only)
func NewValidatorCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "validator",
		Short: "Validator query commands",
		Long:  "Commands for querying validator information",
	}

	cmd.AddCommand(newValidatorInfoCmd())

	return cmd
}

// NewConfigCmd returns the config command group
func NewConfigCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "config",
		Short: "Configuration commands",
		Long:  "Commands for managing configuration settings",
	}

	cmd.AddCommand(
		newSetRpcCmd(),
		newSetContractCmd(),
		newShowConfigCmd(),
	)

	return cmd
}

// newTxSendCmd returns the tx send command group
func newTxSendCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "send",
		Short: "Create, sign and send transactions",
		Long:  "Create, sign and immediately send transactions to the network",
	}

	cmd.AddCommand(
		newTxSendValidatorCmd(),
		newTxSendCollateralCmd(),
	)

	return cmd
}

// newTxCreateCmd returns the tx create command group
func newTxCreateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create transactions (signed or unsigned)",
		Long:  "Create transactions that can be signed immediately or later",
	}

	cmd.AddCommand(
		newTxCreateValidatorCmd(),
		newTxCreateCollateralCmd(),
	)

	return cmd
}

// newTxSendValidatorCmd returns validator commands for tx send
func newTxSendValidatorCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "validator",
		Short: "Validator transaction commands",
		Long:  "Commands for validator-related transactions",
	}

	cmd.AddCommand(
		newTxSendValidatorCreateCmd(),
		newTxSendValidatorUpdateOperatorCmd(),
		newTxSendValidatorUpdateMetadataCmd(),
		newTxSendValidatorUnjailCmd(),
	)

	return cmd
}

// newTxSendCollateralCmd returns collateral commands for tx send
func newTxSendCollateralCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "collateral",
		Short: "Collateral transaction commands",
		Long:  "Commands for collateral-related transactions",
	}

	cmd.AddCommand(
		newTxSendCollateralDepositCmd(),
		newTxSendCollateralWithdrawCmd(),
	)

	return cmd
}

// newTxCreateValidatorCmd returns validator commands for tx create
func newTxCreateValidatorCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "validator",
		Short: "Create validator transactions",
		Long:  "Create validator transactions (signed or unsigned)",
	}

	cmd.AddCommand(
		newTxCreateValidatorCreateCmd(),
		newTxCreateValidatorUpdateOperatorCmd(),
		newTxCreateValidatorUpdateMetadataCmd(),
		newTxCreateValidatorUnjailCmd(),
	)

	return cmd
}

// newTxCreateCollateralCmd returns collateral commands for tx create
func newTxCreateCollateralCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "collateral",
		Short: "Create collateral transactions",
		Long:  "Create collateral transactions (signed or unsigned)",
	}

	cmd.AddCommand(
		newTxCreateCollateralDepositCmd(),
		newTxCreateCollateralWithdrawCmd(),
	)

	return cmd
}
