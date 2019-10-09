package abi

import "fmt"

const (
	I8Ty byte = iota
	I16Ty
	I32Ty
	I64Ty
	U8Ty
	U16Ty
	U32Ty
	U64Ty
	BoolTy
	BytesTy
	StringTy
	AddressTy
	HashTy
)

type Type struct {
	Name string // ABI type name
	T byte      // type code
}

// as string
func (t Type) String() string {
	return t.Name
}

func (t Type) TypeName() string {
	switch t.T {
	case I8Ty:
		return "int8"
	case I16Ty:
		return "int16"
	case I32Ty:
		return "int32"
	case I64Ty:
		return "int64"
	case U8Ty:
		return "uint8"
	case U16Ty:
		return "uint16"
	case U32Ty:
		return "uint32"
	case U64Ty:
		return "uint64"
	case BoolTy:
		return "bool"
	case BytesTy:
		return "[]byte"
	case StringTy:
		return "string"
	case AddressTy:
		return "common.Address"
	case HashTy:
		return "common.Hash"
	default:
		panic(fmt.Errorf("unknown type: %v", t.T))
	}
}

func (t Type) Nil() string {
	switch t.T {
	case I8Ty:
		return "int8(0)"
	case I16Ty:
		return "int16(0)"
	case I32Ty:
		return "int32(0)"
	case I64Ty:
		return "int64(0)"
	case U8Ty:
		return "uint8(0)"
	case U16Ty:
		return "uint16(0)"
	case U32Ty:
		return "uint32(0)"
	case U64Ty:
		return "uint64(0)"
	case StringTy:
		return "\"\""
	case BytesTy:
		return "nil"
	case BoolTy:
		return "false"
	case AddressTy:
		return "common.Address{}"
	case HashTy:
		return "common.Hash{}"
	default:
		panic(fmt.Errorf("unknown type: %v", t.T))
	}
}

func (t Type) ToBoundType(s string) string {
	return t.BoundType() + "(" + s + ")"
}

func (t Type) FromBoundType(s string) string {
	return t.TypeName() + "(" + s + ")"
}

func (t Type) BoundType() string {
	var op string
	switch t.T {
	case I8Ty:
		op = "bind.I8"
	case I16Ty:
		op = "bind.I16"
	case I32Ty:
		op = "bind.I32"
	case I64Ty:
		op = "bind.I64"
	case U8Ty:
		op = "bind.U8"
	case U16Ty:
		op = "bind.U16"
	case U32Ty:
		op = "bind.U32"
	case U64Ty:
		op = "bind.U64"
	case BoolTy:
		op = "bind.Bool"
	case StringTy:
		op = "bind.String"
	case AddressTy:
		op = "bind.Address"
	case HashTy:
		op = "bind.Hash"
	default:
		panic(fmt.Errorf("unknown type: %v", t.T))
	}
	return op
}
