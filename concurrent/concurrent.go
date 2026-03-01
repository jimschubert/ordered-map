// Package concurrent provides a thread-safe version of OrderedMap.
package concurrent

import (
	"sync"

	"github.com/jimschubert/ordered-map"
)

// OrderedMap is an OrderedMap that can be safely used concurrently
// by multiple goroutines. All operations are protected by an internal mutex.
type OrderedMap[K comparable, V any] struct {
	mu *sync.RWMutex
	om *orderedmap.OrderedMap[K, V]
}

// New initializes a new concurrent OrderedMap.
func New[K comparable, V any]() *OrderedMap[K, V] {
	return &OrderedMap[K, V]{
		mu: &sync.RWMutex{},
		om: orderedmap.New[K, V](),
	}
}

// Set a key of type K to a value of type V. If the key exists, the value will be modified.
// This operation is safe for concurrent use.
func (o *OrderedMap[K, V]) Set(key K, value V) *OrderedMap[K, V] {
	o.mu.Lock()
	defer o.mu.Unlock()
	o.om.Set(key, value)
	return o
}

// Get the value stored at the key.
// This operation is safe for concurrent use.
func (o *OrderedMap[K, V]) Get(key K) (*V, bool) {
	o.mu.RLock()
	defer o.mu.RUnlock()
	return o.om.Get(key)
}

// GetOrDefault returns the value stored at key, or defaultValue if the key is not present.
// When V is a pointer or interface type, a stored nil value is returned as-is; only absence of the key uses the default.
// This operation is safe for concurrent use.
func (o *OrderedMap[K, V]) GetOrDefault(key K, defaultValue V) V {
	o.mu.RLock()
	defer o.mu.RUnlock()
	return o.om.GetOrDefault(key, defaultValue)
}

// Remove the key (and value) from the map.
// Returns the removed value and true if the value has been removed.
// Returns nil and false if the item did not exist in the map.
// This operation is safe for concurrent use.
func (o *OrderedMap[K, V]) Remove(key K) (*orderedmap.KeyValuePair[K, V], bool) {
	o.mu.Lock()
	defer o.mu.Unlock()
	return o.om.Remove(key)
}

// First returns the first KeyValuePair contained in the map, or nil.
// This operation is safe for concurrent use.
func (o *OrderedMap[K, V]) First() *orderedmap.KeyValuePair[K, V] {
	o.mu.RLock()
	defer o.mu.RUnlock()
	return o.om.First()
}

// Last returns the last KeyValuePair contained in the map, or nil.
// This operation is safe for concurrent use.
func (o *OrderedMap[K, V]) Last() *orderedmap.KeyValuePair[K, V] {
	o.mu.RLock()
	defer o.mu.RUnlock()
	return o.om.Last()
}

// Len returns the number of key-value pairs in the map.
// This operation is safe for concurrent use.
func (o *OrderedMap[K, V]) Len() int {
	o.mu.RLock()
	defer o.mu.RUnlock()
	return o.om.Len()
}

// Contains reports whether the map contains key.
// This operation is safe for concurrent use.
func (o *OrderedMap[K, V]) Contains(key K) bool {
	o.mu.RLock()
	defer o.mu.RUnlock()
	return o.om.Contains(key)
}

// Iterator returns an initialized *Iterator[K, V] for walking the map's contents in-order.
// NOTE: The returned iterator holds a snapshot of the keys at the time Iterator() is called.
// This operation is safe for concurrent use, but the iterator should not be shared across goroutines.
func (o *OrderedMap[K, V]) Iterator() *Iterator[K, V] {
	o.mu.RLock()
	defer o.mu.RUnlock()

	// snapshot keys to avoid lock during iteration
	keys := o.om.Keys()
	return &Iterator[K, V]{
		om:   o,
		keys: keys,
		idx:  0,
	}
}

// Keys returns the ordered slice of keys for this map.
// This operation is safe for concurrent use.
func (o *OrderedMap[K, V]) Keys() []K {
	o.mu.RLock()
	defer o.mu.RUnlock()
	return o.om.Keys()
}

// MoveToFront moves key (and associated value) to the front of the map.
// If key does not exist in the map, this will return an error.
// This operation is safe for concurrent use.
func (o *OrderedMap[K, V]) MoveToFront(key K) error {
	o.mu.Lock()
	defer o.mu.Unlock()
	return o.om.MoveToFront(key)
}

// MoveToBack moves key (and associated value) to the back of the map.
// If key does not exist in the map, this will return an error.
// This operation is safe for concurrent use.
func (o *OrderedMap[K, V]) MoveToBack(key K) error {
	o.mu.Lock()
	defer o.mu.Unlock()
	return o.om.MoveToBack(key)
}

// MoveAfter moves the pair defined at 'key' after the pair defined at 'after'.
// If either element is not found, this will return an error.
// This operation is safe for concurrent use.
func (o *OrderedMap[K, V]) MoveAfter(key, after K) error {
	o.mu.Lock()
	defer o.mu.Unlock()
	return o.om.MoveAfter(key, after)
}

// MoveBefore moves the pair defined at 'key' before the pair defined at 'before'.
// If either element is not found, this will return an error.
// This operation is safe for concurrent use.
func (o *OrderedMap[K, V]) MoveBefore(key, before K) error {
	o.mu.Lock()
	defer o.mu.Unlock()
	return o.om.MoveBefore(key, before)
}

// InsertAfter inserts the provided key and value after the pair defined at 'after'.
// If 'after' is not found, or if key already exists, this will return an error.
// This operation is safe for concurrent use.
func (o *OrderedMap[K, V]) InsertAfter(key K, value V, after K) error {
	o.mu.Lock()
	defer o.mu.Unlock()
	return o.om.InsertAfter(key, value, after)
}

// InsertBefore inserts the provided key and value before the pair defined at 'before'.
// If 'before' is not found, or if key already exists, this will return an error.
// This operation is safe for concurrent use.
func (o *OrderedMap[K, V]) InsertBefore(key K, value V, before K) error {
	o.mu.Lock()
	defer o.mu.Unlock()
	return o.om.InsertBefore(key, value, before)
}

// String fulfills the fmt.Stringer interface.
// This operation is safe for concurrent use.
func (o *OrderedMap[K, V]) String() string {
	o.mu.RLock()
	defer o.mu.RUnlock()
	return o.om.String()
}

// GoString fulfills the fmt.GoStringer interface.
// This operation is safe for concurrent use.
func (o *OrderedMap[K, V]) GoString() string {
	o.mu.RLock()
	defer o.mu.RUnlock()
	return o.om.GoString()
}

