package contract

import (
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"strconv"
	"strings"

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
	flagContract        = "contract"
	flagFunc            = "func"
	flagSimulate        = "simulate"
	flagRWSetsHash      = "rwsh"
	flagArgs            = "args"
	flagArgTypes        = "argtypes"
	flagReturnValueType = "type"
	flagSilent          = "silent"
)

// type of return value
const (
	Int     = "int"
	Int32   = "int32"
	Int64   = "int64"
	UInt    = "uint"
	UInt32  = "uint32"
	UInt64  = "uint64"
	Bytes   = "bytes"
	Str     = "str"
	Address = "address"
)

func init() {
	contractCmd.AddCommand(callCmd)
	callCmd.Flags().String(helper.FlagAddress, "", "address")
	callCmd.Flags().String(flagContract, "", "contract address")
	callCmd.Flags().String(flagFunc, "", "function name")
	callCmd.Flags().StringSlice(flagArgs, nil, "arguments")
	callCmd.Flags().StringSlice(flagArgTypes, nil, "types of arguments")
	callCmd.Flags().String(flagRWSetsHash, "", "RWSets hash")
	callCmd.Flags().Uint(flagGas, 0, "gas for tx")
	callCmd.Flags().String(flagReturnValueType, Int, "a type of return value")
	callCmd.Flags().Bool(flagSimulate, false, "execute as simulation")
	callCmd.Flags().Bool(flagSilent, false, "if true, suppress unnecessary output")
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
		args, err := serializeArgs(
			viper.GetStringSlice(flagArgs),
			viper.GetStringSlice(flagArgTypes),
		)
		if err != nil {
			return err
		}
		tx := &transaction.ContractCallTx{
			Address:    caddr,
			Func:       viper.GetString(flagFunc),
			Args:       args,
			RWSetsHash: rwh,
			Common: transaction.CommonTx{
				Code:  transaction.CONTRACT_CALL,
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
			if viper.GetBool(flagSilent) {
				fmt.Print(string(res.Returned))
			} else {
				pretty.Println(rs)
				fmt.Printf("RWSetsHash: 0x%x\n", rs.Hash())
				v, err := parseReturnValue(res.Returned, viper.GetString(flagReturnValueType))
				if err != nil {
					return err
				}
				fmt.Println("Result:", v)
			}
			return nil
		}

		if err := ctx.SignAndBroadcastTx(tx, from); err != nil {
			return err
		}

		return nil
	},
}

func serializeArgs(args []string, types []string) ([][]byte, error) {
	if len(args) != len(types) {
		return nil, fmt.Errorf("the number of arguments does not match the number of types: %v != %v", len(args), len(types))
	}

	var bs [][]byte
	for i, arg := range args {
		switch types[i] {
		case Int32, UInt32:
			v, err := strconv.Atoi(arg)
			if err != nil {
				return nil, err
			}
			var b [4]byte
			binary.BigEndian.PutUint32(b[:], uint32(v))
			bs = append(bs, b[:])
		case Int64, UInt64:
			v, err := strconv.Atoi(arg)
			if err != nil {
				return nil, err
			}
			var b [8]byte
			binary.BigEndian.PutUint64(b[:], uint64(v))
			bs = append(bs, b[:])
		case Bytes:
			if strings.HasPrefix(arg, "0x") {
				a, err := hex.DecodeString(arg[2:])
				if err != nil {
					return nil, err
				}
				bs = append(bs, a)
			} else {
				bs = append(bs, []byte(arg))
			}
		case Str:
			bs = append(bs, []byte(arg))
		case Address:
			ha := strings.TrimPrefix(arg, "0x")
			if l := len(ha); l != 20*2 {
				return nil, fmt.Errorf("address: invalid length %v", l)
			}
			a, err := hex.DecodeString(ha)
			if err != nil {
				return nil, err
			}
			bs = append(bs, a)
		default:
			return nil, fmt.Errorf("unknown types: %v", types[i])
		}
	}
	return bs, nil
}

func parseReturnValue(b []byte, tp string) (interface{}, error) {
	switch tp {
	case Int:
		if l := len(b); l == 4 {
			return int32(binary.BigEndian.Uint32(b)), nil
		} else if l == 8 {
			return int64(binary.BigEndian.Uint64(b)), nil
		} else {
			return nil, fmt.Errorf("unexpected bytes: %x", b)
		}
	case UInt:
		if l := len(b); l == 4 {
			return binary.BigEndian.Uint32(b), nil
		} else if l == 8 {
			return binary.BigEndian.Uint64(b), nil
		} else {
			return nil, fmt.Errorf("unexpected bytes: %x", b)
		}
	case Bytes:
		return b, nil
	case Str:
		return string(b), nil
	case Address:
		if l := len(b); l == 20 {
			var addr common.Address
			copy(addr[:], b)
			return addr, nil
		} else {
			return nil, fmt.Errorf("unexpected bytes: %x", b)
		}
	default:
		return nil, fmt.Errorf("unknown type: %v", tp)
	}
}
