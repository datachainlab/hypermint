package validator

import (
	"github.com/tendermint/tendermint/crypto"
	"github.com/tendermint/tendermint/crypto/secp256k1"
	pvm "github.com/tendermint/tendermint/privval"
)

func GenFilePV(keyFilePath, stateFilePath string, prv crypto.PrivKey) *pvm.FilePV {
	privValidator := pvm.GenFilePV(keyFilePath, stateFilePath)
	privValidator.Key.PrivKey = prv
	privValidator.Key.PubKey = prv.PubKey()
	privValidator.Key.Address = prv.PubKey().Address()
	privValidator.Save()
	return privValidator
}

func GenFilePVWithECDSA(keyFilePath, stateFilePath string, prv crypto.PrivKey) *pvm.FilePV {
	var p secp256k1.PrivKeySecp256k1
	copy(p[:], prv.Bytes())
	return GenFilePV(keyFilePath, stateFilePath, p)
}
