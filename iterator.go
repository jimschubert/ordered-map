package orderedmap

import "github.com/jimschubert/ordered-map/internal/list"

// Iterator allows iteration of an OrderedMap
type Iterator[K comparable, V any] struct {
	orderedMap *OrderedMap[K, V]
	pos        *list.Element[*KeyValuePair[K, V]]
}

// Next returns the next KeyValuePair, or nil if there are no more items
func (i *Iterator[K, V]) Next() *KeyValuePair[K, V] {
	if i.pos == nil {
		return nil
	}
	var value *KeyValuePair[K, V]
	if i.pos.Value != nil {
		value = i.pos.Value
		i.pos = i.pos.Next()
	}
	return value
}
