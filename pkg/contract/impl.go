package contract

import (
	"github.com/bluele/hypermint/pkg/util"
	"github.com/ethereum/go-ethereum/common"
)

type Writer interface {
	Write([]byte) int
}

type Reader interface {
	Read() []byte
}

type value struct {
	mem []byte
	pos int64
	len int64
}

func (v *value) Write(b []byte) int {
	if len(b) > int(v.len) {
		return -1
	}
	copy(v.mem[v.pos:], b)
	return len(b)
}

func (v value) Read() []byte {
	b := make([]byte, v.len)
	copy(b, v.mem[v.pos:v.pos+v.len])
	return b
}

func NewWriter(mem []byte, pos, len int64) Writer {
	return &value{mem: mem, pos: pos, len: len}
}

func NewReader(mem []byte, pos, len int64) Reader {
	return &value{mem: mem, pos: pos, len: len}
}

func GetArg(ps Process, idx int, w Writer) int {
	arg := ps.Args().Get(idx)
	return w.Write([]byte(arg))
}

func Log(ps Process, msg Reader) int {
	ps.Logger().Debug(string(msg.Read()))
	return 0
}

func GetSender(ps Process, w Writer) int {
	s := ps.Sender()
	return w.Write(s[:])
}

func ReadState(ps Process, key Reader, buf Writer) int {
	v, err := ps.State().Get(key.Read())
	if err != nil {
		ps.Logger().Debug("fail to execute ReadState", "err", err)
		return -1
	}
	return buf.Write(v)
}

func WriteState(ps Process, key, val Reader) int {
	err := ps.State().Set(key.Read(), val.Read())
	if err != nil {
		ps.Logger().Debug("fail to execute WriteState", "err", err)
		return -1
	}
	return 0
}

func SetResponse(ps Process, val Reader) int {
	ps.SetResponse(val.Read())
	return 0
}

func CallContract(ps Process, addr, entry Reader, args Args, ret Writer) int {
	res, err := ps.Call(common.BytesToAddress(addr.Read()), entry.Read(), args)
	if err != nil {
		ps.Logger().Debug("fail to execute CallContract", "err", err)
		return -1
	}
	return ret.Write(res)
}

func ECRecover(ps Process, h, v, r, s Reader, ret Writer) int {
	pub, err := util.Ecrecover(h.Read(), v.Read(), r.Read(), s.Read())
	if err != nil {
		ps.Logger().Debug("fail to execute ECRecover", "err", err)
		return -1
	}
	return ret.Write(pub)
}

func ECRecoverAddress(ps Process, h, v, r, s Reader, ret Writer) int {
	addr, err := util.EcrecoverAddress(h.Read(), v.Read(), r.Read(), s.Read())
	if err != nil {
		ps.Logger().Debug("fail to execute ECRecoverAddress", "err", err)
		return -1
	}
	return ret.Write(addr[:])
}
