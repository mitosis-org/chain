package types

import (
	"encoding/json"
	"fmt"
)

// NewGenesisState creates a new GenesisState object
func NewGenesisState(
	params Params,
	validators []Validator,
	withdrawals []Withdrawal,
	lastValidatorPowers []LastValidatorPower,
) *GenesisState {
	return &GenesisState{
		Params:              params,
		Validators:          validators,
		Withdrawals:         withdrawals,
		LastValidatorPowers: lastValidatorPowers,
	}
}

// DefaultGenesisState creates a default GenesisState object
func DefaultGenesisState() *GenesisState {
	return &GenesisState{
		Params:              DefaultParams(),
		Validators:          []Validator{},
		Withdrawals:         []Withdrawal{},
		LastValidatorPowers: []LastValidatorPower{},
	}
}

// Validate performs basic genesis state validation
func (gs GenesisState) Validate() error {
	if err := gs.Params.Validate(); err != nil {
		return err
	}

	// Validate validators
	for i, validator := range gs.Validators {
		if len(validator.Pubkey) == 0 {
			return fmt.Errorf("validator %d has no pubkey", i)
		}
		if validator.Collateral.IsNil() || validator.Collateral.IsNegative() {
			return fmt.Errorf("validator %d has invalid collateral: %s", i, validator.Collateral)
		}
		if validator.ExtraVotingPower.IsNil() || validator.ExtraVotingPower.IsNegative() {
			return fmt.Errorf("validator %d has invalid extra voting power: %s", i, validator.ExtraVotingPower)
		}

		// For genesis validators, voting power should be computed
		validator.VotingPower = validator.ComputeVotingPower(gs.Params.MaxLeverageRatio)

		if validator.VotingPower.IsNil() || validator.VotingPower.IsNegative() {
			return fmt.Errorf("validator %d has invalid voting power: %s", i, validator.VotingPower)
		}
	}

	// Validate withdrawals
	for i, withdrawal := range gs.Withdrawals {
		if len(withdrawal.Pubkey) == 0 {
			return fmt.Errorf("withdrawal %d has no pubkey", i)
		}
		if withdrawal.Amount <= 0 {
			return fmt.Errorf("withdrawal %d has invalid amount: %d", i, withdrawal.Amount)
		}
		if withdrawal.Receiver == nil {
			return fmt.Errorf("withdrawal %d has no receiver", i)
		}
		if withdrawal.ReceivesAt == 0 {
			return fmt.Errorf("withdrawal %d has no receives_at timestamp", i)
		}
	}

	// Validate last validator powers
	for i, lastPower := range gs.LastValidatorPowers {
		if len(lastPower.Pubkey) == 0 {
			return fmt.Errorf("last validator power %d has no pubkey", i)
		}
		if lastPower.Power < 0 {
			return fmt.Errorf("last validator power %d has negative power: %d", i, lastPower.Power)
		}
	}

	return nil
}

// GetGenesisStateFromAppState returns the genesis state from the given app state
func GetGenesisStateFromAppState(appState map[string]json.RawMessage) GenesisState {
	var genesisState GenesisState

	if appState[ModuleName] != nil {
		_ = json.Unmarshal(appState[ModuleName], &genesisState)
	}

	return genesisState
}
