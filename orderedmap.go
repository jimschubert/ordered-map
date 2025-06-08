package orderedmap

import (
	"bytes"
	"fmt"

	"github.com/jimschubert/ordered-map/internal/list"
)

// OrderedMap is a data structure relating items keyed by K to values of type V.
// The insertion order is maintained and iterable.
//
// Operations to manipulate the order are exposed which mirror the API of stdlib's List.
//
// NOTE: This map maintains ordering, _not_ sorting.
type OrderedMap[K comparable, V any] struct {
	items map[K]*KeyValuePair[K, V]
	order list.List[*KeyValuePair[K, V]]
}

// Init initializes or clears ordered map o.
func (o *OrderedMap[K, V]) Init() *OrderedMap[K, V] {
	o.items = make(map[K]*KeyValuePair[K, V])
	o.order.Init()
	return o
}

func (o *OrderedMap[K, V]) insertKeyValuePair(key K, value V) *KeyValuePair[K, V] {
	pair := KeyValuePair[K, V]{Key: key, Value: value}
	element := o.order.PushBack(&pair)
	o.items[key] = &pair
	pair.element = element
	return &pair
}

// Set a key of type K to a value of type V. If the key exists, the value will be modified.
func (o *OrderedMap[K, V]) Set(key K, value V) *OrderedMap[K, V] {
	if existing, ok := o.items[key]; ok {
		existing.Value = value
		return o
	}

	_ = o.insertKeyValuePair(key, value)
	return o
}

// Get the value stored at the key.
func (o *OrderedMap[K, V]) Get(key K) (*V, bool) {
	if existing, ok := o.items[key]; ok {
		value := existing.Value
		return &value, true
	}

	return nil, false
}

// GetOrDefault either gets the value stored at key or returns the default value defined by defaultValue
func (o *OrderedMap[K, V]) GetOrDefault(key K, defaultValue V) V {
	value, ok := o.Get(key)
	if value == nil || !ok {
		return defaultValue
	}

	return *value
}

// Remove the key (and value) from the map.
// Returns the removed value and true if the value has been removed.
// Returns nil and false if the item did not exist in the map.
func (o *OrderedMap[K, V]) Remove(key K) (*KeyValuePair[K, V], bool) {
	if kvp, ok := o.items[key]; ok {
		delete(o.items, key)
		o.order.Remove(kvp.element)
		return kvp, true
	}

	return nil, false
}

// First returns the first KeyValuePair contained in the map, or nil.
func (o *OrderedMap[K, V]) First() *KeyValuePair[K, V] {
	front := o.order.Front()
	if front == nil {
		return nil
	}
	return front.Value
}

// Last returns the last KeyValuePair contained in the map, or nil.
func (o *OrderedMap[K, V]) Last() *KeyValuePair[K, V] {
	last := o.order.Back()
	if last == nil {
		return nil
	}
	return last.Value
}

// Iterator returns an initialized *Iterator[K, V] for walking the map's contents in-order.
func (o *OrderedMap[K, V]) Iterator() *Iterator[K, V] {
	return &Iterator[K, V]{
		pos:        o.order.Front(),
		orderedMap: o,
	}
}

// Keys returns the ordered slice of keys for this map
func (o *OrderedMap[K, V]) Keys() []K {
	keys := make([]K, 0)
	it := o.Iterator()
	var kvp *KeyValuePair[K, V]
	for {
		kvp = it.Next()
		if kvp == nil {
			break
		}
		keys = append(keys, kvp.Key)
	}
	return keys
}

// MoveToFront allows for manipulating the order of a map by moving key (and associated value) to the front of the map.
//
// If key does not exist in the map, this will raise a KeyNotFoundError to signal failed intent to the caller.
//
// If key does not exist, the map is unmodified.
func (o *OrderedMap[K, V]) MoveToFront(key K) error {
	if element, ok := o.items[key]; ok {
		o.order.MoveToFront(element.element)
		return nil
	}
	return keyNotFound(key)
}

// MoveToBack allows for manipulating the order of a map by moving key (and associated value) to the back of the map.
//
// If key does not exist in the map, this will raise a KeyNotFoundError to signal failed intent to the caller.
//
// If key does not exist, the map is unmodified.
func (o *OrderedMap[K, V]) MoveToBack(key K) error {
	if element, ok := o.items[key]; ok {
		o.order.MoveToBack(element.element)
		return nil
	}
	return keyNotFound(key)
}

