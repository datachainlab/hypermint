package hmclient

import (
	"context"
	"github.com/bluele/hypermint/pkg/abci/codec"
	hmabcitypes "github.com/bluele/hypermint/pkg/abci/types"
	"github.com/bluele/hypermint/abi/bind"
	"github.com/bluele/hypermint/pkg/contract/event"
	"github.com/bluele/hypermint/pkg/handler"
	"github.com/bluele/hypermint/pkg/transaction"
	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
	"github.com/tendermint/go-amino"
	tmabcitypes "github.com/tendermint/tendermint/abci/types"
	rpcClient "github.com/tendermint/tendermint/rpc/client"
	rpcTypes "github.com/tendermint/tendermint/rpc/core/types"
	"strings"
)

type Signer func(msg []byte, addr []byte) ([]byte, error)

type Client struct {
	Client rpcClient.Client
}

func NewClient(conn string) *Client {
	c := rpcClient.NewHTTP(conn, "/websocket")
	return &Client{
		Client: c,
	}
}

func (c Client) SimulateTx(ctx context.Context, tx transaction.Transaction) (*bind.SimulateResult, error) {
	if data, err := c.simulateTx(tx); err != nil {
		return nil, err
	} else {
		res := new(handler.ContractCallTxResponse)
		if err := amino.UnmarshalBinaryBare(data, res); err != nil {
			return nil, err
		}
		return &bind.SimulateResult{
			Data: res.Returned,
		}, err
	}
}

func (c *Client) simulateTx(tx transaction.Transaction) ([]byte, error) {
	res, err := c.Client.ABCIQuery("/app/simulate", tx.Bytes())
	if err != nil {
		return nil, err
	}
	var result hmabcitypes.Result
	codec.Cdc.MustUnmarshalBinaryLengthPrefixed(res.Response.Value, &result)

	if result.Code != 0 {
		return result.Data, errors.Errorf("Simulate failed: (%d) %s",
			result.Code, result.Log)
	}

	return result.Data, nil
}

func (c *Client) BroadcastTxAsync(ctx context.Context, tx transaction.Transaction) (*bind.AsyncResult, error) {
	res, err := c.Client.BroadcastTxAsync(tx.Bytes())
	if err != nil {
		return nil, err
	}
	return &bind.AsyncResult{
		TxHash: common.BytesToHash(res.Hash.Bytes()),
	}, nil
}

func (c *Client) BroadcastTxSync(ctx context.Context, tx transaction.Transaction) (*bind.SyncResult, error) {
	res, err := c.Client.BroadcastTxSync(tx.Bytes())
	if err != nil {
		return nil, err
	}
	return &bind.SyncResult{
		TxHash: common.BytesToHash(res.Hash.Bytes()),
	}, nil
}

func (c *Client) BroadcastTxCommit(ctx context.Context, tx transaction.Transaction) (*bind.CommitResult, error) {
	res, err := c.Client.BroadcastTxCommit(tx.Bytes())
	if err != nil && res != nil {
		if res.CheckTx.GetCode() == 104 && strings.Contains(err.Error(), "already exists:") {
			return nil, ErrAlreadyExists
		}
	}
	if err != nil {
		return nil, err
	}

	var events []tmabcitypes.Event
	events = res.DeliverTx.Events
	var entries []*event.Entry
	for _, tme := range events {
		es, err := event.GetEntryFromEvent(hmabcitypes.Event(tme))
		if err != nil {
			return nil, err
		}
		entries = append(entries, es...)
	}

	return &bind.CommitResult{
		TxHash:  common.BytesToHash(res.Hash.Bytes()),
		Height:  res.Height,
		Entries: entries,
	}, nil
}

func (c *Client) TransactionEventEntries(ctx context.Context, txHash common.Hash) ([]*event.Entry, error) {
	resultTx, err := c.transactionByHash(txHash.Bytes())
	if err != nil {
		return nil, err
	}
	tx, err := transaction.DecodeTx(resultTx.Tx)
	if err != nil {
		return nil, errors.New(err.Error())
	}
	switch tx := tx.(type) {
	case *transaction.ContractCallTx:
		events, err := event.GetContractEventsFromResultTx(tx.Address, resultTx)
		if err != nil {
			return nil, err
		}
		var entries []*event.Entry
		for _, tme := range events {
			es, err := event.GetEntryFromEvent(tme)
			if err != nil {
				return nil, err
			}
			entries = append(entries, es...)
		}
		return entries, nil
	default:
		return nil, errors.New("invalid transaction type")
	}
}

func (c *Client) transactionByHash(txHash []byte) (*rpcTypes.ResultTx, error) {
	resultTx, err := c.Client.Tx(txHash, false)
	if err != nil {
		return nil, err
	}
	result := resultTx.TxResult
	if result.Code != 0 {
		return resultTx, errors.Errorf("Tx failed: (%d) %s", result.Code, result.Log)
	}
	return resultTx, nil
}
