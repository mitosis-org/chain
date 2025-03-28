package keeper_test

import (
	"context"
	"testing"

	"cosmossdk.io/math"
	cryptocodec "github.com/cosmos/cosmos-sdk/crypto/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	mitotypes "github.com/mitosis-org/chain/types"
	"github.com/mitosis-org/chain/x/evmvalidator/testutil"
	"github.com/mitosis-org/chain/x/evmvalidator/types"
	"github.com/stretchr/testify/suite"
)

// ValidatorSetTestSuite is a test suite to be used with validator set tests
type ValidatorSetTestSuite struct {
	suite.Suite
	tk testutil.TestKeeper
}

// SetupTest initializes the test suite
func (s *ValidatorSetTestSuite) SetupTest() {
	s.tk = testutil.CreateTestInput(&s.Suite)
}

// TestValidatorSetTestSuite runs the validator set test suite
func TestValidatorSetTestSuite(t *testing.T) {
	suite.Run(t, new(ValidatorSetTestSuite))
}

// setupTestParams sets up test parameters
func (s *ValidatorSetTestSuite) setupTestParams() types.Params {
	params := types.Params{
		MaxValidators:    10,
		MaxLeverageRatio: math.LegacyNewDec(10), // 10x leverage
		MinVotingPower:   1,
		WithdrawalLimit:  10,
	}
	err := s.tk.Keeper.SetParams(s.tk.Ctx, params)
	s.Require().NoError(err)
	return params
}

// registerValidator is a helper function to register a validator
func (s *ValidatorSetTestSuite) registerValidator(collateral math.Uint, extraVotingPower math.Uint, jailed bool) ([]byte, mitotypes.EthAddress, types.Validator) {
	_, pubkey, ethAddr := testutil.GenerateSecp256k1Key()
	err := s.tk.Keeper.RegisterValidator(s.tk.Ctx, ethAddr, pubkey, collateral, extraVotingPower, jailed)
	s.Require().NoError(err)

	validator, found := s.tk.Keeper.GetValidator(s.tk.Ctx, ethAddr)
	s.Require().True(found)
	return pubkey, ethAddr, validator
}

func (s *ValidatorSetTestSuite) Test_ApplyAndReturnValidatorSetUpdates_NewValidators() {
	// Set test parameters
	s.setupTestParams()

	// Register validators
	_, addr1, validator1 := s.registerValidator(math.NewUint(5000000000), math.ZeroUint(), false) // 5 MITO, power = 5
	_, addr2, validator2 := s.registerValidator(math.NewUint(3000000000), math.ZeroUint(), false) // 3 MITO, power = 3
	_, addr3, validator3 := s.registerValidator(math.NewUint(2000000000), math.ZeroUint(), false) // 2 MITO, power = 2

	// Mock slashing keeper hooks
	bondedValidators := make(map[string]bool)
	s.tk.MockSlash.AfterValidatorBondedFn = func(ctx context.Context, consAddr sdk.ConsAddress) error {
		bondedValidators[consAddr.String()] = true
		return nil
	}

	// Apply validator set updates
	updates, err := s.tk.Keeper.ApplyAndReturnValidatorSetUpdates(s.tk.Ctx)
	s.Require().NoError(err)

	// Check number of updates (should be 3 validators)
	s.Require().Equal(3, len(updates))

	// Verify updates have correct validators
	var foundAddrs = make(map[string]bool)
	for _, update := range updates {
		// Convert the ABCI public key to SDK public key
		pk, err := cryptocodec.FromCmtProtoPublicKey(update.PubKey)
		s.Require().NoError(err)

		// Get validator by pubkey
		var validator types.Validator
		var found bool
		s.tk.Keeper.IterateValidators_(s.tk.Ctx, func(_ int64, val types.Validator) bool {
			valPk, err := val.ConsPubKey()
			if err != nil {
				return false
			}

			if pk.Equals(valPk) {
				validator = val
				found = true
				return true
			}
			return false
		})
		s.Require().True(found, "validator not found for pubkey")

		// Record found addresses
		foundAddrs[validator.Addr.String()] = true

		// Check power
		expectedPower := validator.ConsensusVotingPower()
		s.Require().Equal(expectedPower, update.Power)
	}

	// Verify all validators were included
	s.Require().True(foundAddrs[addr1.String()])
	s.Require().True(foundAddrs[addr2.String()])
	s.Require().True(foundAddrs[addr3.String()])

	// Check last validator powers were updated
	power1, found := s.tk.Keeper.GetLastValidatorPower(s.tk.Ctx, addr1)
	s.Require().True(found)
	s.Require().Equal(int64(5), power1)

	power2, found := s.tk.Keeper.GetLastValidatorPower(s.tk.Ctx, addr2)
	s.Require().True(found)
	s.Require().Equal(int64(3), power2)

	power3, found := s.tk.Keeper.GetLastValidatorPower(s.tk.Ctx, addr3)
	s.Require().True(found)
	s.Require().Equal(int64(2), power3)

	// Check validators are bonded
	updatedValidator1, found := s.tk.Keeper.GetValidator(s.tk.Ctx, addr1)
	s.Require().True(found)
	s.Require().True(updatedValidator1.Bonded)

	updatedValidator2, found := s.tk.Keeper.GetValidator(s.tk.Ctx, addr2)
	s.Require().True(found)
	s.Require().True(updatedValidator2.Bonded)

	updatedValidator3, found := s.tk.Keeper.GetValidator(s.tk.Ctx, addr3)
	s.Require().True(found)
	s.Require().True(updatedValidator3.Bonded)

	// Check slashing keeper was called for each validator
	consAddr1, err := validator1.ConsAddr()
	s.Require().NoError(err)
	consAddr2, err := validator2.ConsAddr()
	s.Require().NoError(err)
	consAddr3, err := validator3.ConsAddr()
	s.Require().NoError(err)

	s.Require().Equal(3, len(bondedValidators))
	s.Require().True(bondedValidators[consAddr1.String()])
	s.Require().True(bondedValidators[consAddr2.String()])
	s.Require().True(bondedValidators[consAddr3.String()])

	// More important is to verify that the validators are actually bonded in state
	for _, addr := range []mitotypes.EthAddress{addr1, addr2, addr3} {
		val, found := s.tk.Keeper.GetValidator(s.tk.Ctx, addr)
		s.Require().True(found)
		s.Require().True(val.Bonded, "Validator should be bonded")
	}
}

