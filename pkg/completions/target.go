package completions

import (
	"strings"

	"github.com/Salvadego/lathe/internal/config"
	"github.com/Salvadego/lathe/pkg/defaults"
	"github.com/spf13/cobra"
)

func CompletionFuncTarget(
	cmd *cobra.Command,
	args []string,
	toComplete string,
) ([]cobra.Completion, cobra.ShellCompDirective) {

	cfg := defaults.Config()
	cfg, _ = config.ResolveCfg(cfg)

	completions := make([]cobra.Completion, 0, len(cfg.Targets))

	for name := range cfg.Targets {
		if strings.HasPrefix(name, toComplete) {
			completions = append(completions, name)
		}
	}

	return completions, cobra.ShellCompDirectiveNoFileComp
}

func CompletionFuncMode(
	cmd *cobra.Command,
	args []string,
	toComplete string,
) ([]cobra.Completion, cobra.ShellCompDirective) {

	cfg := defaults.Config()
	cfg, _ = config.ResolveCfg(cfg)

	completions := make([]cobra.Completion, 0, len(cfg.Builds))

	for name := range cfg.Builds {
		if strings.HasPrefix(name, toComplete) {
			completions = append(completions, name)
		}
	}

	return completions, cobra.ShellCompDirectiveNoFileComp
}
