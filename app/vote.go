package app

import (
	"context"
	abci "github.com/cometbft/cometbft/abci/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	evmengtypes "github.com/omni-network/omni/octane/evmengine/types"
)

var _ evmengtypes.VoteExtensionProvider = &NoVoteExtensionProvider{}

type NoVoteExtensionProvider struct {
}

func (n NoVoteExtensionProvider) PrepareVotes(ctx context.Context, commit abci.ExtendedCommitInfo) ([]sdk.Msg, error) {
	return nil, nil
}
