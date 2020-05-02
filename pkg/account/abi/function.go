package abi

import (
	"fmt"
	"strings"
)

type Function struct {
	RawName string `json:"name"`
	Name string `json:"-"`
	Type string `json:"type"`
	Simulate bool `json:"simulate"`
	Inputs Arguments `json:"inputs"`
	Outputs Arguments `json:"outputs"`
}

func (f Function) GetName() string {
	return f.Name
}

// TODO
func (f Function) Sig() string {
	types := make([]string, len(f.Inputs))
	for i, input := range f.Inputs {
		types[i] = input.Type.String()
	}
	return fmt.Sprintf("%v(%v)", f.RawName, strings.Join(types, ","))
}

// TODO
func (f Function) String() string {
	inputs := make([]string, len(f.Inputs))
	for i, input := range f.Inputs {
		inputs[i] = fmt.Sprintf("%v %v", input.Type, input.Name)
	}
	outputs := make([]string, len(f.Outputs))
	for i, output := range f.Outputs {
		outputs[i] = output.Type.String()
		if len(output.Name) > 0 {
			outputs[i] += fmt.Sprintf(" %v", output.Name)
		}
	}
	return fmt.Sprintf("function %v(%v) returns(%v)", f.RawName, strings.Join(inputs, ", "), strings.Join(outputs, ", "))
}
