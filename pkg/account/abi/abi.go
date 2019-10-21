package abi

import (
	"encoding/json"
	"fmt"
	"github.com/iancoleman/strcase"
	"io/ioutil"
)

type ABI struct {
	Functions []Function `json:"functions"`
	Events    []Event    `json:"events"`
	Structs   []Struct   `json:"structs"`
}

func (e *ABI) UnmarshalJSON(data []byte) error {
	var fields []struct {
		Type     string    `json:"type"`
		Name     string    `json:"name"`
		Simulate bool      `json:"simulate"`
		Encoding string    `json:"encoding"`
		Inputs   Arguments `json:"inputs"`
		Outputs  Arguments `json:"outputs"`
	}
	if err := json.Unmarshal(data, &fields); err != nil {
		return err
	}
	for _, field := range fields {
		switch field.Type {
		case "function":
			e.Functions = append(e.Functions, Function{
				Type:     field.Type,
				Name:     strcase.ToCamel(field.Name),
				RawName:  field.Name,
				Simulate: field.Simulate,
				Inputs:   field.Inputs,
				Outputs:  field.Outputs,
			})
		case "event":
			e.Events = append(e.Events, Event{
				Type:     field.Type,
				Name:     strcase.ToCamel(field.Name),
				RawName:  field.Name,
				Encoding: field.Encoding,
				Inputs:   field.Inputs,
			})
		case "struct":
			e.Structs = append(e.Structs, Struct{
				Type:     field.Type,
				Name:     strcase.ToCamel(field.Name),
				RawName:  field.Name,
				Encoding: field.Encoding,
				Inputs:   field.Inputs,
			})
		default:
			return fmt.Errorf("unknown type: %v", field.Type)
		}
	}
	return nil
}

func Unmarshal(data []byte) (*ABI, error) {
	var abi ABI
	err := json.Unmarshal(data, &abi)
	return &abi, err
}

func ReadJsonFile(filename string) (*ABI, error) {
	if bytes, err := ioutil.ReadFile(filename); err != nil {
		return nil, err
	} else {
		return Unmarshal(bytes)
	}
}
