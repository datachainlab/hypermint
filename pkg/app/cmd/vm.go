package cmd

import (
	"fmt"
	"io/ioutil"

	"github.com/bluele/hypermint/pkg/abci/store"
	sdk "github.com/bluele/hypermint/pkg/abci/types"
	"github.com/bluele/hypermint/pkg/app"
	"github.com/bluele/hypermint/pkg/client/helper"
	"github.com/bluele/hypermint/pkg/contract"
	"github.com/bluele/hypermint/pkg/db"
	"github.com/bluele/hypermint/pkg/util"
	"github.com/ethereum/go-ethereum/common"
	"github.com/kr/pretty"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	abci "github.com/tendermint/tendermint/abci/types"
	dbm "github.com/tendermint/tendermint/libs/db"
)

const (
	flagWASMPath = "path"
	flagArgs     = "args"
	flagEntry    = "entry"
	flagSimulate = "simulate"
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
			b, err := ioutil.ReadFile(path)
			if err != nil {
				return err
			}
			gdb, err := dbm.NewGoLevelDB("hm", "/tmp")
			if err != nil {
				return err
			}
			defer gdb.Close()
			cms := store.NewCommitMultiStore(gdb)
			var key = sdk.NewKVStoreKey("main")
			cms.MountStoreWithDB(key, sdk.StoreTypeIAVL, nil)
			if err := cms.LoadLatestVersion(); err != nil {
				return err
			}
			var kvs sdk.KVStore
			if viper.GetBool(flagSimulate) {
				kvs = cms.CacheMultiStore().GetKVStore(key)
			} else {
				kvs = cms.GetKVStore(key)
			}
			env := &contract.Env{
				Sender: from,
				Contract: &contract.Contract{
					Owner: from,
					Code:  b,
				},
				DB:   db.NewVersionedDB(kvs, db.Version{1, 1}),
				Args: viper.GetStringSlice(flagArgs),
			}
			c := sdk.NewContext(cms, abci.Header{}, false, nil)
			res, err := env.Exec(c, viper.GetString(flagEntry))
			if err != nil {
				return err
			}
			cms.Commit()
			pretty.Println(res.RWSets)
			fmt.Printf("RWSetsHash is '0x%x'\n", res.RWSets.Hash())
			fmt.Println("response:", string(res.Response))
			return nil
		},
	}
	cmd.Flags().String(flagWASMPath, "", "wasm path")
	cmd.Flags().StringSlice(flagArgs, nil, "arguments")
	cmd.Flags().String(flagEntry, "app_main", "")
	cmd.Flags().String(helper.FlagAddress, "", "address")
	cmd.Flags().Bool(flagSimulate, false, "is simluation")
	util.CheckRequiredFlag(cmd, flagWASMPath)
	return cmd
}
