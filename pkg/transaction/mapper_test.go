package transaction

import (
	"testing"

	"github.com/bluele/hypermint/pkg/abci/types"
	"github.com/bluele/hypermint/pkg/testutil"
	"github.com/stretchr/testify/assert"
	abci "github.com/tendermint/tendermint/abci/types"
)

func TestTxIndexMapper(t *testing.T) {
	assert := assert.New(t)
	k := types.NewTransientStoreKey("test")
	m := NewTxIndexMapper(k)
	cms, err := testutil.GetTestCommitMultiStore(k)
	assert.NoError(err)
	ctx := types.NewContext(cms, abci.Header{}, false, nil)

	assert.Equal(uint32(0), m.Get(ctx))
	m.Incr(ctx)
	assert.Equal(uint32(1), m.Get(ctx))
	m.Incr(ctx)
	assert.Equal(uint32(2), m.Get(ctx))
}
