package config

import (
	"errors"
	"os"
	"path/filepath"
	"time"

	"github.com/spf13/viper"
	cfg "github.com/tendermint/tendermint/config"

	"github.com/bluele/hypermint/pkg/util"
)

var ErrConfigNotFound = errors.New("config not found")

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

func SaveConfig(c *cfg.Config) {
	configFilePath := filepath.Join(c.RootDir, "config/config.toml")
	cfg.EnsureRoot(c.RootDir)
	cfg.WriteConfigFile(configFilePath, c)
}

func CreateConfig(moniker, root string) (*cfg.Config, error) {
	c := cfg.DefaultConfig()
	c.SetRoot(root)
	c.Moniker = moniker
	c.ProfListenAddress = "localhost:6060"
	c.P2P.RecvRate = 5120000
	c.P2P.SendRate = 5120000
	c.Consensus.TimeoutCommit = 5000 * time.Millisecond
	c.TxIndex.IndexTags = "contract.address,event.data,event.name"
	return c, unmarshalWithViper(viper.GetViper(), c)
}

func GetConfig(root string) (*cfg.Config, error) {
	configFilePath := filepath.Join(root, "config/config.toml")
	if _, err := os.Stat(configFilePath); err != nil && !os.IsExist(err) {
		return nil, ErrConfigNotFound
	}
	c := new(cfg.Config)
	return c, unmarshalWithViper(viper.GetViper(), c)
}

func unmarshalWithViper(vp *viper.Viper, c *cfg.Config) error {
	// you can configure tedermint params via environment variables.
	// TM_PARAMS="consensus.timeout_commit=3000,instrumentation.prometheus=true" ./liamd start
	util.SetEnvToViper(vp, "TM_PARAMS")
	if err := vp.Unmarshal(c); err != nil {
		return err
	}
	return nil
}
