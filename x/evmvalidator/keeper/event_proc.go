package keeper

import (
	"context"
	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/mitosis-org/chain/bindings"
	"github.com/mitosis-org/chain/x/evmvalidator/types"
	"github.com/omni-network/omni/lib/errors"
	"github.com/omni-network/omni/lib/k1util"
	evmengtypes "github.com/omni-network/omni/octane/evmengine/types"
)

var (
	_ evmengtypes.EvmEventProcessor = &Keeper{}

	ABI                            = mustGetABI(bindings.ConsensusValidatorEntrypointMetaData)
	EventMsgRegisterValidator      = mustGetEvent(ABI, "MsgRegisterValidator")
	EventMsgDepositCollateral      = mustGetEvent(ABI, "MsgDepositCollateral")
	EventMsgWithdrawCollateral     = mustGetEvent(ABI, "MsgWithdrawCollateral")
	EventMsgUnjail                 = mustGetEvent(ABI, "MsgUnjail")
	EventMsgUpdateExtraVotingPower = mustGetEvent(ABI, "MsgUpdateExtraVotingPower")

	EventsByID = map[common.Hash]abi.Event{
		EventMsgRegisterValidator.ID:      EventMsgRegisterValidator,
		EventMsgDepositCollateral.ID:      EventMsgDepositCollateral,
		EventMsgWithdrawCollateral.ID:     EventMsgWithdrawCollateral,
		EventMsgUnjail.ID:                 EventMsgUnjail,
		EventMsgUpdateExtraVotingPower.ID: EventMsgUpdateExtraVotingPower,
	}
)

// Name returns the name of the module
func (*Keeper) Name() string {
	return types.ModuleName
}

// FilterParams defines the matching EVM log events
func (k *Keeper) FilterParams() ([]common.Address, [][]common.Hash) {
	return []common.Address{k.evmValidatorEntrypointAddr},
		[][]common.Hash{
			{
				EventMsgRegisterValidator.ID,
				EventMsgDepositCollateral.ID,
				EventMsgWithdrawCollateral.ID,
				EventMsgUnjail.ID,
				EventMsgUpdateExtraVotingPower.ID,
			},
		}
}

// Deliver delivers related EVM log events
func (k *Keeper) Deliver(ctx context.Context, _ common.Hash, elog evmengtypes.EVMEvent) error {
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	cacheCtx, writeCache := sdkCtx.CacheContext()

	// If the processing fails, the error will be logged and the state cache will be discarded
	if err := catch(func() error {
		return k.parseAndProcessEvent(cacheCtx, elog)
	}); err != nil {
		sdkCtx.Logger().Error("Delivering event failed",
			"name", eventName(elog),
			"height", cacheCtx.BlockHeight(),
			"err", err,
		)
		return nil
	}

	writeCache()
	return nil
}

// parseAndProcessEvent parses the provided event and processes it
func (k *Keeper) parseAndProcessEvent(ctx sdk.Context, elog evmengtypes.EVMEvent) error {
	ethlog, err := elog.ToEthLog()
	if err != nil {
		return err
	}

	switch ethlog.Topics[0] {
	case EventMsgRegisterValidator.ID:
		event, err := k.evmValidatorEntrypointContract.ParseMsgRegisterValidator(ethlog)
		if err != nil {
			return errors.Wrap(err, "parse register validator event")
		}
		return k.processRegisterValidator(ctx, event)

	case EventMsgDepositCollateral.ID:
		event, err := k.evmValidatorEntrypointContract.ParseMsgDepositCollateral(ethlog)
		if err != nil {
			return errors.Wrap(err, "parse deposit collateral event")
		}
		return k.processDepositCollateral(ctx, event)

	case EventMsgWithdrawCollateral.ID:
		event, err := k.evmValidatorEntrypointContract.ParseMsgWithdrawCollateral(ethlog)
		if err != nil {
			return errors.Wrap(err, "parse withdraw collateral event")
		}
		return k.processWithdrawCollateral(ctx, event)

	case EventMsgUnjail.ID:
		event, err := k.evmValidatorEntrypointContract.ParseMsgUnjail(ethlog)
		if err != nil {
			return errors.Wrap(err, "parse unjail event")
		}
		return k.processUnjail(ctx, event)

	case EventMsgUpdateExtraVotingPower.ID:
		event, err := k.evmValidatorEntrypointContract.ParseMsgUpdateExtraVotingPower(ethlog)
		if err != nil {
			return errors.Wrap(err, "parse update extra voting power event")
		}
		return k.processUpdateExtraVotingPower(ctx, event)

	default:
		return errors.New("unknown event")
	}
}

