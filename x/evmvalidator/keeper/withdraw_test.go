package keeper_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"cosmossdk.io/math"
	"github.com/ethereum/go-ethereum/common"
	mitotypes "github.com/mitosis-org/chain/types"
	"github.com/mitosis-org/chain/x/evmvalidator/testutil"
	"github.com/mitosis-org/chain/x/evmvalidator/types"
	"github.com/stretchr/testify/suite"
)

// WithdrawTestSuite is a test suite for the ProcessMaturedWithdrawals functionality
type WithdrawTestSuite struct {
	suite.Suite
	tk testutil.TestKeeper
}

// SetupTest initializes the test suite
func (s *WithdrawTestSuite) SetupTest() {
	s.tk = testutil.CreateTestInput(&s.Suite)
}

// TestWithdrawTestSuite runs the withdraw test suite
func TestWithdrawTestSuite(t *testing.T) {
	suite.Run(t, new(WithdrawTestSuite))
}

// setupTestParams sets up test parameters
func (s *WithdrawTestSuite) setupTestParams() types.Params {
	params := types.Params{
		MaxValidators:    10,
		MaxLeverageRatio: math.LegacyNewDec(10),
		MinVotingPower:   1,
		WithdrawalLimit:  2, // Set a low limit to test the limit functionality
	}
	err := s.tk.Keeper.SetParams(s.tk.Ctx, params)
	s.Require().NoError(err)
	return params
}

// createTestValidator registers a validator for testing
func (s *WithdrawTestSuite) createTestValidator(collateral math.Uint) types.Validator {
	_, pubkey, ethAddr := testutil.GenerateSecp256k1Key()
	err := s.tk.Keeper.RegisterValidator(s.tk.Ctx, ethAddr, pubkey, collateral, math.ZeroUint(), false)
	s.Require().NoError(err)

	validator, found := s.tk.Keeper.GetValidator(s.tk.Ctx, ethAddr)
	s.Require().True(found)
	return validator
}

// createAndAddWithdrawal creates a withdrawal and adds it to state
func (s *WithdrawTestSuite) createAndAddWithdrawal(
	valAddr mitotypes.EthAddress,
	amount uint64,
	receiver mitotypes.EthAddress,
	maturesAt int64,
) types.Withdrawal {
	withdrawal := &types.Withdrawal{
		ValAddr:        valAddr,
		Amount:         amount,
		Receiver:       receiver,
		MaturesAt:      maturesAt,
		CreationHeight: s.tk.Ctx.BlockHeight(),
	}
	s.tk.Keeper.AddNewWithdrawalWithNextID(s.tk.Ctx, withdrawal)
	return *withdrawal
}

