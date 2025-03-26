package keeper

import (
	"context"

	"cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	"github.com/mitosis-org/chain/x/evmvalidator/types"
)

var _ types.MsgServer = msgServer{}

type msgServer struct {
	k *Keeper
}

// NewMsgServerImpl returns an implementation of the evmvalidator MsgServer interface
// for the provided Keeper.
func NewMsgServerImpl(keeper *Keeper) types.MsgServer {
	return &msgServer{k: keeper}
}

// UpdateParams updates the module parameters
func (m msgServer) UpdateParams(ctx context.Context, msg *types.MsgUpdateParams) (*types.MsgUpdateParamsResponse, error) {
	sdkCtx := sdk.UnwrapSDKContext(ctx)

	if m.k.authority != msg.Authority {
		return nil, errors.Wrapf(govtypes.ErrInvalidSigner, "invalid authority; expected %s, got %s", m.k.authority, msg.Authority)
	}

	if err := msg.Params.Validate(); err != nil {
		return nil, err
	}

	oldParams := m.k.GetParams(sdkCtx)

	if err := m.k.SetParams(sdkCtx, msg.Params); err != nil {
		return nil, err
	}

	if oldParams.MaxLeverageRatio != msg.Params.MaxLeverageRatio || oldParams.MinVotingPower != msg.Params.MinVotingPower {
		validators := m.k.GetAllValidators(sdkCtx)
		for _, validator := range validators {
			m.k.UpdateValidatorState(sdkCtx, &validator, "update params")
		}
	}

	return &types.MsgUpdateParamsResponse{}, nil
}

// UpdateValidatorEntrypointContractAddr updates the address of the ConsensusValidatorEntrypoint contract
func (m msgServer) UpdateValidatorEntrypointContractAddr(ctx context.Context, msg *types.MsgUpdateValidatorEntrypointContractAddr) (*types.MsgUpdateValidatorEntrypointContractAddrResponse, error) {
	sdkCtx := sdk.UnwrapSDKContext(ctx)

	if m.k.authority != msg.Authority {
		return nil, errors.Wrapf(govtypes.ErrInvalidSigner, "invalid authority; expected %s, got %s", m.k.authority, msg.Authority)
	}

	m.k.SetValidatorEntrypointContractAddr(sdkCtx, msg.Addr)

	return &types.MsgUpdateValidatorEntrypointContractAddrResponse{}, nil
}
