package bind

const tmplGoSource = `
// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package {{.Package}}

import (
	"encoding/binary"
	"bytes"
	"encoding/json"

	"github.com/bluele/hypermint/abi/bind"
	"github.com/bluele/hypermint/pkg/transaction"
	"github.com/ethereum/go-ethereum/common"
{{if .Mock}}
	"github.com/stretchr/testify/mock"
{{end}}
)

// Reference imports to suppress errors if they are not otherwise used.
var (
	_ = bind.Bind
	_ = common.Big1
	_ = transaction.ContractInitFunc
	_ = bytes.NewBuffer
	_ = binary.Read
	_ = json.NewEncoder
)

{{range $contract := .Contracts}}
	var {{.Type}}ABI string = "{{ .InputABI }}"

	{{range .Structs}}
		{{if eq .Encoding "json"}}
			type {{$contract.Type}}{{.Name}}Raw struct {
				{{range .Inputs -}}
					{{.PubName}} {{.Type.BoundType}} `+"`"+`json:"{{.Name}}"`+"`"+`
				{{end}}
			}
		{{else}}
			type {{$contract.Type}}{{.Name}}Raw struct {
				{{range .Inputs -}}
					{{.PubName}} {{.Type.TypeName}}
				{{end}}
			}
		{{end}}

		type {{$contract.Type}}{{.Name}} struct {
			{{range .Inputs -}}
				{{.PubName}} {{.Type.TypeName}} `+"`"+`json:"{{.Name}}"`+"`"+`
			{{end}}
		}
		
		{{if eq .Encoding "json"}}
			func (_{{.Name}} *{{$contract.Type}}{{.Name}}) Decode(bs []byte) error {
				var raw {{$contract.Type}}{{.Name}}Raw
				if err := json.Unmarshal(bs, &raw); err != nil {
					return err
				}
				if err := bind.DeepCopy(_{{.Name}}, &raw); err != nil {
					return err
				}
				return nil
			}
		{{end}}
	{{end}}

	{{range .Events}}
		{{if eq .Encoding "json"}}
			type {{$contract.Type}}{{.Name}}Raw struct {
				{{range .Inputs -}}
					{{.PubName}} {{.Type.BoundType}} `+"`"+`json:"{{.Name}}"`+"`"+`
				{{end}}
			}
		{{else}}
			type {{$contract.Type}}{{.Name}}Raw struct {
				{{range .Inputs -}}
					{{.PubName}} {{.Type.TypeName}}
				{{end}}
			}
		{{end}}

		type {{$contract.Type}}{{.Name}} struct {
			{{range .Inputs -}}
				{{.PubName}} {{.Type.TypeName}} `+"`"+`json:"{{.Name}}"`+"`"+`
			{{end}}
		}
		
		{{if eq .Encoding "json"}}
			func (_{{.Name}} *{{$contract.Type}}{{.Name}}) Decode(bs []byte) error {
				var raw {{$contract.Type}}{{.Name}}Raw
				if err := json.Unmarshal(bs, &raw); err != nil {
					return err
				}
				if err := bind.DeepCopy(_{{.Name}}, &raw); err != nil {
					return err
				}
				return nil
			}
		{{end}}
	{{end}}

	var {{.Type}}EventDecoder = bind.NewEventDecoder()

	{{range .Events}}
		var {{$contract.Type}}{{.Name}}Info = bind.EventInfo {
			ID: "{{.Name}}",
			EventCreator: func() bind.Event {
				return &{{$contract.Type}}{{.Name}}{}
			},
		}
	{{end}}

	func init() {
		{{range .Events -}}
			{{$contract.Type}}EventDecoder.Register({{$contract.Type}}{{.Name}}Info)
		{{end -}}
	}

	type {{.Type}}Contract interface {
		{{range .Transacts}}
			{{if .Simulate}}
				{{.Name}}(opts *bind.TransactOpts {{range .Inputs}}, {{.Name}} {{.Type.TypeName}}{{end}}) ({{range .Outputs}}{{.Type.TypeName}}, {{end}}error)
			{{else}}
				{{.Name}}(opts *bind.TransactOpts {{range .Inputs}}, {{.Name}} {{.Type.TypeName}}{{end}}) (*bind.SyncResult, error)
				{{.Name}}Commit(opts *bind.TransactOpts {{range .Inputs}}, {{.Name}} {{.Type.TypeName}}{{end}}) (*bind.CommitResult, error)
			{{end}}
		{{end}}
	}

	type {{.Type}} struct {
	  {{.Type}}Simulator  // Read-only binding to the contract
	  {{.Type}}Transactor // Write-only binding to the contract
	}

	type {{.Type}}Simulator struct {
		contract *bind.BoundContract
	}

	type {{.Type}}Transactor struct {
		contract *bind.BoundContract
	}

	type {{.Type}}Raw struct {
		Contract *{{.Type}}
	}

	type {{.Type}}SimulatorRaw struct {
		Contract *{{.Type}}Simulator
	}

	type {{.Type}}TransactorRaw struct {
		Contract *{{.Type}}Transactor
	}

	func New{{.Type}}(address common.Address, backend bind.ContractBackend) (*{{.Type}}, error) {
	  contract, err := bind{{.Type}}(address, backend, backend)
	  if err != nil {
	    return nil, err
	  }
	  return &{{.Type}}{ {{.Type}}Simulator: {{.Type}}Simulator{contract: contract}, {{.Type}}Transactor: {{.Type}}Transactor{contract: contract} }, nil
	}

	func New{{.Type}}Simulator(address common.Address, simulator bind.ContractSimulator) (*{{.Type}}Simulator, error) {
	  contract, err := bind{{.Type}}(address, simulator, nil)
	  if err != nil {
	    return nil, err
	  }
	  return &{{.Type}}Simulator{contract: contract}, nil
	}

	func New{{.Type}}Transactor(address common.Address, transactor bind.ContractTransactor) (*{{.Type}}Transactor, error) {
	  contract, err := bind{{.Type}}(address, nil, transactor)
	  if err != nil {
	    return nil, err
	  }
	  return &{{.Type}}Transactor{contract: contract}, nil
	}

	func bind{{.Type}}(address common.Address, simulator bind.ContractSimulator, transactor bind.ContractTransactor) (*bind.BoundContract, error) {
	  return bind.NewBoundContract(address, simulator, transactor), nil
	}

	func (_{{$contract.Type}} *{{$contract.Type}}Raw) Transact(opts *bind.TransactOpts, fn string, params ...[]byte) (*bind.SyncResult, error) {
		return _{{$contract.Type}}.Contract.{{$contract.Type}}Transactor.contract.Transact(opts, fn, params...)
	}

	func (_{{$contract.Type}} *{{$contract.Type}}SimulatorRaw) Simulate(opts *bind.TransactOpts, fn string, params ...[]byte) (*bind.SimulateResult, error) {
		return _{{$contract.Type}}.Contract.contract.Simulate(opts, fn, params...)
	}

	{{range .Transacts}}
		{{if .Simulate}}
			func (_{{$contract.Type}} *{{$contract.Type}}Simulator) {{.Name}}(opts *bind.TransactOpts {{range .Inputs}}, {{.Name}} {{.Type.TypeName}}{{end}}) ({{range .Outputs}}{{.Type.TypeName}}, {{end}}error) {
				{{if .Outputs}}result{{else}}_{{end}}, err := _{{$contract.Type}}.contract.Simulate(opts, "{{.RawName}}", bind.Args({{range .Inputs}}
					{{.Type.ToBoundType .Name}},{{end}})...)
				if err != nil {
					return {{range .Outputs}}{{.Type.Nil}}, {{end}}err
				}
				{{if .Outputs}}
				buf := bytes.NewBuffer(result.Data)
				{{end}}
				{{range $i, $v := .Outputs -}}
					{{if eq .Type.TypeName "string"}}
						var v{{ $i }} {{.Type.TypeName}}
						v{{ $i }} = buf.String()
					{{else}}
						var v{{ $i }} {{.Type.TypeName}}
						binary.Read(buf, binary.BigEndian, &v{{ $i }})
					{{end}}
				{{end -}}
				return {{range $i, $v := .Outputs}}v{{ $i }}, {{end}}nil
			}
		{{else}}
			func (_{{$contract.Type}} *{{$contract.Type}}Transactor) {{.Name}}(opts *bind.TransactOpts {{range .Inputs}}, {{.Name}} {{.Type.TypeName}}{{end}}) (*bind.SyncResult, error) {
				return _{{$contract.Type}}.contract.Transact(opts, "{{.RawName}}", bind.Args({{range .Inputs}}
					{{.Type.ToBoundType .Name}},{{end}})...)
			}
			
			func (_{{$contract.Type}} *{{$contract.Type}}Transactor) {{.Name}}Commit(opts *bind.TransactOpts {{range .Inputs}}, {{.Name}} {{.Type.TypeName}}{{end}}) (*bind.CommitResult, error) {
				return _{{$contract.Type}}.contract.TransactCommit(opts, "{{.RawName}}", bind.Args({{range .Inputs}}
					{{.Type.ToBoundType .Name}},{{end}})...)
			}
		{{end}}
	{{end}}
{{end}}

{{if .Mock}}
	{{range $contract := .Contracts}}
		type Mock{{.Type}} struct {
			mock.Mock
		}

		{{range .Transacts}}
			{{if .Simulate}}
				func (_{{$contract.Type}} *Mock{{$contract.Type}}) {{.Name}}(opts *bind.TransactOpts {{range .Inputs}}, {{.Name}} {{.Type.TypeName}}{{end}}) ({{range .Outputs}}{{.Type.TypeName}}, {{end}}error) {
					ret := _{{$contract.Type}}.Called({{range .Inputs}}{{.Name}}, {{end}})
					return {{range $i, $v := .Outputs}}ret.Get({{$i}}).({{.Type.TypeName}}), {{end}}ret.Error({{len .Outputs}})
				}
			{{else}}
				func (_{{$contract.Type}} *Mock{{$contract.Type}}) {{.Name}}(opts *bind.TransactOpts {{range .Inputs}}, {{.Name}} {{.Type.TypeName}}{{end}}) (*bind.SyncResult, error) {
					ret := _{{$contract.Type}}.Called({{range .Inputs}}{{.Name}}, {{end}})
					return ret.Get(0).(*bind.SyncResult), ret.Error(1)
				}
				
				func (_{{$contract.Type}} *Mock{{$contract.Type}}) {{.Name}}Commit(opts *bind.TransactOpts {{range .Inputs}}, {{.Name}} {{.Type.TypeName}}{{end}}) (*bind.CommitResult, error) {
					ret := _{{$contract.Type}}.Called({{range .Inputs}}{{.Name}}, {{end}})
					return ret.Get(0).(*bind.CommitResult), ret.Error(1)
				}
			{{end}}
		{{end}}
	{{end}}
{{end}}
`
