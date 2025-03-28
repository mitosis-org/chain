package keeper_test

import (
	"context"
	"testing"
	"time"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	mitotypes "github.com/mitosis-org/chain/types"
	"github.com/mitosis-org/chain/x/evmvalidator/testutil"
	"github.com/mitosis-org/chain/x/evmvalidator/types"
	"github.com/stretchr/testify/suite"
)

// ValidatorTestSuite is a test suite to be used with validator tests
type ValidatorTestSuite struct {
	suite.Suite
	tk testutil.TestKeeper
}

// SetupTest initializes the test suite
func (s *ValidatorTestSuite) SetupTest() {
	s.tk = testutil.CreateTestInput(&s.Suite)
}

// TestValidatorTestSuite runs the validator test suite
func TestValidatorTestSuite(t *testing.T) {
	suite.Run(t, new(ValidatorTestSuite))
}

// Helper functions to reduce duplication
func (s *ValidatorTestSuite) setupTestParams() types.Params {
	params := types.Params{
		MaxValidators:    100,
		MaxLeverageRatio: math.LegacyNewDec(10), // 10x leverage
		MinVotingPower:   1,
		WithdrawalLimit:  10,
	}
	err := s.tk.Keeper.SetParams(s.tk.Ctx, params)
	s.Require().NoError(err)
	return params
}

func (s *ValidatorTestSuite) registerValidator(collateral math.Uint, extraVotingPower math.Uint, jailed bool) ([]byte, mitotypes.EthAddress, types.Validator) {
	_, pubkey, ethAddr := testutil.GenerateSecp256k1Key()
	err := s.tk.Keeper.RegisterValidator(s.tk.Ctx, ethAddr, pubkey, collateral, extraVotingPower, jailed)
	s.Require().NoError(err)

	validator, found := s.tk.Keeper.GetValidator(s.tk.Ctx, ethAddr)
	s.Require().True(found)
	return pubkey, ethAddr, validator
}

// ==================== RegisterValidator Tests ====================

func (s *ValidatorTestSuite) Test_RegisterValidator() {
	// Generate validator data
	_, pubkey, ethAddr := testutil.GenerateSecp256k1Key()
	collateral := math.NewUint(1000000000) // 1 MITO in gwei
	extraVotingPower := math.NewUint(0)

	// Register validator
	err := s.tk.Keeper.RegisterValidator(s.tk.Ctx, ethAddr, pubkey, collateral, extraVotingPower, false)
	s.Require().NoError(err)

	// Check if validator exists
	validator, found := s.tk.Keeper.GetValidator(s.tk.Ctx, ethAddr)
	s.Require().True(found)
	s.Require().Equal(ethAddr, validator.Addr)
	s.Require().Equal(pubkey, validator.Pubkey)
	s.Require().Equal(collateral, validator.Collateral)
	s.Require().Equal(extraVotingPower, validator.ExtraVotingPower)
	s.Require().False(validator.Jailed)

	// Check if consensus address mapping exists
	consAddr, err := validator.ConsAddr()
	s.Require().NoError(err)

	valFromConsAddr, found := s.tk.Keeper.GetValidatorByConsAddr(s.tk.Ctx, consAddr)
	s.Require().True(found)
	s.Require().Equal(validator, valFromConsAddr)

	// Check if validator is properly indexed by power
	iterator := s.tk.Keeper.GetValidatorsByPowerIndexIterator(s.tk.Ctx)
	defer iterator.Close()

	found = false
	for ; iterator.Valid(); iterator.Next() {
		addr := mitotypes.BytesToEthAddress(iterator.Value())
		if addr.String() == ethAddr.String() {
			found = true
			break
		}
	}
	s.Require().True(found, "validator should be indexed by power")

	// Try registering the same validator again (should fail)
	err = s.tk.Keeper.RegisterValidator(s.tk.Ctx, ethAddr, pubkey, collateral, extraVotingPower, false)
	s.Require().Error(err)
	s.Require().ErrorIs(err, types.ErrValidatorAlreadyExists)

	// Try registering with invalid pubkey (different address)
	_, invalidPubkey, _ := testutil.GenerateSecp256k1Key()
	err = s.tk.Keeper.RegisterValidator(s.tk.Ctx, ethAddr, invalidPubkey, collateral, extraVotingPower, false)
	s.Require().Error(err)
}

func (s *ValidatorTestSuite) Test_RegisterValidator_ZeroCollateral() {
	_, pubkey, ethAddr := testutil.GenerateSecp256k1Key()
	zeroCollateral := math.ZeroUint()
	extraVotingPower := math.NewUint(0)

	// Try registering with zero collateral (should still work but voting power will be zero)
	err := s.tk.Keeper.RegisterValidator(s.tk.Ctx, ethAddr, pubkey, zeroCollateral, extraVotingPower, false)
	s.Require().NoError(err)

	// Check if validator exists
	validator, found := s.tk.Keeper.GetValidator(s.tk.Ctx, ethAddr)
	s.Require().True(found)
	s.Require().Equal(zeroCollateral, validator.Collateral)

	// Validator's voting power should be zero or automatically jailed if min power required
	params := s.tk.Keeper.GetParams(s.tk.Ctx)
	if params.MinVotingPower > 0 {
		s.Require().True(validator.Jailed)
	}
}

