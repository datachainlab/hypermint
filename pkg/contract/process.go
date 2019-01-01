package contract

import (
	"encoding/binary"
	"fmt"
	"log"
	"strconv"

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
	if len(s) < int(sv.size) {
		copy(sv.mem[sv.ptr:], s)
		return nil
	}
	return fmt.Errorf("%v >= %v", len(s), int(sv.size))
}

type Int64Value struct {
	mem []byte
	ptr uint32
}

func (v *Int64Value) Set(s []byte) error {
	i, err := strconv.Atoi(string(s))
	if err != nil {
		return err
	}
	binary.LittleEndian.PutUint64(v.mem[v.ptr:], uint64(i))
	return nil
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
	log.Println("read:", string(key), string(v))
	return len(v), ret.Set(v)
}

func (p *Process) WriteState(key, v []byte) int64 {
	log.Printf("WriteState key=%v value=%v", string(key), string(v))
	p.DB.Set(key, v)
	return 0
}

func (p *Process) Log(msg []byte) {
	fmt.Printf("[app] %s\n", string(msg))
}

func readMem(inp []byte, off uint32, max uint32) []byte {
	var result []byte

	mem := inp[int(off):]
	for i, bt := range mem {
		if uint32(i) == max {
			return result
		}
		if bt == 0 {
			return result
		}

		result = append(result, bt)
	}

	return result
}
