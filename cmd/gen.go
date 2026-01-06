package cmd

import (
	"github.com/Salvadego/lathe/pkg/completions"
	"github.com/spf13/cobra"
)

var genCmd = &cobra.Command{
	Use:           "gen",
	SilenceUsage:  true,
	SilenceErrors: true,
	Short:         "Generate build artifacts",
}

func init() {
	genCmd.PersistentFlags().StringVarP(
		&buildMode,
		"mode", "m",
		"debug",
		"Build mode (debug|release)",
	)

	genCmd.PersistentFlags().StringVarP(
		&buildTarget,
		"target", "t",
		"",
		"Target to build",
	)

	genCmd.RegisterFlagCompletionFunc("target", completions.CompletionFuncTarget)
	genCmd.RegisterFlagCompletionFunc("mode", completions.CompletionFuncMode)

	rootCmd.AddCommand(genCmd)
}