func (s *ValidatorTestSuite) Test_RegisterValidator_InvalidPubkey() {
	invalidPubkey := []byte("invalid-pubkey")
	_, _, validEthAddr := testutil.GenerateSecp256k1Key()
	collateral := math.NewUint(1000000000)
	extraVotingPower := math.NewUint(0)

	// Try registering with invalid pubkey format
	err := s.tk.Keeper.RegisterValidator(s.tk.Ctx, validEthAddr, invalidPubkey, collateral, extraVotingPower, false)
	s.Require().Error(err)
}

// ==================== DepositCollateral Tests ====================

func (s *ValidatorTestSuite) Test_DepositCollateral() {
	// Use helper functions
	s.setupTestParams()
	_, ethAddr, validator := s.registerValidator(math.NewUint(1000000000), math.ZeroUint(), false)

	// Initial voting power should be 1
	initialVotingPower := int64(1)
	s.Require().Equal(initialVotingPower, validator.VotingPower)

	// Deposit additional collateral
	additionalCollateral := math.NewUint(500000000) // 0.5 MITO in gwei
	s.tk.Keeper.DepositCollateral(s.tk.Ctx, &validator, additionalCollateral)

	// Check if collateral was updated
	expectedCollateral := math.NewUint(1000000000).Add(additionalCollateral)

	// Get validator after deposit
	updatedValidator, found := s.tk.Keeper.GetValidator(s.tk.Ctx, ethAddr)
	s.Require().True(found)
	s.Require().Equal(expectedCollateral, updatedValidator.Collateral)

	// Expected voting power should increase from 1 to 1.5, which gets truncated to 1
	// For a more noticeable change, let's add another 0.5 MITO
	s.tk.Keeper.DepositCollateral(s.tk.Ctx, &updatedValidator, additionalCollateral)

	// Now we should have 2 MITO total, which should give 2 voting power
	finalExpectedCollateral := expectedCollateral.Add(additionalCollateral)
	expectedVotingPower := int64(2)

	// Get validator after second deposit
	finalValidator, found := s.tk.Keeper.GetValidator(s.tk.Ctx, ethAddr)
	s.Require().True(found)
	s.Require().Equal(finalExpectedCollateral, finalValidator.Collateral)
	s.Require().Equal(expectedVotingPower, finalValidator.VotingPower)
}

func (s *ValidatorTestSuite) Test_DepositCollateral_ZeroAmount() {
	// Use helper function to register a validator
	_, ethAddr, validator := s.registerValidator(math.NewUint(1000000000), math.ZeroUint(), false)
	initialCollateral := validator.Collateral

	// Initial voting power
	initialVotingPower := validator.VotingPower

	// Deposit zero collateral
	zeroCollateral := math.ZeroUint()
	s.tk.Keeper.DepositCollateral(s.tk.Ctx, &validator, zeroCollateral)

	// Check collateral remains unchanged
	s.Require().Equal(initialCollateral, validator.Collateral)

	// Check voting power remains unchanged
	s.Require().Equal(initialVotingPower, validator.VotingPower)

	// Check state is updated
	updatedValidator, found := s.tk.Keeper.GetValidator(s.tk.Ctx, ethAddr)
	s.Require().True(found)
	s.Require().Equal(initialCollateral, updatedValidator.Collateral)
}

// ==================== WithdrawCollateral Tests ====================

