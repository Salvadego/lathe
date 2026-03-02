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
		req := types.BuildRequest{
			Mode:       types.Mode(buildMode),
			TargetName: buildTarget,
			WorkDir:    workDir,
		}

		ctx, err := build.ResolveBuildContext(req, cfg)
		if err != nil {
			return err
		}

		mk := build.Makefile(ctx, cfg)
		return os.WriteFile("Makefile", []byte(mk), 0644)
	},
}

func init() {
	genCmd.AddCommand(genMakefileCmd)
}
