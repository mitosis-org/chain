package keeper_test

import (
	"sort"
	"testing"
	"time"

	"cosmossdk.io/math"
	"github.com/ethereum/go-ethereum/common"
	mitotypes "github.com/mitosis-org/chain/types"
	"github.com/mitosis-org/chain/x/evmvalidator/testutil"
	"github.com/mitosis-org/chain/x/evmvalidator/types"
	"github.com/stretchr/testify/suite"
)

// KeeperTestSuite is a test suite to be used with keeper tests
type KeeperTestSuite struct {
	suite.Suite
	tk testutil.TestKeeper
}

// SetupTest initializes the test suite
func (s *KeeperTestSuite) SetupTest() {
	s.tk = testutil.NewTestKeeper(&s.Suite)
}

// TestKeeperTestSuite runs the keeper test suite
func TestKeeperTestSuite(t *testing.T) {
	suite.Run(t, new(KeeperTestSuite))
}

// createTestWithdrawal returns a test withdrawal
func (s *KeeperTestSuite) createTestWithdrawal(
	valAddr mitotypes.EthAddress,
	amount uint64,
	maturesAt int64,
	id uint64,
) types.Withdrawal {
	return types.Withdrawal{
		ID:             id,
		ValAddr:        valAddr,
		Amount:         amount,
		Receiver:       valAddr, // Use the same address as receiver for simplicity
		MaturesAt:      maturesAt,
		CreationHeight: s.tk.Ctx.BlockHeight(),
	}
}

func (s *KeeperTestSuite) Test_GetParams() {
	params := s.tk.Keeper.GetParams(s.tk.Ctx)

	// Verify default params are returned
	defaultParams := types.DefaultParams()
	s.Require().Equal(defaultParams, params)
}

func (s *KeeperTestSuite) Test_SetParams() {
	defaultParams := types.DefaultParams()

	// Modify some params
	newParams := types.Params{
		MaxValidators:    defaultParams.MaxValidators + 1,
		MaxLeverageRatio: defaultParams.MaxLeverageRatio.Add(math.LegacyNewDec(1)),
		MinVotingPower:   defaultParams.MinVotingPower + 1,
		WithdrawalLimit:  defaultParams.WithdrawalLimit + 1,
	}

	// Set new params
	err := s.tk.Keeper.SetParams(s.tk.Ctx, newParams)
	s.Require().NoError(err)

	// Verify new params are returned
	params := s.tk.Keeper.GetParams(s.tk.Ctx)
	s.Require().Equal(newParams, params)

	// Test invalid params
	invalidParams := types.Params{
		MaxValidators:    0, // Invalid
		MaxLeverageRatio: defaultParams.MaxLeverageRatio,
		MinVotingPower:   defaultParams.MinVotingPower,
		WithdrawalLimit:  defaultParams.WithdrawalLimit,
	}

	err = s.tk.Keeper.SetParams(s.tk.Ctx, invalidParams)
	s.Require().Error(err)
}

func (s *KeeperTestSuite) Test_GetValidatorEntrypointContractAddr() {
	// Initially, address should be empty
	addr := s.tk.Keeper.GetValidatorEntrypointContractAddr(s.tk.Ctx)
	s.Require().Equal(mitotypes.EthAddress{}, addr)

	// Set a new address
	newAddr := mitotypes.EthAddress(common.HexToAddress("0x1234567890123456789012345678901234567890"))
	s.tk.Keeper.SetValidatorEntrypointContractAddr(s.tk.Ctx, newAddr)

	// Verify new address is returned
	addr = s.tk.Keeper.GetValidatorEntrypointContractAddr(s.tk.Ctx)
	s.Require().Equal(newAddr, addr)
}

func (s *KeeperTestSuite) Test_SetValidatorEntrypointContractAddr() {
	// Set a new address
	newAddr := mitotypes.EthAddress(common.HexToAddress("0x1234567890123456789012345678901234567890"))
	s.tk.Keeper.SetValidatorEntrypointContractAddr(s.tk.Ctx, newAddr)

	// Verify new address is returned
	addr := s.tk.Keeper.GetValidatorEntrypointContractAddr(s.tk.Ctx)
	s.Require().Equal(newAddr, addr)
}

func (s *KeeperTestSuite) Test_GetValidator() {
	// Generate validator data
	_, pubkey, ethAddr := testutil.GenerateSecp256k1Key()
	collateral := math.NewUint(1000000000)
	extraVotingPower := math.NewUint(500000000)

	// Test GetValidator when validator doesn't exist
	_, found := s.tk.Keeper.GetValidator(s.tk.Ctx, ethAddr)
	s.Require().False(found)

	// Create a validator
	validator := types.Validator{
		Addr:             ethAddr,
		Pubkey:           pubkey,
		Collateral:       collateral,
		ExtraVotingPower: extraVotingPower,
		VotingPower:      100,
		Jailed:           false,
		Bonded:           true,
	}

	// Set the validator
	s.tk.Keeper.SetValidator(s.tk.Ctx, validator)

	// Test GetValidator when validator exists
	gotValidator, found := s.tk.Keeper.GetValidator(s.tk.Ctx, ethAddr)
	s.Require().True(found)
	s.Require().Equal(validator, gotValidator)
}