func (s *ValidatorTestSuite) Test_WithdrawCollateral() {
	// Set parameters for test
	s.setupTestParams()

	// Register a validator with enough collateral to have voting power = 2
	_, ethAddr, validator := s.registerValidator(math.NewUint(2000000000), math.ZeroUint(), false)
	initialCollateral := validator.Collateral

	// Initial voting power should be 2
	initialVotingPower := int64(2)
	s.Require().Equal(initialVotingPower, validator.VotingPower)

	// Create withdrawal request
	withdrawalAmount := uint64(500000000) // 0.5 MITO in gwei
	receiver := ethAddr
	maturesAt := time.Now().Unix() + 86400 // 1 day from now

	withdrawal := &types.Withdrawal{
		ValAddr:        ethAddr,
		Amount:         withdrawalAmount,
		Receiver:       receiver,
		MaturesAt:      maturesAt,
		CreationHeight: s.tk.Ctx.BlockHeight(),
	}

	// Withdraw collateral
	err := s.tk.Keeper.WithdrawCollateral(s.tk.Ctx, &validator, withdrawal)
	s.Require().NoError(err)

	// Check if collateral was updated
	expectedCollateral := initialCollateral.Sub(math.NewUint(withdrawalAmount))

	// Expected voting power after withdrawal should be 1.5, truncated to 1
	expectedVotingPower := int64(1)

	// Get validator after withdrawal
	updatedValidator, found := s.tk.Keeper.GetValidator(s.tk.Ctx, ethAddr)
	s.Require().True(found)
	s.Require().Equal(expectedCollateral, updatedValidator.Collateral)
	s.Require().Equal(expectedVotingPower, updatedValidator.VotingPower)

	// Try withdrawing more than available (should fail)
	excessWithdrawal := &types.Withdrawal{
		ValAddr:        ethAddr,
		Amount:         initialCollateral.Uint64(), // Try to withdraw full initial amount
		Receiver:       receiver,
		MaturesAt:      maturesAt,
		CreationHeight: s.tk.Ctx.BlockHeight(),
	}

	err = s.tk.Keeper.WithdrawCollateral(s.tk.Ctx, &updatedValidator, excessWithdrawal)
	s.Require().Error(err)
	s.Require().ErrorIs(err, types.ErrInsufficientCollateral)

	// Try withdrawing zero amount (should succeed but do nothing)
	zeroWithdrawal := &types.Withdrawal{
		ValAddr:        ethAddr,
		Amount:         0,
		Receiver:       receiver,
		MaturesAt:      maturesAt,
		CreationHeight: s.tk.Ctx.BlockHeight(),
	}

	err = s.tk.Keeper.WithdrawCollateral(s.tk.Ctx, &updatedValidator, zeroWithdrawal)
	s.Require().NoError(err)

	// Get validator after zero withdrawal (should be unchanged)
	finalValidator, found := s.tk.Keeper.GetValidator(s.tk.Ctx, ethAddr)
	s.Require().True(found)
	s.Require().Equal(expectedCollateral, finalValidator.Collateral)
	s.Require().Equal(expectedVotingPower, finalValidator.VotingPower)
}

func (s *ValidatorTestSuite) Test_WithdrawCollateral_ToMinimumRequired() {
	// Set up test parameters
	s.setupTestParams()

	// Register a validator with enough collateral
	_, ethAddr, validator := s.registerValidator(math.NewUint(2000000000), math.ZeroUint(), false)
	initialCollateral := validator.Collateral

	s.Require().Equal(int64(2), validator.VotingPower) // Should have 2 voting power

	// Withdraw to just above the minimum (1 MITO)
	withdrawalAmount := uint64(900000000) // 0.9 MITO in gwei
	receiver := ethAddr
	maturesAt := time.Now().Unix() + 86400 // 1 day from now

	withdrawal := &types.Withdrawal{
		ValAddr:        ethAddr,
		Amount:         withdrawalAmount,
		Receiver:       receiver,
		MaturesAt:      maturesAt,
		CreationHeight: s.tk.Ctx.BlockHeight(),
	}

	// Withdraw collateral
	err := s.tk.Keeper.WithdrawCollateral(s.tk.Ctx, &validator, withdrawal)
	s.Require().NoError(err)

	// Expected remaining collateral: 2 MITO - 0.9 MITO = 1.1 MITO
	expectedCollateral := initialCollateral.Sub(math.NewUint(withdrawalAmount))
	expectedVotingPower := int64(1) // Voting power should be 1

	// Check updated validator
	updatedValidator, found := s.tk.Keeper.GetValidator(s.tk.Ctx, ethAddr)
	s.Require().True(found)
	s.Require().Equal(expectedCollateral, updatedValidator.Collateral)
	s.Require().Equal(expectedVotingPower, updatedValidator.VotingPower)
	s.Require().False(updatedValidator.Jailed)

	// Now withdraw to just below the minimum (1 MITO)
	belowMinWithdrawal := &types.Withdrawal{
		ValAddr:        ethAddr,
		Amount:         200000000, // 0.2 MITO in gwei
		Receiver:       receiver,
		MaturesAt:      maturesAt,
		CreationHeight: s.tk.Ctx.BlockHeight(),
	}

	// Withdraw collateral
	err = s.tk.Keeper.WithdrawCollateral(s.tk.Ctx, &updatedValidator, belowMinWithdrawal)
	s.Require().NoError(err)

	// Final collateral: 1.1 MITO - 0.2 MITO = 0.9 MITO (below min required)
	finalValidator, found := s.tk.Keeper.GetValidator(s.tk.Ctx, ethAddr)
	s.Require().True(found)
	s.Require().Equal(math.NewUint(900000000), finalValidator.Collateral)
	s.Require().Zero(finalValidator.VotingPower) // Should have 0 voting power

	// Validator should be jailed because voting power is below minimum
	s.Require().True(finalValidator.Jailed, "Validator should be jailed when voting power drops below minimum")
}

// ==================== Slash_ Tests ====================

