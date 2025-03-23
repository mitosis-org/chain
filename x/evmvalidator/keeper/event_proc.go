package keeper

import (
	"context"
	stderrors "errors"

	sdkmath "cosmossdk.io/math"

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

	cachedValidatorEntrypointAddr     common.Address
	cachedValidatorEntrypointContract *bindings.ConsensusValidatorEntrypoint
)

// Name returns the name of the module
func (*Keeper) Name() string {
	return types.ModuleName
}

// FilterParams defines the matching EVM log events
func (k *Keeper) FilterParams(ctx context.Context) ([]common.Address, [][]common.Hash) {
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	contractAddr := k.GetValidatorEntrypointContractAddr(sdkCtx).Address()
	return []common.Address{contractAddr},
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
func (k *Keeper) processEvent(ctx sdk.Context, blockHash common.Hash, elog evmengtypes.EVMEvent) (error, bool) {
	ethlog, err := elog.ToEthLog()
	if err != nil {
		return err, false
	}

	contract, err := getValidatorEntrypointContract(ethlog.Address)
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
		event, err := contract.ParseMsgRegisterValidator(ethlog)
		if err != nil {
			return errors.Wrap(err, "parse MsgRegisterValidator"), false
		}
		if err, ignore := k.processRegisterValidator(ctx, event); err != nil {
			if !ignore {
				return errors.Wrap(err, "process MsgRegisterValidator"), false
			}

			if errFB := k.fallbackRegisterValidator(ctx, event); errFB != nil {
				return stderrors.Join(
					errors.Wrap(err, "process MsgRegisterValidator"),
					errors.Wrap(errFB, "fallback MsgRegisterValidator"),
				), false
			}

			k.Logger(ctx).Error("Processing failed but fallback succeeded",
				"name", eventName(elog),
				"height", ctx.BlockHeight(),
				"evmBlockHash", blockHash.Hex(),
				"evmLog", elog.String(),
				"err", err)
		}

	// Potential failure cases are:
	// - The validator does not exist (might be verified at the EVM contract level)
	// We must refund the collateral to the user through fallback logic if the primary logic fails.
	// The fallback logic must not fail due to its critical nature and should not fail because it's trivial.
	// Therefore, we raise an error if the fallback fails.
	case EventMsgDepositCollateral.ID:
		event, err := contract.ParseMsgDepositCollateral(ethlog)
		if err != nil {
			return errors.Wrap(err, "parse MsgDepositCollateral"), false
		}
		if err, ignore := k.processDepositCollateral(ctx, event); err != nil {
			if !ignore {
				return errors.Wrap(err, "process MsgDepositCollateral"), false
			}

			if errFB := k.fallbackDepositCollateral(ctx, event); errFB != nil {
				return stderrors.Join(
					errors.Wrap(err, "process MsgDepositCollateral"),
					errors.Wrap(errFB, "fallback MsgDepositCollateral"),
				), false
			}

			k.Logger(ctx).Error("Processing failed but fallback succeeded",
				"name", eventName(elog),
				"height", ctx.BlockHeight(),
				"evmBlockHash", blockHash.Hex(),
				"evmLog", elog.String(),
				"err", err)
		}

	// Potential failure cases are:
	// - The validator does not exist (might be verified at the EVM contract level)
	// - The withdrawal amount is greater than the validator's collateral (could be not verified at the EVM contract level)
	// Fortunately, this logic is not critical. Even if it fails, users won't lose money
	// and the state won't become corrupted. Therefore, we simply ignore errors when they occur.
	case EventMsgWithdrawCollateral.ID:
		event, err := contract.ParseMsgWithdrawCollateral(ethlog)
		if err != nil {
			return errors.Wrap(err, "parse MsgWithdrawCollateral"), false
		}
		if err, ignore := k.processWithdrawCollateral(ctx, event); err != nil {
			return errors.Wrap(err, "process MsgWithdrawCollateral"), ignore
		}

	// Potential failure cases are:
	// - The validator does not exist (might be verified at the EVM contract level)
	// - The validator is not jailed (could be not verified at the EVM contract level)
	// Fortunately, this logic is not critical. Even if it fails, users won't lose money
	// and the state won't become corrupted. Therefore, we simply ignore errors when they occur.
	case EventMsgUnjail.ID:
		event, err := contract.ParseMsgUnjail(ethlog)
		if err != nil {
			return errors.Wrap(err, "parse MsgUnjail"), false
		}
		if err, ignore := k.processUnjail(ctx, event); err != nil {
			return errors.Wrap(err, "process MsgUnjail"), ignore
		}

	// Potential failure cases are:
	// - The validator does not exist (might be verified at the EVM contract level)
	// Fortunately, this logic is not critical. Even if it fails, users won't lose money
	// and the state won't become corrupted. Therefore, we simply ignore errors when they occur.
	case EventMsgUpdateExtraVotingPower.ID:
		event, err := contract.ParseMsgUpdateExtraVotingPower(ethlog)
		if err != nil {
			return errors.Wrap(err, "parse MsgUpdateExtraVotingPower"), false
		}
		if err, ignore := k.processUpdateExtraVotingPower(ctx, event); err != nil {
			return errors.Wrap(err, "process MsgUpdateExtraVotingPower"), ignore
		}

	default:
		return errors.New("unknown event"), false
	}

	return nil, false
}

