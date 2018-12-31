package contract

import (
	"fmt"
	"log"

	"github.com/bluele/hypermint/pkg/abci/types"
	"github.com/perlin-network/life/exec"
)

type Env struct {
	Contract *Contract

	DB   types.KVStore
	Args []string
}

func (env *Env) GetVM() (*VM, error) {
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
func (env *Env) Exec(entry string) error {
	vm, err := env.GetVM()
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
	log.Println(ret)
	return nil
}
