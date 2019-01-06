package contract

import (
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/perlin-network/life/exec"
)

type Process struct {
	vm *exec.VirtualMachine
	*Env
}

func NewProcess(vm *exec.VirtualMachine, env *Env) *Process {
	return &Process{
		vm:  vm,
		Env: env,
	}
}

type Value interface {
	Set([]byte) error
}

type StringValue struct {
	mem  []byte
	ptr  uint32
	size uint32
}

func (sv *StringValue) Set(s []byte) error {
	fmt.Println("StringValue.Set:", s)
	if len(s) <= int(sv.size) {
		copy(sv.mem[sv.ptr:], s)
		return nil
	}
	return fmt.Errorf("error: %v >= %v", len(s), int(sv.size))
}

// GetArg returns read size and error or nil
func (p *Process) GetArg(idx int, ret Value) (int, error) {
	if len(p.Env.Args) <= idx {
		return 0, fmt.Errorf("not found idx: %v %v", len(p.Env.Args), idx)
	}
	arg := []byte(p.Env.Args[idx])
	return len(arg), ret.Set(arg)
}

// ReadState returns read size and error or nil
func (p *Process) ReadState(key []byte, ret Value) (int, error) {
	v := p.DB.Get(key)
	if v == nil {
		return 0, fmt.Errorf("key not found")
	}
	return len(v), ret.Set(v)
}

func (p *Process) WriteState(key, v []byte) int64 {
	p.DB.Set(key, v)
	return 0
}

func (p *Process) GetSender() common.Address {
	return p.Env.Contract.Owner
}
