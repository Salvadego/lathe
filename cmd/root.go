package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/Salvadego/lathe/internal/config"
	"github.com/Salvadego/lathe/pkg/defaults"
)

var (
	cfg config.Config
)

var rootCmd = &cobra.Command{
	Use:   "lathe",
	Short: "Lathe is a simple C build tool",
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		cfg = config.ResolveCfg(defaults.Config())
		return config.Validate(cfg)
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
