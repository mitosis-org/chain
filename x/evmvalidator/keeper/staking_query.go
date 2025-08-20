package keeper

import (
	"context"
	"time"

	"cosmossdk.io/math"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
	mitotypes "github.com/mitosis-org/chain/types"
	"github.com/mitosis-org/chain/x/evmvalidator/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var _ stakingtypes.QueryServer = StakingQueryServer{}

// StakingQueryServer implements the staking QueryServer interface for compatibility
type StakingQueryServer struct {
	k *Keeper
}

// NewStakingQueryServer creates a new StakingQueryServer instance
func NewStakingQueryServer(keeper *Keeper) stakingtypes.QueryServer {
	return &StakingQueryServer{k: keeper}
}

// Validators returns all validators, compatible with staking module
func (q StakingQueryServer) Validators(ctx context.Context, req *stakingtypes.QueryValidatorsRequest) (*stakingtypes.QueryValidatorsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	sdkCtx := sdk.UnwrapSDKContext(ctx)

	// Get all evmvalidator validators
	var stakingValidators []stakingtypes.Validator
	var index int64

	q.k.IterateValidators_(sdkCtx, func(_ int64, validator types.Validator) (stop bool) {
		// Convert evmvalidator to staking validator format
		stakingVal := q.convertToStakingValidator(sdkCtx, validator)
		
		// Apply status filter if specified
		if req.Status != "" {
			if stakingVal.Status.String() != req.Status {
				return false
			}
		}

		stakingValidators = append(stakingValidators, stakingVal)
		index++
		return false
	})

	// Apply pagination manually since we converted the data
	var pageRes *query.PageResponse
	if req.Pagination != nil {
		offset := int(req.Pagination.Offset)
		limit := int(req.Pagination.Limit)
		if limit == 0 {
			limit = 100 // default limit
		}
		if offset < len(stakingValidators) {
			end := offset + limit
			if end > len(stakingValidators) {
				end = len(stakingValidators)
			}
			stakingValidators = stakingValidators[offset:end]
		} else {
			stakingValidators = []stakingtypes.Validator{}
		}
		pageRes = &query.PageResponse{
			Total: uint64(index),
		}
	} else {
		pageRes = &query.PageResponse{
			Total: uint64(index),
		}
	}

	return &stakingtypes.QueryValidatorsResponse{
		Validators: stakingValidators,
		Pagination: pageRes,
	}, nil
}

// Validator returns a validator by operator address
func (q StakingQueryServer) Validator(ctx context.Context, req *stakingtypes.QueryValidatorRequest) (*stakingtypes.QueryValidatorResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	if req.ValidatorAddr == "" {
		return nil, status.Error(codes.InvalidArgument, "validator address cannot be empty")
	}

	sdkCtx := sdk.UnwrapSDKContext(ctx)

	// Convert validator address from bech32 to ethereum address
	valAddr, err := sdk.ValAddressFromBech32(req.ValidatorAddr)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid validator address: %s", err.Error())
	}

	// Convert to ethereum address (assuming validator address maps to ethereum address)
	ethAddr := mitotypes.EthAddress(valAddr.Bytes())
	
	validator, found := q.k.GetValidator(sdkCtx, ethAddr)
	if !found {
		return nil, status.Errorf(codes.NotFound, "validator %s not found", req.ValidatorAddr)
	}

	stakingVal := q.convertToStakingValidator(sdkCtx, validator)

	return &stakingtypes.QueryValidatorResponse{
		Validator: stakingVal,
	}, nil
}

