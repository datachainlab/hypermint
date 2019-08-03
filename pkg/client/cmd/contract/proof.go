package contract

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/bluele/hypermint/pkg/app"
	"github.com/bluele/hypermint/pkg/client"
	"github.com/bluele/hypermint/pkg/db"
	"github.com/bluele/hypermint/pkg/proof"
	"github.com/bluele/hypermint/pkg/util"

	"github.com/ethereum/go-ethereum/common"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	cmn "github.com/tendermint/tendermint/libs/common"
	rpcclient "github.com/tendermint/tendermint/rpc/client"
)

func init() {
	contractCmd.AddCommand(ProofCMD())
}

func ProofCMD() *cobra.Command {
	const (
		flagContractAddress = "address"
		flagKey             = "key"
		flagValue           = "value"
		flagHeight          = "height"
		flagOutputPath      = "out"
		flagInputPath       = "in"
	)

	var proofCmd = &cobra.Command{
		Use:   "proof",
		Short: "proof utility",
	}

	var getCmd = &cobra.Command{
		Use:   "get",
		Short: "Get a proof of data existence",
		RunE: func(cmd *cobra.Command, args []string) error {
			var (
				key    cmn.HexBytes
				value  cmn.HexBytes
				height int64
				err    error
			)
			viper.BindPFlags(cmd.Flags())
			ctx, err := client.NewClientContextFromViper()
			if err != nil {
				return err
			}
			if h := int64(viper.GetInt(flagHeight)); h > 0 {
				height = h
			} else if h == 0 {
				height = 0
			} else {
				return fmt.Errorf("invalid height %v", h)
			}

			if v := viper.GetString(flagKey); strings.HasPrefix(v, "0x") {
				key, err = hex.DecodeString(v[2:])
				if err != nil {
					return err
				}
			} else {
				key = cmn.HexBytes(v)
			}
			if v := viper.GetString(flagValue); strings.HasPrefix(v, "0x") {
				value, err = hex.DecodeString(v[2:])
				if err != nil {
					return err
				}
			} else {
				value = cmn.HexBytes(v)
			}

			contractAddr := common.HexToAddress(viper.GetString(flagContractAddress))
			path := fmt.Sprintf("/store/%v/key", app.ContractStoreKey.Name())
			res, err := ctx.Client.ABCIQueryWithOptions(
				path,
				append(contractAddr.Bytes(), key.Bytes()...),
				rpcclient.ABCIQueryOptions{
					Height: height,
					Prove:  true,
				},
			)
			if err != nil {
				return err
			}
			vo, err := db.BytesToValueObject(res.Response.Value)
			if err != nil {
				return err
			}
			if value != nil && !bytes.Equal(value, vo.Value) {
				return fmt.Errorf("value is mismatch: %v(%v) != %v(%v)",
					string(value), value.Bytes(),
					string(vo.Value), vo.Value,
				)
			}

			h := res.Response.Height + 1
			c, err := ctx.Client.Commit(&h)
			if err != nil {
				return err
			}
			header := c.SignedHeader.Header
			op, err := proof.MakeKVProofOp(header)
			if err != nil {
				return err
			}
			p := res.Response.Proof
			p.Ops = append(p.Ops, op)

			kvp := proof.MakeKVProofInfo(
				header.Height,
				p,
				contractAddr,
				key,
				vo,
			)
			if err := kvp.VerifyWithHeader(header); err != nil {
				return err
			}
			b, err := kvp.Marshal()
			if err != nil {
				return err
			}
			return ioutil.WriteFile(viper.GetString(flagOutputPath), b, 0644)
		},
	}
	getCmd.Flags().String(flagContractAddress, "", "contract address")
	getCmd.Flags().String(flagKey, "", "key string(if this value is hex, decoded as byte array)")
	getCmd.Flags().String(flagValue, "", "expected value")
	getCmd.Flags().Int64(flagHeight, 0, "height")
	getCmd.Flags().String(flagOutputPath, "", "output path to proof info")
	util.CheckRequiredFlag(getCmd, flagContractAddress, flagKey, flagValue, flagOutputPath)

	var verifyCmd = &cobra.Command{
		Use:   "verify",
		Short: "verify data existence from proof file",
		RunE: func(cmd *cobra.Command, args []string) error {
			viper.BindPFlags(cmd.Flags())
			ctx, err := client.NewClientContextFromViper()
			if err != nil {
				return err
			}

			in := viper.GetString(flagInputPath)
			b, err := ioutil.ReadFile(in)
			if err != nil {
				return err
			}
			kvp := new(proof.KVProofInfo)
			if err := kvp.Unmarshal(b); err != nil {
				return err
			}
			c, err := ctx.Client.Commit(&kvp.Height)
			if err != nil {
				return err
			}
			if err := kvp.VerifyWithHeader(c.SignedHeader.Header); err != nil {
				return err
			}
			fmt.Println("ok")
			return nil
		},
	}
	verifyCmd.Flags().String(flagInputPath, "", "path to proof file")
	util.CheckRequiredFlag(verifyCmd, flagInputPath)

	var showCmd = &cobra.Command{
		Use:   "show",
		Short: "pretty print a proof info",
		RunE: func(cmd *cobra.Command, args []string) error {
			viper.BindPFlags(cmd.Flags())
			in := viper.GetString(flagInputPath)
			b, err := ioutil.ReadFile(in)
			if err != nil {
				return err
			}
			kvp := new(proof.KVProofInfo)
			if err := kvp.Unmarshal(b); err != nil {
				return err
			}
			fmt.Println(kvp.String())
			return nil
		},
	}
	showCmd.Flags().String(flagInputPath, "", "path to proof file")
	util.CheckRequiredFlag(showCmd, flagInputPath)

	proofCmd.AddCommand(getCmd, verifyCmd, showCmd)
	return proofCmd
}
