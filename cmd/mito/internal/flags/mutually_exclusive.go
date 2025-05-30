package flags

import (
	"fmt"

	"github.com/spf13/cobra"
)

// MutuallyExclusiveGroup represents a group of mutually exclusive flags
type MutuallyExclusiveGroup struct {
	Name        string
	Description string
	Flags       []Flag
	Required    bool
}

func (g MutuallyExclusiveGroup) GetName() string {
	return g.Name
}

func (g MutuallyExclusiveGroup) GetFlags() []string {
	flags := make([]string, 0)
	for _, flag := range g.Flags {
		flags = append(flags, flag.GetFlags()...)
	}
	return flags
}

func (g MutuallyExclusiveGroup) IsRequired() bool {
	return g.Required
}

func (g MutuallyExclusiveGroup) Validate(cmd *cobra.Command) error {
	setFlags := make([]string, 0)
	setGroups := make([]string, 0)

	// Check which flags or groups in the group are set
	for _, flag := range g.Flags {
		if flag.IsGroup() {
			// For nested groups, check if any of their flags are set
			for _, nestedFlag := range flag.group.GetFlags() {
				if cmd.Flags().Changed(nestedFlag) {
					setGroups = append(setGroups, flag.group.GetName())
					break
				}
			}
			// Validate the nested group
			if err := flag.group.Validate(cmd); err != nil {
				return fmt.Errorf("in group %s: %w", g.Name, err)
			}
		} else {
			if cmd.Flags().Changed(flag.name) {
				setFlags = append(setFlags, flag.name)
			}
		}
	}

	// Combine set flags and groups for mutual exclusivity check
	totalSet := append(setFlags, setGroups...)

	// Validate mutual exclusivity
	if len(totalSet) > 1 {
		return fmt.Errorf("flags/groups %v are mutually exclusive (from group: %s)", totalSet, g.Name)
	}

	// Validate required constraint
	if g.Required && len(totalSet) == 0 {
		return fmt.Errorf("one of the following is required: %v (group: %s)", g.GetFlags(), g.Name)
	}

	return nil
}
