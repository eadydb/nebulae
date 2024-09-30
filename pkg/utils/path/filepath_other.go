//go:build !windows

package path

import (
	"path/filepath"
)

// Join is a wrapper around filepath.Join.
func Join(elem ...string) string {
	return filepath.Join(elem...)
}
