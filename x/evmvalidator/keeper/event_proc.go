package keeper

import (
	"context"
	sdkmath "cosmossdk.io/math"
	stderrors "errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/mitosis-org/chain/bindings"
	mitotypes "github.com/mitosis-org/chain/types"
	"github.com/mitosis-org/chain/x/evmvalidator/types"
	"github.com/omni-network/omni/lib/errors"
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

	err, ignored := k.processEvent(cacheCtx, elog)
	if err != nil {
		if ignored {
			// If the processing fails but needs to be ignored, the error will be logged and
			// the state cache will be discarded.
			k.Logger(sdkCtx).Error("Processing event failed but ignored",
				"name", eventName(elog),
				"height", cacheCtx.BlockHeight(),
				"err", err,
			)
			return nil
		} else {
			return errors.Wrap(err, "failed to process event")
		}
	}

	writeCache()
	return nil
}

// processEvent parses the provided event and processes it.
// If the second return value is true, the error will be ignored.
func (k *Keeper) processEvent(ctx sdk.Context, elog evmengtypes.EVMEvent) (error, bool) {
	ethlog, err := elog.ToEthLog()
	if err != nil {
		return err, false
	}

	switch ethlog.Topics[0] {
	// Potential failure cases are:
	// - The validator already exist (might be verified at the EVM contract level)
	// - valAddr and pubKey are not consistent (might be verified at the EVM contract level)
	// We must refund the collateral to the user through fallback logic if the primary logic fails.
	// The fallback logic must not fail due to its critical nature and should not fail because it's trivial.
	// Therefore, we raise an error if the fallback fails.
	case EventMsgRegisterValidator.ID:
		event, err := k.evmValidatorEntrypointContract.ParseMsgRegisterValidator(ethlog)
		if err != nil {
			return errors.Wrap(err, "parse MsgRegisterValidator"), false
		}
		if err := k.processRegisterValidator(ctx, event); err != nil {
			if errFB := k.fallbackRegisterValidator(ctx, event); errFB != nil {
				return stderrors.Join(
					errors.Wrap(err, "process MsgRegisterValidator"),
					errors.Wrap(errFB, "fallback MsgRegisterValidator"),
				), false
			}

			k.Logger(ctx).Error("Processing failed but fallback succeeded",
				"name", eventName(elog),
				"height", ctx.BlockHeight(),
				"err", err)
		}

	// Potential failure cases are:
	// - The validator does not exist (might be verified at the EVM contract level)
	// We must refund the collateral to the user through fallback logic if the primary logic fails.
	// The fallback logic must not fail due to its critical nature and should not fail because it's trivial.
	// Therefore, we raise an error if the fallback fails.
	case EventMsgDepositCollateral.ID:
		event, err := k.evmValidatorEntrypointContract.ParseMsgDepositCollateral(ethlog)
		if err != nil {
			return errors.Wrap(err, "parse MsgDepositCollateral"), false
		}
		if err := k.processDepositCollateral(ctx, event); err != nil {
			if errFB := k.fallbackDepositCollateral(ctx, event); errFB != nil {
				return stderrors.Join(
					errors.Wrap(err, "process MsgDepositCollateral"),
					errors.Wrap(errFB, "fallback MsgDepositCollateral"),
				), false
			}

			k.Logger(ctx).Error("Processing failed but fallback succeeded",
				"name", eventName(elog),
				"height", ctx.BlockHeight(),
				"err", err)
		}

	// Potential failure cases are:
	// - The validator does not exist (might be verified at the EVM contract level)
	// - The withdrawal amount is greater than the validator's collateral (could be not verified at the EVM contract level)
	// Fortunately, this logic is not critical. Even if it fails, users won't lose money
	// and the state won't become corrupted. Therefore, we simply ignore errors when they occur.
	case EventMsgWithdrawCollateral.ID:
		event, err := k.evmValidatorEntrypointContract.ParseMsgWithdrawCollateral(ethlog)
		if err != nil {
			return errors.Wrap(err, "parse MsgWithdrawCollateral"), false
		}
		if err := k.processWithdrawCollateral(ctx, event); err != nil {
			return errors.Wrap(err, "process MsgWithdrawCollateral"), true
		}

	// Potential failure cases are:
	// - The validator does not exist (might be verified at the EVM contract level)
	// - The validator is not jailed (could be not verified at the EVM contract level)
	// Fortunately, this logic is not critical. Even if it fails, users won't lose money
	// and the state won't become corrupted. Therefore, we simply ignore errors when they occur.
	case EventMsgUnjail.ID:
		event, err := k.evmValidatorEntrypointContract.ParseMsgUnjail(ethlog)
		if err != nil {
			return errors.Wrap(err, "parse MsgUnjail"), false
		}
		if err := k.processUnjail(ctx, event); err != nil {
			// NOTE: It is not critical so ignore the error.
			return errors.Wrap(err, "process MsgUnjail"), true
		}

	// Potential failure cases are:
	// - The validator does not exist (might be verified at the EVM contract level)
	// Fortunately, this logic is not critical. Even if it fails, users won't lose money
	// and the state won't become corrupted. Therefore, we simply ignore errors when they occur.
	case EventMsgUpdateExtraVotingPower.ID:
		event, err := k.evmValidatorEntrypointContract.ParseMsgUpdateExtraVotingPower(ethlog)
		if err != nil {
			return errors.Wrap(err, "parse MsgUpdateExtraVotingPower"), false
		}
		if err := k.processUpdateExtraVotingPower(ctx, event); err != nil {
			// NOTE: It is not critical so ignore the error.
			return errors.Wrap(err, "process MsgUpdateExtraVotingPower"), true
		}

	default:
		return errors.New("unknown event"), false
	}

	return nil, false
}

