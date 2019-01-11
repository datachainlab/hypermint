package validator

import (
	"crypto/ecdsa"

	gcrypto "github.com/ethereum/go-ethereum/crypto"
	"github.com/tendermint/tendermint/crypto"
	"github.com/tendermint/tendermint/crypto/secp256k1"
	pvm "github.com/tendermint/tendermint/privval"
)

func GenFilePV(path string, prv crypto.PrivKey) *pvm.FilePV {
	privValidator := pvm.GenFilePV(path)
	privValidator.PrivKey = prv
	privValidator.PubKey = prv.PubKey()
	privValidator.Address = prv.PubKey().Address()
	privValidator.Save()
	return privValidator
}

func GenFilePVWithECDSA(path string, prv *ecdsa.PrivateKey) *pvm.FilePV {
	pb := gcrypto.FromECDSA(prv)
	var p secp256k1.PrivKeySecp256k1
	copy(p[:], pb)
	return GenFilePV(path, p)
}

func bytesToECDSAPrvKey(b []byte) *ecdsa.PrivateKey {
	pv, err := gcrypto.ToECDSA(b)
	if err != nil {
		panic(err)
	}
	return pv
}
