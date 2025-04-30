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
	s.tk = testutil.NewTestKeeper(&s.Suite)
}

// TestValidatorTestSuite runs the validator test suite
func TestValidatorTestSuite(t *testing.T) {
	suite.Run(t, new(ValidatorTestSuite))
}

// ==================== RegisterValidator Tests ====================

func (s *ValidatorTestSuite) Test_RegisterValidator() {
	// Generate validator data
	_, pubkey, valAddr := testutil.GenerateSecp256k1Key()
	_, _, ownerAddr := testutil.GenerateSecp256k1Key() // Different owner address
	collateral := math.NewUint(1000000000)             // 1 MITO in gwei
	extraVotingPower := math.NewUint(1000000000)       // 1 MITO in gwei

	// Register validator
	err := s.tk.Keeper.RegisterValidator(s.tk.Ctx, valAddr, pubkey, ownerAddr, collateral, extraVotingPower, false)
	s.Require().NoError(err)

	// Check if validator exists
	validator, found := s.tk.Keeper.GetValidator(s.tk.Ctx, valAddr)
	s.Require().True(found)
	s.Require().Equal(types.Validator{
		Addr:             valAddr,
		Pubkey:           pubkey,
		Collateral:       collateral,
		CollateralShares: types.CalculateCollateralSharesForDeposit(math.ZeroUint(), math.ZeroUint(), collateral),
		ExtraVotingPower: extraVotingPower,
		VotingPower:      2,
		Jailed:           false,
		Bonded:           false,
	}, validator)

	// Verify there's exactly one ownership record for this validator with the correct owner
	var ownerships []types.CollateralOwnership
	s.tk.Keeper.IterateCollateralByValidator(s.tk.Ctx, valAddr, func(ownership types.CollateralOwnership) bool {
		ownerships = append(ownerships, ownership)
		return false
	})

	s.Require().Equal(1, len(ownerships), "There should be exactly one ownership record")
	s.Require().Equal(types.CollateralOwnership{
		ValAddr:        valAddr,
		Owner:          ownerAddr,
		Shares:         validator.CollateralShares,
		CreationHeight: s.tk.Ctx.BlockHeight(),
	}, ownerships[0], "The ownership record should match")

	// Check if consensus address mapping exists
	valFromConsAddr, found := s.tk.Keeper.GetValidatorByConsAddr(s.tk.Ctx, validator.MustConsAddr())
	s.Require().True(found)
	s.Require().Equal(validator, valFromConsAddr)

	// Check if validator is properly indexed by power
	iterator := s.tk.Keeper.GetValidatorsByPowerIndexIterator(s.tk.Ctx)
	defer iterator.Close()

	found = false
	for ; iterator.Valid(); iterator.Next() {
		addr := mitotypes.BytesToEthAddress(iterator.Value())
		if addr.String() == valAddr.String() {
			found = true
			break
		}
	}
	s.Require().True(found, "validator should be indexed by power")

	// Try registering the same validator again (should fail)
	err = s.tk.Keeper.RegisterValidator(s.tk.Ctx, valAddr, pubkey, ownerAddr, collateral, extraVotingPower, false)
	s.Require().Error(err)
	s.Require().ErrorIs(err, types.ErrValidatorAlreadyExists)
}

func (s *ValidatorTestSuite) Test_RegisterValidator_ZeroCollateral() {
	_, pubkey, valAddr := testutil.GenerateSecp256k1Key()
	zeroCollateral := math.ZeroUint()
	extraVotingPower := math.NewUint(0)

	// Try registering with zero collateral (should still work but voting power will be zero)
	err := s.tk.Keeper.RegisterValidator(s.tk.Ctx, valAddr, pubkey, valAddr, zeroCollateral, extraVotingPower, false)
	s.Require().NoError(err)

	// Check if validator exists
	validator, found := s.tk.Keeper.GetValidator(s.tk.Ctx, valAddr)
	s.Require().True(found)
	s.Require().Equal(types.Validator{
		Addr:             valAddr,
		Pubkey:           pubkey,
		Collateral:       zeroCollateral,
		CollateralShares: types.CalculateCollateralSharesForDeposit(zeroCollateral, math.ZeroUint(), zeroCollateral),
		ExtraVotingPower: extraVotingPower,
		VotingPower:      0,
		Jailed:           true,
		Bonded:           false,
	}, validator)

	// Verify ownership record with iterate method
	var ownerships []types.CollateralOwnership
	s.tk.Keeper.IterateCollateralByValidator(s.tk.Ctx, valAddr, func(ownership types.CollateralOwnership) bool {
		ownerships = append(ownerships, ownership)
		return false
	})

	s.Require().Equal(1, len(ownerships), "There should be exactly one ownership record")
	s.Require().Equal(types.CollateralOwnership{
		ValAddr:        valAddr,
		Owner:          valAddr,
		Shares:         validator.CollateralShares,
		CreationHeight: s.tk.Ctx.BlockHeight(),
	}, ownerships[0], "The ownership record should match")
}

