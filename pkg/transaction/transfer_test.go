package transaction

import (
	"fmt"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/assert"
	cmn "github.com/tendermint/tendermint/libs/common"
)

func TestTransferTxEncoding(t *testing.T) {
	var cases = []struct {
		tx          *TransferTx
		decodeError bool
	}{
		{
			&TransferTx{
				To:     common.BytesToAddress(cmn.RandBytes(20)),
				Amount: 1024,
				Common: CommonTx{
					Code:      TRANSFER,
					From:      common.BytesToAddress(cmn.RandBytes(20)),
					Nonce:     1,
					Gas:       1,
					Signature: cmn.RandBytes(65),
				},
			},
			false,
		},
		{
			&TransferTx{
				To:     common.BytesToAddress(cmn.RandBytes(20)),
				Amount: 1024,
				Common: CommonTx{
					Code:      0,
					From:      common.BytesToAddress(cmn.RandBytes(20)),
					Nonce:     1,
					Gas:       1,
					Signature: cmn.RandBytes(65),
				},
			},
			true,
		},
	}

	for i, cs := range cases {
		t.Run(fmt.Sprint(i), func(t *testing.T) {
			assert := assert.New(t)
			b := cs.tx.Bytes()
			tx1, err := DecodeTransferTx(b)
			assert.NoError(err)
			assert.Equal(cs.tx, tx1)

			tx2, err := DecodeTx(b)
			if cs.decodeError {
				assert.Error(err)
				return
			}
			assert.NoError(err)
			tx3, ok := tx2.(*TransferTx)
			assert.True(ok)
			assert.NotNil(tx3)
		})
	}
}

func TestTransferTxSignature(t *testing.T) {
	fprv, err := crypto.GenerateKey()
	if err != nil {
		t.Fatal(err)
	}
	tprv, err := crypto.GenerateKey()
	if err != nil {
		t.Fatal(err)
	}

	var cases = []struct {
		tx    *TransferTx
		valid bool
	}{
		{
			&TransferTx{
				Common: CommonTx{
					From: crypto.PubkeyToAddress(fprv.PublicKey),
					Gas:  0,
				},
				To:     crypto.PubkeyToAddress(tprv.PublicKey),
				Amount: 100,
			},
			true,
		},
		{
			&TransferTx{
				Common: CommonTx{
					From: crypto.PubkeyToAddress(fprv.PublicKey),
					Gas:  0,
				},
				To:     crypto.PubkeyToAddress(tprv.PublicKey),
				Amount: 0,
			},
			false,
		},
		{
			&TransferTx{
				Common: CommonTx{
					From: crypto.PubkeyToAddress(fprv.PublicKey),
					Gas:  0,
				},
				Amount: 100,
			},
			false,
		},
		{
			&TransferTx{
				Common: CommonTx{
					Gas: 0,
				},
				To:     crypto.PubkeyToAddress(tprv.PublicKey),
				Amount: 100,
			},
			false,
		},
	}

	for id, cs := range cases {
		t.Run(fmt.Sprint(id), func(t *testing.T) {
			assert := assert.New(t)
			tx := cs.tx
			sig, err := crypto.Sign(tx.GetSignBytes(), fprv)
			assert.NoError(err)
			tx.SetSignature(sig)

			terr := tx.ValidateBasic()
			if cs.valid {
				assert.Nil(terr)
			} else {
				assert.NotNil(terr)
			}
		})
	}
}
