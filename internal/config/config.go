package config

import (
	"fmt"
	"os"
)

type (
	CompilerConfig struct {
		CC string `toml:"cc"`
	}

	FlagsConfig struct {
		CFlags  []string `toml:"cflags"`
		LdFlags []string `toml:"ldflags"`
	}

	BuildConfig struct {
		CFlags   []string `toml:"cflags"`
		LdFlags  []string `toml:"ldflags"`
		Sanitize bool     `toml:"sanitize"`
	}

	TargetConfig struct {
		Sources  []string `toml:"sources"`
		Includes []string `toml:"includes"`
		Defines  []string `toml:"defines"`
		Aliases  []string `toml:"aliases"`
	}

	AliasConfig struct {
		Includes []string `toml:"includes"`
		Defines  []string `toml:"defines"`
		CFlags   []string `toml:"cflags"`
		LdFlags  []string `toml:"ldflags"`
	}

	Config struct {
		Compiler CompilerConfig          `toml:"compiler"`
		Flags    FlagsConfig             `toml:"flags"`
		Builds   map[string]BuildConfig  `toml:"build"`
		Targets  map[string]TargetConfig `toml:"target"`
		Aliases  map[string]AliasConfig  `toml:"alias"`
	}
)

func ResolveCfg(cfg Config) Config {
	if p := GlobalConfigPath(); p != "" {
		if _, err := os.Stat(p); err == nil {
			if g, err := Load(p); err == nil {
				cfg = Merge(cfg, g)
			}
		}
	}

	if _, err := os.Stat(ProjectConfigPath()); err == nil {
		if p, err := Load(ProjectConfigPath()); err == nil {
			cfg = Merge(cfg, p)
		}
	}

	return cfg
}

func Validate(cfg Config) error {
	for name, target := range cfg.Targets {
		for _, alias := range target.Aliases {
			if _, ok := cfg.Aliases[alias]; !ok {
				return fmt.Errorf(
					"target %q references unknown alias %q",
					name, alias,
				)
			}
		}
	}
	return nil
}
