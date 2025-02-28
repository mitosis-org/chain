package app

import (
	runtimev1alpha1 "cosmossdk.io/api/cosmos/app/runtime/v1alpha1"
	appv1alpha1 "cosmossdk.io/api/cosmos/app/v1alpha1"
	authmodulev1 "cosmossdk.io/api/cosmos/auth/module/v1"
	bankmodulev1 "cosmossdk.io/api/cosmos/bank/module/v1"
	consensusmodulev1 "cosmossdk.io/api/cosmos/consensus/module/v1"
	evidencemodulev1 "cosmossdk.io/api/cosmos/evidence/module/v1"
	genutilmodulev1 "cosmossdk.io/api/cosmos/genutil/module/v1"
	slashingmodulev1 "cosmossdk.io/api/cosmos/slashing/module/v1"
	txconfigv1 "cosmossdk.io/api/cosmos/tx/config/v1"
	upgrademodulev1 "cosmossdk.io/api/cosmos/upgrade/module/v1"
	"cosmossdk.io/core/appconfig"
	"cosmossdk.io/depinject"
	evidencetypes "cosmossdk.io/x/evidence/types"
	upgradetypes "cosmossdk.io/x/upgrade/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/cosmos/cosmos-sdk/x/genutil"
	evmgovmodulev1 "github.com/mitosis-org/chain/api/mitosis/evmgov/module/v1"
	evmvalmodulev1 "github.com/mitosis-org/chain/api/mitosis/evmvalidator/module/v1"
	evmvaltypes "github.com/mitosis-org/chain/x/evmvalidator/types"
	evmengmodule "github.com/omni-network/omni/octane/evmengine/module"

	"github.com/cosmos/cosmos-sdk/runtime"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	consensustypes "github.com/cosmos/cosmos-sdk/x/consensus/types"
	genutiltypes "github.com/cosmos/cosmos-sdk/x/genutil/types"
	slashingtypes "github.com/cosmos/cosmos-sdk/x/slashing/types"
	evmgovtypes "github.com/mitosis-org/chain/x/evmgov/types"
	evmengtypes "github.com/omni-network/omni/octane/evmengine/types"
)

// AppConfig returns the default app depinject config.
func AppConfig() depinject.Config {
	return depinject.Configs(
		appConfig,
		depinject.Supply(
			// supply custom module basics
			map[string]module.AppModuleBasic{
				genutiltypes.ModuleName: genutil.NewAppModuleBasic(genutiltypes.DefaultMessageValidator),
			},
		),
	)
}

//nolint:gochecknoglobals // Cosmos-style
var (
	genesisModuleOrder = []string{
		authtypes.ModuleName,
		banktypes.ModuleName,
		evmvaltypes.ModuleName, // should be before slashing & genutil & evidence
		slashingtypes.ModuleName,
		genutiltypes.ModuleName,
		upgradetypes.ModuleName,
		evidencetypes.ModuleName,
		evmengtypes.ModuleName,
		evmgovtypes.ModuleName,
	}

	preBlockers = []string{
		upgradetypes.ModuleName, // NOTE: upgrade module must come first, as upgrades might break state schema.
	}

	beginBlockers = []string{
		slashingtypes.ModuleName,
		evidencetypes.ModuleName,
	}

	endBlockers = []string{
		evmvaltypes.ModuleName,
	}

	blockAccAddrs = []string{
		authtypes.FeeCollectorName,
	}

	moduleAccPerms = []*authmodulev1.ModuleAccountPermission{
		{Account: authtypes.FeeCollectorName},
		{Account: evmgovtypes.ModuleName, Permissions: []string{authtypes.Burner}},
	}

	// appConfig application configuration (used by depinject).
	appConfig = appconfig.Compose(&appv1alpha1.Config{
		Modules: []*appv1alpha1.ModuleConfig{
			{
				Name: runtime.ModuleName,
				Config: appconfig.WrapAny(&runtimev1alpha1.Module{
					AppName:       "mitosisd",
					PreBlockers:   preBlockers,
					BeginBlockers: beginBlockers,
					EndBlockers:   endBlockers,
					InitGenesis:   genesisModuleOrder,
					OverrideStoreKeys: []*runtimev1alpha1.StoreKeyConfig{
						{
							ModuleName: authtypes.ModuleName,
							KvStoreKey: "acc",
						},
					},
				}),
			},
			{
				Name: authtypes.ModuleName,
				Config: appconfig.WrapAny(&authmodulev1.Module{
					ModuleAccountPermissions: moduleAccPerms,
					Bech32Prefix:             Bech32Prefix,
					Authority:                evmgovtypes.ModuleName,
				}),
			},
			{
				Name: "tx",
				Config: appconfig.WrapAny(&txconfigv1.Config{
					SkipAnteHandler: true, // Disable ante handler since there is only EVM payload tx.
					SkipPostHandler: true,
				}),
			},
			{
				Name: banktypes.ModuleName,
				Config: appconfig.WrapAny(&bankmodulev1.Module{
					BlockedModuleAccountsOverride: blockAccAddrs,
					Authority:                     evmgovtypes.ModuleName,
				}),
			},
			{
				Name: consensustypes.ModuleName,
				Config: appconfig.WrapAny(&consensusmodulev1.Module{
					Authority: evmgovtypes.ModuleName,
				}),
			},
			{
				Name: slashingtypes.ModuleName,
				Config: appconfig.WrapAny(&slashingmodulev1.Module{
					Authority: evmgovtypes.ModuleName,
				}),
			},
			{
				Name:   genutiltypes.ModuleName,
				Config: appconfig.WrapAny(&genutilmodulev1.Module{}),
			},
			{
				Name:   evidencetypes.ModuleName,
				Config: appconfig.WrapAny(&evidencemodulev1.Module{}),
			},
			{
				Name: upgradetypes.ModuleName,
				Config: appconfig.WrapAny(&upgrademodulev1.Module{
					Authority: evmgovtypes.ModuleName,
				}),
			},
			{
				Name: evmengtypes.ModuleName,
				Config: appconfig.WrapAny(&evmengmodule.Module{
					Authority: evmgovtypes.ModuleName,
				}),
			},
			{
				Name: evmvaltypes.ModuleName,
				Config: appconfig.WrapAny(&evmvalmodulev1.Module{
					Authority:                  evmgovtypes.ModuleName,
					EvmValidatorEntrypointAddr: EVMValidatorEntrypointAddr,
				}),
			},
			{
				Name: evmgovtypes.ModuleName,
				Config: appconfig.WrapAny(&evmgovmodulev1.Module{
					Authority:                   evmgovtypes.ModuleName,
					EvmGovernanceEntrypointAddr: EVMGovernanceEntrypointAddr,
				}),
			},
		},
	})
)
