package registry

import (
	"fmt"
	"maps"
	"slices"
	"sync"
)

//
// ────────────────────────────────────────
// generic name registry.
//

// Registry maps unique names to registered values.
type Registry[T any] struct {
	mu       sync.RWMutex
	m        map[string]T
	validate func(T) error
}

// New returns an empty registry configured by opts.
func New[T any](opts ...Option[T]) *Registry[T] {
	r := &Registry[T]{m: make(map[string]T)}
	for _, opt := range opts {
		opt(r)
	}
	return r
}

// Register stores item under name.
func (r *Registry[T]) Register(name string, item T) error {
	if r.validate != nil {
		if err := r.validate(item); err != nil {
			return fmt.Errorf("registry: validate %q: %w", name, err)
		}
	}
	r.mu.Lock()
	r.m[name] = item
	r.mu.Unlock()
	return nil
}

// MustRegister panics when [Register] fails.
func (r *Registry[T]) MustRegister(name string, item T) {
	if err := r.Register(name, item); err != nil {
		panic(err)
	}
}

// Unregister removes name when present and reports whether it existed.
func (r *Registry[T]) Unregister(name string) bool {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.m[name]; !ok {
		return false
	}
	delete(r.m, name)
	return true
}

// Has reports whether name is registered.
func (r *Registry[T]) Has(name string) bool {
	r.mu.RLock()
	defer r.mu.RUnlock()
	_, ok := r.m[name]
	return ok
}

// Len returns the number of registered names.
func (r *Registry[T]) Len() int {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return len(r.m)
}

// Names returns registered names in sorted order.
func (r *Registry[T]) Names() []string {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return slices.Sorted(maps.Keys(r.m))
}

// Get returns the item registered under name.
func (r *Registry[T]) Get(name string) (T, error) {
	r.mu.RLock()
	item, ok := r.m[name]
	r.mu.RUnlock()
	if !ok {
		var zero T
		return zero, fmt.Errorf("registry: lookup %q: %w", name, ErrNotFound)
	}
	return item, nil
}

// Values returns all registered items in arbitrary order.
func (r *Registry[T]) Values() []T {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return slices.Collect(maps.Values(r.m))
}
