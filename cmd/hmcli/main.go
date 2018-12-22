package main

import (
	"github.com/bluele/hypermint/pkg/client/cmd"
	"github.com/spf13/cobra"
)

func main() {
	cobra.EnableCommandSorting = false
	cmd.Execute()
}
