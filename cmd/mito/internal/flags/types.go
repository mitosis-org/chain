package flags

import "github.com/spf13/cobra"

// FlagGroup represents a group of flags with validation rules
type FlagGroup interface {
	Validate(cmd *cobra.Command) error
	GetName() string
	GetFlags() []string
	IsRequired() bool
}

// Flag represents either a string flag name or a FlagGroup
type Flag struct {
	name  string    // flag name if it's a string flag
	group FlagGroup // group if it's a nested group
}

// NewStringFlag creates a Flag from a string flag name
func NewStringFlag(name string) Flag {
	return Flag{name: name}
}

// NewGroupFlag creates a Flag from a FlagGroup
func NewGroupFlag(group FlagGroup) Flag {
	return Flag{group: group}
}

// IsGroup returns true if the Flag represents a group
func (f Flag) IsGroup() bool {
	return f.group != nil
}

// String returns the flag name or group name
func (f Flag) String() string {
	if f.IsGroup() {
		return f.group.GetName()
	}
	return f.name
}

// GetFlags returns all flag names, either the single flag or all flags from a group
func (f Flag) GetFlags() []string {
	if f.IsGroup() {
		return f.group.GetFlags()
	}
	return []string{f.name}
}
