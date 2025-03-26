package types

import sdkmath "cosmossdk.io/math"

// VotingPowerReduction is the default amount of collateral required for 1 unit of consensus-engine power.
// 1e9 collateral (in gwei unit) == 1 MITO == 1 unit of consensus voting power
var VotingPowerReduction = sdkmath.NewInt(1e9)