func (s *ValidatorTestSuite) Test_RegisterValidator_InvalidPubkey() {
	invalidPubkey := []byte("invalid-pubkey")
	_, _, validEthAddr := testutil.GenerateSecp256k1Key()
	collateral := math.NewUint(1000000000)
	extraVotingPower := math.NewUint(0)

	// Try registering with invalid pubkey format
	err := s.tk.Keeper.RegisterValidator(s.tk.Ctx, validEthAddr, invalidPubkey, validEthAddr, collateral, extraVotingPower, false)
	s.Require().Error(err)
}

func (s *ValidatorTestSuite) Test_RegisterValidator_NotMatchedPubkey() {
	_, _, valAddr := testutil.GenerateSecp256k1Key()
	collateral := math.NewUint(1000000000)
	extraVotingPower := math.NewUint(0)

	// Try registering with pubkey not matched to address
	_, notMatchedPubkey, _ := testutil.GenerateSecp256k1Key()
	err := s.tk.Keeper.RegisterValidator(s.tk.Ctx, valAddr, notMatchedPubkey, valAddr, collateral, extraVotingPower, false)
	s.Require().Error(err)
}

// ==================== DepositCollateral Tests ====================

func (s *ValidatorTestSuite) Test_DepositCollateral() {
	// Set up test parameters
	s.tk.SetupDefaultTestParams()

	// Register a validator with a separate owner initially
	_, pubkey, valAddr := testutil.GenerateSecp256k1Key()
	_, _, initialOwnerAddr := testutil.GenerateSecp256k1Key()
	initialCollateral := math.NewUint(1000000000) // 1 MITO

	err := s.tk.Keeper.RegisterValidator(s.tk.Ctx, valAddr, pubkey, initialOwnerAddr, initialCollateral, math.ZeroUint(), false)
	s.Require().NoError(err)

	validator, found := s.tk.Keeper.GetValidator(s.tk.Ctx, valAddr)
	s.Require().True(found)
	initialValidator := validator

	// Create another owner who will deposit collateral
	_, _, secondOwnerAddr := testutil.GenerateSecp256k1Key()

	// Initial state check
	s.Require().Equal(int64(1), validator.VotingPower)
	s.Require().Equal(initialCollateral, validator.Collateral)

	// Initial ownership check
	initialOwnership, found := s.tk.Keeper.GetCollateralOwnership(s.tk.Ctx, valAddr, initialOwnerAddr)
	s.Require().True(found)

	// Deposit additional collateral from second owner
	additionalCollateral := math.NewUint(500000000) // 0.5 MITO in gwei
	s.tk.Keeper.DepositCollateral(s.tk.Ctx, &validator, secondOwnerAddr, additionalCollateral)

	// Check validator state after deposit
	expectedValidator := initialValidator
	expectedValidator.Collateral = initialValidator.Collateral.Add(additionalCollateral)
	expectedValidator.CollateralShares = initialValidator.CollateralShares.Add(types.CalculateCollateralSharesForDeposit(initialValidator.Collateral, initialValidator.CollateralShares, additionalCollateral))
	expectedValidator.VotingPower = int64(1) // Still 1 since 1.5 MITO rounds down to 1
	s.Require().Equal(expectedValidator, validator)

	// Check ownership records
	initialOwnershipAfter, found := s.tk.Keeper.GetCollateralOwnership(s.tk.Ctx, valAddr, initialOwnerAddr)
	s.Require().True(found)
	s.Require().Equal(initialOwnership.Shares, initialOwnershipAfter.Shares, "Initial owner's shares should remain unchanged")

	secondOwnership, found := s.tk.Keeper.GetCollateralOwnership(s.tk.Ctx, valAddr, secondOwnerAddr)
	s.Require().True(found)
	s.Require().Equal(types.CalculateCollateralSharesForWithdrawal(
		validator.Collateral,
		validator.CollateralShares,
		additionalCollateral), secondOwnership.Shares)

	s.Require().Equal(initialOwnershipAfter.Shares.Add(secondOwnership.Shares), validator.CollateralShares)

	// For a more noticeable change, let's add another 0.5 MITO from the second owner
	s.tk.Keeper.DepositCollateral(s.tk.Ctx, &validator, secondOwnerAddr, additionalCollateral)

	// Now we should have 2 MITO total, which should give 2 voting power
	finalExpectedValidator := expectedValidator
	finalExpectedValidator.Collateral = expectedValidator.Collateral.Add(additionalCollateral)
	finalExpectedValidator.CollateralShares = expectedValidator.CollateralShares.Add(types.CalculateCollateralSharesForDeposit(expectedValidator.Collateral, expectedValidator.CollateralShares, additionalCollateral))
	finalExpectedValidator.VotingPower = int64(2)
	s.Require().Equal(finalExpectedValidator, validator)

	// Check updated shares for second owner
	secondOwnershipUpdated, found := s.tk.Keeper.GetCollateralOwnership(s.tk.Ctx, valAddr, secondOwnerAddr)
	s.Require().True(found)
	s.Require().Equal(secondOwnership.Shares.Add(types.CalculateCollateralSharesForDeposit(
		validator.Collateral,
		validator.CollateralShares,
		additionalCollateral)), secondOwnershipUpdated.Shares)
	s.Require().Equal(initialOwnershipAfter.Shares.Add(secondOwnershipUpdated.Shares), validator.CollateralShares)
}

