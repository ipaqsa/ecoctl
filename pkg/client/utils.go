package client

// PtrTo returns a pointer to the provided input.
func PtrTo[T any](v T) *T {
	return &v
}
