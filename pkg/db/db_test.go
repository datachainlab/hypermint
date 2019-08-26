package db

import (
	"fmt"
	"testing"

	"github.com/bluele/hypermint/pkg/abci/types"
	"github.com/bluele/hypermint/pkg/testutil"
	"github.com/stretchr/testify/assert"
	abci "github.com/tendermint/tendermint/abci/types"
)

func TestVersionedDB(t *testing.T) {
	type ri struct {
		key    string
		expect string
	}
	type wi struct {
		key   string
		value string
	}
	type B = []byte

	const blockHeight uint32 = 1

	var cases = []struct {
		state   [][]Write
		ops     []interface{}
		expects RWSetItems
	}{
		{
			nil,
			[]interface{}{ri{"a", ""}, wi{"a", "A"}, ri{"a", ""}, wi{"b", "B"}},
			RWSetItems{
				nil,
				[]Write{{B("a"), B("A")}, {B("b"), B("B")}},
			},
		},
		{
			[][]Write{{{B("a"), B("A")}}},
			[]interface{}{ri{"a", "A"}, wi{"a", "A1"}},
			RWSetItems{
				[]Read{{B("a"), Version{blockHeight, 0}}},
				[]Write{{B("a"), B("A1")}},
			},
		},
		{
			[][]Write{{{B("a"), B("A")}}, {{B("b"), B("B")}}},
			[]interface{}{ri{"a", "A"}, wi{"a", "A1"}, ri{"b", "B"}},
			RWSetItems{
				[]Read{{B("a"), Version{blockHeight, 0}}, {B("b"), Version{blockHeight, 1}}},
				[]Write{{B("a"), B("A1")}},
			},
		},
		{
			[][]Write{nil, {{B("a"), B("A")}}, {{B("a"), B("A1")}}, {{B("b"), B("B")}}},
			[]interface{}{ri{"a", "A1"}, wi{"a", "A2"}, ri{"b", "B"}},
			RWSetItems{
				[]Read{{B("a"), Version{blockHeight, 2}}, {B("b"), Version{blockHeight, 3}}},
				[]Write{{B("a"), B("A2")}},
			},
		},
	}

	testStoreKey := types.NewKVStoreKey("test")
	for i, cs := range cases {
		t.Run(fmt.Sprint(i), func(t *testing.T) {
			assert := assert.New(t)
			cms, err := testutil.GetTestCommitMultiStore(testStoreKey)
			assert.NoError(err)
			ctx := types.NewContext(cms, abci.Header{}, false, nil)

			for i, ws := range cs.state {
				vdb := NewVersionedDB(ctx.KVStore(testStoreKey), Version{blockHeight, uint32(i)})
				vdb.rwm.ws = ws
				commitState(ctx.KVStore(testStoreKey), vdb.RWSetItems(), Version{blockHeight, uint32(i)})
			}
			vdb := NewVersionedDB(ctx.KVStore(testStoreKey), Version{blockHeight, uint32(len(cs.state))})

			for _, op := range cs.ops {
				switch op := op.(type) {
				case ri:
					v, err := vdb.Get([]byte(op.key))
					if op.expect != "" {
						assert.NoError(err)
					} else {
						assert.Equal(op.expect, string(v))
					}
				case wi:
					assert.NoError(vdb.Set([]byte(op.key), []byte(op.value)))
				default:
					t.Fatalf("unknown type %T", op)
				}
			}

			set := vdb.RWSetItems()
			assert.Equal(cs.expects, *set)
		})
	}
}
