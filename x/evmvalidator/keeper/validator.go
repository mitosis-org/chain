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
	initialCollateralOwner mitotypes.EthAddress,
	initialCollateral sdkmath.Uint,
	extraVotingPower sdkmath.Uint,
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

	// Calculate initial shares
	initialCollateralShares := types.CalculateCollateralSharesForDeposit(initialCollateral, sdkmath.ZeroUint(), initialCollateral)

	// Create a new validator
	validator := types.Validator{
		Addr:             valAddr,
		Pubkey:           pubkey,
		Collateral:       initialCollateral,
		CollateralShares: initialCollateralShares,
		ExtraVotingPower: extraVotingPower,
		VotingPower:      0, // will be calculated later
		Jailed:           jailed,
		Bonded:           false,
	}

	// Get consensus public key and address
	consPubKey := validator.MustConsPubKey()
	consAddr := validator.MustConsAddr()

	// Set the validator in state
	k.SetValidator(ctx, validator)

	// Set the validator in consensus address index
	k.SetValidatorByConsAddr(ctx, consAddr, validator.Addr)

	// Call slashing hook
	if err = k.slashingKeeper.AfterValidatorCreated(ctx, consPubKey); err != nil {
		return errors.Wrap(err, "failed to call AfterValidatorCreated hook")
	}

	// Create a new ownership record
	ownership := types.CollateralOwnership{
		ValAddr:        valAddr,
		Owner:          initialCollateralOwner,
		Shares:         initialCollateralShares,
		CreationHeight: ctx.BlockHeight(),
	}

	// Set the ownership record
	k.SetCollateralOwnership(ctx, ownership)

	// Emit event
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeRegisterValidator,
			sdk.NewAttribute(types.AttributeKeyValAddr, valAddr.String()),
			sdk.NewAttribute(types.AttributeKeyPubkey, hex.EncodeToString(validator.Pubkey)),
			sdk.NewAttribute(types.AttributeKeyCollateralOwner, initialCollateralOwner.String()),
			sdk.NewAttribute(types.AttributeKeyCollateral, validator.Collateral.String()),
			sdk.NewAttribute(types.AttributeKeyCollateralShares, validator.CollateralShares.String()),
			sdk.NewAttribute(types.AttributeKeyExtraVotingPower, validator.ExtraVotingPower.String()),
			sdk.NewAttribute(types.AttributeKeyJailed, strconv.FormatBool(validator.Jailed)),
		),
	)

	k.Logger(ctx).Debug("🆕 Validator Registered",
		"height", ctx.BlockHeight(),
		"addr", valAddr.String(),
		"consAddr", consAddr.String(),
		"pubkey", hex.EncodeToString(validator.Pubkey),
		"collateralOwner", initialCollateralOwner.String(),
		"collateral", validator.Collateral,
		"collateralShares", validator.CollateralShares,
		"extraVotingPower", validator.ExtraVotingPower,
		"jailed", validator.Jailed,
	)

	// Update the validator state to calculate voting power
	k.UpdateValidatorState(ctx, &validator, "register validator")

	return nil
}

func (k Keeper) DepositCollateral(ctx sdk.Context, validator *types.Validator, owner mitotypes.EthAddress, amount sdkmath.Uint) {
	if amount.IsZero() {
		return // nothing to deposit
	}

	// Calculate shares for this deposit
	shares := types.CalculateCollateralSharesForDeposit(validator.Collateral, validator.CollateralShares, amount)

	// Update validator's collateral and shares
	validator.Collateral = validator.Collateral.Add(amount)
	validator.CollateralShares = validator.CollateralShares.Add(shares)

	// Update or create the ownership record
	ownership, found := k.GetCollateralOwnership(ctx, validator.Addr, owner)
	if !found {
		ownership = types.CollateralOwnership{
			ValAddr:        validator.Addr,
			Owner:          owner,
			Shares:         shares,
			CreationHeight: ctx.BlockHeight(),
		}
	} else {
		ownership.Shares = ownership.Shares.Add(shares)
	}

	// Save the ownership record
	k.SetCollateralOwnership(ctx, ownership)

	// Emit event
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeDepositCollateral,
			sdk.NewAttribute(types.AttributeKeyValAddr, validator.Addr.String()),
			sdk.NewAttribute(types.AttributeKeyCollateralOwner, owner.String()),
			sdk.NewAttribute(types.AttributeKeyAmount, amount.String()),
			sdk.NewAttribute(types.AttributeKeyShares, shares.String()),
			sdk.NewAttribute(types.AttributeKeyCollateral, validator.Collateral.String()),
		),
	)

	k.Logger(ctx).Debug("💵 Validator Collateral Deposited",
		"height", ctx.BlockHeight(),
		"validator", validator.Addr.String(),
		"collateralOwner", owner.String(),
		"amount", amount.String(),
		"shares", shares.String(),
		"collateral", validator.Collateral.String(),
	)

	// Update the validator state
	k.UpdateValidatorState(ctx, validator, "deposit collateral")
}

