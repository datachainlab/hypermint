package abi

import (
	"encoding/json"
	"fmt"
	"github.com/iancoleman/strcase"
)

type Argument struct {
	Type Type
	Name string
	PubName string
}

type Arguments []Argument

func typeFromName(typeName string) (*Type, error) {
	var t byte
	switch typeName {
	case "i8":
		t = I8Ty
	case "i16":
		t = I16Ty
	case "i32":
		t = I32Ty
	case "i64":
		t = I64Ty
	case "u8":
		t = U8Ty
	case "u16":
		t = U16Ty
	case "u32":
		t = U32Ty
	case "u64":
		t = U64Ty
	case "bool":
		t = BoolTy
	case "bytes":
		t = BytesTy
	case "str":
		t = StringTy
	case "address":
		t = AddressTy
	case "hash":
		t = HashTy
	default:
		return nil, fmt.Errorf("unknown type: %v", typeName)
	}
	return &Type{
		Name: typeName,
		T:    t,
	}, nil
}

func (a *Argument) UnmarshalJSON(data []byte) error {
	var fields struct {
		Type string `json:"type"`
		Name string `json:"name"`
	}
	err := json.Unmarshal(data, &fields)
	if err != nil {
		return err
	}
	if ty, err := typeFromName(fields.Type); err != nil {
		return err
	} else {
		a.Type = *ty
	}
	a.Name = fields.Name
	if a.Name == "" {
		a.Name = "_"
	}
	a.PubName = strcase.ToCamel(fields.Name)
	return nil
}

func (a *Argument) MarshalJSON() ([]byte, error) {
	var fields struct {
		Type string `json:"type"`
		Name string `json:"name"`
	}
	fields.Type = a.Type.Name
	fields.Name = a.Name
	return json.Marshal(fields)
}
