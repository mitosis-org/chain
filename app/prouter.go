package app

import (
	"context"
	"strings"

	"github.com/omni-network/omni/lib/errors"
	"github.com/omni-network/omni/lib/log"
	etypes "github.com/omni-network/omni/octane/evmengine/types"

	abci "github.com/cometbft/cometbft/abci/types"
	cmttypes "github.com/cometbft/cometbft/proto/tendermint/types"

	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/client"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type MsgHandlerManager interface {
	Handler(msg sdk.Msg) baseapp.MsgServiceHandler
}

type MsgHandlerManagerForOctane struct {
	normal    *baseapp.MsgServiceRouter
	forOctane *baseapp.MsgServiceRouter
}

func NewMsgHandlerManagerForOctane(app *MitosisApp) MsgHandlerManagerForOctane {
	forOctane := baseapp.NewMsgServiceRouter()
	forOctane.SetInterfaceRegistry(app.interfaceRegistry)
	app.EVMEngKeeper.RegisterProposalService(forOctane)

	return MsgHandlerManagerForOctane{
		normal:    app.MsgServiceRouter(),
		forOctane: forOctane,
	}
}

func (m MsgHandlerManagerForOctane) Handler(msg sdk.Msg) baseapp.MsgServiceHandler {
	isOctaneMsg := sdk.MsgTypeURL(msg) == sdk.MsgTypeURL(&etypes.MsgExecutionPayload{})

	if isOctaneMsg {
		return m.forOctane.Handler(msg)
	} else {
		return m.normal.Handler(msg)
	}
}

// makeProcessProposalHandler creates a new process proposal handler.
// It ensures all messages included in a cpayload proposal are valid.
// It also updates some external state.
func makeProcessProposalHandler(app *MitosisApp, txConfig client.TxConfig) sdk.ProcessProposalHandler {
	handlerManager := NewMsgHandlerManagerForOctane(app)

	return func(ctx sdk.Context, req *abci.RequestProcessProposal) (*abci.ResponseProcessProposal, error) {
		// Ensure the proposal includes quorum vote extensions (unless first block).
		if req.Height > 1 {
			var totalPower, votedPower int64
			for _, vote := range req.ProposedLastCommit.Votes {
				totalPower += vote.Validator.Power
				if vote.BlockIdFlag != cmttypes.BlockIDFlagCommit {
					continue
				}
				votedPower += vote.Validator.Power
			}
			if totalPower*2/3 >= votedPower {
				return rejectProposal(ctx, errors.New("proposed doesn't include quorum votes extensions"))
			}
		}

		// Ensure only expected messages types are included the expected number of times.
		allowedMsgCounts := map[string]int{
			sdk.MsgTypeURL(&etypes.MsgExecutionPayload{}): 1, // Only a single EVM execution payload is allowed.
		}

		for _, rawTX := range req.Txs {
			tx, err := txConfig.TxDecoder()(rawTX)
			if err != nil {
				return rejectProposal(ctx, errors.Wrap(err, "decode transaction"))
			}

			for _, msg := range tx.GetMsgs() {
				typeURL := sdk.MsgTypeURL(msg)

				// TODO(thai): should revise it again.
				//  For now, we just allow ibc messages for integration with ethos.
				if !strings.HasPrefix(typeURL, "/ibc.") {
					// Ensure the message type is expected and not included too many times.
					if i, ok := allowedMsgCounts[typeURL]; !ok {
						return rejectProposal(ctx, errors.New("unexpected message type", "msg_type", typeURL))
					} else if i <= 0 {
						return rejectProposal(ctx, errors.New("message type included too many times", "msg_type", typeURL))
					}
					allowedMsgCounts[typeURL]--
				}

				handler := handlerManager.Handler(msg)
				if handler == nil {
					return rejectProposal(ctx, errors.New("msg handler not found [BUG]", "msg_type", typeURL))
				}

				_, err := handler(ctx, msg)
				if err != nil {
					return rejectProposal(ctx, errors.Wrap(err, "execute message"))
				}
			}
		}

		return &abci.ResponseProcessProposal{Status: abci.ResponseProcessProposal_ACCEPT}, nil
	}
}

func rejectProposal(ctx context.Context, err error) (*abci.ResponseProcessProposal, error) {
	log.Error(ctx, "Rejecting process proposal", err)
	return &abci.ResponseProcessProposal{Status: abci.ResponseProcessProposal_REJECT}, nil
}
