package bind

import (
	"context"
	"github.com/bluele/hypermint/pkg/contract/event"
	"github.com/bluele/hypermint/pkg/transaction"
	"github.com/ethereum/go-ethereum/common"
)

type AsyncResult struct {
	TxHash common.Hash
}

type SyncResult struct {
	TxHash common.Hash
}

type CommitResult struct {
	TxHash common.Hash
	Height int64
	Entries []*event.Entry
}

type ContractTransactor interface {
	BroadcastTxAsync(ctx context.Context, tx transaction.Transaction) (*AsyncResult, error)
	BroadcastTxSync(ctx context.Context, tx transaction.Transaction) (*SyncResult, error)
	BroadcastTxCommit(ctx context.Context, tx transaction.Transaction) (*CommitResult, error)
}

type SimulateResult struct {
	Data []byte
}

type ContractSimulator interface {
	SimulateTx(ctx context.Context, tx transaction.Transaction) (*SimulateResult, error)
}

type ContractBackend interface {
	ContractSimulator
	ContractTransactor
}
