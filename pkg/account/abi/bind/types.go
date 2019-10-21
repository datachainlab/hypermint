package bind

import (
	"context"
	"github.com/bluele/hypermint/pkg/transaction"
	"github.com/ethereum/go-ethereum/common"
)

type SignerFn func(tx transaction.Transaction, addr common.Address) (transaction.Transaction, error)

type TransactOpts struct {
	From    common.Address
	Signer  SignerFn
	Context context.Context
}
