// Package concurrent provides a thread-safe wrapper around OrderedMap.
//
// The concurrent package exports an OrderedMap type that wraps the main orderedmap.OrderedMap
// with a sync.RWMutex to make it safe for concurrent use across multiple goroutines.
//
// # Usage
//
// To use the concurrent OrderedMap, simply import the concurrent package and use New()
// to create an instance:
//
//	import "github.com/jimschubert/ordered-map/concurrent"
//
//	// Create a new concurrent map
//	m := concurrent.New[string, int]()
//
//	// Use it safely from multiple goroutines
//	m.Set("key", 42)
//	value, ok := m.Get("key")
//
// # Key Differences from orderedmap.OrderedMap
//
// - All operations are protected by an internal RWMutex
// - Read operations (Get, Contains, Len, First, Last, Keys, Iterator) use RLock()
// - Write operations (Set, Remove, MoveToFront, MoveToBack, MoveAfter, MoveBefore, InsertAfter, InsertBefore) use Lock()
// - The Iterator() method returns a snapshot of keys at the time it was called, allowing safe concurrent iteration
//   even if the map is modified during iteration
//
// # Performance Considerations
//
// The concurrent version adds lock overhead to every operation. If you:
// - Only use the map from a single goroutine, use the non-concurrent orderedmap.OrderedMap
// - Use the map from multiple goroutines, use this concurrent package
//
// The lock is held only for the duration of the operation, not during subsequent iterator work.
//
// # Iterator Semantics
//
// Iterators snapshot the keys at the time Iterator() is called. If keys are removed from the map
// during iteration, the iterator will skip those entries gracefully. This provides safe concurrent
// iteration without holding locks for extended periods.
package concurrent

