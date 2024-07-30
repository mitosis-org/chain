package app

import (
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	evmengtypes "github.com/omni-network/omni/octane/evmengine/types"
)

var (
	_ evmengtypes.FeeRecipientProvider = &burnEVMFees{}

	// burnAddress is used as the EVM fee recipient resulting in burned execution fees.
	burnAddress = common.HexToAddress("0x000000000000000000000000000000000000dEaD")
)

// burnEVMFees is a fee recipient provider that burns all execution fees.
type burnEVMFees struct{}

func (burnEVMFees) LocalFeeRecipient() common.Address {
	return burnAddress
}

func (burnEVMFees) VerifyFeeRecipient(address common.Address) error {
	if address != burnAddress {
		return fmt.Errorf("fee recipient not the burn address. addr: %s", address.Hex())
	}

	return nil
}
