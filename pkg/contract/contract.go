package contract

import (
	"github.com/bluele/hypermint/pkg/transaction"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/rlp"
)

type Contract struct {
	Owner common.Address
	Code  []byte
}

func (c *Contract) Bytes() []byte {
	return c.Code
}

func (c *Contract) Decode(b []byte) error {
	return rlp.DecodeBytes(b, c)
}

func (c *Contract) Encode() ([]byte, error) {
	return rlp.EncodeToBytes(c)
}

func (c *Contract) Address() common.Address {
	return common.BytesToAddress(crypto.Keccak256(c.Code)[12:])
}

func TxToContract(tx *transaction.ContractDeployTx) *Contract {
	return &Contract{
		Owner: tx.Common.From,
		Code:  tx.Code,
	}
}
