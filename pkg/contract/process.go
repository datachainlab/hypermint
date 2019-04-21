package contract

import (
	"github.com/bluele/hypermint/pkg/db"
	"github.com/bluele/hypermint/pkg/logger"

	"github.com/ethereum/go-ethereum/common"
	"github.com/perlin-network/life/exec"
)

var defaultLogger = logger.GetDefaultLogger("*:debug")

type Process interface {
	Sender() common.Address
	Args() Args
	State() db.StateDB
	SetResponse([]byte)
	Call(addr common.Address, entry []byte, args Args) ([]byte, error)
	Logger() logger.Logger
}

type process struct {
	vm *exec.VirtualMachine
	*Env
	rwm    db.RWSetMap
	logger logger.Logger
}

func NewProcess(vm *exec.VirtualMachine, env *Env) Process {
	return &process{
		vm:  vm,
		Env: env,
	}
}

func (p process) Sender() common.Address {
	return p.Env.Sender
}

func (p process) Args() Args {
	return p.Env.Args
}

func (p process) State() db.StateDB {
	return p.DB
}

func (p *process) Call(addr common.Address, entry []byte, args Args) ([]byte, error) {
	env, err := p.EnvManager.Get(p.Env.Context, p.Env.Sender, addr, args)
	if err != nil {
		return nil, err
	}
	res, err := env.Exec(p.Env.Context, string(entry))
	if err != nil {
		return nil, err
	}
	p.state.Add(res.RWSets)
	return res.Response, nil
}

func (p *process) Logger() logger.Logger {
	if p.logger != nil {
		return p.logger
	}
	return defaultLogger
}
