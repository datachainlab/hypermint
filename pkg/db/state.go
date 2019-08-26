package db

import (
	"fmt"

	sdk "github.com/bluele/hypermint/pkg/abci/types"
)

type StateManager struct {
	key sdk.StoreKey
}

func NewStateManager(key sdk.StoreKey) *StateManager {
	return &StateManager{key: key}
}

func (sm StateManager) CommitState(ctx sdk.Context, sets RWSets) {
	db := ctx.KVStore(sm.key)
	version := Version{uint32(ctx.BlockHeight()), ctx.TxIndex()}
	for _, s := range sets {
		commitState(db.Prefix(s.Address.Bytes()), s.Items, version)
	}
}

func commitState(db sdk.KVStore, items *RWSetItems, version Version) {
	for _, w := range items.WriteSet {
		fmt.Printf("Save %v => %v\n", string(w.Key), string(w.Value))
		db.Set(w.Key, (&ValueObject{Value: w.Value, Version: version}).Marshal())
	}
}
