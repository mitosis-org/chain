package types

import (
	"fmt"
	mitotypes "github.com/mitosis-org/chain/types"
	"github.com/omni-network/omni/lib/errors"
	"github.com/omni-network/omni/lib/k1util"
)

// ValidatePubkeyWithEthAddress validates the public key format with the given address
func ValidatePubkeyWithEthAddress(pubkey []byte, addr mitotypes.EthAddress) error {
	pbPubkey, err := k1util.PBPubKeyFromBytes(pubkey)
	if err != nil {
		return err
	}

	expectedAddr, err := k1util.PubKeyPBToAddress(pbPubkey)
	if err != nil {
		return err
	}

	if addr != mitotypes.EthAddress(expectedAddr) {
		return errors.New("mismatched address",
			"pubkey", fmt.Sprintf("%X", pubkey),
			"expected", addr.String(), "actual", expectedAddr.String(),
		)
	}

	return nil
}
