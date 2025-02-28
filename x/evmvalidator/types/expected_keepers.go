package types

import (
	"context"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"

	"cosmossdk.io/core/address"
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
)

// SlashingKeeper defines the expected slashing keeper
type SlashingKeeper interface {
	UnjailFromConsAddr(ctx context.Context, consAddr sdk.ConsAddress) error
}

// StakingKeeper defines the expected staking keeper (now implemented by Keeper)
type StakingKeeper interface {
	ValidatorAddressCodec() address.Codec
	ConsensusAddressCodec() address.Codec

	// iterate through validators by operator address, execute func for each validator
	IterateValidators(context.Context,
		func(index int64, validator ValidatorI) (stop bool)) error

	Validator(context.Context, sdk.ValAddress) (ValidatorI, error)            // get a particular validator by operator address
	ValidatorByConsAddr(context.Context, sdk.ConsAddress) (ValidatorI, error) // get a particular validator by consensus address

	// slash the validator and delegators of the validator, specifying offense height, offense power, and slash fraction
	Slash(context.Context, sdk.ConsAddress, int64, int64, math.LegacyDec) (math.Int, error)
	SlashWithInfractionReason(context.Context, sdk.ConsAddress, int64, int64, math.LegacyDec, stakingtypes.Infraction) (math.Int, error)
	Jail(context.Context, sdk.ConsAddress) error   // jail a validator
	Unjail(context.Context, sdk.ConsAddress) error // unjail a validator

	// MaxValidators returns the maximum amount of bonded validators
	MaxValidators(context.Context) (uint32, error)

	// IsValidatorJailed returns if the validator is jailed.
	IsValidatorJailed(ctx context.Context, addr sdk.ConsAddress) (bool, error)
}

// ValidatorI expected validator interface (implemented by Validator)
type ValidatorI interface {
	IsJailed() bool                          // whether the validator is jailed
	ConsPubKey() (cryptotypes.PubKey, error) // validation consensus pubkey (cryptotypes.PubKey)
	GetConsAddr() ([]byte, error)            // validation consensus address
}
