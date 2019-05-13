package app

import (
	"github.com/spf13/viper"
	cfg "github.com/tendermint/tendermint/config"
	tmcli "github.com/tendermint/tendermint/libs/cli"
	"github.com/tendermint/tendermint/libs/common"
	"github.com/tendermint/tendermint/libs/log"

	"github.com/bluele/hypermint/pkg/config"
	"github.com/bluele/hypermint/pkg/logger"
)

// server context
type Context struct {
	Config *cfg.Config
	Logger log.Logger
}

// SetupContext initializes config object and bind its to context
func SetupContext(ctx *Context) error {
	root := viper.GetString(tmcli.HomeFlag)
	c, err := config.GetConfig(root)
	if err == config.ErrConfigNotFound {
		c, err = config.CreateConfig(common.RandStr(8), root)
		if err != nil {
			return err
		}
		config.SaveConfig(c)
	}
	if err != nil {
		return err
	}
	c.SetRoot(root)
	lg := logger.GetDefaultLogger(c.LogLevel)
	ctx.Config = c
	ctx.Logger = lg
	return nil
}
