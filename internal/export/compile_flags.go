package export

import (
	"os"
	"strings"

	"github.com/Salvadego/lathe/internal/types"
)

func WriteCompileFlags(path string, ctx *types.BuildContext) error {
	var flags []string
	flags = append(flags, ctx.CFlags...)

	for _, d := range ctx.Defines {
		flags = append(flags, "-D"+d)
	}
	for _, i := range ctx.Includes {
		flags = append(flags, "-I"+i)
	}

	return os.WriteFile(path, []byte(strings.Join(flags, "\n")), 0644)
}
