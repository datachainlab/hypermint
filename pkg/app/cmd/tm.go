package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	tcmd "github.com/tendermint/tendermint/cmd/tendermint/commands"
	"github.com/tendermint/tendermint/crypto"
	"github.com/tendermint/tendermint/libs/bech32"
	"github.com/tendermint/tendermint/p2p"
	pvm "github.com/tendermint/tendermint/privval"

	"github.com/bluele/hypermint/pkg/app"
)

const (
	FlagJson = "json"
)

// showNodeIDCmd - ported from Tendermint, dump node ID to stdout
func showNodeIDCmd(ctx *app.Context) *cobra.Command {
	return &cobra.Command{
		Use:   "show-node-id",
		Short: "Show this node's ID",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg := ctx.Config
			nodeKey, err := p2p.LoadOrGenNodeKey(cfg.NodeKeyFile())
			if err != nil {
				return err
			}
			fmt.Println(nodeKey.ID())
			return nil
		},
	}
}

// showValidator - ported from Tendermint, show this node's validator info
func showValidatorCmd(ctx *app.Context) *cobra.Command {
	cmd := cobra.Command{
		Use:   "show-validator",
		Short: "Show this node's tendermint validator info",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg := ctx.Config
			privValidator := pvm.LoadOrGenFilePV(cfg.PrivValidatorKeyFile(), cfg.PrivValidatorStateFile())
			valPubKey := privValidator.Key.PubKey

			if viper.GetBool(FlagJson) {
				return printlnJSON(valPubKey)
			}

			pubkey, err := Bech32ifyConsPub(valPubKey)
			if err != nil {
				return err
			}

			fmt.Println(pubkey)
			return nil
		},
	}
	cmd.Flags().Bool(FlagJson, false, "get machine parseable output")
	return &cmd
}

// showAddressCmd - show this node's validator address
func showAddressCmd(ctx *app.Context) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "show-address",
		Short: "Shows this node's tendermint validator address",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg := ctx.Config
			privValidator := pvm.LoadOrGenFilePV(cfg.PrivValidatorKeyFile(), cfg.PrivValidatorStateFile())
			valAddr := privValidator.Key.Address

			if viper.GetBool(FlagJson) {
				return printlnJSON(valAddr)
			}

			fmt.Printf("0x%X", valAddr)
			return nil
		},
	}

	cmd.Flags().Bool(FlagJson, false, "get machine parseable output")
	return cmd
}

func printlnJSON(v interface{}) error {
	cdc := app.GetCodec()
	marshalled, err := cdc.MarshalJSON(v)
	if err != nil {
		return err
	}
	fmt.Println(string(marshalled))
	return nil
}

// UnsafeResetAllCmd - extension of the tendermint command, resets initialization
func UnsafeResetAllCmd(ctx *app.Context) *cobra.Command {
	return &cobra.Command{
		Use:   "unsafe-reset-all",
		Short: "Resets the blockchain database, removes address book files, and resets priv_validator.json to the genesis state",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg := ctx.Config
			tcmd.ResetAll(cfg.DBDir(), cfg.P2P.AddrBookFile(), cfg.PrivValidatorKeyFile(), cfg.PrivValidatorStateFile(), ctx.Logger)
			return nil
		},
	}
}

const Bech32PrefixConsPub = "cosmosvalconspub"

// Bech32ifyConsPub returns a Bech32 encoded string containing the
// Bech32PrefixConsPub prefixfor a given consensus node's PubKey.
func Bech32ifyConsPub(pub crypto.PubKey) (string, error) {
	return bech32.ConvertAndEncode(Bech32PrefixConsPub, pub.Bytes())
}
