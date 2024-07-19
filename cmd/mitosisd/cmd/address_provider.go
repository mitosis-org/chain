package cmd

import (
	"github.com/cometbft/cometbft/crypto"
	"github.com/cometbft/cometbft/crypto/ed25519"
	evmengtypes "github.com/omni-network/omni/octane/evmengine/types"
)

var (
	_ evmengtypes.AddressProvider = &NoAddressProvider{}
	_ evmengtypes.AddressProvider = &SimpleAddressProvider{}
)

type NoAddressProvider struct{}

func (s NoAddressProvider) PubKey() crypto.PubKey {
	dummy := [ed25519.PubKeySize]byte{}
	return ed25519.PubKey(dummy[:])
}

type SimpleAddressProvider struct {
	pubKey crypto.PubKey
}

func (s SimpleAddressProvider) PubKey() crypto.PubKey {
	return s.pubKey
}
