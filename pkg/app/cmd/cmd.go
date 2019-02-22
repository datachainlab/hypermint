package cmd

import (
	"github.com/spf13/cobra"
	amino "github.com/tendermint/go-amino"

	"github.com/bluele/hypermint/pkg/app"
)

// lineBreak can be included in a command list to provide a blank line
// to help with readability
var lineBreak = &cobra.Command{Run: func(*cobra.Command, []string) {}}

// add server commands
func AddServerCommands(
	ctx *app.Context, cdc *amino.Codec,
	rootCmd *cobra.Command, appInit app.AppInit,
	appCreator app.AppCreator, appExport app.AppExporter) {

	tendermintCmd := &cobra.Command{
		Use:   "tendermint",
		Short: "Tendermint subcommands",
	}

	tendermintCmd.AddCommand(
		showNodeIDCmd(ctx),
		showValidatorCmd(ctx),
		showAddressCmd(ctx),
		validatorCmd(ctx),
	)

	rootCmd.AddCommand(
		initCmd(ctx, cdc, appInit),
		createCmd(ctx),
		testnetFilesCmd(ctx, cdc, appInit),
		startCmd(ctx, appCreator),
		UnsafeResetAllCmd(ctx),
		vmCmd(ctx),
		lineBreak,
		tendermintCmd,
		// TODO impl export cmd
		// server.ExportCmd(ctx, cdc, appExport),
		lineBreak,
		versionCmd,
	)
}