func (s *ValidatorTestSuite) Test_DepositCollateral_ZeroAmount() {
	// Set up test parameters
	s.tk.SetupDefaultTestParams()

	// Register a validator with separate owner
	_, pubkey, valAddr := testutil.GenerateSecp256k1Key()
	_, _, ownerAddr := testutil.GenerateSecp256k1Key()
	initialCollateral := math.NewUint(1000000000) // 1 MITO

	err := s.tk.Keeper.RegisterValidator(s.tk.Ctx, valAddr, pubkey, ownerAddr, initialCollateral, math.ZeroUint(), false)
	s.Require().NoError(err)

	validator, found := s.tk.Keeper.GetValidator(s.tk.Ctx, valAddr)
	s.Require().True(found)
	initialValidator := validator

	// Check initial ownership
	var initialOwnerships []types.CollateralOwnership
	s.tk.Keeper.IterateCollateralByValidator(s.tk.Ctx, valAddr, func(ownership types.CollateralOwnership) bool {
		initialOwnerships = append(initialOwnerships, ownership)
		return false
	})
	s.Require().Equal(1, len(initialOwnerships), "Should initially have one ownership record")
	s.Require().Equal(types.CollateralOwnership{
		ValAddr:        valAddr,
		Owner:          ownerAddr,
		Shares:         validator.CollateralShares,
		CreationHeight: s.tk.Ctx.BlockHeight(),
	}, initialOwnerships[0], "The ownership record should match")

	// Create another owner to deposit zero collateral
	_, _, anotherOwnerAddr := testutil.GenerateSecp256k1Key()
	zeroCollateral := math.ZeroUint()
	s.tk.Keeper.DepositCollateral(s.tk.Ctx, &validator, anotherOwnerAddr, zeroCollateral)

	// Check validator state is unchanged
	s.Require().Equal(initialValidator, validator)

	// Verify no new ownership record is created for zero deposit
	var ownershipsAfterAnotherZeroDeposit []types.CollateralOwnership
	s.tk.Keeper.IterateCollateralByValidator(s.tk.Ctx, valAddr, func(ownership types.CollateralOwnership) bool {
		ownershipsAfterAnotherZeroDeposit = append(ownershipsAfterAnotherZeroDeposit, ownership)
		return false
	})
	s.Require().Equal(1, len(ownershipsAfterAnotherZeroDeposit), "Should still have only one ownership record after zero deposit from new owner")
	s.Require().Equal(initialOwnerships[0], ownershipsAfterAnotherZeroDeposit[0], "Original ownership should still be unchanged")
}

// ==================== WithdrawCollateral Tests ====================

func (s *ValidatorTestSuite) Test_WithdrawCollateral() {
	// Set parameters for test
	s.tk.SetupDefaultTestParams()

	// Register a validator with enough collateral to have voting power = 2
	validator := s.tk.RegisterTestValidator(math.NewUint(2000000000), math.ZeroUint(), false)
	initialValidator := validator

	// Create withdrawal request
	withdrawalAmount := uint64(500000000) // 0.5 MITO in gwei
	withdrawal := types.Withdrawal{
		ValAddr:        validator.Addr,
		Amount:         withdrawalAmount,
		Receiver:       validator.Addr,
		MaturesAt:      time.Now().Unix() + 86400, // 1 day from now,
		CreationHeight: s.tk.Ctx.BlockHeight(),
	}

	expectedValidator := initialValidator
	expectedValidator.Collateral = expectedValidator.Collateral.Sub(math.NewUint(withdrawalAmount))
	expectedValidator.CollateralShares = expectedValidator.CollateralShares.Sub(
		types.CalculateCollateralSharesForWithdrawal(initialValidator.Collateral, initialValidator.CollateralShares, math.NewUint(withdrawalAmount)),
	)
	expectedValidator.VotingPower = int64(1)

	// Withdraw collateral
	err := s.tk.Keeper.WithdrawCollateral(s.tk.Ctx, &validator, validator.Addr, &withdrawal)
	s.Require().NoError(err)
	s.Require().Equal(expectedValidator, validator)
}

