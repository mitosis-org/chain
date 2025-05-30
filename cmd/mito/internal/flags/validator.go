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

// ValidateSigningFlags validates signing-related flags for commands that require signing
func ValidateSigningFlags(cmd *cobra.Command) error {
	validator := NewFlagValidator()
	validator.AddGroup(SigningMethodGroup)
	return validator.ValidateFlags(cmd)
}
