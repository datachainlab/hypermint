package transaction

import (
	"time"

	"github.com/ethereum/go-ethereum/common"
)

func GetNonceByAddress(addr common.Address) (uint64, error) {
	return uint64(time.Now().UnixNano()), nil
}
