package build

import (
	"fmt"

	"github.com/Salvadego/lathe/internal/config"
	"github.com/Salvadego/lathe/internal/types"
)

func resolveCompiler(cfg config.Config) string {
	if cfg.Compiler.CC != "" {
		return cfg.Compiler.CC
	}
	return "cc"
}

func resolveSources(
	req types.BuildRequest,
	cfg config.Config,
) ([]string, error) {

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

	return expandGlobs(sources)
}

func resolveFlags(
	req types.BuildRequest,
	cfg config.Config,
) (cflags, ldflags, includes, defines []string, err error) {

	cflags = append(cflags, cfg.Flags.CFlags...)
	ldflags = append(ldflags, cfg.Flags.LdFlags...)

	if buildCfg, ok := cfg.Builds[string(req.Mode)]; ok {

		cflags = append(cflags, buildCfg.CFlags...)
		ldflags = append(ldflags, buildCfg.LdFlags...)

		for _, aliasName := range buildCfg.Aliases {
			alias, ok := cfg.Aliases[aliasName]
			if !ok {
				return nil, nil, nil, nil,
					fmt.Errorf("unknown alias %q", aliasName)
			}

			cflags = append(cflags, alias.CFlags...)
			ldflags = append(ldflags, alias.LdFlags...)
			includes = append(includes, alias.Includes...)
			defines = append(defines, alias.Defines...)
		}
	}

	if target, ok := cfg.Targets[req.TargetName]; ok {
		includes = append(includes, target.Includes...)
		defines = append(defines, target.Defines...)

		for _, aliasName := range target.Aliases {
			alias, ok := cfg.Aliases[aliasName]
			if !ok {
				return nil, nil, nil, nil,
					fmt.Errorf("unknown alias %q", aliasName)
			}

			cflags = append(cflags, alias.CFlags...)
			ldflags = append(ldflags, alias.LdFlags...)
			includes = append(includes, alias.Includes...)
			defines = append(defines, alias.Defines...)
		}
	}

	return
}

func ResolveBuildContext(
	req types.BuildRequest,
	cfg config.Config,
) (*types.BuildContext, error) {

	cc := resolveCompiler(cfg)

	sources, err := resolveSources(req, cfg)
	if err != nil {
		return nil, err
	}

	cflags, ldflags, includes, defines, err :=
		resolveFlags(req, cfg)
	if err != nil {
		return nil, err
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
