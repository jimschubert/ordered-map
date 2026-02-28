// Package list provides a minimal generic wrapper around container/list.
// It delegates to the standard library and performs type assertions as needed.
package list

import "container/list"

// List is a generic doubly-linked list backed by container/list.
type List[T any] struct {
	inner *list.List
}

// Element is an element in a generic list.
type Element[T any] struct {
	inner *list.Element
}

// New returns an initialized list.
func New[T any]() *List[T] {
	return &List[T]{inner: list.New()}
}

// Init initializes or clears the list. The zero value is ready for Init.
func (l *List[T]) Init() *List[T] {
	if l.inner == nil {
		l.inner = list.New()
	} else {
		l.inner.Init()
	}
	return l
}

// Len returns the number of elements.
func (l *List[T]) Len() int {
	if l.inner == nil {
		return 0
	}
	return l.inner.Len()
}

// Front returns the first element, or nil.
func (l *List[T]) Front() *Element[T] {
	if l.inner == nil {
		return nil
	}
	return l.wrap(l.inner.Front())
}

// Back returns the last element, or nil.
func (l *List[T]) Back() *Element[T] {
	if l.inner == nil {
		return nil
	}
	return l.wrap(l.inner.Back())
}

// PushBack inserts v at the back and returns the new element.
func (l *List[T]) PushBack(v T) *Element[T] {
	l.lazyInit()
	return l.wrap(l.inner.PushBack(v))
}

// Remove removes e and returns its value.
func (l *List[T]) Remove(e *Element[T]) T {
	if e == nil || e.inner == nil {
		var zero T
		return zero
	}
	v := l.inner.Remove(e.inner)
	if v == nil {
		var zero T
		return zero
	}
	return v.(T)
}

// MoveToFront moves e to the front.
func (l *List[T]) MoveToFront(e *Element[T]) {
	if e != nil && e.inner != nil && l.inner != nil {
		l.inner.MoveToFront(e.inner)
	}
}

// MoveToBack moves e to the back.
func (l *List[T]) MoveToBack(e *Element[T]) {
	if e != nil && e.inner != nil && l.inner != nil {
		l.inner.MoveToBack(e.inner)
	}
}

// MoveAfter moves e to after mark.
func (l *List[T]) MoveAfter(e, mark *Element[T]) {
	if e != nil && mark != nil && l.inner != nil {
		l.inner.MoveAfter(e.inner, mark.inner)
	}
}

// MoveBefore moves e to before mark.
func (l *List[T]) MoveBefore(e, mark *Element[T]) {
	if e != nil && mark != nil && l.inner != nil {
		l.inner.MoveBefore(e.inner, mark.inner)
	}
}

// Next returns the next list element, or nil.
func (e *Element[T]) Next() *Element[T] {
	if e == nil || e.inner == nil {
		return nil
	}
	return wrapElement[T](e.inner.Next())
}

// Value returns the value stored in e.
func (e *Element[T]) Value() T {
	if e == nil || e.inner == nil || e.inner.Value == nil {
		var empty T
		return empty
	}
	return e.inner.Value.(T)
}

func (l *List[T]) lazyInit() {
	if l.inner == nil {
		l.inner = list.New()
	}
}

func (l *List[T]) wrap(e *list.Element) *Element[T] {
	return wrapElement[T](e)
}

func wrapElement[T any](e *list.Element) *Element[T] {
	if e == nil {
		return nil
	}
	return &Element[T]{inner: e}
}
