package bind

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/bluele/hypermint/pkg/account/abi"
	"github.com/iancoleman/strcase"
	"go/format"
	"strings"
	"text/template"
)

type tmplFunction struct {
	abi.Function
	Use string
	Short string
}

type tmplContract struct {
	Type      string
	Name      string
	Use       string
	Short     string
	Functions []tmplFunction
	Events    []abi.Event
	Structs   []abi.Struct
	InputABI  string
}

type tmplData struct {
	Package   string
	Contracts map[string]tmplContract
	Endpoint  string
}

func Bind(packageName string, name, jsonABI string) (string, error) {
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

	var functions []tmplFunction
	for _, f := range abi.Functions {
		functions = append(functions, tmplFunction{
			Function: f,
			Use: strcase.ToDelimited(f.Name, '-'),
			Short: strcase.ToDelimited(f.Name, ' '),
		})
	}

	contracts := make(map[string]tmplContract)
	c := tmplContract{
		Type:      strcase.ToCamel(name),
		Name:      name,
		Use:       name,
		Short:     name,
		Functions: functions,
		Events:    abi.Events,
		Structs:   abi.Structs,
		InputABI:  strings.Replace(inputABI, "\"", "\\\"", -1),
	}

	contracts[name] = c

	data := &tmplData{
		Package:   packageName,
		Contracts: contracts,
		Endpoint:  "tcp://localhost:26657",
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