func (s *KeeperTestSuite) Test_HasValidator() {
	// Generate validator data
	_, pubkey, ethAddr := testutil.GenerateSecp256k1Key()

	// Test HasValidator when validator doesn't exist
	has := s.tk.Keeper.HasValidator(s.tk.Ctx, ethAddr)
	s.Require().False(has)

	// Create a validator
	validator := types.Validator{
		Addr:             ethAddr,
		Pubkey:           pubkey,
		Collateral:       math.NewUint(1000000000),
		ExtraVotingPower: math.NewUint(0),
		VotingPower:      100,
		Jailed:           false,
		Bonded:           true,
	}

	// Set the validator
	s.tk.Keeper.SetValidator(s.tk.Ctx, validator)

	// Test HasValidator when validator exists
	has = s.tk.Keeper.HasValidator(s.tk.Ctx, ethAddr)
	s.Require().True(has)
}

func (s *KeeperTestSuite) Test_SetValidator() {
	// Generate validator data
	_, pubkey, ethAddr := testutil.GenerateSecp256k1Key()
	collateral := math.NewUint(1000000000)
	extraVotingPower := math.NewUint(500000000)

	// Create a validator
	validator := types.Validator{
		Addr:             ethAddr,
		Pubkey:           pubkey,
		Collateral:       collateral,
		ExtraVotingPower: extraVotingPower,
		VotingPower:      100,
		Jailed:           false,
		Bonded:           true,
	}

	// Set the validator
	s.tk.Keeper.SetValidator(s.tk.Ctx, validator)

	// Test GetValidator when validator exists
	gotValidator, found := s.tk.Keeper.GetValidator(s.tk.Ctx, ethAddr)
	s.Require().True(found)
	s.Require().Equal(validator, gotValidator)
}

func (s *KeeperTestSuite) Test_IterateValidators_() {
	// Setup validators
	// Generate validator 1
	_, pubkey1, ethAddr1 := testutil.GenerateSecp256k1Key()
	validator1 := types.Validator{
		Addr:             ethAddr1,
		Pubkey:           pubkey1,
		Collateral:       math.NewUint(5000000000),
		ExtraVotingPower: math.ZeroUint(),
		VotingPower:      5,
		Jailed:           false,
		Bonded:           true,
	}
	s.tk.Keeper.SetValidator(s.tk.Ctx, validator1)

	// Generate validator 2
	_, pubkey2, ethAddr2 := testutil.GenerateSecp256k1Key()
	validator2 := types.Validator{
		Addr:             ethAddr2,
		Pubkey:           pubkey2,
		Collateral:       math.NewUint(3000000000),
		ExtraVotingPower: math.ZeroUint(),
		VotingPower:      3,
		Jailed:           false,
		Bonded:           true,
	}
	s.tk.Keeper.SetValidator(s.tk.Ctx, validator2)

	// Generate validator 3
	_, pubkey3, ethAddr3 := testutil.GenerateSecp256k1Key()
	validator3 := types.Validator{
		Addr:             ethAddr3,
		Pubkey:           pubkey3,
		Collateral:       math.NewUint(2000000000),
		ExtraVotingPower: math.ZeroUint(),
		VotingPower:      2,
		Jailed:           true, // Jailed validator
		Bonded:           false,
	}
	s.tk.Keeper.SetValidator(s.tk.Ctx, validator3)

	// Test iteration with accumulation
	validators := make(map[string]types.Validator)
	s.tk.Keeper.IterateValidators_(s.tk.Ctx, func(index int64, validator types.Validator) (stop bool) {
		validators[validator.Addr.String()] = validator
		return false // Continue iteration
	})

	// Verify all validators were iterated
	s.Require().Equal(3, len(validators), "Should have iterated through all 3 validators")
	s.Require().Equal(validator1, validators[ethAddr1.String()], "Validator 1 should match")
	s.Require().Equal(validator2, validators[ethAddr2.String()], "Validator 2 should match")
	s.Require().Equal(validator3, validators[ethAddr3.String()], "Validator 3 should match")

	// Test early termination
	count := 0
	s.tk.Keeper.IterateValidators_(s.tk.Ctx, func(index int64, validator types.Validator) (stop bool) {
		count++
		return true // Stop after first item
	})
	s.Require().Equal(1, count, "Should have stopped after first validator")

	// Test index correctness
	indices := make([]int64, 0)
	s.tk.Keeper.IterateValidators_(s.tk.Ctx, func(index int64, validator types.Validator) (stop bool) {
		indices = append(indices, index)
		return false
	})

	// Verify indices are sequential
	s.Require().Equal(3, len(indices), "Should have 3 indices")
	s.Require().Equal(int64(0), indices[0], "First index should be 0")
	s.Require().Equal(int64(1), indices[1], "Second index should be 1")
	s.Require().Equal(int64(2), indices[2], "Third index should be 2")
}

