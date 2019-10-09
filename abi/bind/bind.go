package bind

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/bluele/hypermint/abi"
	"github.com/iancoleman/strcase"
	"go/format"
	"strings"
	"text/template"
)

type tmplContract struct {
	Type      string
	Transacts []abi.Function
	Events    []abi.Event
	Structs   []abi.Struct
	InputABI  string
}

type tmplData struct {
	Package   string
	Contracts map[string]tmplContract
	Mock      bool
}

func Bind(packageName string, name, jsonABI string, mock bool) (string, error) {
	buffer := new(bytes.Buffer)

	abi, err := abi.Unmarshal([]byte(jsonABI))
	if err != nil {
		return "", err
	}

	var packBuf bytes.Buffer
	packEnc := json.NewEncoder(&packBuf)
	err = packEnc.Encode(abi)
	if err != nil {
		return "", err
	}
	inputABI := strings.TrimRight(packBuf.String(), "\n")

	contracts := make(map[string]tmplContract)
	c := tmplContract{
		Type:      strcase.ToCamel(name),
		Transacts: abi.Functions,
		Events:    abi.Events,
		Structs:   abi.Structs,
		InputABI:  strings.Replace(inputABI, "\"", "\\\"", -1),
	}

	contracts[name] = c

	data := &tmplData{
		Package:   packageName,
		Contracts: contracts,
		Mock:      mock,
	}

	tmpl := template.Must(template.New("").Parse(tmplGoSource))
	if err := tmpl.Execute(buffer, data); err != nil {
		return "", err
	}

	code, err := format.Source(buffer.Bytes())
	if err != nil {
		return "", fmt.Errorf("%v\n%s", err, buffer)
	}
	return string(code), nil
}