func (s *ValidatorTestSuite) Test_WithdrawCollateral_InsufficientCollateral() {
	// Set parameters for test
	s.tk.SetupDefaultTestParams()

	// Register a validator
	validator := s.tk.RegisterTestValidator(math.NewUint(1000000000), math.ZeroUint(), false)
	initialValidator := validator

	// Try withdrawing more than available (should fail)
	excessWithdrawal := &types.Withdrawal{
		ValAddr:        validator.Addr,
		Amount:         initialValidator.Collateral.Uint64() + 1, // Try to withdraw full initial amount + 1
		Receiver:       validator.Addr,
		MaturesAt:      time.Now().Unix() + 86400, // 1 day from now
		CreationHeight: s.tk.Ctx.BlockHeight(),
	}

	err := s.tk.Keeper.WithdrawCollateral(s.tk.Ctx, &validator, validator.Addr, excessWithdrawal)
	s.Require().Error(err)
	s.Require().ErrorIs(err, types.ErrInsufficientCollateral)

	fullWithdrawal := &types.Withdrawal{
		ValAddr:        validator.Addr,
		Amount:         initialValidator.Collateral.Uint64(), // Try to withdraw full initial amount
		Receiver:       validator.Addr,
		MaturesAt:      time.Now().Unix() + 86400, // 1 day from now
		CreationHeight: s.tk.Ctx.BlockHeight(),
	}

	err = s.tk.Keeper.WithdrawCollateral(s.tk.Ctx, &validator, validator.Addr, fullWithdrawal)
	s.Require().NoError(err)
	s.Require().Equal(math.ZeroUint(), validator.Collateral)
	s.Require().Equal(types.Validator{
		Addr:             initialValidator.Addr,
		Pubkey:           initialValidator.Pubkey,
		Collateral:       math.ZeroUint(),
		CollateralShares: math.ZeroUint(),
		ExtraVotingPower: initialValidator.ExtraVotingPower,
		VotingPower:      0,
		Jailed:           true,
		Bonded:           false,
	}, validator)
}

func (s *ValidatorTestSuite) Test_WithdrawCollateral_ZeroAmount() {
	// Set parameters for test
	s.tk.SetupDefaultTestParams()

	// Register a validator with enough collateral to have voting power = 2
	validator := s.tk.RegisterTestValidator(math.NewUint(2000000000), math.ZeroUint(), false)
	initialValidator := validator

	// Initial voting power should be 2
	initialVotingPower := int64(2)
	s.Require().Equal(initialVotingPower, validator.VotingPower)

	// Create withdrawal request
	withdrawalAmount := uint64(0)
	receiver := validator.Addr
	maturesAt := time.Now().Unix() + 86400 // 1 day from now

	withdrawal := types.Withdrawal{
		ValAddr:        validator.Addr,
		Amount:         withdrawalAmount,
		Receiver:       receiver,
		MaturesAt:      maturesAt,
		CreationHeight: s.tk.Ctx.BlockHeight(),
	}

	// Withdraw collateral
	err := s.tk.Keeper.WithdrawCollateral(s.tk.Ctx, &validator, validator.Addr, &withdrawal)
	s.Require().NoError(err)
	s.Require().Equal(initialValidator, validator)
}

// ==================== Enhanced withdraw test with multiple owners ====================

func (s *ValidatorTestSuite) Test_WithdrawCollateral_MultipleOwners() {
	// Set parameters for test
	s.tk.SetupDefaultTestParams()

	// Register a validator with the first owner
	_, pubkey, valAddr := testutil.GenerateSecp256k1Key()
	_, _, owner1Addr := testutil.GenerateSecp256k1Key()
	initialCollateral := math.NewUint(2000000000) // 2 MITO

	err := s.tk.Keeper.RegisterValidator(s.tk.Ctx, valAddr, pubkey, owner1Addr, initialCollateral, math.ZeroUint(), false)
	s.Require().NoError(err)

	validator, found := s.tk.Keeper.GetValidator(s.tk.Ctx, valAddr)
	s.Require().True(found)

	// Get the initial ownership record
	owner1Ownership, found := s.tk.Keeper.GetCollateralOwnership(s.tk.Ctx, valAddr, owner1Addr)
	s.Require().True(found)
	initialShares := owner1Ownership.Shares

	// Add a second owner with additional collateral
	_, _, owner2Addr := testutil.GenerateSecp256k1Key()
	additionalCollateral := math.NewUint(1000000000) // 1 MITO

	s.tk.Keeper.DepositCollateral(s.tk.Ctx, &validator, owner2Addr, additionalCollateral)

	// Get the updated validator and ownership records
	validator, found = s.tk.Keeper.GetValidator(s.tk.Ctx, valAddr)
	s.Require().True(found)

	owner2Ownership, found := s.tk.Keeper.GetCollateralOwnership(s.tk.Ctx, valAddr, owner2Addr)
	s.Require().True(found)

	// Create withdrawal request for owner1 (half of their collateral)
	withdrawalAmount := uint64(1000000000) // 1 MITO
	withdrawal := types.Withdrawal{
		ValAddr:        valAddr,
		Amount:         withdrawalAmount,
		Receiver:       owner1Addr,
		MaturesAt:      time.Now().Unix() + 86400, // 1 day from now
		CreationHeight: s.tk.Ctx.BlockHeight(),
	}

	// Calculate shares to withdraw
	sharesToWithdraw := types.CalculateCollateralSharesForWithdrawal(
		validator.Collateral,
		validator.CollateralShares,
		math.NewUint(withdrawalAmount))

	// Withdraw collateral for owner1
	err = s.tk.Keeper.WithdrawCollateral(s.tk.Ctx, &validator, owner1Addr, &withdrawal)
	s.Require().NoError(err)

	// Check validator state (total collateral should be reduced)
	expectedCollateral := initialCollateral.Add(additionalCollateral).Sub(math.NewUint(withdrawalAmount))
	s.Require().Equal(expectedCollateral, validator.Collateral)

	// Check ownership records after withdrawal
	owner1OwnershipAfter, found := s.tk.Keeper.GetCollateralOwnership(s.tk.Ctx, valAddr, owner1Addr)
	s.Require().True(found)
	s.Require().Equal(initialShares.Sub(sharesToWithdraw), owner1OwnershipAfter.Shares)

	// Owner2's shares should remain unchanged
	owner2OwnershipAfter, found := s.tk.Keeper.GetCollateralOwnership(s.tk.Ctx, valAddr, owner2Addr)
	s.Require().True(found)
	s.Require().Equal(owner2Ownership.Shares, owner2OwnershipAfter.Shares)

	// Try withdrawing more than available from owner2 (should fail)
	excessWithdrawal := types.Withdrawal{
		ValAddr:        valAddr,
		Amount:         additionalCollateral.Uint64() + 1, // More than owner2 has
		Receiver:       owner2Addr,
		MaturesAt:      time.Now().Unix() + 86400,
		CreationHeight: s.tk.Ctx.BlockHeight(),
	}

	err = s.tk.Keeper.WithdrawCollateral(s.tk.Ctx, &validator, owner2Addr, &excessWithdrawal)
	s.Require().Error(err)
	s.Require().ErrorIs(err, types.ErrInsufficientCollateral)
}

