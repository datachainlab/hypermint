package cmd

import (
	hnode "github.com/bluele/hypermint/pkg/node"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	abcis "github.com/tendermint/tendermint/abci/server"
	tcmd "github.com/tendermint/tendermint/cmd/tendermint/commands"
	cmn "github.com/tendermint/tendermint/libs/common"
	"github.com/tendermint/tendermint/node"
	pvm "github.com/tendermint/tendermint/privval"
	"github.com/tendermint/tendermint/proxy"

	"github.com/bluele/hypermint/pkg/app"
)

const (
	flagWithTendermint = "with-tendermint"
	flagAddress        = "address"
	flagTraceStore     = "trace-store"
	flagPruning        = "pruning"
	flagBlkNumPerTerm  = "term_block_num"
)

// startCmd runs the service passed in, either stand-alone or in-process with
// Tendermint.
func startCmd(ctx *app.Context, appCreator app.AppCreator) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "start",
		Short: "Run the full node",
		RunE: func(cmd *cobra.Command, args []string) error {
			if !viper.GetBool(flagWithTendermint) {
				ctx.Logger.Info("Starting ABCI without Tendermint")
				return startStandAlone(ctx, appCreator)
			}

			ctx.Logger.Info("Starting ABCI with Tendermint")

			_, err := startInProcess(ctx, appCreator)
			return err
		},
	}

	// core flags for the ABCI application
	cmd.Flags().Bool(flagWithTendermint, true, "Run abci app embedded in-process with tendermint")
	cmd.Flags().String(flagAddress, "tcp://0.0.0.0:26658", "Listen address")
	cmd.Flags().String(flagTraceStore, "", "Enable KVStore tracing to an output file")
	cmd.Flags().String(flagPruning, "syncable", "Pruning strategy: syncable, nothing, everything")
	cmd.Flags().Uint(flagBlkNumPerTerm, 100, "The number of block in term")

	// add support for all Tendermint-specific command line options
	tcmd.AddNodeFlags(cmd)
	return cmd
}

func startStandAlone(ctx *app.Context, appCreator app.AppCreator) error {
	addr := viper.GetString(flagAddress)
	home := viper.GetString("home")
	traceStore := viper.GetString(flagTraceStore)

	app, err := appCreator(home, ctx.Logger, traceStore)
	if err != nil {
		return err
	}

	svr, err := abcis.NewServer(addr, "socket", app)
	if err != nil {
		return errors.Errorf("error creating listener: %v\n", err)
	}

	svr.SetLogger(ctx.Logger.With("module", "abci-server"))

	err = svr.Start()
	if err != nil {
		cmn.Exit(err.Error())
	}

	// wait forever
	cmn.TrapSignal(func() {
		// cleanup
		err = svr.Stop()
		if err != nil {
			cmn.Exit(err.Error())
		}
	})
	return nil
}

func startInProcess(ctx *app.Context, appCreator app.AppCreator) (*node.Node, error) {
	cfg := ctx.Config
	home := cfg.RootDir
	traceStore := viper.GetString(flagTraceStore)

	app, err := appCreator(home, ctx.Logger, traceStore)
	if err != nil {
		return nil, err
	}
	nodeKey, err := hnode.LoadNodeKey(cfg.NodeKeyFile())
	if err != nil {
		return nil, err
	}
	pv := pvm.LoadFilePV(cfg.PrivValidatorKeyFile(), cfg.PrivValidatorStateFile())

	// create & start tendermint node
	tmNode, err := node.NewNode(
		cfg,
		pv,
		nodeKey,
		proxy.NewLocalClientCreator(app),
		node.DefaultGenesisDocProviderFunc(cfg),
		node.DefaultDBProvider,
		node.DefaultMetricsProvider(cfg.Instrumentation),
		ctx.Logger.With("module", "node"),
	)
	if err != nil {
		return nil, err
	}

	err = tmNode.Start()
	if err != nil {
		return nil, err
	}

	// Sleep forever and then...
	cmn.TrapSignal(func() {
		tmNode.Stop()
	})

	return tmNode, nil
}
