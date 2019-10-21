package bind

import "github.com/ethereum/go-ethereum/accounts/keystore"

func NewKeyStore(keydir string) *keystore.KeyStore {
	return keystore.NewKeyStore(keydir, keystore.StandardScryptN, keystore.StandardScryptP)
}