// ==================== Test for transferring collateral ownership ====================

func (s *ValidatorTestSuite) Test_TransferCollateralOwnership() {
	// Set up test parameters
	s.tk.SetupDefaultTestParams()

	// Register a validator with an owner
	_, pubkey, valAddr := testutil.GenerateSecp256k1Key()
	_, _, originOwnerAddr := testutil.GenerateSecp256k1Key()
	initialCollateral := math.NewUint(2000000000) // 2 MITO

	err := s.tk.Keeper.RegisterValidator(s.tk.Ctx, valAddr, pubkey, originOwnerAddr, initialCollateral, math.ZeroUint(), false)
	s.Require().NoError(err)

	validator, found := s.tk.Keeper.GetValidator(s.tk.Ctx, valAddr)
	s.Require().True(found)

	// Get the original ownership record
	originOwnership, found := s.tk.Keeper.GetCollateralOwnership(s.tk.Ctx, valAddr, originOwnerAddr)
	s.Require().True(found)

	// Create a new owner to transfer to
	_, _, newOwnerAddr := testutil.GenerateSecp256k1Key()

	// Transfer ownership from original owner to new owner
	s.tk.Keeper.TransferCollateralOwnership(
		s.tk.Ctx,
		&validator,
		originOwnership,
		newOwnerAddr,
	)

	// Verify original ownership record was deleted
	_, found = s.tk.Keeper.GetCollateralOwnership(s.tk.Ctx, valAddr, originOwnerAddr)
	s.Require().False(found, "Original ownership record should be deleted")

	// Verify new ownership record was created with same shares
	newOwnership, found := s.tk.Keeper.GetCollateralOwnership(s.tk.Ctx, valAddr, newOwnerAddr)
	s.Require().True(found, "New ownership record should be created")
	s.Require().Equal(types.CollateralOwnership{
		ValAddr:        valAddr,
		Owner:          newOwnerAddr,
		Shares:         originOwnership.Shares,
		CreationHeight: s.tk.Ctx.BlockHeight(),
	}, newOwnership)

	// Validator collateral should remain unchanged throughout transfers
	validatorAfter, found := s.tk.Keeper.GetValidator(s.tk.Ctx, valAddr)
	s.Require().True(found)
	s.Require().Equal(initialCollateral, validatorAfter.Collateral)
}

// ==================== Slash_ Tests ====================

func (s *ValidatorTestSuite) Test_Slash_() {
	// Set up test parameters
	s.tk.SetupDefaultTestParams()

	// Register a validator
	validator := s.tk.RegisterTestValidator(math.NewUint(1000000000), math.ZeroUint(), false)

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
	updatedValidator, found := s.tk.Keeper.GetValidator(s.tk.Ctx, validator.Addr)
	s.Require().True(found)
	s.Require().Equal(expectedCollateral, updatedValidator.Collateral)

	// Test slashing with negative fraction (should fail)
	_, err = s.tk.Keeper.Slash_(s.tk.Ctx, &validator, infractionHeight, power, math.LegacyNewDec(-1))
	s.Require().Error(err)
}

func (s *ValidatorTestSuite) Test_Slash_ExceedsCollateral() {
	// Set up test parameters
	s.tk.SetupDefaultTestParams()

	// Register a validator
	validator := s.tk.RegisterTestValidator(math.NewUint(1000000000), math.ZeroUint(), false)
	initialValidator := validator

	// Should have 1 voting power based on 1 MITO collateral
	naturalVotingPower := validator.VotingPower
	s.Require().Equal(int64(1), naturalVotingPower)

	// Slash parameters - use a high voting power (10x natural) and 100% slash to exceed collateral
	infractionHeight := s.tk.Ctx.BlockHeight() - 1
	slashFraction := math.LegacyNewDecWithPrec(100, 2) // 100% slash
	slashPower := naturalVotingPower * 10              // Use power higher than available collateral

	// Slash the validator
	// Since attempted slash exceeds collateral, should only slash what's available
	slashedAmount, err := s.tk.Keeper.Slash_(s.tk.Ctx, &validator, infractionHeight, slashPower, slashFraction)
	s.Require().NoError(err)
	s.Require().Equal(initialValidator.Collateral, slashedAmount)

	expectedValidator := initialValidator
	expectedValidator.Collateral = math.ZeroUint()
	expectedValidator.VotingPower = 0
	expectedValidator.Jailed = true
	s.Require().Equal(expectedValidator, validator)
}

