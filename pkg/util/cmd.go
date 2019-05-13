package util

import "github.com/spf13/cobra"

func CheckRequiredFlag(cmd *cobra.Command, names ...string) {
	for _, name := range names {
		if err := cmd.MarkFlagRequired(name); err != nil {
			panic(err)
		}
	}
}
