package abi

type Struct struct {
	RawName string `json:"name"`
	Name string `json:"-"`
	Type string `json:"type"`
	Encoding string `json:"encoding"`
	Inputs Arguments `json:"inputs"`
}

func (s Struct) GetName() string {
	return s.Name
}