// processRegisterValidator processes MsgRegisterValidator event
func (k *Keeper) processRegisterValidator(ctx sdk.Context, event *bindings.ConsensusValidatorEntrypointMsgRegisterValidator) error {
	pubkey := event.Valkey

	// Validate the public key
	if err := validatePubkey(pubkey); err != nil {
		return errors.Wrap(err, "invalid validator pubkey")
	}

	// Convert the amount to math.Int
	collateral := sdkmath.NewIntFromBigInt(event.InitialCollateralAmount)

	// Register the validator
	return k.registerValidator(ctx, pubkey, collateral, sdkmath.ZeroInt(), false)
}

// processDepositCollateral processes MsgDepositCollateral event
func (k *Keeper) processDepositCollateral(ctx sdk.Context, event *bindings.ConsensusValidatorEntrypointMsgDepositCollateral) error {
	pubkey := event.Valkey

	// Validate the public key
	if err := validatePubkey(pubkey); err != nil {
		return errors.Wrap(err, "invalid validator pubkey")
	}

	// Check if validator exists
	validator, found := k.GetValidator(ctx, pubkey)
	if !found {
		return types.ErrValidatorNotFound
	}

	// Convert the amount to math.Int
	amount := sdkmath.NewIntFromBigInt(event.Amount)

	// Update validator's collateral
	if err := k.depositCollateral(ctx, &validator, amount); err != nil {
		return errors.Wrap(err, "failed to deposit collateral")
	}

	return nil
}

// processWithdrawCollateral processes MsgWithdrawCollateral event
func (k *Keeper) processWithdrawCollateral(ctx sdk.Context, event *bindings.ConsensusValidatorEntrypointMsgWithdrawCollateral) error {
	pubkey := event.Valkey

	// Validate the public key
	if err := validatePubkey(pubkey); err != nil {
		return errors.Wrap(err, "invalid validator pubkey")
	}

	// Check if validator exists
	validator, found := k.GetValidator(ctx, pubkey)
	if !found {
		return types.ErrValidatorNotFound
	}

	// Convert the amount to math.Int
	amount := sdkmath.NewIntFromBigInt(event.Amount)

	// Create a withdrawal
	withdrawal := types.NewWithdrawal(
		pubkey,
		amount,
		event.Receiver.String(),
		event.ReceivesAt.Uint64(),
	)

	// Request withdrawal
	if err := k.withdrawCollateral(ctx, &validator, withdrawal); err != nil {
		return errors.Wrap(err, "failed to withdraw collateral")
	}

	return nil
}

// processUnjail processes MsgUnjail event
func (k *Keeper) processUnjail(ctx sdk.Context, event *bindings.ConsensusValidatorEntrypointMsgUnjail) error {
	pubkey := event.Valkey

	// Validate the public key
	if err := validatePubkey(pubkey); err != nil {
		return errors.Wrap(err, "invalid validator pubkey")
	}

	// Check if validator exists
	validator, found := k.GetValidator(ctx, pubkey)
	if !found {
		return types.ErrValidatorNotFound
	}

	// Check if validator is jailed
	if !validator.Jailed {
		// NOTE: There is no verification logic in EVM. So, just ignore instead of returning an error
		// because there could be many error logs otherwise.
		return nil // No-op if not jailed
	}

	// Get consensus address
	consAddr, err := validator.ConsAddr()
	if err != nil {
		return errors.Wrap(err, "failed to get consensus address")
	}

	// unjail validator through slashing keeper
	if err = k.slashingKeeper.UnjailFromConsAddr(ctx, consAddr); err != nil {
		return errors.Wrap(err, "failed to unjail validator")
	}

	return nil
}

// processUpdateExtraVotingPower processes MsgUpdateExtraVotingPower event
func (k *Keeper) processUpdateExtraVotingPower(ctx sdk.Context, event *bindings.ConsensusValidatorEntrypointMsgUpdateExtraVotingPower) error {
	pubkey := event.Valkey

	// Validate the public key
	if err := validatePubkey(pubkey); err != nil {
		return errors.Wrap(err, "invalid validator pubkey")
	}

	// Check if validator exists
	validator, found := k.GetValidator(ctx, pubkey)
	if !found {
		return types.ErrValidatorNotFound
	}

	// Convert the extra voting power to math.Int
	extraVotingPower := sdkmath.NewIntFromBigInt(event.ExtraVotingPower)

	// Update extra voting power
	if err := k.updateExtraVotingPower(ctx, &validator, extraVotingPower); err != nil {
		return errors.Wrap(err, "failed to update extra voting power")
	}

	return nil
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

// validatePubkey validates the public key format
func validatePubkey(pubkey []byte) error {
	if len(pubkey) != 33 { // Compressed secp256k1 pubkey is 33 bytes
		return types.ErrInvalidPubKey
	}

	// Additional validation if needed
	// Try to convert to cosmos pubkey
	_, err := k1util.PubKeyBytesToCosmos(pubkey)
	if err != nil {
		return errors.Wrap(err, "invalid pubkey format")
	}

	return nil
}
