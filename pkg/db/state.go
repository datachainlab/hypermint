package db

import (
	sdk "github.com/bluele/hypermint/pkg/abci/types"
)

type StateManager struct {
	key sdk.StoreKey
}

func NewStateManager(key sdk.StoreKey) *StateManager {
	return &StateManager{key: key}
}

func (sm StateManager) CommitState(ctx sdk.Context, sets RWSets) {
	CommitState(ctx.KVStore(sm.key), sets, Version{uint32(ctx.BlockHeight()), ctx.TxIndex()})
}

func CommitState(db sdk.KVStore, sets RWSets, version Version) {
	for _, s := range sets {
		commitState(db.Prefix(s.Address.Bytes()), s.Items, version)
	}
}

func commitState(db sdk.KVStore, items *RWSetItems, version Version) {
	for _, w := range items.WriteSet {
		db.Set(w.Key, (&ValueObject{Value: w.Value, Version: version}).Marshal())
	}
}
