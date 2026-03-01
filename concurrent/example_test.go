package concurrent

import (
	"fmt"
	"sync"
)

// Example_basicConcurrentUsage demonstrates basic usage of the concurrent OrderedMap.
func Example_basicConcurrentUsage() {
	// Create a concurrent OrderedMap
	m := New[string, string]()

	// Safe to use from multiple goroutines
	m.Set("First", "1st")
	m.Set("Second", "2nd")
	m.Set("Third", "3rd")

	// Retrieve values
	value, ok := m.Get("Second")
	if ok {
		fmt.Printf("Value at 'Second': %s\n", *value)
	}

	// Iterate safely
	it := m.Iterator()
	for {
		kvp := it.Next()
		if kvp == nil {
			break
		}
		fmt.Printf("%s=%s ", kvp.Key, kvp.Value)
	}
	fmt.Println()

	// Output:
	// Value at 'Second': 2nd
	// First=1st Second=2nd Third=3rd
}

// Example_concurrentReadsAndWrites demonstrates concurrent reads and writes.
func Example_concurrentReadsAndWrites() {
	m := New[int, int]()

	var wg sync.WaitGroup

	// Writers
	wg.Add(2)
	go func() {
		defer wg.Done()
		for i := 0; i < 5; i++ {
			m.Set(i, i*10)
		}
	}()

	go func() {
		defer wg.Done()
		for i := 5; i < 10; i++ {
			m.Set(i, i*10)
		}
	}()

	// Readers
	wg.Add(2)
	go func() {
		defer wg.Done()
		for i := 0; i < 10; i++ {
			if v, ok := m.Get(i); ok {
				_ = v // use value
			}
		}
	}()

	go func() {
		defer wg.Done()
		for i := 0; i < 10; i++ {
			_ = m.Contains(i)
		}
	}()

	wg.Wait()
	fmt.Printf("Final map size: %d\n", m.Len())
	// Output:
	// Final map size: 10
}

// Example_iteratorSnapshot demonstrates that the iterator captures a snapshot.
func Example_iteratorSnapshot() {
	m := New[int, int]()

	// Populate
	for i := 0; i < 3; i++ {
		m.Set(i, i*10)
	}

	// Create iterator - captures snapshot at this moment
	it := m.Iterator()

	// Remove an item after iterator is created
	m.Remove(1)

	// Iterator still iterates the original snapshot
	var keys []int
	for {
		kvp := it.Next()
		if kvp == nil {
			break
		}
		keys = append(keys, kvp.Key)
	}

	for i, k := range keys {
		if i > 0 {
			fmt.Print(" ")
		}
		fmt.Print(k)
	}
	fmt.Println()

	fmt.Printf("Map now has %d items, iterator saw snapshot\n", m.Len())
	// Output:
	// 0 2
	// Map now has 2 items, iterator saw snapshot
}

