package testutil

import (
	"context"

	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/mitosis-org/chain/x/evmvalidator/types"
)

// MockSlashingKeeper is a mock of SlashingKeeper interface for testing
type MockSlashingKeeper struct {
	AddPubkeyFn             func(ctx context.Context, pubkey cryptotypes.PubKey) error
	UnjailFromConsAddrFn    func(ctx context.Context, consAddr sdk.ConsAddress) error
	AfterValidatorBondedFn  func(ctx context.Context, consAddr sdk.ConsAddress) error
	AfterValidatorCreatedFn func(ctx context.Context, consPubKey cryptotypes.PubKey) error
}

func (m MockSlashingKeeper) AddPubkey(ctx context.Context, pubkey cryptotypes.PubKey) error {
	if m.AddPubkeyFn != nil {
		return m.AddPubkeyFn(ctx, pubkey)
	}
	return nil
}

func (m MockSlashingKeeper) UnjailFromConsAddr(ctx context.Context, consAddr sdk.ConsAddress) error {
	if m.UnjailFromConsAddrFn != nil {
		return m.UnjailFromConsAddrFn(ctx, consAddr)
	}
	return nil
}

func (m MockSlashingKeeper) AfterValidatorBonded(ctx context.Context, consAddr sdk.ConsAddress) error {
	if m.AfterValidatorBondedFn != nil {
		return m.AfterValidatorBondedFn(ctx, consAddr)
	}
	return nil
}

func (m MockSlashingKeeper) AfterValidatorCreated(ctx context.Context, consPubKey cryptotypes.PubKey) error {
	if m.AfterValidatorCreatedFn != nil {
		return m.AfterValidatorCreatedFn(ctx, consPubKey)
	}
	return nil
}

// MockEvmEngineKeeper is a mock of EvmEngineKeeper interface for testing
type MockEvmEngineKeeper struct {
	InsertWithdrawalFn func(ctx context.Context, withdrawalAddr common.Address, amountGwei uint64) error
}

func (m MockEvmEngineKeeper) InsertWithdrawal(ctx context.Context, withdrawalAddr common.Address, amountGwei uint64) error {
	if m.InsertWithdrawalFn != nil {
		return m.InsertWithdrawalFn(ctx, withdrawalAddr, amountGwei)
	}
	return nil
}

var (
	_ types.SlashingKeeper  = MockSlashingKeeper{}
	_ types.EvmEngineKeeper = MockEvmEngineKeeper{}
)