func (s *WithdrawTestSuite) Test_ProcessMaturedWithdrawals() {
	// Setup test parameters
	s.setupTestParams()

	// Create a validator
	validator := s.createTestValidator(math.NewUint(10000000000)) // 10 MITO
	valAddr := validator.Addr
	receiverAddr := valAddr // Using the same address for receiver for simplicity

	// Set the current block time
	now := time.Now().Unix()
	ctx := s.tk.Ctx.WithBlockTime(time.Unix(now, 0))

	// Create matured withdrawals with different maturity times
	maturedWithdrawal1 := s.createAndAddWithdrawal(valAddr, 3000000000, receiverAddr, now-86400*3) // 3 MITO, Oldest
	maturedWithdrawal2 := s.createAndAddWithdrawal(valAddr, 2000000000, receiverAddr, now-86400*2) // 2 MITO
	maturedWithdrawal3 := s.createAndAddWithdrawal(valAddr, 1000000000, receiverAddr, now-86400)   // 1 MITO, Newest
	futureWithdrawal := s.createAndAddWithdrawal(valAddr, 4000000000, receiverAddr, now+86400)     // 4 MITO, Future

	// Track inserted withdrawals
	insertedWithdrawals := make(map[string][]uint64)
	s.tk.MockEvmEng.InsertWithdrawalFn = func(ctx context.Context, withdrawalAddr common.Address, amountGwei uint64) error {
		addr := withdrawalAddr.String()
		insertedWithdrawals[addr] = append(insertedWithdrawals[addr], amountGwei)
		return nil
	}

	// Process matured withdrawals
	err := s.tk.Keeper.ProcessMaturedWithdrawals(ctx)
	s.Require().NoError(err)

	// Verify withdrawals were processed correctly
	// According to the withdrawal limit of 2 (from setupTestParams), we should process 2 withdrawals
	receiverAddrStr := receiverAddr.String()
	s.Require().Equal(2, len(insertedWithdrawals[receiverAddrStr]), "Two withdrawals should be processed")

	// The actual values will be the first two matured withdrawals processed in order of maturity
	// We expect the oldest two (maturedWithdrawal1 and maturedWithdrawal2) to be processed
	receiverWithdrawals := insertedWithdrawals[receiverAddrStr]
	s.Require().Contains(receiverWithdrawals, uint64(3000000000), "First matured withdrawal should be processed")
	s.Require().Contains(receiverWithdrawals, uint64(2000000000), "Second matured withdrawal should be processed")

	// Check which withdrawals still exist
	allWithdrawals := s.tk.Keeper.GetAllWithdrawals(ctx)
	s.Require().Equal(2, len(allWithdrawals), "There should be 2 withdrawals in total")
	s.Require().NotContains(allWithdrawals, maturedWithdrawal1, "First matured withdrawal should be deleted")
	s.Require().NotContains(allWithdrawals, maturedWithdrawal2, "Second matured withdrawal should be deleted")
	s.Require().Contains(allWithdrawals, maturedWithdrawal3, "Third matured withdrawal should still exist")
	s.Require().Contains(allWithdrawals, futureWithdrawal, "Future withdrawal should still exist")
}

func (s *WithdrawTestSuite) Test_ProcessMaturedWithdrawals_NoMaturedWithdrawals() {
	// Setup test parameters
	s.setupTestParams()

	// Create a validator
	validator := s.createTestValidator(math.NewUint(10000000000)) // 10 MITO
	valAddr := validator.Addr

	// Set the current block time
	now := time.Now().Unix()
	ctx := s.tk.Ctx.WithBlockTime(time.Unix(now, 0))

	// Create only future withdrawals
	futureWithdrawal1 := s.createAndAddWithdrawal(valAddr, 1000000000, valAddr, now+86400)
	futureWithdrawal2 := s.createAndAddWithdrawal(valAddr, 2000000000, valAddr, now+86400*2)

	// Track inserted withdrawals
	insertedWithdrawals := make(map[string]uint64)
	s.tk.MockEvmEng.InsertWithdrawalFn = func(ctx context.Context, withdrawalAddr common.Address, amountGwei uint64) error {
		insertedWithdrawals[withdrawalAddr.String()] = amountGwei
		return nil
	}

	// Process matured withdrawals
	err := s.tk.Keeper.ProcessMaturedWithdrawals(ctx)
	s.Require().NoError(err)

	// No withdrawals should be processed
	s.Require().Equal(0, len(insertedWithdrawals), "No withdrawals should be processed")

	// Future withdrawals should still exist
	allWithdrawals := s.tk.Keeper.GetAllWithdrawals(ctx)
	s.Require().Equal(2, len(allWithdrawals), "There should be 2 withdrawals in total")
	s.Require().Contains(allWithdrawals, futureWithdrawal1, "First future withdrawal should still exist")
	s.Require().Contains(allWithdrawals, futureWithdrawal2, "Second future withdrawal should still exist")
}

func (s *WithdrawTestSuite) Test_ProcessMaturedWithdrawals_InsertError() {
	// Setup test parameters
	s.setupTestParams()

	// Create a validator
	validator := s.createTestValidator(math.NewUint(10000000000)) // 10 MITO
	valAddr := validator.Addr

	// Set the current block time
	now := time.Now().Unix()
	ctx := s.tk.Ctx.WithBlockTime(time.Unix(now, 0))

	// Create a matured withdrawal
	_ = s.createAndAddWithdrawal(valAddr, 1000000000, valAddr, now-86400)

	// Mock an error in the InsertWithdrawal function
	expectedError := "mocked insertion error"
	s.tk.MockEvmEng.InsertWithdrawalFn = func(ctx context.Context, withdrawalAddr common.Address, amountGwei uint64) error {
		return fmt.Errorf("%s", expectedError)
	}

	// Process matured withdrawals - should return the error
	err := s.tk.Keeper.ProcessMaturedWithdrawals(ctx)
	s.Require().Error(err)
	s.Require().Contains(err.Error(), expectedError)
}

