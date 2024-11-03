package app

import (
	"github.com/ethereum/go-ethereum/common"
	evmengtypes "github.com/omni-network/omni/octane/evmengine/types"
)

var (
	_ evmengtypes.AddressProvider      = &ValidatorAddressProvider{}
	_ evmengtypes.FeeRecipientProvider = &ValidatorAddressProvider{}
)

type ValidatorAddressProvider struct {
	Addr common.Address
}

func (s ValidatorAddressProvider) LocalAddress() common.Address {
	return s.Addr
}

func (s ValidatorAddressProvider) LocalFeeRecipient() common.Address {
	return s.Addr
}

func (s ValidatorAddressProvider) VerifyFeeRecipient(address common.Address) error {
	return nil
}