func (s *KeeperTestSuite) Test_GetAllValidators() {
	// Initial state should have no validators
	initialValidators := s.tk.Keeper.GetAllValidators(s.tk.Ctx)
	s.Require().Equal(0, len(initialValidators), "Should start with no validators")

	// Setup validators
	// Generate validator 1
	_, pubkey1, ethAddr1 := testutil.GenerateSecp256k1Key()
	validator1 := types.Validator{
		Addr:             ethAddr1,
		Pubkey:           pubkey1,
		Collateral:       math.NewUint(5000000000),
		ExtraVotingPower: math.ZeroUint(),
		VotingPower:      5,
		Jailed:           false,
		Bonded:           true,
	}
	s.tk.Keeper.SetValidator(s.tk.Ctx, validator1)

	// Generate validator 2
	_, pubkey2, ethAddr2 := testutil.GenerateSecp256k1Key()
	validator2 := types.Validator{
		Addr:             ethAddr2,
		Pubkey:           pubkey2,
		Collateral:       math.NewUint(3000000000),
		ExtraVotingPower: math.ZeroUint(),
		VotingPower:      3,
		Jailed:           false,
		Bonded:           true,
	}
	s.tk.Keeper.SetValidator(s.tk.Ctx, validator2)

	// Generate validator 3 (jailed)
	_, pubkey3, ethAddr3 := testutil.GenerateSecp256k1Key()
	validator3 := types.Validator{
		Addr:             ethAddr3,
		Pubkey:           pubkey3,
		Collateral:       math.NewUint(2000000000),
		ExtraVotingPower: math.ZeroUint(),
		VotingPower:      2,
		Jailed:           true, // Jailed validator
		Bonded:           false,
	}
	s.tk.Keeper.SetValidator(s.tk.Ctx, validator3)

	// Get all validators
	allValidators := s.tk.Keeper.GetAllValidators(s.tk.Ctx)

	// Verify all validators are returned
	s.Require().Equal(3, len(allValidators), "Should return all 3 validators")

	// Create map for easier lookup
	validatorMap := make(map[string]types.Validator)
	for _, val := range allValidators {
		validatorMap[val.Addr.String()] = val
	}

	// Verify each validator
	s.Require().Equal(validator1, validatorMap[ethAddr1.String()], "Validator 1 should match")
	s.Require().Equal(validator2, validatorMap[ethAddr2.String()], "Validator 2 should match")
	s.Require().Equal(validator3, validatorMap[ethAddr3.String()], "Validator 3 should match")

	// Verify jailed validator is included (unlike GetNotJailedValidatorsByPower)
	s.Require().True(validatorMap[ethAddr3.String()].Jailed, "Jailed validator should be included")

	// Verify GetAllValidators returns what we would get by iterating
	var iteratedValidators []types.Validator
	s.tk.Keeper.IterateValidators_(s.tk.Ctx, func(_ int64, validator types.Validator) bool {
		iteratedValidators = append(iteratedValidators, validator)
		return false
	})

	// Sort both lists by address to ensure fair comparison
	sortByAddr := func(validators []types.Validator) {
		sort.Slice(validators, func(i, j int) bool {
			return validators[i].Addr.String() < validators[j].Addr.String()
		})
	}

	sortByAddr(allValidators)
	sortByAddr(iteratedValidators)

	s.Require().Equal(iteratedValidators, allValidators, "GetAllValidators should match iterating through validators")
}

