package cmd

import (
	"fmt"
	"log"
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

// Common signing method groups
var (
	SigningMethodGroup = MutuallyExclusiveGroup{
		Name:        "signing-method",
		Description: "Method for signing transactions",
		Flags:       []string{"private-key", "keyfile"},
		Required:    false, // Will be set to true for commands that require signing
	}

	// Transaction type group (for future use)
	TransactionTypeGroup = MutuallyExclusiveGroup{
		Name:        "transaction-type",
		Description: "Type of transaction to create",
		Flags:       []string{"signed", "unsigned"},
		Required:    false,
	}

	// Output format group (for future use)
	OutputFormatGroup = MutuallyExclusiveGroup{
		Name:        "output-format",
		Description: "Output format for transaction data",
		Flags:       []string{"json", "raw", "hex"},
		Required:    false,
	}

	// Network information group for offline mode
	NetworkInfoGroup = MutuallyExclusiveGroup{
		Name:        "network-info",
		Description: "Network information for offline transaction creation",
		Flags:       []string{"rpc-url", "chain-id"},
		Required:    false, // Will be validated conditionally
	}
)

// mustMarkFlagRequired marks a flag as required and panics if it fails
func mustMarkFlagRequired(cmd *cobra.Command, flag string) {
	if err := cmd.MarkFlagRequired(flag); err != nil {
		log.Fatalf("Failed to mark flag '%s' as required: %v", flag, err)
	}
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

// containsAnyFlag checks if the error message contains any of the specified flags
func containsAnyFlag(errorMsg string, flags []string) bool {
	for _, flag := range flags {
		if strings.Contains(errorMsg, flag) {
			return true
		}
	}
	return false
}