// convertToStakingValidator converts evmvalidator.Validator to staking.Validator
func (q StakingQueryServer) convertToStakingValidator(ctx sdk.Context, evmVal types.Validator) stakingtypes.Validator {
	// Convert status
	var status stakingtypes.BondStatus
	switch {
	case evmVal.Jailed:
		status = stakingtypes.Unbonded
	case evmVal.Bonded:
		status = stakingtypes.Bonded
	default:
		status = stakingtypes.Unbonding
	}

	// Convert validator address to bech32 format
	valAddr := sdk.ValAddress(evmVal.Addr.Bytes())

	// Convert pubkey to consensus pubkey format
	var consensusPubkey *codectypes.Any
	if len(evmVal.Pubkey) > 0 {
		// Create secp256k1 pubkey from raw bytes
		pubkey := &secp256k1.PubKey{Key: evmVal.Pubkey}
		// Convert to Any type
		pubkeyAny, err := codectypes.NewAnyWithValue(pubkey)
		if err == nil {
			consensusPubkey = pubkeyAny
		}
	}

	return stakingtypes.Validator{
		OperatorAddress: valAddr.String(),
		ConsensusPubkey: consensusPubkey,
		Jailed:          evmVal.Jailed,
		Status:          status,
		Tokens:          math.NewIntFromBigInt(evmVal.Collateral.BigInt()),
		DelegatorShares: math.LegacyNewDecFromBigInt(evmVal.CollateralShares.BigInt()),
		Description: stakingtypes.Description{
			Moniker:         "validator", // evmvalidator doesn't have description fields
			Identity:        "",
			Website:         "",
			SecurityContact: "",
			Details:         "",
		},
		UnbondingHeight: 0, // Not applicable for EVM validators
		UnbondingTime:   ctx.BlockTime(), // Use current block time
		Commission: stakingtypes.Commission{
			CommissionRates: stakingtypes.CommissionRates{
				Rate:          math.LegacyZeroDec(), // evmvalidator doesn't have commission
				MaxRate:       math.LegacyOneDec(), // Default max rate
				MaxChangeRate: math.LegacyOneDec(), // Default max change rate
			},
			UpdateTime: ctx.BlockTime(),
		},
		MinSelfDelegation: math.ZeroInt(),
	}
}

// ValidatorDelegations returns delegations to a validator (not implemented for EVM validators)
func (q StakingQueryServer) ValidatorDelegations(ctx context.Context, req *stakingtypes.QueryValidatorDelegationsRequest) (*stakingtypes.QueryValidatorDelegationsResponse, error) {
	return &stakingtypes.QueryValidatorDelegationsResponse{
		DelegationResponses: []stakingtypes.DelegationResponse{},
		Pagination:          &query.PageResponse{},
	}, nil
}

// ValidatorUnbondingDelegations returns unbonding delegations from a validator (not implemented for EVM validators)
func (q StakingQueryServer) ValidatorUnbondingDelegations(ctx context.Context, req *stakingtypes.QueryValidatorUnbondingDelegationsRequest) (*stakingtypes.QueryValidatorUnbondingDelegationsResponse, error) {
	return &stakingtypes.QueryValidatorUnbondingDelegationsResponse{
		UnbondingResponses: []stakingtypes.UnbondingDelegation{},
		Pagination:         &query.PageResponse{},
	}, nil
}

// Delegation returns a delegation (not implemented for EVM validators)
func (q StakingQueryServer) Delegation(ctx context.Context, req *stakingtypes.QueryDelegationRequest) (*stakingtypes.QueryDelegationResponse, error) {
	return nil, status.Error(codes.Unimplemented, "delegation queries not supported for EVM validators")
}

// UnbondingDelegation returns an unbonding delegation (not implemented for EVM validators)
func (q StakingQueryServer) UnbondingDelegation(ctx context.Context, req *stakingtypes.QueryUnbondingDelegationRequest) (*stakingtypes.QueryUnbondingDelegationResponse, error) {
	return nil, status.Error(codes.Unimplemented, "unbonding delegation queries not supported for EVM validators")
}

// DelegatorDelegations returns all delegations of a delegator (not implemented for EVM validators)
func (q StakingQueryServer) DelegatorDelegations(ctx context.Context, req *stakingtypes.QueryDelegatorDelegationsRequest) (*stakingtypes.QueryDelegatorDelegationsResponse, error) {
	return &stakingtypes.QueryDelegatorDelegationsResponse{
		DelegationResponses: []stakingtypes.DelegationResponse{},
		Pagination:          &query.PageResponse{},
	}, nil
}

