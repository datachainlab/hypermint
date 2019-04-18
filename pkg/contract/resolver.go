package contract

import (
	"bytes"
	"fmt"
	"log"

	"github.com/bluele/hypermint/pkg/util"
	"github.com/ethereum/go-ethereum/common"
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
	switch module {
	case "env":
		switch field {
		case "__get_sender":
			return r.getF(func(vm *exec.VirtualMachine, ps *Process) int64 {
				cf := vm.GetCurrentFrame()
				sender := ps.GetSender()
				ret := &BytesValue{
					mem:  vm.Memory,
					ptr:  uint32(cf.Locals[0]),
					size: uint32(cf.Locals[1]),
				}
				if err := ret.Set(sender[:]); err != nil {
					log.Println("error: ", err)
					return -1
				}
				return 0
			})
		case "__get_arg":
			return r.getF(func(vm *exec.VirtualMachine, ps *Process) int64 {
				cf := vm.GetCurrentFrame()
				idx := cf.Locals[0]
				ret := &BytesValue{
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
				ret := &BytesValue{
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
				ret, err := ps.WriteState(key, value)
				if err != nil {
					return -1
				}
				return ret
			})
		case "__log":
			return r.getF(func(vm *exec.VirtualMachine, _ *Process) int64 {
				msg := readBytes(vm, 0, 1)
				log.Printf("__log: %v(%v)", string(msg), msg)
				return 0
			})
		case "__set_response":
			return r.getF(func(vm *exec.VirtualMachine, ps *Process) int64 {
				value := readBytes(vm, 0, 1)
				ps.SetResponse(value)
				return 0
			})
		case "__call_contract":
			return r.getF(func(vm *exec.VirtualMachine, ps *Process) int64 {
				addr := common.BytesToAddress(readBytes(vm, 0, 1))
				entry := string(readBytes(vm, 2, 3))
				cf := vm.GetCurrentFrame()
				ret := &BytesValue{
					mem:  vm.Memory,
					ptr:  uint32(cf.Locals[4]),
					size: uint32(cf.Locals[5]),
				}
				args, err := readArgs(vm, int(cf.Locals[6]), uint32(cf.Locals[7]))
				if err != nil {
					log.Println("error: ", err)
					return -1
				}
				env, err := ps.EnvManager.Get(ps.Env.Context, r.env.Sender, addr, args)
				if err != nil {
					log.Println("error: ", err)
					return -1
				}
				res, err := env.Exec(ps.Env.Context, entry)
				if err != nil {
					log.Println("error: ", err)
					return -1
				}
				if err := ret.Set(res.Response); err != nil {
					log.Println("error: ", err)
					return -1
				}
				env.state.Add(res.RWSets)
				return int64(len(res.Response))
			})
		case "__ecrecover":
			return r.getF(func(vm *exec.VirtualMachine, ps *Process) int64 {
				h := readBytes(vm, 0, 1)
				v := readBytes(vm, 2, 3)
				r := readBytes(vm, 4, 5)
				s := readBytes(vm, 6, 7)
				cf := vm.GetCurrentFrame()
				ret := &BytesValue{
					mem:  vm.Memory,
					ptr:  uint32(cf.Locals[8]),
					size: uint32(cf.Locals[9]),
				}
				pub, err := util.Ecrecover(h, v, r, s)
				if err != nil {
					log.Println("error: ", err)
					return -1
				}
				if err := ret.Set(pub[:]); err != nil {
					log.Println("error: ", err)
					return -1
				}
				return 0
			})
		case "__ecrecover_address":
			return r.getF(func(vm *exec.VirtualMachine, ps *Process) int64 {
				h := readBytes(vm, 0, 1)
				v := readBytes(vm, 2, 3)
				r := readBytes(vm, 4, 5)
				s := readBytes(vm, 6, 7)
				cf := vm.GetCurrentFrame()
				ret := &BytesValue{
					mem:  vm.Memory,
					ptr:  uint32(cf.Locals[8]),
					size: uint32(cf.Locals[9]),
				}
				addr, err := util.EcrecoverAddress(h, v, r, s)
				if err != nil {
					log.Println("error: ", err)
					return -1
				}
				if err := ret.Set(addr[:]); err != nil {
					log.Println("error: ", err)
					return -1
				}
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
	b := make([]byte, msgLen)
	copy(b, vm.Memory[ptr:ptr+msgLen])
	return b
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