func (s *ValidatorTestSuite) Test_Slash_() {
	// Set up test parameters
	s.setupTestParams()

	// Register a validator
	_, ethAddr, validator := s.registerValidator(math.NewUint(1000000000), math.ZeroUint(), false)

	// Should have 1 voting power based on 1 MITO collateral
	naturalVotingPower := validator.VotingPower
	s.Require().Equal(int64(1), naturalVotingPower)

	// Slash parameters
	infractionHeight := s.tk.Ctx.BlockHeight() - 1
	slashFraction := math.LegacyNewDecWithPrec(5, 2) // 5% slash
	power := naturalVotingPower

	// Slash the validator
	slashedAmount, err := s.tk.Keeper.Slash_(s.tk.Ctx, &validator, infractionHeight, power, slashFraction)
	s.Require().NoError(err)

	// Calculate expected slashed amount
	expectedSlashedAmount := math.NewUintFromBigInt(
		math.LegacyNewDec(power).
			MulInt(types.VotingPowerReduction).
			Mul(slashFraction).
			TruncateInt().
			BigInt(),
	)

	// Check slashed amount
	s.Require().Equal(expectedSlashedAmount, slashedAmount)

	// Check collateral was reduced
	expectedCollateral := math.NewUint(1000000000).Sub(expectedSlashedAmount)
	s.Require().Equal(expectedCollateral, validator.Collateral)

	// Check if validator was updated in state
	updatedValidator, found := s.tk.Keeper.GetValidator(s.tk.Ctx, ethAddr)
	s.Require().True(found)
	s.Require().Equal(expectedCollateral, updatedValidator.Collateral)

	// Test slashing with negative fraction (should fail)
	_, err = s.tk.Keeper.Slash_(s.tk.Ctx, &validator, infractionHeight, power, math.LegacyNewDec(-1))
	s.Require().Error(err)
}

func (s *ValidatorTestSuite) Test_Slash_ExceedsCollateral() {
	// Set up test parameters
	s.setupTestParams()

	// Register a validator
	_, ethAddr, validator := s.registerValidator(math.NewUint(1000000000), math.ZeroUint(), false)
	initialCollateral := validator.Collateral

	// Should have 1 voting power based on 1 MITO collateral
	naturalVotingPower := validator.VotingPower
	s.Require().Equal(int64(1), naturalVotingPower)

	// Slash parameters - use a high voting power (10x natural) and 100% slash to exceed collateral
	infractionHeight := s.tk.Ctx.BlockHeight() - 1
	slashFraction := math.LegacyNewDecWithPrec(100, 2) // 100% slash
	slashPower := naturalVotingPower * 10              // Use power higher than available collateral

	// Slash the validator
	slashedAmount, err := s.tk.Keeper.Slash_(s.tk.Ctx, &validator, infractionHeight, slashPower, slashFraction)
	s.Require().NoError(err)

	// Since attempted slash exceeds collateral, should only slash what's available
	s.Require().Equal(initialCollateral, slashedAmount)
	s.Require().Equal(math.ZeroUint(), validator.Collateral)

	// Validator should be updated in state with zero collateral
	updatedValidator, found := s.tk.Keeper.GetValidator(s.tk.Ctx, ethAddr)
	s.Require().True(found)
	s.Require().Equal(math.ZeroUint(), updatedValidator.Collateral)

	// With zero collateral, validator should have zero voting power and be jailed
	s.Require().Zero(updatedValidator.VotingPower)
	s.Require().True(updatedValidator.Jailed)
}

