package transaction

import (
	"encoding/binary"

	"github.com/bluele/hypermint/pkg/abci/types"
)

var (
	txIndexKey = []byte("k")
)

type TxIndexMapper interface {
	Get(types.Context) uint32
	Incr(types.Context) uint32
}

type txIndexMapper struct {
	storeKey types.StoreKey
}

func NewTxIndexMapper(k types.StoreKey) TxIndexMapper {
	return &txIndexMapper{storeKey: k}
}

func (m *txIndexMapper) Get(ctx types.Context) uint32 {
	return m.get(m.getStore(ctx))
}

func (m *txIndexMapper) get(kvs types.KVStore) uint32 {
	height := kvs.Get(txIndexKey)
	if height == nil {
		return 0
	}
	return binary.BigEndian.Uint32(height)
}

func (m *txIndexMapper) set(kvs types.KVStore, height uint32) {
	b := make([]byte, 4)
	binary.BigEndian.PutUint32(b, height)
	kvs.Set(txIndexKey, b)
}

func (m *txIndexMapper) Incr(ctx types.Context) uint32 {
	s := m.getStore(ctx)
	height := m.get(s) + 1
	m.set(s, height)
	return height
}

func (m *txIndexMapper) getStore(ctx types.Context) types.KVStore {
	return ctx.KVStore(m.storeKey)
}
