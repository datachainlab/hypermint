package transaction

import (
	"fmt"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
	cmn "github.com/tendermint/tendermint/libs/common"
)

func TestContractDeployTxEncoding(t *testing.T) {
	var cases = []struct {
		tx          *ContractDeployTx
		decodeError bool
	}{
		{
			&ContractDeployTx{
				Code: cmn.RandBytes(1024),
				Common: CommonTx{
					Code:      CONTRACT_DEPLOY,
					From:      common.BytesToAddress(cmn.RandBytes(20)),
					Nonce:     1,
					Gas:       1,
					Signature: cmn.RandBytes(65),
				},
			},
			false,
		},
		{
			&ContractDeployTx{
				Code: cmn.RandBytes(1024),
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
			tx1, err := DecodeContractDeployTx(b)
			assert.NoError(err)
			assert.Equal(cs.tx, tx1)

			tx2, err := DecodeTx(b)
			if cs.decodeError {
				assert.Error(err)
				return
			}
			assert.NoError(err)
			tx3, ok := tx2.(*ContractDeployTx)
			assert.True(ok)
			assert.NotNil(tx3)
		})
	}
}
