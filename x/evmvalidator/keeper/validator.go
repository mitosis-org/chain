package keeper

import (
	sdkmath "cosmossdk.io/math"
	"encoding/hex"
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	mitotypes "github.com/mitosis-org/chain/types"
	"github.com/mitosis-org/chain/x/evmvalidator/types"
	"github.com/omni-network/omni/lib/errors"
	"strconv"
	"time"
)

func (k Keeper) registerValidator(
	ctx sdk.Context,
	valAddr mitotypes.EthAddress,
	pubkey []byte,
	collateral sdkmath.Int,
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

	// Create a new validator
	validator := types.Validator{
		Addr:             valAddr,
		Pubkey:           pubkey,
		Collateral:       collateral,
		ExtraVotingPower: extraVotingPower,
		VotingPower:      0, // will be computed later
		Jailed:           jailed,
		Bonded:           false,
	}

	// Compute voting power
	params := k.GetParams(ctx)
	validator.VotingPower = validator.ComputeVotingPower(params.MaxLeverageRatio)

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

	// Set the validator in power index
	if !validator.Jailed {
		k.SetValidatorByPowerIndex(ctx, validator.VotingPower, validator.Addr)
	}

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
			sdk.NewAttribute(types.AttributeKeyVotingPower, fmt.Sprintf("%d", validator.VotingPower)),
			sdk.NewAttribute(types.AttributeKeyJailed, strconv.FormatBool(jailed)),
		),
	)

	// If min voting power requirement is not met, jail the validator
	if validator.VotingPower < params.MinVotingPower {
		if err := k.jail(ctx, &validator, "min voting power requirement is not met during registration"); err != nil {
			return errors.Wrap(err, "failed to jail validator")
		}
	}

	return nil
}

func (k Keeper) depositCollateral(ctx sdk.Context, validator *types.Validator, amount sdkmath.Int) {
	// Update validator's collateral
	validator.Collateral = validator.Collateral.Add(amount)

	// Recompute voting power
	params := k.GetParams(ctx)
	oldVotingPower := validator.VotingPower
	validator.VotingPower = validator.ComputeVotingPower(params.MaxLeverageRatio)

	// Update the validator in state
	k.SetValidator(ctx, *validator)

	// Update the validator in power index
	if !validator.Jailed {
		k.DeleteValidatorByPowerIndex(ctx, validator.VotingPower, validator.Addr)
		k.SetValidatorByPowerIndex(ctx, validator.VotingPower, validator.Addr)
	}

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

	// If voting power changed, emit update event
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
}

func (k Keeper) withdrawCollateral(ctx sdk.Context, validator *types.Validator, withdrawal *types.Withdrawal) error {
	amount := sdkmath.NewIntFromUint64(withdrawal.Amount)

	// Ensure validator has enough collateral
	if validator.Collateral.LT(amount) {
		return errors.Wrap(types.ErrInsufficientCollateral,
			"validator does not have enough collateral to withdraw",
			"collateral", validator.Collateral.String(), "amount", amount.String(),
		)
	}

	// Add a new withdrawal
	k.AddNewWithdrawalWithNextID(ctx, withdrawal)

	// Update validator's collateral (immediately reduce to prevent multiple withdrawals)
	validator.Collateral = validator.Collateral.Sub(amount)

	// Recompute voting power
	params := k.GetParams(ctx)
	oldVotingPower := validator.VotingPower
	validator.VotingPower = validator.ComputeVotingPower(params.MaxLeverageRatio)

	// Update the validator in state
	k.SetValidator(ctx, *validator)

	// Update the validator in power index
	if !validator.Jailed {
		k.DeleteValidatorByPowerIndex(ctx, oldVotingPower, validator.Addr)
		k.SetValidatorByPowerIndex(ctx, validator.VotingPower, validator.Addr)
	}

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

	// If voting power changed, emit update event
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

	// If min voting power requirement is not met, jail the validator
	if validator.VotingPower < params.MinVotingPower {
		if err := k.jail(ctx, validator, "min voting power requirement is not met due to withdrawal"); err != nil {
			return errors.Wrap(err, "failed to jail validator")
		}
	}

	return nil
}

