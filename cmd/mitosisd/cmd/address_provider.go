package cmd

import (
	"github.com/cometbft/cometbft/crypto"
	evmengtypes "github.com/omni-network/omni/octane/evmengine/types"
)

var (
	_ evmengtypes.AddressProvider = &SimpleAddressProvider{}
)

type SimpleAddressProvider struct {
	pubKey crypto.PubKey
}

func (s SimpleAddressProvider) PubKey() crypto.PubKey {
	return s.pubKey
}
