package contract

import (
	"github.com/spf13/cobra"
)

const (
	flagGas = "gas"
)

var contractCmd = &cobra.Command{
	Use:   "contract",
	Short: "wasm contract command",
}

func Setup(cmd *cobra.Command) {
	cmd.AddCommand(contractCmd)
}
