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
func (s *WithdrawTestSuite) createTestValidator(collateral math.Uint) (mitotypes.EthAddress, types.Validator) {
	_, pubkey, ethAddr := testutil.GenerateSecp256k1Key()
	err := s.tk.Keeper.RegisterValidator(s.tk.Ctx, ethAddr, pubkey, collateral, math.ZeroUint(), false)
	s.Require().NoError(err)

	validator, found := s.tk.Keeper.GetValidator(s.tk.Ctx, ethAddr)
	s.Require().True(found)
	return ethAddr, validator
}

// createAndAddWithdrawal creates a withdrawal and adds it to state
func (s *WithdrawTestSuite) createAndAddWithdrawal(
	valAddr mitotypes.EthAddress,
	amount uint64,
	receiver mitotypes.EthAddress,
	maturesAt int64,
) *types.Withdrawal {
	withdrawal := &types.Withdrawal{
		ValAddr:        valAddr,
		Amount:         amount,
		Receiver:       receiver,
		MaturesAt:      maturesAt,
		CreationHeight: s.tk.Ctx.BlockHeight(),
	}
	s.tk.Keeper.AddNewWithdrawalWithNextID(s.tk.Ctx, withdrawal)
	return withdrawal
}

func (s *WithdrawTestSuite) Test_ProcessMaturedWithdrawals() {
	// Setup test parameters
	s.setupTestParams()

	// Create a validator
	valAddr, _ := s.createTestValidator(math.NewUint(10000000000)) // 10 MITO
	receiverAddr := valAddr                                        // Using the same address for receiver for simplicity

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
	var found1, found2, found3, found4 bool
	s.tk.Keeper.IterateWithdrawalsByMaturesAt(ctx, func(w types.Withdrawal) bool {
		if w.ID == maturedWithdrawal1.ID {
			found1 = true
		}
		if w.ID == maturedWithdrawal2.ID {
			found2 = true
		}
		if w.ID == maturedWithdrawal3.ID {
			found3 = true
		}
		if w.ID == futureWithdrawal.ID {
			found4 = true
		}
		return false
	})

	// First two matured withdrawals should be deleted, third and future should exist
	s.Require().False(found1, "First matured withdrawal should be deleted")
	s.Require().False(found2, "Second matured withdrawal should be deleted")
	s.Require().True(found3, "Third matured withdrawal should still exist")
	s.Require().True(found4, "Future withdrawal should still exist")
}

func (s *WithdrawTestSuite) Test_ProcessMaturedWithdrawals_NoMaturedWithdrawals() {
	// Setup test parameters
	s.setupTestParams()

	// Create a validator
	valAddr, _ := s.createTestValidator(math.NewUint(10000000000)) // 10 MITO

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
	var found1, found2 bool
	s.tk.Keeper.IterateWithdrawalsByMaturesAt(ctx, func(w types.Withdrawal) bool {
		if w.ID == futureWithdrawal1.ID {
			found1 = true
		}
		if w.ID == futureWithdrawal2.ID {
			found2 = true
		}
		return false
	})
	s.Require().True(found1, "Future withdrawal 1 should still exist")
	s.Require().True(found2, "Future withdrawal 2 should still exist")
}

func (s *WithdrawTestSuite) Test_ProcessMaturedWithdrawals_InsertError() {
	// Setup test parameters
	s.setupTestParams()

	// Create a validator
	valAddr, _ := s.createTestValidator(math.NewUint(10000000000)) // 10 MITO

	// Set the current block time
	now := time.Now().Unix()
	ctx := s.tk.Ctx.WithBlockTime(time.Unix(now, 0))

	// Create a matured withdrawal
	maturedWithdrawal := s.createAndAddWithdrawal(valAddr, 1000000000, valAddr, now-86400)

	// Mock an error in the InsertWithdrawal function
	expectedError := "mocked insertion error"
	s.tk.MockEvmEng.InsertWithdrawalFn = func(ctx context.Context, withdrawalAddr common.Address, amountGwei uint64) error {
		return fmt.Errorf("%s", expectedError)
	}

	// Process matured withdrawals - should return the error
	err := s.tk.Keeper.ProcessMaturedWithdrawals(ctx)
	s.Require().Error(err)
	s.Require().Contains(err.Error(), expectedError)

	// The withdrawal should still exist since processing failed
	var found bool
	s.tk.Keeper.IterateWithdrawalsByMaturesAt(ctx, func(w types.Withdrawal) bool {
		if w.ID == maturedWithdrawal.ID {
			found = true
			return true
		}
		return false
	})
	s.Require().True(found, "Matured withdrawal should still exist")
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
	valAddr, _ := s.createTestValidator(math.NewUint(10000000000)) // 10 MITO

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

	// The other withdrawals should still exist
	var found2, found3 bool
	s.tk.Keeper.IterateWithdrawalsByMaturesAt(ctx, func(w types.Withdrawal) bool {
		if w.ID == maturedWithdrawal2.ID {
			found2 = true
		}
		if w.ID == maturedWithdrawal3.ID {
			found3 = true
		}
		return false
	})
	s.Require().True(found2, "Second matured withdrawal should still exist")
	s.Require().True(found3, "Third matured withdrawal should still exist")

	// Process withdrawals again - should process the next oldest one
	err = s.tk.Keeper.ProcessMaturedWithdrawals(ctx)
	s.Require().NoError(err)

	// Verify the second oldest withdrawal was processed
	s.Require().Equal(2, len(insertedWithdrawals[valAddrStr]), "Two withdrawals should be processed in total")
	s.Require().Contains(insertedWithdrawals[valAddrStr], uint64(2000000000), "Second withdrawal should be processed")

	// Only the newest withdrawal should still exist
	var found1Next, found2Next, found3Next bool
	s.tk.Keeper.IterateWithdrawalsByMaturesAt(ctx, func(w types.Withdrawal) bool {
		if w.ID == maturedWithdrawal1.ID {
			found1Next = true
		}
		if w.ID == maturedWithdrawal2.ID {
			found2Next = true
		}
		if w.ID == maturedWithdrawal3.ID {
			found3Next = true
		}
		return false
	})
	s.Require().False(found1Next, "First matured withdrawal should be deleted")
	s.Require().False(found2Next, "Second matured withdrawal should be deleted")
	s.Require().True(found3Next, "Third matured withdrawal should still exist")
}
