package keeper

import (
	"encoding/hex"
	"fmt"
	"strconv"
	"time"

	sdkmath "cosmossdk.io/math"

	sdk "github.com/cosmos/cosmos-sdk/types"
	mitotypes "github.com/mitosis-org/chain/types"
	"github.com/mitosis-org/chain/x/evmvalidator/types"
	"github.com/omni-network/omni/lib/errors"
)

func (k Keeper) RegisterValidator(
	ctx sdk.Context,
	valAddr mitotypes.EthAddress,
	pubkey []byte,
	collateral sdkmath.Uint,
	extraVotingPower sdkmath.LegacyDec,
	jailed bool,
) error {
	// Validate pubkey with address
	err := types.ValidatePubkeyWithEthAddress(pubkey, valAddr)
	if err != nil {
		return errors.Wrap(err, "failed to validate pubkey with address")
	}

	// Check if validator already exists
	if k.HasValidator(ctx, valAddr) {
		return errors.Wrap(types.ErrValidatorAlreadyExists, valAddr.String())
	}

	// Ensure extra voting power is non-negative
	if extraVotingPower.IsNegative() {
		return errors.New("extra voting power cannot be negative: %s", extraVotingPower.String())
	}

	// Create a new validator
	validator := types.Validator{
		Addr:             valAddr,
		Pubkey:           pubkey,
		Collateral:       collateral,
		ExtraVotingPower: extraVotingPower,
		VotingPower:      0, // will be calculated later
		Jailed:           jailed,
		Bonded:           false,
	}

	// Get consensus public key and address
	consPubKey, err := validator.ConsPubKey()
	if err != nil {
		return errors.Wrap(err, "failed to get consensus public key")
	}
	consAddr, err := validator.ConsAddr()
	if err != nil {
		return errors.Wrap(err, "failed to get consensus address")
	}

	// Set the validator in state
	k.SetValidator(ctx, validator)

	// Set the validator in consensus address index
	k.SetValidatorByConsAddr(ctx, consAddr, validator.Addr)

	// Call slashing hook
	if err = k.slashingKeeper.AfterValidatorCreated(ctx, consPubKey); err != nil {
		return errors.Wrap(err, "failed to call AfterValidatorCreated hook")
	}

	// Emit event
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeRegisterValidator,
			sdk.NewAttribute(types.AttributeKeyValAddr, valAddr.String()),
			sdk.NewAttribute(types.AttributeKeyPubkey, hex.EncodeToString(pubkey)),
			sdk.NewAttribute(types.AttributeKeyCollateral, collateral.String()),
			sdk.NewAttribute(types.AttributeKeyExtraVotingPower, extraVotingPower.String()),
			sdk.NewAttribute(types.AttributeKeyJailed, strconv.FormatBool(jailed)),
		),
	)

	// Update the validator state to calculate voting power
	k.UpdateValidatorState(ctx, &validator, "register validator")

	return nil
}

func (k Keeper) DepositCollateral(ctx sdk.Context, validator *types.Validator, amount sdkmath.Uint) error {
	// Update validator's collateral
	validator.Collateral = validator.Collateral.Add(amount)

	// Update the validator state
	k.UpdateValidatorState(ctx, validator, "deposit collateral")

	// Emit event
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeDepositCollateral,
			sdk.NewAttribute(types.AttributeKeyValAddr, validator.Addr.String()),
			sdk.NewAttribute(types.AttributeKeyAmount, amount.String()),
			sdk.NewAttribute(types.AttributeKeyCollateral, validator.Collateral.String()),
			sdk.NewAttribute(types.AttributeKeyVotingPower, fmt.Sprintf("%d", validator.VotingPower)),
		),
	)

	return nil
}

