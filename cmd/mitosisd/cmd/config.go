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
	ValidatorMode bool   `mapstructure:"validator-mode"`
	Mock          bool   `mapstructure:"mock"`
	Endpoint      string `mapstructure:"endpoint"`
	JWTFile       string `mapstructure:"jwt-file"`
}

func DefaultAppConfig() AppConfig {
	srvCfg := serverconfig.DefaultConfig()
	srvCfg.StateSync.SnapshotInterval = 1000
	srvCfg.StateSync.SnapshotKeepRecent = 10

	return AppConfig{
		Config: *srvCfg,
		Engine: &EngineConfig{
			ValidatorMode: false,
			Mock:          true,
			Endpoint:      "",
			JWTFile:       "",
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

# If you're running a validator node, must set this to true.
validator-mode = {{ .Engine.ValidatorMode }}

# If it is true, the execution client will be mocked and endpoint and jwt-file will be ignored.
# Otherwise, it will be connect to a real execution client.
mock = {{ .Engine.Mock }}

# Execution client Engine API http endpoint.
endpoint = "{{ .Engine.Endpoint }}"

# Execution client JWT file used for authentication.
jwt-file = "{{ .Engine.JWTFile }}"
`

	return defaultAppTemplate, appConfig
}

func initTendermintConfig() *tmcfg.Config {
	cfg := tmcfg.DefaultConfig()

	return cfg
}
