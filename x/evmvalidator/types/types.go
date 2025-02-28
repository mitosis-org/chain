package types

import sdkmath "cosmossdk.io/math"

// NewLastValidatorPower creates a new LastValidatorPower instance
func NewLastValidatorPower(pubkey []byte, power int64) LastValidatorPower {
	return LastValidatorPower{
		Pubkey: pubkey,
		Power:  power,
	}
}

// NewWithdrawal creates a new Withdrawal instance
func NewWithdrawal(pubkey []byte, amount sdkmath.Int, receiver string, receivesAt uint64) Withdrawal {
	return Withdrawal{
		Pubkey:     pubkey,
		Amount:     amount,
		Receiver:   receiver,
		ReceivesAt: receivesAt,
	}
}
