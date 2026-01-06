package config

import (
	"github.com/BurntSushi/toml"
)

func Load(path string) (Config, error) {
	var cfg Config

	_, err := toml.DecodeFile(path, &cfg)
	if err != nil {
		return Config{}, err
	}

	return cfg, nil
}
