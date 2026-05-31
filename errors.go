package registry

import "errors"

//
// ────────────────────────────────────────
// sentinel errors.
//

// ErrNotFound is returned when a name is not registered.
var ErrNotFound = errors.New("registry: not found")
