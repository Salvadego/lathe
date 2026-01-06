package build

import (
	"fmt"

	"github.com/Salvadego/lathe/internal/config"
	"github.com/Salvadego/lathe/internal/types"
)

func ResolveBuildContext(
	req types.BuildRequest,
	cfg config.Config,
) (*types.BuildContext, error) {
	cc := cfg.Compiler.CC
	if cc == "" {
		cc = "cc"
	}

	sources := req.Sources

	if len(sources) == 0 && req.TargetName != "" {
		target, ok := cfg.Targets[req.TargetName]
		if !ok {
			return nil, fmt.Errorf("unknown target %q", req.TargetName)
		}
		sources = target.Sources
	}

	if len(sources) == 0 {
		return nil, fmt.Errorf("no source files provided")
	}

	sources, err := expandGlobs(sources)
	if err != nil {
		return nil, err
	}

	var cflags []string
	var ldflags []string
	var includes []string
	var defines []string

	cflags = append(cflags, cfg.Flags.CFlags...)

	buildCfg, ok := cfg.Builds[string(req.Mode)]
	if ok {
		cflags = append(cflags, buildCfg.CFlags...)
		ldflags = append(ldflags, buildCfg.LdFlags...)
		for _, aliasName := range buildCfg.Aliases {
			alias, ok := cfg.Aliases[aliasName]
			if !ok {
				return nil, fmt.Errorf("unknown alias %q", aliasName)
			}

			cflags = append(cflags, alias.CFlags...)
			ldflags = append(ldflags, alias.LdFlags...)
			includes = append(includes, alias.Includes...)
			defines = append(defines, alias.Defines...)
		}
	}

	ldflags = append(ldflags, cfg.Flags.LdFlags...)

	if target, ok := cfg.Targets[req.TargetName]; ok {
		includes = append(includes, target.Includes...)
		defines = append(defines, target.Defines...)

		for _, aliasName := range target.Aliases {
			alias, ok := cfg.Aliases[aliasName]
			if !ok {
				return nil, fmt.Errorf("unknown alias %q", aliasName)
			}

			cflags = append(cflags, alias.CFlags...)
			ldflags = append(ldflags, alias.LdFlags...)
			includes = append(includes, alias.Includes...)
			defines = append(defines, alias.Defines...)
		}
	}

	output := "a.out"
	if req.TargetName != "" {
		output = req.TargetName
	}

	workDir := ".lathe"
	if req.WorkDir != "" {
		workDir = req.WorkDir
	}

	ctx := &types.BuildContext{
		Mode:        req.Mode,
		Compiler:    cc,
		Sources:     sources,
		Includes:    includes,
		Defines:     defines,
		CFlags:      cflags,
		LdFlags:     ldflags,
		Output:      output,
		Verbose:     req.Verbose,
		WorkDir:     workDir,
		Incremental: req.Incremental,
	}

	return ctx, nil
}
