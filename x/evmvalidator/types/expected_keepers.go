package types

import (
	"context"

	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	"github.com/ethereum/go-ethereum/common"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// SlashingKeeper defines the expected slashing keeper
type SlashingKeeper interface {
	AddPubkey(ctx context.Context, pubkey cryptotypes.PubKey) error
	UnjailFromConsAddr(ctx context.Context, consAddr sdk.ConsAddress) error

	//========= Replacements of Hooks =========//

	AfterValidatorBonded(ctx context.Context, consAddr sdk.ConsAddress) error
	AfterValidatorCreated(ctx context.Context, consPubKey cryptotypes.PubKey) error
}

type EvmEngineKeeper interface {
	InsertWithdrawal(ctx context.Context, withdrawalAddr common.Address, amountGwei uint64) error
}
