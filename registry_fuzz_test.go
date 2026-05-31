package registry_test

import (
	"errors"
	"testing"

	"github.com/brandonkramer/registry"
)

func FuzzRegistryRegisterGet(f *testing.F) {
	f.Add("alpha")
	f.Fuzz(func(t *testing.T, name string) {
		reg := registry.New[string]()
		if name == "" {
			return
		}
		if err := reg.Register(name, name); err != nil {
			return
		}
		got, err := reg.Get(name)
		if err != nil || got != name {
			t.Fatalf("Get(%q)=%q err=%v", name, got, err)
		}
		if !reg.Has(name) {
			t.Fatalf("Has(%q)=false", name)
		}
	})
}

func FuzzRegistryNotFound(f *testing.F) {
	f.Add("missing")
	f.Fuzz(func(t *testing.T, name string) {
		reg := registry.New[int]()
		if reg.Has(name) {
			t.Skip("unexpected preexisting name")
		}
		_, err := reg.Get(name)
		if !errors.Is(err, registry.ErrNotFound) {
			t.Fatalf("err=%v", err)
		}
	})
}
