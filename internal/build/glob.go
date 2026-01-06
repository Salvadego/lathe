package build

import (
	"fmt"
	"path/filepath"
	"sort"
)

func expandGlobs(patterns []string) ([]string, error) {
	seen := make(map[string]struct{})
	var out []string

	for _, p := range patterns {
		matches, err := filepath.Glob(p)
		if err != nil {
			return nil, fmt.Errorf("invalid glob %q: %w", p, err)
		}

		if len(matches) == 0 {
			return nil, fmt.Errorf("glob %q matched no files", p)
		}

		sort.Strings(matches)

		for _, m := range matches {
			if _, ok := seen[m]; ok {
				continue
			}
			seen[m] = struct{}{}
			out = append(out, m)
		}
	}

	return out, nil
}
