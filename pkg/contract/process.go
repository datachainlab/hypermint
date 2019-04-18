package contract

import (
	"fmt"

	"github.com/bluele/hypermint/pkg/db"

	"github.com/ethereum/go-ethereum/common"
	"github.com/perlin-network/life/exec"
)

type Process struct {
	vm *exec.VirtualMachine
	*Env
	rwm db.RWSetMap
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

type BytesValue struct {
	mem  []byte
	ptr  uint32
	size uint32
}

func (bv *BytesValue) Set(b []byte) error {
	if len(b) <= int(bv.size) {
		copy(bv.mem[bv.ptr:], b)
		return nil
	}
	return fmt.Errorf("allocation error: %v >= %v", len(b), int(bv.size))
}

// GetArg returns read size and error or nil
func (p *Process) GetArg(idx int, ret Value) (int, error) {
	if p.Env.Args.Len() <= idx {
		return 0, fmt.Errorf("not found idx: %v %v", p.Env.Args.Len(), idx)
	}
	arg := p.Env.Args.Get(idx)
	return len(arg), ret.Set(arg)
}

// ReadState returns read size and error or nil
func (p *Process) ReadState(key []byte, ret Value) (int, error) {
	v, err := p.DB.Get(key)
	if err != nil {
		return -1, err
	}
	return len(v), ret.Set(v)
}

func (p *Process) WriteState(key, value []byte) (int64, error) {
	if err := p.DB.Set(key, value); err != nil {
		return -1, err
	}
	return 0, nil
}

func (p *Process) GetSender() common.Address {
	return p.Env.Sender
}
