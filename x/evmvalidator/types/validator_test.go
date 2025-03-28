package types

import (
	"testing"

	"cosmossdk.io/math"
	"github.com/stretchr/testify/require"
)

func TestComputeVotingPower(t *testing.T) {
	testCases := []struct {
		name             string
		collateral       math.Uint
		extraVotingPower math.Uint
		maxLeverageRatio math.LegacyDec
		expectedPower    int64
	}{
		{
			name:             "collateral only",
			collateral:       math.NewUint(1e9), // 1 MITO
			extraVotingPower: math.ZeroUint(),
			maxLeverageRatio: math.LegacyNewDec(1),
			expectedPower:    1,
		},
		{
			name:             "collateral and extra power",
			collateral:       math.NewUint(1e9), // 1 MITO
			extraVotingPower: math.NewUint(1e9), // 1 MITO
			maxLeverageRatio: math.LegacyNewDec(2),
			expectedPower:    2,
		},
		{
			name:             "leverage ratio limiting",
			collateral:       math.NewUint(1e9), // 1 MITO
			extraVotingPower: math.NewUint(5e9), // 5 MITO
			maxLeverageRatio: math.LegacyNewDec(3),
			expectedPower:    3,
		},
		{
			name:             "zero collateral",
			collateral:       math.ZeroUint(),
			extraVotingPower: math.NewUint(1e9), // 1 MITO
			maxLeverageRatio: math.LegacyNewDec(1000),
			expectedPower:    0,
		},
		{
			name:             "fractional voting power",
			collateral:       math.NewUint(1e9),                // 1 MITO
			extraVotingPower: math.NewUint(2e9),                // 5 MITO
			maxLeverageRatio: math.LegacyNewDecWithPrec(29, 1), // x2.9
			expectedPower:    2,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			validator := Validator{
				Collateral:       tc.collateral,
				ExtraVotingPower: tc.extraVotingPower,
			}
			power := validator.ComputeVotingPower(tc.maxLeverageRatio)
			require.Equal(t, tc.expectedPower, power)
		})
	}
}
