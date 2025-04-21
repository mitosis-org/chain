package types

import (
	"cosmossdk.io/math"
)

// SharePrecision defines the precision factor for share calculations
var SharePrecision = math.NewUint(1e18)

// CalculateCollateralSharesForDeposit calculates how much shares correspond to a collateral amount
// Uses floor division (truncate) for deposits - system favorable
func CalculateCollateralSharesForDeposit(
	totalCollateral math.Uint,
	totalShares math.Uint,
	amount math.Uint,
) math.Uint {
	// If the amount is zero, return zero shares
	if amount.IsZero() {
		return math.ZeroUint()
	}

	// If there are no shares yet, or no collateral, initialize 1:1 with precision
	if totalShares.IsZero() || totalCollateral.IsZero() {
		return amount.Mul(SharePrecision)
	}

	// Calculate based on the current exchange rate
	// shares = (amount * totalShares) / totalCollateral
	return amount.Mul(totalShares).Quo(totalCollateral)
}

// CalculateCollateralSharesForWithdrawal calculates shares for withdrawal - uses ceiling division
// This ensures the system doesn't give out more collateral than it should
func CalculateCollateralSharesForWithdrawal(
	totalCollateral math.Uint,
	totalShares math.Uint,
	amount math.Uint,
) math.Uint {
	// If the amount is zero, return zero shares
	if amount.IsZero() {
		return math.ZeroUint()
	}

	// If there are no shares yet, or no collateral, initialize 1:1 with precision
	if totalShares.IsZero() || totalCollateral.IsZero() {
		return amount.Mul(SharePrecision)
	}

	// Calculate using ceiling division to prevent rounding errors during withdrawals
	product := amount.Mul(totalShares)
	if product.Mod(totalCollateral).IsZero() {
		return product.Quo(totalCollateral)
	}
	return product.Quo(totalCollateral).Add(math.NewUint(1))
}
