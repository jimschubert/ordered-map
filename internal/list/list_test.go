package list

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	l := New[int]()
	assert.NotNil(t, l)
	assert.Equal(t, 0, l.Len())
	assert.Nil(t, l.Front())
	assert.Nil(t, l.Back())
}

func TestZeroValue(t *testing.T) {
	var l List[int]
	assert.Equal(t, 0, l.Len())
	assert.Nil(t, l.Front())
	assert.Nil(t, l.Back())
	l.Init()
	if l.inner == nil {
		t.Fatal("Init() did not initialize inner list")
	}
	e := l.PushBack(1)
	assert.NotNil(t, e)
	assert.Equal(t, 1, e.Value())
	assert.Equal(t, 1, l.Len())
}

func TestPushBackAndIterate(t *testing.T) {
	l := New[int]()
	e1 := l.PushBack(1)
	e2 := l.PushBack(2)
	e3 := l.PushBack(3)
	e4 := l.PushBack(4)

	assert.Equal(t, 4, l.Len())
	assert.Equal(t, 1, e1.Value())
	assert.Equal(t, 2, e2.Value())
	assert.Equal(t, 3, e3.Value())
	assert.Equal(t, 4, e4.Value())

	var got []int
	for e := l.Front(); e != nil; e = e.Next() {
		got = append(got, e.Value())
	}
	assert.Equal(t, []int{1, 2, 3, 4}, got)
}

func TestInitClears(t *testing.T) {
	l := New[string]()
	l.PushBack("a")
	l.PushBack("b")
	l.PushBack("c")
	assert.Equal(t, 3, l.Len())
	l.Init()
	assert.Equal(t, 0, l.Len())
	assert.Nil(t, l.Front())
	assert.Nil(t, l.Back())
}

func TestRemove(t *testing.T) {
	l := New[int]()
	e1 := l.PushBack(1)
	e2 := l.PushBack(2)
	e3 := l.PushBack(3)

	assert.Equal(t, 2, l.Remove(e2))
	assert.Equal(t, 2, l.Len())
	expectOrder(t, l, []int{1, 3})

	l.Remove(e1)
	assert.Equal(t, 1, l.Len())
	expectOrder(t, l, []int{3})

	l.Remove(e3)
	assert.Equal(t, 0, l.Len())
	assert.Nil(t, l.Front())
	assert.Nil(t, l.Back())
}

func TestRemoveReturnsValue(t *testing.T) {
	l := New[string]()
	e := l.PushBack("hello")
	assert.Equal(t, "hello", l.Remove(e))
}

func TestMoveToFront(t *testing.T) {
	l := New[int]()
	l.PushBack(1)
	l.PushBack(2)
	e3 := l.PushBack(3)
	l.PushBack(4)

	l.MoveToFront(e3)
	expectOrder(t, l, []int{3, 1, 2, 4})

	e4 := l.Back()
	l.MoveToFront(e4)
	expectOrder(t, l, []int{4, 3, 1, 2})
}

func TestMoveToBack(t *testing.T) {
	l := New[int]()
	e1 := l.PushBack(1)
	l.PushBack(2)
	l.PushBack(3)
	l.PushBack(4)

	l.MoveToBack(e1)
	expectOrder(t, l, []int{2, 3, 4, 1})
}

func TestMoveAfter(t *testing.T) {
	l := New[int]()
	e1 := l.PushBack(1)
	e2 := l.PushBack(2)
	e3 := l.PushBack(3)
	_ = l.PushBack(4)

	l.MoveAfter(e3, e2) // move 3 after 2 -> 1, 2, 3, 4 (no change)
	expectOrder(t, l, []int{1, 2, 3, 4})

	l.MoveAfter(e1, e3) // move 1 after 3 -> 2, 3, 1, 4
	expectOrder(t, l, []int{2, 3, 1, 4})
}

func TestMoveBefore(t *testing.T) {
	l := New[int]()
	e1 := l.PushBack(1)
	e2 := l.PushBack(2)
	e3 := l.PushBack(3)
	e4 := l.PushBack(4)

	l.MoveBefore(e4, e2) // move 4 before 2 -> 1, 4, 2, 3
	expectOrder(t, l, []int{1, 4, 2, 3})

	l.MoveBefore(e3, e1) // move 3 before 1 -> 3, 1, 4, 2
	expectOrder(t, l, []int{3, 1, 4, 2})
}

func TestElementValueNilSafe(t *testing.T) {
	var e *Element[int]
	assert.Equal(t, 0, e.Value())
	assert.Nil(t, e.Next())
}

func TestRemoveNilElement(t *testing.T) {
	l := New[int]()
	l.PushBack(1)
	assert.Equal(t, 0, l.Remove(nil))
	assert.Equal(t, 1, l.Len())
}

func expectOrder[T comparable](t *testing.T, l *List[T], want []T) {
	t.Helper()
	var got []T
	for e := l.Front(); e != nil; e = e.Next() {
		got = append(got, e.Value())
	}
	assert.Equal(t, want, got)
}