// Test_ApplyAndReturnValidatorSetUpdates_NoChanges tests that no updates are returned when there are no changes
func (s *ValidatorSetTestSuite) Test_ApplyAndReturnValidatorSetUpdates_PowerChange() {
	// Set test parameters
	s.setupTestParams()

	// Register validators
	_, addr1, _ := s.registerValidator(math.NewUint(5000000000), math.ZeroUint(), false) // 5 MITO, power = 5
	_, addr2, _ := s.registerValidator(math.NewUint(3000000000), math.ZeroUint(), false) // 3 MITO, power = 3

	// Initial update
	initialUpdates, err := s.tk.Keeper.ApplyAndReturnValidatorSetUpdates(s.tk.Ctx)
	s.Require().NoError(err)
	s.Require().Equal(2, len(initialUpdates))

	// Call again with no changes
	updates, err := s.tk.Keeper.ApplyAndReturnValidatorSetUpdates(s.tk.Ctx)
	s.Require().NoError(err)
	s.Require().Equal(0, len(updates), "should have no updates when there are no changes")

	// Change a validator power
	validator, found := s.tk.Keeper.GetValidator(s.tk.Ctx, addr1)
	s.Require().True(found)

	// Deposit more collateral to change power
	s.tk.Keeper.DepositCollateral(s.tk.Ctx, &validator, math.NewUint(2000000000))
	s.tk.Keeper.UpdateValidatorState(s.tk.Ctx, &validator, "testing")

	// Apply updates with changes
	updates, err = s.tk.Keeper.ApplyAndReturnValidatorSetUpdates(s.tk.Ctx)
	s.Require().NoError(err)
	s.Require().Equal(1, len(updates), "should have one update for the changed validator")

	// Verify the correct validator was updated
	var foundUpdate bool
	for _, update := range updates {
		pk, err := cryptocodec.FromCmtProtoPublicKey(update.PubKey)
		s.Require().NoError(err)

		valPk, err := validator.ConsPubKey()
		s.Require().NoError(err)

		if pk.Equals(valPk) {
			foundUpdate = true
			s.Require().Equal(int64(7), update.Power) // 5 + 2 = 7 MITO
		}
	}
	s.Require().True(foundUpdate, "update for validator 1 not found")

	// Check last validator power was updated
	power1, found := s.tk.Keeper.GetLastValidatorPower(s.tk.Ctx, addr1)
	s.Require().True(found)
	s.Require().Equal(int64(7), power1)

	// Validator 2 power should remain unchanged
	power2, found := s.tk.Keeper.GetLastValidatorPower(s.tk.Ctx, addr2)
	s.Require().True(found)
	s.Require().Equal(int64(3), power2)
}