func (s *KeeperTestSuite) Test_GetNotJailedValidatorsByPower() {
	// Set test parameters
	params := types.Params{
		MaxValidators:    10,
		MaxLeverageRatio: math.LegacyNewDec(10), // 10x leverage
		MinVotingPower:   1,
		WithdrawalLimit:  10,
	}
	err := s.tk.Keeper.SetParams(s.tk.Ctx, params)
	s.Require().NoError(err)

	// Register validators with different powers
	// Generate validator 1 with power 5
	_, pubkey1, ethAddr1 := testutil.GenerateSecp256k1Key()
	err = s.tk.Keeper.RegisterValidator(s.tk.Ctx, ethAddr1, pubkey1, math.NewUint(5000000000), math.ZeroUint(), false)
	s.Require().NoError(err)
	validator1, found := s.tk.Keeper.GetValidator(s.tk.Ctx, ethAddr1)
	s.Require().True(found)

	// Generate validator 2 with power 3
	_, pubkey2, ethAddr2 := testutil.GenerateSecp256k1Key()
	err = s.tk.Keeper.RegisterValidator(s.tk.Ctx, ethAddr2, pubkey2, math.NewUint(3000000000), math.ZeroUint(), false)
	s.Require().NoError(err)
	validator2, found := s.tk.Keeper.GetValidator(s.tk.Ctx, ethAddr2)
	s.Require().True(found)

	// Generate validator 3 with power 2
	_, pubkey3, ethAddr3 := testutil.GenerateSecp256k1Key()
	err = s.tk.Keeper.RegisterValidator(s.tk.Ctx, ethAddr3, pubkey3, math.NewUint(2000000000), math.ZeroUint(), false)
	s.Require().NoError(err)
	validator3, found := s.tk.Keeper.GetValidator(s.tk.Ctx, ethAddr3)
	s.Require().True(found)

	// Generate validator 4 with power 1, jailed
	_, pubkey4, ethAddr4 := testutil.GenerateSecp256k1Key()
	err = s.tk.Keeper.RegisterValidator(s.tk.Ctx, ethAddr4, pubkey4, math.NewUint(1000000000), math.ZeroUint(), true)
	s.Require().NoError(err)
	_, found = s.tk.Keeper.GetValidator(s.tk.Ctx, ethAddr4)
	s.Require().True(found)

	// Get not jailed validators with max of 10
	validators := s.tk.Keeper.GetNotJailedValidatorsByPower(s.tk.Ctx, 10)
	s.Require().Equal(3, len(validators))
	s.Require().Equal(validator1.Addr, validators[0].Addr)
	s.Require().Equal(validator2.Addr, validators[1].Addr)
	s.Require().Equal(validator3.Addr, validators[2].Addr)

	// Limit the max validators to 2
	limitedValidators := s.tk.Keeper.GetNotJailedValidatorsByPower(s.tk.Ctx, 2)
	s.Require().Equal(2, len(limitedValidators))
	s.Require().Equal(validator1.Addr, limitedValidators[0].Addr)
	s.Require().Equal(validator2.Addr, limitedValidators[1].Addr)

	// Delete validator 2 from power index
	s.tk.Keeper.DeleteValidatorByPowerIndex(s.tk.Ctx, validator2.VotingPower, validator2.Addr)

	// Get validators again
	updatedValidators := s.tk.Keeper.GetNotJailedValidatorsByPower(s.tk.Ctx, 10)
	s.Require().Equal(2, len(updatedValidators))
	s.Require().Equal(validator1.Addr, updatedValidators[0].Addr)
	s.Require().Equal(validator3.Addr, updatedValidators[1].Addr)
}

func (s *KeeperTestSuite) Test_GetValidatorByConsAddr() {
	// Generate validator data
	_, pubkey, ethAddr := testutil.GenerateSecp256k1Key()
	collateral := math.NewUint(1000000000)
	extraVotingPower := math.NewUint(500000000)

	// Create a validator
	validator := types.Validator{
		Addr:             ethAddr,
		Pubkey:           pubkey,
		Collateral:       collateral,
		ExtraVotingPower: extraVotingPower,
		VotingPower:      100,
		Jailed:           false,
		Bonded:           true,
	}

	// Set the validator
	s.tk.Keeper.SetValidator(s.tk.Ctx, validator)

	// Get consensus address
	consAddr := validator.MustConsAddr()

	// Before setting validator by consensus address
	_, found := s.tk.Keeper.GetValidatorByConsAddr(s.tk.Ctx, consAddr)
	s.Require().False(found)

	// Set validator by consensus address
	s.tk.Keeper.SetValidatorByConsAddr(s.tk.Ctx, consAddr, ethAddr)

	// Get validator by consensus address
	gotValidator, found := s.tk.Keeper.GetValidatorByConsAddr(s.tk.Ctx, consAddr)
	s.Require().True(found)
	s.Require().Equal(validator, gotValidator)
}

func (s *KeeperTestSuite) Test_SetValidatorByConsAddr() {
	// Generate validator data
	_, pubkey, ethAddr := testutil.GenerateSecp256k1Key()

	// Create a validator
	validator := types.Validator{
		Addr:             ethAddr,
		Pubkey:           pubkey,
		Collateral:       math.NewUint(1000000000),
		ExtraVotingPower: math.NewUint(0),
		VotingPower:      100,
		Jailed:           false,
		Bonded:           true,
	}

	// Set the validator
	s.tk.Keeper.SetValidator(s.tk.Ctx, validator)

	// Get consensus address
	consAddr := validator.MustConsAddr()

	// Set validator by consensus address
	s.tk.Keeper.SetValidatorByConsAddr(s.tk.Ctx, consAddr, ethAddr)

	// Get validator by consensus address
	gotValidator, found := s.tk.Keeper.GetValidatorByConsAddr(s.tk.Ctx, consAddr)
	s.Require().True(found)
	s.Require().Equal(validator, gotValidator)
}

