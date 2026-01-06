package config

import (
	"fmt"

	"github.com/BurntSushi/toml"
)

func Load(path string) (Config, error) {
	var cfg Config

	md, err := toml.DecodeFile(path, &cfg)
	if err != nil {
		return Config{}, err
	}

	if undecoded := md.Undecoded(); len(undecoded) > 0 {
		return Config{}, fmt.Errorf(
			"unknown config keys in %s: %v",
			path,
			undecoded,
		)
	}

	return cfg, nil
}