func (s *WithdrawTestSuite) Test_ProcessMaturedWithdrawals_WithdrawalLimit() {
	// Setup test parameters with a very low withdrawal limit
	params := types.Params{
		MaxValidators:    10,
		MaxLeverageRatio: math.LegacyNewDec(10),
		MinVotingPower:   1,
		WithdrawalLimit:  1, // Only process 1 withdrawal per block
	}
	err := s.tk.Keeper.SetParams(s.tk.Ctx, params)
	s.Require().NoError(err)

	// Create a validator
	validator := s.createTestValidator(math.NewUint(10000000000)) // 10 MITO
	valAddr := validator.Addr

	// Set the current block time
	now := time.Now().Unix()
	ctx := s.tk.Ctx.WithBlockTime(time.Unix(now, 0))

	// Create multiple matured withdrawals
	maturedWithdrawal1 := s.createAndAddWithdrawal(valAddr, 1000000000, valAddr, now-86400*3) // Oldest
	maturedWithdrawal2 := s.createAndAddWithdrawal(valAddr, 2000000000, valAddr, now-86400*2)
	maturedWithdrawal3 := s.createAndAddWithdrawal(valAddr, 3000000000, valAddr, now-86400) // Newest

	// Track inserted withdrawals
	insertedWithdrawals := make(map[string][]uint64)
	s.tk.MockEvmEng.InsertWithdrawalFn = func(ctx context.Context, withdrawalAddr common.Address, amountGwei uint64) error {
		addr := withdrawalAddr.String()
		insertedWithdrawals[addr] = append(insertedWithdrawals[addr], amountGwei)
		return nil
	}

	// Process matured withdrawals - should only process the oldest one
	err = s.tk.Keeper.ProcessMaturedWithdrawals(ctx)
	s.Require().NoError(err)

	// Verify only the oldest withdrawal was processed
	valAddrStr := valAddr.String()
	s.Require().Equal(1, len(insertedWithdrawals), "Only one validator should have withdrawals")
	s.Require().Equal(1, len(insertedWithdrawals[valAddrStr]), "Only one withdrawal should be processed")
	s.Require().Equal(uint64(1000000000), insertedWithdrawals[valAddrStr][0], "Oldest withdrawal should be processed")

	// Verify withdrawals state is correct
	allWithdrawals := s.tk.Keeper.GetAllWithdrawals(ctx)
	s.Require().Equal(2, len(allWithdrawals), "There should be 2 withdrawals in total")
	s.Require().NotContains(allWithdrawals, maturedWithdrawal1, "First matured withdrawal should be deleted")
	s.Require().Contains(allWithdrawals, maturedWithdrawal2, "Second matured withdrawal should still exist")
	s.Require().Contains(allWithdrawals, maturedWithdrawal3, "Third matured withdrawal should still exist")

	// Process withdrawals again - should process the next oldest one
	err = s.tk.Keeper.ProcessMaturedWithdrawals(ctx)
	s.Require().NoError(err)

	// Verify the second oldest withdrawal was processed
	s.Require().Equal(2, len(insertedWithdrawals[valAddrStr]), "Two withdrawals should be processed in total")
	s.Require().Contains(insertedWithdrawals[valAddrStr], uint64(2000000000), "Second withdrawal should be processed")

	// Only the newest withdrawal should still exist
	allWithdrawals = s.tk.Keeper.GetAllWithdrawals(ctx)
	s.Require().Equal(1, len(allWithdrawals), "There should be 1 withdrawal in total")
	s.Require().NotContains(allWithdrawals, maturedWithdrawal1, "First matured withdrawal should be deleted")
	s.Require().NotContains(allWithdrawals, maturedWithdrawal2, "Second matured withdrawal should be deleted")
	s.Require().Contains(allWithdrawals, maturedWithdrawal3, "Third matured withdrawal should still exist")
}
