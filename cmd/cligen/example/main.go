package main

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
	"strings"
)

var rootCmd = &cobra.Command{
	Use:   "example",
	Short: "Example CLI",
}

const (
	TokenContractAddress = "0x3b74892E655c9E2e47f3442Ef64dA6766fd2a62c"
	TokenCallerAddress = "0x933062aC91b12b27A88F0d8F00BE9eD0513f0D14"
)

func init() {
	rootCmd.AddCommand(TokenCmd(TokenContractAddress, TokenCallerAddress))
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
