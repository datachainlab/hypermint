package transaction

import (
	"github.com/bluele/hypermint/pkg/abci/types"
	"github.com/bluele/hypermint/pkg/util"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/tendermint/tendermint/crypto"
)

var _ Transaction = &ValidatorAddTx{}

type ValidatorAddTx struct {
	PubKey     crypto.PubKey
	Amount     uint64
	Commitment []byte
	CommonTx
}

func (tx *ValidatorAddTx) Decode(b []byte) error {
	return rlp.DecodeBytes(b, tx)
}

func (tx *ValidatorAddTx) ValidateBasic() types.Error {
	if err := tx.CommonTx.ValidateBasic(); err != nil {
		return err
	}
	if tx.Amount == 0 {
		return ErrInvalidTransfer(DefaultCodespace, "tx.Amount == 0")
	}
	return tx.VerifySignature(tx.GetSignBytes())
}

func (tx *ValidatorAddTx) GetSignBytes() []byte {
	ntx := *tx
	ntx.Signature = nil
	b, err := rlp.EncodeToBytes(ntx)
	if err != nil {
		panic(err)
	}
	return util.TxHash(b)
}

func (tx *ValidatorAddTx) Bytes() []byte {
	b, err := rlp.EncodeToBytes(tx)
	if err != nil {
		panic(err)
	}
	return append([]byte{VALIDATOR_ADD}, b...)
}

var _ Transaction = &ValidatorRemoveTx{}

type ValidatorRemoveTx struct {
	PubKey crypto.PubKey
	CommonTx
}

func (tx *ValidatorRemoveTx) Decode(b []byte) error {
	return rlp.DecodeBytes(b, tx)
}

func (tx *ValidatorRemoveTx) ValidateBasic() types.Error {
	if err := tx.CommonTx.ValidateBasic(); err != nil {
		return err
	}
	return tx.VerifySignature(tx.GetSignBytes())
}

func (tx *ValidatorRemoveTx) GetSignBytes() []byte {
	ntx := *tx
	ntx.Signature = nil
	b, err := rlp.EncodeToBytes(ntx)
	if err != nil {
		panic(err)
	}
	return util.TxHash(b)
}

func (tx *ValidatorRemoveTx) Bytes() []byte {
	b, err := rlp.EncodeToBytes(tx)
	if err != nil {
		panic(err)
	}
	return append([]byte{VALIDATOR_REMOVE}, b...)
}