func (k Keeper) WithdrawCollateral(ctx sdk.Context, validator *types.Validator, withdrawal *types.Withdrawal) error {
	amount := sdkmath.NewUint(withdrawal.Amount)

	if amount.IsZero() {
		return nil // nothing to withdraw
	}

	// Ensure validator has enough collateral
	if validator.Collateral.LT(amount) {
		return errors.Wrap(types.ErrInsufficientCollateral,
			"validator does not have enough collateral to withdraw",
			"collateral", validator.Collateral.String(), "amount", amount.String(),
		)
	}

	// Update validator's collateral
	validator.Collateral = validator.Collateral.Sub(amount)

	// Update the validator state
	k.UpdateValidatorState(ctx, validator, "withdraw collateral")

	// Add a new withdrawal
	k.AddNewWithdrawalWithNextID(ctx, withdrawal)

	// Emit events
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeWithdrawCollateral,
			sdk.NewAttribute(types.AttributeKeyWithdrawalID, fmt.Sprintf("%d", withdrawal.ID)),
			sdk.NewAttribute(types.AttributeKeyValAddr, withdrawal.ValAddr.String()),
			sdk.NewAttribute(types.AttributeKeyAmount, amount.String()),
			sdk.NewAttribute(types.AttributeKeyReceiver, withdrawal.Receiver.String()),
			sdk.NewAttribute(types.AttributeKeyMaturesAt, time.Unix(withdrawal.MaturesAt, 0).String()),
		),
	)

	return nil
}

// Slash_ slashes a validator's collateral by a fraction
func (k Keeper) Slash_(ctx sdk.Context, validator *types.Validator, infractionHeight int64, power int64, slashFraction sdkmath.LegacyDec) (sdkmath.Uint, error) {
	currentTime := ctx.BlockTime().Unix()

	if slashFraction.IsNegative() {
		return sdkmath.ZeroUint(), fmt.Errorf("attempted to Slash_ with a negative Slash_ factor: %s", slashFraction.String())
	}

	// Calculate the collateral amount to Slash_
	targetSlashAmount := sdkmath.NewUintFromBigInt(
		sdkmath.LegacyNewDec(power).
			MulInt(types.VotingPowerReductionForGwei).
			Mul(slashFraction).
			TruncateInt().
			BigInt(),
	)

	remainingSlashAmount := targetSlashAmount

	// Slash the not matured withdrawals from the oldest to the newest
	// NOTE: The implementation differs from x/staking. In the case of x/staking, slashing is applied
	// proportionally to each unbonding entry, and only to entries that contributed at the infraction height.
	// However, in x/evmvalidator, since the amounts delegated by users are not subject to slashing,
	// we thought it would be acceptable to use a simpler policy.
	// Therefore, we decided to apply slashing sequentially from the oldest withdrawal up to the collateral.
	k.IterateWithdrawalsForValidator(ctx, validator.Addr, func(w types.Withdrawal) bool {
		if remainingSlashAmount.IsZero() {
			return true
		}

		// If withdrawal is matured, it is not subject to slashing
		if w.MaturesAt <= currentTime {
			return false
		}

		// Slash the withdrawal
		withdrawalAmount := sdkmath.NewUint(w.Amount)

		if withdrawalAmount.GTE(remainingSlashAmount) {
			w.Amount = withdrawalAmount.Sub(remainingSlashAmount).Uint64()
			remainingSlashAmount = sdkmath.ZeroUint()
			k.SetWithdrawal(ctx, w) // overwrite the withdrawal
		} else {
			remainingSlashAmount = remainingSlashAmount.Sub(withdrawalAmount)
			k.DeleteWithdrawal(ctx, w)
		}

		return false
	})

	// Slash the collateral
	if validator.Collateral.GTE(remainingSlashAmount) {
		validator.Collateral = validator.Collateral.Sub(remainingSlashAmount)
		remainingSlashAmount = sdkmath.ZeroUint()
	} else {
		k.Logger(ctx).Error("Slash amount exceeds validator's collateral",
			"validator", validator.Addr.String(),
			"slashAmount", remainingSlashAmount.String(),
			"collateral", validator.Collateral.String(),
		)
		remainingSlashAmount = remainingSlashAmount.Sub(validator.Collateral)
		validator.Collateral = sdkmath.ZeroUint()
	}

	actualSlashAmount := targetSlashAmount.Sub(remainingSlashAmount)

	// Update the validator state
	k.UpdateValidatorState(ctx, validator, "slash")

	// Emit event
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeSlashValidator,
			sdk.NewAttribute(types.AttributeKeyValAddr, validator.Addr.String()),
			sdk.NewAttribute(types.AttributeKeyAmount, actualSlashAmount.String()),
			sdk.NewAttribute(types.AttributeKeySlashFraction, slashFraction.String()),
			sdk.NewAttribute(types.AttributeKeyInfractionHeight, fmt.Sprintf("%d", infractionHeight)),
			sdk.NewAttribute(types.AttributeKeyInfractionPower, fmt.Sprintf("%d", power)),
		),
	)

	return actualSlashAmount, nil
}

