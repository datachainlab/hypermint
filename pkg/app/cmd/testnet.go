package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	gcrypto "github.com/ethereum/go-ethereum/crypto"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	amino "github.com/tendermint/go-amino"
	cfg "github.com/tendermint/tendermint/config"
	cmn "github.com/tendermint/tendermint/libs/common"
	tmtime "github.com/tendermint/tendermint/types/time"

	"github.com/bluele/hypermint/pkg/app"
	"github.com/bluele/hypermint/pkg/config"
	"github.com/bluele/hypermint/pkg/util"
	"github.com/bluele/hypermint/pkg/validator"
)

var (
	nodeDirPrefix  = "node-dir-prefix"
	nValidators    = "validators-num"
	nNonValidators = "non-validators-num"
	outputDir      = "output-dir"

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
	cmd.Flags().IntP(nNonValidators, "n", 0,
		"Number of non-validators to initialize the testnet with")
	cmd.Flags().String(startingIPAddress, "192.168.0.1",
		"Starting IP address (192.168.0.1 results in persistent peers list ID0@192.168.0.1:46656, ID1@192.168.0.2:46656, ...)")
	cmd.Flags().AddFlagSet(appInit.FlagsAppGenTx)
	return cmd
}

func testnetWithConfig(c *cfg.Config, cdc *amino.Codec, appInit app.AppInit) error {
	outDir := viper.GetString(outputDir)
	numValidators := viper.GetInt(nValidators)
	numNonValidators := viper.GetInt(nNonValidators)

	// Generate private key, node ID, initial transaction
	for i := 0; i < numValidators; i++ {
		nodeDirName := fmt.Sprintf("%s%d", viper.GetString(nodeDirPrefix), i)
		di := getDirsInfo(outDir, i)
		c.SetRoot(di.NodeDir())

		err := os.MkdirAll(di.ConfigDir(), nodeDirPerm)
		if err != nil {
			_ = os.RemoveAll(outDir)
			return err
		}

		err = os.MkdirAll(di.ClientDir(), nodeDirPerm)
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
			di.ClientDir(),
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
		err = writeFile(name, di.ClientDir(), cliPrint)
		if err != nil {
			return err
		}

		// Gather gentxs folder
		name = fmt.Sprintf("%v.json", nodeDirName)
		err = writeFile(name, di.GenTxsDir(), genTxFile)
		if err != nil {
			return err
		}
	}

	// Generate genesis.json and config.toml
	chainID := "chain-" + cmn.RandStr(6)
	genTime := tmtime.Now()
	var genesisFilePath string
	for i := 0; i < numValidators; i++ {
		di := getDirsInfo(outDir, i)
		initConfig := InitConfig{
			chainID,
			true,
			di.GenTxsDir(),
			true,
			genTime,
		}
		c.Moniker = di.DirName()
		c.SetRoot(di.NodeDir())

		// Run `init` and generate genesis.json and config.toml
		_, _, _, err := initWithConfig(cdc, appInit, c, initConfig)
		if err != nil {
			return err
		}
		if i == 0 {
			genesisFilePath = c.GenesisFile()
		}
	}

	for i := 0; i < numNonValidators; i++ {
		id := i + numValidators
		di := getDirsInfo(outDir, id)
		c.Moniker = di.DirName()
		c.SetRoot(di.NodeDir())

		err := os.MkdirAll(di.ConfigDir(), nodeDirPerm)
		if err != nil {
			_ = os.RemoveAll(outDir)
			return err
		}
		err = os.MkdirAll(di.ClientDir(), nodeDirPerm)
		if err != nil {
			_ = os.RemoveAll(outDir)
			return err
		}

		prv, err := gcrypto.GenerateKey()
		if err != nil {
			return err
		}
		validator.GenFilePVWithECDSA(c.PrivValidatorKeyFile(), c.PrivValidatorStateFile(), prv)
		if err := util.CopyFile(genesisFilePath, filepath.Join(c.RootDir, "config/genesis.json")); err != nil {
			return err
		}
	}

	fmt.Printf("Successfully initialized node directories val=%v nval=%v\n", viper.GetInt(nValidators), viper.GetInt(nNonValidators))
	return nil
}

type dirsInfo struct {
	rootDir string
	dirName string
}

func (di dirsInfo) DirName() string {
	return di.dirName
}

func (di dirsInfo) NodeRootDir() string {
	return filepath.Join(di.rootDir, di.dirName)
}

func (di dirsInfo) ClientDir() string {
	return filepath.Join(di.NodeRootDir(), "hmcli")
}

func (di dirsInfo) NodeDir() string {
	return filepath.Join(di.NodeRootDir(), "hmd")
}

func (di dirsInfo) ConfigDir() string {
	return filepath.Join(di.NodeDir(), "config")
}

func (di dirsInfo) GenTxsDir() string {
	return filepath.Join(di.rootDir, "gentxs")
}

func getDirsInfo(rootDir string, id int) dirsInfo {
	dirName := fmt.Sprintf("%s%d", viper.GetString(nodeDirPrefix), id)
	return dirsInfo{
		rootDir: rootDir,
		dirName: dirName,
	}
}
