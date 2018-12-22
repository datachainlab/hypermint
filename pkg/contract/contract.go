package contract

import (
	"github.com/bluele/hypermint/pkg/transaction"
	"github.com/ethereum/go-ethereum/common"
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

func TxToContract(tx *transaction.ContractDeployTx) *Contract {
	return &Contract{
		Owner: tx.From,
		Code:  tx.Code,
	}
}
