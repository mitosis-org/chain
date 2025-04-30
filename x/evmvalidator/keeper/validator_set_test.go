package keeper_test

import (
	"context"
	"testing"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
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
	s.tk = testutil.NewTestKeeper(&s.Suite)
}

// TestValidatorSetTestSuite runs the validator set test suite
func TestValidatorSetTestSuite(t *testing.T) {
	suite.Run(t, new(ValidatorSetTestSuite))
}

func (s *ValidatorSetTestSuite) Test_ApplyAndReturnValidatorSetUpdates_NewValidators() {
	// Set test parameters
	s.tk.SetupDefaultTestParams()

	// Register validators
	validator1 := s.tk.RegisterTestValidator(math.NewUint(5000000000), math.ZeroUint(), false) // 5 MITO, power = 5
	validator2 := s.tk.RegisterTestValidator(math.NewUint(3000000000), math.ZeroUint(), false) // 3 MITO, power = 3
	validator3 := s.tk.RegisterTestValidator(math.NewUint(2000000000), math.ZeroUint(), false) // 2 MITO, power = 2

	// Mock slashing keeper hooks
	bondedValidators := make(map[string]bool)
	s.tk.MockSlash.AfterValidatorBondedFn = func(ctx context.Context, consAddr sdk.ConsAddress) error {
		bondedValidators[consAddr.String()] = true
		return nil
	}

	// Apply validator set updates
	updates, err := s.tk.Keeper.ApplyAndReturnValidatorSetUpdates(s.tk.Ctx)
	s.Require().NoError(err)
	s.Require().Equal(3, len(updates)) // should be 3 validators
	s.Require().Contains(updates, validator1.MustABCIValidatorUpdate())
	s.Require().Contains(updates, validator2.MustABCIValidatorUpdate())
	s.Require().Contains(updates, validator3.MustABCIValidatorUpdate())

	// Check last validator powers were updated
	power1, found := s.tk.Keeper.GetLastValidatorPower(s.tk.Ctx, validator1.Addr)
	s.Require().True(found)
	s.Require().Equal(int64(5), power1)

	power2, found := s.tk.Keeper.GetLastValidatorPower(s.tk.Ctx, validator2.Addr)
	s.Require().True(found)
	s.Require().Equal(int64(3), power2)

	power3, found := s.tk.Keeper.GetLastValidatorPower(s.tk.Ctx, validator3.Addr)
	s.Require().True(found)
	s.Require().Equal(int64(2), power3)

	// Check validators are bonded
	updatedValidator1, found := s.tk.Keeper.GetValidator(s.tk.Ctx, validator1.Addr)
	s.Require().True(found)
	expectedValidator1 := validator1
	expectedValidator1.Bonded = true
	s.Require().Equal(expectedValidator1, updatedValidator1)

	updatedValidator2, found := s.tk.Keeper.GetValidator(s.tk.Ctx, validator2.Addr)
	s.Require().True(found)
	expectedValidator2 := validator2
	expectedValidator2.Bonded = true
	s.Require().Equal(expectedValidator2, updatedValidator2)

	updatedValidator3, found := s.tk.Keeper.GetValidator(s.tk.Ctx, validator3.Addr)
	s.Require().True(found)
	expectedValidator3 := validator3
	expectedValidator3.Bonded = true
	s.Require().Equal(expectedValidator3, updatedValidator3)

	// Check slashing keeper was called for each validator
	s.Require().Equal(3, len(bondedValidators))
	s.Require().True(bondedValidators[validator1.MustConsAddr().String()])
	s.Require().True(bondedValidators[validator2.MustConsAddr().String()])
	s.Require().True(bondedValidators[validator3.MustConsAddr().String()])
}