// DelegatorUnbondingDelegations returns all unbonding delegations of a delegator (not implemented for EVM validators)
func (q StakingQueryServer) DelegatorUnbondingDelegations(ctx context.Context, req *stakingtypes.QueryDelegatorUnbondingDelegationsRequest) (*stakingtypes.QueryDelegatorUnbondingDelegationsResponse, error) {
	return &stakingtypes.QueryDelegatorUnbondingDelegationsResponse{
		UnbondingResponses: []stakingtypes.UnbondingDelegation{},
		Pagination:         &query.PageResponse{},
	}, nil
}

// Redelegations returns redelegations (not implemented for EVM validators)
func (q StakingQueryServer) Redelegations(ctx context.Context, req *stakingtypes.QueryRedelegationsRequest) (*stakingtypes.QueryRedelegationsResponse, error) {
	return &stakingtypes.QueryRedelegationsResponse{
		RedelegationResponses: []stakingtypes.RedelegationResponse{},
		Pagination:            &query.PageResponse{},
	}, nil
}

// DelegatorValidators returns validators that a delegator is bonded to (not implemented for EVM validators)
func (q StakingQueryServer) DelegatorValidators(ctx context.Context, req *stakingtypes.QueryDelegatorValidatorsRequest) (*stakingtypes.QueryDelegatorValidatorsResponse, error) {
	return &stakingtypes.QueryDelegatorValidatorsResponse{
		Validators: []stakingtypes.Validator{},
		Pagination: &query.PageResponse{},
	}, nil
}

// DelegatorValidator returns a validator that a delegator is bonded to (not implemented for EVM validators)
func (q StakingQueryServer) DelegatorValidator(ctx context.Context, req *stakingtypes.QueryDelegatorValidatorRequest) (*stakingtypes.QueryDelegatorValidatorResponse, error) {
	return nil, status.Error(codes.Unimplemented, "delegator validator queries not supported for EVM validators")
}

// Pool returns the current staking pool (approximated from EVM validators)
func (q StakingQueryServer) Pool(ctx context.Context, req *stakingtypes.QueryPoolRequest) (*stakingtypes.QueryPoolResponse, error) {
	sdkCtx := sdk.UnwrapSDKContext(ctx)

	var bondedTokens, notBondedTokens math.Int
	bondedTokens = math.ZeroInt()
	notBondedTokens = math.ZeroInt()

	q.k.IterateValidators_(sdkCtx, func(_ int64, validator types.Validator) (stop bool) {
		tokens := math.NewIntFromBigInt(validator.Collateral.BigInt())
		if validator.Bonded && !validator.Jailed {
			bondedTokens = bondedTokens.Add(tokens)
		} else {
			notBondedTokens = notBondedTokens.Add(tokens)
		}
		return false
	})

	return &stakingtypes.QueryPoolResponse{
		Pool: stakingtypes.Pool{
			BondedTokens:    bondedTokens,
			NotBondedTokens: notBondedTokens,
		},
	}, nil
}

// Params returns the staking module parameters (adapted from evmvalidator params)
func (q StakingQueryServer) Params(ctx context.Context, req *stakingtypes.QueryParamsRequest) (*stakingtypes.QueryParamsResponse, error) {
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	evmParams := q.k.GetParams(sdkCtx)

	return &stakingtypes.QueryParamsResponse{
		Params: stakingtypes.Params{
			UnbondingTime:     time.Hour * 24 * 21, // Default 21 days unbonding period
			MaxValidators:     evmParams.MaxValidators,
			MaxEntries:        7, // Default value
			HistoricalEntries: 10000, // Default value
			BondDenom:         "stake", // Default denom
			MinCommissionRate: math.LegacyZeroDec(),
		},
	}, nil
}

// HistoricalInfo returns historical info (not implemented for EVM validators)
func (q StakingQueryServer) HistoricalInfo(ctx context.Context, req *stakingtypes.QueryHistoricalInfoRequest) (*stakingtypes.QueryHistoricalInfoResponse, error) {
	return nil, status.Error(codes.Unimplemented, "historical info queries not supported for EVM validators")
}