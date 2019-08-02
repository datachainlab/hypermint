package main

import (
	"encoding/json"
	"io"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/libs/cli"
	"github.com/tendermint/tendermint/libs/log"
	"github.com/tendermint/tendermint/types"
	db "github.com/tendermint/tm-db"

	"github.com/bluele/hypermint/pkg/app"
	"github.com/bluele/hypermint/pkg/app/cmd"
	"github.com/bluele/hypermint/pkg/logger"
)

const (
	appName = "hm"
	confDir = "$HOME/.hmd"
)

func main() {
	cobra.EnableCommandSorting = false
	ctx := new(app.Context)
	rootCmd := &cobra.Command{
		Use:   "hmd",
		Short: "hm node",
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			return app.SetupContext(ctx)
		},
	}
	rootCmd.PersistentFlags().String("log_level", "debug", "Log level")

	cmd.AddServerCommands(
		ctx,
		app.GetCodec(),
		rootCmd,
		app.NewAppInit(),
		app.ConstructAppCreator(
			newApp,
			appName,
		),
		app.ConstructAppExporter(
			exportAppState,
			appName,
		),
	)

	viper.BindPFlags(rootCmd.Flags())
	// prepare and add flags
	rootDir := os.ExpandEnv(confDir)
	executor := cli.PrepareBaseCmd(rootCmd, "PC", rootDir)
	executor.Execute()
}

func newApp(lg log.Logger, db db.DB, traceStore io.Writer) abci.Application {
	logger.SetLogger(lg)
	return app.NewChain(lg, db, traceStore)
}

func exportAppState(lg log.Logger, db db.DB, traceStore io.Writer) (json.RawMessage, []types.GenesisValidator, error) {
	logger.SetLogger(lg)
	return app.NewChain(lg, db, traceStore).ExportAppStateJSON()
}