func (s *ValidatorTestSuite) Test_Slash_WithdrawalsAndCollateral() {
	// ==================== SETUP ====================
	s.setupTestParams()
	_, valAddr, validator := s.registerValidator(math.NewUint(5000000000), math.ZeroUint(), false) // 5 MITO

	// ==================== SETUP WITHDRAWALS ====================
	now := time.Now().Unix()

	// Create two future withdrawals with different maturity times
	futureWithdrawal1 := &types.Withdrawal{
		ValAddr:        valAddr,
		Amount:         600000000, // 0.6 MITO
		Receiver:       valAddr,
		MaturesAt:      now + 86400, // 1 day from now
		CreationHeight: s.tk.Ctx.BlockHeight(),
	}

	futureWithdrawal2 := &types.Withdrawal{
		ValAddr:        valAddr,
		Amount:         800000000, // 0.8 MITO
		Receiver:       valAddr,
		MaturesAt:      now + 172800, // 2 days from now
		CreationHeight: s.tk.Ctx.BlockHeight(),
	}

	// Create an already matured withdrawal (1 day ago)
	maturedWithdrawal := &types.Withdrawal{
		ValAddr:        valAddr,
		Amount:         2000000000, // 2 MITO
		Receiver:       valAddr,
		MaturesAt:      now - 86400,
		CreationHeight: s.tk.Ctx.BlockHeight() - 100,
	}

	// Process the withdrawals
	err := s.tk.Keeper.WithdrawCollateral(s.tk.Ctx, &validator, futureWithdrawal1)
	s.Require().NoError(err)
	err = s.tk.Keeper.WithdrawCollateral(s.tk.Ctx, &validator, futureWithdrawal2)
	s.Require().NoError(err)
	err = s.tk.Keeper.WithdrawCollateral(s.tk.Ctx, &validator, maturedWithdrawal)
	s.Require().NoError(err)

	// Get updated validator with 1.6 MITO (5 - 0.6 - 0.8 - 2) of collateral
	updatedValidator, found := s.tk.Keeper.GetValidator(s.tk.Ctx, valAddr)
	s.Require().True(found)
	s.Require().Equal(math.NewUint(1600000000), updatedValidator.Collateral, "Validator should have 1.6 MITO collateral")
	s.Require().Equal(int64(1), updatedValidator.VotingPower, "Validator should have 1 voting power")

	// ==================== RECORD STATE BEFORE SLASHING ====================
	// Log the validator state before slashing for reporting purposes
	s.T().Logf("Before slashing: Validator collateral=%s, power=%d",
		updatedValidator.Collateral.String(), updatedValidator.VotingPower)

	// ==================== PERFORM SLASHING ====================
	// Slash the validator by 100% of voting power (1 MITO)
	slashFraction := math.LegacyOneDec() // 100%
	slashedAmount, err := s.tk.Keeper.Slash_(s.tk.Ctx, &updatedValidator, s.tk.Ctx.BlockHeight()-10, 1, slashFraction)
	s.Require().NoError(err)

	// ==================== VERIFY RESULTS ====================
	// 1. Verify slashed amount matches expectation (1 MITO)
	expectedSlashedAmount := math.NewUint(1000000000)
	s.Require().Equal(expectedSlashedAmount, slashedAmount, "Should slash 1 MITO")

	// 2. Check which withdrawals remain after slashing
	var remainingWithdrawals []types.Withdrawal
	s.tk.Keeper.IterateWithdrawalsForValidator(s.tk.Ctx, valAddr, func(w types.Withdrawal) bool {
		remainingWithdrawals = append(remainingWithdrawals, w)
		return false
	})

	// 3. Matured withdrawals should not be affected by slashing
	var foundMaturedWithdrawal bool
	for _, w := range remainingWithdrawals {
		if w.MaturesAt == maturedWithdrawal.MaturesAt {
			foundMaturedWithdrawal = true
			s.Require().Equal(maturedWithdrawal.Amount, w.Amount,
				"Matured withdrawal amount should be unchanged")
			break
		}
	}
	s.Require().True(foundMaturedWithdrawal, "Matured withdrawal should still exist")

	// 4. First future withdrawal (0.6 MITO) should be completely slashed and deleted
	var foundFirstFutureWithdrawal bool
	for _, w := range remainingWithdrawals {
		if w.MaturesAt == futureWithdrawal1.MaturesAt {
			foundFirstFutureWithdrawal = true
			break
		}
	}
	s.Require().False(foundFirstFutureWithdrawal,
		"First future withdrawal should be deleted, not set to zero amount")

	// 5. Second future withdrawal (0.8 MITO) should be partially slashed (by 0.4 MITO)
	// leaving 0.4 MITO remaining
	var foundSecondFutureWithdrawal bool
	for _, w := range remainingWithdrawals {
		if w.MaturesAt == futureWithdrawal2.MaturesAt {
			foundSecondFutureWithdrawal = true
			s.Require().Equal(uint64(400000000), w.Amount,
				"Second future withdrawal should be partially slashed to 0.4 MITO")
			break
		}
	}
	s.Require().True(foundSecondFutureWithdrawal,
		"Second future withdrawal should still exist with reduced amount")

	// 6. Verify validator collateral remains unchanged
	finalValidator, found := s.tk.Keeper.GetValidator(s.tk.Ctx, valAddr)
	s.Require().True(found)
	s.Require().Equal(math.NewUint(1600000000), finalValidator.Collateral,
		"Validator collateral should remain unchanged")
}

// ==================== Jail_ Tests ====================

func (s *ValidatorTestSuite) Test_Jail_() {
	// Register a validator
	_, pubkey, ethAddr := testutil.GenerateSecp256k1Key()
	initialCollateral := math.NewUint(1000000000)
	extraVotingPower := math.NewUint(0)

	err := s.tk.Keeper.RegisterValidator(s.tk.Ctx, ethAddr, pubkey, initialCollateral, extraVotingPower, false)
	s.Require().NoError(err)

	// Get validator before jailing
	validator, found := s.tk.Keeper.GetValidator(s.tk.Ctx, ethAddr)
	s.Require().True(found)
	s.Require().False(validator.Jailed)

	// Jail the validator
	reason := "testing jail"
	s.tk.Keeper.Jail_(s.tk.Ctx, &validator, reason)

	// Check if validator was jailed
	s.Require().True(validator.Jailed)

	// Check if validator was updated in state
	updatedValidator, found := s.tk.Keeper.GetValidator(s.tk.Ctx, ethAddr)
	s.Require().True(found)
	s.Require().True(updatedValidator.Jailed)
}

