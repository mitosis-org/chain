package types

import (
	"testing"

	"cosmossdk.io/math"
	"github.com/stretchr/testify/require"
)

func TestCalculateCollateralSharesForDeposit(t *testing.T) {
	tests := []struct {
		name            string
		totalCollateral math.Uint
		totalShares     math.Uint
		amount          math.Uint
		expected        math.Uint
	}{
		{
			name:            "zero amount",
			totalCollateral: math.NewUint(1000),
			totalShares:     math.NewUint(1000),
			amount:          math.ZeroUint(),
			expected:        math.ZeroUint(),
		},
		{
			name:            "zero collateral and shares",
			totalCollateral: math.ZeroUint(),
			totalShares:     math.ZeroUint(),
			amount:          math.NewUint(500),
			expected:        math.NewUint(500).Mul(SharePrecision),
		},
		{
			// not possible case in practice - existing shares are respected.
			name:            "zero collateral with shares",
			totalCollateral: math.ZeroUint(),
			totalShares:     math.NewUint(1000),
			amount:          math.NewUint(500),
			expected:        math.NewUint(500).Mul(SharePrecision),
		},
		{
			// not possible case in practice - new shares have all existing collateral
			name:            "zero shares with collateral",
			totalCollateral: math.NewUint(1000),
			totalShares:     math.ZeroUint(),
			amount:          math.NewUint(500),
			expected:        math.NewUint(500).Mul(SharePrecision),
		},
		{
			name:            "1:1 ratio",
			totalCollateral: math.NewUint(1000),
			totalShares:     math.NewUint(1000),
			amount:          math.NewUint(500),
			expected:        math.NewUint(500),
		},
		{
			name:            "2:1 shares to collateral ratio",
			totalCollateral: math.NewUint(1000),
			totalShares:     math.NewUint(2000),
			amount:          math.NewUint(500),
			expected:        math.NewUint(1000),
		},
		{
			name:            "1:2 shares to collateral ratio",
			totalCollateral: math.NewUint(2000),
			totalShares:     math.NewUint(1000),
			amount:          math.NewUint(500),
			expected:        math.NewUint(250),
		},
		{
			name:            "large numbers",
			totalCollateral: math.NewUintFromString("1000000000000000000"),
			totalShares:     math.NewUintFromString("3000000000000000000"),
			amount:          math.NewUintFromString("500000000000000000"),
			expected:        math.NewUintFromString("1500000000000000000"),
		},
		{
			name:            "non-divisible result",
			totalCollateral: math.NewUint(3),
			totalShares:     math.NewUint(10),
			amount:          math.NewUint(1),
			expected:        math.NewUint(3), // 1 * 10 / 3 = 3 (floor division for deposit)
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := CalculateCollateralSharesForDeposit(tc.totalCollateral, tc.totalShares, tc.amount)
			require.Equal(t, tc.expected, result)
		})
	}
}

func TestCalculateCollateralSharesForWithdrawal(t *testing.T) {
	tests := []struct {
		name            string
		totalCollateral math.Uint
		totalShares     math.Uint
		amount          math.Uint
		expected        math.Uint
	}{
		{
			name:            "zero amount",
			totalCollateral: math.NewUint(1000),
			totalShares:     math.NewUint(1000),
			amount:          math.ZeroUint(),
			expected:        math.ZeroUint(),
		},
		{
			name:            "zero collateral and shares",
			totalCollateral: math.ZeroUint(),
			totalShares:     math.ZeroUint(),
			amount:          math.NewUint(500),
			expected:        math.NewUint(500).Mul(SharePrecision),
		},
		{
			// not possible case in practice - existing shares are respected.
			name:            "zero collateral with shares",
			totalCollateral: math.ZeroUint(),
			totalShares:     math.NewUint(1000),
			amount:          math.NewUint(500),
			expected:        math.NewUint(500).Mul(SharePrecision),
		},
		{
			// not possible case in practice - new shares have all existing collateral
			name:            "zero shares with collateral",
			totalCollateral: math.NewUint(1000),
			totalShares:     math.ZeroUint(),
			amount:          math.NewUint(500),
			expected:        math.NewUint(500).Mul(SharePrecision),
		},
		{
			name:            "1:1 ratio",
			totalCollateral: math.NewUint(1000),
			totalShares:     math.NewUint(1000),
			amount:          math.NewUint(500),
			expected:        math.NewUint(500),
		},
		{
			name:            "2:1 shares to collateral ratio",
			totalCollateral: math.NewUint(1000),
			totalShares:     math.NewUint(2000),
			amount:          math.NewUint(500),
			expected:        math.NewUint(1000),
		},
		{
			name:            "1:2 shares to collateral ratio",
			totalCollateral: math.NewUint(2000),
			totalShares:     math.NewUint(1000),
			amount:          math.NewUint(500),
			expected:        math.NewUint(250),
		},
		{
			name:            "exact division",
			totalCollateral: math.NewUint(10),
			totalShares:     math.NewUint(20),
			amount:          math.NewUint(5),
			expected:        math.NewUint(10), // 5 * 20 / 10 = 10
		},
		{
			name:            "non-divisible result",
			totalCollateral: math.NewUint(3),
			totalShares:     math.NewUint(10),
			amount:          math.NewUint(1),
			expected:        math.NewUint(4), // 1 * 10 / 3 = 3 + 1 = 4 (ceiling division for withdrawal)
		},
		{
			name:            "another non-divisible result",
			totalCollateral: math.NewUint(10),
			totalShares:     math.NewUint(7),
			amount:          math.NewUint(3),
			expected:        math.NewUint(3), // 3 * 7 / 10 = 2.1 -> 3 (ceiling division for withdrawal)
		},
		{
			name:            "large numbers",
			totalCollateral: math.NewUintFromString("1000000000000000000"),
			totalShares:     math.NewUintFromString("3000000000000000000"),
			amount:          math.NewUintFromString("500000000000000000"),
			expected:        math.NewUintFromString("1500000000000000000"),
		},
		{
			name:            "large numbers with remainder",
			totalCollateral: math.NewUintFromString("1000000000000000003"),
			totalShares:     math.NewUintFromString("3000000000000000000"),
			amount:          math.NewUintFromString("500000000000000001"),
			expected:        math.NewUintFromString("1499999999999999999"),
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := CalculateCollateralSharesForWithdrawal(tc.totalCollateral, tc.totalShares, tc.amount)
			require.Equal(t, tc.expected, result)
		})
	}
}
