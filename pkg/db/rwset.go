package db

import (
	"encoding/binary"
	"errors"

	"github.com/bluele/hypermint/pkg/util"
	"github.com/ethereum/go-ethereum/common"
)

type Read struct {
	Key     []byte
	Version Version
}

type Write struct {
	Key   []byte
	Value []byte
}

type RWSet struct {
	ReadSet  []Read
	WriteSet []Write
}

type RWSets struct {
	Address common.Address
	RWSet   *RWSet
	Childs  []*RWSets
}

func (rs RWSets) Hash() []byte {
	b, err := rs.Bytes()
	if err != nil {
		panic(err)
	}
	return util.TxHash(b)
}

func (rs RWSets) Bytes() ([]byte, error) {
	return cdc.MarshalBinaryBare(rs)
}

func (rs *RWSets) FromBytes(b []byte) error {
	return cdc.UnmarshalBinaryBare(b, rs)
}

type State struct {
	Childs []*RWSets
}

func (s *State) Add(rws *RWSets) {
	s.Childs = append(s.Childs, rws)
}

type ValueObject struct {
	Value   []byte
	Version Version
}

func BytesToValueObject(b []byte) (*ValueObject, error) {
	vo := new(ValueObject)
	return vo, vo.Unmarshal(b)
}

func (vo *ValueObject) Marshal() []byte {
	buf := make([]byte, len(vo.Value)+4+4)
	copy(buf[:], vo.Value)
	binary.BigEndian.PutUint32(buf[len(vo.Value):], vo.Version.Height)
	binary.BigEndian.PutUint32(buf[len(vo.Value)+4:], vo.Version.TxIdx)
	return buf
}

func (vo *ValueObject) Unmarshal(b []byte) error {
	if len(b) < 8 {
		return errors.New("length of value is too short")
	}
	vo.Value = b[:len(b)-8]
	vo.Version.Height = binary.BigEndian.Uint32(b[len(b)-8:])
	vo.Version.TxIdx = binary.BigEndian.Uint32(b[len(b)-4:])
	return nil
}

type RWSetMap struct {
	rmap map[string]int
	rs   []Read
	wmap map[string]int
	ws   []Write
}

func NewRWSetMap() *RWSetMap {
	return &RWSetMap{
		rmap: make(map[string]int),
		wmap: make(map[string]int),
	}
}

func (m *RWSetMap) AddRead(key []byte, version Version) bool {
	s := string(key)
	if _, ok := m.rmap[s]; ok {
		return false
	}
	r := Read{Key: key, Version: version}
	m.rs = append(m.rs, r)
	m.rmap[s] = len(m.rs) - 1
	return true
}

func (m *RWSetMap) AddWrite(key, value []byte) {
	s := string(key)
	w := Write{Key: key, Value: value}
	if idx, ok := m.wmap[s]; ok {
		m.ws[idx] = w
	} else {
		m.ws = append(m.ws, w)
		m.wmap[s] = len(m.ws) - 1
	}
}

func (m *RWSetMap) ToSet() *RWSet {
	return &RWSet{ReadSet: m.rs, WriteSet: m.ws}
}

func (m *RWSetMap) GetRead(key []byte) (Read, bool) {
	idx, ok := m.rmap[string(key)]
	if !ok {
		return Read{}, false
	}
	return m.rs[idx], true
}

func (m *RWSetMap) GetWrite(key []byte) (Write, bool) {
	idx, ok := m.wmap[string(key)]
	if !ok {
		return Write{}, false
	}
	return m.ws[idx], true
}
