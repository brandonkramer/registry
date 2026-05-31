package registry

import "errors"

//
// ────────────────────────────────────────
// sentinel errors.
//

// ErrNotFound is returned when a name is not registered.
var ErrNotFound = errors.New("registry: not found")

// ErrExists is returned when [WithRejectDuplicates] is set and the name is already registered.
var ErrExists = errors.New("registry: already exists")

// ErrNoKeyFrom is returned when [Registry.RegisterItem] is called without [WithKeyFrom].
var ErrNoKeyFrom = errors.New("registry: key function not configured")