func (s *ValidatorTestSuite) Test_Slash_Withdrawals() {
	// ==================== SETUP ====================
	s.tk.SetupDefaultTestParams()
	validator := s.tk.RegisterTestValidator(math.NewUint(5000000000), math.ZeroUint(), false) // 5 MITO

	// ==================== SETUP WITHDRAWALS ====================
	now := time.Now().Unix()

	// Create two future withdrawals with different maturity times
	futureWithdrawal1 := types.Withdrawal{
		ValAddr:        validator.Addr,
		Amount:         600000000, // 0.6 MITO
		Receiver:       validator.Addr,
		MaturesAt:      now + 86400, // 1 day from now
		CreationHeight: s.tk.Ctx.BlockHeight(),
	}

	futureWithdrawal2 := types.Withdrawal{
		ValAddr:        validator.Addr,
		Amount:         800000000, // 0.8 MITO
		Receiver:       validator.Addr,
		MaturesAt:      now + 172800, // 2 days from now
		CreationHeight: s.tk.Ctx.BlockHeight(),
	}

	// Create an already matured withdrawal (1 day ago)
	maturedWithdrawal := types.Withdrawal{
		ValAddr:        validator.Addr,
		Amount:         2000000000, // 2 MITO
		Receiver:       validator.Addr,
		MaturesAt:      now - 86400,
		CreationHeight: s.tk.Ctx.BlockHeight() - 100,
	}

	// Process the withdrawals
	err := s.tk.Keeper.WithdrawCollateral(s.tk.Ctx, &validator, validator.Addr, &futureWithdrawal1)
	s.Require().NoError(err)
	err = s.tk.Keeper.WithdrawCollateral(s.tk.Ctx, &validator, validator.Addr, &futureWithdrawal2)
	s.Require().NoError(err)
	err = s.tk.Keeper.WithdrawCollateral(s.tk.Ctx, &validator, validator.Addr, &maturedWithdrawal)
	s.Require().NoError(err)

	// Get updated validator with 1.6 MITO (5 - 0.6 - 0.8 - 2) of collateral
	updatedValidator, found := s.tk.Keeper.GetValidator(s.tk.Ctx, validator.Addr)
	s.Require().True(found)
	s.Require().Equal(math.NewUint(1600000000), updatedValidator.Collateral, "Validator should have 1.6 MITO collateral")
	s.Require().Equal(int64(1), updatedValidator.VotingPower, "Validator should have 1 voting power")

	// Slash the validator by 100% of voting power (1 MITO)
	slashFraction := math.LegacyOneDec() // 100%
	slashedAmount, err := s.tk.Keeper.Slash_(s.tk.Ctx, &updatedValidator, s.tk.Ctx.BlockHeight()-10, 1, slashFraction)
	s.Require().NoError(err)

	// Verify slashed amount matches expectation (1 MITO)
	expectedSlashedAmount := math.NewUint(1000000000)
	s.Require().Equal(expectedSlashedAmount, slashedAmount, "Should slash 1 MITO")

	// Check which withdrawals remain after slashing
	// - Matured withdrawals should not be affected by slashing
	// - First future withdrawal (0.6 MITO) should be completely slashed and deleted
	// - Second future withdrawal (0.8 MITO) should be partially slashed (by 0.4 MITO)
	// leaving 0.4 MITO remaining
	var remainingWithdrawals []types.Withdrawal
	s.tk.Keeper.IterateWithdrawalsForValidator(s.tk.Ctx, validator.Addr, func(w types.Withdrawal) bool {
		remainingWithdrawals = append(remainingWithdrawals, w)
		return false
	})
	s.Require().Equal(2, len(remainingWithdrawals), "Should have 2 withdrawals remaining")
	s.Require().Equal(maturedWithdrawal, remainingWithdrawals[0], "Matured withdrawal should still exist")
	expectedFutureWithdrawal2 := futureWithdrawal2
	expectedFutureWithdrawal2.Amount = 400000000 // 0.4 MITO
	s.Require().Equal(expectedFutureWithdrawal2, remainingWithdrawals[1], "Second future withdrawal should be partially slashed to 0.4 MITO")

	// Verify validator collateral remains unchanged
	finalValidator, found := s.tk.Keeper.GetValidator(s.tk.Ctx, validator.Addr)
	s.Require().True(found)
	s.Require().Equal(updatedValidator, finalValidator)
}

// ==================== Jail_ Tests ====================