// processRegisterValidator processes MsgRegisterValidator event
// The second return value indicates whether it is okay to ignore the error
func (k *Keeper) processRegisterValidator(ctx sdk.Context, event *bindings.ConsensusValidatorEntrypointMsgRegisterValidator) (error, bool) {
	valAddr := mitotypes.EthAddress(event.ValAddr)

	// Convert the amount to math.Int
	collateral := sdkmath.NewIntFromBigInt(event.InitialCollateralAmountGwei)

	// Register the validator
	if err := k.registerValidator(ctx, valAddr, event.PubKey, collateral, sdkmath.LegacyZeroDec(), false); err != nil {
		ignore := errors.Is(err, types.ErrValidatorAlreadyExists) ||
			errors.Is(err, types.ErrInvalidPubKey)
		return errors.Wrap(err, "failed to register validator"), ignore
	}

	return nil, false
}

// fallbackRegisterValidator handles the case when the MsgRegisterValidator event fails to process
func (k *Keeper) fallbackRegisterValidator(ctx sdk.Context, event *bindings.ConsensusValidatorEntrypointMsgRegisterValidator) error {
	return k.evmEngKeeper.InsertWithdrawal(ctx, event.CollateralRefundAddr, event.InitialCollateralAmountGwei.Uint64())
}

// processDepositCollateral processes MsgDepositCollateral event
// The second return value indicates whether it is okay to ignore the error
func (k *Keeper) processDepositCollateral(ctx sdk.Context, event *bindings.ConsensusValidatorEntrypointMsgDepositCollateral) (error, bool) {
	valAddr := mitotypes.EthAddress(event.ValAddr)

	// Check if validator exists
	validator, found := k.GetValidator(ctx, valAddr)
	if !found {
		return types.ErrValidatorNotFound, true
	}

	// Convert the amount to math.Int
	amount := sdkmath.NewIntFromBigInt(event.AmountGwei)

	// Update validator's collateral
	k.depositCollateral(ctx, &validator, amount)

	return nil, false
}

// fallbackDepositCollateral handles the case when the MsgDepositCollateral event fails to process
func (k *Keeper) fallbackDepositCollateral(ctx sdk.Context, event *bindings.ConsensusValidatorEntrypointMsgDepositCollateral) error {
	return k.evmEngKeeper.InsertWithdrawal(ctx, event.CollateralRefundAddr, event.AmountGwei.Uint64())
}

// processWithdrawCollateral processes MsgWithdrawCollateral event
// The second return value indicates whether it is okay to ignore the error
func (k *Keeper) processWithdrawCollateral(ctx sdk.Context, event *bindings.ConsensusValidatorEntrypointMsgWithdrawCollateral) (error, bool) {
	valAddr := mitotypes.EthAddress(event.ValAddr)

	// Check if validator exists
	validator, found := k.GetValidator(ctx, valAddr)
	if !found {
		return types.ErrValidatorNotFound, true
	}

	// Create a withdrawal
	withdrawal := types.Withdrawal{
		ID:             0, // ID will be set later
		ValAddr:        valAddr,
		Amount:         event.AmountGwei.Uint64(),
		Receiver:       mitotypes.EthAddress(event.Receiver),
		MaturesAt:      event.MaturesAt.Int64(),
		CreationHeight: ctx.BlockHeight(),
	}

	// Request withdrawal
	if err := k.withdrawCollateral(ctx, &validator, &withdrawal); err != nil {
		ignore := errors.Is(err, types.ErrInsufficientCollateral)
		return errors.Wrap(err, "failed to withdraw collateral"), ignore
	}

	return nil, false
}

// processUnjail processes MsgUnjail event
// The second return value indicates whether it is okay to ignore the error
func (k *Keeper) processUnjail(ctx sdk.Context, event *bindings.ConsensusValidatorEntrypointMsgUnjail) (error, bool) {
	valAddr := mitotypes.EthAddress(event.ValAddr)

	// Check if validator exists
	validator, found := k.GetValidator(ctx, valAddr)
	if !found {
		return types.ErrValidatorNotFound, true
	}

	// Check if validator is jailed
	if !validator.Jailed {
		return errors.New("validator is not jailed", "validator", valAddr), true
	}

	// Get consensus address
	consAddr, err := validator.ConsAddr()
	if err != nil {
		return errors.Wrap(err, "failed to get consensus address"), false
	}

	// unjail validator through slashing keeper
	if err = k.slashingKeeper.UnjailFromConsAddr(ctx, consAddr); err != nil {
		return errors.Wrap(err, "failed to unjail validator"), false
	}

	return nil, false
}

// processUpdateExtraVotingPower processes MsgUpdateExtraVotingPower event
// The second return value indicates whether it is okay to ignore the error
func (k *Keeper) processUpdateExtraVotingPower(ctx sdk.Context, event *bindings.ConsensusValidatorEntrypointMsgUpdateExtraVotingPower) (error, bool) {
	valAddr := mitotypes.EthAddress(event.ValAddr)

	// Check if validator exists
	validator, found := k.GetValidator(ctx, valAddr)
	if !found {
		return types.ErrValidatorNotFound, true
	}

	// Convert the extra voting power
	extraVotingPower := sdkmath.LegacyNewDecFromBigInt(event.ExtraVotingPowerWei).QuoInt(types.VotingPowerReductionForWei) // 1 wei = 1 / 1e18

	// Update extra voting power
	if err := k.updateExtraVotingPower(ctx, &validator, extraVotingPower); err != nil {
		return errors.Wrap(err, "failed to update extra voting power"), false
	}

	return nil, false
}

func getValidatorEntrypointContract(addr common.Address) (*bindings.ConsensusValidatorEntrypoint, error) {
	if addr != cachedValidatorEntrypointAddr {
		var err error

		cachedValidatorEntrypointContract, err = bindings.NewConsensusValidatorEntrypoint(addr, nil)
		if err != nil {
			return nil, errors.Wrap(err, "failed to create consensus validator entrypoint contract")
		}
		cachedValidatorEntrypointAddr = addr
	}

	return cachedValidatorEntrypointContract, nil
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
