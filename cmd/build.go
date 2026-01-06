package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/Salvadego/lathe/internal/build"
	"github.com/Salvadego/lathe/internal/types"
	"github.com/Salvadego/lathe/pkg/completions"
)

var (
	buildMode   string
	buildTarget string
	workDir     string
	verbose     bool
)

var buildCmd = &cobra.Command{
	Use:           "build [sources...]",
	Short:         "Build a target or sources",
	SilenceUsage:  true,
	SilenceErrors: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		req := types.BuildRequest{
			Mode:       types.Mode(buildMode),
			Sources:    args,
			TargetName: buildTarget,
			Verbose:    verbose,
			WorkDir:    workDir,
		}

		ctx, err := build.ResolveBuildContext(req, cfg)
		if err != nil {
			return err
		}

		output, err := build.Build(ctx)
		if err != nil {
			return err
		}

		fmt.Println(output)
		return nil
	},
}

func init() {
	buildCmd.Flags().StringVarP(
		&buildMode,
		"mode", "m",
		"debug",
		"Build mode (debug|release)",
	)

	buildCmd.Flags().StringVarP(
		&buildTarget,
		"target", "t",
		"",
		"Target to build",
	)

	buildCmd.Flags().StringVar(
		&workDir,
		"workdir",
		".lathe",
		"Build output directory",
	)

	buildCmd.Flags().BoolVarP(
		&verbose,
		"verbose", "v",
		false,
		"Verbose build output",
	)

	buildCmd.RegisterFlagCompletionFunc("target", completions.CompletionFuncTarget)
	buildCmd.RegisterFlagCompletionFunc("mode", completions.CompletionFuncMode)
	rootCmd.AddCommand(buildCmd)
}
