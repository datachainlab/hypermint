package cmd

import (
	"fmt"

	"github.com/bluele/hypermint/pkg/consts"
	"github.com/spf13/cobra"
)

var (
	// VersionCmd prints out the current sdk version
	versionCmd = &cobra.Command{
		Use:   "version",
		Short: "Print the app version",
		Run:   printVersion,
	}
)

// return version of CLI/node and commit hash
func GetVersion() string {
	v := consts.Version
	if consts.GitCommit != "" {
		v = v + "-" + consts.GitCommit
	}
	return v
}

// CMD
func printVersion(cmd *cobra.Command, args []string) {
	v := GetVersion()
	fmt.Println(v)
}
