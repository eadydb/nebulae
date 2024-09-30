//go:build windows

package path

import (
	"path/filepath"
	"strings"
)

// Join is a wrapper around filepath.Join that converts all path separators to
// forward slashes. This is useful for Windows paths that are used in URLs.
func Join(elem ...string) string {
	return strings.Replace(filepath.Join(elem...), `\`, `/`, -1)
}
