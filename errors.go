package orderedmap

import "fmt"

// KeyNotFoundError conveys to the caller that a key was requested but not found in the map
type KeyNotFoundError[K comparable] struct {
	Key K
}

// Error provides a string representation of this error.
func (k *KeyNotFoundError[K]) Error() string {
	return fmt.Sprintf("key not found: %v", k.Key)
}

func keyNotFound[K comparable](key K) *KeyNotFoundError[K] {
	return &KeyNotFoundError[K]{Key: key}
}

// DuplicateKeyValueError conveys to the caller that a key and value were requested to be inserted via manipulation
// function such as InsertBefore or InsertAfter, but the key already existed in the map.
// This is raised to avoid unintentional modification of Value.
// If the caller intends to modify Value for an existing key, call Set, followed by one of the Insert or Move functions.
type DuplicateKeyValueError[K comparable, V any] struct {
	Key   K
	Value V
}

// Error provides a string representation of this error.
func (k *DuplicateKeyValueError[K, V]) Error() string {
	return fmt.Sprintf("key %v already exists with value %v", k.Key, k.Value)
}

func duplicateValue[K comparable, V any](key K, value V) *DuplicateKeyValueError[K, V] {
	return &DuplicateKeyValueError[K, V]{
		Key:   key,
		Value: value,
	}
}
