package concurrent

import (
	"sync"
	"sync/atomic"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestNew creates an initialized concurrent OrderedMap.
func TestNew(t *testing.T) {
	m := New[string, string]()
	assert.NotNil(t, m)
	assert.Equal(t, 0, m.Len())
	assert.NotNil(t, m.mu)
}

// TestSetGet basic set and get operations.
func TestSetGet(t *testing.T) {
	m := New[string, int]()
	m.Set("a", 1)
	m.Set("b", 2)
	m.Set("c", 3)

	assert.Equal(t, 3, m.Len())

	val, ok := m.Get("a")
	assert.True(t, ok)
	assert.Equal(t, 1, *val)

	val, ok = m.Get("b")
	assert.True(t, ok)
	assert.Equal(t, 2, *val)
}

// TestSetOverwrite overwrites existing values.
func TestSetOverwrite(t *testing.T) {
	m := New[string, int]()
	m.Set("a", 1)
	m.Set("a", 10)

	assert.Equal(t, 1, m.Len())

	val, ok := m.Get("a")
	assert.True(t, ok)
	assert.Equal(t, 10, *val)
}

// TestGetOrDefault returns default for missing keys.
func TestGetOrDefault(t *testing.T) {
	m := New[string, string]()
	m.Set("exists", "value")

	assert.Equal(t, "value", m.GetOrDefault("exists", "default"))
	assert.Equal(t, "default", m.GetOrDefault("missing", "default"))
}

// TestRemove removes keys from map.
func TestRemove(t *testing.T) {
	m := New[string, int]()
	m.Set("a", 1)
	m.Set("b", 2)

	kvp, ok := m.Remove("a")
	assert.True(t, ok)
	assert.NotNil(t, kvp)
	assert.Equal(t, "a", kvp.Key)
	assert.Equal(t, 1, kvp.Value)
	assert.Equal(t, 1, m.Len())

	_, ok = m.Get("a")
	assert.False(t, ok)
}

// TestRemoveNonexistent removes a missing key.
func TestRemoveNonexistent(t *testing.T) {
	m := New[string, int]()
	kvp, ok := m.Remove("missing")
	assert.False(t, ok)
	assert.Nil(t, kvp)
}

// TestContains checks key membership.
func TestContains(t *testing.T) {
	m := New[string, int]()
	m.Set("a", 1)

	assert.True(t, m.Contains("a"))
	assert.False(t, m.Contains("missing"))
}

// TestFirstLast gets first and last elements.
func TestFirstLast(t *testing.T) {
	m := New[string, int]()

	assert.Nil(t, m.First())
	assert.Nil(t, m.Last())

	m.Set("a", 1)
	assert.Equal(t, "a", m.First().Key)
	assert.Equal(t, "a", m.Last().Key)

	m.Set("b", 2)
	m.Set("c", 3)
	assert.Equal(t, "a", m.First().Key)
	assert.Equal(t, "c", m.Last().Key)
}

// TestKeys returns ordered keys.
func TestKeys(t *testing.T) {
	m := New[string, int]()
	m.Set("a", 1).Set("b", 2).Set("c", 3)

	keys := m.Keys()
	assert.Equal(t, []string{"a", "b", "c"}, keys)
}

// TestIterator iterates through map.
func TestIterator(t *testing.T) {
	m := New[string, int]()
	m.Set("a", 1).Set("b", 2).Set("c", 3)

	it := m.Iterator()
	var results []string
	for {
		kvp := it.Next()
		if kvp == nil {
			break
		}
		results = append(results, kvp.Key)
	}

	assert.Equal(t, []string{"a", "b", "c"}, results)
}

// TestMoveToFront moves element to front.
func TestMoveToFront(t *testing.T) {
	m := New[string, int]()
	m.Set("a", 1).Set("b", 2).Set("c", 3)

	assert.NoError(t, m.MoveToFront("c"))
	assert.Equal(t, "c", m.First().Key)
	assert.Error(t, m.MoveToFront("missing"))
}

// TestMoveToBack moves element to back.
func TestMoveToBack(t *testing.T) {
	m := New[string, int]()
	m.Set("a", 1).Set("b", 2).Set("c", 3)

	assert.NoError(t, m.MoveToBack("a"))
	assert.Equal(t, "a", m.Last().Key)
	assert.Error(t, m.MoveToBack("missing"))
}

// TestMoveAfter moves element after another.
func TestMoveAfter(t *testing.T) {
	m := New[string, int]()
	m.Set("a", 1).Set("b", 2).Set("c", 3)

	assert.NoError(t, m.MoveAfter("a", "b"))
	assert.Equal(t, []string{"b", "a", "c"}, m.Keys())
}

// TestMoveBefore moves element before another.
func TestMoveBefore(t *testing.T) {
	m := New[string, int]()
	m.Set("a", 1).Set("b", 2).Set("c", 3)

	assert.NoError(t, m.MoveBefore("c", "b"))
	assert.Equal(t, []string{"a", "c", "b"}, m.Keys())
}

// TestInsertAfter inserts element after another.
func TestInsertAfter(t *testing.T) {
	m := New[string, int]()
	m.Set("a", 1).Set("c", 3)

	assert.NoError(t, m.InsertAfter("b", 2, "a"))
	assert.Equal(t, []string{"a", "b", "c"}, m.Keys())
}

// TestInsertBefore inserts element before another.
func TestInsertBefore(t *testing.T) {
	m := New[string, int]()
	m.Set("a", 1).Set("c", 3)

	assert.NoError(t, m.InsertBefore("b", 2, "c"))
	assert.Equal(t, []string{"a", "b", "c"}, m.Keys())
}

// TestStringAndGoString string representations.
func TestStringAndGoString(t *testing.T) {
	m := New[string, int]()
	m.Set("a", 1).Set("b", 2)

	assert.NotEmpty(t, m.String())
	assert.NotEmpty(t, m.GoString())
}

// TestConcurrentWrites concurrent writes are safe.
func TestConcurrentWrites(t *testing.T) {
	m := New[int, int]()
	numGoroutines := 10
	opsPerGoroutine := 100

	var wg sync.WaitGroup
	wg.Add(numGoroutines)

	for i := 0; i < numGoroutines; i++ {
		go func(id int) {
			defer wg.Done()
			for j := 0; j < opsPerGoroutine; j++ {
				key := id*opsPerGoroutine + j
				m.Set(key, key*2)
			}
		}(i)
	}

	wg.Wait()

	assert.Equal(t, numGoroutines*opsPerGoroutine, m.Len())
}

// TestConcurrentReadsWrites concurrent reads and writes are safe.
func TestConcurrentReadsWrites(t *testing.T) {
	m := New[int, int]()

	for i := 0; i < 100; i++ {
		m.Set(i, i*2)
	}

	numReaders := 5
	numWriters := 5
	opsPerGoroutine := 50
	var wg sync.WaitGroup

	for i := 0; i < numReaders; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < opsPerGoroutine; j++ {
				key := j % 100
				m.Get(key)
				m.Contains(key)
			}
		}()
	}

	for i := 0; i < numWriters; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for j := 0; j < opsPerGoroutine; j++ {
				key := 100 + id*opsPerGoroutine + j
				m.Set(key, key*2)
			}
		}(i)
	}

	wg.Wait()

	expectedLen := 100 + numWriters*opsPerGoroutine
	assert.Equal(t, expectedLen, m.Len())
}

