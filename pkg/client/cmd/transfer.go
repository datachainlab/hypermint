package cmd

import (
	"errors"

	"github.com/bluele/hypermint/pkg/client"
	"github.com/bluele/hypermint/pkg/client/helper"
	"github.com/bluele/hypermint/pkg/transaction"
	"github.com/bluele/hypermint/pkg/util"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	flagTo     = "to"
	flagAmount = "amount"
	flagGas    = "gas"
)

func init() {
	rootCmd.AddCommand(transferCmd)
	transferCmd.Flags().String(flagTo, "", "Addresse sending to")
	transferCmd.Flags().Uint(flagAmount, 0, "Amount to be spent")
	transferCmd.Flags().Uint(flagGas, 0, "gas for tx")
	transferCmd.Flags().String(helper.FlagAddress, "", "Address to sign with")
	util.CheckRequiredFlag(transferCmd, flagAmount)
	util.CheckRequiredFlag(transferCmd, flagGas)
}

var transferCmd = &cobra.Command{
	Use:   "transfer",
	Short: "Build, Sign, and Send transactions",
	RunE: func(cmd *cobra.Command, args []string) error {
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
		tos, err := helper.ParseAddrs(viper.GetString(flagTo))
		if err != nil {
			return err
		}
		if len(tos) == 0 {
			return errors.New("must provide an address to send to")
		}

		nonce, err := transaction.GetNonceByAddress(from)
		if err != nil {
			return err
		}
		tx := &transaction.TransferTx{
			Common: transaction.CommonTx{
				Code:  transaction.TRANSFER,
				From:  from,
				Gas:   uint64(viper.GetInt(flagGas)),
				Nonce: nonce,
			},
			To:     tos[0],
			Amount: uint64(viper.GetInt(flagAmount)),
		}
		return ctx.SignAndBroadcastTx(tx, from)
	},
}
