package version

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/version"
	"github.com/spf13/cobra"
)

// NewVersionCmd returns the version command
func NewVersionCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "version",
		Short: "Show version information",
		Long:  "Display version, git commit, and build date information",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("mito version %s\n", version.Version)
			fmt.Printf("Git commit: %s\n", version.Commit)
			fmt.Printf("Build date: %s\n", version.BuildTags)
		},
	}

	return cmd
}
