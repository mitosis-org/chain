package validation

import (
	"fmt"

	"github.com/mitosis-org/chain/cmd/mito/internal/config"
	"github.com/mitosis-org/chain/cmd/mito/internal/utils"
	"github.com/spf13/cobra"
)

// ValidatorCreateFields contains fields for validator creation validation
type ValidatorCreateFields struct {
	PubKey            string
	Operator          string
	RewardManager     string
	CommissionRate    string
	Metadata          string
	InitialCollateral string
}

// ValidatorUpdateFields contains fields for validator update validation
type ValidatorUpdateFields struct {
	ValidatorAddress string
	FieldName        string
	NewValue         string
}

// ValidateCreateTxFlagGroups validates flag groups for create transactions
func ValidateCreateTxFlagGroups(cmd *cobra.Command) error {
	// Check mutually exclusive flags
	signed := cmd.Flags().Changed("signed")
	unsigned := cmd.Flags().Changed("unsigned")

	if signed && unsigned {
		return fmt.Errorf("cannot specify both --signed and --unsigned")
	}

	// Check signing method flags (mutually exclusive)
	privateKey := cmd.Flags().Changed("private-key")
	keyfile := cmd.Flags().Changed("keyfile")
	account := cmd.Flags().Changed("account")
	privValidatorKey := cmd.Flags().Changed("priv-validator-key")

	signingMethods := 0
	if privateKey {
		signingMethods++
	}
	if keyfile {
		signingMethods++
	}
	if account {
		signingMethods++
	}
	if privValidatorKey {
		signingMethods++
	}

	if signingMethods > 1 {
		return fmt.Errorf("cannot specify multiple signing methods (use only one of --private-key, --keyfile, --account, or --priv-validator-key)")
	}

	// Check keyfile password options
	keyfilePassword := cmd.Flags().Changed("keyfile-password")
	keyfilePasswordFile := cmd.Flags().Changed("keyfile-password-file")

	if keyfilePassword && keyfilePasswordFile {
		return fmt.Errorf("cannot specify both --keyfile-password and --keyfile-password-file")
	}

	// Check unsigned dependencies
	if unsigned {
		if !cmd.Flags().Changed("nonce") {
			return fmt.Errorf("--nonce is required when --unsigned is specified")
		}
	}

	// Check if signing is required for signed transactions
	if signed || (!signed && !unsigned) { // Default to signed
		if signingMethods == 0 {
			return fmt.Errorf("signing method is required for signed transactions (use --private-key, --keyfile, --account, or --priv-validator-key)")
		}
	}

	return nil
}

// ValidateSendTxFlagGroups validates flag groups for send transactions
func ValidateSendTxFlagGroups(cmd *cobra.Command) error {
	// Send transactions are always signed, so signing method is required
	privateKey := cmd.Flags().Changed("private-key")
	keyfile := cmd.Flags().Changed("keyfile")
	account := cmd.Flags().Changed("account")
	privValidatorKey := cmd.Flags().Changed("priv-validator-key")

	signingMethods := 0
	if privateKey {
		signingMethods++
	}
	if keyfile {
		signingMethods++
	}
	if account {
		signingMethods++
	}
	if privValidatorKey {
		signingMethods++
	}

	if signingMethods == 0 {
		return fmt.Errorf("signing method is required for sending transactions (use --private-key, --keyfile, --account, or --priv-validator-key)")
	}

	if signingMethods > 1 {
		return fmt.Errorf("cannot specify multiple signing methods (use only one of --private-key, --keyfile, --account, or --priv-validator-key)")
	}

	// Check keyfile password options
	keyfilePassword := cmd.Flags().Changed("keyfile-password")
	keyfilePasswordFile := cmd.Flags().Changed("keyfile-password-file")

	if keyfilePassword && keyfilePasswordFile {
		return fmt.Errorf("cannot specify both --keyfile-password and --keyfile-password-file")
	}

	return nil
}

// ValidateValidatorCreateFields validates validator creation fields
func ValidateValidatorCreateFields(config *config.ResolvedConfig, fields *ValidatorCreateFields) error {
	// Validate required fields
	if fields.PubKey == "" {
		return fmt.Errorf("public key is required (use --pubkey)")
	}
	if fields.Operator == "" {
		return fmt.Errorf("operator address is required (use --operator)")
	}
	if fields.RewardManager == "" {
		return fmt.Errorf("reward manager address is required (use --reward-manager)")
	}
	if fields.CommissionRate == "" {
		return fmt.Errorf("commission rate is required (use --commission-rate)")
	}
	if fields.Metadata == "" {
		return fmt.Errorf("metadata is required (use --metadata)")
	}
	if fields.InitialCollateral == "" {
		return fmt.Errorf("initial collateral is required (use --initial-collateral)")
	}

	// Validate address formats
	if _, err := utils.ValidateAddress(fields.Operator); err != nil {
		return fmt.Errorf("invalid operator address: %w", err)
	}
	if _, err := utils.ValidateAddress(fields.RewardManager); err != nil {
		return fmt.Errorf("invalid reward manager address: %w", err)
	}

	return nil
}

// ValidateValidatorUpdateFields validates validator update fields
func ValidateValidatorUpdateFields(config *config.ResolvedConfig, fields *ValidatorUpdateFields) error {
	if fields.ValidatorAddress == "" {
		return fmt.Errorf("validator address is required (use --validator)")
	}
	if fields.NewValue == "" {
		return fmt.Errorf("new %s is required", fields.FieldName)
	}

	// Validate validator address format
	if _, err := utils.ValidateAddress(fields.ValidatorAddress); err != nil {
		return fmt.Errorf("invalid validator address: %w", err)
	}

	return nil
}

// ValidateNetworkFields validates network-related fields
func ValidateNetworkFields(config *config.ResolvedConfig, requireNetwork bool) error {
	if requireNetwork {
		if config.RPCURL == "" {
			return fmt.Errorf("RPC URL is required (use --rpc-url or set with 'mito config set-rpc')")
		}
		if config.ValidatorManagerContractAddr == "" {
			return fmt.Errorf("ValidatorManager contract address is required (use --contract or set with 'mito config set-contract')")
		}
	}
	return nil
}

// ValidateSigningFields validates signing-related fields
func ValidateSigningFields(config *config.ResolvedConfig, requireSigning bool) error {
	if requireSigning && !config.HasSigningMethod() {
		return fmt.Errorf("signing method is required (use --private-key, --keyfile, --account, or --priv-validator-key)")
	}
	return nil
}
