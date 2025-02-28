package types

import (
	"cosmossdk.io/errors"
)

// x/evmvalidator module sentinel errors
var (
	ErrValidatorNotFound       = errors.Register(ModuleName, 1, "validator not found")
	ErrValidatorAlreadyExists  = errors.Register(ModuleName, 2, "validator already exists")
	ErrValidatorJailed         = errors.Register(ModuleName, 3, "validator jailed")
	ErrValidatorPubKeyExists   = errors.Register(ModuleName, 4, "validator pubkey already exists")
	ErrInvalidPubKey           = errors.Register(ModuleName, 5, "invalid validator pubkey")
	ErrInvalidVotingPower      = errors.Register(ModuleName, 6, "invalid voting power")
	ErrInvalidCollateral       = errors.Register(ModuleName, 7, "invalid collateral amount")
	ErrInvalidExtraVotingPower = errors.Register(ModuleName, 8, "invalid extra voting power")
	ErrInsufficientCollateral  = errors.Register(ModuleName, 9, "insufficient collateral")
	ErrInvalidWithdrawal       = errors.Register(ModuleName, 10, "invalid withdrawal")
	ErrInvalidReceiver         = errors.Register(ModuleName, 11, "invalid receiver address")
	ErrInvalidTimestamp        = errors.Register(ModuleName, 12, "invalid timestamp")
	ErrInvalidParams           = errors.Register(ModuleName, 13, "invalid parameters")
)
