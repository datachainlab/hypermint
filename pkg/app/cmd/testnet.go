package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	amino "github.com/tendermint/go-amino"
	cfg "github.com/tendermint/tendermint/config"
	cmn "github.com/tendermint/tendermint/libs/common"
	tmtime "github.com/tendermint/tendermint/types/time"

	"github.com/bluele/hypermint/pkg/app"
	"github.com/bluele/hypermint/pkg/config"
)

var (
	nodeDirPrefix = "node-dir-prefix"
	nValidators   = "validators-num"
	outputDir     = "output-dir"

	startingIPAddress = "starting-ip-address"
)

const nodeDirPerm = 0755

// get cmd to initialize all files for tendermint testnet and application
func testnetFilesCmd(ctx *app.Context, cdc *amino.Codec, appInit app.AppInit) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "testnet",
		Short: "Initialize files for a hmd testnet",
		Long: `testnet will create "v" number of directories and populate each with
necessary files (private validator, genesis, config, etc.).

Note, strict routability for addresses is turned off in the config file.

Example:

	hmd testnet --validators-num 4 --output-dir ./output --starting-ip-address 192.168.10.2
	`,
		RunE: func(_ *cobra.Command, _ []string) error {
			config := ctx.Config
			err := testnetWithConfig(config, cdc, appInit)
			return err
		},
	}
	cmd.Flags().IntP(nValidators, "v", 4,
		"Number of validators to initialize the testnet with")
	cmd.Flags().StringP(outputDir, "o", "./mytestnet",
		"Directory to store initialization data for the testnet")
	cmd.Flags().String(nodeDirPrefix, "node",
		"Prefix the directory name for each node with (node results in node0, node1, ...)")
	cmd.Flags().String(startingIPAddress, "192.168.0.1",
		"Starting IP address (192.168.0.1 results in persistent peers list ID0@192.168.0.1:46656, ID1@192.168.0.2:46656, ...)")
	cmd.Flags().AddFlagSet(appInit.FlagsAppGenTx)
	return cmd
}

func testnetWithConfig(c *cfg.Config, cdc *amino.Codec, appInit app.AppInit) error {
	outDir := viper.GetString(outputDir)
	numValidators := viper.GetInt(nValidators)

	// Generate private key, node ID, initial transaction
	for i := 0; i < numValidators; i++ {
		nodeDirName := fmt.Sprintf("%s%d", viper.GetString(nodeDirPrefix), i)
		nodeDir := filepath.Join(outDir, nodeDirName, "hmd")
		clientDir := filepath.Join(outDir, nodeDirName, "hmcli")
		gentxsDir := filepath.Join(outDir, "gentxs")
		c.SetRoot(nodeDir)

		err := os.MkdirAll(filepath.Join(nodeDir, "config"), nodeDirPerm)
		if err != nil {
			_ = os.RemoveAll(outDir)
			return err
		}

		err = os.MkdirAll(clientDir, nodeDirPerm)
		if err != nil {
			_ = os.RemoveAll(outDir)
			return err
		}

		c.Moniker = nodeDirName
		ip, err := getIP(i)
		if err != nil {
			return err
		}

		genTxConfig := config.GenTx{
			nodeDirName,
			clientDir,
			true,
			ip,
		}

		// Run `init gen-tx` and generate initial transactions
		cliPrint, genTxFile, err := gentxWithConfig(cdc, appInit, c, genTxConfig)
		if err != nil {
			return err
		}

		// Save private key seed words
		name := fmt.Sprintf("%v.json", "key_seed")
		err = writeFile(name, clientDir, cliPrint)
		if err != nil {
			return err
		}

		// Gather gentxs folder
		name = fmt.Sprintf("%v.json", nodeDirName)
		err = writeFile(name, gentxsDir, genTxFile)
		if err != nil {
			return err
		}
	}

	// Generate genesis.json and config.toml
	chainID := "chain-" + cmn.RandStr(6)
	genTime := tmtime.Now()
	for i := 0; i < numValidators; i++ {
		nodeDirName := fmt.Sprintf("%s%d", viper.GetString(nodeDirPrefix), i)
		nodeDir := filepath.Join(outDir, nodeDirName, "hmd")
		gentxsDir := filepath.Join(outDir, "gentxs")
		initConfig := InitConfig{
			chainID,
			true,
			gentxsDir,
			true,
			genTime,
		}
		c.Moniker = nodeDirName
		c.SetRoot(nodeDir)

		// Run `init` and generate genesis.json and config.toml
		_, _, _, err := initWithConfig(cdc, appInit, c, initConfig)
		if err != nil {
			return err
		}
	}

	fmt.Printf("Successfully initialized %v node directories\n", viper.GetInt(nValidators))
	return nil
}
