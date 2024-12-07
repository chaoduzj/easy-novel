package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(versionCmd)
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of easy-novel",
	Long:  `All software has versions. This is easy-novel's`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("easy-novel v0.1 -- HEAD")
	},
}
