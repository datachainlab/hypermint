package util

import (
	"errors"
	"fmt"
	"testing"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/assert"
	"github.com/tendermint/tendermint/libs/common"
)

func TestEcrecover(t *testing.T) {
	for i := 0; i < 100; i++ {
		t.Run(fmt.Sprint(i), func(t *testing.T) {
			assert := assert.New(t)
			msgHash := crypto.Keccak256(common.RandBytes(32))
			privateKey, err := crypto.GenerateKey()
			assert.NoError(err)

			sig, err := crypto.Sign(msgHash, privateKey)
			assert.NoError(err)

			epub := crypto.FromECDSAPub(&privateKey.PublicKey)

			v, r, s, err := splitSigParams(sig)
			assert.NoError(err)

			apub, err := Ecrecover(msgHash, v, r, s)
			assert.NoError(err)
			assert.Equal(len(epub), len(apub))
			assert.Equal(epub, apub)

			addr, err := EcrecoverAddress(msgHash, v, r, s)
			assert.NoError(err)

			eaddr := crypto.PubkeyToAddress(privateKey.PublicKey)
			assert.Equal(eaddr, addr)
		})
	}
}

func splitSigParams(sig []byte) (v []byte, r []byte, s []byte, err error) {
	if len(sig) != 65 {
		err = errors.New("signature length should be 65")
		return
	}

	v = make([]byte, 1)
	r = make([]byte, 32)
	s = make([]byte, 32)

	copy(r, sig[0:32])
	copy(s, sig[32:64])
	v[0] = sig[64]
	return
}