func (s *KeeperTestSuite) Test_GetValidatorsByPowerIndexIterator() {
	// Generate validator data
	_, pubkey1, ethAddr1 := testutil.GenerateSecp256k1Key()
	_, pubkey2, ethAddr2 := testutil.GenerateSecp256k1Key()

	// Create validators
	validator1 := types.Validator{
		Addr:             ethAddr1,
		Pubkey:           pubkey1,
		Collateral:       math.NewUint(1000000000),
		ExtraVotingPower: math.NewUint(0),
		VotingPower:      100,
		Jailed:           false,
		Bonded:           true,
	}

	validator2 := types.Validator{
		Addr:             ethAddr2,
		Pubkey:           pubkey2,
		Collateral:       math.NewUint(2000000000),
		ExtraVotingPower: math.NewUint(0),
		VotingPower:      200,
		Jailed:           false,
		Bonded:           true,
	}

	// Set validators
	s.tk.Keeper.SetValidator(s.tk.Ctx, validator1)
	s.tk.Keeper.SetValidator(s.tk.Ctx, validator2)

	// Set validators by power index
	s.tk.Keeper.SetValidatorByPowerIndex(s.tk.Ctx, validator1.VotingPower, validator1.Addr)
	s.tk.Keeper.SetValidatorByPowerIndex(s.tk.Ctx, validator2.VotingPower, validator2.Addr)

	// Get validator iterator by power index
	iterator := s.tk.Keeper.GetValidatorsByPowerIndexIterator(s.tk.Ctx)
	defer iterator.Close()

	// Count validators and check they exist
	count := 0
	for ; iterator.Valid(); iterator.Next() {
		count++
	}

	s.Require().Equal(2, count, "Expected 2 validators in power index")
}

func (s *KeeperTestSuite) Test_SetValidatorByPowerIndex() {
	// Generate validator data
	_, pubkey, ethAddr := testutil.GenerateSecp256k1Key()
	votingPower := int64(100)

	// Create a validator
	validator := types.Validator{
		Addr:             ethAddr,
		Pubkey:           pubkey,
		Collateral:       math.NewUint(1000000000),
		ExtraVotingPower: math.NewUint(0),
		VotingPower:      votingPower,
		Jailed:           false,
		Bonded:           true,
	}

	// Set the validator
	s.tk.Keeper.SetValidator(s.tk.Ctx, validator)

	// Set validator by power index
	s.tk.Keeper.SetValidatorByPowerIndex(s.tk.Ctx, votingPower, ethAddr)

	// Get validator iterator by power index
	iterator := s.tk.Keeper.GetValidatorsByPowerIndexIterator(s.tk.Ctx)
	defer iterator.Close()

	// Check if validator exists in power index
	found := false
	for ; iterator.Valid(); iterator.Next() {
		// Get the validator address from the iterator value
		value := iterator.Value()
		addr := mitotypes.BytesToEthAddress(value)

		if addr.String() == ethAddr.String() {
			found = true
			break
		}
	}

	s.Require().True(found, "Validator should exist in power index")
}

func (s *KeeperTestSuite) Test_DeleteValidatorByPowerIndex() {
	// Generate validator data
	_, pubkey, ethAddr := testutil.GenerateSecp256k1Key()
	votingPower := int64(100)

	// Create a validator
	validator := types.Validator{
		Addr:             ethAddr,
		Pubkey:           pubkey,
		Collateral:       math.NewUint(1000000000),
		ExtraVotingPower: math.NewUint(0),
		VotingPower:      votingPower,
		Jailed:           false,
		Bonded:           true,
	}

	// Set the validator
	s.tk.Keeper.SetValidator(s.tk.Ctx, validator)

	// Set validator by power index
	s.tk.Keeper.SetValidatorByPowerIndex(s.tk.Ctx, votingPower, ethAddr)

	// Delete validator from power index
	s.tk.Keeper.DeleteValidatorByPowerIndex(s.tk.Ctx, votingPower, ethAddr)

	// Get validator iterator by power index
	iterator := s.tk.Keeper.GetValidatorsByPowerIndexIterator(s.tk.Ctx)
	defer iterator.Close()

	// Check if validator exists in power index
	found := false
	for ; iterator.Valid(); iterator.Next() {
		// Get the validator address from the iterator value
		value := iterator.Value()
		addr := mitotypes.BytesToEthAddress(value)

		if addr.String() == ethAddr.String() {
			found = true
			break
		}
	}

	s.Require().False(found, "Validator should not exist in power index after deletion")
}

func (s *KeeperTestSuite) Test_GetLastValidatorPower() {
	// Generate validator address
	_, _, ethAddr := testutil.GenerateSecp256k1Key()

	// Test GetLastValidatorPower when power doesn't exist
	_, found := s.tk.Keeper.GetLastValidatorPower(s.tk.Ctx, ethAddr)
	s.Require().False(found)

	// Set validator power
	power := int64(1000)
	s.tk.Keeper.SetLastValidatorPower(s.tk.Ctx, ethAddr, power)

	// Test GetLastValidatorPower when power exists
	gotPower, found := s.tk.Keeper.GetLastValidatorPower(s.tk.Ctx, ethAddr)
	s.Require().True(found)
	s.Require().Equal(power, gotPower)
}

func (s *KeeperTestSuite) Test_SetLastValidatorPower() {
	// Generate validator address
	_, _, ethAddr := testutil.GenerateSecp256k1Key()

	// Set validator power
	power := int64(1000)
	s.tk.Keeper.SetLastValidatorPower(s.tk.Ctx, ethAddr, power)

	// Verify power was set correctly
	gotPower, found := s.tk.Keeper.GetLastValidatorPower(s.tk.Ctx, ethAddr)
	s.Require().True(found)
	s.Require().Equal(power, gotPower)
}

