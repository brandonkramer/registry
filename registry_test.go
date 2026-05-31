package registry_test

import (
	"errors"
	"fmt"
	"sync"
	"testing"

	"github.com/brandonkramer/registry"
)

type item struct {
	Name string
}

func validator(v item) error {
	if v.Name == "" {
		return fmt.Errorf("empty name")
	}
	return nil
}

func TestRegistryRoundTrip(t *testing.T) {
	t.Parallel()
	reg := registry.New[item](registry.WithValidator(validator))
	if err := reg.Register("a", item{Name: "a"}); err != nil {
		t.Fatal(err)
	}
	got, err := reg.Get("a")
	if err != nil || got.Name != "a" {
		t.Fatalf("Get: %+v err=%v", got, err)
	}
	if !reg.Has("a") || reg.Len() != 1 {
		t.Fatalf("has=%v len=%d", reg.Has("a"), reg.Len())
	}
	if names := reg.Names(); len(names) != 1 || names[0] != "a" {
		t.Fatalf("names=%v", names)
	}
	if vals := reg.Values(); len(vals) != 1 || vals[0].Name != "a" {
		t.Fatalf("values=%v", vals)
	}
}

func TestRegistryNotFound(t *testing.T) {
	t.Parallel()
	reg := registry.New[item]()
	_, err := reg.Get("missing")
	if !errors.Is(err, registry.ErrNotFound) {
		t.Fatalf("err=%v", err)
	}
}

func TestRegistryValidation(t *testing.T) {
	t.Parallel()
	reg := registry.New[item](registry.WithValidator(validator))
	if err := reg.Register("bad", item{}); err == nil {
		t.Fatal("expected validation error")
	}
}

func TestMustRegisterPanics(t *testing.T) {
	t.Parallel()
	reg := registry.New[item](registry.WithValidator(validator))
	defer func() {
		if recover() == nil {
			t.Fatal("expected panic")
		}
	}()
	reg.MustRegister("bad", item{})
}

func TestMustRegisterItem(t *testing.T) {
	t.Parallel()
	reg := registry.New[item](registry.WithKeyFrom(func(v item) string { return v.Name }))
	reg.MustRegisterItem(item{Name: "a"})
	got, err := reg.Get("a")
	if err != nil || got.Name != "a" {
		t.Fatalf("Get: %+v err=%v", got, err)
	}
}

func TestMustRegisterItemPanics(t *testing.T) {
	t.Parallel()
	reg := registry.New[item]()
	defer func() {
		if recover() == nil {
			t.Fatal("expected panic")
		}
	}()
	reg.MustRegisterItem(item{Name: "a"})
}

func TestUnregister(t *testing.T) {
	t.Parallel()
	reg := registry.New[item]()
	reg.MustRegister("a", item{Name: "a"})
	if !reg.Unregister("a") {
		t.Fatal("expected true")
	}
	if reg.Unregister("a") {
		t.Fatal("expected false on second remove")
	}
	if reg.Len() != 0 {
		t.Fatal("expected empty registry")
	}
}

func TestRegistryConcurrentRegister(t *testing.T) {
	t.Parallel()
	reg := registry.New[int]()
	var wg sync.WaitGroup
	for i := range 32 {
		wg.Add(1)
		go func(n int) {
			defer wg.Done()
			name := fmt.Sprintf("n%d", n)
			if err := reg.Register(name, n); err != nil {
				t.Errorf("register %s: %v", name, err)
			}
		}(i)
	}
	wg.Wait()
	if reg.Len() != 32 {
		t.Fatalf("len=%d", reg.Len())
	}
}

func TestRegistryNotFoundWrapped(t *testing.T) {
	t.Parallel()
	reg := registry.New[int]()
	_, err := reg.Get("missing")
	if err == nil {
		t.Fatal("expected error")
	}
	if !errors.Is(err, registry.ErrNotFound) {
		t.Fatalf("err=%v", err)
	}
}

func TestRegisterItemWithKeyFrom(t *testing.T) {
	t.Parallel()
	reg := registry.New[item](
		registry.WithValidator(validator),
		registry.WithKeyFrom(func(v item) string { return v.Name }),
	)
	if err := reg.RegisterItem(item{Name: "a"}); err != nil {
		t.Fatal(err)
	}
	got, err := reg.Get("a")
	if err != nil || got.Name != "a" {
		t.Fatalf("Get: %+v err=%v", got, err)
	}
}

func TestRegisterItemNoKeyFrom(t *testing.T) {
	t.Parallel()
	reg := registry.New[item]()
	if err := reg.RegisterItem(item{Name: "a"}); !errors.Is(err, registry.ErrNoKeyFrom) {
		t.Fatalf("err=%v", err)
	}
}

func TestRegisterItemEmptyKey(t *testing.T) {
	t.Parallel()
	reg := registry.New[item](registry.WithKeyFrom(func(v item) string { return v.Name }))
	if err := reg.RegisterItem(item{}); err == nil {
		t.Fatal("expected empty key error")
	}
}

func TestRejectDuplicates(t *testing.T) {
	t.Parallel()
	reg := registry.New[item](registry.WithRejectDuplicates[item]())
	if err := reg.Register("a", item{Name: "a"}); err != nil {
		t.Fatal(err)
	}
	err := reg.Register("a", item{Name: "a"})
	if !errors.Is(err, registry.ErrExists) {
		t.Fatalf("err=%v", err)
	}
}

func TestRegisterReplacesByDefault(t *testing.T) {
	t.Parallel()
	reg := registry.New[int]()
	if err := reg.Register("a", 1); err != nil {
		t.Fatal(err)
	}
	if err := reg.Register("a", 2); err != nil {
		t.Fatal(err)
	}
	got, err := reg.Get("a")
	if err != nil || got != 2 {
		t.Fatalf("got=%d err=%v", got, err)
	}
}

func TestSnapshotCopy(t *testing.T) {
	t.Parallel()
	reg := registry.New[int]()
	reg.MustRegister("a", 1)
	snap := reg.Snapshot()
	snap["b"] = 2
	if reg.Len() != 1 || reg.Has("b") {
		t.Fatalf("snapshot mutation affected registry: len=%d hasB=%v", reg.Len(), reg.Has("b"))
	}
}
