package app

import (
	cfg "github.com/tendermint/tendermint/config"
	"github.com/tendermint/tendermint/libs/log"

	"github.com/bluele/hypermint/pkg/config"
	"github.com/bluele/hypermint/pkg/logger"
)

// server context
type Context struct {
	Config *cfg.Config
	Logger log.Logger
}

func SetupContext(ctx *Context) {
	c := config.MustGetConfig()
	lg := logger.GetDefaultLogger(c.LogLevel)
	ctx.Config = c
	ctx.Logger = lg
}

func NewContext(config *cfg.Config, logger log.Logger) *Context {
	return &Context{config, logger}
}
