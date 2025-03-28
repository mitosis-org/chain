package keeper_test

import (
	"testing"

	"cosmossdk.io/math"
	"github.com/mitosis-org/chain/x/evmvalidator/keeper"
	"github.com/mitosis-org/chain/x/evmvalidator/testutil"
	"github.com/mitosis-org/chain/x/evmvalidator/types"
	"github.com/stretchr/testify/suite"
)

// MsgServerTestSuite is a test suite to be used with msg server tests
type MsgServerTestSuite struct {
	suite.Suite
	tk testutil.TestKeeper
}

// SetupTest initializes the test suite
func (s *MsgServerTestSuite) SetupTest() {
	s.tk = testutil.NewTestKeeper(&s.Suite)
}

// TestMsgServerTestSuite runs the msg server test suite
func TestMsgServerTestSuite(t *testing.T) {
	suite.Run(t, new(MsgServerTestSuite))
}

// Test_UpdateParams tests the UpdateParams message handler
func (s *MsgServerTestSuite) Test_UpdateParams() {
	// Set up initial params
	initialParams := s.tk.SetupDefaultTestParams()

	// Create msg server
	msgServer := keeper.NewMsgServerImpl(s.tk.Keeper)

	// Case 1: Test with valid authority
	newParams := initialParams
	newParams.MaxValidators = 200
	newParams.MaxLeverageRatio = math.LegacyNewDec(15)
	newParams.MinVotingPower = 2
	newParams.WithdrawalLimit = 20

	// Create the message with the correct authority
	msg := &types.MsgUpdateParams{
		Authority: "evmgov", // This should match the authority set in the keeper
		Params:    newParams,
	}

	// Test the message handler
	resp, err := msgServer.UpdateParams(s.tk.Ctx, msg)
	s.Require().NoError(err)
	s.Require().NotNil(resp)

	// Verify the params were updated
	updatedParams := s.tk.Keeper.GetParams(s.tk.Ctx)
	s.Require().Equal(newParams, updatedParams)

	// Case 2: Test with invalid authority
	invalidMsg := &types.MsgUpdateParams{
		Authority: "invalid-authority",
		Params:    newParams,
	}

	// Test the message handler with invalid authority
	_, err = msgServer.UpdateParams(s.tk.Ctx, invalidMsg)
	s.Require().Error(err)
	s.Require().Contains(err.Error(), "invalid authority")

	// Case 3: Test with invalid params
	invalidParamsMsg := &types.MsgUpdateParams{
		Authority: "evmgov",
		Params: types.Params{
			MaxValidators:    100,
			MaxLeverageRatio: math.LegacyNewDec(-1), // Invalid negative leverage ratio
			MinVotingPower:   1,
			WithdrawalLimit:  10,
		},
	}

	// Test the message handler with invalid params
	_, err = msgServer.UpdateParams(s.tk.Ctx, invalidParamsMsg)
	s.Require().Error(err)
}

// Test_UpdateParams_UpdatesValidatorStates_MaxLeverageRatio tests that updating parameters
// with changes to MaxLeverageRatio triggers validator state updates
func (s *MsgServerTestSuite) Test_UpdateParams_UpdatesValidatorStates_MaxLeverageRatio() {
	// Set up initial params
	initialParams := types.Params{
		MaxValidators:    100,
		MaxLeverageRatio: math.LegacyNewDec(10), // 10x leverage
		MinVotingPower:   10,
		WithdrawalLimit:  10,
	}
	err := s.tk.Keeper.SetParams(s.tk.Ctx, initialParams)
	s.Require().NoError(err)

	// Register a validator with exactly enough voting power to meet the initial requirement
	initialValidator := s.tk.RegisterTestValidator(math.NewUint(1000000000), math.NewUint(1000000000*10), false) // 1 MITO collateral, 10 MITO extra voting power, 10 voting power

	s.Require().False(initialValidator.Jailed, "Validator should not be jailed initially")
	s.Require().Equal(int64(10), initialValidator.VotingPower, "Initial voting power should be 10")

	// Create msg server
	msgServer := keeper.NewMsgServerImpl(s.tk.Keeper)

	// Create message with changes to MaxLeverageRatio
	newParams := initialParams
	newParams.MaxLeverageRatio = math.LegacyNewDec(5) // Decrease from 10 to 5

	msg := &types.MsgUpdateParams{
		Authority: "evmgov",
		Params:    newParams,
	}

	// Test the message handler
	resp, err := msgServer.UpdateParams(s.tk.Ctx, msg)
	s.Require().NoError(err)
	s.Require().NotNil(resp)

	// Verify the validator state was updated
	updatedValidator, found := s.tk.Keeper.GetValidator(s.tk.Ctx, initialValidator.Addr)
	s.Require().True(found)
	expectedValidator := initialValidator
	expectedValidator.VotingPower = 5
	expectedValidator.Jailed = true
	s.Require().Equal(expectedValidator, updatedValidator)
}

