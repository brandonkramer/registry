// Package registry provides a thread-safe generic name registry with optional validation.
//
// Use [New] with [WithValidator] to reject invalid entries at registration time.
// [Registry.Get] returns [ErrNotFound] for unknown names; use [errors.Is] to detect it.
//
// # Example
//
//	reg := registry.New[Driver](registry.WithValidator(func(d Driver) error {
//	    if d.Name == "" {
//	        return fmt.Errorf("empty name")
//	    }
//	    return nil
//	}))
//	reg.MustRegister("alpha", Driver{Name: "alpha"})
//	d, err := reg.Get("alpha")
package registry
