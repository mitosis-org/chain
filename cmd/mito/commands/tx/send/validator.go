package send

import (
	"github.com/mitosis-org/chain/cmd/mito/commands/tx/send/validator"
	"github.com/spf13/cobra"
)

// NewValidatorCmd returns validator commands for tx send
func NewValidatorCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "validator",
		Short: "Send validator transactions",
		Long: `Create, sign and send validator transactions to the network

This command provides access to validator transaction operations:
- validator create: Create and send a new validator transaction
- validator update-metadata: Send validator metadata update transaction
- validator update-operator: Send validator operator update transaction
- validator update-reward-config: Send validator reward config update transaction
- validator update-reward-manager: Send validator reward manager update transaction
- validator unjail: Send validator unjail transaction`,
	}

	// Add subcommands
	cmd.AddCommand(
		validator.NewCreateCmd(),
		validator.NewUpdateMetadataCmd(),
		validator.NewUpdateOperatorCmd(),
		validator.NewUpdateRewardConfigCmd(),
		validator.NewUpdateRewardManagerCmd(),
		validator.NewUnjailCmd(),
	)

	return cmd
}
