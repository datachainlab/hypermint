package contract

import (
	"bytes"
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

func (r *Resolver) getF(cb func(*exec.VirtualMachine, Process) int64) exec.FunctionImport {
	return func(vm *exec.VirtualMachine) int64 {
		ps := NewProcess(vm, r.env, r.env.Logger)
		return cb(vm, ps)
	}
}

// ResolveFunc defines a set of import functions that may be called within a WebAssembly module.
func (r *Resolver) ResolveFunc(module, field string) exec.FunctionImport {
	switch module {
	case "env":
		switch field {
		case "__get_sender":
			return r.getF(func(vm *exec.VirtualMachine, ps Process) int64 {
				cf := vm.GetCurrentFrame()
				w := NewWriter(vm.Memory, cf.Locals[0], cf.Locals[1])
				return int64(GetSender(ps, w))
			})
		case "__get_arg":
			return r.getF(func(vm *exec.VirtualMachine, ps Process) int64 {
				cf := vm.GetCurrentFrame()
				idx := int(cf.Locals[0])
				w := NewWriter(vm.Memory, cf.Locals[1], cf.Locals[2])
				return int64(GetArg(ps, idx, w))
			})
		case "__read_state":
			return r.getF(func(vm *exec.VirtualMachine, ps Process) int64 {
				cf := vm.GetCurrentFrame()
				key := NewReader(vm.Memory, cf.Locals[0], cf.Locals[1])
				buf := NewWriter(vm.Memory, cf.Locals[2], cf.Locals[3])
				return int64(ReadState(ps, key, buf))
			})
		case "__write_state":
			return r.getF(func(vm *exec.VirtualMachine, ps Process) int64 {
				cf := vm.GetCurrentFrame()
				key := NewReader(vm.Memory, cf.Locals[0], cf.Locals[1])
				val := NewReader(vm.Memory, cf.Locals[2], cf.Locals[3])
				return int64(WriteState(ps, key, val))
			})
		case "__log":
			return r.getF(func(vm *exec.VirtualMachine, ps Process) int64 {
				cf := vm.GetCurrentFrame()
				msg := NewReader(vm.Memory, cf.Locals[0], cf.Locals[1])
				return int64(Log(ps, msg))
			})
		case "__set_response":
			return r.getF(func(vm *exec.VirtualMachine, ps Process) int64 {
				cf := vm.GetCurrentFrame()
				return int64(SetResponse(ps, NewReader(vm.Memory, cf.Locals[0], cf.Locals[1])))
			})
		case "__call_contract":
			return r.getF(func(vm *exec.VirtualMachine, ps Process) int64 {
				cf := vm.GetCurrentFrame()
				addr := NewReader(vm.Memory, cf.Locals[0], cf.Locals[1])
				entry := NewReader(vm.Memory, cf.Locals[2], cf.Locals[3])
				ret := NewWriter(vm.Memory, cf.Locals[4], cf.Locals[5])
				args, err := readArgs(vm, int(cf.Locals[6]), uint32(cf.Locals[7]))
				if err != nil {
					log.Println("error: ", err)
					return -1
				}
				return int64(CallContract(ps, addr, entry, args, ret))
			})
		case "__ecrecover":
			return r.getF(func(vm *exec.VirtualMachine, ps Process) int64 {
				cf := vm.GetCurrentFrame()
				h := NewReader(vm.Memory, cf.Locals[0], cf.Locals[1])
				v := NewReader(vm.Memory, cf.Locals[2], cf.Locals[3])
				r := NewReader(vm.Memory, cf.Locals[4], cf.Locals[5])
				s := NewReader(vm.Memory, cf.Locals[6], cf.Locals[7])
				ret := NewWriter(vm.Memory, cf.Locals[8], cf.Locals[9])
				return int64(ECRecover(ps, h, v, r, s, ret))
			})
		case "__ecrecover_address":
			return r.getF(func(vm *exec.VirtualMachine, ps Process) int64 {
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

func readArgs(vm *exec.VirtualMachine, argc int, argvPtr uint32) (Args, error) {
	var args Args

	buf := bytes.NewBuffer(nil)
	cur := argvPtr
	num := 0
	for num < argc {
		b := vm.Memory[cur]
		cur++
		if b == 0 {
			args.PushBytes(buf.Bytes())
			buf.Reset()
			num++
			continue
		}
		buf.WriteByte(b)
	}

	return args, nil
}

// ResolveGlobal defines a set of global variables for use within a WebAssembly module.
func (r *Resolver) ResolveGlobal(module, field string) int64 {
	panic(fmt.Errorf("not supported module: %s %s", module, field))
}
