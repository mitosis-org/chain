package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/mitosis-org/chain/bindings"
	"github.com/mitosis-org/chain/x/evmgov/types"
	"github.com/omni-network/omni/lib/errors"
	evmengtypes "github.com/omni-network/omni/octane/evmengine/types"
)

var (
	_ evmengtypes.EvmEventProcessor = &Keeper{}

	ABI             = mustGetABI(bindings.ConsensusGovernanceEntrypointMetaData)
	EventMsgExecute = mustGetEvent(ABI, "MsgExecute")

	EventsByID = map[common.Hash]abi.Event{
		EventMsgExecute.ID: EventMsgExecute,
	}
)

func (*Keeper) Name() string {
	return types.ModuleName
}

// FilterParams defines the matching EVM log events.
func (k *Keeper) FilterParams(_ context.Context) ([]common.Address, [][]common.Hash) {
	return []common.Address{k.evmGovernanceEntrypointAddr},
		[][]common.Hash{{EventMsgExecute.ID}}
}

// Deliver delivers related EVM log events.
func (k *Keeper) Deliver(ctx context.Context, blockHash common.Hash, elog evmengtypes.EVMEvent) error {
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	cacheCtx, writeCache := sdkCtx.CacheContext()

	err, ignore := k.processEvent(cacheCtx, blockHash, elog)
	if err != nil {
		if ignore {
			// If the processing fails but needs to be ignored, the error will be logged and
			// the state cache will be discarded.
			k.Logger(sdkCtx).Error("Processing event failed but ignored",
				"name", eventName(elog),
				"height", cacheCtx.BlockHeight(),
				"evmBlockHash", blockHash.Hex(),
				"evmLog", elog.String(),
				"err", err,
			)
			return nil
		} else {
			return errors.Wrap(err, "failed to process event",
				"name", eventName(elog),
				"height", cacheCtx.BlockHeight(),
				"evmBlockHash", blockHash.Hex(),
				"evmLog", elog.String(),
			)
		}
	}

	writeCache()
	return nil
}

// processEvent parses the provided event and processes it.
// If the second return value is true, the error will be ignored.
func (k *Keeper) processEvent(ctx sdk.Context, _ common.Hash, elog evmengtypes.EVMEvent) (error, bool) {
	ethlog, err := elog.ToEthLog()
	if err != nil {
		return err, false
	}

	switch ethlog.Topics[0] {
	// Potential failure cases are:
	// - Messages has invalid format. (could be not verified at the EVM contract level)
	// - Failed to execute messages. (could be not verified at the EVM contract level)
	// Fortunately, this logic is not critical. If it fails, we can identify the cause
	// and then proceed with the governance process again at the EVM contract level.
	case EventMsgExecute.ID:
		event, err := k.evmGovernanceEntrypointContract.ParseMsgExecute(ethlog)
		if err != nil {
			return errors.Wrap(err, "parse MsgExecute"), false
		}

		if err := k.processMsgExecute(ctx, event); err != nil {
			return errors.Wrap(err, "process MsgExecute"), true
		}
	default:
		return errors.New("unknown event"), false
	}

	return nil, false
}

// processMsgExecute processes the MsgExecute event.
func (k *Keeper) processMsgExecute(ctx sdk.Context, event *bindings.ConsensusGovernanceEntrypointMsgExecute) error {
	msgs, err := k.ParseMessages(event.Messages)
	if err != nil {
		return err
	}

	return k.ExecuteMessages(ctx, msgs)
}

// mustGetABI returns the metadata's ABI as an abi.ABI type.
// It panics on error.
func mustGetABI(metadata *bind.MetaData) *abi.ABI {
	abi, err := metadata.GetAbi()
	if err != nil {
		panic(err)
	}

	return abi
}

// mustGetEvent returns the event with the given name from the ABI.
// It panics if the event is not found.
func mustGetEvent(abi *abi.ABI, name string) abi.Event {
	event, ok := abi.Events[name]
	if !ok {
		panic("event not found")
	}

	return event
}

// eventName returns the name of the EVM event log or "unknown".
func eventName(elog evmengtypes.EVMEvent) string {
	ethlog, err := elog.ToEthLog()
	if err != nil {
		return "unknown"
	}

	event, ok := EventsByID[ethlog.Topics[0]]
	if !ok {
		return "unknown"
	}

	return event.Name
}
