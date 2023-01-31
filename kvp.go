package orderedmap

import (
	"fmt"

	"github.com/jimschubert/ordered-map/internal/list"
)

// KeyValuePair holds the ordered map pair represented by Key and Value
type KeyValuePair[K comparable, V any] struct {
	Key     K
	Value   V
	element *list.Element[*KeyValuePair[K, V]]
}

// String representation of this KeyValuePair
func (k *KeyValuePair[K, V]) String() string {
	return fmt.Sprintf("%v=%+v", k.Key, k.Value)
}
