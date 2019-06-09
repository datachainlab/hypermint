package transaction

import (
	"github.com/bluele/hypermint/pkg/abci/types"
	"github.com/bluele/hypermint/pkg/util"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/rlp"
)

var emptyAddr common.Address

type TransferTx struct {
	Common CommonTx
	To     common.Address
	Amount uint64
}

func DecodeTransferTx(b []byte) (*TransferTx, error) {
	tx := new(TransferTx)
	return tx, rlp.DecodeBytes(b, tx)
}

func (tx *TransferTx) SetSignature(sig []byte) {
	tx.Common.SetSignature(sig)
}

func (tx *TransferTx) ValidateBasic() types.Error {
	if err := tx.Common.ValidateBasic(); err != nil {
		return err
	}
	if tx.Amount == 0 {
		return ErrInvalidTransfer(DefaultCodespace, "tx.Amount == 0")
	}
	if isEmptyAddr(tx.To) {
		return ErrInvalidTransfer(DefaultCodespace, "tx.To == empty")
	}
	return tx.Common.VerifySignature(tx.GetSignBytes())
}

func (tx *TransferTx) GetSignBytes() []byte {
	ntx := *tx
	ntx.SetSignature(nil)
	return util.TxHash(ntx.Bytes())
}

func (tx *TransferTx) Bytes() []byte {
	b, err := rlp.EncodeToBytes(tx)
	if err != nil {
		panic(err)
	}
	return b
}

func isEmptyAddr(addr common.Address) bool {
	return addr == emptyAddr
}