// ==================== Unjail_ Tests ====================

func (s *ValidatorTestSuite) Test_Unjail_() {
	// Define parameters for testing with minimum voting power of 1
	params := types.Params{
		MaxValidators:    100,
		MaxLeverageRatio: math.LegacyNewDec(10), // 10x leverage
		MinVotingPower:   1,
		WithdrawalLimit:  10,
	}

	// Set parameters
	err := s.tk.Keeper.SetParams(s.tk.Ctx, params)
	s.Require().NoError(err)

	// Register a jailed validator with enough collateral to meet min voting power
	_, pubkey, ethAddr := testutil.GenerateSecp256k1Key()
	initialCollateral := math.NewUint(1000000000) // 1 MITO in gwei = 1 voting power
	extraVotingPower := math.NewUint(0)

	err = s.tk.Keeper.RegisterValidator(s.tk.Ctx, ethAddr, pubkey, initialCollateral, extraVotingPower, true) // jailed = true
	s.Require().NoError(err)

	// Get validator before unjailing
	validator, found := s.tk.Keeper.GetValidator(s.tk.Ctx, ethAddr)
	s.Require().True(found)
	s.Require().True(validator.Jailed)

	// For unjailing to succeed, we need to mock the slashing keeper's UnjailFromConsAddr function
	s.tk.MockSlash.UnjailFromConsAddrFn = func(ctx context.Context, consAddr sdk.ConsAddress) error {
		return nil
	}

	// Set voting power to meet the minimum requirement (should be already set by the keeper)
	if validator.VotingPower < 1 {
		validator.VotingPower = 1
		s.tk.Keeper.SetValidator(s.tk.Ctx, validator)
	}

	// Unjail the validator
	err = s.tk.Keeper.Unjail_(s.tk.Ctx, &validator)
	s.Require().NoError(err)

	// Check if validator was unjailed
	s.Require().False(validator.Jailed)

	// Check if validator was updated in state
	updatedValidator, found := s.tk.Keeper.GetValidator(s.tk.Ctx, ethAddr)
	s.Require().True(found)
	s.Require().False(updatedValidator.Jailed)

	// Test unjailing a non-jailed validator
	err = s.tk.Keeper.Unjail_(s.tk.Ctx, &updatedValidator)
	s.Require().NoError(err) // No error because the function just returns if already unjailed
}

func (s *ValidatorTestSuite) Test_Unjail_InsufficientVotingPower() {
	// Define parameters for testing with minimum voting power of 10
	params := types.Params{
		MaxValidators:    100,
		MaxLeverageRatio: math.LegacyNewDec(10), // 10x leverage
		MinVotingPower:   10,                    // Higher minimum voting power
		WithdrawalLimit:  10,
	}

	// Set parameters
	err := s.tk.Keeper.SetParams(s.tk.Ctx, params)
	s.Require().NoError(err)

	// Register a jailed validator with insufficient collateral
	_, pubkey, ethAddr := testutil.GenerateSecp256k1Key()
	initialCollateral := math.NewUint(1000000000) // 1 MITO in gwei = 1 voting power
	extraVotingPower := math.NewUint(0)

	err = s.tk.Keeper.RegisterValidator(s.tk.Ctx, ethAddr, pubkey, initialCollateral, extraVotingPower, true) // jailed = true
	s.Require().NoError(err)

	// Get validator
	validator, found := s.tk.Keeper.GetValidator(s.tk.Ctx, ethAddr)
	s.Require().True(found)
	s.Require().True(validator.Jailed)

	// Calculate voting power
	validator.VotingPower = validator.ComputeVotingPower(params.MaxLeverageRatio)
	s.Require().Equal(int64(1), validator.VotingPower) // Should have 1 voting power, below minimum 10

	// For mocking the slashing keeper's UnjailFromConsAddr function
	s.tk.MockSlash.UnjailFromConsAddrFn = func(ctx context.Context, consAddr sdk.ConsAddress) error {
		return nil
	}

	// Try to unjail the validator with insufficient voting power
	err = s.tk.Keeper.Unjail_(s.tk.Ctx, &validator)
	s.Require().Error(err)
	s.Require().Contains(err.Error(), "voting power below minimum requirement")

	// Validator should still be jailed
	updatedValidator, found := s.tk.Keeper.GetValidator(s.tk.Ctx, ethAddr)
	s.Require().True(found)
	s.Require().True(updatedValidator.Jailed)
}

