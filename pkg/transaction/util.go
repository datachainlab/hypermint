package transaction

import (
	"time"

	"github.com/ethereum/go-ethereum/common"
)

var emptyAddr common.Address

func GetNonceByAddress(addr common.Address) (uint64, error) {
	return uint64(time.Now().UnixNano()), nil
}

func isEmptyAddr(addr common.Address) bool {
	return addr == emptyAddr
}
