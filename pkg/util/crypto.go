package util

import (
	"crypto/ecdsa"

	"github.com/tendermint/tendermint/crypto"
	"github.com/tendermint/tendermint/crypto/secp256k1"
)

func PrvKeyToCryptoKey(prv *ecdsa.PrivateKey) crypto.PrivKey {
	var p secp256k1.PrivKeySecp256k1
	copy(p[:], prv.D.Bytes())
	return p
}