// Jail_ jails a validator
func (k Keeper) Jail_(ctx sdk.Context, validator *types.Validator, reason string) {
	if validator.Jailed {
		return // already jailed
	}

	// Update the validator in state
	validator.Jailed = true
	k.SetValidator(ctx, *validator)

	// Delete the validator in power index
	k.DeleteValidatorByPowerIndex(ctx, validator.VotingPower, validator.Addr)

	// Emit event
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeJailValidator,
			sdk.NewAttribute(types.AttributeKeyValAddr, validator.Addr.String()),
			sdk.NewAttribute(types.AttributeKeyReason, reason),
		),
	)
}

// Unjail_ unjails a validator
func (k Keeper) Unjail_(ctx sdk.Context, validator *types.Validator) error {
	if !validator.Jailed {
		return nil // already unjailed
	}

	// Check if voting power meets minimum requirement
	params := k.GetParams(ctx)
	if validator.VotingPower < params.MinVotingPower {
		return errors.Wrap(types.ErrInvalidVotingPower, "voting power below minimum requirement")
	}

	// Update the validator in state
	validator.Jailed = false
	k.SetValidator(ctx, *validator)

	// Set the validator back in power index
	k.SetValidatorByPowerIndex(ctx, validator.VotingPower, validator.Addr)

	// Emit event
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeUnjailValidator,
			sdk.NewAttribute(types.AttributeKeyValAddr, validator.Addr.String()),
		),
	)

	return nil
}

func (k Keeper) UpdateExtraVotingPower(ctx sdk.Context, validator *types.Validator, extraVotingPower sdkmath.LegacyDec) error {
	if extraVotingPower.IsNegative() {
		return errors.New("extra voting power cannot be negative: %s", extraVotingPower.String())
	}

	// Update validator's extra voting power
	oldExtraVotingPower := validator.ExtraVotingPower
	validator.ExtraVotingPower = extraVotingPower

	// Update the validator state
	k.UpdateValidatorState(ctx, validator, "update extra voting power")

	// Emit events
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeUpdateExtraVotingPower,
			sdk.NewAttribute(types.AttributeKeyValAddr, validator.Addr.String()),
			sdk.NewAttribute(types.AttributeKeyOldExtraVotingPower, oldExtraVotingPower.String()),
			sdk.NewAttribute(types.AttributeKeyExtraVotingPower, extraVotingPower.String()),
		),
	)

	return nil
}

func (k Keeper) UpdateValidatorState(ctx sdk.Context, validator *types.Validator, context string) {
	params := k.GetParams(ctx)
	oldVotingPower := validator.VotingPower

	// Recompute voting power
	validator.VotingPower = validator.ComputeVotingPower(params.MaxLeverageRatio)

	// Update the validator in state
	k.SetValidator(ctx, *validator)

	// Update the validator in power index
	if !validator.Jailed {
		k.DeleteValidatorByPowerIndex(ctx, oldVotingPower, validator.Addr)
		k.SetValidatorByPowerIndex(ctx, validator.VotingPower, validator.Addr)
	}

	// Emit update event if voting power changed
	if validator.VotingPower != oldVotingPower {
		ctx.EventManager().EmitEvent(
			sdk.NewEvent(
				types.EventTypeUpdateVotingPower,
				sdk.NewAttribute(types.AttributeKeyValAddr, validator.Addr.String()),
				sdk.NewAttribute(types.AttributeKeyOldVotingPower, fmt.Sprintf("%d", oldVotingPower)),
				sdk.NewAttribute(types.AttributeKeyVotingPower, fmt.Sprintf("%d", validator.VotingPower)),
			),
		)
	}

	// Check min voting power requirement
	if validator.VotingPower < params.MinVotingPower {
		k.Jail_(ctx, validator, fmt.Sprintf("min voting power requirement is not met: %s", context))
	}
}
