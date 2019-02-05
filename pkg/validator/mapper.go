package validator

import (
	"errors"

	"github.com/bluele/hypermint/pkg/abci/types"
	"github.com/ethereum/go-ethereum/common"
)

var (
	ErrKeyNotFound     = errors.New("key not found")
	ErrKeyAlreadyExist = errors.New("key already exist")
)

type ValidatorMapper interface {
	Set(types.Context, *Validator) error
	Get(types.Context, common.Address) (*Validator, error)
	Remove(types.Context, common.Address)
}

type validatorMapper struct {
	storeKey types.StoreKey
}

func NewValidatorMapper(storeKey types.StoreKey) ValidatorMapper {
	return &validatorMapper{storeKey: storeKey}
}

func (vm *validatorMapper) Set(ctx types.Context, val *Validator) error {
	_, err := vm.Get(ctx, val.Address)
	if err == nil {
		return ErrKeyAlreadyExist
	} else if err != ErrKeyNotFound {
		return err
	}

	store := vm.getStore(ctx)
	b, err := cdc.MarshalBinaryBare(*val)
	if err != nil {
		return err
	}
	store.Set(val.Address.Bytes(), b)
	return nil
}

func (vm *validatorMapper) Get(ctx types.Context, addr common.Address) (*Validator, error) {
	store := vm.getStore(ctx)
	b := store.Get(addr.Bytes())
	if b == nil {
		return nil, ErrKeyNotFound
	}
	val := new(Validator)
	return val, cdc.UnmarshalBinaryBare(b, val)
}

func (vm *validatorMapper) Remove(ctx types.Context, addr common.Address) {
	store := vm.getStore(ctx)
	store.Delete(addr.Bytes())
}

func (vm *validatorMapper) getStore(ctx types.Context) types.KVStore {
	return ctx.KVStore(vm.storeKey)
}