func (k Keeper) WithdrawCollateral(
	ctx sdk.Context,
	validator *types.Validator,
	owner mitotypes.EthAddress,
	withdrawal *types.Withdrawal,
) error {
	amount := sdkmath.NewUint(withdrawal.Amount)

	if amount.IsZero() {
		return nil // nothing to withdraw
	}

	// Get the ownership record
	ownership, found := k.GetCollateralOwnership(ctx, validator.Addr, owner)
	if !found {
		return errors.Wrap(types.ErrInsufficientCollateral,
			"collateral owner does not have collateral for this validator",
			"validator", validator.Addr.String(), "collateralOwner", owner.String(),
		)
	}

	// Calculate how many shares to withdraw
	sharesToWithdraw := types.CalculateCollateralSharesForWithdrawal(validator.Collateral, validator.CollateralShares, amount)

	// Ensure owner has enough shares
	if ownership.Shares.LT(sharesToWithdraw) {
		return errors.Wrap(types.ErrInsufficientCollateral,
			"collateral owner does not have enough collateral to withdraw",
			"validator", validator.Addr.String(), "collateralOwner", owner.String(),
			"shares", ownership.Shares.String(), "requiredShares", sharesToWithdraw.String(),
		)
	}

	// Ensure validator has enough total collateral
	if validator.Collateral.LT(amount) {
		return errors.Wrap(types.ErrInsufficientCollateral,
			"validator does not have enough collateral to withdraw",
			"collateral", validator.Collateral.String(), "amount", amount.String(),
		)
	}

	// Update validator's collateral and shares
	validator.Collateral = validator.Collateral.Sub(amount)
	validator.CollateralShares = validator.CollateralShares.Sub(sharesToWithdraw)

	// Update ownership record
	ownership.Shares = ownership.Shares.Sub(sharesToWithdraw)

	// If no shares left, delete the ownership record, otherwise update it
	if ownership.Shares.IsZero() {
		k.DeleteCollateralOwnership(ctx, validator.Addr, owner)
	} else {
		k.SetCollateralOwnership(ctx, ownership)
	}

	// Add a new withdrawal
	k.AddNewWithdrawalWithNextID(ctx, withdrawal)

	// Emit events
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeWithdrawCollateral,
			sdk.NewAttribute(types.AttributeKeyWithdrawalID, fmt.Sprintf("%d", withdrawal.ID)),
			sdk.NewAttribute(types.AttributeKeyValAddr, withdrawal.ValAddr.String()),
			sdk.NewAttribute(types.AttributeKeyCollateralOwner, owner.String()),
			sdk.NewAttribute(types.AttributeKeyAmount, amount.String()),
			sdk.NewAttribute(types.AttributeKeyShares, sharesToWithdraw.String()),
			sdk.NewAttribute(types.AttributeKeyReceiver, withdrawal.Receiver.String()),
			sdk.NewAttribute(types.AttributeKeyMaturesAt, time.Unix(withdrawal.MaturesAt, 0).String()),
		),
	)

	k.Logger(ctx).Debug("💸 Validator Collateral Withdrawal Requested",
		"height", ctx.BlockHeight(),
		"validator", validator.Addr.String(),
		"collateralOwner", owner.String(),
		"withdrawalID", withdrawal.ID,
		"amount", amount.String(),
		"shares", sharesToWithdraw.String(),
		"receiver", withdrawal.Receiver.String(),
		"maturesAt", time.Unix(withdrawal.MaturesAt, 0),
	)

	// Update the validator state
	k.UpdateValidatorState(ctx, validator, "withdraw collateral")

	return nil
}

// TransferCollateralOwnership transfers collateral ownership from one owner to another
func (k Keeper) TransferCollateralOwnership(
	ctx sdk.Context,
	validator *types.Validator,
	prevOwnership types.CollateralOwnership,
	newOwner mitotypes.EthAddress,
) {
	// If the previous owner is the same as the new owner, do nothing
	if prevOwnership.Owner.Equal(newOwner) {
		return
	}

	// Transfer all shares from the previous owner to the new owner
	sharesToTransfer := prevOwnership.Shares

	// Get or create ownership record for the new owner
	newOwnership, found := k.GetCollateralOwnership(ctx, validator.Addr, newOwner)
	if !found {
		newOwnership = types.CollateralOwnership{
			ValAddr:        validator.Addr,
			Owner:          newOwner,
			Shares:         sharesToTransfer,
			CreationHeight: ctx.BlockHeight(),
		}
	} else {
		newOwnership.Shares = newOwnership.Shares.Add(sharesToTransfer)
	}

	// Save the new ownership record
	k.SetCollateralOwnership(ctx, newOwnership)

	// Remove the previous owner's record
	k.DeleteCollateralOwnership(ctx, validator.Addr, prevOwnership.Owner)

	// Emit event
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeTransferCollateralOwnership,
			sdk.NewAttribute(types.AttributeKeyValAddr, validator.Addr.String()),
			sdk.NewAttribute(types.AttributeKeyCollateralOwner, prevOwnership.Owner.String()),
			sdk.NewAttribute(types.AttributeKeyCollateralNewOwner, newOwner.String()),
			sdk.NewAttribute(types.AttributeKeyShares, sharesToTransfer.String()),
		),
	)

	k.Logger(ctx).Debug("🔄 Validator Collateral Ownership Transferred",
		"height", ctx.BlockHeight(),
		"validator", validator.Addr.String(),
		"prevOwner", prevOwnership.Owner.String(),
		"newOwner", newOwner.String(),
		"shares", sharesToTransfer.String(),
	)
}