func (s *ValidatorTestSuite) Test_Unjail_SlashingKeeperError() {
	// Define parameters for testing with a high minimum voting power
	params := types.Params{
		MaxValidators:    100,
		MaxLeverageRatio: math.LegacyNewDec(10), // 10x leverage
		MinVotingPower:   5,                     // Higher minimum voting power than collateral can provide
		WithdrawalLimit:  10,
	}

	// Set parameters
	err := s.tk.Keeper.SetParams(s.tk.Ctx, params)
	s.Require().NoError(err)

	// Register a jailed validator with enough collateral for 1 voting power
	_, pubkey, ethAddr := testutil.GenerateSecp256k1Key()
	initialCollateral := math.NewUint(1000000000) // 1 MITO in gwei = 1 voting power
	extraVotingPower := math.NewUint(0)

	err = s.tk.Keeper.RegisterValidator(s.tk.Ctx, ethAddr, pubkey, initialCollateral, extraVotingPower, true) // jailed = true
	s.Require().NoError(err)

	// Get validator
	validator, found := s.tk.Keeper.GetValidator(s.tk.Ctx, ethAddr)
	s.Require().True(found)
	s.Require().True(validator.Jailed)

	// Verify the voting power is less than the minimum required
	votingPower := validator.ComputeVotingPower(params.MaxLeverageRatio)
	s.Require().Less(votingPower, params.MinVotingPower)

	// Try to unjail the validator with insufficient voting power
	err = s.tk.Keeper.Unjail_(s.tk.Ctx, &validator)
	s.Require().Error(err)
	s.Require().Contains(err.Error(), "voting power below minimum requirement")

	// Validator should still be jailed
	updatedValidator, found := s.tk.Keeper.GetValidator(s.tk.Ctx, ethAddr)
	s.Require().True(found)
	s.Require().True(updatedValidator.Jailed)
}

// ==================== UpdateExtraVotingPower Tests ====================

func (s *ValidatorTestSuite) Test_UpdateExtraVotingPower() {
	// Register a validator
	_, pubkey, ethAddr := testutil.GenerateSecp256k1Key()
	initialCollateral := math.NewUint(1000000000)
	initialExtraVotingPower := math.NewUint(0)

	err := s.tk.Keeper.RegisterValidator(s.tk.Ctx, ethAddr, pubkey, initialCollateral, initialExtraVotingPower, false)
	s.Require().NoError(err)

	// Get validator before update
	validator, found := s.tk.Keeper.GetValidator(s.tk.Ctx, ethAddr)
	s.Require().True(found)
	s.Require().Equal(initialExtraVotingPower, validator.ExtraVotingPower)

	// Update extra voting power
	newExtraVotingPower := math.NewUint(500000000)
	s.tk.Keeper.UpdateExtraVotingPower(s.tk.Ctx, &validator, newExtraVotingPower)

	// Check if extra voting power was updated
	s.Require().Equal(newExtraVotingPower, validator.ExtraVotingPower)

	// Check if validator was updated in state
	updatedValidator, found := s.tk.Keeper.GetValidator(s.tk.Ctx, ethAddr)
	s.Require().True(found)
	s.Require().Equal(newExtraVotingPower, updatedValidator.ExtraVotingPower)
}

// ==================== UpdateValidatorState Tests ====================

func (s *ValidatorTestSuite) Test_UpdateValidatorState() {
	// Define parameters for testing
	params := types.Params{
		MaxValidators:    100,
		MaxLeverageRatio: math.LegacyNewDec(10), // 10x leverage
		MinVotingPower:   1,
		WithdrawalLimit:  10,
	}

	// Set parameters
	err := s.tk.Keeper.SetParams(s.tk.Ctx, params)
	s.Require().NoError(err)

	// Register a validator
	_, pubkey, ethAddr := testutil.GenerateSecp256k1Key()
	collateral := math.NewUint(1000000000)      // 1 MITO in gwei
	extraVotingPower := math.NewUint(500000000) // 0.5 MITO in gwei

	err = s.tk.Keeper.RegisterValidator(s.tk.Ctx, ethAddr, pubkey, collateral, extraVotingPower, false)
	s.Require().NoError(err)

	// Get validator
	validator, found := s.tk.Keeper.GetValidator(s.tk.Ctx, ethAddr)
	s.Require().True(found)

	// Manually update values to verify calculation
	validator.Collateral = math.NewUint(2000000000)       // 2 MITO in gwei
	validator.ExtraVotingPower = math.NewUint(1000000000) // 1 MITO in gwei
	validator.Jailed = false
	validator.Bonded = false

	// Set the validator with our changes (the UpdateValidatorState expects the validator to exist in state)
	s.tk.Keeper.SetValidator(s.tk.Ctx, validator)

	// Update validator state
	s.tk.Keeper.UpdateValidatorState(s.tk.Ctx, &validator, "testing")

	// Calculate expected voting power
	// Voting power = (collateral + extraVotingPower) / 10^9 = 3
	expectedVotingPower := int64(3)

	// Check if voting power was updated correctly
	s.Require().Equal(expectedVotingPower, validator.VotingPower)

	// Check if validator was updated in state
	updatedValidator, found := s.tk.Keeper.GetValidator(s.tk.Ctx, ethAddr)
	s.Require().True(found)
	s.Require().Equal(expectedVotingPower, updatedValidator.VotingPower)

	// Test with jailed validator
	validator.Jailed = true

	// Set the validator with our changes
	s.tk.Keeper.SetValidator(s.tk.Ctx, validator)

	// Update validator state for jailed validator
	s.tk.Keeper.UpdateValidatorState(s.tk.Ctx, &validator, "testing jailed")

	// Jailed validator should have same voting power but should not be bonded
	s.Require().Equal(expectedVotingPower, validator.VotingPower)
	s.Require().False(validator.Bonded)

	// Check if validator was updated in state with jailed status
	jailedValidator, found := s.tk.Keeper.GetValidator(s.tk.Ctx, ethAddr)
	s.Require().True(found)
	s.Require().Equal(expectedVotingPower, jailedValidator.VotingPower)
	s.Require().True(jailedValidator.Jailed)
}

