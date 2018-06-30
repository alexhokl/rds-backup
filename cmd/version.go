package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var version = "development"
var tag = ""

func init() {
	var versionCmd = &cobra.Command{
		Use:   "version",
		Short: "Show version",
		Long:  "Show version",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("rds-backup version (%s, %s)\n", tag, version)
		},
	}
	RootCmd.AddCommand(versionCmd)
}
