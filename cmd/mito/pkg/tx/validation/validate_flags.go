package validation

import (
	"github.com/mitosis-org/chain/cmd/mito/internal/flags"
	"github.com/spf13/cobra"
)

func ValidateCreateTxFlagGroups(cmd *cobra.Command) error {
	groups := []flags.FlagGroup{
		SigningMethodGroup,
		TransactionTypeGroup,
		NetworkModeGroup,
	}

	return flags.ValidateGroups(cmd, groups)
}

func ValidateSendTxFlagGroups(cmd *cobra.Command) error {
	groups := []flags.FlagGroup{
		SigningMethodGroup,
		TransactionTypeGroup,
		RequireOnlineModeGroup,
	}
	return flags.ValidateGroups(cmd, groups)
}