func (s *KeeperTestSuite) Test_DeleteLastValidatorPower() {
	// Generate validator address
	_, _, ethAddr := testutil.GenerateSecp256k1Key()

	// Set validator power
	power := int64(1000)
	s.tk.Keeper.SetLastValidatorPower(s.tk.Ctx, ethAddr, power)

	// Verify power exists
	_, found := s.tk.Keeper.GetLastValidatorPower(s.tk.Ctx, ethAddr)
	s.Require().True(found)

	// Delete validator power
	s.tk.Keeper.DeleteLastValidatorPower(s.tk.Ctx, ethAddr)

	// Test GetLastValidatorPower after deletion
	_, found = s.tk.Keeper.GetLastValidatorPower(s.tk.Ctx, ethAddr)
	s.Require().False(found)
}

func (s *KeeperTestSuite) Test_IterateLastValidatorPowers() {
	// Generate validator addresses
	_, _, ethAddr1 := testutil.GenerateSecp256k1Key()
	_, _, ethAddr2 := testutil.GenerateSecp256k1Key()
	_, _, ethAddr3 := testutil.GenerateSecp256k1Key()

	// Set validator powers
	s.tk.Keeper.SetLastValidatorPower(s.tk.Ctx, ethAddr1, 100)
	s.tk.Keeper.SetLastValidatorPower(s.tk.Ctx, ethAddr2, 200)
	s.tk.Keeper.SetLastValidatorPower(s.tk.Ctx, ethAddr3, 300)

	// Test IterateLastValidatorPowers
	powers := make(map[string]int64)
	s.tk.Keeper.IterateLastValidatorPowers(s.tk.Ctx, func(valAddr mitotypes.EthAddress, power int64) bool {
		powers[valAddr.String()] = power
		return false
	})

	s.Require().Equal(3, len(powers))
	s.Require().Equal(int64(100), powers[ethAddr1.String()])
	s.Require().Equal(int64(200), powers[ethAddr2.String()])
	s.Require().Equal(int64(300), powers[ethAddr3.String()])

	// Test early termination
	count := 0
	s.tk.Keeper.IterateLastValidatorPowers(s.tk.Ctx, func(valAddr mitotypes.EthAddress, power int64) bool {
		count++
		return true // Stop after first item
	})

	s.Require().Equal(1, count)
}

func (s *KeeperTestSuite) Test_IterateLastValidators() {
	// Generate validator data
	_, pubkey1, ethAddr1 := testutil.GenerateSecp256k1Key()
	_, pubkey2, ethAddr2 := testutil.GenerateSecp256k1Key()

	// Create validators
	validator1 := types.Validator{
		Addr:             ethAddr1,
		Pubkey:           pubkey1,
		Collateral:       math.NewUint(1000000000),
		ExtraVotingPower: math.NewUint(0),
		VotingPower:      100,
		Jailed:           false,
		Bonded:           true,
	}

	validator2 := types.Validator{
		Addr:             ethAddr2,
		Pubkey:           pubkey2,
		Collateral:       math.NewUint(2000000000),
		ExtraVotingPower: math.NewUint(0),
		VotingPower:      200,
		Jailed:           false,
		Bonded:           true,
	}

	// Set validators
	s.tk.Keeper.SetValidator(s.tk.Ctx, validator1)
	s.tk.Keeper.SetValidator(s.tk.Ctx, validator2)

	// Set validator powers
	s.tk.Keeper.SetLastValidatorPower(s.tk.Ctx, ethAddr1, validator1.VotingPower)
	s.tk.Keeper.SetLastValidatorPower(s.tk.Ctx, ethAddr2, validator2.VotingPower)

	// Test IterateLastValidators
	validators := make(map[string]types.Validator)

	err := s.tk.Keeper.IterateLastValidators(s.tk.Ctx, func(index int64, validator types.Validator) bool {
		validators[validator.Addr.String()] = validator
		return false
	})

	s.Require().NoError(err)
	s.Require().Equal(2, len(validators))
	s.Require().Equal(validator1, validators[ethAddr1.String()])
	s.Require().Equal(validator2, validators[ethAddr2.String()])

	// Test early termination
	count := 0
	err = s.tk.Keeper.IterateLastValidators(s.tk.Ctx, func(index int64, validator types.Validator) bool {
		count++
		return true // Stop after first item
	})

	s.Require().NoError(err)
	s.Require().Equal(1, count)
}

