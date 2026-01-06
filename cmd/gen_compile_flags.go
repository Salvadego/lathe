package cmd

import (
	"github.com/spf13/cobra"

	"github.com/Salvadego/lathe/internal/build"
	"github.com/Salvadego/lathe/internal/export"
	"github.com/Salvadego/lathe/internal/types"
)

var genCompileFlagsCmd = &cobra.Command{
	Use:           "compile-flags",
	Short:         "Generate compile_flags.txt",
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

		return export.WriteCompileFlags("compile_flags.txt", ctx)
	},
}

func init() {
	genCmd.AddCommand(genCompileFlagsCmd)
}
