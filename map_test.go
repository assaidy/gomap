package gomap

import (
	"math/rand"
	"strconv"
	"testing"
)

type CustomKey struct {
	ID   int
	Name string
}

func TestMap(t *testing.T) {
	hashFunc := func(k int) int {
		return k
	}
	// equalFunc := func(a, b int) bool {
	// 	return a == b
	// }

	m := New[int, string](hashFunc, ComparableEqualFunc[int], 10)

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
		stringMap := New[string, int](func(k string) int {
			sum := 0
			for _, c := range k {
				sum += int(c)
			}
			return sum
		}, ComparableEqualFunc[string], 10)

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
		structMap := New[CustomKey, string](func(k CustomKey) int {
			return k.ID
		}, func(a, b CustomKey) bool {
			return a.ID == b.ID && a.Name == b.Name
		}, 10)

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

func TestIterator(t *testing.T) {
	hashFunc := func(k int) int {
		return k
	}
	// equalFunc := func(a, b int) bool {
	// 	return a == b
	// }

	m := New[int, string](hashFunc, ComparableEqualFunc[int])

	m.Set(1, "one")
	m.Set(2, "two")
	m.Set(3, "three")

	tests := []struct {
		key int
		val string
	}{
		{1, "one"},
		{2, "two"},
		{3, "three"},
	}

	i := 0
	for e := range m.Iterator() {
		if tests[i].key != e.key || tests[i].val != e.val {
			t.Fatalf("expected: %+v got: %+v", tests[i], e)
		}
		i++
	}
}

// Equality function for integers
func intHash(a int) int {
	return a
}

// Equality function for integers
func intEqual(a, b int) bool {
	return a == b
}

// Benchmark setup
func randomInts(n int) []int {
	nums := make([]int, n)
	for i := 0; i < n; i++ {
		nums[i] = rand.Intn(1000000)
	}
	return nums
}

func randomStringMap(n int) map[int]string {
	m := make(map[int]string)
	for i := 0; i < n; i++ {
		m[rand.Intn(1000000)] = strconv.Itoa(i)
	}
	return m
}

// Benchmark Custom Map
func BenchmarkCustomMap_Set(b *testing.B) {
	nums := randomInts(10000)
	m := New[int, string](intHash, intEqual)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		m.Set(nums[i%10000], strconv.Itoa(nums[i%10000]))
	}
}

func BenchmarkCustomMap_Get(b *testing.B) {
	nums := randomInts(10000)
	m := New[int, string](intHash, intEqual)
	for _, n := range nums {
		m.Set(n, strconv.Itoa(n))
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		m.Get(nums[i%10000])
	}
}

func BenchmarkCustomMap_Delete(b *testing.B) {
	nums := randomInts(10000)
	m := New[int, string](intHash, intEqual)
	for _, n := range nums {
		m.Set(n, strconv.Itoa(n))
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		m.Delete(nums[i%10000])
	}
}

// Benchmark Built-in Map
func BenchmarkBuiltinMap_Set(b *testing.B) {
	nums := randomInts(10000)
	m := make(map[int]string)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		m[nums[i%10000]] = strconv.Itoa(nums[i%10000])
	}
}

func BenchmarkBuiltinMap_Get(b *testing.B) {
	nums := randomInts(10000)
	m := randomStringMap(10000)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = m[nums[i%10000]]
	}
}

func BenchmarkBuiltinMap_Delete(b *testing.B) {
	nums := randomInts(10000)
	m := randomStringMap(10000)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		delete(m, nums[i%10000])
	}
}