func (s *KeeperTestSuite) Test_GetLastValidatorPowers() {
	// Generate validator addresses
	_, _, ethAddr1 := testutil.GenerateSecp256k1Key()
	_, _, ethAddr2 := testutil.GenerateSecp256k1Key()

	// Set validator powers
	s.tk.Keeper.SetLastValidatorPower(s.tk.Ctx, ethAddr1, 100)
	s.tk.Keeper.SetLastValidatorPower(s.tk.Ctx, ethAddr2, 200)

	// Test GetLastValidatorPowers
	powers := s.tk.Keeper.GetLastValidatorPowers(s.tk.Ctx)

	// Create map for easy lookup
	powerMap := make(map[string]int64)
	for _, power := range powers {
		powerMap[power.ValAddr.String()] = power.Power
	}

	s.Require().Equal(2, len(powers))
	s.Require().Equal(int64(100), powerMap[ethAddr1.String()])
	s.Require().Equal(int64(200), powerMap[ethAddr2.String()])
}

func (s *KeeperTestSuite) Test_GetWithdrawalLastID() {
	// Initially, ID should be 0
	id := s.tk.Keeper.GetWithdrawalLastID(s.tk.Ctx)
	s.Require().Equal(uint64(0), id)
}

func (s *KeeperTestSuite) Test_SetWithdrawalLastID() {
	// Set a new ID
	newID := uint64(1234)
	s.tk.Keeper.SetWithdrawalLastID(s.tk.Ctx, newID)

	// Verify new ID is returned
	id := s.tk.Keeper.GetWithdrawalLastID(s.tk.Ctx)
	s.Require().Equal(newID, id)
}

func (s *KeeperTestSuite) Test_SetWithdrawal() {
	// Generate validator address
	_, _, valAddr := testutil.GenerateSecp256k1Key()

	// Create withdrawal
	withdrawal := types.Withdrawal{
		ID:             1,
		ValAddr:        valAddr,
		Amount:         1000000000,
		Receiver:       valAddr,
		MaturesAt:      1000,
		CreationHeight: 100,
	}

	// Set withdrawal
	s.tk.Keeper.SetWithdrawal(s.tk.Ctx, withdrawal)

	// Verify withdrawal was stored by iterating through all withdrawals
	var found bool
	s.tk.Keeper.IterateWithdrawalsByMaturesAt(s.tk.Ctx, func(w types.Withdrawal) bool {
		if w.ID == withdrawal.ID {
			found = true
			s.Require().Equal(withdrawal.ValAddr, w.ValAddr)
			s.Require().Equal(withdrawal.Amount, w.Amount)
			s.Require().Equal(withdrawal.Receiver, w.Receiver)
			s.Require().Equal(withdrawal.MaturesAt, w.MaturesAt)
			s.Require().Equal(withdrawal.CreationHeight, w.CreationHeight)
			return true
		}
		return false
	})
	s.Require().True(found, "withdrawal not found")
}

func (s *KeeperTestSuite) Test_AddNewWithdrawalWithNextID() {
	// Setup
	_, _, valAddr := testutil.GenerateSecp256k1Key()

	// Create withdrawal without ID
	withdrawal := s.createTestWithdrawal(valAddr, 1000000000, time.Now().Unix()+86400, 0)

	// Add withdrawal with next ID
	s.tk.Keeper.AddNewWithdrawalWithNextID(s.tk.Ctx, &withdrawal)

	// ID should be 1 (since we start from 0)
	s.Require().Equal(uint64(1), withdrawal.ID)

	// Last ID should be updated to 1
	lastID := s.tk.Keeper.GetWithdrawalLastID(s.tk.Ctx)
	s.Require().Equal(uint64(1), lastID)

	// Add another withdrawal
	withdrawal2 := s.createTestWithdrawal(valAddr, 2000000000, time.Now().Unix()+86400*2, 0)
	s.tk.Keeper.AddNewWithdrawalWithNextID(s.tk.Ctx, &withdrawal2)

	// ID should be 2
	s.Require().Equal(uint64(2), withdrawal2.ID)

	// Last ID should be updated to 2
	lastID = s.tk.Keeper.GetWithdrawalLastID(s.tk.Ctx)
	s.Require().Equal(uint64(2), lastID)
}

func (s *KeeperTestSuite) Test_DeleteWithdrawal() {
	// Setup
	_, _, valAddr := testutil.GenerateSecp256k1Key()

	// Create and add withdrawal
	withdrawal := s.createTestWithdrawal(valAddr, 1000000000, time.Now().Unix()+86400, 0)
	s.tk.Keeper.AddNewWithdrawalWithNextID(s.tk.Ctx, &withdrawal)

	// Verify withdrawal exists
	var foundBefore bool
	s.tk.Keeper.IterateWithdrawalsByMaturesAt(s.tk.Ctx, func(w types.Withdrawal) bool {
		if w.ID == withdrawal.ID {
			foundBefore = true
			return true
		}
		return false
	})
	s.Require().True(foundBefore, "withdrawal should exist before deletion")

	// Delete withdrawal
	s.tk.Keeper.DeleteWithdrawal(s.tk.Ctx, withdrawal)

	// Verify withdrawal was deleted
	var foundAfter bool
	s.tk.Keeper.IterateWithdrawalsByMaturesAt(s.tk.Ctx, func(w types.Withdrawal) bool {
		if w.ID == withdrawal.ID {
			foundAfter = true
			return true
		}
		return false
	})
	s.Require().False(foundAfter, "withdrawal should not exist after deletion")
}