// Test_UpdateParams_UpdatesValidatorStates_MinVotingPower tests that updating parameters
// with changes to MinVotingPower triggers validator state updates
func (s *MsgServerTestSuite) Test_UpdateParams_UpdatesValidatorStates_MinVotingPower() {
	// Set up initial params with MinVotingPower = 1
	initialParams := s.tk.SetupDefaultTestParams()

	// Register a validator with exactly enough voting power to meet the initial requirement
	initialValidator := s.tk.RegisterTestValidator(math.NewUint(1000000000), math.ZeroUint(), false) // 1 MITO collateral

	// Verify initial validator state
	s.Require().False(initialValidator.Jailed, "Validator should not be jailed initially")
	s.Require().Equal(int64(1), initialValidator.VotingPower, "Initial voting power should be 10")

	// Create msg server
	msgServer := keeper.NewMsgServerImpl(s.tk.Keeper)

	// Create message with changes to MinVotingPower (increase requirement above validator's power)
	newParams := initialParams
	newParams.MinVotingPower = 2 // Increase from 1 to 2

	msg := &types.MsgUpdateParams{
		Authority: "evmgov",
		Params:    newParams,
	}

	// Test the message handler
	_, err := msgServer.UpdateParams(s.tk.Ctx, msg)
	s.Require().NoError(err)

	// Verify the validator state was updated
	updatedValidator, found := s.tk.Keeper.GetValidator(s.tk.Ctx, initialValidator.Addr)
	s.Require().True(found)
	expectedValidator := initialValidator
	expectedValidator.Jailed = true
	s.Require().Equal(expectedValidator, updatedValidator)
}

// Test_UpdateValidatorEntrypointContractAddr tests the UpdateValidatorEntrypointContractAddr message handler
func (s *MsgServerTestSuite) Test_UpdateValidatorEntrypointContractAddr() {
	// Set up initial params
	s.tk.SetupDefaultTestParams()

	// Create msg server
	msgServer := keeper.NewMsgServerImpl(s.tk.Keeper)

	// Create a sample EthAddress
	_, _, ethAddr := testutil.GenerateSecp256k1Key()

	// Case 1: Test with valid authority
	msg := &types.MsgUpdateValidatorEntrypointContractAddr{
		Authority: "evmgov",
		Addr:      ethAddr,
	}

	// Test the message handler
	resp, err := msgServer.UpdateValidatorEntrypointContractAddr(s.tk.Ctx, msg)
	s.Require().NoError(err)
	s.Require().NotNil(resp)

	// Verify the address was updated
	updatedAddr := s.tk.Keeper.GetValidatorEntrypointContractAddr(s.tk.Ctx)
	s.Require().Equal(ethAddr, updatedAddr)

	// Case 2: Test with invalid authority
	invalidMsg := &types.MsgUpdateValidatorEntrypointContractAddr{
		Authority: "invalid-authority",
		Addr:      ethAddr,
	}

	// Test the message handler with invalid authority
	_, err = msgServer.UpdateValidatorEntrypointContractAddr(s.tk.Ctx, invalidMsg)
	s.Require().Error(err)
	s.Require().Contains(err.Error(), "invalid authority")
}
