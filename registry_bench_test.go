package registry_test

import (
	"fmt"
	"testing"

	"github.com/brandonkramer/registry"
)

func BenchmarkRegisterGet(b *testing.B) {
	reg := registry.New[int]()
	for b.Loop() {
		name := fmt.Sprintf("n%d", b.N)
		if err := reg.Register(name, b.N); err != nil {
			b.Fatal(err)
		}
		if _, err := reg.Get(name); err != nil {
			b.Fatal(err)
		}
	}
}
