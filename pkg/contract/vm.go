package contract

import (
	"errors"
	"fmt"

	"github.com/bluele/hypermint/pkg/abci/types"
	sdk "github.com/bluele/hypermint/pkg/abci/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/perlin-network/life/exec"
)

type Env struct {
	Context    sdk.Context
	EnvManager *EnvManager
	Contract   *Contract
	VMProvider VMProvider

	Response []byte

	DB   types.KVStore
	Args []string
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

// TODO calc gas cost
func (env *Env) Exec(ctx sdk.Context, entry string) error {
	vmProvider := env.VMProvider
	if vmProvider == nil {
		vmProvider = DefaultVMProvider
	}
	vm, err := vmProvider(env)
	if err != nil {
		return err
	}
	id, ok := vm.GetFunctionExport(entry)
	if !ok {
		return fmt.Errorf("entry point not found")
	}
	ret, err := vm.Run(id)
	if err != nil {
		vm.PrintStackTrace()
		return err
	}
	if ret == -1 {
		return errors.New("execute contract error")
	}
	fmt.Printf("response is %v\n", string(env.GetReponse()))
	return nil
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

func (em *EnvManager) Get(ctx sdk.Context, addr common.Address, args []string) (*Env, error) {
	c, err := em.cm.Get(ctx, addr)
	if err != nil {
		return nil, err
	}
	return &Env{
		Context:    ctx,
		EnvManager: em,
		Contract:   c,
		DB:         ctx.KVStore(em.key).Prefix(addr.Bytes()),
		Args:       args,
	}, nil
}
