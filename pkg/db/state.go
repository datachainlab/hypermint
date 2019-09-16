package db

import (
	"fmt"

	sdk "github.com/bluele/hypermint/pkg/abci/types"
	"github.com/ethereum/go-ethereum/common"
)

type StateManager struct {
	key sdk.StoreKey
}

func NewStateManager(key sdk.StoreKey) *StateManager {
	return &StateManager{key: key}
}

func (sm StateManager) CommitState(ctx sdk.Context, sets RWSets) error {
	return CommitState(ctx.KVStore(sm.key), sets, Version{uint32(ctx.BlockHeight()), ctx.TxIndex()}, NewKeyMaps())
}

func CommitState(db sdk.KVStore, sets RWSets, version Version, m KeyMaps) error {
	for _, s := range sets {
		if err := commitState(s.Address, db.Prefix(s.Address.Bytes()), s.Items, version, m.GetReadKeyMap(s.Address), m.GetWriteKeyMap(s.Address)); err != nil {
			return err
		}
	}
	return nil
}

func commitState(addr common.Address, db sdk.KVStore, items *RWSetItems, version Version, rm, wm KeyMap) error {
	rkeys := make([]string, 0, len(items.ReadSet))
	for _, r := range items.ReadSet {
		if len(wm) > 0 && wm.Has(string(r.Key)) {
			return fmt.Errorf("ReadSet: conflicted updates exist: address=%v key=%x(%v)", addr, r.Key, string(r.Key))
		}
		rkeys = append(rkeys, string(r.Key))
	}

	wkeys := make([]string, 0, len(items.WriteSet))
	for _, w := range items.WriteSet {
		if len(rm) > 0 && rm.Has(string(w.Key)) {
			return fmt.Errorf("ReadSet: conflicted updates exist: address=%v key=%x(%v)", addr, w.Key, string(w.Key))
		}
		db.Set(w.Key, (&ValueObject{Value: w.Value, Version: version}).Marshal())
		wkeys = append(wkeys, string(w.Key))
	}

	for _, k := range rkeys {
		rm.Set(k)
	}
	for _, k := range wkeys {
		wm.Set(k)
	}

	return nil
}

type KeyMaps struct {
	Read  map[common.Address]KeyMap
	Write map[common.Address]KeyMap
}

func NewKeyMaps() KeyMaps {
	return KeyMaps{
		Read:  make(map[common.Address]KeyMap),
		Write: make(map[common.Address]KeyMap),
	}
}

func (m *KeyMaps) GetReadKeyMap(addr common.Address) KeyMap {
	if _, ok := m.Read[addr]; !ok {
		m.Read[addr] = NewKeyMap()
	}
	return m.Read[addr]
}

func (m *KeyMaps) GetWriteKeyMap(addr common.Address) KeyMap {
	if _, ok := m.Write[addr]; !ok {
		m.Write[addr] = NewKeyMap()
	}
	return m.Write[addr]
}

// KeyMap contains updated key
type KeyMap map[string]struct{}

func NewKeyMap() KeyMap {
	return make(KeyMap)
}

func (m KeyMap) Has(k string) bool {
	_, ok := m[k]
	return ok
}

func (m KeyMap) Set(k string) {
	m[k] = struct{}{}
}
