package flags

// Common signing method groups
var (
	// Network information group for offline mode
	NetworkInfoGroup = MutuallyExclusiveGroup{
		Name:        "network-info",
		Description: "Network information for offline transaction creation",
		Flags: []Flag{
			NewStringFlag("rpc-url"),
			NewStringFlag("chain-id"),
		},
		Required: false, // Will be validated conditionally
	}
)
