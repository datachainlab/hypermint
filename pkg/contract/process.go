package contract

import (
	"fmt"
	"log"
	"strconv"

	"github.com/bluele/hypermint/pkg/abci/types"
	"github.com/perlin-network/life/exec"
)

type Process struct {
	vm *exec.VirtualMachine
	db types.KVStore
}

func NewProcess(vm *exec.VirtualMachine, db types.KVStore) *Process {
	return &Process{
		vm: vm,
		db: db,
	}
}

func (p *Process) ReadStr(key []byte, valPtr, valLen uint32) int64 {
	v := p.db.Get(key)
	if v == nil {
		return -1
	}

	log.Println("read:", string(key), string(v))

	if len(v) < int(valLen) {
		copy(p.vm.Memory[valPtr:], []byte(v))
		return 0
	}
	return 1
}

func (p *Process) _ReadInt(keyPtr, keyLen, valPtr, valLen uint32) int {
	key := string(readMem(p.vm.Memory, keyPtr, keyLen))

	v := p.db.Get([]byte(key))
	if v == nil {
		return 0
	}
	i, err := strconv.Atoi(string(v))
	if err != nil {
		panic(err)
	}

	return i
}

func (p *Process) WriteStr(key, v []byte) int64 {
	log.Printf("key=%v value=%v", string(key), string(v))
	p.db.Set(key, v)
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
