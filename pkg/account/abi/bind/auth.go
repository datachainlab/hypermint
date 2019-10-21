package bind

import (
	"crypto/ecdsa"
	"errors"
	"github.com/bluele/hypermint/pkg/transaction"
	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

func NewKeyedTransactor(key *ecdsa.PrivateKey) *TransactOpts {
	keyAddr := crypto.PubkeyToAddress(key.PublicKey)
	return &TransactOpts{
		From: keyAddr,
		Signer: func(tx transaction.Transaction, address common.Address) (transaction.Transaction, error) {
			if address != keyAddr {
				return nil, errors.New("not authorized to sign this account")
			}
			signature, err := crypto.Sign(tx.GetSignBytes(), key)
			if err != nil {
				return nil, err
			}
			tx.SetSignature(signature)
			return tx, nil
		},
	}
}

func NewKeyStoreTransactor(keystore *keystore.KeyStore, account accounts.Account) (*TransactOpts, error) {
	return &TransactOpts{
		From: account.Address,
		Signer: func(tx transaction.Transaction, address common.Address) (transaction.Transaction, error) {
			if address != account.Address {
				return nil, errors.New("not authorized to sign this account")
			}
			signature, err := keystore.SignHash(account, tx.GetSignBytes())
			if err != nil {
				return nil, err
			}
			tx.SetSignature(signature)
			return tx, nil
		},
	}, nil
}
