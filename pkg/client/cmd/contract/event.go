package contract

import (
	"context"
	"fmt"

	"github.com/bluele/hypermint/pkg/abci/types"
	"github.com/bluele/hypermint/pkg/client"
	"github.com/bluele/hypermint/pkg/contract/event"
	"github.com/bluele/hypermint/pkg/util"

	ethcmn "github.com/ethereum/go-ethereum/common"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/tendermint/tendermint/libs/common"
	tmtypes "github.com/tendermint/tendermint/types"
)

func init() {
	contractCmd.AddCommand(EventCMD())
}

func EventCMD() *cobra.Command {
	var eventCmd = &cobra.Command{
		Use:   "event",
		Short: "This provides you to pub/sub events",
	}

	// common
	const (
		flagContractAddress = "address"
		flagEventName       = "event.name"
		flagEventValue      = "event.value"
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
			q := fmt.Sprintf("tm.event='Tx' AND contract.address='%v' AND contract.event.name='%v'", viper.GetString(flagContractAddress), viper.GetString(flagEventName))
			fmt.Printf("subscription-id=%#v query=%#v\n", id, q)
			out, err := cl.Subscribe(context.Background(), id, q)
			if err != nil {
				return err
			}
			for ev := range out {
				etx := ev.Data.(tmtypes.EventDataTx)
				fmt.Printf("TxID=0x%x\n", etx.Tx.Hash())
				for _, ev := range etx.Result.Events {
					if ev.Type != "contract" {
						continue
					}
					printEvents([]types.Event{types.Event(ev)})
				}
			}
			return nil
		},
	}
	subscribeCmd.Flags().String(flagContractAddress, "", "contract address for subscription")
	subscribeCmd.Flags().String(flagEventName, "", "event name for subscription")
	util.CheckRequiredFlag(subscribeCmd, flagContractAddress, flagEventName)
	eventCmd.AddCommand(subscribeCmd)

	// search
	const (
		flagCount = "count"
	)

	var searchCmd = &cobra.Command{
		Use:   "search",
		Short: "Search Txs using events",
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
			contractAddr := ethcmn.HexToAddress(viper.GetString(flagContractAddress))
			q, err := event.MakeEventSearchQuery(
				contractAddr,
				viper.GetString(flagEventName),
				viper.GetString(flagEventValue),
			)
			if err != nil {
				return err
			}
			res, err := cl.TxSearch(q, true, 0, 0)
			if err != nil {
				return err
			}
			if viper.GetBool(flagCount) {
				var count int
				for _, tx := range res.Txs {
					events, err := event.GetContractEventsFromResultTx(contractAddr, tx)
					if err != nil {
						return err
					}
					if len(events) == 0 {
						continue
					}
					events, err = event.FilterContractEvents(
						events,
						viper.GetString(flagEventName),
						viper.GetString(flagEventValue),
					)
					if err != nil {
						return err
					}
					if len(events) == 0 {
						continue
					}
					count++
				}
				fmt.Print(count)
				return nil
			} else {
				for _, tx := range res.Txs {
					fmt.Printf("Tx=0x%x\n", tx.Tx.Hash())
					events, err := event.GetContractEventsFromResultTx(contractAddr, tx)
					if err != nil {
						return err
					}
					printEvents(events)
				}
			}
			return nil
		},
	}

	searchCmd.Flags().String(flagContractAddress, "", "contract address for subscription")
	searchCmd.Flags().String(flagEventName, "", "event name")
	searchCmd.Flags().String(flagEventValue, "", "event value as hex string")
	searchCmd.Flags().Bool(flagCount, false, "if true, only print count of matched txs")
	util.CheckRequiredFlag(searchCmd, flagContractAddress, flagEventName)
	eventCmd.AddCommand(searchCmd)

	return eventCmd
}

func printEvents(events []types.Event) {
	for _, ev := range events {
		fmt.Printf("event type=%v\n", ev.Type)
		es, err := event.GetEntryFromEvent(ev)
		if err != nil {
			panic(err)
		}
		for _, entry := range es {
			fmt.Printf("\tname=%v value=%v\n", string(entry.Name), string(entry.Value))
		}
	}
}
