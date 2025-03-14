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

	// ValidatorKeyPrefix is the prefix for a validator
	ValidatorKeyPrefix = []byte{0x02}

	// ValidatorByConsAddrKeyPrefix is the prefix for a validator index, by consensus address
	ValidatorByConsAddrKeyPrefix = []byte{0x03}

	// ValidatorByPowerIndexKeyPrefix is the prefix for a validator index, sorted by power
	ValidatorByPowerIndexKeyPrefix = []byte{0x04}

	// LastValidatorPowerKeyPrefix is the prefix for last validator powers
	LastValidatorPowerKeyPrefix = []byte{0x05}

	// WithdrawalQueueKeyPrefix is the prefix for the withdrawal queue
	WithdrawalQueueKeyPrefix = []byte{0x06}

	// WithdrawalByValidatorKeyPrefix is the prefix for an index to withdrawals by validator and maturesAt
	WithdrawalByValidatorKeyPrefix = []byte{0x07}
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
	binary.BigEndian.PutUint64(powerBytes, uint64(^power))
	return append(ValidatorByPowerIndexKeyPrefix, append(powerBytes, address.MustLengthPrefix(valAddr.Bytes())...)...)
}

// GetLastValidatorPowerKey creates key for a validator from address
func GetLastValidatorPowerKey(valAddr mitotypes.EthAddress) []byte {
	return append(LastValidatorPowerKeyPrefix, address.MustLengthPrefix(valAddr.Bytes())...)
}

// GetWithdrawalQueueKey creates key for withdrawals at a timestamp
func GetWithdrawalQueueKey(maturesAt int64) []byte {
	maturesAtBytes := make([]byte, 8)
	binary.BigEndian.PutUint64(maturesAtBytes, uint64(maturesAt))
	return append(WithdrawalQueueKeyPrefix, maturesAtBytes...)
}

// GetWithdrawalByValidatorKey creates a key for indexing withdrawals by validator and maturesAt
func GetWithdrawalByValidatorKey(valAddr mitotypes.EthAddress, maturesAt int64) []byte {
	maturesAtBytes := make([]byte, 8)
	binary.BigEndian.PutUint64(maturesAtBytes, uint64(maturesAt))
	return append(append(WithdrawalByValidatorKeyPrefix, address.MustLengthPrefix(valAddr.Bytes())...), maturesAtBytes...)
}

// GetWithdrawalByValidatorIterationKey creates a key for iterating withdrawals by validator and maturesAt
func GetWithdrawalByValidatorIterationKey(valAddr mitotypes.EthAddress) []byte {
	return append(WithdrawalByValidatorKeyPrefix, address.MustLengthPrefix(valAddr.Bytes())...)
}
