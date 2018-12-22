package contract

import (
	"io/ioutil"
	"log"
	"os"

	"github.com/bluele/hypermint/pkg/client"
	"github.com/bluele/hypermint/pkg/client/helper"
	"github.com/bluele/hypermint/pkg/transaction"
	"github.com/bluele/hypermint/pkg/util"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	flagCode = "path"
)

func init() {
	contractCmd.AddCommand(deployCmd)
	deployCmd.Flags().String(helper.FlagAddress, "", "address")
	deployCmd.Flags().String(flagCode, "", "contract code path")
	deployCmd.Flags().Uint(flagGas, 0, "gas for tx")
	util.CheckRequiredFlag(deployCmd, helper.FlagAddress, flagCode, flagGas)
}

var deployCmd = &cobra.Command{
	Use:   "deploy",
	Short: "deploy contract code",
	RunE: func(cmd *cobra.Command, args []string) error {
		viper.BindPFlags(cmd.Flags())
		ctx, err := client.NewClientContextFromViper()
		if err != nil {
			return err
		}

		path := viper.GetString(flagCode)
		code, err := getCode(path)
		if err != nil {
			return err
		}

		addrs, err := ctx.GetInputAddresses()
		if err != nil {
			return err
		}
		from := addrs[0]

		nonce, err := transaction.GetNonceByAddress(from)
		if err != nil {
			return err
		}
		tx := &transaction.ContractDeployTx{
			Code: code,
			CommonTx: transaction.CommonTx{
				From:  from,
				Gas:   uint64(viper.GetInt(flagGas)),
				Nonce: nonce,
			},
		}
		if err := ctx.SignAndBroadcastTx(tx, from); err != nil {
			return err
		}
		log.Printf("address=%v", tx.Address().Hex())
		return nil
	},
}

func getCode(path string) ([]byte, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return ioutil.ReadAll(f)
}
