package validation

import (
	"github.com/mitosis-org/chain/cmd/mito/internal/flags"
)

var (
	SigningMethodGroup = flags.MutuallyExclusiveGroup{
		Name:        "signing-method",
		Description: "Method for signing transactions",
		Flags: []flags.Flag{
			flags.NewStringFlag("private-key"),
			flags.NewGroupFlag(SigningWithKeyfileGroup),
		},
		Required: false, // Will be set to true for commands that require signing
	}

	KeyfilePasswordMethodGroup = flags.MutuallyExclusiveGroup{
		Name:        "keyfile-password-method",
		Description: "Method for open keyfile with password",
		Flags: []flags.Flag{
			flags.NewStringFlag("keyfile-password"),
			flags.NewStringFlag("keyfile-password-file"),
		},
		Required: true,
	}

	SigningWithKeyfileGroup = flags.DependentGroup{
		Name:        "signing-with-keyfile",
		Description: "Requires when signing with keyfiles",
		Flags: []flags.Flag{
			flags.NewStringFlag("keyfile"),
			flags.NewGroupFlag(KeyfilePasswordMethodGroup),
		},
		Required: false,
	}

	TransactionTypeGroup = flags.MutuallyExclusiveGroup{
		Name:        "transaction-type",
		Description: "Type of transaction to create",
		Flags: []flags.Flag{
			flags.NewStringFlag("signed"),
			flags.NewStringFlag("unsigned"),
		},
		Required: false,
	}

	OfflineModeGroup = flags.DependentGroup{
		Name:        "offline-mode",
		Description: "Requires when offline mode is enabled",
		Flags: []flags.Flag{
			flags.NewStringFlag("chain-id"),
			flags.NewStringFlag("gas-limit"),
			flags.NewStringFlag("gas-price"),
			flags.NewStringFlag("nonce"),
			flags.NewStringFlag("contract-fee"),
		},
		Required: false,
	}

	OnlineModeGroup = flags.DependentGroup{
		Name:        "online-mode",
		Description: "Requires when online mode is enabled",
		Flags: []flags.Flag{
			flags.NewStringFlag("rpc-url"),
			flags.NewStringFlag("validator-manager-contract-addr"),
		},
		Required: false,
	}

	NetworkModeGroup = flags.MutuallyExclusiveGroup{
		Name:        "network-mode",
		Description: "Mode of network to use",
		Flags: []flags.Flag{
			flags.NewGroupFlag(OfflineModeGroup),
			flags.NewGroupFlag(OnlineModeGroup),
		},
		Required: true,
	}

	RequireOnlineModeGroup = flags.DependentGroup{
		Name:        "require-online-mode",
		Description: "Required online mode flags",
		Flags: []flags.Flag{
			flags.NewStringFlag("rpc-url"),
			flags.NewStringFlag("validator-manager-contract-addr"),
		},
		Required: true,
	}

	// Output format group (for future use)
	OutputFormatGroup = flags.MutuallyExclusiveGroup{
		Name:        "output-format",
		Description: "Output format for transaction data",
		Flags: []flags.Flag{
			flags.NewStringFlag("json"),
			flags.NewStringFlag("raw"),
			flags.NewStringFlag("hex"),
		},
		Required: false,
	}
)
