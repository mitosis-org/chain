package flags

// Common signing method groups
var (
	SigningMethodGroup = MutuallyExclusiveGroup{
		Name:        "signing-method",
		Description: "Method for signing transactions",
		Flags: []Flag{
			NewStringFlag("private-key"),
			NewGroupFlag(SigningWithKeyfileGroup),
		},
		Required: false, // Will be set to true for commands that require signing
	}

	KeyfilePasswordMethodGroup = MutuallyExclusiveGroup{
		Name:        "keyfile-password-method",
		Description: "Method for open keyfile with password",
		Flags: []Flag{
			NewStringFlag("keyfile-password"),
			NewStringFlag("keyfile-password-file"),
		},
		Required: true,
	}

	SigningWithKeyfileGroup = DependentGroup{
		Name:        "signing-with-keyfile",
		Description: "Requires when signing with keyfiles",
		Flags: []Flag{
			NewStringFlag("keyfile"),
			NewGroupFlag(KeyfilePasswordMethodGroup),
		},
		Required: false,
	}

	// Transaction type group (for future use)
	TransactionTypeGroup = MutuallyExclusiveGroup{
		Name:        "transaction-type",
		Description: "Type of transaction to create",
		Flags: []Flag{
			NewStringFlag("signed"),
			NewStringFlag("unsigned"),
		},
		Required: false,
	}

	// Output format group (for future use)
	OutputFormatGroup = MutuallyExclusiveGroup{
		Name:        "output-format",
		Description: "Output format for transaction data",
		Flags: []Flag{
			NewStringFlag("json"),
			NewStringFlag("raw"),
			NewStringFlag("hex"),
		},
		Required: false,
	}

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
