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
func (k *Keeper) FilterParams() ([]common.Address, [][]common.Hash) {
	return []common.Address{k.evmGovernanceEntrypointAddr},
		[][]common.Hash{{EventMsgExecute.ID}}
}

// Deliver delivers related EVM log events.
func (k *Keeper) Deliver(ctx context.Context, _ common.Hash, elog evmengtypes.EVMEvent) error {
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	cacheCtx, writeCache := sdkCtx.CacheContext()

	// If the processing fails, the error will be logged and the state cache will be discarded.
	if err := catch(func() error { //nolint:contextcheck // False positive wrt ctx
		return k.parseAndProcessEvent(cacheCtx, elog)
	}); err != nil {
		k.Logger(sdkCtx).Error("Delivering event failed",
			"name", eventName(elog),
			"height", cacheCtx.BlockHeight(),
			"err", err,
		)
		return nil
	}

	writeCache()
	return nil
}

// parseAndProcessEvent parses the provided event and processes it.
func (k *Keeper) parseAndProcessEvent(ctx sdk.Context, elog evmengtypes.EVMEvent) error {
	ethlog, err := elog.ToEthLog()
	if err != nil {
		return err
	}

	switch ethlog.Topics[0] {
	case EventMsgExecute.ID:
		event, err := k.evmGovernanceEntrypointContract.ParseMsgExecute(ethlog)
		if err != nil {
			return errors.Wrap(err, "parse msg execute")
		}

		if err := k.processMsgExecute(ctx, event); err != nil {
			return errors.Wrap(err, "process msg execute")
		}
	default:
		return errors.New("unknown event")
	}

	return nil
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

// catch executes the function, returning an error if it panics.
func catch(fn func() error) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = errors.New("recovered", "panic", r)
		}
	}()

	return fn()
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
