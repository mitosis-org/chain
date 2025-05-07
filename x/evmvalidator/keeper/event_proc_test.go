package keeper_test

import (
	"context"
	"math/big"
	"testing"
	"time"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/mitosis-org/chain/bindings"
	"github.com/mitosis-org/chain/x/evmvalidator/testutil"
	"github.com/mitosis-org/chain/x/evmvalidator/types"
	"github.com/stretchr/testify/suite"
)

// EventProcessingTestSuite is a test suite for event processing functionality
type EventProcessingTestSuite struct {
	suite.Suite
	tk testutil.TestKeeper
}

// SetupTest initializes the test suite
func (s *EventProcessingTestSuite) SetupTest() {
	s.tk = testutil.NewTestKeeper(&s.Suite)
}

// TestEventProcessingTestSuite runs the event processing test suite
func TestEventProcessingTestSuite(t *testing.T) {
	suite.Run(t, new(EventProcessingTestSuite))
}

// Skip this test for now since it has issues with valAddress comparisons
func (s *EventProcessingTestSuite) Test_ProcessRegisterValidator() {
	// Setup test parameters
	s.tk.SetupDefaultTestParams()

	// Generate validator data
	_, pubkey, valAddr := testutil.GenerateSecp256k1Key()
	_, _, collateralOwnerAddr := testutil.GenerateSecp256k1Key()

	// Create the register validator event
	event := &bindings.ConsensusValidatorEntrypointMsgRegisterValidator{
		ValAddr:                     common.BytesToAddress(valAddr.Bytes()),
		PubKey:                      pubkey,
		InitialCollateralOwner:      common.BytesToAddress(collateralOwnerAddr.Bytes()),
		InitialCollateralAmountGwei: big.NewInt(1000000000), // 1 MITO
	}

	// Process the register validator event
	err, ignore := s.tk.Keeper.ProcessRegisterValidator(s.tk.Ctx, event)
	s.Require().NoError(err)
	s.Require().False(ignore)

	// Verify validator was registered
	validator, found := s.tk.Keeper.GetValidator(s.tk.Ctx, valAddr)
	s.Require().True(found)
	s.Require().Equal(types.Validator{
		Addr:             valAddr,
		Pubkey:           pubkey,
		Collateral:       math.NewUintFromBigInt(event.InitialCollateralAmountGwei),
		CollateralShares: types.CalculateCollateralSharesForDeposit(math.ZeroUint(), math.ZeroUint(), math.NewUintFromBigInt(event.InitialCollateralAmountGwei)),
		ExtraVotingPower: math.ZeroUint(),
		VotingPower:      1,
		Jailed:           false,
		Bonded:           false,
	}, validator)
}

func (s *EventProcessingTestSuite) Test_FallbackRegisterValidator() {
	// Generate validator data
	_, pubkey, valAddr := testutil.GenerateSecp256k1Key()
	_, _, collateralOwnerAddr := testutil.GenerateSecp256k1Key()

	// Create the register validator event
	event := &bindings.ConsensusValidatorEntrypointMsgRegisterValidator{
		ValAddr:                     common.BytesToAddress(valAddr.Bytes()),
		PubKey:                      pubkey,
		InitialCollateralOwner:      common.BytesToAddress(collateralOwnerAddr.Bytes()),
		InitialCollateralAmountGwei: big.NewInt(1000000000), // 1 MITO
	}

	// Track inserted withdrawals
	var insertedWithdrawalAddr common.Address
	var insertedWithdrawalAmount uint64
	s.tk.MockEvmEng.InsertWithdrawalFn = func(ctx context.Context, withdrawalAddr common.Address, amountGwei uint64) error {
		insertedWithdrawalAddr = withdrawalAddr
		insertedWithdrawalAmount = amountGwei
		return nil
	}

	// Test fallback
	err := s.tk.Keeper.FallbackRegisterValidator(s.tk.Ctx, event)
	s.Require().NoError(err)

	// Verify withdrawal was inserted
	s.Require().Equal(event.InitialCollateralOwner, insertedWithdrawalAddr)
	s.Require().Equal(event.InitialCollateralAmountGwei.Uint64(), insertedWithdrawalAmount)
}