// processRegisterValidator processes MsgRegisterValidator event
func (k *Keeper) processRegisterValidator(ctx sdk.Context, event *bindings.ConsensusValidatorEntrypointMsgRegisterValidator) error {
	valAddr := mitotypes.EthAddress(event.ValAddr)

	// Convert the amount to math.Int
	collateral := sdkmath.NewIntFromBigInt(event.InitialCollateralAmountGwei)

	// Register the validator
	return k.registerValidator(ctx, valAddr, event.PubKey, collateral, sdkmath.ZeroInt(), false)
}

// fallbackRegisterValidator handles the case when the MsgRegisterValidator event fails to process
func (k *Keeper) fallbackRegisterValidator(ctx sdk.Context, event *bindings.ConsensusValidatorEntrypointMsgRegisterValidator) error {
	return k.evmEngKeeper.InsertWithdrawal(ctx, event.CollateralReturnAddr, event.InitialCollateralAmountGwei.Uint64())
}

// processDepositCollateral processes MsgDepositCollateral event
func (k *Keeper) processDepositCollateral(ctx sdk.Context, event *bindings.ConsensusValidatorEntrypointMsgDepositCollateral) error {
	valAddr := mitotypes.EthAddress(event.ValAddr)

	// Check if validator exists
	validator, found := k.GetValidator(ctx, valAddr)
	if !found {
		return types.ErrValidatorNotFound
	}

	// Convert the amount to math.Int
	amount := sdkmath.NewIntFromBigInt(event.AmountGwei)

	// Update validator's collateral
	k.depositCollateral(ctx, &validator, amount)

	return nil
}

// fallbackDepositCollateral handles the case when the MsgDepositCollateral event fails to process
func (k *Keeper) fallbackDepositCollateral(ctx sdk.Context, event *bindings.ConsensusValidatorEntrypointMsgDepositCollateral) error {
	return k.evmEngKeeper.InsertWithdrawal(ctx, event.CollateralReturnAddr, event.AmountGwei.Uint64())
}

// processWithdrawCollateral processes MsgWithdrawCollateral event
func (k *Keeper) processWithdrawCollateral(ctx sdk.Context, event *bindings.ConsensusValidatorEntrypointMsgWithdrawCollateral) error {
	valAddr := mitotypes.EthAddress(event.ValAddr)

	// Check if validator exists
	validator, found := k.GetValidator(ctx, valAddr)
	if !found {
		return types.ErrValidatorNotFound
	}

	// Create a withdrawal
	withdrawal := types.Withdrawal{
		ValAddr:   valAddr,
		Amount:    event.AmountGwei.Uint64(),
		Receiver:  mitotypes.EthAddress(event.Receiver),
		MaturesAt: event.MaturesAt.Uint64(),
	}

	// Request withdrawal
	if err := k.withdrawCollateral(ctx, &validator, withdrawal); err != nil {
		return errors.Wrap(err, "failed to withdraw collateral")
	}

	return nil
}

// processUnjail processes MsgUnjail event
func (k *Keeper) processUnjail(ctx sdk.Context, event *bindings.ConsensusValidatorEntrypointMsgUnjail) error {
	valAddr := mitotypes.EthAddress(event.ValAddr)

	// Check if validator exists
	validator, found := k.GetValidator(ctx, valAddr)
	if !found {
		return types.ErrValidatorNotFound
	}

	// Check if validator is jailed
	if !validator.Jailed {
		return errors.New("validator is not jailed", "validator", valAddr)
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
	valAddr := mitotypes.EthAddress(event.ValAddr)

	// Check if validator exists
	validator, found := k.GetValidator(ctx, valAddr)
	if !found {
		return types.ErrValidatorNotFound
	}

	// Convert the extra voting power to math.Int
	extraVotingPower := sdkmath.NewIntFromBigInt(event.ExtraVotingPowerGwei)

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
