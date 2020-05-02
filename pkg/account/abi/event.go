package abi

type Event struct {
	RawName string `json:"name"`
	Name string `json:"-"`
	Type string `json:"type"`
	Encoding string `json:"encoding"`
	Inputs Arguments `json:"inputs"`
}

func (e Event) GetName() string {
	return e.Name
}