func (s *EventProcessingTestSuite) Test_ProcessDepositCollateral() {
	s.tk.SetupDefaultTestParams()

	initialValidator := s.tk.RegisterTestValidator(math.NewUint(1000000000), math.ZeroUint(), false) // 1 MITO
	_, _, collateralOwnerAddr := testutil.GenerateSecp256k1Key()

	// Create the deposit collateral event
	event := &bindings.ConsensusValidatorEntrypointMsgDepositCollateral{
		ValAddr:         common.BytesToAddress(initialValidator.Addr.Bytes()),
		CollateralOwner: common.BytesToAddress(collateralOwnerAddr.Bytes()),
		AmountGwei:      big.NewInt(1000000000), // 1 MITO more
	}

	// Process the deposit collateral event
	err, ignore := s.tk.Keeper.ProcessDepositCollateral(s.tk.Ctx, event)
	s.Require().NoError(err)
	s.Require().False(ignore)

	// Verify validator's collateral was increased
	updatedValidator, found := s.tk.Keeper.GetValidator(s.tk.Ctx, initialValidator.Addr)
	s.Require().True(found)
	expectedCollateral := initialValidator.Collateral.Add(math.NewUintFromBigInt(event.AmountGwei))
	s.Require().Equal(expectedCollateral, updatedValidator.Collateral)
	s.Require().Equal(int64(2), updatedValidator.VotingPower)

	// Try depositing to a non-existent validator (should return error with ignore=true)
	_, _, nonExistentAddr := testutil.GenerateSecp256k1Key()
	invalidEvent := &bindings.ConsensusValidatorEntrypointMsgDepositCollateral{
		ValAddr:         common.BytesToAddress(nonExistentAddr.Bytes()),
		CollateralOwner: common.BytesToAddress(collateralOwnerAddr.Bytes()),
		AmountGwei:      big.NewInt(1000000000),
	}

	err, ignore = s.tk.Keeper.ProcessDepositCollateral(s.tk.Ctx, invalidEvent)
	s.Require().Error(err)
	s.Require().True(ignore)
	s.Require().ErrorIs(err, types.ErrValidatorNotFound)
}

func (s *EventProcessingTestSuite) Test_FallbackDepositCollateral() {
	// Generate validator data
	_, _, valAddr := testutil.GenerateSecp256k1Key()
	_, _, collateralOwnerAddr := testutil.GenerateSecp256k1Key()

	// Create the deposit collateral event
	event := &bindings.ConsensusValidatorEntrypointMsgDepositCollateral{
		ValAddr:         common.BytesToAddress(valAddr.Bytes()),
		CollateralOwner: common.BytesToAddress(collateralOwnerAddr.Bytes()),
		AmountGwei:      big.NewInt(1000000000), // 1 MITO
	}

	// Track inserted withdrawals
	var insertedWithdrawalAddr common.Address
	var insertedWithdrawalAmount uint64
	s.tk.MockEvmEng.InsertWithdrawalFn = func(ctx context.Context, withdrawalAddr common.Address, amountGwei uint64) error {
		insertedWithdrawalAddr = withdrawalAddr
		insertedWithdrawalAmount = amountGwei
		return nil
	}

	// Test fallback
	err := s.tk.Keeper.FallbackDepositCollateral(s.tk.Ctx, event)
	s.Require().NoError(err)

	// Verify withdrawal was inserted
	s.Require().Equal(event.CollateralOwner, insertedWithdrawalAddr)
	s.Require().Equal(event.AmountGwei.Uint64(), insertedWithdrawalAmount)
}

func (s *EventProcessingTestSuite) Test_ProcessWithdrawCollateral() {
	// Setup test parameters
	s.tk.SetupDefaultTestParams()

	// Register a validator first with 3 MITO collateral
	initialValidator := s.tk.RegisterTestValidator(math.NewUint(3000000000), math.ZeroUint(), false) // 3 MITO
	valAddr := initialValidator.Addr
	_, _, receiverAddr := testutil.GenerateSecp256k1Key()

	// Create withdraw collateral event
	now := time.Now().Unix()
	event := &bindings.ConsensusValidatorEntrypointMsgWithdrawCollateral{
		ValAddr:         common.BytesToAddress(valAddr.Bytes()),
		CollateralOwner: common.BytesToAddress(valAddr.Bytes()), // Using validator as owner since RegisterTestValidator uses the validator as owner
		AmountGwei:      big.NewInt(1000000000),                 // 1 MITO
		Receiver:        common.BytesToAddress(receiverAddr.Bytes()),
		MaturesAt:       big.NewInt(now + 86400), // 1 day from now
	}

	// Process the withdraw collateral event
	err, ignore := s.tk.Keeper.ProcessWithdrawCollateral(s.tk.Ctx, event)
	s.Require().NoError(err)
	s.Require().False(ignore)

	// Verify validator's collateral was reduced
	updatedValidator, found := s.tk.Keeper.GetValidator(s.tk.Ctx, valAddr)
	s.Require().True(found)
	expectedValidator := initialValidator
	expectedValidator.Collateral = initialValidator.Collateral.Sub(math.NewUintFromBigInt(event.AmountGwei))
	expectedValidator.CollateralShares = initialValidator.CollateralShares.Sub(types.CalculateCollateralSharesForWithdrawal(initialValidator.Collateral, initialValidator.CollateralShares, math.NewUintFromBigInt(event.AmountGwei)))
	expectedValidator.VotingPower = int64(2)
	s.Require().Equal(expectedValidator, updatedValidator)

	// Verify withdrawal was created
	expectedWithdrawal := types.Withdrawal{
		ID:             1,
		ValAddr:        valAddr,
		Amount:         event.AmountGwei.Uint64(),
		Receiver:       receiverAddr,
		MaturesAt:      event.MaturesAt.Int64(),
		CreationHeight: s.tk.Ctx.BlockHeight(),
	}
	found = false
	s.tk.Keeper.IterateWithdrawalsForValidator(s.tk.Ctx, valAddr, func(w types.Withdrawal) bool {
		if w == expectedWithdrawal {
			found = true
			return true
		}
		return false
	})
	s.Require().True(found)

	// Try withdrawing more than available collateral (should fail with ignore=true)
	invalidEvent := &bindings.ConsensusValidatorEntrypointMsgWithdrawCollateral{
		ValAddr:         common.BytesToAddress(valAddr.Bytes()),
		CollateralOwner: common.BytesToAddress(valAddr.Bytes()),
		AmountGwei:      big.NewInt(3000000000), // 3 MITO (more than remaining collateral)
		Receiver:        common.BytesToAddress(receiverAddr.Bytes()),
		MaturesAt:       big.NewInt(now + 86400),
	}

	err, ignore = s.tk.Keeper.ProcessWithdrawCollateral(s.tk.Ctx, invalidEvent)
	s.Require().Error(err)
	s.Require().True(ignore)
	s.Require().Contains(err.Error(), "failed to withdraw collateral")
}

