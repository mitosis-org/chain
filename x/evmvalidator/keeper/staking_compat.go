package keeper

import (
	"cosmossdk.io/core/address"
)

// This file implements the StakingKeeper interface expected by other modules like x/slashing.
// It provides compatibility functions for modules that expect x/staking behavior.

// ValidatorAddressCodec returns the address codec for validators
func (k Keeper) ValidatorAddressCodec() address.Codec {
	return k.validatorAddressCodec
}

// ConsensusAddressCodec returns the address codec for consensus nodes
func (k Keeper) ConsensusAddressCodec() address.Codec {
	return k.consensusAddressCodec
}