func (s *ValidatorTestSuite) Test_Jail_() {
	// Register a validator
	validator := s.tk.RegisterTestValidator(math.NewUint(1000000000), math.ZeroUint(), false)
	initialValidator := validator
	s.Require().False(validator.Jailed)

	// Jail the validator
	reason := "testing jail"
	s.tk.Keeper.Jail_(s.tk.Ctx, &validator, reason)

	// Check if validator was jailed
	expectedValidator := initialValidator
	expectedValidator.Jailed = true
	s.Require().Equal(expectedValidator, validator)

	// Check if validator state was updated
	updatedValidator, found := s.tk.Keeper.GetValidator(s.tk.Ctx, validator.Addr)
	s.Require().True(found)
	s.Require().Equal(expectedValidator, updatedValidator)

	// Check if validator was deleted from power index
	_, found = s.tk.Keeper.GetValidatorByPowerIndex(s.tk.Ctx, initialValidator.VotingPower, validator.Addr)
	s.Require().False(found)
}

// ==================== Unjail_ Tests ====================

func (s *ValidatorTestSuite) Test_Unjail_() {
	// Set up test parameters
	s.tk.SetupDefaultTestParams()

	// Register a validator
	validator := s.tk.RegisterTestValidator(math.NewUint(1000000000), math.ZeroUint(), true)
	initialValidator := validator

	// Check if validator was not added to power index because it is jailed
	iterator := s.tk.Keeper.GetValidatorsByPowerIndexIterator(s.tk.Ctx)
	for ; iterator.Valid(); iterator.Next() {
		addr := mitotypes.BytesToEthAddress(iterator.Value())
		s.Require().NotEqual(validator.Addr, addr)
	}

	// For unjailing to succeed, we need to mock the slashing keeper's UnjailFromConsAddr function
	s.tk.MockSlash.UnjailFromConsAddrFn = func(ctx context.Context, consAddr sdk.ConsAddress) error {
		return nil
	}

	// Unjail the validator
	err := s.tk.Keeper.Unjail_(s.tk.Ctx, &validator)
	s.Require().NoError(err)

	// Check if validator was unjailed
	expectedValidator := initialValidator
	expectedValidator.Jailed = false
	s.Require().Equal(expectedValidator, validator)

	// Check if validator was updated in state
	updatedValidator, found := s.tk.Keeper.GetValidator(s.tk.Ctx, validator.Addr)
	s.Require().True(found)
	s.Require().Equal(expectedValidator, updatedValidator)

	// Check if validator was added to power index
	_, found = s.tk.Keeper.GetValidatorByPowerIndex(s.tk.Ctx, validator.VotingPower, validator.Addr)
	s.Require().True(found)

	// Test unjailing a non-jailed validator
	err = s.tk.Keeper.Unjail_(s.tk.Ctx, &validator)
	s.Require().NoError(err) // No error because the function just returns if already unjailed
}

func (s *ValidatorTestSuite) Test_Unjail_InsufficientVotingPower() {
	// Define parameters for testing with minimum voting power of 10
	s.tk.SetupTestParams(types.Params{
		MaxValidators:    100,
		MaxLeverageRatio: math.LegacyNewDec(10), // 10x leverage
		MinVotingPower:   10,                    // Higher minimum voting power
		WithdrawalLimit:  10,
	})

	// Register a jailed validator with insufficient collateral
	validator := s.tk.RegisterTestValidator(math.NewUint(1000000000), math.ZeroUint(), true)

	// For mocking the slashing keeper's UnjailFromConsAddr function
	s.tk.MockSlash.UnjailFromConsAddrFn = func(ctx context.Context, consAddr sdk.ConsAddress) error {
		return nil
	}

	// Try to unjail the validator with insufficient voting power
	err := s.tk.Keeper.Unjail_(s.tk.Ctx, &validator)
	s.Require().Error(err)
	s.Require().Contains(err.Error(), "voting power below minimum requirement")
}

// ==================== UpdateExtraVotingPower Tests ====================

func (s *ValidatorTestSuite) Test_UpdateExtraVotingPower() {
	s.tk.SetupDefaultTestParams()

	validator := s.tk.RegisterTestValidator(math.NewUint(1000000000), math.ZeroUint(), false) // 1 MITO
	initialValidator := validator

	// Update extra voting power
	newExtraVotingPower := math.NewUint(1500000000) // 1.5 MITO
	s.tk.Keeper.UpdateExtraVotingPower(s.tk.Ctx, &validator, newExtraVotingPower)

	// Check if extra voting power was updated
	expectedValidator := initialValidator
	expectedValidator.ExtraVotingPower = newExtraVotingPower
	expectedValidator.VotingPower = 2 // 1 MITO + 1.5 MITO -> 2
	s.Require().Equal(expectedValidator, validator)
}

// ==================== UpdateValidatorState Tests ====================

