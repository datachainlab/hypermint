package account

import (
	"testing"

	"github.com/bluele/hypermint/pkg/abci/types"
	amino "github.com/tendermint/go-amino"
)

var (
	testStoreKey = types.NewKVStoreKey("test")
	cdc          = amino.NewCodec()
)

func TestTransfer(t *testing.T) {
	// am := NewAccountMapper(testStoreKey, cdc)

	// am.AddBalance()
}
