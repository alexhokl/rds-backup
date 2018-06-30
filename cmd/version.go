package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var version = "development"

func init() {
	var versionCmd = &cobra.Command{
		Use:   "version",
		Short: "Show version",
		Long:  "Show version",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("rds-backup version (%s)\n", version)
		},
	}
	RootCmd.AddCommand(versionCmd)
}
