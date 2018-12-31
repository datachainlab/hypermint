package transaction

import (
	"github.com/bluele/hypermint/pkg/abci/types"
	"github.com/bluele/hypermint/pkg/util"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/rlp"
)

type ContractCallTx struct {
	Address common.Address
	Func    string // function name
	Args    []string
	CommonTx
}

func (tx *ContractCallTx) ValidateBasic() types.Error {
	if err := tx.CommonTx.ValidateBasic(); err != nil {
		return err
	}
	return tx.VerifySignature(tx.GetSignBytes())
}

func (tx *ContractCallTx) Decode(b []byte) error {
	return rlp.DecodeBytes(b, tx)
}

func (tx *ContractCallTx) GetSignBytes() []byte {
	ntx := *tx
	ntx.Signature = nil
	return util.TxHash(ntx.Bytes())
}

func (tx *ContractCallTx) Bytes() []byte {
	b, err := rlp.EncodeToBytes(tx)
	if err != nil {
		panic(err)
	}
	return append([]byte{CONTRACT_CALL}, b...)
}
