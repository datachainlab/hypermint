package contract

import (
	"errors"
	"fmt"

	"github.com/bluele/hypermint/pkg/db"

	sdk "github.com/bluele/hypermint/pkg/abci/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/perlin-network/life/exec"
)

type Env struct {
	Context  sdk.Context
	Sender   common.Address
	Args     []string
	Response []byte

	EnvManager *EnvManager
	Contract   *Contract
	VMProvider VMProvider

	DB    *db.VersionedDB
	state db.State
}

type VMProvider func(*Env) (*VM, error)

func DefaultVMProvider(env *Env) (*VM, error) {
	v, err := exec.NewVirtualMachine(env.Contract.Code, exec.VMConfig{
		EnableJIT:          false,
		DefaultMemoryPages: 128,
		DefaultTableSize:   65536,
	}, NewResolver(env), nil)
	if err != nil {
		return nil, err
	}
	return &VM{VirtualMachine: v}, nil
}

type VM struct {
	*exec.VirtualMachine
}

type Result struct {
	Response []byte
	RWSets   *db.RWSets
}

func (env *Env) Exec(ctx sdk.Context, entry string) (*Result, error) {
	vmProvider := env.VMProvider
	if vmProvider == nil {
		vmProvider = DefaultVMProvider
	}
	vm, err := vmProvider(env)
	if err != nil {
		return nil, err
	}
	id, ok := vm.GetFunctionExport(entry)
	if !ok {
		return nil, fmt.Errorf("entry point not found")
	}
	ret, err := vm.Run(id)
	if err != nil {
		vm.PrintStackTrace()
		return nil, err
	}
	if ret == -1 {
		return nil, errors.New("execute contract error")
	}
	set, err := env.DB.Commit()
	if err != nil {
		return nil, err
	}
	return &Result{
		Response: env.GetReponse(),
		RWSets: &db.RWSets{
			Address: env.Contract.Address(),
			RWSet:   set,
			Childs:  env.state.Childs,
		},
	}, nil
}

func (env *Env) SetResponse(v []byte) {
	env.Response = v
}

func (env *Env) GetReponse() []byte {
	return env.Response
}

type EnvManager struct {
	key sdk.StoreKey
	cm  ContractMapper
}

func NewEnvManager(key sdk.StoreKey, cm ContractMapper) *EnvManager {
	return &EnvManager{
		key: key,
		cm:  cm,
	}
}

func (em *EnvManager) Get(ctx sdk.Context, sender, addr common.Address, args []string) (*Env, error) {
	c, err := em.cm.Get(ctx, addr)
	if err != nil {
		return nil, err
	}
	return &Env{
		Context:    ctx,
		Sender:     sender,
		EnvManager: em,
		Contract:   c,
		DB:         db.NewVersionedDB(ctx.KVStore(em.key).Prefix(addr.Bytes()), db.Version{uint32(ctx.BlockHeight()), ctx.TxIndex()}),
		Args:       args,
	}, nil
}
