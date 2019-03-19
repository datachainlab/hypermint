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
	Address    common.Address
	Func       string // function name
	Args       []string
	RWSetsHash []byte
	CommonTx
}

func (tx *ContractCallTx) ValidateBasic() types.Error {
	if err := tx.CommonTx.ValidateBasic(); err != nil {
		return err
	}
	if tx.Func == ContractInitFunc {
		return ErrInvalidCall(DefaultCodespace, fmt.Sprintf("func '%v' is reserved by contract initializer", ContractInitFunc))
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
