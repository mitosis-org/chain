package keeper

import (
	"context"
	mitotypes "github.com/mitosis-org/chain/types"

	"cosmossdk.io/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	"github.com/mitosis-org/chain/x/evmvalidator/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var _ types.QueryServer = QueryServer{}

// QueryServer implements the QueryServer interface for the evmvalidator module
type QueryServer struct {
	k Keeper
}

// NewQueryServer creates a new QueryServer instance
func NewQueryServer(keeper *Keeper) types.QueryServer {
	return &QueryServer{k: *keeper}
}

// Params returns the module's parameters
func (q QueryServer) Params(ctx context.Context, req *types.QueryParamsRequest) (*types.QueryParamsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	sdkCtx := sdk.UnwrapSDKContext(ctx)
	params := q.k.GetParams(sdkCtx)
	return &types.QueryParamsResponse{Params: params}, nil
}

// Validator returns a specific validator by validator address
func (q QueryServer) Validator(ctx context.Context, req *types.QueryValidatorRequest) (*types.QueryValidatorResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	if req.ValAddr == nil {
		return nil, status.Error(codes.InvalidArgument, "validator address cannot be empty")
	}

	sdkCtx := sdk.UnwrapSDKContext(ctx)
	valAddr := mitotypes.BytesToEthAddress(req.ValAddr)

	validator, found := q.k.GetValidator(sdkCtx, valAddr)
	if !found {
		return nil, status.Errorf(codes.NotFound, "validator %s not found", valAddr.String())
	}

	return &types.QueryValidatorResponse{Validator: validator}, nil
}

// ValidatorByConsAddr returns a validator by consensus address
func (q QueryServer) ValidatorByConsAddr(ctx context.Context, req *types.QueryValidatorByConsAddrRequest) (*types.QueryValidatorByConsAddrResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	if req.ConsAddr == "" {
		return nil, status.Error(codes.InvalidArgument, "consensus address cannot be empty")
	}

	sdkCtx := sdk.UnwrapSDKContext(ctx)

	consAddr, err := q.k.consensusAddressCodec.StringToBytes(req.ConsAddr)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid consensus address: %s", err.Error())
	}

	validator, found := q.k.GetValidatorByConsAddr(sdkCtx, consAddr)
	if !found {
		return nil, status.Errorf(codes.NotFound, "validator with consensus address %s not found", req.ConsAddr)
	}

	return &types.QueryValidatorByConsAddrResponse{Validator: validator}, nil
}

// Validators returns all validators
func (q QueryServer) Validators(ctx context.Context, req *types.QueryValidatorsRequest) (*types.QueryValidatorsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	sdkCtx := sdk.UnwrapSDKContext(ctx)

	store := sdkCtx.KVStore(q.k.storeKey)
	valStore := prefix.NewStore(store, types.ValidatorKeyPrefix)

	var validators []types.Validator
	pageRes, err := query.Paginate(valStore, req.Pagination, func(key []byte, value []byte) error {
		var validator types.Validator
		q.k.cdc.MustUnmarshal(value, &validator)
		validators = append(validators, validator)
		return nil
	})

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryValidatorsResponse{
		Validators: validators,
		Pagination: pageRes,
	}, nil
}

// Withdrawal returns a specific withdrawal by ID
func (q QueryServer) Withdrawal(ctx context.Context, req *types.QueryWithdrawalRequest) (*types.QueryWithdrawalResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	if req.Id == 0 {
		return nil, status.Error(codes.InvalidArgument, "withdrawal ID cannot be 0")
	}

	sdkCtx := sdk.UnwrapSDKContext(ctx)

	// Since withdrawals are indexed by maturesAt and ID, we need to iterate all withdrawals
	// to find one with the requested ID
	var foundWithdrawal *types.Withdrawal
	q.k.IterateWithdrawalsByMaturesAt(sdkCtx, func(withdrawal types.Withdrawal) bool {
		if withdrawal.ID == req.Id {
			foundWithdrawal = &withdrawal
			return true
		}
		return false
	})

	if foundWithdrawal == nil {
		return nil, status.Errorf(codes.NotFound, "withdrawal with ID %d not found", req.Id)
	}

	return &types.QueryWithdrawalResponse{Withdrawal: *foundWithdrawal}, nil
}

// Withdrawals returns all withdrawals
func (q QueryServer) Withdrawals(ctx context.Context, req *types.QueryWithdrawalsRequest) (*types.QueryWithdrawalsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	sdkCtx := sdk.UnwrapSDKContext(ctx)

	store := sdkCtx.KVStore(q.k.storeKey)
	withdrawalStore := prefix.NewStore(store, types.WithdrawalByMaturesAtKeyPrefix)

	var withdrawals []types.Withdrawal
	pageRes, err := query.Paginate(withdrawalStore, req.Pagination, func(key []byte, value []byte) error {
		var withdrawal types.Withdrawal
		q.k.cdc.MustUnmarshal(value, &withdrawal)
		withdrawals = append(withdrawals, withdrawal)
		return nil
	})

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryWithdrawalsResponse{
		Withdrawals: withdrawals,
		Pagination:  pageRes,
	}, nil
}

// WithdrawalsByValidator returns withdrawals for a specific validator
func (q QueryServer) WithdrawalsByValidator(ctx context.Context, req *types.QueryWithdrawalsByValidatorRequest) (*types.QueryWithdrawalsByValidatorResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	if req.ValAddr == nil {
		return nil, status.Error(codes.InvalidArgument, "validator address cannot be empty")
	}

	sdkCtx := sdk.UnwrapSDKContext(ctx)
	valAddr := mitotypes.BytesToEthAddress(req.ValAddr)

	store := sdkCtx.KVStore(q.k.storeKey)
	prefixKey := types.GetWithdrawalByValidatorIterationKey(valAddr)
	withdrawalStore := prefix.NewStore(store, prefixKey)

	var withdrawals []types.Withdrawal
	pageRes, err := query.Paginate(withdrawalStore, req.Pagination, func(key []byte, value []byte) error {
		var withdrawal types.Withdrawal
		q.k.cdc.MustUnmarshal(value, &withdrawal)
		withdrawals = append(withdrawals, withdrawal)
		return nil
	})

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryWithdrawalsByValidatorResponse{
		Withdrawals: withdrawals,
		Pagination:  pageRes,
	}, nil
}
