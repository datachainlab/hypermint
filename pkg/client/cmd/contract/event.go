package contract

import (
	"context"
	"fmt"

	"github.com/bluele/hypermint/pkg/client"
	"github.com/bluele/hypermint/pkg/contract"
	"github.com/bluele/hypermint/pkg/util"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/tendermint/tendermint/libs/common"
	"github.com/tendermint/tendermint/types"
)

func init() {
	contractCmd.AddCommand(EventCMD())
}

func EventCMD() *cobra.Command {
	var eventCmd = &cobra.Command{
		Use:   "event",
		Short: "This provides you to pub/sub events",
	}

	const (
		flagContractAddress = "address"
		flagEventName       = "event"
	)

	var subscribeCmd = &cobra.Command{
		Use:   "subscribe",
		Short: "Subscribe Txs using events",
		RunE: func(cmd *cobra.Command, args []string) error {
			viper.BindPFlags(cmd.Flags())
			ctx, err := client.NewClientContextFromViper()
			if err != nil {
				return err
			}
			cl, err := ctx.GetNode()
			if err != nil {
				return err
			}
			if err := cl.OnStart(); err != nil {
				return err
			}
			defer cl.Stop()
			id := common.RandStr(8)
			q := fmt.Sprintf("contract.address='%v' AND event.name='%v'", viper.GetString(flagContractAddress), viper.GetString(flagEventName))
			fmt.Printf("subscription-id=%#v query=%#v\n", id, q)
			out, err := cl.Subscribe(context.Background(), id, q)
			if err != nil {
				return err
			}
			for ev := range out {
				etx := ev.Data.(types.EventDataTx)
				fmt.Printf("TxID=0x%x\n", etx.Tx.Hash())
				for _, ev := range etx.Result.Events {
					for _, tag := range ev.Attributes {
						k := string(tag.GetKey())
						if k == "event.data" {
							ev, err := contract.ParseEventData(tag.GetValue())
							if err != nil {
								return err
							}
							fmt.Println(ev.String())
						} else if k == "event.name" || k == "contract.address" {
							// skip
						} else {
							fmt.Printf("unknown event: %v\n", tag)
						}
					}
				}
			}
			return nil
		},
	}
	subscribeCmd.Flags().String(flagContractAddress, "", "contract address for subscription")
	subscribeCmd.Flags().String(flagEventName, "", "event name for subscription")
	util.CheckRequiredFlag(subscribeCmd, flagContractAddress, flagEventName)

	eventCmd.AddCommand(subscribeCmd)
	return eventCmd
}
