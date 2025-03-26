package types

import (
	"cosmossdk.io/errors"
)

// x/evmvalidator module sentinel errors
var (
	ErrValidatorNotFound      = errors.Register(ModuleName, 1, "validator not found")
	ErrValidatorAlreadyExists = errors.Register(ModuleName, 2, "validator already exists")
	ErrInvalidPubKey          = errors.Register(ModuleName, 3, "invalid validator pubkey")
	ErrInvalidVotingPower     = errors.Register(ModuleName, 4, "invalid voting power")
	ErrInsufficientCollateral = errors.Register(ModuleName, 5, "insufficient collateral")
)