func (s *EventProcessingTestSuite) Test_ProcessUnjail() {
	// Setup test parameters
	s.tk.SetupDefaultTestParams()

	// Register a validator
	validator := s.tk.RegisterTestValidator(math.NewUint(1000000000), math.ZeroUint(), false) // 1 MITO
	valAddr := validator.Addr

	// Jail validator directly
	s.tk.Keeper.Jail_(s.tk.Ctx, &validator, "test")

	// Create unjail event
	event := &bindings.ConsensusValidatorEntrypointMsgUnjail{
		ValAddr: common.BytesToAddress(valAddr.Bytes()),
	}

	// Mock the slashing keeper's UnjailFromConsAddr function
	s.tk.MockSlash.UnjailFromConsAddrFn = func(ctx context.Context, consAddr sdk.ConsAddress) error {
		return s.tk.Keeper.Unjail(ctx, consAddr)
	}

	// Process unjail event
	err, ignore := s.tk.Keeper.ProcessUnjail(s.tk.Ctx, event)
	s.Require().NoError(err)
	s.Require().False(ignore)

	// Verify validator is now unjailed
	updatedValidator, found := s.tk.Keeper.GetValidator(s.tk.Ctx, valAddr)
	s.Require().True(found)
	expectedValidator := validator
	expectedValidator.Jailed = false
	s.Require().Equal(expectedValidator, updatedValidator)

	// Try unjailing an already unjailed validator (should return error with ignore=true)
	err, ignore = s.tk.Keeper.ProcessUnjail(s.tk.Ctx, event)
	s.Require().Error(err)
	s.Require().True(ignore)
	s.Require().Contains(err.Error(), "validator is not jailed")
}

func (s *EventProcessingTestSuite) Test_ProcessUpdateExtraVotingPower() {
	// Setup test parameters
	s.tk.SetupDefaultTestParams()

	// Register a validator
	initialValidator := s.tk.RegisterTestValidator(math.NewUint(1000000000), math.ZeroUint(), false) // 1 MITO
	valAddr := initialValidator.Addr

	// Verify initial state
	s.Require().Equal(math.ZeroUint(), initialValidator.ExtraVotingPower)
	s.Require().Equal(int64(1), initialValidator.VotingPower) // 1 MITO = 1 voting power

	// Create update extra voting power event (1 MITO = 1e18 wei / 1e9 = 1e9 gwei)
	extraPowerWei := big.NewInt(0).Mul(big.NewInt(1e9), big.NewInt(1e9))
	event := &bindings.ConsensusValidatorEntrypointMsgUpdateExtraVotingPower{
		ValAddr:             common.BytesToAddress(valAddr.Bytes()),
		ExtraVotingPowerWei: extraPowerWei, // 1 MITO in wei
	}

	// Process update extra voting power event
	err, ignore := s.tk.Keeper.ProcessUpdateExtraVotingPower(s.tk.Ctx, event)
	s.Require().NoError(err)
	s.Require().False(ignore)

	// Verify validator's extra voting power and total voting power were updated
	updatedValidator, found := s.tk.Keeper.GetValidator(s.tk.Ctx, valAddr)
	s.Require().True(found)
	expectedValidator := initialValidator
	expectedValidator.ExtraVotingPower = math.NewUint(1e9)
	expectedValidator.VotingPower = 2
	s.Require().Equal(expectedValidator, updatedValidator)

	// Try updating a non-existent validator (should return error with ignore=true)
	_, _, nonExistentAddr := testutil.GenerateSecp256k1Key()
	invalidEvent := &bindings.ConsensusValidatorEntrypointMsgUpdateExtraVotingPower{
		ValAddr:             common.BytesToAddress(nonExistentAddr.Bytes()),
		ExtraVotingPowerWei: extraPowerWei,
	}

	err, ignore = s.tk.Keeper.ProcessUpdateExtraVotingPower(s.tk.Ctx, invalidEvent)
	s.Require().Error(err)
	s.Require().True(ignore)
	s.Require().ErrorIs(err, types.ErrValidatorNotFound)
}
