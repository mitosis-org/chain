package cmd

import (
	tmcfg "github.com/cometbft/cometbft/config"
	serverconfig "github.com/cosmos/cosmos-sdk/server/config"
)

func initAppConfig() (string, interface{}) {
	type CustomAppConfig struct {
		serverconfig.Config
	}

	srvCfg := serverconfig.DefaultConfig()
	srvCfg.StateSync.SnapshotInterval = 1000
	srvCfg.StateSync.SnapshotKeepRecent = 10

	customAppConfig := CustomAppConfig{
		Config: *srvCfg,
	}

	defaultAppTemplate := serverconfig.DefaultConfigTemplate

	return defaultAppTemplate, customAppConfig
}

func initTendermintConfig() *tmcfg.Config {
	cfg := tmcfg.DefaultConfig()

	return cfg
}
