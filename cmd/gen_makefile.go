package cmd

import (
	"os"

	"github.com/spf13/cobra"

	"github.com/Salvadego/lathe/internal/build"
	"github.com/Salvadego/lathe/internal/types"
)

var genMakefileCmd = &cobra.Command{
	Use:           "makefile",
	Short:         "Generate Makefile",
	SilenceUsage:  true,
	SilenceErrors: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		mode := types.Mode(buildMode)
		mk, err := build.Makefile(cfg, mode)
		if err != nil {
			return err
		}
		return os.WriteFile("Makefile", []byte(mk), 0644)
	},
}

func init() {
	genCmd.AddCommand(genMakefileCmd)
}
