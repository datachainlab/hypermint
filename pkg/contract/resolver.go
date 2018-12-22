package contract

import (
	"fmt"
	"log"

	"github.com/bluele/hypermint/pkg/abci/types"
	"github.com/perlin-network/life/exec"
)

type Resolver struct {
	db types.KVStore
}

func NewResolver(db types.KVStore) *Resolver {
	return &Resolver{db: db}
}

func (r *Resolver) getF(cb func(*exec.VirtualMachine, *Process) int64) exec.FunctionImport {
	return func(vm *exec.VirtualMachine) int64 {
		ps := NewProcess(vm, r.db)
		return cb(vm, ps)
	}
}

// ResolveFunc defines a set of import functions that may be called within a WebAssembly module.
func (r *Resolver) ResolveFunc(module, field string) exec.FunctionImport {
	log.Printf("Resolve func: %s %s\n", module, field)
	switch module {
	case "env":
		switch field {
		case "__read_str":
			return r.getF(func(vm *exec.VirtualMachine, ps *Process) int64 {
				var key []byte
				{
					ptr := int(uint32(vm.GetCurrentFrame().Locals[0]))
					msgLen := int(uint32(vm.GetCurrentFrame().Locals[1]))
					key = vm.Memory[ptr : ptr+msgLen]
				}
				valPtr := uint32(vm.GetCurrentFrame().Locals[2])
				msgLen := uint32(vm.GetCurrentFrame().Locals[3])

				return ps.ReadStr(key, valPtr, msgLen)
			})
		case "__write_str":
			return r.getF(func(vm *exec.VirtualMachine, ps *Process) int64 {
				var key, value []byte
				{
					ptr := int(uint32(vm.GetCurrentFrame().Locals[0]))
					msgLen := int(uint32(vm.GetCurrentFrame().Locals[1]))
					key = vm.Memory[ptr : ptr+msgLen]
				}
				{
					ptr := int(uint32(vm.GetCurrentFrame().Locals[2]))
					msgLen := int(uint32(vm.GetCurrentFrame().Locals[3]))
					value = vm.Memory[ptr : ptr+msgLen]
				}
				return ps.WriteStr(key, value)
			})
		default:
			panic(fmt.Errorf("unknown field: %s", field))
		}
	default:
		panic(fmt.Errorf("unknown module: %s", module))
	}
}

// ResolveGlobal defines a set of global variables for use within a WebAssembly module.
func (r *Resolver) ResolveGlobal(module, field string) int64 {
	log.Printf("Resolve global: %s %s\n", module, field)
	switch module {
	case "env":
		switch field {
		case "__life_magic":
			return 424
		default:
			panic(fmt.Errorf("unknown field: %s", field))
		}
	default:
		panic(fmt.Errorf("unknown module: %s", module))
	}
}
