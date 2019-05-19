package contract

import (
	"fmt"

	"github.com/perlin-network/life/exec"
)

type Resolver struct {
	env *Env
	vt  ValueTable
}

func NewResolver(env *Env) *Resolver {
	return &Resolver{env: env, vt: make(valueT)}
}

func (r *Resolver) withProcess(cb func(*exec.VirtualMachine, Process) int64) exec.FunctionImport {
	return func(vm *exec.VirtualMachine) int64 {
		ps := NewProcess(vm, r.env, r.env.Logger, r.vt)
		return cb(vm, ps)
	}
}

// ResolveFunc defines a set of import functions that may be called within a WebAssembly module.
func (r *Resolver) ResolveFunc(module, field string) exec.FunctionImport {
	switch module {
	case "env":
		switch field {
		case "__get_sender":
			return r.withProcess(func(vm *exec.VirtualMachine, ps Process) int64 {
				cf := vm.GetCurrentFrame()
				w := NewWriter(vm.Memory, cf.Locals[0], cf.Locals[1])
				return int64(GetSender(ps, w))
			})
		case "__get_arg":
			return r.withProcess(func(vm *exec.VirtualMachine, ps Process) int64 {
				cf := vm.GetCurrentFrame()
				idx := int(cf.Locals[0])
				offset := int(cf.Locals[1])
				w := NewWriter(vm.Memory, cf.Locals[2], cf.Locals[3])
				return int64(GetArg(ps, idx, offset, w))
			})
		case "__read_state":
			return r.withProcess(func(vm *exec.VirtualMachine, ps Process) int64 {
				cf := vm.GetCurrentFrame()
				key := NewReader(vm.Memory, cf.Locals[0], cf.Locals[1])
				offset := int(cf.Locals[2])
				buf := NewWriter(vm.Memory, cf.Locals[3], cf.Locals[4])
				return int64(ReadState(ps, key, offset, buf))
			})
		case "__write_state":
			return r.withProcess(func(vm *exec.VirtualMachine, ps Process) int64 {
				cf := vm.GetCurrentFrame()
				key := NewReader(vm.Memory, cf.Locals[0], cf.Locals[1])
				val := NewReader(vm.Memory, cf.Locals[2], cf.Locals[3])
				return int64(WriteState(ps, key, val))
			})
		case "__log":
			return r.withProcess(func(vm *exec.VirtualMachine, ps Process) int64 {
				cf := vm.GetCurrentFrame()
				msg := NewReader(vm.Memory, cf.Locals[0], cf.Locals[1])
				return int64(Log(ps, msg))
			})
		case "__set_response":
			return r.withProcess(func(vm *exec.VirtualMachine, ps Process) int64 {
				cf := vm.GetCurrentFrame()
				return int64(SetResponse(ps, NewReader(vm.Memory, cf.Locals[0], cf.Locals[1])))
			})
		case "__call_contract":
			return r.withProcess(func(vm *exec.VirtualMachine, ps Process) int64 {
				cf := vm.GetCurrentFrame()
				addr := NewReader(vm.Memory, cf.Locals[0], cf.Locals[1])
				entry := NewReader(vm.Memory, cf.Locals[2], cf.Locals[3])
				argb := NewReader(vm.Memory, cf.Locals[4], cf.Locals[5])
				return int64(CallContract(ps, addr, entry, argb))
			})
		case "__read":
			return r.withProcess(func(vm *exec.VirtualMachine, ps Process) int64 {
				cf := vm.GetCurrentFrame()
				id := int(cf.Locals[0])
				offset := int(cf.Locals[1])
				buf := NewWriter(vm.Memory, cf.Locals[2], cf.Locals[3])
				return int64(Read(ps, id, offset, buf))
			})
		case "__keccak256":
			return r.withProcess(func(vm *exec.VirtualMachine, ps Process) int64 {
				cf := vm.GetCurrentFrame()
				msg := NewReader(vm.Memory, cf.Locals[0], cf.Locals[1])
				buf := NewWriter(vm.Memory, cf.Locals[2], cf.Locals[3])
				return int64(Keccak256(ps, msg, buf))
			})
		case "__sha256":
			return r.withProcess(func(vm *exec.VirtualMachine, ps Process) int64 {
				cf := vm.GetCurrentFrame()
				msg := NewReader(vm.Memory, cf.Locals[0], cf.Locals[1])
				buf := NewWriter(vm.Memory, cf.Locals[2], cf.Locals[3])
				return int64(Sha256(ps, msg, buf))
			})
		case "__ecrecover":
			return r.withProcess(func(vm *exec.VirtualMachine, ps Process) int64 {
				cf := vm.GetCurrentFrame()
				h := NewReader(vm.Memory, cf.Locals[0], cf.Locals[1])
				v := NewReader(vm.Memory, cf.Locals[2], cf.Locals[3])
				r := NewReader(vm.Memory, cf.Locals[4], cf.Locals[5])
				s := NewReader(vm.Memory, cf.Locals[6], cf.Locals[7])
				ret := NewWriter(vm.Memory, cf.Locals[8], cf.Locals[9])
				return int64(ECRecover(ps, h, v, r, s, ret))
			})
		case "__ecrecover_address":
			return r.withProcess(func(vm *exec.VirtualMachine, ps Process) int64 {
				cf := vm.GetCurrentFrame()
				h := NewReader(vm.Memory, cf.Locals[0], cf.Locals[1])
				v := NewReader(vm.Memory, cf.Locals[2], cf.Locals[3])
				r := NewReader(vm.Memory, cf.Locals[4], cf.Locals[5])
				s := NewReader(vm.Memory, cf.Locals[6], cf.Locals[7])
				ret := NewWriter(vm.Memory, cf.Locals[8], cf.Locals[9])
				return int64(ECRecoverAddress(ps, h, v, r, s, ret))
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
	panic(fmt.Errorf("not supported module: %s %s", module, field))
}
