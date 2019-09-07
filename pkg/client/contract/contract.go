package contract

import (
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"strconv"
	"strings"

	"github.com/ethereum/go-ethereum/common"
)

// type of return value
const (
	Int     = "int"
	Int32   = "int32"
	Int64   = "int64"
	UInt    = "uint"
	UInt32  = "uint32"
	UInt64  = "uint64"
	Bytes   = "bytes"
	Str     = "str"
	Address = "address"
)

func SerializeCallArgs(args []string, types []string) ([][]byte, error) {
	if len(args) != len(types) {
		return nil, fmt.Errorf("the number of arguments does not match the number of types: %v != %v", len(args), len(types))
	}

	var bs [][]byte
	for i, arg := range args {
		switch types[i] {
		case Int32, UInt32:
			v, err := strconv.Atoi(arg)
			if err != nil {
				return nil, err
			}
			var b [4]byte
			binary.BigEndian.PutUint32(b[:], uint32(v))
			bs = append(bs, b[:])
		case Int64, UInt64:
			v, err := strconv.Atoi(arg)
			if err != nil {
				return nil, err
			}
			var b [8]byte
			binary.BigEndian.PutUint64(b[:], uint64(v))
			bs = append(bs, b[:])
		case Bytes:
			if strings.HasPrefix(arg, "0x") {
				a, err := hex.DecodeString(arg[2:])
				if err != nil {
					return nil, err
				}
				bs = append(bs, a)
			} else {
				bs = append(bs, []byte(arg))
			}
		case Str:
			bs = append(bs, []byte(arg))
		case Address:
			ha := strings.TrimPrefix(arg, "0x")
			if l := len(ha); l != 20*2 {
				return nil, fmt.Errorf("address: invalid length %v", l)
			}
			a, err := hex.DecodeString(ha)
			if err != nil {
				return nil, err
			}
			bs = append(bs, a)
		default:
			return nil, fmt.Errorf("unknown types: %v", types[i])
		}
	}
	return bs, nil
}

func DeserializeValue(b []byte, tp string) (interface{}, error) {
	switch tp {
	case Int, Int32, Int64:
		if l := len(b); l == 4 {
			return int32(binary.BigEndian.Uint32(b)), nil
		} else if l == 8 {
			return int64(binary.BigEndian.Uint64(b)), nil
		} else {
			return nil, fmt.Errorf("unexpected bytes: %x", b)
		}
	case UInt, UInt32, UInt64:
		if l := len(b); l == 4 {
			return binary.BigEndian.Uint32(b), nil
		} else if l == 8 {
			return binary.BigEndian.Uint64(b), nil
		} else {
			return nil, fmt.Errorf("unexpected bytes: %x", b)
		}
	case Bytes:
		return b, nil
	case Str:
		return string(b), nil
	case Address:
		if l := len(b); l == 20 {
			var addr common.Address
			copy(addr[:], b)
			return addr, nil
		} else {
			return nil, fmt.Errorf("unexpected bytes: %x", b)
		}
	default:
		return nil, fmt.Errorf("unknown type: %v", tp)
	}
}
