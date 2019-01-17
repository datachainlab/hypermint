package contract

import (
	"github.com/bluele/hypermint/pkg/abci/types"
	"github.com/ethereum/go-ethereum/common"
)

type ContractMapper interface {
	Put(ctx types.Context, addr common.Address, c *Contract)
	Get(ctx types.Context, addr common.Address) (*Contract, error)
}

type contractMapper struct {
	storeKey types.StoreKey
}

func NewContractMapper(storeKey types.StoreKey) ContractMapper {
	return &contractMapper{
		storeKey: storeKey,
	}
}

func (cm *contractMapper) Put(ctx types.Context, addr common.Address, c *Contract) {
	cm.put(cm.getStore(ctx), addr, c)
}

func (cm *contractMapper) put(kvs types.KVStore, addr common.Address, c *Contract) {
	b, err := c.Encode()
	if err != nil {
		panic(err)
	}
	kvs.Set(addr.Bytes(), b)
}

func (cm *contractMapper) Get(ctx types.Context, addr common.Address) (*Contract, error) {
	return cm.get(cm.getStore(ctx), addr)
}

func (cm *contractMapper) get(kvs types.KVStore, addr common.Address) (*Contract, error) {
	v := kvs.Get(addr.Bytes())
	if v == nil {
		return nil, ErrContractNotFound
	}
	c := new(Contract)
	if err := c.Decode(v); err != nil {
		return nil, err
	}
	return c, nil
}

func (cm *contractMapper) getStore(ctx types.Context) types.KVStore {
	return ctx.KVStore(cm.storeKey)
}
