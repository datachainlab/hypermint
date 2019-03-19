package transaction

import (
	"github.com/bluele/hypermint/pkg/abci/types"
	"github.com/bluele/hypermint/pkg/util"
	"github.com/ethereum/go-ethereum/rlp"
)

type ContractDeployTx struct {
	Code []byte
	CommonTx
}

func (tx *ContractDeployTx) Decode(b []byte) error {
	return rlp.DecodeBytes(b, tx)
}

func (tx *ContractDeployTx) ValidateBasic() types.Error {
	if err := tx.CommonTx.ValidateBasic(); err != nil {
		return err
	}
	if tx.Code == nil {
		return ErrInvalidDeploy(DefaultCodespace, "tx.Code == nil")
	}
	return tx.VerifySignature(tx.GetSignBytes())
}

func (tx *ContractDeployTx) GetSignBytes() []byte {
	ntx := *tx
	ntx.Signature = nil
	return util.TxHash(ntx.Bytes())
}

func (tx *ContractDeployTx) Bytes() []byte {
	b, err := rlp.EncodeToBytes(tx)
	if err != nil {
		panic(err)
	}
	return append([]byte{CONTRACT_DEPLOY}, b...)
}