// Slash_ slashes a validator's collateral by a fraction
func (k Keeper) Slash_(ctx sdk.Context, validator *types.Validator, infractionHeight int64, power int64, slashFraction sdkmath.LegacyDec) (sdkmath.Uint, error) {
	currentTime := ctx.BlockTime().Unix()

	// Ensure power and slash fraction are non-negative
	if power < 0 {
		return sdkmath.ZeroUint(), fmt.Errorf("attempted to slash with a negative power: %d", power)
	}
	if slashFraction.IsNegative() {
		return sdkmath.ZeroUint(), fmt.Errorf("attempted to slash with a negative slash fraction: %s", slashFraction.String())
	}

	// Calculate the collateral amount to slash
	targetSlashAmount := sdkmath.NewUintFromBigInt(
		sdkmath.LegacyNewDec(power).
			MulInt(types.VotingPowerReduction).
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

		if withdrawalAmount.GT(remainingSlashAmount) {
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

	k.Logger(ctx).Debug("💥 Validator Slashed",
		"height", ctx.BlockHeight(),
		"validator", validator.Addr.String(),
		"slashAmount", actualSlashAmount.String(),
		"slashFraction", slashFraction.String(),
		"infractionHeight", infractionHeight,
		"infractionPower", power,
	)

	// Update the validator state
	k.UpdateValidatorState(ctx, validator, "slash")

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

	k.Logger(ctx).Debug("🔒 Validator Jailed",
		"height", ctx.BlockHeight(),
		"validator", validator.Addr.String(),
		"reason", reason,
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

	k.Logger(ctx).Debug("🔓 Validator Unjailed",
		"height", ctx.BlockHeight(),
		"validator", validator.Addr.String(),
	)

	return nil
}

func (k Keeper) UpdateExtraVotingPower(ctx sdk.Context, validator *types.Validator, extraVotingPower sdkmath.Uint) {
	// Update validator's extra voting power
	oldExtraVotingPower := validator.ExtraVotingPower
	validator.ExtraVotingPower = extraVotingPower

	// Emit event
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeUpdateExtraVotingPower,
			sdk.NewAttribute(types.AttributeKeyValAddr, validator.Addr.String()),
			sdk.NewAttribute(types.AttributeKeyOldExtraVotingPower, oldExtraVotingPower.String()),
			sdk.NewAttribute(types.AttributeKeyExtraVotingPower, extraVotingPower.String()),
		),
	)

	k.Logger(ctx).Debug("💳 Validator Extra Voting Power Updated",
		"height", ctx.BlockHeight(),
		"validator", validator.Addr.String(),
		"oldExtraVotingPower", oldExtraVotingPower,
		"newExtraVotingPower", extraVotingPower,
	)

	// Update the validator state
	k.UpdateValidatorState(ctx, validator, "update extra voting power")
}

func (k Keeper) UpdateValidatorState(ctx sdk.Context, validator *types.Validator, context string) {
	params := k.GetParams(ctx)
	oldVotingPower := validator.VotingPower

	// Recompute voting power
	validator.VotingPower = validator.ComputeVotingPower(params.MaxLeverageRatio)

	// Update the validator in state
	k.SetValidator(ctx, *validator)

	// Update the validator in power index
	k.DeleteValidatorByPowerIndex(ctx, oldVotingPower, validator.Addr)
	if !validator.Jailed {
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

		k.Logger(ctx).Debug("🔋 Validator Voting Power Updated",
			"height", ctx.BlockHeight(),
			"context", context,
			"validator", validator.Addr.String(),
			"oldVotingPower", oldVotingPower,
			"newVotingPower", validator.VotingPower,
		)
	}

	// Check min voting power requirement
	if validator.VotingPower < params.MinVotingPower {
		k.Jail_(ctx, validator, fmt.Sprintf("min voting power requirement is not met: %s", context))
	}
}
