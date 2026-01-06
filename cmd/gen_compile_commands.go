package cmd

import (
	"github.com/spf13/cobra"

	"github.com/Salvadego/lathe/internal/build"
	"github.com/Salvadego/lathe/internal/export"
	"github.com/Salvadego/lathe/internal/types"
)

var genCompileCmd = &cobra.Command{
	Use:           "compile-commands",
	Short:         "Generate compile_commands.json",
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

		cmds, err := build.CompileCommands(ctx)
		if err != nil {
			return err
		}

		return export.WriteCompileCommands("compile_commands.json", cmds)
	},
}

func init() {
	genCmd.AddCommand(genCompileCmd)
}
