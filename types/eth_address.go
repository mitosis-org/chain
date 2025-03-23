package types

import (
	"bytes"
	"fmt"

	"github.com/ethereum/go-ethereum/common"
)

// EthAddress is a wrapper for common.Address that implements proto custom type interfaces
type EthAddress common.Address

// BytesToEthAddress converts bytes to an EthAddress
func BytesToEthAddress(b []byte) EthAddress {
	return EthAddress(common.BytesToAddress(b))
}

// Marshal converts the EthAddress to bytes for protobuf serialization
func (a EthAddress) Marshal() ([]byte, error) {
	return common.Address(a).Bytes(), nil
}

// Unmarshal sets the EthAddress from bytes from protobuf deserialization
func (a *EthAddress) Unmarshal(data []byte) error {
	if len(data) != common.AddressLength {
		return fmt.Errorf("invalid address length: got %d, want %d", len(data), common.AddressLength)
	}
	copy((*a)[:], data)
	return nil
}

// MarshalTo implements the protobuf marshaler interface
func (a EthAddress) MarshalTo(data []byte) (int, error) {
	copy(data, a[:])
	return common.AddressLength, nil
}

// MarshalJSON implements the json.Marshaler interface
func (a EthAddress) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf(`"%s"`, a.String())), nil
}

// UnmarshalJSON implements the json.Unmarshaler interface
func (a *EthAddress) UnmarshalJSON(data []byte) error {
	if len(data) < 2 || data[0] != '"' || data[len(data)-1] != '"' {
		return fmt.Errorf("invalid address format, expected quoted string")
	}

	// Remove quotes
	hexStr := string(data[1 : len(data)-1])

	// Validate 0x prefix
	if len(hexStr) < 2 || hexStr[:2] != "0x" {
		return fmt.Errorf("invalid address format, expected 0x prefix")
	}

	// Validate length (excluding the 0x prefix)
	if len(hexStr) != (common.AddressLength*2 + 2) {
		return fmt.Errorf("invalid address length: got %d chars (with 0x prefix), want %d",
			len(hexStr), common.AddressLength*2+2)
	}

	// Convert hex string to address
	*a = EthAddress(common.HexToAddress(hexStr))

	return nil
}

// Size returns the size of the EthAddress in bytes
func (a EthAddress) Size() int {
	return common.AddressLength
}

// Equal compares two EthAddresses for equality
func (a EthAddress) Equal(other EthAddress) bool {
	return bytes.Equal(a[:], other[:])
}

// Address returns the common.Address representation
func (a EthAddress) Address() common.Address {
	return common.Address(a)
}

// String returns the string representation of the EthAddress
func (a EthAddress) String() string {
	return a.Address().String()
}

// Bytes returns the byte representation of the EthAddress
func (a EthAddress) Bytes() []byte {
	return a.Address().Bytes()
}
