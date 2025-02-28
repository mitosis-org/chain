package types

import (
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	"github.com/omni-network/omni/lib/k1util"
)

var _ ValidatorI = (*Validator)(nil)

// IsJailed implements ValidatorI
func (v Validator) IsJailed() bool {
	return v.Jailed
}

func (v Validator) ConsPubKey() (cryptotypes.PubKey, error) {
	return k1util.PubKeyBytesToCosmos(v.Pubkey)
}

// GetConsAddr implements ValidatorI
func (v Validator) GetConsAddr() ([]byte, error) {
	consAddr, err := v.ConsAddr()
	if err != nil {
		return nil, err
	}
	return consAddr.Bytes(), nil
}