// TestConcurrentRemoves concurrent removes are safe.
func TestConcurrentRemoves(t *testing.T) {
	m := New[int, int]()

	for i := 0; i < 100; i++ {
		m.Set(i, i*2)
	}

	numGoroutines := 5
	keysPerGoroutine := 20
	var wg sync.WaitGroup
	wg.Add(numGoroutines)

	for i := 0; i < numGoroutines; i++ {
		go func(id int) {
			defer wg.Done()
			for j := 0; j < keysPerGoroutine; j++ {
				key := id*keysPerGoroutine + j
				m.Remove(key)
			}
		}(i)
	}

	wg.Wait()

	expectedLen := 100 - (numGoroutines * keysPerGoroutine)
	assert.Equal(t, expectedLen, m.Len())
}

// TestConcurrentIterators multiple iterators work concurrently.
func TestConcurrentIterators(t *testing.T) {
	m := New[int, int]()

	for i := 0; i < 50; i++ {
		m.Set(i, i*2)
	}

	numIterators := 5
	var wg sync.WaitGroup
	var successCount atomic.Int32

	for i := 0; i < numIterators; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			count := 0
			it := m.Iterator()
			for {
				kvp := it.Next()
				if kvp == nil {
					break
				}
				count++
			}
			if count == 50 {
				successCount.Add(1)
			}
		}()
	}

	wg.Wait()

	assert.Equal(t, int32(numIterators), successCount.Load())
}

