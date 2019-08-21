package testutil

import (
	"github.com/bluele/hypermint/pkg/abci/store"
	"github.com/bluele/hypermint/pkg/abci/types"
	"github.com/tendermint/tm-db"
)

func GetTestCommitMultiStore(key types.StoreKey) (types.CommitMultiStore, error) {
	memdb := db.NewMemDB()
	cms := store.NewCommitMultiStore(memdb)
	cms.MountStoreWithDB(key, types.StoreTypeIAVL, nil)
	if err := cms.LoadLatestVersion(); err != nil {
		return nil, err
	}
	return cms, nil
}
