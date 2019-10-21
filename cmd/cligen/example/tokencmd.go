// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package main

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"math/big"
	"os"
	"strings"

	"github.com/bluele/hypermint/pkg/account/abi/bind"
	clibind "github.com/bluele/hypermint/pkg/account/cli/bind"
	"github.com/bluele/hypermint/pkg/account/client"
	"github.com/bluele/hypermint/pkg/transaction"
	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/common"
)

// Reference imports to suppress errors if they are not otherwise used.
var (
	_ = big.NewInt
	_ = strings.NewReader
	_ = bind.Bind
	_ = clibind.Bind
	_ = common.Big1
	_ = transaction.ContractInitFunc
	_ = bytes.NewBuffer
	_ = binary.Read
	_ = json.NewEncoder
)

func TokenCmd(contractAddress, callerAddress string) *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "token",
		Short: "token",
	}
	cmd.PersistentFlags().Bool("verbose", false, "verbose")
	cmd.PersistentFlags().String("passphrase", "xxx", "passphrase")
	cmd.PersistentFlags().String("caller", callerAddress, "caller")
	cmd.PersistentFlags().String("contract", contractAddress, "contract")
	cmd.PersistentFlags().String("endpoint", "tcp://localhost:26657", "endpoint")
	cmd.AddCommand(tokenGetBalanceCmd)
	cmd.AddCommand(tokenTransferCmd)
	cmd.AddCommand(tokenInitCmd)

	return cmd
}

func init() {
}

var tokenGetBalanceCmd = &cobra.Command{
	Use:   "get-balance",
	Short: "get balance",
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := viper.BindPFlags(cmd.Flags()); err != nil {
			return err
		}
		verbose := viper.GetBool("verbose")

		if verbose {
			fmt.Fprintf(os.Stderr, "passphrase=%v\n", viper.GetString("passphrase"))
		}

		ks := clibind.NewKeyStore("keystore")
		opts, c, err := tokenContractFromFlags(ks)
		if err != nil {
			return err
		}

		if verbose {
			fmt.Fprintf(os.Stderr, "simulating GetBalance...")
		}
		v0, err := c.GetBalance(
			opts,
		)
		if err != nil {
			return err
		}
		if verbose {
			fmt.Fprintf(os.Stderr, "done\n")
		}

		if _, err := fmt.Fprintln(os.Stdout, v0); err != nil {
			return err
		}

		return nil

	},
}

func init() {

	cmd := tokenTransferCmd

	cmd.Flags().String("to", "", "to")

	cmd.Flags().String("amount", "", "amount")

}

var tokenTransferCmd = &cobra.Command{
	Use:   "transfer",
	Short: "transfer",
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := viper.BindPFlags(cmd.Flags()); err != nil {
			return err
		}
		verbose := viper.GetBool("verbose")

		if len(viper.GetString("to")) == 0 {
			return errors.New("invalid address")
		}
		argto := common.HexToAddress(viper.GetString("to"))

		if verbose {
			fmt.Fprintf(os.Stderr, "to=%v\n", argto)
		}

		argamount := int64(viper.GetInt64("amount"))

		if verbose {
			fmt.Fprintf(os.Stderr, "amount=%v\n", argamount)
		}

		if verbose {
			fmt.Fprintf(os.Stderr, "passphrase=%v\n", viper.GetString("passphrase"))
		}

		ks := clibind.NewKeyStore("keystore")
		opts, c, err := tokenContractFromFlags(ks)
		if err != nil {
			return err
		}

		if verbose {
			fmt.Fprintf(os.Stderr, "committing Transfer...")
		}
		r, err := c.TransferCommit(
			opts,
			argto,
			argamount,
		)
		if err != nil {
			return err
		}
		if verbose {
			fmt.Fprintf(os.Stderr, "%v\n", r.TxHash.Hex())
		}
		ed := TokenEventDecoder
		if ed == nil {
			return errors.New("TokenEventDecoder is nil")
		}

		return nil

	},
}

func init() {
}

var tokenInitCmd = &cobra.Command{
	Use:   "init",
	Short: "init",
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := viper.BindPFlags(cmd.Flags()); err != nil {
			return err
		}
		verbose := viper.GetBool("verbose")

		if verbose {
			fmt.Fprintf(os.Stderr, "passphrase=%v\n", viper.GetString("passphrase"))
		}

		ks := clibind.NewKeyStore("keystore")
		opts, c, err := tokenContractFromFlags(ks)
		if err != nil {
			return err
		}

		if verbose {
			fmt.Fprintf(os.Stderr, "committing Init...")
		}
		r, err := c.InitCommit(
			opts,
		)
		if err != nil {
			return err
		}
		if verbose {
			fmt.Fprintf(os.Stderr, "%v\n", r.TxHash.Hex())
		}
		ed := TokenEventDecoder
		if ed == nil {
			return errors.New("TokenEventDecoder is nil")
		}

		return nil

	},
}

func tokenContractFromFlags(ks *keystore.KeyStore) (*bind.TransactOpts, TokenContract, error) {
	pp, err := clibind.GetFlagPassphrase(func() (string, error) {
		pp := viper.GetString("passphrase")
		if pp == "xxx" {
			return "", errors.New("passphrase required")
		}
		return pp, nil
	})()
	if err != nil {
		return nil, nil, err
	}

	caller := common.HexToAddress(viper.GetString("caller"))
	if err := clibind.Unlock(ks, caller, pp); err != nil {
		return nil, nil, err
	}

	contractAddress, err := clibind.GetFlagAddress(func() (string, error) {
		return viper.GetString("contract"), nil
	}, nil)()
	if err != nil {
		return nil, nil, err
	}

	opts, err := bind.NewKeyStoreTransactor(ks, accounts.Account{Address: caller})
	if err != nil {
		return nil, nil, err
	}

	endpoint := viper.GetString("endpoint")
	cl := hmclient.NewClient(endpoint)
	cToken, err := NewToken(contractAddress, cl)
	if err != nil {
		return nil, nil, err
	}

	return opts, cToken, nil
}
