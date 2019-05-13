package context

import (
	"fmt"

	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
	rpclient "github.com/tendermint/tendermint/rpc/client"
	ctypes "github.com/tendermint/tendermint/rpc/core/types"

	"github.com/bluele/hypermint/pkg/abci/codec"
	"github.com/bluele/hypermint/pkg/abci/types"
	"github.com/bluele/hypermint/pkg/client/helper"
	"github.com/bluele/hypermint/pkg/transaction"
	"github.com/bluele/hypermint/pkg/util"
)

type Context struct {
	HomeDir        string
	NodeURI        string
	InputAddresses []common.Address
	Client         rpclient.Client
	Verbose        bool
}

// Prepares a simple rpc.Client
func (ctx *Context) GetNode() (rpclient.Client, error) {
	if ctx.Client == nil {
		return nil, errors.New("must define node URI")
	}
	return ctx.Client, nil
}

// Get the from address from the name flag
func (ctx *Context) GetInputAddresses() ([]common.Address, error) {
	ks := keystore.NewKeyStore(ctx.HomeDir, keystore.StandardScryptN, keystore.StandardScryptP)
	for _, addr := range ctx.InputAddresses {
		if !ks.HasAddress(addr) {
			return nil, errors.Errorf("no account for: %s", addr.Hex())
		}
	}
	return ctx.InputAddresses, nil
}

func (ctx *Context) GetPassphrase(addr common.Address) (string, error) {
	pass := viper.GetString(helper.FlagPassword)
	if pass == "" {
		return ctx.getPassphraseFromStdin(addr)
	}
	return pass, nil
}

// Get passphrase from std input
func (ctx *Context) getPassphraseFromStdin(addr common.Address) (string, error) {
	buf := helper.BufferStdin()
	prompt := fmt.Sprintf("Password to sign with '%s':", addr.Hex())
	return helper.GetPassword(prompt, buf)
}

// Broadcast the transaction bytes to Tendermint
func (ctx *Context) BroadcastTx(tx []byte) (*ctypes.ResultBroadcastTxCommit, error) {
	node, err := ctx.GetNode()
	if err != nil {
		return nil, err
	}

	res, err := node.BroadcastTxCommit(tx)
	if err != nil {
		return res, err
	}

	if res.CheckTx.Code != uint32(0) {
		return res, errors.Errorf("CheckTx failed: (%d) %s",
			res.CheckTx.Code, res.CheckTx.Log)
	}
	if res.DeliverTx.Code != uint32(0) {
		return res, errors.Errorf("DeliverTx failed: (%d) %s",
			res.DeliverTx.Code, res.DeliverTx.Log)
	}
	return res, err
}

func (ctx *Context) Sign(msg []byte, addr common.Address) ([]byte, error) {
	passphrase, err := ctx.GetPassphrase(addr)
	if err != nil {
		return nil, err
	}
	ks := keystore.NewKeyStore(ctx.HomeDir, keystore.StandardScryptN, keystore.StandardScryptP)
	acc := accounts.Account{
		Address: addr,
	}
	acct, err := ks.Find(acc)
	if err != nil {
		return nil, err
	}
	return ks.SignHashWithPassphrase(acct, passphrase, msg)
}

func (ctx *Context) SignAndBroadcastTx(tx transaction.Transaction, addr common.Address) error {
	sig, err := ctx.Sign(tx.GetSignBytes(), addr)
	if err != nil {
		return err
	}
	tx.SetSignature(sig)

	res, err := ctx.BroadcastTx(tx.Bytes())
	if err != nil {
		return err
	}
	if ctx.Verbose {
		fmt.Printf("txHash=%v BlockHeight=%v\n", res.Hash.String(), res.Height)
	}
	return nil
}

func (ctx *Context) SignAndSimulateTx(tx transaction.Transaction, addr common.Address) ([]byte, error) {
	sig, err := ctx.Sign(tx.GetSignBytes(), addr)
	if err != nil {
		return nil, err
	}
	tx.SetSignature(sig)

	res, err := ctx.Client.ABCIQuery("/app/simulate", tx.Bytes())
	if err != nil {
		return nil, err
	}
	var result types.Result
	codec.Cdc.MustUnmarshalBinaryLengthPrefixed(res.Response.Value, &result)

	if result.Code != 0 {
		return result.Data, errors.Errorf("Simulate failed: (%d) %s",
			result.Code, result.Log)
	}

	return result.Data, nil
}

func (ctx *Context) GetBalanceByAddress(addr common.Address) (uint64, error) {
	cl, err := ctx.GetNode()
	if err != nil {
		return 0, err
	}
	res, err := cl.ABCIQuery("/store/main/key", addr.Bytes())
	if err != nil {
		return 0, err
	}
	if res.Response.IsErr() {
		return 0, errors.New(res.Response.String())
	}
	if res.Response.Value == nil {
		return 0, errors.New("response is nil")
	}
	return util.BytesToUint64(res.Response.Value)
}
