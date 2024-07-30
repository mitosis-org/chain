package app

import (
	"sync"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	Bech32Prefix = "mito"

	Bech32PrefixAccAddr  = Bech32Prefix
	Bech32PrefixAccPub   = Bech32Prefix + "pub"
	Bech32PrefixValAddr  = Bech32Prefix + "valoper"
	Bech32PrefixValPub   = Bech32Prefix + "valoperpub"
	Bech32PrefixConsAddr = Bech32Prefix + "valcons"
	Bech32PrefixConsPub  = Bech32Prefix + "valconspub"

	Bip44CoinType uint32 = 118 // TODO(thai): consider to change it to 60.
	Bip44Purpose  uint32 = 44
)

var initConfig sync.Once

// SetupConfig sets up the Cosmos SDK configuration to be compatible with the semantics of ethereum.
func SetupConfig() {
	initConfig.Do(
		func() {
			// set the address prefixes
			config := sdk.GetConfig()
			SetBech32Prefixes(config)
			SetBip44CoinType(config)
			config.Seal()
		},
	)
}

// SetBech32Prefixes sets the global prefixes to be used when serializing addresses and public keys to Bech32 strings.
func SetBech32Prefixes(config *sdk.Config) {
	config.SetBech32PrefixForAccount(Bech32PrefixAccAddr, Bech32PrefixAccPub)
	config.SetBech32PrefixForValidator(Bech32PrefixValAddr, Bech32PrefixValPub)
	config.SetBech32PrefixForConsensusNode(Bech32PrefixConsAddr, Bech32PrefixConsPub)
}

// SetBip44CoinType sets the global coin type to be used in hierarchical deterministic wallets.
func SetBip44CoinType(config *sdk.Config) {
	config.SetPurpose(Bip44Purpose)
	config.SetCoinType(Bip44CoinType)
}
