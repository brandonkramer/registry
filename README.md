# registry

Thread-safe generic name registry with optional validation, key extraction, and duplicate policy.

## Quick start

```go
reg := registry.New[Plugin](
    registry.WithValidator(func(p Plugin) error {
        if p.Name == "" {
            return fmt.Errorf("empty name")
        }
        return nil
    }),
    registry.WithKeyFrom(func(p Plugin) string { return p.Name }),
)
reg.MustRegisterItem(Plugin{Name: "alpha"})

p, err := reg.Get("alpha")
if errors.Is(err, registry.ErrNotFound) {
    // unknown name
}

snap := reg.Snapshot() // safe map copy for iteration
```

## Options

| Option | Effect |
|--------|--------|
| `WithValidator(fn)` | Validate before register |
| `WithKeyFrom(fn)` | Enable `RegisterItem` / `MustRegisterItem` |
| `WithRejectDuplicates()` | `Register` returns `ErrExists` instead of replacing |

By default, registering the same name **replaces** the previous value.
