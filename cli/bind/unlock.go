package bind

import (
	"fmt"
	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/common"
)

func Unlock(ks *keystore.KeyStore, address common.Address, passphrase string) error {
	if !ks.HasAddress(address) {
		return fmt.Errorf("account not found: %v", address.Hex())
	}
	return ks.Unlock(accounts.Account{Address: address}, passphrase)
}
