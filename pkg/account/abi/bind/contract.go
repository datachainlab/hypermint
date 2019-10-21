package bind

import (
	"errors"
	"github.com/ethereum/go-ethereum/common"
)

type BoundContract struct {
	address common.Address
	transactor ContractTransactor
	simulator ContractSimulator
}

func NewBoundContract(address common.Address, simulator ContractSimulator, transactor ContractTransactor) *BoundContract {
	return &BoundContract{
		address: address,
		transactor: transactor,
		simulator: simulator,
	}
}

func (b *BoundContract) Transact(opt *TransactOpts, fn string, args... []byte) (*SyncResult, error) {
	return b.TransactWithRWH(opt, fn, nil, args...)
}

func (b *BoundContract) TransactWithRWH(opt *TransactOpts, fn string, rwh []byte, args... []byte) (*SyncResult, error) {
	tx, err := MakeContractCallTx(opt.From, b.address, fn, args, rwh)
	if err != nil {
		return nil, err
	}
	if opt.Signer == nil {
		return nil, errors.New("no signer")
	}
	signedTx, err := opt.Signer(tx, opt.From)
	if err != nil {
		return nil, err
	}
	return b.transactor.BroadcastTxSync(opt.Context, signedTx)
}

func (b *BoundContract) TransactCommit(opt *TransactOpts, fn string, args... []byte) (*CommitResult, error) {
	return b.TransactCommitWithRWH(opt, fn, nil, args...)
}

func (b *BoundContract) TransactCommitWithRWH(opt *TransactOpts, fn string, rwh []byte, args... []byte) (*CommitResult, error) {
	tx, err := MakeContractCallTx(opt.From, b.address, fn, args, rwh)
	if err != nil {
		return nil, err
	}
	if opt.Signer == nil {
		return nil, errors.New("no signer")
	}
	signedTx, err := opt.Signer(tx, opt.From)
	if err != nil {
		return nil, err
	}
	return b.transactor.BroadcastTxCommit(opt.Context, signedTx)
}

func (b *BoundContract) Simulate(opt *TransactOpts, fn string, args ...[]byte) (*SimulateResult, error) {
	return b.SimulateTxWithRWH(opt, fn, nil, args...)
}

func (b *BoundContract) SimulateTxWithRWH(opt *TransactOpts, fn string, rwh []byte, args ...[]byte) (*SimulateResult, error) {
	tx, err := MakeContractCallTx(opt.From, b.address, fn, args, rwh)
	if err != nil {
		return nil, err
	}
	if opt.Signer == nil {
		return nil, errors.New("no signer")
	}
	signedTx, err := opt.Signer(tx, opt.From)
	if err != nil {
		return nil, err
	}
	return b.simulator.SimulateTx(opt.Context, signedTx)
}
