package cmd

import (
	"os"
	"os/exec"

	"github.com/spf13/cobra"

	"github.com/Salvadego/lathe/internal/build"
	"github.com/Salvadego/lathe/internal/types"
	"github.com/Salvadego/lathe/pkg/completions"
)

var (
	runMode    string
	runTarget  string
	runWorkDir string
	runVerbose bool
)

var runCmd = &cobra.Command{
	Use:           "run [sources...] -- [program args]",
	Short:         "Build and run a target",
	Args:          cobra.ArbitraryArgs,
	SilenceUsage:  true,
	SilenceErrors: true,
	RunE: func(cmd *cobra.Command, args []string) error {

		dash := cmd.ArgsLenAtDash()

		var srcs []string
		var progArgs []string

		if dash >= 0 {
			srcs = args[:dash]
			progArgs = args[dash:]
		} else {
			srcs = args
		}

		req := types.BuildRequest{
			Mode:        types.Mode(runMode),
			Sources:     srcs,
			TargetName:  runTarget,
			Verbose:     runVerbose,
			WorkDir:     runWorkDir,
			Incremental: incremental,
		}

		ctx, err := build.ResolveBuildContext(req, cfg)
		if err != nil {
			return err
		}

		output, err := build.Build(ctx)
		if err != nil {
			return err
		}

		execCmd := exec.Command(output, progArgs...)
		execCmd.Stdin = os.Stdin
		execCmd.Stdout = os.Stdout
		execCmd.Stderr = os.Stderr

		return execCmd.Run()
	},
}

func init() {
	runCmd.Flags().StringVarP(
		&runMode,
		"mode", "m",
		"debug",
		"Build mode",
	)

	runCmd.Flags().StringVarP(
		&runTarget,
		"target", "t",
		"",
		"Target to run",
	)

	runCmd.Flags().StringVar(
		&runWorkDir,
		"workdir",
		".lathe",
		"Build output directory",
	)

	runCmd.Flags().BoolVarP(
		&runVerbose,
		"verbose", "v",
		false,
		"Verbose build output",
	)

	runCmd.Flags().BoolVar(
		&incremental,
		"incremental",
		false,
		"Enable incremental builds",
	)

	runCmd.RegisterFlagCompletionFunc("target", completions.CompletionFuncTarget)
	runCmd.RegisterFlagCompletionFunc("mode", completions.CompletionFuncMode)

	rootCmd.AddCommand(runCmd)
}
