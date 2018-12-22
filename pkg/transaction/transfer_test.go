package transaction

import (
	"fmt"
	"testing"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/assert"
)

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
		tx            *TransferTx
		expectedError bool
	}{
		{
			&TransferTx{
				CommonTx: CommonTx{
					From: crypto.PubkeyToAddress(fprv.PublicKey),
					Gas:  1,
				},
				To:     crypto.PubkeyToAddress(tprv.PublicKey),
				Amount: 100,
			},
			false,
		},
		{
			&TransferTx{
				CommonTx: CommonTx{
					From: crypto.PubkeyToAddress(fprv.PublicKey),
					Gas:  0,
				},
				To:     crypto.PubkeyToAddress(tprv.PublicKey),
				Amount: 100,
			},
			true,
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
			if cs.expectedError {
				assert.NotNil(terr)
			} else {
				assert.Nil(terr)
			}
		})
	}
}
