package types

import (
	evidencetypes "cosmossdk.io/x/evidence/types"
	slashingtypes "github.com/cosmos/cosmos-sdk/x/slashing/types"
)

var (
	_ slashingtypes.ValidatorI = (*Validator)(nil)
	_ evidencetypes.ValidatorI = (*Validator)(nil)
)

// GetConsAddr implements ValidatorI
func (v Validator) GetConsAddr() ([]byte, error) {
	consAddr, err := v.ConsAddr()
	if err != nil {
		return nil, err
	}
	return consAddr.Bytes(), nil
}

// IsJailed implements ValidatorI
func (v Validator) IsJailed() bool {
	return v.Jailed
}
