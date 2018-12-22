package cmd

import (
	"os"

	"github.com/bluele/hypermint/pkg/client/cmd/contract"
	"github.com/bluele/hypermint/pkg/client/helper"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var homeDir = os.ExpandEnv("$HOME/.hmcli/keys")

var rootCmd = &cobra.Command{
	Use:   "hmcli",
	Short: "Blockchain Client",
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	// initConfig to be ran when Execute is called
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringP(helper.FlagHomeDir, "", homeDir, "directory for keystore")
	rootCmd.PersistentFlags().BoolP(helper.FlagVerbose, "v", false, "enable verbose output")
	rootCmd.PersistentFlags().String(helper.FlagNode, "tcp://localhost:26657", "<host>:<port> to tendermint rpc interface for this chain")
	rootCmd.PersistentFlags().StringP(helper.FlagPassword, "p", "", "password for signing tx")
	contract.Setup(rootCmd)
	viper.BindPFlags(rootCmd.Flags())
	viper.Set(helper.FlagHomeDir, homeDir)
}

// initConfig reads in config file and ENV variables if set
func initConfig() {
	viper.AutomaticEnv()
}
