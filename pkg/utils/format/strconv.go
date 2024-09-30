package format

import (
	"fmt"
)

// Parse returns the value of the string representation of u.
func Parse[T any](s string) (t T, err error) {
	_, err = fmt.Sscan(s, &t)
	return t, err
}

// String returns the string representation of u.
func String(i any) string {
	return fmt.Sprintf("%v", i)
}
