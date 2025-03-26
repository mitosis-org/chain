package types

import (
	"fmt"

	mitotypes "github.com/mitosis-org/chain/types"

	"github.com/omni-network/omni/lib/errors"
)

// NewGenesisState creates a new GenesisState object
func NewGenesisState(
	params Params,
	validatorEntrypointContractAddr mitotypes.EthAddress,
	validators []Validator,
	withdrawals []Withdrawal,
	lastValidatorPowers []LastValidatorPower,
) *GenesisState {
	return &GenesisState{
		Params:                          params,
		ValidatorEntrypointContractAddr: validatorEntrypointContractAddr,
		Validators:                      validators,
		Withdrawals:                     withdrawals,
		LastValidatorPowers:             lastValidatorPowers,
	}
}

// DefaultGenesisState creates a default GenesisState object
func DefaultGenesisState() *GenesisState {
	return &GenesisState{
		Params:                          DefaultParams(),
		ValidatorEntrypointContractAddr: mitotypes.EthAddress{},
		Validators:                      []Validator{},
		Withdrawals:                     []Withdrawal{},
		LastValidatorPowers:             []LastValidatorPower{},
	}
}

// Validate performs basic genesis state validation
func (gs GenesisState) Validate() error {
	if err := gs.Params.Validate(); err != nil {
		return err
	}

	// Validate validators
	for i, validator := range gs.Validators {
		if err := ValidatePubkeyWithEthAddress(validator.Pubkey, validator.Addr); err != nil {
			return errors.Wrap(err, fmt.Sprintf("validator %d has not matched addr and pubkey: %s, %X", i, validator.Addr.String(), validator.Pubkey))
		}
		if validator.Collateral.IsNil() {
			return fmt.Errorf("validator %d has invalid collateral: %s", i, validator.Collateral)
		}
		if validator.ExtraVotingPower.IsNil() {
			return fmt.Errorf("validator %d has invalid extra voting power: %s", i, validator.ExtraVotingPower)
		}

		// NOTE: voting power will be recomputed in InitGenesis
	}

	// Validate withdrawals
	for i, withdrawal := range gs.Withdrawals {
		if withdrawal.Amount <= 0 {
			return fmt.Errorf("withdrawal %d has invalid amount: %d", i, withdrawal.Amount)
		}
		if withdrawal.MaturesAt == 0 {
			return fmt.Errorf("withdrawal %d has no matures_at timestamp", i)
		}
		if withdrawal.CreationHeight == 0 {
			return fmt.Errorf("withdrawal %d has no creation_height", i)
		}
	}

	// Validate last validator powers
	for i, lastPower := range gs.LastValidatorPowers {
		if lastPower.Power < 0 {
			return fmt.Errorf("last validator power %d has negative power: %d", i, lastPower.Power)
		}
	}

	return nil
}
