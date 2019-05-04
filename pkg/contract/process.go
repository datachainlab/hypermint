package contract

import (
	"encoding/binary"
	"fmt"

	"github.com/bluele/hypermint/pkg/db"
	"github.com/bluele/hypermint/pkg/logger"

	"github.com/ethereum/go-ethereum/common"
	"github.com/perlin-network/life/exec"
)

var defaultLogger = logger.GetDefaultLogger("*:debug").With("module", "process")

type Process interface {
	Logger() logger.Logger
	Sender() common.Address
	Args() Args
	State() db.StateDB
	SetResponse([]byte)
	Call(addr common.Address, entry []byte, args Args) (int, error)
	Read(id int) ([]byte, error)
	ValueTable() ValueTable
}

type ValueTable interface {
	Get(id int) ([]byte, error)
	Put(v []byte) (int, error)
}

type process struct {
	vm     *exec.VirtualMachine
	env    *Env
	rwm    db.RWSetMap
	logger logger.Logger
	vt     ValueTable
}

func NewProcess(vm *exec.VirtualMachine, env *Env, logger logger.Logger, vt ValueTable) Process {
	if logger == nil {
		logger = defaultLogger
	}
	return &process{
		vm:     vm,
		env:    env,
		logger: logger,
		vt:     vt,
	}
}

func (p process) Sender() common.Address {
	return p.env.Sender
}

func (p process) Args() Args {
	return p.env.Args
}

func (p process) State() db.StateDB {
	return p.env.DB
}

func (p *process) Call(addr common.Address, entry []byte, args Args) (int, error) {
	env, err := p.env.EnvManager.Get(p.env.Context, p.env.Sender, addr, args)
	if err != nil {
		return -1, err
	}
	res, err := env.Exec(p.env.Context, string(entry))
	if err != nil {
		return -1, err
	}
	p.env.state.Add(res.RWSets)
	return p.ValueTable().Put(res.Response)
}

func (p *process) Read(id int) ([]byte, error) {
	return p.ValueTable().Get(id)
}

func (p *process) SetResponse(v []byte) {
	p.env.SetResponse(v)
}

func (p *process) Logger() logger.Logger {
	if p.logger != nil {
		return p.logger
	}
	return defaultLogger
}

func (p process) ValueTable() ValueTable {
	return p.vt
}

type valueT map[int][]byte

func (vt valueT) Get(id int) ([]byte, error) {
	v, ok := vt[id]
	if !ok {
		return nil, fmt.Errorf("id '%v' not found", id)
	}
	return v, nil
}

func (vt valueT) Put(v []byte) (int, error) {
	cid := len(vt)
	if _, ok := vt[cid]; ok {
		return -1, fmt.Errorf("id '%v' already exists", cid)
	}
	vt[cid] = v
	return cid, nil
}

func DeserializeArgs(b []byte) (Args, error) {
	var args Args
	argc := int(binary.BigEndian.Uint32(b[0:4]))
	var offset uint32 = 4
	for i := 0; i < argc; i++ {
		size := binary.BigEndian.Uint32(b[offset : offset+4])
		bs := make([]byte, 4+size)
		copy(bs[:], b[offset+4:offset+4+size])
		offset += 4 + size
		args.PushBytes(bs)
	}
	return args, nil
}
