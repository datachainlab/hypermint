package cmd

import (
	"io/ioutil"
	"log"
	"os"

	"github.com/bluele/hypermint/pkg/abci/store"
	sdk "github.com/bluele/hypermint/pkg/abci/types"
	"github.com/bluele/hypermint/pkg/app"
	"github.com/bluele/hypermint/pkg/contract"
	"github.com/bluele/hypermint/pkg/util"
	"github.com/ethereum/go-ethereum/common"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	dbm "github.com/tendermint/tendermint/libs/db"
)

const (
	flagWASMPath = "path"
)

func vmCmd(ctx *app.Context) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "vm",
		Short: "exec wasm on vm",
		RunE: func(cmd *cobra.Command, args []string) error {
			viper.BindPFlags(cmd.Flags())
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
			value := kvs.Get([]byte("key"))
			log.Println("value=>", string(value))
			vmn := contract.NewVMManager()
			v, err := vmn.GetVM(kvs, &contract.Contract{
				Owner: common.Address{},
				Code:  b,
			})
			if err != nil {
				return err
			}
			if err := v.ExecContract("app_main"); err != nil {
				return err
			}
			cms.Commit()
			return nil
		},
	}
	cmd.Flags().String(flagWASMPath, "", "wasm path")
	util.CheckRequiredFlag(cmd, flagWASMPath)
	return cmd
}
