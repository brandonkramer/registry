package registry

//
// ────────────────────────────────────────
// registry options.
//

// Option configures a [Registry].
type Option[T any] func(*Registry[T])

// WithValidator registers a validation hook run before each [Registry.Register].
func WithValidator[T any](validate func(T) error) Option[T] {
	return func(r *Registry[T]) { r.validate = validate }
}

// WithKeyFrom extracts registry keys from values for [Registry.RegisterItem].
func WithKeyFrom[T any](keyFrom func(T) string) Option[T] {
	return func(r *Registry[T]) { r.keyFrom = keyFrom }
}

// WithRejectDuplicates makes [Registry.Register] return [ErrExists] instead of replacing an entry.
func WithRejectDuplicates[T any]() Option[T] {
	return func(r *Registry[T]) { r.rejectDup = true }
}
