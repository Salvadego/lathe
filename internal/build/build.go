package build

import (
	"github.com/Salvadego/lathe/internal/types"
	"os"
)

func needsRebuild(src, obj string) bool {
	so, err1 := os.Stat(src)
	oo, err2 := os.Stat(obj)

	if err2 != nil {
		return true
	}

	if err1 != nil {
		return true
	}

	return so.ModTime().After(oo.ModTime())
}

func Build(ctx *types.BuildContext) (string, error) {
	objects, err := CompileObjects(ctx)
	if err != nil {
		return "", err
	}

	return Link(ctx, objects)
}
