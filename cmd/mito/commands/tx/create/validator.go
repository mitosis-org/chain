package create

import (
	"github.com/mitosis-org/chain/cmd/mito/commands/tx/create/validator"
	"github.com/spf13/cobra"
)

// NewValidatorCmd returns validator commands for tx create
func NewValidatorCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "validator",
		Short: "Create validator transactions",
		Long: `Create validator transactions (signed or unsigned)

This command provides access to validator transaction creation operations:
- validator create: Create a new validator transaction
- validator update-metadata: Create validator metadata update transaction
- validator update-operator: Create validator operator update transaction
- validator update-reward-config: Create validator reward config update transaction
- validator update-reward-manager: Create validator reward manager update transaction
- validator unjail: Create validator unjail transaction`,
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
