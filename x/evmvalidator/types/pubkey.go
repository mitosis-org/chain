package types

import (
	ethcrypto "github.com/ethereum/go-ethereum/crypto"
	mitotypes "github.com/mitosis-org/chain/types"
	"github.com/omni-network/omni/lib/errors"
	"github.com/omni-network/omni/lib/k1util"
)

// ValidatePubkey validates the public key format
func ValidatePubkey(pubkey []byte) error {
	if len(pubkey) != 33 { // Compressed secp256k1 pubkey is 33 bytes
		return ErrInvalidPubKey
	}

	// Additional validation if needed
	// Try to convert to cosmos pubkey
	_, err := k1util.PubKeyBytesToCosmos(pubkey)
	if err != nil {
		return errors.Wrap(err, "invalid pubkey format")
	}

	return nil
}

// ValidatePubkeyWithEthAddress validates the public key format with the given address
func ValidatePubkeyWithEthAddress(pubkey []byte, addr mitotypes.EthAddress) error {
	if len(pubkey) != 33 { // Compressed secp256k1 pubkey is 33 bytes
		return ErrInvalidPubKey
	}

	derivedAddr, err := PubkeyToEthAddress(pubkey)
	if err != nil {
		return errors.Wrap(err, "failed to derive address from pubkey")
	}

	if addr != derivedAddr {
		return errors.New("address mismatch")
	}

	return nil
}

// PubkeyToEthAddress converts a 33-byte pubkey to a 20-byte EVM address. It decompresses the pubkey,
// applies keccak256, then take the last 20 bytes to get the corresponding EVM address.
func PubkeyToEthAddress(pubkey []byte) (mitotypes.EthAddress, error) {
	key, err := ethcrypto.DecompressPubkey(pubkey)
	if err != nil {
		return mitotypes.EthAddress{}, err
	}
	return mitotypes.EthAddress(ethcrypto.PubkeyToAddress(*key)), nil
}
