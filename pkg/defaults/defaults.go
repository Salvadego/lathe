package defaults

import "github.com/Salvadego/lathe/internal/config"

func Config() config.Config {
	return config.Config{
		Compiler: config.CompilerConfig{
			CC: "gcc",
		},

		Builds: map[string]config.BuildConfig{
			"debug": {
				CFlags: []string{
					"-ggdb",
					"-O0",
				},
				Sanitize: true,
			},
			"release": {
				CFlags: []string{
					"-O3",
				},
			},
		},
		Targets: make(map[string]config.TargetConfig),
		Aliases: make(map[string]config.AliasConfig),
		Flags:   config.FlagsConfig{},
	}
}
