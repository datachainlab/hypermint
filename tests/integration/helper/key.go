package helper

import (
	"crypto/ecdsa"

	"github.com/bluele/hypermint/pkg/util/wallet"
	"github.com/ethereum/go-ethereum/accounts/keystore"
)

func GetPrivKey(ks *keystore.KeyStore, mnemonic, path string) *ecdsa.PrivateKey {
	prv, err := wallet.GetPrivKeyWithMnemonic(mnemonic, path)
	if err != nil {
		panic(err)
	}
	if ks != nil {
		_, err = ks.ImportECDSA(prv, "password")
		if err != nil {
			panic(err)
		}
	}
	return prv
}
