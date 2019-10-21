package main

import (
	"errors"
	"fmt"
	"github.com/bluele/hypermint/pkg/account/abi/bind"
	clibind "github.com/bluele/hypermint/pkg/account/cli/bind"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"io/ioutil"
	"os"
	"path"
	"strings"
)

var rootCmd = &cobra.Command{
	Use:   "cligen",
	Short: "Code generation from ABI",
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := viper.BindPFlags(cmd.Flags()); err != nil {
			return err
		}
		pkg := viper.GetString("package")
		outdir := viper.GetString("outdir")
		if len(outdir) == 0 {
			outdir = pkg
		}
		name := viper.GetString("name")
		if len(name) == 0 {
			return errors.New("name not specified")
		}
		abi := viper.GetString("abi")
		return Generate(outdir, pkg, name, abi)
	},
}

func Generate(outdir, pkg, name, abiJsonFilename string) error {
	if abiJson, err := ioutil.ReadFile(abiJsonFilename); err != nil {
		return err
	} else if boundSrc, err := bind.Bind(pkg, name, string(abiJson), false); err != nil {
		return err
	} else if src, err := clibind.Bind(pkg, name, string(abiJson)); err != nil {
		return err
	} else {
		if err := os.MkdirAll(outdir, 0700); err != nil {
			return err
		}
		if err := ioutil.WriteFile(path.Join(outdir, name+".go"), []byte(boundSrc), 0700); err != nil {
			return err
		}
		if err := ioutil.WriteFile(path.Join(outdir, name+"cmd.go"), []byte(src), 0700); err != nil {
			return err
		}
		return nil
	}
}


func init() {
	rootCmd.Flags().String("package", "contract", "package name")
	rootCmd.Flags().String("name", "", "contract name")
	rootCmd.Flags().String("abi", "example.json", "abi json")
	rootCmd.Flags().String("outdir", "", "output dir")
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
