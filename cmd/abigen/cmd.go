package main

import (
	"errors"
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"io/ioutil"
	"os"
	"path"
	"strings"
	"github.com/bluele/hypermint/abi/bind"
)

var rootCmd = &cobra.Command{
	Use:   "abigen",
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
		mock := viper.GetBool("mock")
		return Generate(outdir, pkg, name, abi, mock)
	},
}

func Generate(outdir, pkg, name, abiJsonFilename string, mock bool) error {
	if abiJson, err := ioutil.ReadFile(abiJsonFilename); err != nil {
		return err
	} else if src, err := bind.Bind(pkg, name, string(abiJson), mock); err != nil {
		return err
	} else {
		if err := os.MkdirAll(pkg, 0700); err != nil {
			return err
		}
		if err := ioutil.WriteFile(path.Join(outdir, name+".go"), []byte(src), 0700); err != nil {
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
	rootCmd.Flags().Bool("mock", false, "generate mock")
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}