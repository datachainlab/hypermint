package transaction

import (
	"fmt"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
	cmn "github.com/tendermint/tendermint/libs/common"
)

func TestContractCallTxEncoding(t *testing.T) {
	var cases = []struct {
		tx          *ContractCallTx
		decodeError bool
	}{
		{
			&ContractCallTx{
				Address:    common.BytesToAddress(cmn.RandBytes(20)),
				Func:       cmn.RandStr(20),
				Args:       [][]byte{},
				RWSetsHash: []byte{},
				Common: CommonTx{
					Code:      CONTRACT_CALL,
					From:      common.BytesToAddress(cmn.RandBytes(20)),
					Nonce:     1,
					Gas:       1,
					Signature: cmn.RandBytes(65),
				},
			},
			false,
		},
		{
			&ContractCallTx{
				Address:    common.BytesToAddress(cmn.RandBytes(20)),
				Func:       cmn.RandStr(20),
				Args:       [][]byte{},
				RWSetsHash: []byte{},
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
			tx1, err := DecodeContractCallTx(b)
			assert.NoError(err)
			assert.Equal(cs.tx, tx1)

			tx2, err := DecodeTx(b)
			if cs.decodeError {
				assert.Error(err)
				return
			}
			assert.NoError(err)
			tx3, ok := tx2.(*ContractCallTx)
			assert.True(ok)
			assert.NotNil(tx3)
		})
	}
}
