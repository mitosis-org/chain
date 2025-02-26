package keeper

import (
	"fmt"
	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/mitosis-org/chain/bindings"
	"github.com/omni-network/omni/lib/errors"
)

type Keeper struct {
	cdc                             codec.Codec
	router                          baseapp.MessageRouter
	evmGovernanceEntrypointAddr     common.Address
	evmGovernanceEntrypointContract *bindings.ConsensusGovernanceEntrypoint
}

func NewKeeper(cdc codec.Codec, router baseapp.MessageRouter, evmGovernanceEntrypointAddr common.Address) (*Keeper, error) {
	evmGovernanceEntrypointContract, err := bindings.NewConsensusGovernanceEntrypoint(evmGovernanceEntrypointAddr, nil)
	if err != nil {
		return nil, err
	}

	return &Keeper{
		cdc,
		router,
		evmGovernanceEntrypointAddr,
		evmGovernanceEntrypointContract,
	}, nil
}

func (k *Keeper) ParseMessages(rawMsgs []string) ([]sdk.Msg, error) {
	// Example of rawMsg: {"@type": "/cosmos.upgrade.v1beta1.MsgSoftwareUpgrade", "authority": "...", "plan": {...}}

	var msgs []sdk.Msg

	// TODO(thai): There is no error even though there is missing field. How can we make sure there is no missing field?
	for _, rawMsg := range rawMsgs {
		var protoMsg codectypes.Any
		err := k.cdc.UnmarshalJSON([]byte(rawMsg), &protoMsg)
		if err != nil {
			return nil, errors.Wrap(err, "failed to parse message")
		}

		var msg sdk.Msg
		err = k.cdc.UnpackAny(&protoMsg, &msg)
		if err != nil {
			return nil, errors.Wrap(err, fmt.Sprintf("failed to unpack message of %s", protoMsg.TypeUrl))
		}

		msgs = append(msgs, msg)
	}

	return msgs, nil
}

func (k *Keeper) ExecuteMessages(ctx sdk.Context, msgs []sdk.Msg) error {
	for _, msg := range msgs {
		handler := k.router.Handler(msg)

		_, err := safeExecuteHandler(ctx, msg, handler)
		if err != nil {
			return errors.Wrap(err, "failed to execute message")
		}
	}

	return nil
}

// safeExecuteHandler executes handle(msg) and recovers from panic.
func safeExecuteHandler(ctx sdk.Context, msg sdk.Msg, handler baseapp.MsgServiceHandler,
) (res *sdk.Result, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("handling msg [%s] PANICKED: %v", msg, r)
		}
	}()
	res, err = handler(ctx, msg)
	return
}
