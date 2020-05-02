package bind

import (
	"encoding/binary"
	"encoding/json"
	"github.com/ethereum/go-ethereum/common"
)

type Bool bool

func (b Bool) Bytes() []byte {
	if b {
		return []byte{1}
	} else {
		return []byte{0}
	}
}

type I8 int8

func (i I8) Bytes() []byte {
	return []byte{byte(i)}
}

type I16 int16

func (i I16) Bytes() []byte {
	v := make([]byte, 2)
	binary.BigEndian.PutUint16(v, uint16(i))
	return v
}

type I32 int32

func (i I32) Bytes() []byte {
	v := make([]byte, 4)
	binary.BigEndian.PutUint32(v, uint32(i))
	return v
}

type I64 int64

func (i I64) Bytes() []byte {
	v := make([]byte, 8)
	binary.BigEndian.PutUint64(v, uint64(i))
	return v
}

type U8 uint8

func (u U8) Bytes() []byte {
	return []byte{byte(u)}
}

type U16 uint16

func (u U16) Bytes() []byte {
	v := make([]byte, 2)
	binary.BigEndian.PutUint16(v, uint16(u))
	return v
}

type U32 uint32

func (u U32) Bytes() []byte {
	v := make([]byte, 4)
	binary.BigEndian.PutUint32(v, uint32(u))
	return v
}

type U64 uint64

func (u U64) Bytes() []byte {
	v := make([]byte, 8)
	binary.BigEndian.PutUint64(v, uint64(u))
	return v
}

type String string

func (s String) Bytes() []byte {
	return []byte(s)
}

type Bytes []byte

func (s Bytes) Bytes() []byte {
	return s
}

type Address common.Address

func (a Address) Bytes() []byte {
	return common.Address(a).Bytes()
}

func (a *Address) MarshalJSON() ([]byte, error) {
	ua := []uint8(a.Bytes())
	return json.Marshal(ua)
}

func (a *Address) UnmarshalJSON(bs []byte) error {
	var ua [20]uint8
	if err := json.Unmarshal(bs, &ua); err != nil {
		return err
	}
	*a = Address(common.Address(ua))
	return nil
}

type Hash common.Hash

func (a Hash) Bytes() []byte {
	return common.Hash(a).Bytes()
}