// Test_ApplyAndReturnValidatorSetUpdates_PowerChange tests that no updates are returned when there are no changes
func (s *ValidatorSetTestSuite) Test_ApplyAndReturnValidatorSetUpdates_PowerChange() {
	// Set test parameters
	s.tk.SetupDefaultTestParams()

	// Register validators
	validator1 := s.tk.RegisterTestValidator(math.NewUint(5000000000), math.ZeroUint(), false) // 5 MITO, power = 5
	validator2 := s.tk.RegisterTestValidator(math.NewUint(3000000000), math.ZeroUint(), false) // 3 MITO, power = 3

	// Initial update
	initialUpdates, err := s.tk.Keeper.ApplyAndReturnValidatorSetUpdates(s.tk.Ctx)
	s.Require().NoError(err)
	s.Require().Equal(2, len(initialUpdates))
	s.Require().Contains(initialUpdates, validator1.MustABCIValidatorUpdate())
	s.Require().Contains(initialUpdates, validator2.MustABCIValidatorUpdate())

	// Call again with no changes
	updates, err := s.tk.Keeper.ApplyAndReturnValidatorSetUpdates(s.tk.Ctx)
	s.Require().NoError(err)
	s.Require().Equal(0, len(updates), "should have no updates when there are no changes")

	// Get validators again
	var found bool
	validator1, found = s.tk.Keeper.GetValidator(s.tk.Ctx, validator1.Addr)
	s.Require().True(found)
	validator2, found = s.tk.Keeper.GetValidator(s.tk.Ctx, validator2.Addr)
	s.Require().True(found)

	// Deposit more collateral to change power
	s.tk.Keeper.DepositCollateral(s.tk.Ctx, &validator1, validator1.Addr, math.NewUint(2000000000))

	// Apply updates with changes
	updates, err = s.tk.Keeper.ApplyAndReturnValidatorSetUpdates(s.tk.Ctx)
	s.Require().NoError(err)
	s.Require().Equal(1, len(updates), "should have one update for the changed validator")
	s.Require().Contains(updates, validator1.MustABCIValidatorUpdate())

	// Check last validator power was updated
	power1, found := s.tk.Keeper.GetLastValidatorPower(s.tk.Ctx, validator1.Addr)
	s.Require().True(found)
	s.Require().Equal(int64(7), power1)

	// Validator 2 power should remain unchanged
	power2, found := s.tk.Keeper.GetLastValidatorPower(s.tk.Ctx, validator2.Addr)
	s.Require().True(found)
	s.Require().Equal(int64(3), power2)
}

