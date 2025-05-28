package flags

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
