package types

var _ ValidatorI = (*Validator)(nil)

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