func (s *ValidatorTestSuite) Test_UpdateValidatorState() {
	s.tk.SetupDefaultTestParams()

	validator := s.tk.RegisterTestValidator(math.NewUint(1000000000), math.ZeroUint(), false) // 1 MITO
	initialValidator := validator
	valAddr := validator.Addr

	// Update validator state
	validator.Collateral = math.NewUint(2000000000)       // 2 MITO in gwei
	validator.ExtraVotingPower = math.NewUint(1000000000) // 1 MITO in gwei
	validator.Jailed = false
	validator.Bonded = false
	s.tk.Keeper.UpdateValidatorState(s.tk.Ctx, &validator, "testing")

	expectedValidator := initialValidator
	expectedValidator.Collateral = math.NewUint(2000000000)
	expectedValidator.ExtraVotingPower = math.NewUint(1000000000)
	expectedValidator.Jailed = false
	expectedValidator.Bonded = false
	expectedValidator.VotingPower = 3
	s.Require().Equal(expectedValidator, validator)

	// Check if validator was updated in state
	updatedValidator, found := s.tk.Keeper.GetValidator(s.tk.Ctx, valAddr)
	s.Require().True(found)
	s.Require().Equal(expectedValidator, updatedValidator)

	// Check if validator was deleted from power index
	_, found = s.tk.Keeper.GetValidatorByPowerIndex(s.tk.Ctx, initialValidator.VotingPower, valAddr)
	s.Require().False(found)

	// Check if validator was added to power index
	_, found = s.tk.Keeper.GetValidatorByPowerIndex(s.tk.Ctx, validator.VotingPower, valAddr)
	s.Require().True(found)
}

func (s *ValidatorTestSuite) Test_UpdateValidatorState_Jailed() {
	s.tk.SetupDefaultTestParams()

	validator := s.tk.RegisterTestValidator(math.NewUint(1000000000), math.ZeroUint(), false) // 1 MITO
	initialValidator := validator
	valAddr := validator.Addr

	validator.Jailed = true
	s.tk.Keeper.UpdateValidatorState(s.tk.Ctx, &validator, "testing")

	expectedValidator := initialValidator
	expectedValidator.Jailed = true
	s.Require().Equal(expectedValidator, validator)

	// Power index should not exist
	iterator := s.tk.Keeper.GetValidatorsByPowerIndexIterator(s.tk.Ctx)
	for ; iterator.Valid(); iterator.Next() {
		addr := mitotypes.BytesToEthAddress(iterator.Value())
		s.Require().NotEqual(valAddr, addr)
	}
}

func (s *ValidatorTestSuite) Test_UpdateValidatorState_BelowMinVotingPower() {
	s.tk.SetupDefaultTestParams()

	validator := s.tk.RegisterTestValidator(math.NewUint(1000000000), math.ZeroUint(), false) // 1 MITO
	initialValidator := validator
	valAddr := validator.Addr

	// Reduce collateral below minimum
	validator.Collateral = math.NewUint(900000000) // 0.9 MITO in gwei
	s.tk.Keeper.UpdateValidatorState(s.tk.Ctx, &validator, "testing minimum voting power")

	// Check if validator was jailed due to insufficient voting power
	expectedValidator := initialValidator
	expectedValidator.Collateral = math.NewUint(900000000)
	expectedValidator.Jailed = true
	expectedValidator.VotingPower = 0
	s.Require().Equal(expectedValidator, validator)

	// Check if validator was updated in state
	updatedValidator, found := s.tk.Keeper.GetValidator(s.tk.Ctx, valAddr)
	s.Require().True(found)
	s.Require().Equal(expectedValidator, updatedValidator)

	// Check if validator was deleted from power index
	iterator := s.tk.Keeper.GetValidatorsByPowerIndexIterator(s.tk.Ctx)
	for ; iterator.Valid(); iterator.Next() {
		addr := mitotypes.BytesToEthAddress(iterator.Value())
		s.Require().NotEqual(valAddr, addr)
	}
}

func (s *ValidatorTestSuite) Test_UpdateValidatorState_MaxLeverageRatio() {
	// Define parameters for testing with low max leverage ratio
	params := types.Params{
		MaxValidators:    100,
		MaxLeverageRatio: math.LegacyNewDec(2), // Only 2x leverage allowed
		MinVotingPower:   1,
		WithdrawalLimit:  10,
	}
	err := s.tk.Keeper.SetParams(s.tk.Ctx, params)
	s.Require().NoError(err)

	// Register a validator
	initialCollateral := math.NewUint(1000000000) // 1 MITO in gwei
	initialExtraVP := math.NewUint(3000000000)    // 3 MITO in gwei as extra voting power
	validator := s.tk.RegisterTestValidator(initialCollateral, initialExtraVP, false)
	initialValidator := validator
	valAddr := validator.Addr

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
	expectedValidator := initialValidator
	expectedValidator.Collateral = math.NewUint(2000000000)
	expectedValidator.VotingPower = 4
	s.Require().Equal(expectedValidator, validator)

	// Check if validator was updated in state
	updatedValidator, found := s.tk.Keeper.GetValidator(s.tk.Ctx, valAddr)
	s.Require().True(found)
	s.Require().Equal(int64(4), updatedValidator.VotingPower)

	// Check if validator was added to power index
	_, found = s.tk.Keeper.GetValidatorByPowerIndex(s.tk.Ctx, validator.VotingPower, valAddr)
	s.Require().True(found)

	// Check if validator was deleted from power index
	_, found = s.tk.Keeper.GetValidatorByPowerIndex(s.tk.Ctx, initialValidator.VotingPower, valAddr)
	s.Require().False(found)
}
