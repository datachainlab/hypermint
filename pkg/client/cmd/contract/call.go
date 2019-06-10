package contract

import (
	"encoding/hex"
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/kr/pretty"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/tendermint/go-amino"

	"github.com/bluele/hypermint/pkg/client"
	"github.com/bluele/hypermint/pkg/client/helper"
	"github.com/bluele/hypermint/pkg/db"
	"github.com/bluele/hypermint/pkg/handler"
	"github.com/bluele/hypermint/pkg/transaction"
	"github.com/bluele/hypermint/pkg/util"
)

const (
	flagContract   = "contract"
	flagFunc       = "func"
	flagSimulate   = "simulate"
	flagRWSetsHash = "rwsh"
	flagArgs       = "args"
)

func init() {
	contractCmd.AddCommand(callCmd)
	callCmd.Flags().String(helper.FlagAddress, "", "address")
	callCmd.Flags().String(flagContract, "", "contract address")
	callCmd.Flags().String(flagFunc, "", "function name")
	callCmd.Flags().StringSlice(flagArgs, nil, "arguments")
	callCmd.Flags().String(flagRWSetsHash, "", "RWSets hash")
	callCmd.Flags().Uint(flagGas, 0, "gas for tx")
	callCmd.Flags().Bool(flagSimulate, false, "execute as simulation")
	util.CheckRequiredFlag(callCmd, helper.FlagAddress, flagGas)
}

var callCmd = &cobra.Command{
	Use:   "call",
	Short: "call contract",
	RunE: func(cmd *cobra.Command, _ []string) error {
		viper.BindPFlags(cmd.Flags())
		ctx, err := client.NewClientContextFromViper()
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

		caddr := common.HexToAddress(viper.GetString(flagContract))

		var rwh []byte
		if hs := viper.GetString(flagRWSetsHash); hs != "" {
			rwh, err = hex.DecodeString(hs)
			if err != nil {
				return err
			}
		}
		var args [][]byte
		for _, arg := range viper.GetStringSlice(flagArgs) {
			args = append(args, []byte(arg))
		}
		tx := &transaction.ContractCallTx{
			Address:    caddr,
			Func:       viper.GetString(flagFunc),
			Args:       args,
			RWSetsHash: rwh,
			Common: transaction.CommonTx{
				From:  from,
				Gas:   uint64(viper.GetInt(flagGas)),
				Nonce: nonce,
			},
		}
		if viper.GetBool(flagSimulate) {
			r, err := ctx.SignAndSimulateTx(tx, from)
			if err != nil {
				return err
			}
			res := new(handler.ContractCallTxResponse)
			if err := amino.UnmarshalBinaryBare(r, res); err != nil {
				return err
			}
			rs := new(db.RWSets)
			if err := rs.FromBytes(res.RWSetsBytes); err != nil {
				return err
			}
			pretty.Println(rs)
			fmt.Printf("RWSetsHash: 0x%x\n", rs.Hash())
			fmt.Println("Result:", string(res.Returned))
			return nil
		}

		if err := ctx.SignAndBroadcastTx(tx, from); err != nil {
			return err
		}

		return nil
	},
}
