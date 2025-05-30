package flags

import (
	"github.com/spf13/cobra"
)

// StringFlagWithValue defines a string flag and sets its value
func StringFlagWithValue(cmd *cobra.Command, name string, value string, usage string) error {
	cmd.Flags().String(name, "", usage)
	return cmd.Flags().Set(name, value)
}
