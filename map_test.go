package gomap

import (
	"testing"
)

// TODO: do bencmarks and comppare with default map type

type CustomKey struct {
	ID   int
	Name string
}

func TestMap(t *testing.T) {
	hashFunc := func(k int) int {
		return k
	}
	equalFunc := func(a, b int) bool {
		return a == b
	}

	m := New[int, string](10, hashFunc, equalFunc)

	t.Run("Set and Get", func(t *testing.T) {
		m.Set(1, "one")
		m.Set(2, "two")
		m.Set(11, "eleven") // Collision with key 1

		if val, ok := m.Get(1); !ok || val != "one" {
			t.Errorf("Expected value 'one', got '%v' (ok: %v)", val, ok)
		}

		if val, ok := m.Get(2); !ok || val != "two" {
			t.Errorf("Expected value 'two', got '%v' (ok: %v)", val, ok)
		}

		if val, ok := m.Get(11); !ok || val != "eleven" {
			t.Errorf("Expected value 'eleven', got '%v' (ok: %v)", val, ok)
		}

		if _, ok := m.Get(3); ok {
			t.Errorf("Expected key 3 to not exist")
		}
	})

	t.Run("Update Value", func(t *testing.T) {
		m.Set(1, "uno")
		if val, ok := m.Get(1); !ok || val != "uno" {
			t.Errorf("Expected updated value 'uno', got '%v' (ok: %v)", val, ok)
		}
	})

	t.Run("Delete", func(t *testing.T) {
		m.Delete(1)
		if _, ok := m.Get(1); ok {
			t.Errorf("Expected key 1 to be deleted")
		}
		if m.Size() != 2 {
			t.Errorf("Expected size 2, got %d", m.Size())
		}
	})

	t.Run("Size", func(t *testing.T) {
		m.Set(3, "three")
		m.Set(4, "four")
		if m.Size() != 4 {
			t.Errorf("Expected size 4, got %d", m.Size())
		}
	})

	t.Run("String Keys", func(t *testing.T) {
		stringMap := New[string, int](10, func(k string) int {
			sum := 0
			for _, c := range k {
				sum += int(c)
			}
			return sum
		}, func(a, b string) bool {
			return a == b
		})

		stringMap.Set("one", 1)
		stringMap.Set("two", 2)
		if val, ok := stringMap.Get("one"); !ok || val != 1 {
			t.Errorf("Expected value 1 for key 'one', got '%v' (ok: %v)", val, ok)
		}
		if val, ok := stringMap.Get("two"); !ok || val != 2 {
			t.Errorf("Expected value 2 for key 'two', got '%v' (ok: %v)", val, ok)
		}
	})

	t.Run("Custom Struct Keys", func(t *testing.T) {
		structMap := New[CustomKey, string](10, func(k CustomKey) int {
			return k.ID
		}, func(a, b CustomKey) bool {
			return a.ID == b.ID && a.Name == b.Name
		})

		key1 := CustomKey{ID: 1, Name: "Alice"}
		key2 := CustomKey{ID: 2, Name: "Bob"}
		structMap.Set(key1, "Developer")
		structMap.Set(key2, "Designer")

		if val, ok := structMap.Get(key1); !ok || val != "Developer" {
			t.Errorf("Expected value 'Developer' for key %+v, got '%v' (ok: %v)", key1, val, ok)
		}
		if val, ok := structMap.Get(key2); !ok || val != "Designer" {
			t.Errorf("Expected value 'Designer' for key %+v, got '%v' (ok: %v)", key2, val, ok)
		}
	})
}
