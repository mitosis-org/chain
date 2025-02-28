package types

import (
	"encoding/binary"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/omni-network/omni/lib/k1util"
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

	// ValidatorKeyPrefix is the prefix for validator store
	ValidatorKeyPrefix = []byte{0x02}

	// ValidatorPowerRankStoreKeyPrefix is the prefix for validator power rank store
	ValidatorPowerRankStoreKeyPrefix = []byte{0x03}

	// LastValidatorPowerKeyPrefix is the prefix for last validator powers
	LastValidatorPowerKeyPrefix = []byte{0x04}

	// WithdrawalQueueKeyPrefix is the prefix for the withdrawal queue
	WithdrawalQueueKeyPrefix = []byte{0x05}
)

// GetValidatorKey creates key for a validator from consensus public key
func GetValidatorKey(pubkey []byte) []byte {
	return append(ValidatorKeyPrefix, pubkey...)
}

// GetValidatorPowerRankKey creates the key for the validator power rank store from power and pubkey
func GetValidatorPowerRankKey(power int64, pubkey []byte) []byte {
	// NOTE: power is the voting power, not the tokens amount
	powerBytes := make([]byte, 8)
	// power is converted to descending order for the key (higher power first)
	// because we want to iterate from highest to lowest power in EndBlocker
	binary.BigEndian.PutUint64(powerBytes, uint64(^power))
	return append(ValidatorPowerRankStoreKeyPrefix, append(powerBytes, pubkey...)...)
}

// GetLastValidatorPowerKey creates key for a validator from pubkey
func GetLastValidatorPowerKey(pubkey []byte) []byte {
	return append(LastValidatorPowerKeyPrefix, pubkey...)
}

// GetWithdrawalQueueKey creates key for withdrawals at a timestamp
func GetWithdrawalQueueKey(receivesAt uint64) []byte {
	receivesAtBytes := make([]byte, 8)
	binary.BigEndian.PutUint64(receivesAtBytes, receivesAt)
	return append(WithdrawalQueueKeyPrefix, receivesAtBytes...)
}

// GetWithdrawalQueueKeyByTime returns the key for all withdrawals made for a given time
func GetWithdrawalQueueKeyByTime(receivesAt uint64) []byte {
	receivesAtBytes := make([]byte, 8)
	binary.BigEndian.PutUint64(receivesAtBytes, receivesAt)
	return append(WithdrawalQueueKeyPrefix, receivesAtBytes...)
}

// GetAddressFromConsensusPublicKey derives an SDK consensus address from a consensus public key
func GetAddressFromConsensusPublicKey(pubkey []byte) (sdk.ConsAddress, error) {
	// Convert pubkey to cosmos format
	cosmosKey, err := k1util.PubKeyBytesToCosmos(pubkey)
	if err != nil {
		return nil, err
	}
	return sdk.ConsAddress(cosmosKey.Address()), nil
}
