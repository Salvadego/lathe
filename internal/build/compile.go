package build

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"

	"crypto/sha1"
	"crypto/sha256"
	"encoding/hex"

	"github.com/Salvadego/lathe/internal/config"
	"github.com/Salvadego/lathe/internal/types"
)

func objectPath(ctx *types.BuildContext, src string) string {
	h := sha1.Sum([]byte(src))
	name := hex.EncodeToString(h[:8])
	return filepath.Join(
		ctx.WorkDir,
		string(ctx.Mode),
		name+".o",
	)
}

func buildSignature(ctx *types.BuildContext, src string) []byte {
	h := sha256.New()

	io.WriteString(h, ctx.Compiler)
	io.WriteString(h, string(ctx.Mode))
	io.WriteString(h, src)

	for _, f := range ctx.CFlags {
		io.WriteString(h, f)
	}
	for _, d := range ctx.Defines {
		io.WriteString(h, d)
	}
	for _, i := range ctx.Includes {
		io.WriteString(h, i)
	}

	return h.Sum(nil)
}

func sigPath(obj string) string {
	return obj + ".sig"
}

func CompileObjects(ctx *types.BuildContext) ([]string, error) {
	cmds, err := CompileCommands(ctx)
	if err != nil {
		return nil, err
	}

	var objects []string

	for _, c := range cmds {
		objects = append(objects, c.Output)

		if ctx.Incremental {
			if !needsRebuild(ctx, c.File, c.Output) {
				continue
			}
		}

		if err := os.MkdirAll(filepath.Dir(c.Output), 0o755); err != nil {
			return nil, err
		}

		cmd := exec.Command(c.Command[0], c.Command[1:]...)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		if ctx.Verbose {
			fmt.Println(c.Command)
		}

		if err := cmd.Run(); err != nil {
			return nil, err
		}

		sig := buildSignature(ctx, c.File)
		if err := os.WriteFile(sigPath(c.Output), sig, 0644); err != nil {
			return nil, err
		}
	}

	return objects, nil
}

func CompileCommands(ctx *types.BuildContext) ([]types.CompileCommand, error) {

	if err := os.MkdirAll(ctx.WorkDir, 0o755); err != nil {
		return nil, err
	}

	cwd, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	var cmds []types.CompileCommand

	for _, src := range ctx.Sources {
		args := []string{}

		args = append(args, ctx.CFlags...)

		for _, def := range ctx.Defines {
			args = append(args, "-D"+def)
		}
		for _, inc := range ctx.Includes {
			args = append(args, "-I"+inc)
		}

		out := objectPath(ctx, src)

		args = append(args,
			"-c", src,
			"-o", out,
		)

		cmds = append(cmds, types.CompileCommand{
			Directory: cwd,
			Command:   append([]string{ctx.Compiler}, args...),
			File:      src,
			Output:    out,
		})
	}

	return cmds, nil
}

func Link(ctx *types.BuildContext, objects []string) (string, error) {
	out := filepath.Join(ctx.WorkDir, ctx.Output)
	args := append(objects, "-o", out)
	args = append(args, ctx.LdFlags...)

	cmd := exec.Command(ctx.Compiler, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if ctx.Verbose {
		fmt.Println(append([]string{ctx.Compiler}, args...))
	}

	return out, cmd.Run()
}

func Makefile(cfg config.Config, mode types.Mode) (string, error) {
	if len(cfg.Targets) == 0 {
		return "", fmt.Errorf("no targets defined in config")
	}

	var b strings.Builder

	cc := resolveCompiler(cfg)

	b.WriteString("CC := " + cc + "\n")
	b.WriteString("BUILD_DIR := build\n")
	b.WriteString("MODE := " + string(mode) + "\n\n")

	// Collect targets
	var targetNames []string
	for name := range cfg.Targets {
		targetNames = append(targetNames, name)
	}
	sort.Strings(targetNames)

	b.WriteString("TARGETS := " + strings.Join(targetNames, " ") + "\n")
	b.WriteString(".DEFAULT_GOAL := all\n")
	b.WriteString(".PHONY: all clean run help list\n\n")

	b.WriteString("all: $(TARGETS)\n\n")

	b.WriteString("ALL_OBJS :=\n\n")

	for _, name := range targetNames {

		target := cfg.Targets[name]
		if len(target.Sources) == 0 {
			return "", fmt.Errorf("target %q has no sources", name)
		}

		req := types.BuildRequest{
			Mode:       mode,
			TargetName: name,
		}

		cflags, ldflags, includes, defines, err :=
			resolveFlags(req, cfg)
		if err != nil {
			return "", err
		}

		sort.Strings(target.Sources)

		// Per-target flags
		b.WriteString(name + "_CFLAGS := " + strings.Join(cflags, " ") + "\n")
		for _, d := range defines {
			b.WriteString(name + "_CFLAGS += -D" + d + "\n")
		}
		for _, i := range includes {
			b.WriteString(name + "_CFLAGS += -I" + i + "\n")
		}
		b.WriteString(name + "_CFLAGS += -MMD -MP\n")

		b.WriteString(name + "_LDFLAGS := " + strings.Join(ldflags, " ") + "\n\n")

		// Sources
		b.WriteString(name + "_SRCS := \\\n")
		for i, s := range target.Sources {
			if i == len(target.Sources)-1 {
				b.WriteString("\t" + s + "\n")
			} else {
				b.WriteString("\t" + s + " \\\n")
			}
		}
		b.WriteString("\n")

		// Objects
		b.WriteString(name + "_OBJS := $(" + name + "_SRCS:%.c=$(BUILD_DIR)/$(MODE)/" + name + "/%.o)\n\n")
		b.WriteString("ALL_OBJS += $(" + name + "_OBJS)\n\n")

		// Link rule
		b.WriteString(name + ": $(" + name + "_OBJS)\n")
		b.WriteString("\t$(CC) $^ -o $@ $(" + name + "_LDFLAGS)\n\n")

		// Compile rule per target
		b.WriteString("$(BUILD_DIR)/$(MODE)/" + name + "/%.o: %.c\n")
		b.WriteString("\t@mkdir -p $(dir $@)\n")
		b.WriteString("\t$(CC) $(" + name + "_CFLAGS) -c $< -o $@\n\n")
	}

	b.WriteString("DEPS := $(ALL_OBJS:.o=.d)\n\n")
	b.WriteString("-include $(DEPS)\n\n")

	// Utilities
	first := targetNames[0]

	b.WriteString("run: " + first + "\n")
	b.WriteString("\t./" + first + "\n\n")

	b.WriteString("clean:\n")
	b.WriteString("\trm -rf $(BUILD_DIR) $(TARGETS)\n\n")

	b.WriteString("list:\n")
	b.WriteString("\t@echo \"Available targets:\"\n")
	b.WriteString("\t@for t in $(TARGETS); do echo \"  $$t\"; done\n\n")

	b.WriteString("help:\n")
	b.WriteString("\t@echo \"\"\n")
	b.WriteString("\t@echo \"Lathe Makefile\"\n")
	b.WriteString("\t@echo \"\"\n")
	b.WriteString("\t@echo \"  make              Build all targets\"\n")
	b.WriteString("\t@echo \"  make <target>     Build specific binary\"\n")
	b.WriteString("\t@echo \"  make run          Run first target\"\n")
	b.WriteString("\t@echo \"  make clean        Remove artifacts\"\n")
	b.WriteString("\t@echo \"  make list         List targets\"\n")
	b.WriteString("\t@echo \"\"\n")

	return b.String(), nil
}
