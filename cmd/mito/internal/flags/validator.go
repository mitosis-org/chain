package flags

import (
	"github.com/spf13/cobra"
)

// FlagValidator manages validation of flag groups
type FlagValidator struct {
	groups []FlagGroup
}

// NewFlagValidator creates a new flag validator
func NewFlagValidator() *FlagValidator {
	return &FlagValidator{
		groups: make([]FlagGroup, 0),
	}
}

func NewFlagValidatorWithGroups(groups []FlagGroup) *FlagValidator {
	return &FlagValidator{
		groups: groups,
	}
}

func ValidateGroups(cmd *cobra.Command, groups []FlagGroup) error {
	return NewFlagValidatorWithGroups(groups).ValidateFlags(cmd)
}

// AddGroup adds a flag group
func (fv *FlagValidator) AddGroup(group FlagGroup) {
	fv.groups = append(fv.groups, group)
}

// ValidateFlags validates all flag groups for a command
func (fv *FlagValidator) ValidateFlags(cmd *cobra.Command) error {
	for _, group := range fv.groups {
		if err := group.Validate(cmd); err != nil {
			return err
		}
	}
	return nil
}

func (fv *FlagValidator) ValidateGroups(cmd *cobra.Command, groups []FlagGroup) error {
	fv.groups = groups
	return fv.ValidateFlags(cmd)
}
