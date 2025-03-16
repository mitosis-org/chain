package app

import (
	"context"
	sdklog "cosmossdk.io/log"
	"fmt"
	abci "github.com/cometbft/cometbft/abci/types"
	cmtproto "github.com/cometbft/cometbft/proto/tendermint/types"
	servertypes "github.com/cosmos/cosmos-sdk/server/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

type postFinalizeCallback func(sdk.Context) error

type ABCIWrappedApplication struct {
	servertypes.Application
	postFinalize postFinalizeCallback
	logger       sdklog.Logger
}

func NewABCIWrappedApplication(app *MitosisApp) *ABCIWrappedApplication {
	return &ABCIWrappedApplication{
		Application:  app,
		postFinalize: app.EVMEngKeeper.PostFinalize,
		logger:       app.Logger(),
	}
}

func (a ABCIWrappedApplication) Info(info *abci.RequestInfo) (*abci.ResponseInfo, error) {
	a.logger.Info("😈 ABCI call: Info")

	resp, err := a.Application.Info(info)
	if err != nil {
		a.logger.Error("Info failed [BUG]", "err", err)
	}

	return resp, err
}

func (a ABCIWrappedApplication) Query(ctx context.Context, query *abci.RequestQuery) (*abci.ResponseQuery, error) {
	return a.Application.Query(ctx, query) // No log here since this can be very noisy
}

func (a ABCIWrappedApplication) CheckTx(tx *abci.RequestCheckTx) (*abci.ResponseCheckTx, error) {
	a.logger.Info("😈 ABCI call: CheckTx")
	return a.Application.CheckTx(tx)
}

func (a ABCIWrappedApplication) InitChain(chain *abci.RequestInitChain) (*abci.ResponseInitChain, error) {
	a.logger.Info("😈 ABCI call: InitChain")

	resp, err := a.Application.InitChain(chain)
	if err != nil {
		a.logger.Error("InitChain failed [BUG]", "err", err)
	}

	return resp, err
}

func (a ABCIWrappedApplication) PrepareProposal(proposal *abci.RequestPrepareProposal) (*abci.ResponsePrepareProposal, error) {
	logger := a.logger.With("height", proposal.Height, "proposer", hex7(proposal.ProposerAddress))
	logger.Info("😈 ABCI call: PrepareProposal")

	resp, err := a.Application.PrepareProposal(proposal)
	if err != nil {
		logger.Error("PrepareProposal failed [BUG]", "err", err)
	}

	return resp, err
}

func (a ABCIWrappedApplication) ProcessProposal(proposal *abci.RequestProcessProposal) (*abci.ResponseProcessProposal, error) {
	logger := a.logger.With("height", proposal.Height, "proposer", hex7(proposal.ProposerAddress))
	logger.Info("😈 ABCI call: ProcessProposal")

	resp, err := a.Application.ProcessProposal(proposal)
	if err != nil {
		logger.Error("ProcessProposal failed [BUG]", "err", err)
	}

	return resp, err
}

func (a ABCIWrappedApplication) FinalizeBlock(req *abci.RequestFinalizeBlock) (*abci.ResponseFinalizeBlock, error) {
	logger := a.logger.With("height", req.Height, "proposer", hex7(req.ProposerAddress))
	logger.Info("😈 ABCI call: FinalizeBlock")

	resp, err := a.Application.FinalizeBlock(req)
	if err != nil {
		logger.Error("Finalize req failed [BUG]", "err", err)
		return resp, err
	}

	// Call custom `PostFinalize` callback after the block is finalized.
	header := cmtproto.Header{
		Height:             req.Height,
		Time:               req.Time,
		ProposerAddress:    req.ProposerAddress,
		NextValidatorsHash: req.NextValidatorsHash,
		AppHash:            resp.AppHash, // The app hash after the block is finalized.
	}
	sdkCtx := sdk.NewContext(a.Application.CommitMultiStore().CacheMultiStore(), header, false, nil)
	if err := a.postFinalize(sdkCtx); err != nil {
		logger.Error("PostFinalize callback failed [BUG]", "err", err)
		return resp, err
	}

	attrs := []any{
		"val_updates", len(resp.ValidatorUpdates),
	}
	for i, update := range resp.ValidatorUpdates {
		attrs = append(attrs, fmt.Sprintf("pubkey_%d", i), hex7(update.PubKey.GetSecp256K1()))
		attrs = append(attrs, fmt.Sprintf("power_%d", i), update.Power)
	}
	logger.Info("😈 ABCI response: FinalizeBlock", attrs...)

	for i, res := range resp.TxResults {
		if res.Code == 0 {
			continue
		}
		logger.Error("FinalizeBlock contains unexpected failed transaction [BUG]",
			"info", res.Info, "code", res.Code, "log", res.Log,
			"code_space", res.Codespace, "index", i, "height", req.Height)
	}

	return resp, err
}

func (a ABCIWrappedApplication) ExtendVote(ctx context.Context, vote *abci.RequestExtendVote) (*abci.ResponseExtendVote, error) {
	logger := a.logger.With("height", vote.Height, "proposer", hex7(vote.ProposerAddress))
	logger.Info("😈 ABCI call: ExtendVote")

	resp, err := a.Application.ExtendVote(ctx, vote)
	if err != nil {
		logger.Error("ExtendVote failed [BUG]", "err", err)
	}

	return resp, err
}

func (a ABCIWrappedApplication) VerifyVoteExtension(extension *abci.RequestVerifyVoteExtension) (*abci.ResponseVerifyVoteExtension, error) {
	logger := a.logger.With("height", extension.Height, "validator", hex7(extension.ValidatorAddress))
	logger.Info("😈 ABCI call: VerifyVoteExtension")

	resp, err := a.Application.VerifyVoteExtension(extension)
	if err != nil {
		logger.Error("VerifyVoteExtension failed [BUG]", "err", err)
	}

	return resp, err
}

func (a ABCIWrappedApplication) Commit() (*abci.ResponseCommit, error) {
	a.logger.Info("😈 ABCI call: Commit")
	return a.Application.Commit()
}

func (a ABCIWrappedApplication) ListSnapshots(listSnapshots *abci.RequestListSnapshots) (*abci.ResponseListSnapshots, error) {
	a.logger.Info("😈 ABCI call: ListSnapshots")
	return a.Application.ListSnapshots(listSnapshots)
}

func (a ABCIWrappedApplication) OfferSnapshot(snapshot *abci.RequestOfferSnapshot) (*abci.ResponseOfferSnapshot, error) {
	a.logger.Info("😈 ABCI call: OfferSnapshot")
	return a.Application.OfferSnapshot(snapshot)
}

func (a ABCIWrappedApplication) LoadSnapshotChunk(chunk *abci.RequestLoadSnapshotChunk) (*abci.ResponseLoadSnapshotChunk, error) {
	a.logger.Info("😈 ABCI call: LoadSnapshotChunk")
	return a.Application.LoadSnapshotChunk(chunk)
}

func (a ABCIWrappedApplication) ApplySnapshotChunk(chunk *abci.RequestApplySnapshotChunk) (*abci.ResponseApplySnapshotChunk, error) {
	a.logger.Info("😈 ABCI call: ApplySnapshotChunk")
	return a.Application.ApplySnapshotChunk(chunk)
}

func hex7(value []byte) string {
	h := fmt.Sprintf("%X", value)

	const maxLen = 7
	if len(h) > maxLen {
		h = h[:maxLen]
	}

	return h
}
