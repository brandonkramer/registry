# registry

Thread-safe generic name registry with optional validation.

## Install

From [pkg.go.dev](https://pkg.go.dev/github.com/brandonkramer/registry):

```bash
go get github.com/brandonkramer/registry
```

## Quick start

```go
reg := registry.New[Plugin](registry.WithValidator(func(p Plugin) error {
    if p.Name == "" {
        return fmt.Errorf("empty name")
    }
    return nil
}))
reg.MustRegister("alpha", Plugin{Name: "alpha"})
p, err := reg.Get("alpha")
if errors.Is(err, registry.ErrNotFound) {
    // unknown name
}
```

## Development

```bash
make check
```
