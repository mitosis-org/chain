package types

import (
	"fmt"
	mitotypes "github.com/mitosis-org/chain/types"
	"github.com/omni-network/omni/lib/errors"
	"github.com/omni-network/omni/lib/k1util"
)

// ValidatePubkeyWithEthAddress validates the public key format with the given address
func ValidatePubkeyWithEthAddress(pubkey []byte, addr mitotypes.EthAddress) error {
	expectedAddr, err := PubkeyToEthAddress(pubkey)
	if err != nil {
		return nil
	}

	if addr != expectedAddr {
		return errors.New("mismatched address",
			"pubkey", fmt.Sprintf("%X", pubkey),
			"expected", addr.String(), "actual", expectedAddr.String(),
		)
	}

	return nil
}

// PubkeyToEthAddress converts the given public key (33-byte compressed) to Ethereum address
func PubkeyToEthAddress(pubkey []byte) (mitotypes.EthAddress, error) {
	pbPubkey, err := k1util.PBPubKeyFromBytes(pubkey)
	if err != nil {
		return mitotypes.EthAddress{}, err
	}

	addr, err := k1util.PubKeyPBToAddress(pbPubkey)
	if err != nil {
		return mitotypes.EthAddress{}, err
	}

	return mitotypes.EthAddress(addr), nil
}