// Test_ApplyAndReturnValidatorSetUpdates_JailedValidator tests validators are excluded when jailed
func (s *ValidatorSetTestSuite) Test_ApplyAndReturnValidatorSetUpdates_JailedValidator() {
	// Set test parameters
	s.setupTestParams()

	// Register validators - one normal, one jailed
	_, addr1, _ := s.registerValidator(math.NewUint(5000000000), math.ZeroUint(), false)          // 5 MITO, power = 5
	_, addr2, validator2 := s.registerValidator(math.NewUint(3000000000), math.ZeroUint(), false) // 3 MITO, power = 3

	// Initial update
	initialUpdates, err := s.tk.Keeper.ApplyAndReturnValidatorSetUpdates(s.tk.Ctx)
	s.Require().NoError(err)
	s.Require().Equal(2, len(initialUpdates))

	// Jail the second validator
	s.tk.Keeper.Jail_(s.tk.Ctx, &validator2, "testing jail")

	// Apply updates
	updates, err := s.tk.Keeper.ApplyAndReturnValidatorSetUpdates(s.tk.Ctx)
	s.Require().NoError(err)

	// Should see one update - the jailed validator with zero power
	s.Require().Equal(1, len(updates), "should have one update for the jailed validator")

	// Verify the update sets power to zero for the jailed validator
	var foundUpdate bool
	for _, update := range updates {
		pk, err := cryptocodec.FromCmtProtoPublicKey(update.PubKey)
		s.Require().NoError(err)

		valPk, err := validator2.ConsPubKey()
		s.Require().NoError(err)

		if pk.Equals(valPk) {
			foundUpdate = true
			s.Require().Equal(int64(0), update.Power, "jailed validator should have zero power")
		}
	}
	s.Require().True(foundUpdate, "update for jailed validator not found")

	// Check validator 2 is unbonded
	updatedValidator2, found := s.tk.Keeper.GetValidator(s.tk.Ctx, addr2)
	s.Require().True(found)
	s.Require().False(updatedValidator2.Bonded, "jailed validator should not be bonded")

	// Validator 1 should still be bonded
	updatedValidator1, found := s.tk.Keeper.GetValidator(s.tk.Ctx, addr1)
	s.Require().True(found)
	s.Require().True(updatedValidator1.Bonded, "non-jailed validator should remain bonded")

	// Last validator power should be removed for jailed validator
	_, found = s.tk.Keeper.GetLastValidatorPower(s.tk.Ctx, addr2)
	s.Require().False(found, "last power should be removed for jailed validator")
}

// Test_ApplyAndReturnValidatorSetUpdates_MaxValidators tests the max validators limit
func (s *ValidatorSetTestSuite) Test_ApplyAndReturnValidatorSetUpdates_MaxValidators() {
	// Set test parameters with only 2 max validators
	params := types.Params{
		MaxValidators:    2,
		MaxLeverageRatio: math.LegacyNewDec(10),
		MinVotingPower:   1,
		WithdrawalLimit:  10,
	}
	err := s.tk.Keeper.SetParams(s.tk.Ctx, params)
	s.Require().NoError(err)

	// Register 3 validators with decreasing power
	_, addr1, _ := s.registerValidator(math.NewUint(5000000000), math.ZeroUint(), false) // 5 MITO, power = 5
	_, addr2, _ := s.registerValidator(math.NewUint(3000000000), math.ZeroUint(), false) // 3 MITO, power = 3
	_, addr3, _ := s.registerValidator(math.NewUint(2000000000), math.ZeroUint(), false) // 2 MITO, power = 2

	// Apply validator set updates
	updates, err := s.tk.Keeper.ApplyAndReturnValidatorSetUpdates(s.tk.Ctx)
	s.Require().NoError(err)

	// Check number of updates (should be 2 validators - the two highest powered ones)
	s.Require().Equal(2, len(updates))

	// Check validator 1 and 2 should be bonded
	updatedValidator1, found := s.tk.Keeper.GetValidator(s.tk.Ctx, addr1)
	s.Require().True(found)
	s.Require().True(updatedValidator1.Bonded)

	updatedValidator2, found := s.tk.Keeper.GetValidator(s.tk.Ctx, addr2)
	s.Require().True(found)
	s.Require().True(updatedValidator2.Bonded)

	// Validator 3 should not be bonded as it's not in the top MaxValidators
	updatedValidator3, found := s.tk.Keeper.GetValidator(s.tk.Ctx, addr3)
	s.Require().True(found)
	s.Require().False(updatedValidator3.Bonded)

	// Last validator powers should reflect this
	power1, found := s.tk.Keeper.GetLastValidatorPower(s.tk.Ctx, addr1)
	s.Require().True(found)
	s.Require().Equal(int64(5), power1)

	power2, found := s.tk.Keeper.GetLastValidatorPower(s.tk.Ctx, addr2)
	s.Require().True(found)
	s.Require().Equal(int64(3), power2)

	_, found = s.tk.Keeper.GetLastValidatorPower(s.tk.Ctx, addr3)
	s.Require().False(found, "validator 3 should not have last power as it's not in top MaxValidators")
}
