package bind

import (
	"github.com/bluele/hypermint/pkg/transaction"
	"github.com/ethereum/go-ethereum/common"
)

func MakeContractCallTx(callerAddress, contractAddress common.Address, fn string, args [][]byte, rwh []byte) (transaction.Transaction, error) {
	nonce, err := transaction.GetNonceByAddress(callerAddress)
	if err != nil {
		return nil, err
	}
	var byteArgs [][]byte
	for _, arg := range args {
		byteArgs = append(byteArgs, arg)
	}
	return &transaction.ContractCallTx{
		Common: transaction.CommonTx{
			Code:      transaction.CONTRACT_CALL,
			From:      callerAddress,
			Nonce:     nonce,
			Gas:       1,
			Signature: nil,
		},
		Address:    contractAddress,
		Func:       fn,
		Args:       byteArgs,
		RWSetsHash: rwh,
	}, nil
}
