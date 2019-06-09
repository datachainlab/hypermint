package transaction

import (
	"fmt"
	"io"

	"github.com/bluele/hypermint/pkg/abci/types"
	perrors "github.com/bluele/hypermint/pkg/errors"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/pkg/errors"
)

type CommonTx struct {
	Code      uint8
	From      common.Address
	Nonce     uint64
	Gas       uint64
	Signature []byte
}

func (tx *CommonTx) ValidateBasic() types.Error {
	if isEmptyAddr(tx.From) {
		return ErrInvalidTx(DefaultCodespace, "tx.From == nil")
	}
	if tx.Gas == 0 {
		return ErrInvalidTx(DefaultCodespace, "tx.Gas == 0")
	}
	if len(tx.Signature) == 0 {
		return ErrInvalidTx(DefaultCodespace, "len(tx.Signature) == 0")
	}
	return nil
}

func (tx *CommonTx) SetSignature(sig []byte) {
	tx.Signature = sig
}

func (tx *CommonTx) verifySignature(hash []byte) error {
	rawPub, err := crypto.Ecrecover(hash, tx.Signature)
	if err != nil {
		return errors.Wrap(err, "crypto.Ecrecover")
	}
	pub, err := crypto.UnmarshalPubkey(rawPub)
	if err != nil {
		return errors.Wrap(err, "crypto.DecompressPubkey")
	}
	signer := crypto.PubkeyToAddress(*pub)
	if signer != tx.From {
		return fmt.Errorf("signer mismatch: %v != %v", signer.Hex(), tx.From.Hex())
	}
	return nil
}

func (tx *CommonTx) VerifySignature(hash []byte) types.Error {
	err := tx.verifySignature(hash)
	if err == nil {
		return nil
	}
	return perrors.ErrInvalidSignature(perrors.DefaultCodespace, err.Error())
}

func (tx CommonTx) EncodeRLP(w io.Writer) error {
	if err := tx.EncodeNoSig(w); err != nil {
		return err
	}
	if err := rlp.Encode(w, tx.Signature); err != nil {
		return err
	}
	return nil
}

func (tx CommonTx) EncodeNoSig(w io.Writer) error {
	if err := rlp.Encode(w, tx.Code); err != nil {
		return err
	}
	if err := rlp.Encode(w, tx.From); err != nil {
		return err
	}
	if err := rlp.Encode(w, tx.Nonce); err != nil {
		return err
	}
	if err := rlp.Encode(w, tx.Gas); err != nil {
		return err
	}
	return nil
}

func (tx *CommonTx) DecodeRLP(s *rlp.Stream) error {
	if err := s.Decode(&tx.Code); err != nil {
		return err
	}
	if err := s.Decode(&tx.From); err != nil {
		return err
	}
	if err := s.Decode(&tx.Nonce); err != nil {
		return err
	}
	if err := s.Decode(&tx.Gas); err != nil {
		return err
	}
	b, err := s.Bytes()
	if err != nil {
		return err
	}
	tx.Signature = b
	return nil
}
