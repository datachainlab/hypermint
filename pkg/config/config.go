package config

import (
	"os"
	"path/filepath"
	"time"

	"github.com/spf13/viper"
	cfg "github.com/tendermint/tendermint/config"
	tmcli "github.com/tendermint/tendermint/libs/cli"

	"github.com/bluele/hypermint/pkg/util"
)

//_____________________________________________________________________

// Configuration structure for command functions that share configuration.
// For example: init, init gen-tx and testnet commands need similar input and run the same code

// Storage for init gen-tx command input parameters
type GenTx struct {
	Name      string
	CliRoot   string
	Overwrite bool
	IP        string
}

func GetConfig() (*cfg.Config, error) {
	c := cfg.DefaultConfig()
	c.SetRoot(viper.GetString(tmcli.HomeFlag))
	c.ProfListenAddress = "localhost:6060"
	c.P2P.RecvRate = 5120000
	c.P2P.SendRate = 5120000
	c.Consensus.TimeoutCommit = 5000 * time.Millisecond

	vp := viper.GetViper()
	// you can configure tedermint params via environment variables.
	// TM_PARAMS="consensus.timeout_commit=3000,instrumentation.prometheus=true" ./liamd start
	util.SetEnvToViper(vp, "TM_PARAMS")
	if err := vp.Unmarshal(c); err != nil {
		return nil, err
	}

	configFilePath := filepath.Join(c.RootDir, "config/config.toml")
	if _, err := os.Stat(configFilePath); os.IsNotExist(err) {
		cfg.EnsureRoot(c.RootDir)
		cfg.WriteConfigFile(configFilePath, c)
		// Fall through, just so that its parsed into memory.
	}

	return c, nil
}

func MustGetConfig() *cfg.Config {
	c, err := GetConfig()
	if err != nil {
		panic(err)
	}
	return c
}
