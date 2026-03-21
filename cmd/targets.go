package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var targetsCmd = &cobra.Command{
	Use:   "targets",
	Short: "List out the targets from the current config",
	Run: func(cmd *cobra.Command, args []string) {
		for name := range cfg.Targets {
			fmt.Println(name)
		}
	},
}

func init() {
	rootCmd.AddCommand(targetsCmd)
}
