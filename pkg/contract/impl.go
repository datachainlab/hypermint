package contract

import (
	"github.com/bluele/hypermint/pkg/util"
	"github.com/ethereum/go-ethereum/common"
)

type Writer interface {
	Len() int
	Write([]byte) int
}

type Reader interface {
	Len() int
	Read() []byte
}

type value struct {
	mem []byte
	pos int64
	len int64
}

func (v *value) Len() int {
	return int(v.len)
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

func ReadState(ps Process, key Reader, offset int, buf Writer) int {
	v, err := ps.State().Get(key.Read())
	if err != nil {
		ps.Logger().Debug("fail to execute ReadState", "err", err)
		return -1
	}
	if offset < 0 {
		ps.Logger().Debug("offset must be positive", "offset", offset)
		return -1
	} else if len(v) <= offset {
		ps.Logger().Debug("offset is over value length", "offset", offset, "length", len(v))
		return 0
	}
	return buf.Write(v[offset:min(offset+buf.Len(), len(v))])
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

func min(vs ...int) int {
	if len(vs) == 0 {
		panic("length of vs should be greater than 0")
	}
	min := vs[0]
	for _, v := range vs {
		if v < min {
			min = v
		}
	}
	return min
}