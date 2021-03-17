package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

const version = "v0.4.2"

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of Golic",
	Long:  "",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(version)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}