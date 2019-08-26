package contract

import (
	"fmt"

	"github.com/bluele/hypermint/pkg/db"
	"github.com/bluele/hypermint/pkg/logger"

	sdk "github.com/bluele/hypermint/pkg/abci/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/perlin-network/life/exec"
)

type Env struct {
	Context  sdk.Context
	Logger   logger.Logger
	Sender   common.Address
	Args     Args
	response []byte

	EnvManager *EnvManager
	Contract   *Contract
	VMProvider VMProvider

	DB     *db.VersionedDB
	state  db.RWSets
	events []*Event
}

type Args struct {
	values [][]byte
}

func (a *Args) Len() int {
	return len(a.values)
}

func (a *Args) PushString(s string) {
	a.values = append(a.values, []byte(s))
}

func (a *Args) PushBytes(b []byte) {
	a.values = append(a.values, b)
}

func (a Args) Get(idx int) []byte {
	return a.values[idx]
}

func NewArgs(bs [][]byte) Args {
	return Args{values: bs}
}

func NewArgsFromStrings(ss []string) Args {
	values := make([][]byte, len(ss))
	for i, s := range ss {
		values[i] = []byte(s)
	}
	return Args{values: values}
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
	Code     int32
	Response []byte
	RWSets   db.RWSets
	Events   []*Event
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
	code := int32(ret)
	if code < 0 {
		return &Result{Code: code}, fmt.Errorf("execute contract error(exit code: %v)", code)
	}
	env.state.Add(&db.RWSet{
		Address: env.Contract.Address(),
		Items:   env.DB.RWSetItems(),
	})
	return &Result{
		Code:     code,
		Response: env.GetReponse(),
		RWSets:   env.state,
		Events:   env.events,
	}, nil
}

func (env *Env) SetResponse(v []byte) {
	env.response = v
}

func (env *Env) GetReponse() []byte {
	return env.response
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

func (em *EnvManager) Get(ctx sdk.Context, sender, addr common.Address, args Args) (*Env, error) {
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
