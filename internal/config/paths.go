package config

import (
	"os"
	"path/filepath"
)

func GlobalConfigPath() string {
	dir, err := os.UserConfigDir()
	if err != nil {
		return ""
	}
	return filepath.Join(dir, "lathe", "lathe.toml")
}

func ProjectConfigPath() string {
	return "lathe.toml"
}
