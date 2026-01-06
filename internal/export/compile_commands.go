package export

import (
	"encoding/json"
	"os"

	"github.com/Salvadego/lathe/internal/types"
)

func WriteCompileCommands(
	path string,
	cmds []types.CompileCommand,
) error {

	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	enc := json.NewEncoder(f)
	enc.SetIndent("", "  ")

	return enc.Encode(cmds)
}
