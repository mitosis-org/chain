package version

import (
	"fmt"

	"github.com/spf13/cobra"
)

// Version information
var (
	Version   = "0.0.1"
	GitCommit = "unknown"
	BuildDate = "unknown"
)

// NewVersionCmd returns the version command
func NewVersionCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "version",
		Short: "Show version information",
		Long:  "Display version, git commit, and build date information",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("mito version %s\n", Version)
			fmt.Printf("Git commit: %s\n", GitCommit)
			fmt.Printf("Build date: %s\n", BuildDate)
		},
	}

	return cmd
}
