package concurrent

import (
	"github.com/jimschubert/ordered-map"
)

// Iterator allows for walking the map's contents in-order.
// The iterator holds a snapshot of keys from when Iterator() was called,
// allowing safe iteration even if the map is modified concurrently.
type Iterator[K comparable, V any] struct {
	om   *OrderedMap[K, V]
	keys []K
	idx  int
}

// Next returns the next KeyValuePair in the map, or nil if the iteration is complete.
// The returned KeyValuePair contains the current key and its value at the time Next() is called.
// If a key has been removed from the map since Iterator() was called, Next() will return nil for that entry.
func (it *Iterator[K, V]) Next() *orderedmap.KeyValuePair[K, V] {
	if it.idx >= len(it.keys) {
		return nil
	}

	key := it.keys[it.idx]
	it.idx++

	value, ok := it.om.Get(key)
	if !ok {
		// key was removed, skip
		return it.Next()
	}

	return &orderedmap.KeyValuePair[K, V]{
		Key:   key,
		Value: *value,
	}
}