// MoveAfter allows for manipulating the order of a map by moving the pair defined at 'key' after the pair defined at 'after'.
//
// If either element is not found, this will raise a KeyNotFoundError to signal failed intent to the caller.
//
// This differs from behavior one might expect from container/list in the standard library, because we operate on
// user defined types rather than directly on an element of the ordered list of KeyValuePair.
func (o *OrderedMap[K, V]) MoveAfter(key, after K) error {
	if element, ok := o.items[key]; ok {
		if mark, exists := o.items[after]; exists {
			o.order.MoveAfter(element.element, mark.element)
			return nil
		}

		return keyNotFound(after)
	}

	return keyNotFound(key)
}

// MoveBefore allows for manipulating the order of a map by moving the pair defined at 'key' before the pair defined at 'before'.
//
// If either element is not found, this will raise a KeyNotFoundError to signal failed intent to the caller.
//
// This differs from behavior one might expect from container/list in the standard library, because we operate on
// user defined types rather than directly on an element of the ordered list of KeyValuePair.
func (o *OrderedMap[K, V]) MoveBefore(key, before K) error {
	if element, ok := o.items[key]; ok {
		if mark, exists := o.items[before]; exists {
			o.order.MoveBefore(element.element, mark.element)
			return nil
		}

		return keyNotFound(before)
	}

	return keyNotFound(key)
}

// InsertAfter allows for manipulating the order of a map by inserting the provided key and value after the pair defined at 'after'.
//
// If either element is not found, this will raise a KeyNotFoundError to signal failed intent to the caller.
// If key and after are the same or if key already exists, this will raise a DuplicateKeyValueError.
//
// This differs from behavior one might expect from container/list in the standard library, because we operate on
// user defined types rather than directly on an element of the ordered list of KeyValuePair.
func (o *OrderedMap[K, V]) InsertAfter(key K, value V, after K) error {
	if mark, ok := o.items[after]; ok {
		if exists, precondition := o.items[key]; precondition {
			return duplicateValue(exists.Key, exists.Value)
		}
		if key == after {
			return duplicateValue(mark.Key, mark.Value)
		}
		newElement := o.insertKeyValuePair(key, value)
		o.order.MoveAfter(newElement.element, mark.element)
		return nil
	}

	return keyNotFound(key)
}

// InsertBefore allows for manipulating the order of a map by inserting the provided key and value before the pair defined at 'before'.
//
// If either element is not found, this will raise a KeyNotFoundError to signal failed intent to the caller.
// If key and before are the same or if key already exists, this will raise a DuplicateKeyValueError.
//
// This differs from behavior one might expect from container/list in the standard library, because we operate on
// user defined types rather than directly on an element of the ordered list of KeyValuePair.
func (o *OrderedMap[K, V]) InsertBefore(key K, value V, before K) error {
	if mark, ok := o.items[before]; ok {
		if exists, precondition := o.items[key]; precondition {
			return duplicateValue(exists.Key, exists.Value)
		}
		if key == before {
			return duplicateValue(mark.Key, mark.Value)
		}
		newElement := o.insertKeyValuePair(key, value)
		o.order.MoveBefore(newElement.element, mark.element)
		return nil
	}
	return keyNotFound(key)
}

// String fulfills the fmt.Stringer interface
func (o *OrderedMap[K, V]) String() string {
	buf := bytes.Buffer{}
	buf.WriteString(fmt.Sprintf("OrderedMap[%T,%T]", *new(K), *new(V)))
	if o != nil && o.order.Len() > 0 {
		l := o.order
		for e := l.Front(); e != nil; e = e.Next() {
			buf.WriteString(fmt.Sprintf("\t%v=%v,\n", e.Value.Key, e.Value.Value))
		}
	} else {
		buf.WriteString("{}")
	}
	return buf.String()
}

// GoString fulfills the fmt.GoStringer interface and can be coupled with go-cmp for easier diffs.
func (o *OrderedMap[K, V]) GoString() string {
	if o == nil {
		return fmt.Sprintf("(*orderedmap.OrderedMap[%T, %T])(nil)", *new(K), *new(V))
	}
	buf := bytes.Buffer{}
	buf.WriteString(fmt.Sprintf("orderedmap.New[%T,%T]()", *new(K), *new(V)))
	if o != nil && o.order.Len() > 0 {
		buf.WriteString(".\n")
		l := o.order
		for e := l.Front(); e != nil; e = e.Next() {
			buf.WriteString(fmt.Sprintf("\tSet(%#v, %#v)", e.Value.Key, e.Value.Value))
			if e.Next() != nil {
				buf.WriteString(".\n")
			}
		}
	}
	return buf.String()
}

// New initializes a new OrderedMap
func New[K comparable, V any]() *OrderedMap[K, V] {
	m := new(OrderedMap[K, V])
	l := list.New[*KeyValuePair[K, V]]()
	m.order = *l
	m.Init()
	return m
}
