package contract

import (
	"fmt"
	"log"

	"github.com/perlin-network/life/exec"
)

type Resolver struct {
	env *Env
}

func NewResolver(env *Env) *Resolver {
	return &Resolver{env: env}
}

func (r *Resolver) getF(cb func(*exec.VirtualMachine, *Process) int64) exec.FunctionImport {
	return func(vm *exec.VirtualMachine) int64 {
		ps := NewProcess(vm, r.env)
		return cb(vm, ps)
	}
}

// ResolveFunc defines a set of import functions that may be called within a WebAssembly module.
func (r *Resolver) ResolveFunc(module, field string) exec.FunctionImport {
	log.Printf("Resolve func: %s %s\n", module, field)
	switch module {
	case "env":
		switch field {
		case "__get_arg":
			return r.getF(func(vm *exec.VirtualMachine, ps *Process) int64 {
				cf := vm.GetCurrentFrame()
				idx := cf.Locals[0]
				ret := &StringValue{
					mem:  vm.Memory,
					ptr:  uint32(cf.Locals[1]),
					size: uint32(cf.Locals[2]),
				}
				size, err := ps.GetArg(int(idx), ret)
				if err != nil {
					log.Println("error: ", err)
					return -1
				}
				return int64(size)
			})
		case "__read_state":
			return r.getF(func(vm *exec.VirtualMachine, ps *Process) int64 {
				cf := vm.GetCurrentFrame()
				key := readBytes(vm, 0, 1)
				ret := &StringValue{
					mem:  vm.Memory,
					ptr:  uint32(cf.Locals[2]),
					size: uint32(cf.Locals[3]),
				}
				size, err := ps.ReadState(key, ret)
				if err != nil {
					return -1
				}
				return int64(size)
			})
		case "__write_state":
			return r.getF(func(vm *exec.VirtualMachine, ps *Process) int64 {
				key := readBytes(vm, 0, 1)
				value := readBytes(vm, 2, 3)
				return ps.WriteState(key, value)
			})
		case "__log":
			return r.getF(func(vm *exec.VirtualMachine, _ *Process) int64 {
				msg := readBytes(vm, 0, 1)
				log.Printf("__log: %v", string(msg))
				return 0
			})
		case "__set_response":
			return r.getF(func(vm *exec.VirtualMachine, ps *Process) int64 {
				value := readBytes(vm, 0, 1)
				ps.SetResponse(value)
				return 0
			})
		default:
			panic(fmt.Errorf("unknown field: %s", field))
		}
	default:
		panic(fmt.Errorf("unknown module: %s", module))
	}
}

func readBytes(vm *exec.VirtualMachine, ptrIdx, sizeIdx int) []byte {
	ptr := int(uint32(vm.GetCurrentFrame().Locals[ptrIdx]))
	msgLen := int(uint32(vm.GetCurrentFrame().Locals[sizeIdx]))
	return vm.Memory[ptr : ptr+msgLen]
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