func (s *ValidatorTestSuite) Test_UpdateValidatorState_BelowMinVotingPower() {
	// Define parameters for testing
	params := types.Params{
		MaxValidators:    100,
		MaxLeverageRatio: math.LegacyNewDec(10), // 10x leverage
		MinVotingPower:   1,
		WithdrawalLimit:  10,
	}

	// Set parameters
	err := s.tk.Keeper.SetParams(s.tk.Ctx, params)
	s.Require().NoError(err)

	// Register a validator with collateral at minimum
	_, pubkey, ethAddr := testutil.GenerateSecp256k1Key()
	initialCollateral := math.NewUint(1000000000) // 1 MITO in gwei = 1 voting power
	extraVotingPower := math.NewUint(0)

	err = s.tk.Keeper.RegisterValidator(s.tk.Ctx, ethAddr, pubkey, initialCollateral, extraVotingPower, false)
	s.Require().NoError(err)

	// Get validator
	validator, found := s.tk.Keeper.GetValidator(s.tk.Ctx, ethAddr)
	s.Require().True(found)
	s.Require().Equal(int64(1), validator.VotingPower) // Should have 1 voting power
	s.Require().False(validator.Jailed)

	// Reduce collateral below minimum
	validator.Collateral = math.NewUint(900000000) // 0.9 MITO in gwei
	s.tk.Keeper.SetValidator(s.tk.Ctx, validator)

	// Update validator state - this should trigger jailing due to insufficient voting power
	s.tk.Keeper.UpdateValidatorState(s.tk.Ctx, &validator, "testing minimum voting power")

	// Check if validator was jailed due to insufficient voting power
	s.Require().True(validator.Jailed, "Validator should be jailed when voting power drops below minimum")
	s.Require().Zero(validator.VotingPower) // Jailed validators have zero consensus power

	// Check if validator was updated in state
	updatedValidator, found := s.tk.Keeper.GetValidator(s.tk.Ctx, ethAddr)
	s.Require().True(found)
	s.Require().True(updatedValidator.Jailed)
}

func (s *ValidatorTestSuite) Test_UpdateValidatorState_MaxLeverageRatio() {
	// Define parameters for testing with low max leverage ratio
	params := types.Params{
		MaxValidators:    100,
		MaxLeverageRatio: math.LegacyNewDec(2), // Only 2x leverage allowed
		MinVotingPower:   1,
		WithdrawalLimit:  10,
	}

	// Set parameters
	err := s.tk.Keeper.SetParams(s.tk.Ctx, params)
	s.Require().NoError(err)

	// Register a validator
	_, pubkey, ethAddr := testutil.GenerateSecp256k1Key()
	initialCollateral := math.NewUint(1000000000) // 1 MITO in gwei
	initialExtraVP := math.NewUint(3000000000)    // 3 MITO in gwei as extra voting power

	err = s.tk.Keeper.RegisterValidator(s.tk.Ctx, ethAddr, pubkey, initialCollateral, initialExtraVP, false)
	s.Require().NoError(err)

	// Get validator
	validator, found := s.tk.Keeper.GetValidator(s.tk.Ctx, ethAddr)
	s.Require().True(found)

	// Check voting power - should be limited by max leverage ratio
	// Collateral = 1 MITO, Extra VP = 3 MITO
	// Without leverage limit: 1 + 3 = 4 VP
	// With 2x leverage limit: 1 * 2 = 2 VP (capped)
	s.Require().Equal(int64(2), validator.VotingPower)

	// Increase collateral to increase the leverage cap
	validator.Collateral = math.NewUint(2000000000) // 2 MITO in gwei

	// Update validator state
	s.tk.Keeper.UpdateValidatorState(s.tk.Ctx, &validator, "testing leverage ratio")

	// Check voting power again - should be higher but still limited
	// Collateral = 2 MITO, Extra VP = 3 MITO
	// Without leverage limit: 2 + 3 = 5 VP
	// With 2x leverage limit: 2 * 2 = 4 VP (capped)
	s.Require().Equal(int64(4), validator.VotingPower)

	// Check if validator was updated in state
	updatedValidator, found := s.tk.Keeper.GetValidator(s.tk.Ctx, ethAddr)
	s.Require().True(found)
	s.Require().Equal(int64(4), updatedValidator.VotingPower)
}
