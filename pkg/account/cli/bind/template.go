package bind

var tmplGoSource = `
// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package {{.Package}}

import (
	"errors"
	"os"
	"fmt"
	"math/big"
	"strings"
	"encoding/binary"
	"bytes"
	"encoding/json"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/bluele/hypermint/pkg/account/abi/bind"
	clibind "github.com/bluele/hypermint/pkg/account/cli/bind"
	"github.com/bluele/hypermint/pkg/hmclient"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/accounts"
	"github.com/bluele/hypermint/pkg/transaction"
	"github.com/ethereum/go-ethereum/accounts/keystore"
)

// Reference imports to suppress errors if they are not otherwise used.
var (
	_ = big.NewInt
	_ = strings.NewReader
	_ = bind.Bind
	_ = clibind.Bind
	_ = common.Big1
	_ = transaction.ContractInitFunc
	_ = bytes.NewBuffer
	_ = binary.Read
	_ = json.NewEncoder
)

{{$endpoint := .Endpoint}}
{{range $contract := .Contracts -}}
	func {{.Type}}Cmd(contractAddress, callerAddress string) *cobra.Command {
		var cmd = &cobra.Command{
			Use:   "{{.Use}}",
			Short: "{{.Short}}",
		}
		cmd.PersistentFlags().Bool("verbose", false, "verbose")
		cmd.PersistentFlags().String("passphrase", "xxx", "passphrase")
		cmd.PersistentFlags().String("caller", callerAddress, "caller")
		cmd.PersistentFlags().String("contract", contractAddress, "contract")
		cmd.PersistentFlags().String("endpoint", "{{$endpoint}}", "endpoint")
		cmd.PersistentFlags().String("ksdir", "keystore", "keystore directory")
		{{range .Functions -}}
			cmd.AddCommand({{$contract.Name}}{{.Name}}Cmd)
		{{end}}
		return cmd
	}

	{{$events := .Events}}
	{{range .Functions}}
		func init() {
		    {{if .Inputs}}
                cmd := {{$contract.Name}}{{.Name}}Cmd
                {{range .Inputs -}}
                    {{if eq .Type.Name "bool"}}
                        cmd.Flags().Bool("{{.Name}}", false, "{{.Name}}")
                    {{else}}
                        cmd.Flags().String("{{.Name}}", "", "{{.Name}}")
                    {{end}}
                {{end -}}
			{{end -}}
		}

		var {{$contract.Name}}{{.Name}}Cmd = &cobra.Command{
			Use:   "{{.Use}}",
			Short: "{{.Short}}",
			RunE: func(cmd *cobra.Command, args []string) error {
				if cmd.OutOrStderr() == os.Stderr {
					cmd.SetOut(os.Stdout)
				}
				if err := viper.BindPFlags(cmd.Flags()); err != nil {
					return err
				}
				verbose := viper.GetBool("verbose")

				{{range .Inputs}}
					{{if eq .Type.Name "address"}}
						if len(viper.GetString("{{.Name}}")) == 0 {
							return errors.New("invalid address")
						}
						arg{{.Name}} := common.HexToAddress(viper.GetString("{{.Name}}"))
					{{else if eq .Type.Name "hash"}}
						if len(viper.GetString("{{.Name}}")) == 0 {
							return errors.New("invalid hash")
						}
						arg{{.Name}} := common.HexToHash(viper.GetString("{{.Name}}"))
					{{else if eq .Type.Name "bool"}}
						arg{{.Name}} := viper.GetBool("{{.Name}}")
					{{else}}
						arg{{.Name}} := {{.Type.TypeName}}(viper.GetInt64("{{.Name}}"))
					{{end}}
					if verbose {
						fmt.Fprintf(os.Stderr, "{{.Name}}=%v\n", arg{{.Name}})
					}
				{{end}}
				if verbose {
					fmt.Fprintf(os.Stderr, "passphrase=%v\n", viper.GetString("passphrase"))
				}

				ks := clibind.NewKeyStore(viper.GetString("ksdir"))
				opts, c, err := {{$contract.Name}}ContractFromFlags(ks)
				if err != nil {
					return err
				}
				{{if .Simulate}}
					if verbose {
						fmt.Fprintf(os.Stderr, "simulating {{.Name}}...")
					}
					{{range $i, $o := .Outputs}}v{{$i}}{{end}}, err := c.{{.Name}}(
						opts,
						{{range .Inputs -}}
							arg{{.Name}},
						{{end}}
					)
					if err != nil {
						return err
					}
					if verbose {
						fmt.Fprintf(os.Stderr, "done\n")
					}
					{{range $i, $o := .Outputs}}
						{{if eq .Type.Name "address"}}
							cmd.Println(v{{$i}}.Hex())
						{{else}}
							cmd.Println(v{{$i}})
						{{end}}
					{{end}}
					return nil
				{{else}}
					if verbose {
						fmt.Fprintf(os.Stderr, "committing {{.Name}}...")
					}
					r, err := c.{{.Name}}Commit(
						opts,
						{{range .Inputs -}}
							arg{{.Name}},
						{{end}}
					)
					if err != nil {
						return err
					}
					if verbose {
						fmt.Fprintf(os.Stderr, "%v\n", r.TxHash.Hex())
					}
					ed := {{$contract.Type}}EventDecoder
					if ed == nil {
						return errors.New("{{$contract.Type}}EventDecoder is nil")
					}
					{{range $events}}
						if verbose {
							fmt.Fprintf(os.Stderr, "looking up {{$contract.Type}}{{.Name}}...")
						}
						var _{{.Name}} {{$contract.Type}}{{.Name}}
						if err := ed.FindFirst(&_{{.Name}}, r.Entries); err != nil {
							if verbose {
								fmt.Fprintf(os.Stderr, "%v\n", err.Error())
							}
						} else {
							if verbose {
								fmt.Fprintf(os.Stderr, "found\n")
							}
							json.NewEncoder(os.Stdout).Encode(&_{{.Name}})
						}
					{{end}}
					return nil
				{{end}}
			},
		}
	{{end}}

	func {{.Name}}ContractFromFlags(ks *keystore.KeyStore) (*bind.TransactOpts, {{.Type}}Contract, error) {
		pp, err := clibind.GetFlagPassphrase(func() (string, error) {
			pp := viper.GetString("passphrase")
			if pp == "xxx" {
				return "", errors.New("passphrase required")
			}
			return pp, nil
		})()
		if err != nil {
			return nil, nil, err
		}

		caller := common.HexToAddress(viper.GetString("caller"))
		if err := clibind.Unlock(ks, caller, pp); err != nil {
			return nil, nil, err
		}
	
		contractAddress, err := clibind.GetFlagAddress(func() (string, error) {
			return viper.GetString("contract"), nil
		}, nil)()
		if err != nil {
			return nil, nil, err
		}
	
		opts, err := bind.NewKeyStoreTransactor(ks, accounts.Account{Address:caller})
		if err != nil {
			return nil, nil, err
		}
	
		endpoint := viper.GetString("endpoint")
		cl := hmclient.NewClient(endpoint)
		c{{.Type}}, err := New{{.Type}}(contractAddress, cl)
		if err != nil {
			return nil, nil, err
		}
	
		return opts, c{{.Type}}, nil
	}

{{end}}
`
