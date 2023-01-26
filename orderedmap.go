package orderedmap

type KeyValuePair[K comparable, V any] struct {
	Key   K
	Value V
}

type OrderedMap[K comparable, V any] struct {
	items map[K]*KeyValuePair[K, V]
}

func New[K comparable, V any]() *OrderedMap[K, V] {
	m := OrderedMap[K, V]{}
	return &m
}
