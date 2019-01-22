package cmd

import (
	"io/ioutil"
	"os"

	"github.com/bluele/hypermint/pkg/abci/store"
	"github.com/bluele/hypermint/pkg/abci/types"
	sdk "github.com/bluele/hypermint/pkg/abci/types"
	"github.com/bluele/hypermint/pkg/app"
	"github.com/bluele/hypermint/pkg/client/helper"
	"github.com/bluele/hypermint/pkg/contract"
	"github.com/bluele/hypermint/pkg/util"
	"github.com/ethereum/go-ethereum/common"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	abci "github.com/tendermint/tendermint/abci/types"
	dbm "github.com/tendermint/tendermint/libs/db"
)

const (
	flagWASMPath = "path"
	flagArgs     = "args"
	flagEntry    = "entry"
)

func vmCmd(ctx *app.Context) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "vm",
		Short: "exec wasm on vm",
		RunE: func(cmd *cobra.Command, args []string) error {
			viper.BindPFlags(cmd.Flags())
			addr := viper.GetString(helper.FlagAddress)
			from := common.HexToAddress(addr)

			path := viper.GetString(flagWASMPath)
			f, err := os.Open(path)
			if err != nil {
				return err
			}
			defer f.Close()
			b, err := ioutil.ReadAll(f)
			if err != nil {
				return err
			}
			db, err := dbm.NewGoLevelDB("hm", "/tmp")
			if err != nil {
				return err
			}
			defer db.Close()
			cms := store.NewCommitMultiStore(db)
			var key = sdk.NewKVStoreKey("main")
			cms.MountStoreWithDB(key, sdk.StoreTypeIAVL, nil)
			if err := cms.LoadLatestVersion(); err != nil {
				return err
			}
			kvs := cms.GetKVStore(key)
			env := &contract.Env{
				Contract: &contract.Contract{
					Owner: from,
					Code:  b,
				},
				DB:   kvs,
				Args: viper.GetStringSlice(flagArgs),
			}
			c := types.NewContext(cms, abci.Header{}, false, nil)
			if err := env.Exec(c, viper.GetString(flagEntry)); err != nil {
				return err
			}
			cms.Commit()
			return nil
		},
	}
	cmd.Flags().String(flagWASMPath, "", "wasm path")
	cmd.Flags().StringSlice(flagArgs, nil, "arguments")
	cmd.Flags().String(flagEntry, "app_main", "")
	cmd.Flags().String(helper.FlagAddress, "", "address")
	util.CheckRequiredFlag(cmd, flagWASMPath)
	return cmd
}
