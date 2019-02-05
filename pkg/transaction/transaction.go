package transaction

import "github.com/bluele/hypermint/pkg/abci/types"

const (
	TRANSFER uint8 = 1 + iota
	CONTRACT_DEPLOY
	CONTRACT_CALL
	VALIDATOR_ADD
	VALIDATOR_REMOVE
	VALIDATOR_DELEGATE
)

type Transaction interface {
	types.Tx
	GetSignBytes() []byte
	SetSignature([]byte)
	Decode(b []byte) error
	Bytes() []byte
}

func Decode(b []byte) (types.Tx, types.Error) {
	if len(b) <= 1 {
		return nil, types.ErrTxDecode("tx must have leading tx type")
	}
	switch b[0] {
	case TRANSFER:
		tx := new(TransferTx)
		if err := tx.Decode(b[1:]); err != nil {
			return nil, types.ErrTxDecode("rlp Decode error:" + err.Error())
		}
		return tx, nil
	case CONTRACT_DEPLOY:
		tx := new(ContractDeployTx)
		if err := tx.Decode(b[1:]); err != nil {
			return nil, types.ErrTxDecode("rlp Decode error:" + err.Error())
		}
		return tx, nil
	case CONTRACT_CALL:
		tx := new(ContractCallTx)
		if err := tx.Decode(b[1:]); err != nil {
			return nil, types.ErrTxDecode("rlp Decode error:" + err.Error())
		}
		return tx, nil
	default:
		return nil, types.ErrTxDecode("unknown tx type")
	}
}