// Test_ApplyAndReturnValidatorSetUpdates_JailedValidator tests validators are excluded when jailed
func (s *ValidatorSetTestSuite) Test_ApplyAndReturnValidatorSetUpdates_JailedValidator() {
	// Set test parameters
	s.tk.SetupDefaultTestParams()

	// Register validators - one normal, one jailed
	validator1 := s.tk.RegisterTestValidator(math.NewUint(5000000000), math.ZeroUint(), false) // 5 MITO, power = 5
	validator2 := s.tk.RegisterTestValidator(math.NewUint(3000000000), math.ZeroUint(), false) // 3 MITO, power = 3

	// Initial update
	initialUpdates, err := s.tk.Keeper.ApplyAndReturnValidatorSetUpdates(s.tk.Ctx)
	s.Require().NoError(err)
	s.Require().Equal(2, len(initialUpdates))
	s.Require().Contains(initialUpdates, validator1.MustABCIValidatorUpdate())
	s.Require().Contains(initialUpdates, validator2.MustABCIValidatorUpdate())

	// Get validators again
	var found bool
	validator1, found = s.tk.Keeper.GetValidator(s.tk.Ctx, validator1.Addr)
	s.Require().True(found)
	validator2, found = s.tk.Keeper.GetValidator(s.tk.Ctx, validator2.Addr)
	s.Require().True(found)

	// Jail the second validator
	s.tk.Keeper.Jail_(s.tk.Ctx, &validator2, "testing jail")

	// Apply updates
	updates, err := s.tk.Keeper.ApplyAndReturnValidatorSetUpdates(s.tk.Ctx)
	s.Require().NoError(err)

	// Should see one update - the jailed validator with zero power
	s.Require().Equal(1, len(updates), "should have one update for the jailed validator")
	s.Require().Contains(updates, validator2.MustABCIValidatorUpdateForUnbonding())

	// Validator 1 should still be bonded
	updatedValidator1, found := s.tk.Keeper.GetValidator(s.tk.Ctx, validator1.Addr)
	s.Require().True(found)
	s.Require().Equal(validator1, updatedValidator1)

	// Check validator 2 is unbonded
	updatedValidator2, found := s.tk.Keeper.GetValidator(s.tk.Ctx, validator2.Addr)
	s.Require().True(found)
	expectedValidator2 := validator2
	expectedValidator2.Bonded = false
	s.Require().Equal(expectedValidator2, updatedValidator2)

	// Last validator power should be kept for non-jailed validator
	power1, found := s.tk.Keeper.GetLastValidatorPower(s.tk.Ctx, validator1.Addr)
	s.Require().True(found)
	s.Require().Equal(int64(5), power1)

	// Last validator power should be removed for jailed validator
	_, found = s.tk.Keeper.GetLastValidatorPower(s.tk.Ctx, validator2.Addr)
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
	validator1 := s.tk.RegisterTestValidator(math.NewUint(5000000000), math.ZeroUint(), false) // 5 MITO, power = 5
	validator2 := s.tk.RegisterTestValidator(math.NewUint(3000000000), math.ZeroUint(), false) // 3 MITO, power = 3
	validator3 := s.tk.RegisterTestValidator(math.NewUint(2000000000), math.ZeroUint(), false) // 2 MITO, power = 2
	initialValidator1 := validator1
	initialValidator2 := validator2
	initialValidator3 := validator3
	valAddr1 := validator1.Addr
	valAddr2 := validator2.Addr
	valAddr3 := validator3.Addr

	// Apply validator set updates
	updates, err := s.tk.Keeper.ApplyAndReturnValidatorSetUpdates(s.tk.Ctx)
	s.Require().NoError(err)

	// Check number of updates (should be 2 validators - the two highest powered ones)
	s.Require().Equal(2, len(updates))
	s.Require().Contains(updates, validator1.MustABCIValidatorUpdate())
	s.Require().Contains(updates, validator2.MustABCIValidatorUpdate())

	var found bool

	// Check validator 1 and 2 should be bonded
	validator1, found = s.tk.Keeper.GetValidator(s.tk.Ctx, valAddr1)
	s.Require().True(found)
	expectedValidator1 := initialValidator1
	expectedValidator1.Bonded = true
	s.Require().Equal(expectedValidator1, validator1)

	validator2, found = s.tk.Keeper.GetValidator(s.tk.Ctx, valAddr2)
	s.Require().True(found)
	expectedValidator2 := initialValidator2
	expectedValidator2.Bonded = true
	s.Require().Equal(expectedValidator2, validator2)

	// Validator 3 should not be bonded as it's not in the top MaxValidators
	validator3, found = s.tk.Keeper.GetValidator(s.tk.Ctx, valAddr3)
	s.Require().True(found)
	s.Require().Equal(initialValidator3, validator3)

	// Last validator powers should reflect this
	power1, found := s.tk.Keeper.GetLastValidatorPower(s.tk.Ctx, valAddr1)
	s.Require().True(found)
	s.Require().Equal(int64(5), power1)

	power2, found := s.tk.Keeper.GetLastValidatorPower(s.tk.Ctx, valAddr2)
	s.Require().True(found)
	s.Require().Equal(int64(3), power2)

	_, found = s.tk.Keeper.GetLastValidatorPower(s.tk.Ctx, valAddr3)
	s.Require().False(found, "validator 3 should not have last power as it's not in top MaxValidators")

	/////////////////////
	// Max Validators Change due to power change
	/////////////////////

	s.tk.Keeper.UpdateExtraVotingPower(s.tk.Ctx, &validator3, math.NewUint(2000000000)) // power = 4

	updates, err = s.tk.Keeper.ApplyAndReturnValidatorSetUpdates(s.tk.Ctx)
	s.Require().NoError(err)
	s.Require().Equal(2, len(updates))
	s.Require().Contains(updates, validator2.MustABCIValidatorUpdateForUnbonding())
	s.Require().Contains(updates, validator3.MustABCIValidatorUpdate())

	expectedValidator2 = validator2
	expectedValidator2.Bonded = false
	validator2, found = s.tk.Keeper.GetValidator(s.tk.Ctx, valAddr2)
	s.Require().True(found)
	s.Require().Equal(expectedValidator2, validator2)

	expectedValidator3 := validator3
	expectedValidator3.Bonded = true
	validator3, found = s.tk.Keeper.GetValidator(s.tk.Ctx, valAddr3)
	s.Require().True(found)
	s.Require().Equal(expectedValidator3, validator3)

	// Last validator power should reflect this
	power1, found = s.tk.Keeper.GetLastValidatorPower(s.tk.Ctx, valAddr1)
	s.Require().True(found)
	s.Require().Equal(int64(5), power1)

	_, found = s.tk.Keeper.GetLastValidatorPower(s.tk.Ctx, valAddr2)
	s.Require().False(found, "validator 2 should not have last power as it's not in top MaxValidators")

	power3, found := s.tk.Keeper.GetLastValidatorPower(s.tk.Ctx, valAddr3)
	s.Require().True(found)
	s.Require().Equal(int64(4), power3)
}
