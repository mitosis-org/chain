package app

import (
	abci "github.com/cometbft/cometbft/abci/types"
	cmttypes "github.com/cometbft/cometbft/proto/tendermint/types"
	"github.com/omni-network/omni/lib/errors"
	"github.com/omni-network/omni/lib/log"

	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/client"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func makePrepareProposalHandler(
	app *MitosisApp,
	txConfig client.TxConfig,
	prevHandler sdk.PrepareProposalHandler,
) sdk.PrepareProposalHandler {
	return func(ctx sdk.Context, req *abci.RequestPrepareProposal) (*abci.ResponsePrepareProposal, error) {
		// Use evm engine to create block proposals.
		// Note that we do not check MaxTxBytes since all EngineEVM transaction MUST be included since we cannot
		// postpone them to the next block.
		reqForEVM := *req
		reqForEVM.Txs = nil
		resForEVM, err := app.EVMEngKeeper.PrepareProposal(ctx, &reqForEVM)
		if err != nil {
			return nil, err
		}

		var nonEVMTxs [][]byte
		for _, rawTX := range req.Txs {
			tx, err := txConfig.TxDecoder()(rawTX)
			if err != nil {
				return nil, errors.Wrap(err, "decode transaction")
			}

			if isTxForEVM(tx) {
				// Leave logs and ignore EVM payload transactions.
				log.Warn(ctx, "EVM payload transaction should be not included in a prepare proposal", nil)
			} else {
				nonEVMTxs = append(nonEVMTxs, rawTX)
			}
		}

		reqForNonEVM := *req
		reqForNonEVM.Txs = nonEVMTxs

		// We should decrease MaxTxBytes by the size of the EVM payload transaction.
		for _, evmTx := range resForEVM.Txs {
			reqForNonEVM.MaxTxBytes -= int64(len(evmTx))
		}
		if reqForNonEVM.MaxTxBytes <= 0 {
			// It means that we can't include any non-EVM transactions more.
			return resForEVM, nil
		}

		resForNonEVM, err := prevHandler(ctx, &reqForNonEVM)
		if err != nil {
			return nil, err
		}

		return &abci.ResponsePrepareProposal{Txs: append(resForEVM.Txs, resForNonEVM.Txs...)}, nil
	}
}

// makeProcessProposalRouter creates a new process proposal router that only routes
// expected messages to expected modules.
func makeProcessProposalRouter(app *MitosisApp) *baseapp.MsgServiceRouter {
	router := baseapp.NewMsgServiceRouter()
	router.SetInterfaceRegistry(app.interfaceRegistry)
	app.EVMEngKeeper.RegisterProposalService(router) // EVMEngine calls NewPayload on proposals to verify it.

	return router
}

func makeProcessProposalHandler(
	router *baseapp.MsgServiceRouter,
	txConfig client.TxConfig,
	prevHandler sdk.ProcessProposalHandler,
) sdk.ProcessProposalHandler {
	return func(ctx sdk.Context, req *abci.RequestProcessProposal) (*abci.ResponseProcessProposal, error) {
		nonEVMTxs, err := processTxForEVM(ctx, req, router, txConfig)
		if err != nil {
			return nil, err
		}

		reqForNonEVMTxs := *req
		reqForNonEVMTxs.Txs = nonEVMTxs

		return prevHandler(ctx, &reqForNonEVMTxs)
	}
}

func processTxForEVM(
	ctx sdk.Context,
	req *abci.RequestProcessProposal,
	router *baseapp.MsgServiceRouter,
	txConfig client.TxConfig,
) ([][]byte, error) {
	// TODO(thai): Is this code necessary? I just wanted to use same code as Omni.
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
			return nil, errors.New("proposed doesn't include quorum votes extensions")
		}
	}

	var nonEVMTxs [][]byte

	for _, rawTX := range req.Txs {
		tx, err := txConfig.TxDecoder()(rawTX)
		if err != nil {
			return nil, errors.Wrap(err, "decode transaction")
		}

		if isTxForEVM(tx) {
			if len(tx.GetMsgs()) != 1 {
				return nil, errors.New("EVM payload transaction must contain exactly one message")
			}

			msg := tx.GetMsgs()[0]

			handler := router.Handler(msg)
			if handler == nil {
				return nil, errors.New("EVM msg handler not found [BUG]")
			}

			_, err := handler(ctx, msg)
			if err != nil {
				return nil, errors.Wrap(err, "execute message")
			}
		} else {
			nonEVMTxs = append(nonEVMTxs, rawTX)
		}
	}

	return nonEVMTxs, nil
}
