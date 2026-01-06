package build

import (
	"bytes"
	"os"
	"path/filepath"

	"github.com/Salvadego/lathe/internal/types"
)

func needsRebuildBinary(output string, objects []string) bool {
	out, err := os.Stat(output)
	if err != nil {
		return true
	}

	for _, o := range objects {
		oi, err := os.Stat(o)
		if err != nil {
			return true
		}
		if oi.ModTime().After(out.ModTime()) {
			return true
		}
	}
	return false
}

func needsRebuild(ctx *types.BuildContext, src, obj string) bool {
	sig := buildSignature(ctx, src)

	old, err := os.ReadFile(sigPath(obj))
	if err != nil {
		return true
	}

	if !bytes.Equal(sig, old) {
		return true
	}

	if _, err := os.Stat(obj); err != nil {
		return true
	}

	return false
}

func Build(ctx *types.BuildContext) (string, error) {
	objects, err := CompileObjects(ctx)
	if err != nil {
		return "", err
	}

	out := filepath.Join(ctx.WorkDir, ctx.Output)

	if ctx.Incremental {
		if !needsRebuildBinary(out, objects) {
			return out, nil
		}
	}

	return Link(ctx, objects)
}
