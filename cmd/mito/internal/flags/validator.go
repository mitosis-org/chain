package flags

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

// MutuallyExclusiveGroup represents a group of mutually exclusive flags
type MutuallyExclusiveGroup struct {
	Name        string
	Description string
	Flags       []string
	Required    bool
}

// FlagValidator manages validation of mutually exclusive flags
type FlagValidator struct {
	groups []MutuallyExclusiveGroup
}

// NewFlagValidator creates a new flag validator
func NewFlagValidator() *FlagValidator {
	return &FlagValidator{
		groups: make([]MutuallyExclusiveGroup, 0),
	}
}

// AddMutuallyExclusiveGroup adds a mutually exclusive group
func (fv *FlagValidator) AddMutuallyExclusiveGroup(group MutuallyExclusiveGroup) {
	fv.groups = append(fv.groups, group)
}

// ValidateFlags validates all mutually exclusive groups for a command
func (fv *FlagValidator) ValidateFlags(cmd *cobra.Command) error {
	for _, group := range fv.groups {
		setFlags := make([]string, 0)

		// Check which flags in the group are set
		for _, flagName := range group.Flags {
			if cmd.Flags().Changed(flagName) {
				setFlags = append(setFlags, flagName)
			}
		}

		// Validate mutual exclusivity
		if len(setFlags) > 1 {
			return fmt.Errorf("flags %v are mutually exclusive (from group: %s)", setFlags, group.Name)
		}

		// Validate required constraint
		if group.Required && len(setFlags) == 0 {
			return fmt.Errorf("one of the following flags is required: %v (group: %s)", group.Flags, group.Name)
		}
	}

	return nil
}

// AddMutuallyExclusiveValidation is a helper function to add validation to any command
func AddMutuallyExclusiveValidation(cmd *cobra.Command, groups ...MutuallyExclusiveGroup) {
	validator := NewFlagValidator()

	for _, group := range groups {
		validator.AddMutuallyExclusiveGroup(group)
	}

	// Store the original PreRun
	existingPreRun := cmd.PreRun

	cmd.PreRun = func(cmd *cobra.Command, args []string) {
		// Validate mutually exclusive flags
		if err := validator.ValidateFlags(cmd); err != nil {
			fmt.Printf("Error: %v\n", err)

			// Provide helpful usage information
			for _, group := range groups {
				if containsAnyFlag(err.Error(), group.Flags) {
					fmt.Printf("\nGroup '%s': %s\n", group.Name, group.Description)
					fmt.Printf("  Use only one of: %v\n", group.Flags)
				}
			}
			os.Exit(1)
		}

		// Call existing PreRun if it exists
		if existingPreRun != nil {
			existingPreRun(cmd, args)
		}
	}
}

// ValidateSigningFlags validates signing-related flags for commands that require signing
func ValidateSigningFlags(cmd *cobra.Command) error {
	validator := NewFlagValidator()
	validator.AddMutuallyExclusiveGroup(SigningMethodGroup)
	return validator.ValidateFlags(cmd)
}

// containsAnyFlag checks if the error message contains any of the specified flags
func containsAnyFlag(errorMsg string, flags []string) bool {
	for _, flag := range flags {
		if strings.Contains(errorMsg, flag) {
			return true
		}
	}
	return false
}
