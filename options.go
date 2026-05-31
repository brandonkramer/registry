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
