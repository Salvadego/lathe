package build

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/Salvadego/lathe/internal/types"
)

func objectPath(ctx *types.BuildContext, src string) string {
	base := filepath.Clean(src)
	return filepath.Join(
		ctx.WorkDir,
		string(ctx.Mode),
		base+".o",
	)
}

func CompileObjects(ctx *types.BuildContext) ([]string, error) {
	cmds, err := CompileCommands(ctx)
	if err != nil {
		return nil, err
	}

	var objects []string

	for _, c := range cmds {
		objects = append(objects, c.Output)

		if !needsRebuild(c.File, c.Output) {
			continue
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

func Makefile(ctx *types.BuildContext) string {
	var b strings.Builder

	b.WriteString("CC := " + ctx.Compiler + "\n")
	b.WriteString("CFLAGS := " + strings.Join(ctx.CFlags, " ") + "\n")
	b.WriteString("LDFLAGS := " + strings.Join(ctx.LdFlags, " ") + "\n\n")

	var objs []string
	for _, s := range ctx.Sources {
		objs = append(objs, objectPath(ctx, s))
	}

	b.WriteString("OBJS := " + strings.Join(objs, " ") + "\n\n")
	b.WriteString(ctx.Output + ": $(OBJS)\n")
	b.WriteString("\t$(CC) $(OBJS) -o $@ $(LDFLAGS)\n\n")

	for _, s := range ctx.Sources {
		o := objectPath(ctx, s)
		b.WriteString(o + ": " + s + "\n")
		b.WriteString("\t$(CC) $(CFLAGS) -c " + s + " -o " + o + "\n\n")
	}

	return b.String()
}
