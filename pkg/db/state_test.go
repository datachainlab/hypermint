package db

import (
	"fmt"
	"testing"

	"github.com/bluele/hypermint/pkg/abci/types"
	"github.com/bluele/hypermint/pkg/testutil"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
	abci "github.com/tendermint/tendermint/abci/types"
)

func TestCommitState(t *testing.T) {
	var first = common.Address([20]byte{19: 1})
	var second = common.Address([20]byte{19: 2})
	_ = second

	ver := Version{Height: 1, TxIdx: 1}

	var R = func(key string) Read {
		return makeRead(key, ver)
	}
	var W = func(key, value string) Write {
		return makeWrite(key, value)
	}

	var cases = []struct {
		sets  []*RWSet
		valid bool
	}{
		{
			[]*RWSet{
				makeRWSet(
					first,
					ver,
					[]Read{R("a1")},
					[]Write{W("b1", "v1")},
				),
				makeRWSet(
					first,
					ver,
					[]Read{R("a2")},
					[]Write{},
				),
			},
			true,
		},
		{
			[]*RWSet{
				makeRWSet(
					first,
					ver,
					[]Read{R("a1")},
					[]Write{W("a1", "v1")},
				),
				makeRWSet(
					second,
					ver,
					[]Read{R("a1")},
					[]Write{},
				),
			},
			true,
		},
		{
			[]*RWSet{
				makeRWSet(
					first,
					ver,
					[]Read{R("a1")},
					[]Write{W("a1", "v1")},
				),
				makeRWSet(
					first,
					ver,
					[]Read{R("a1")},
					[]Write{},
				),
			},
			false,
		},
		{
			[]*RWSet{
				makeRWSet(
					first,
					ver,
					[]Read{R("a1")},
					[]Write{},
				),
				makeRWSet(
					first,
					ver,
					[]Read{R("a1")},
					[]Write{W("a1", "v1")},
				),
			},
			false,
		},
		{
			[]*RWSet{
				makeRWSet(
					first,
					ver,
					[]Read{R("a1"), R("a2")},
					[]Write{W("a1", "v1")},
				),
				makeRWSet(
					first,
					ver,
					[]Read{R("a1"), R("a3")},
					[]Write{W("a3", "v1")},
				),
			},
			false,
		},
		{
			[]*RWSet{
				makeRWSet(
					first,
					ver,
					[]Read{R("a1"), R("a2")},
					[]Write{W("a3", "v1")},
				),
				makeRWSet(
					first,
					ver,
					[]Read{R("a1"), R("a4")},
					[]Write{W("a5", "v1")},
				),
			},
			true,
		},
		{
			[]*RWSet{
				makeRWSet(
					first,
					ver,
					[]Read{R("a1"), R("a2")},
					[]Write{W("a2", "v1")},
				),
				makeRWSet(
					first,
					ver,
					[]Read{R("a1"), R("a3")},
					[]Write{W("a3", "v1")},
				),
			},
			true,
		},
		{
			[]*RWSet{
				makeRWSet(
					first,
					ver,
					[]Read{R("a1")},
					[]Write{W("b1", "v1")},
				),
				makeRWSet(
					first,
					ver,
					[]Read{R("a2")},
					[]Write{W("b1", "v2")},
				),
			},
			true,
		},
		{
			[]*RWSet{
				makeRWSet(
					first,
					ver,
					[]Read{R("a1")},
					[]Write{W("b1", "v1")},
				),
				makeRWSet(
					first,
					ver,
					[]Read{R("b1")},
					[]Write{W("b1", "v2")},
				),
			},
			false,
		},
	}

	testStoreKey := types.NewKVStoreKey("test")
	for i, cs := range cases {
		t.Run(fmt.Sprint(i), func(t *testing.T) {
			assert := assert.New(t)

			cms, err := testutil.GetTestCommitMultiStore(testStoreKey)
			assert.NoError(err)
			ctx := types.NewContext(cms, abci.Header{}, false, nil)

			db := ctx.KVStore(testStoreKey)
			akm := NewKeyMaps()
			err = CommitState(db, cs.sets, ver, akm)
			if cs.valid {
				assert.NoError(err)
			} else {
				assert.Error(err)
			}
		})
	}

}

func makeRWSet(addr common.Address, version Version, rs []Read, ws []Write) *RWSet {
	rsm := NewRWSetMap()
	for _, r := range rs {
		rsm.AddRead(r.Key, r.Version)
	}
	for _, w := range ws {
		rsm.AddWrite(w.Key, w.Value)
	}
	return &RWSet{
		Address: addr,
		Items:   rsm.ToItems(),
	}
}

func makeRead(key string, ver Version) Read {
	return Read{Key: []byte(key), Version: ver}
}

func makeWrite(key, value string) Write {
	return Write{Key: []byte(key), Value: []byte(value)}
}
