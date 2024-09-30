package format

// Ptr returns the pointer of v.
func Ptr[T any](v T) *T {
	return &v
}
