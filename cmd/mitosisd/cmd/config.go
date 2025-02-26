package cmd

import (
	tmcfg "github.com/cometbft/cometbft/config"
	serverconfig "github.com/cosmos/cosmos-sdk/server/config"
)

type AppConfig struct {
	serverconfig.Config `mapstructure:",squash"`
	Engine              *EngineConfig `mapstructure:"engine"`
}

type EngineConfig struct {
	Mock            bool   `mapstructure:"mock"`
	Endpoint        string `mapstructure:"endpoint"`
	JWTFile         string `mapstructure:"jwt-file"`
	BuildDelay      string `mapstructure:"build-delay"`
	BuildOptimistic bool   `mapstructure:"build-optimistic"`
}

func DefaultAppConfig() AppConfig {
	srvCfg := serverconfig.DefaultConfig()
	srvCfg.StateSync.SnapshotInterval = 1000
	srvCfg.StateSync.SnapshotKeepRecent = 10

	return AppConfig{
		Config: *srvCfg,
		Engine: &EngineConfig{
			Mock:            false,
			Endpoint:        "http://127.0.0.1:8551",
			JWTFile:         "",
			BuildDelay:      "600ms", // it should be slightly longer than geth's --miner.recommit=500ms.
			BuildOptimistic: true,
		},
	}
}

func initAppConfig() (string, AppConfig) {
	appConfig := DefaultAppConfig()

	defaultAppTemplate := serverconfig.DefaultConfigTemplate + `
###############################################################################
###                          Engine                                         ###
###############################################################################

[engine]

# If it is true, the execution client will be mocked and endpoint and jwt-file will be ignored.
# Otherwise, it will connect to a real execution client.
mock = {{ .Engine.Mock }}

# Execution client Engine API http endpoint.
endpoint = "{{ .Engine.Endpoint }}"

# Execution client JWT file used for authentication.
jwt-file = "{{ .Engine.JWTFile }}"

# Build delay is the time to wait before building a block.
# Slightly longer value is recommended than --miner.recommit in case of geth.
# e.g., 600ms if --miner.recommit=500ms.
build-delay = "{{ .Engine.BuildDelay }}"

# If it is true, build a block optimistically.
build-optimistic = {{ .Engine.BuildOptimistic }}
`

	return defaultAppTemplate, appConfig
}

func initTendermintConfig() *tmcfg.Config {
	cfg := tmcfg.DefaultConfig()

	return cfg
}