// TestIteratorWithConcurrentRemovals iterator with concurrent removals.
func TestIteratorWithConcurrentRemovals(t *testing.T) {
	m := New[int, int]()

	for i := 0; i < 100; i++ {
		m.Set(i, i*2)
	}

	var wg sync.WaitGroup

	wg.Add(1)
	var iteratedKeys []int
	go func() {
		defer wg.Done()
		it := m.Iterator()
		for {
			kvp := it.Next()
			if kvp == nil {
				break
			}
			iteratedKeys = append(iteratedKeys, kvp.Key)
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		for i := 0; i < 100; i += 2 {
			m.Remove(i)
		}
	}()

	wg.Wait()

	// note: iterator snapshots keys at creation time
	assert.NotEmpty(t, iteratedKeys)
}

// TestConcurrentMoveOperations concurrent move operations are safe.
func TestConcurrentMoveOperations(t *testing.T) {
	m := New[int, int]()

	for i := 0; i < 50; i++ {
		m.Set(i, i*2)
	}

	numGoroutines := 5
	var wg sync.WaitGroup
	wg.Add(numGoroutines)

	for g := 0; g < numGoroutines; g++ {
		go func(id int) {
			defer wg.Done()
			for i := 0; i < 10; i++ {
				key := (id * 10) + i
				if key < 50 {
					if i%2 == 0 {
						m.MoveToFront(key)
					} else {
						m.MoveToBack(key)
					}
				}
			}
		}(g)
	}

	wg.Wait()

	assert.Equal(t, 50, m.Len())
}

// TestChainedOperations method chaining.
func TestChainedOperations(t *testing.T) {
	m := New[string, int]()

	m.Set("a", 1).Set("b", 2).Set("c", 3)

	assert.Equal(t, 3, m.Len())

	val, ok := m.Get("b")
	assert.True(t, ok)
	assert.Equal(t, 2, *val)
}

// BenchmarkSet benchmarks Set operations.
func BenchmarkSet(b *testing.B) {
	m := New[int, int]()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		m.Set(i, i*2)
	}
}

// BenchmarkGet benchmarks Get operations.
func BenchmarkGet(b *testing.B) {
	m := New[int, int]()
	for i := 0; i < 1000; i++ {
		m.Set(i, i*2)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		m.Get(i % 1000)
	}
}

// BenchmarkConcurrentSetGet benchmarks concurrent Set and Get operations.
func BenchmarkConcurrentSetGet(b *testing.B) {
	m := New[int, int]()
	numGoroutines := 4
	opsPerGoroutine := b.N / numGoroutines

	b.ResetTimer()
	var wg sync.WaitGroup
	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for j := 0; j < opsPerGoroutine; j++ {
				key := id*opsPerGoroutine + j
				m.Set(key, key*2)
				m.Get(key)
			}
		}(i)
	}
	wg.Wait()
}
