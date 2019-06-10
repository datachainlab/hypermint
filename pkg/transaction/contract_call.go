package transaction

import (
	"fmt"

	"github.com/bluele/hypermint/pkg/abci/types"
	"github.com/bluele/hypermint/pkg/util"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/rlp"
)

const (
	ContractInitFunc = "init"
)

type ContractCallTx struct {
	Common     CommonTx
	Address    common.Address
	Func       string // function name
	Args       [][]byte
	RWSetsHash []byte
}

func DecodeContractCallTx(b []byte) (*ContractCallTx, error) {
	tx := new(ContractCallTx)
	return tx, rlp.DecodeBytes(b, tx)
}

func (tx *ContractCallTx) SetSignature(sig []byte) {
	tx.Common.SetSignature(sig)
}

func (tx *ContractCallTx) ValidateBasic() types.Error {
	if err := tx.Common.ValidateBasic(); err != nil {
		return err
	}
	if tx.Func == ContractInitFunc {
		return ErrInvalidCall(DefaultCodespace, fmt.Sprintf("func '%v' is reserved by contract initializer", ContractInitFunc))
	}
	return tx.Common.VerifySignature(tx.GetSignBytes())
}

func (tx *ContractCallTx) Decode(b []byte) error {
	return rlp.DecodeBytes(b, tx)
}

func (tx *ContractCallTx) GetSignBytes() []byte {
	ntx := *tx
	ntx.SetSignature(nil)
	return util.TxHash(ntx.Bytes())
}

func (tx *ContractCallTx) Bytes() []byte {
	b, err := rlp.EncodeToBytes(tx)
	if err != nil {
		panic(err)
	}
	return b
}