// slash slashes a validator's collateral by a fraction
func (k Keeper) slash(ctx sdk.Context, consAddr sdk.ConsAddress, infractionHeight int64, power int64, slashFraction sdkmath.LegacyDec) (sdkmath.Int, error) {
	currentTime := ctx.BlockTime().Unix()

	// Find the validator by consensus address
	validator, found := k.GetValidatorByConsAddr(ctx, consAddr)
	if !found {
		return sdkmath.ZeroInt(), errors.Wrap(types.ErrValidatorNotFound, consAddr.String())
	}

	if slashFraction.IsNegative() {
		return sdkmath.ZeroInt(), fmt.Errorf("attempted to slash with a negative slash factor: %s", slashFraction.String())
	}

	// Calculate the collateral amount to slash
	targetSlashAmount := sdkmath.LegacyNewDec(power).MulInt(types.VotingPowerReductionForGwei).Mul(slashFraction).TruncateInt()

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

		// slash the withdrawal
		withdrawalAmount := sdkmath.NewIntFromUint64(w.Amount)

		if withdrawalAmount.GTE(remainingSlashAmount) {
			w.Amount = withdrawalAmount.Sub(remainingSlashAmount).Uint64()
			remainingSlashAmount = sdkmath.ZeroInt()
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
		remainingSlashAmount = sdkmath.ZeroInt()
	} else {
		k.Logger(ctx).Error("Slash amount exceeds validator's collateral",
			"validator", validator.Addr.String(),
			"slashAmount", remainingSlashAmount.String(),
			"collateral", validator.Collateral.String(),
		)
		remainingSlashAmount = remainingSlashAmount.Sub(validator.Collateral)
		validator.Collateral = sdkmath.ZeroInt()
	}

	actualSlashAmount := targetSlashAmount.Sub(remainingSlashAmount)

	// Recompute voting power
	params := k.GetParams(ctx)
	oldVotingPower := validator.VotingPower
	validator.VotingPower = validator.ComputeVotingPower(params.MaxLeverageRatio)

	// Update the validator in state
	k.SetValidator(ctx, validator)

	// Update the validator in power index
	if !validator.Jailed {
		k.DeleteValidatorByPowerIndex(ctx, oldVotingPower, validator.Addr)
		k.SetValidatorByPowerIndex(ctx, validator.VotingPower, validator.Addr)
	}

	// Emit events
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

	// If voting power changed, emit update event
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

	// If min voting power requirement is not met, jail the validator
	if validator.VotingPower < params.MinVotingPower {
		if err := k.jail(ctx, &validator, "min voting power requirement is not met due to slashing"); err != nil {
			return sdkmath.ZeroInt(), errors.Wrap(err, "failed to jail validator")
		}
	}

	return actualSlashAmount, nil
}

// jail jails a validator
func (k Keeper) jail(ctx sdk.Context, validator *types.Validator, reason string) error {
	if validator.Jailed {
		return nil // already jailed
	}

	validator.Jailed = true
	oldVotingPower := validator.VotingPower

	// Update the validator in state
	k.SetValidator(ctx, *validator)

	// Delete the validator in power index
	k.DeleteValidatorByPowerIndex(ctx, oldVotingPower, validator.Addr)

	// Emit event
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeJailValidator,
			sdk.NewAttribute(types.AttributeKeyValAddr, validator.Addr.String()),
			sdk.NewAttribute(types.AttributeKeyReason, reason),
		),
	)

	return nil
}

// unjail unjails a validator
func (k Keeper) unjail(ctx sdk.Context, validator *types.Validator) error {
	if !validator.Jailed {
		return nil // already unjailed
	}

	// Check if voting power meets minimum requirement
	params := k.GetParams(ctx)
	if validator.VotingPower < params.MinVotingPower {
		return errors.Wrap(types.ErrInvalidVotingPower, "voting power below minimum requirement")
	}

	validator.Jailed = false

	// Update the validator in state
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

func (k Keeper) updateExtraVotingPower(ctx sdk.Context, validator *types.Validator, extraVotingPower sdkmath.LegacyDec) error {
	// Update validator's extra voting power
	oldVotingPower := validator.VotingPower
	oldExtraVotingPower := validator.ExtraVotingPower
	validator.ExtraVotingPower = extraVotingPower

	// Recompute voting power
	params := k.GetParams(ctx)
	validator.VotingPower = validator.ComputeVotingPower(params.MaxLeverageRatio)

	// Update the validator in state
	k.SetValidator(ctx, *validator)

	// Update the validator in power index
	if !validator.Jailed {
		k.DeleteValidatorByPowerIndex(ctx, oldVotingPower, validator.Addr)
		k.SetValidatorByPowerIndex(ctx, validator.VotingPower, validator.Addr)
	}

	// Emit events
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeUpdateExtraVotingPower,
			sdk.NewAttribute(types.AttributeKeyValAddr, validator.Addr.String()),
			sdk.NewAttribute(types.AttributeKeyOldExtraVotingPower, oldExtraVotingPower.String()),
			sdk.NewAttribute(types.AttributeKeyExtraVotingPower, extraVotingPower.String()),
		),
	)

	// If voting power changed, emit update event
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

	// If min voting power requirement is not met, jail the validator
	if validator.VotingPower < params.MinVotingPower {
		if err := k.jail(ctx, validator, "min voting power requirement is not met due to extra voting power update"); err != nil {
			return errors.Wrap(err, "failed to jail validator")
		}
	}

	return nil
}
