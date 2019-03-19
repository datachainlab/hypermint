package db

import (
	"errors"

	"github.com/bluele/hypermint/pkg/abci/types"
)

var (
	ErrKeyNotFound = errors.New("key not found")
)

type Version struct {
	Height uint32
	TxIdx  uint32
}

type VersionedDB struct {
	store   types.KVStore
	rwm     *RWSetMap
	version Version
}

func NewVersionedDB(store types.KVStore, version Version) *VersionedDB {
	return &VersionedDB{
		store:   store,
		rwm:     NewRWSetMap(),
		version: version,
	}
}

func (db *VersionedDB) get(k []byte) (*ValueObject, error) {
	b := db.store.Get(k)
	if b == nil {
		return nil, ErrKeyNotFound
	}
	return BytesToValueObject(b)
}

func (db *VersionedDB) Set(k, v []byte) error {
	db.rwm.AddWrite(k, v)
	return nil
}

func (db *VersionedDB) set(k, v []byte, version Version) {
	db.store.Set(k, (&ValueObject{Value: v, Version: version}).Marshal())
}

func (db *VersionedDB) Get(k []byte) ([]byte, error) {
	vo, err := db.get(k)
	if err != nil {
		return nil, err
	}
	if _, ok := db.rwm.GetRead(k); !ok {
		db.rwm.AddRead(k, vo.Version)
	}
	return vo.Value, nil
}

func (db *VersionedDB) Commit() (*RWSet, error) {
	set := db.rwm.ToSet()
	for _, w := range set.WriteSet {
		db.set(w.Key, w.Value, db.version)
	}
	db.rwm = NewRWSetMap()
	return set, nil
}
