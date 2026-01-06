package config

import "maps"

func Merge(base, override Config) Config {
	out := base

	if override.Compiler.CC != "" {
		out.Compiler.CC = override.Compiler.CC
	}

	out.Flags.CFlags = append(out.Flags.CFlags, override.Flags.CFlags...)
	out.Flags.LdFlags = append(out.Flags.LdFlags, override.Flags.LdFlags...)

	maps.Copy(out.Builds, override.Builds)
	maps.Copy(out.Targets, override.Targets)
	maps.Copy(out.Aliases, override.Aliases)

	return out
}