func (s *KeeperTestSuite) Test_IterateWithdrawalsByMaturesAt() {
	// Setup
	_, _, valAddr := testutil.GenerateSecp256k1Key()

	// Create withdrawals with different maturity times
	now := time.Now().Unix()
	withdrawal1 := s.createTestWithdrawal(valAddr, 1000000000, now+86400, 0)   // 1 day from now
	withdrawal2 := s.createTestWithdrawal(valAddr, 2000000000, now+86400*2, 0) // 2 days from now
	withdrawal3 := s.createTestWithdrawal(valAddr, 3000000000, now+86400*3, 0) // 3 days from now

	// Add withdrawals
	s.tk.Keeper.AddNewWithdrawalWithNextID(s.tk.Ctx, &withdrawal1)
	s.tk.Keeper.AddNewWithdrawalWithNextID(s.tk.Ctx, &withdrawal2)
	s.tk.Keeper.AddNewWithdrawalWithNextID(s.tk.Ctx, &withdrawal3)

	// Iterate and collect withdrawals
	var withdrawals []types.Withdrawal
	s.tk.Keeper.IterateWithdrawalsByMaturesAt(s.tk.Ctx, func(w types.Withdrawal) bool {
		withdrawals = append(withdrawals, w)
		return false
	})

	// Should have 3 withdrawals
	s.Require().Equal(3, len(withdrawals))

	// Withdrawals should be sorted by maturesAt (ascending)
	s.Require().Equal(withdrawal1.ID, withdrawals[0].ID)
	s.Require().Equal(withdrawal2.ID, withdrawals[1].ID)
	s.Require().Equal(withdrawal3.ID, withdrawals[2].ID)
}

func (s *KeeperTestSuite) Test_GetAllWithdrawals() {
	// Setup
	_, _, valAddr := testutil.GenerateSecp256k1Key()

	// Create and add withdrawals
	now := time.Now().Unix()
	withdrawal1 := s.createTestWithdrawal(valAddr, 1000000000, now+86400, 0)
	withdrawal2 := s.createTestWithdrawal(valAddr, 2000000000, now+86400*2, 0)
	s.tk.Keeper.AddNewWithdrawalWithNextID(s.tk.Ctx, &withdrawal1)
	s.tk.Keeper.AddNewWithdrawalWithNextID(s.tk.Ctx, &withdrawal2)

	// Get all withdrawals
	withdrawals := s.tk.Keeper.GetAllWithdrawals(s.tk.Ctx)

	// Should have 2 withdrawals
	s.Require().Equal(2, len(withdrawals))

	// Withdrawals should be sorted by maturesAt
	s.Require().Equal(withdrawal1.ID, withdrawals[0].ID)
	s.Require().Equal(withdrawal2.ID, withdrawals[1].ID)
}

func (s *KeeperTestSuite) Test_IterateWithdrawalsForValidator() {
	// Setup
	_, _, valAddr1 := testutil.GenerateSecp256k1Key()
	_, _, valAddr2 := testutil.GenerateSecp256k1Key()

	// Create withdrawals for different validators
	now := time.Now().Unix()
	withdrawal1 := s.createTestWithdrawal(valAddr1, 1000000000, now+86400, 0)
	withdrawal2 := s.createTestWithdrawal(valAddr1, 2000000000, now+86400*2, 0)
	withdrawal3 := s.createTestWithdrawal(valAddr2, 3000000000, now+86400, 0)

	// Add withdrawals
	s.tk.Keeper.AddNewWithdrawalWithNextID(s.tk.Ctx, &withdrawal1)
	s.tk.Keeper.AddNewWithdrawalWithNextID(s.tk.Ctx, &withdrawal2)
	s.tk.Keeper.AddNewWithdrawalWithNextID(s.tk.Ctx, &withdrawal3)

	// Iterate for validator 1
	var withdrawalsVal1 []types.Withdrawal
	s.tk.Keeper.IterateWithdrawalsForValidator(s.tk.Ctx, valAddr1, func(w types.Withdrawal) bool {
		withdrawalsVal1 = append(withdrawalsVal1, w)
		return false
	})

	// Should have 2 withdrawals for validator 1
	s.Require().Equal(2, len(withdrawalsVal1))
	s.Require().Equal(withdrawal1.ID, withdrawalsVal1[0].ID)
	s.Require().Equal(withdrawal2.ID, withdrawalsVal1[1].ID)

	// Iterate for validator 2
	var withdrawalsVal2 []types.Withdrawal
	s.tk.Keeper.IterateWithdrawalsForValidator(s.tk.Ctx, valAddr2, func(w types.Withdrawal) bool {
		withdrawalsVal2 = append(withdrawalsVal2, w)
		return false
	})

	// Should have 1 withdrawal for validator 2
	s.Require().Equal(1, len(withdrawalsVal2))
	s.Require().Equal(withdrawal3.ID, withdrawalsVal2[0].ID)
}
