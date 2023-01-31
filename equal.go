package orderedmap

import "reflect"

// Equal is a lock-free evaluation of two OrderedMap values. It is up to the user to
// lock these maps for thread-safe equality check.
//
// This optimizes equality of key/value pairs, ignoring the internals of the data structure.
// If the caller invokes reflect.DeepEqual on equivalent maps, the result should be the same.
// However, reflect.DeepEqual evaluates both exported and unexported fields which unnecessary overhead.
//
// This implementation will incur the overhead of reflect.DeepEqual mentioned above if any key in the OrderedMap refers
// to an OrderedMap value.
func Equal[K comparable, V any](x, y *OrderedMap[K, V]) bool {
	if (x == nil && y != nil) || (y == nil && x != nil) {
		return false
	}
	if x.order.Len() != y.order.Len() {
		return false
	}

	xIt := x.Iterator()
	yIt := y.Iterator()

	var xCurrent *KeyValuePair[K, V]
	var yCurrent *KeyValuePair[K, V]

	for {
		xCurrent = xIt.Next()
		yCurrent = yIt.Next()

		if xCurrent == nil && yCurrent == nil {
			// we've reached the end at the same time without hitting a negative condition
			break
		}

		// one side finished before the other.
		// this can happen if maps were modified after precondition check above.
		if (xCurrent == nil && yCurrent != nil) ||
			(yCurrent == nil && xCurrent != nil) {
			return false
		}

		if xCurrent.Key != yCurrent.Key {
			return false
		}

		if !reflect.DeepEqual(xCurrent.Value, yCurrent.Value) {
			return false
		}
	}

	return true
}
