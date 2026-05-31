// Package registry provides a thread-safe generic name registry with optional validation.
//
// Use [New] with [WithValidator] to reject invalid entries at registration time.
// Use [WithKeyFrom] and [Registry.RegisterItem] when the value carries its own key.
// Use [WithRejectDuplicates] when re-registering the same name should return [ErrExists].
// [Registry.Snapshot] returns a map copy for iteration outside the registry lock.
//
// # Example
//
//	reg := registry.New[Driver](
//	    registry.WithValidator(func(d Driver) error {
//	        if d.Name == "" {
//	            return fmt.Errorf("empty name")
//	        }
//	        return nil
//	    }),
//	    registry.WithKeyFrom(func(d Driver) string { return d.Name }),
//	)
//	reg.MustRegisterItem(Driver{Name: "alpha"})
//	d, err := reg.Get("alpha")
package registry
