package transaction

import (
	"github.com/bluele/hypermint/pkg/abci/types"
	"github.com/bluele/hypermint/pkg/util"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/rlp"
)

var emptyAddr common.Address

type TransferTx struct {
	To     common.Address
	Amount uint64
	CommonTx
}

func (tx *TransferTx) Decode(b []byte) error {
	return rlp.DecodeBytes(b, tx)
}

func (tx *TransferTx) ValidateBasic() types.Error {
	if err := tx.CommonTx.ValidateBasic(); err != nil {
		return err
	}
	if tx.Amount == 0 {
		return ErrInvalidTransfer(DefaultCodespace, "tx.Amount == 0")
	}
	if isEmptyAddr(tx.To) {
		return ErrInvalidTransfer(DefaultCodespace, "tx.To == empty")
	}
	return tx.VerifySignature(tx.GetSignBytes())
}

func (tx *TransferTx) GetSignBytes() []byte {
	ntx := *tx
	ntx.Signature = nil
	b, err := rlp.EncodeToBytes(ntx)
	if err != nil {
		panic(err)
	}
	return util.TxHash(b)
}

func (tx *TransferTx) Bytes() []byte {
	b, err := rlp.EncodeToBytes(tx)
	if err != nil {
		panic(err)
	}
	return append([]byte{TRANSFER}, b...)
}

func isEmptyAddr(addr common.Address) bool {
	return addr == emptyAddr
}
