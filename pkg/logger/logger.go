package logger

import (
	"os"

	"github.com/spf13/viper"
	cfg "github.com/tendermint/tendermint/config"
	"github.com/tendermint/tendermint/libs/cli"
	tmflags "github.com/tendermint/tendermint/libs/cli/flags"
	"github.com/tendermint/tendermint/libs/log"
)

var logger log.Logger

// this method is goroutine unsafe
func SetLogger(lg log.Logger) {
	if logger == nil {
		logger = lg
	}
}

func GetLogger() log.Logger {
	return logger
}

func GetDefaultLogger(lv string) log.Logger {
	logger := log.NewTMLogger(log.NewSyncWriter(os.Stdout))
	logger, err := tmflags.ParseLogLevel(lv, logger, cfg.DefaultLogLevel())
	if err != nil {
		panic(err)
	}
	if viper.GetBool(cli.TraceFlag) {
		logger = log.NewTracingLogger(logger)
	}
	return logger.With("module", "main")
}

func Debug(msg string, keyvals ...interface{}) {
	logger.Debug(msg, keyvals...)
}

func Info(msg string, keyvals ...interface{}) {
	logger.Info(msg, keyvals...)
}

func Error(msg string, keyvals ...interface{}) {
	logger.Error(msg, keyvals...)
}

func With(keyvals ...interface{}) log.Logger {
	return logger.With(keyvals...)
}
