package contract

import (
	"fmt"
	"log"

	"github.com/bluele/hypermint/pkg/abci/types"
	"github.com/perlin-network/life/exec"
)

type VMManager struct{}

func NewVMManager() *VMManager {
	return &VMManager{}
}

func (vm *VMManager) GetVM(db types.KVStore, c *Contract) (*VM, error) {
	v, err := exec.NewVirtualMachine(c.Code, exec.VMConfig{
		EnableJIT:          false,
		DefaultMemoryPages: 128,
		DefaultTableSize:   65536,
	}, NewResolver(db), nil)
	if err != nil {
		return nil, err
	}
	return &VM{VirtualMachine: v}, nil
}

type VM struct {
	*exec.VirtualMachine
}

// TODO calc gas cost
func (vm *VM) ExecContract(entry string) error {
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
