package types

import (
	"encoding/binary"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/address"
	mitotypes "github.com/mitosis-org/chain/types"
)

const (
	// ModuleName defines the module name
	ModuleName = "evmvalidator"

	// StoreKey defines the primary module store key
	StoreKey = ModuleName
)

var (
	// ParamsKey is the key for module parameters
	ParamsKey = []byte{0x01}

	// ValidatorEntrypointContractAddrKey is the prefix for a ConsensusValidatorEntrypoint contract address
	ValidatorEntrypointContractAddrKey = []byte{0x02}

	// ValidatorKeyPrefix is the prefix for a validator
	ValidatorKeyPrefix = []byte{0x03}

	// ValidatorByConsAddrKeyPrefix is the prefix for a validator index, by consensus address
	ValidatorByConsAddrKeyPrefix = []byte{0x04}

	// ValidatorByPowerIndexKeyPrefix is the prefix for a validator index, sorted by power
	ValidatorByPowerIndexKeyPrefix = []byte{0x05}

	// LastValidatorPowerKeyPrefix is the prefix for last validator powers
	LastValidatorPowerKeyPrefix = []byte{0x06}

	// WithdrawalLastIDKeyPrefix is the key for the last withdrawal ID
	WithdrawalLastIDKeyPrefix = []byte{0x07}

	// WithdrawalByMaturesAtKeyPrefix is the prefix for a withdrawal by maturesAt and ID
	WithdrawalByMaturesAtKeyPrefix = []byte{0x08}

	// WithdrawalByValidatorKeyPrefix is the prefix for a withdrawal by validator address, maturesAt, and ID
	WithdrawalByValidatorKeyPrefix = []byte{0x09}

	// CollateralOwnershipKeyPrefix is the prefix for a collateral ownership by validator and owner
	CollateralOwnershipKeyPrefix = []byte{0x0A}
)

// GetValidatorKey creates key for a validator from validator address
func GetValidatorKey(valAddr mitotypes.EthAddress) []byte {
	return append(ValidatorKeyPrefix, address.MustLengthPrefix(valAddr.Bytes())...)
}

// GetValidatorByConsAddrKey creates key for a validator from consensus address
func GetValidatorByConsAddrKey(consAddr sdk.ConsAddress) []byte {
	return append(ValidatorByConsAddrKeyPrefix, address.MustLengthPrefix(consAddr)...)
}

// GetValidatorByPowerIndexKey creates the key for a validator from power and address
func GetValidatorByPowerIndexKey(power int64, valAddr mitotypes.EthAddress) []byte {
	// NOTE: power is the voting power, not the tokens amount
	powerBytes := make([]byte, 8)
	// power is converted to descending order for the key (higher power first)
	// because we want to iterate from highest to lowest power in EndBlocker
	binary.BigEndian.PutUint64(powerBytes, uint64(^power)) //nolint:gosec
	return append(ValidatorByPowerIndexKeyPrefix, append(powerBytes, address.MustLengthPrefix(valAddr.Bytes())...)...)
}

// GetLastValidatorPowerKey creates key for a validator from address
func GetLastValidatorPowerKey(valAddr mitotypes.EthAddress) []byte {
	return append(LastValidatorPowerKeyPrefix, address.MustLengthPrefix(valAddr.Bytes())...)
}

// GetWithdrawalLastIDKey creates key for a withdrawal from ID
func GetWithdrawalLastIDKey() []byte {
	return WithdrawalLastIDKeyPrefix
}

// GetWithdrawalByMaturesAtKey creates a key for a withdrawal by maturesAt and ID
func GetWithdrawalByMaturesAtKey(maturesAt int64, id uint64) []byte {
	maturesAtBytes := make([]byte, 8)
	binary.BigEndian.PutUint64(maturesAtBytes, uint64(maturesAt)) //nolint:gosec
	idBytes := make([]byte, 8)
	binary.BigEndian.PutUint64(idBytes, id)
	return append(WithdrawalByMaturesAtKeyPrefix, append(maturesAtBytes, idBytes...)...)
}

// GetWithdrawalByValidatorKey creates a key for a withdrawal by validator and maturesAt
func GetWithdrawalByValidatorKey(valAddr mitotypes.EthAddress, maturesAt int64, id uint64) []byte {
	maturesAtBytes := make([]byte, 8)
	binary.BigEndian.PutUint64(maturesAtBytes, uint64(maturesAt)) //nolint:gosec
	idBytes := make([]byte, 8)
	binary.BigEndian.PutUint64(idBytes, id)
	return append(WithdrawalByValidatorKeyPrefix, append(address.MustLengthPrefix(valAddr.Bytes()), append(maturesAtBytes, idBytes...)...)...)
}

// GetWithdrawalByValidatorIterationKey creates a key for iterating withdrawals by validator and maturesAt
func GetWithdrawalByValidatorIterationKey(valAddr mitotypes.EthAddress) []byte {
	return append(WithdrawalByValidatorKeyPrefix, address.MustLengthPrefix(valAddr.Bytes())...)
}

// GetCollateralOwnershipKey creates a key for collateral ownership by validator and owner
func GetCollateralOwnershipKey(valAddr mitotypes.EthAddress, owner mitotypes.EthAddress) []byte {
	return append(
		CollateralOwnershipKeyPrefix,
		append(
			address.MustLengthPrefix(valAddr.Bytes()),
			address.MustLengthPrefix(owner.Bytes())...,
		)...,
	)
}

// GetCollateralOwnershipByValidatorIterationKey creates a key for iterating collateral ownerships by validator
func GetCollateralOwnershipByValidatorIterationKey(valAddr mitotypes.EthAddress) []byte {
	return append(
		CollateralOwnershipKeyPrefix,
		address.MustLengthPrefix(valAddr.Bytes())...,
	)
}
