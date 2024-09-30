package path

import (
	"os"
)

// ExpandHome expands home directory in file paths.
func ExpandHome(path string) string {
	if len(path) == 0 {
		return path
	}
	if path[0] != '~' {
		return path
	}
	if len(path) > 1 && path[1] != '/' && path[1] != '\\' {
		return path
	}

	home, err := os.UserHomeDir()
	if err != nil {
		return path
	}

	return Join(home, path[1:])
}
