package validator

import (
	"github.com/tendermint/tendermint/crypto"
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
