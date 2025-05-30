package flags

import (
	"fmt"

	"github.com/spf13/cobra"
)

// DependentGroup represents a group of flags where some flags depend on others
type DependentGroup struct {
	Name        string
	Description string
	Flags       []Flag
	Required    bool
}

func (g DependentGroup) GetName() string {
	return g.Name
}

func (g DependentGroup) GetFlags() []string {
	flags := make([]string, 0)
	for _, flag := range g.Flags {
		flags = append(flags, flag.GetFlags()...)
	}
	return flags
}

func (g DependentGroup) IsRequired() bool {
	return g.Required
}

func (g DependentGroup) Validate(cmd *cobra.Command) error {
	// Check if any primary flag is set
	primaryFlagSet := false
	for _, flag := range g.Flags {
		if !flag.IsGroup() && cmd.Flags().Changed(flag.name) {
			primaryFlagSet = true
			break
		}
	}

	// If no primary flag is set and the group is not required, skip validation
	if !primaryFlagSet && !g.Required {
		return nil
	}

	// If primary flag is set, validate all dependent groups
	if primaryFlagSet {
		for _, flag := range g.Flags {
			if flag.IsGroup() {
				if err := flag.group.Validate(cmd); err != nil {
					return fmt.Errorf("in group %s: %w", g.Name, err)
				}
			}
		}
	}

	return nil
}
