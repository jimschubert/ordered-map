package orderedmap

import (
	"reflect"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func ptr[K any](input K) *K {
	return &input
}

type pair[K comparable, V any] struct {
	Key   K
	Value V
}

func (k *pair[K, V]) Equals(other *KeyValuePair[K, V]) bool {
	if k == nil && other == nil {
		return true
	}
	if k == nil && other != nil {
		return false
	}
	if other == nil && k != nil {
		return false
	}
	if (*k).Key != (*other).Key {
		return false
	}

	var kValue, otherValue any
	kValue = (*k).Value
	otherValue = (*other).Value
	if kValue == nil && otherValue == nil {
		return true
	}
	if kValue == nil && otherValue != nil {
		return false
	}
	if kValue != nil && otherValue == nil {
		return false
	}
	return reflect.DeepEqual(kValue, otherValue)
}

func newFromPairs[K comparable, V any](item ...*pair[K, V]) *OrderedMap[K, V] {
	m := New[K, V]()
	for _, k := range item {
		m.Set(k.Key, k.Value)
	}
	return m
}

func kvp[K comparable, V any](key K, value V) *pair[K, V] {
	return &pair[K, V]{Key: key, Value: value}
}

func TestNew(t *testing.T) {
	type testCase[K comparable, V any] struct {
		name string
		want *OrderedMap[K, V]
	}
	tests := []testCase[string, int]{
		{
			name: "initializes an empty map",
			want: new(OrderedMap[string, int]).Init(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := New[string, int](); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("New() = %#v, want %#v", got, tt.want)
			}
		})
	}
}

func TestOrderedMap_First(t *testing.T) {

	type testCase[K comparable, V any] struct {
		name string
		o    OrderedMap[K, V]
		want *pair[K, V]
	}
	tests := []testCase[string, int]{
		{
			name: "First is nil for empty map",
			o:    *New[string, int](),
			want: nil,
		},
		{
			name: "First is first element in single element map",
			o:    *newFromPairs[string, int](kvp("First", 1)),
			want: kvp("First", 1),
		},
		{
			name: "First is first element in multiple element map",
			o:    *newFromPairs[string, int](kvp("Z", 1), kvp("A", 2)),
			want: kvp("Z", 1),
		},
		{
			name: "First is first element in a manipulated map",
			o: func() OrderedMap[string, int] {
				original := newFromPairs[string, int](kvp("Z", 1), kvp("A", 2))
				_ = original.MoveBefore("A", "Z")
				return *original
			}(),
			want: kvp("A", 2),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.o.First()
			if (tt.want == nil && got != nil) || !tt.want.Equals(got) {
				t.Errorf("First() = %s, want %+v", got, tt.want)
			}
		})
	}
}

func TestOrderedMap_Get(t *testing.T) {
	type args[K comparable] struct {
		key K
	}
	type testCase[K comparable, V any] struct {
		name   string
		o      OrderedMap[K, V]
		args   args[K]
		want   *V
		wantOk bool
	}
	tests := []testCase[string, int]{
		{
			name:   "Get is nil on empty map",
			o:      *New[string, int](),
			args:   args[string]{key: "a"},
			want:   nil,
			wantOk: false,
		},
		{
			name:   "Get expected value on single entry map",
			o:      *newFromPairs[string, int](kvp("a", 1)),
			args:   args[string]{key: "a"},
			want:   ptr(1),
			wantOk: true,
		},
		{
			name: "Get expected value on multiple entry map",
			o: *newFromPairs[string, int](
				kvp("a", 1),
				kvp("cat", 10),
				kvp("z", 1000),
				kvp("p", 2134),
				kvp("dog", 0),
			),
			args:   args[string]{key: "p"},
			want:   ptr(2134),
			wantOk: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, ok := tt.o.Get(tt.args.key)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Get() got = %#v, want %#v", got, tt.want)
			}
			if ok != tt.wantOk {
				t.Errorf("Get() ok = %#v, want %v", ok, tt.wantOk)
			}
		})
	}
}

func TestOrderedMap_GetOrDefault(t *testing.T) {
	type args[K comparable, V any] struct {
		key          K
		defaultValue V
	}
	type testCase[K comparable, V any] struct {
		name string
		o    *OrderedMap[K, V]
		args args[K, V]
		want V
	}
	tests := []testCase[string, string]{
		{
			name: "Provides a default value if key not found in empty map",
			o:    New[string, string](),
			args: args[string, string]{
				key:          "first",
				defaultValue: "1st",
			},
			want: "1st",
		},
		{
			name: "Provides a default value if key found in single element map",
			o:    newFromPairs(kvp("first", "1st")),
			args: args[string, string]{
				key:          "first",
				defaultValue: "not 1st",
			},
			want: "1st",
		},
		{
			name: "Provides a default value if key found in multiple element map",
			o:    newFromPairs(kvp("first", "1st"), kvp("second", "2nd"), kvp("third", "3rd")),
			args: args[string, string]{
				key:          "second",
				defaultValue: "not 2nd",
			},
			want: "2nd",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.o.GetOrDefault(tt.args.key, tt.args.defaultValue); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetOrDefault() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestOrderedMap_Init(t *testing.T) {
	type testCase[K comparable, V any] struct {
		name string
		o    *OrderedMap[K, V]
		want *OrderedMap[K, V]
	}
	tests := []testCase[string, string]{
		{
			name: "Init clears/re-initializes an ordered map",
			o:    newFromPairs(kvp("first", "1st"), kvp("second", "2nd"), kvp("third", "3rd")),
			want: New[string, string](),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.o.Init(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Init() = %#v, want %#v", got, tt.want)
			}
		})
	}
}

func TestOrderedMap_InsertAfter(t *testing.T) {
	type testCase[K comparable, V any] struct {
		name    string
		o       *OrderedMap[K, V]
		key     K
		value   V
		after   K
		wantErr bool
		expect  *OrderedMap[K, V]
	}
	tests := []testCase[string, string]{
		{
			name:    "inserts into correct location of multiple element map",
			o:       newFromPairs(kvp("first", "1st"), kvp("second", "2nd"), kvp("fourth", "4th"), kvp("fifth", "5th")),
			key:     "third",
			value:   "3rd",
			after:   "second",
			wantErr: false,
			expect:  newFromPairs(kvp("first", "1st"), kvp("second", "2nd"), kvp("third", "3rd"), kvp("fourth", "4th"), kvp("fifth", "5th")),
		},
		{
			name:    "errors if target key and value already exist",
			o:       newFromPairs(kvp("first", "1st"), kvp("second", "2nd"), kvp("fourth", "4th"), kvp("fifth", "5th")),
			key:     "first",
			value:   "1st",
			after:   "second",
			wantErr: true,
			expect:  newFromPairs(kvp("first", "1st"), kvp("second", "2nd"), kvp("fourth", "4th"), kvp("fifth", "5th")),
		},
		{
			name:    "errors if after key is not found in empty map",
			o:       New[string, string](),
			key:     "third",
			value:   "3rd",
			after:   "alphabet",
			wantErr: true,
			expect:  New[string, string](),
		},
		{
			name:    "errors if after key is not found in populated map",
			o:       newFromPairs(kvp("first", "1st"), kvp("second", "2nd"), kvp("fourth", "4th"), kvp("fifth", "5th")),
			key:     "third",
			value:   "3rd",
			after:   "alphabet",
			wantErr: true,
			expect:  newFromPairs(kvp("first", "1st"), kvp("second", "2nd"), kvp("fourth", "4th"), kvp("fifth", "5th")),
		},
		{
			name:    "errors if 'key' and 'after' are the same",
			o:       newFromPairs(kvp("first", "1st"), kvp("second", "2nd"), kvp("fourth", "4th"), kvp("fifth", "5th")),
			key:     "third",
			value:   "3rd",
			after:   "third",
			wantErr: true,
			expect:  newFromPairs(kvp("first", "1st"), kvp("second", "2nd"), kvp("fourth", "4th"), kvp("fifth", "5th")),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.o.InsertAfter(tt.key, tt.value, tt.after)
			if (err != nil) != tt.wantErr {
				t.Errorf("InsertAfter() error = %s, wantErr %v", err.Error(), tt.wantErr)
			}

			if err != nil {
				t.Logf("InsertAfter() error was: %s", err.Error())
			}

			if !Equal(tt.expect, tt.o) {
				diff := cmp.Diff(tt.expect.GoString(), tt.o.GoString())
				t.Errorf("InsertAfter() expected state mismatch:\n%s", diff)
			}
		})
	}
}

func TestOrderedMap_InsertBefore(t *testing.T) {
	type testCase[K comparable, V any] struct {
		name    string
		o       *OrderedMap[K, V]
		key     K
		value   V
		before  K
		wantErr bool
		expect  *OrderedMap[K, V]
	}
	tests := []testCase[int, string]{

		{
			name:    "inserts into correct location of multiple element map",
			o:       newFromPairs(kvp(987, "Employee 1"), kvp(443, "Employee 2"), kvp(101, "Employee 4"), kvp(814, "Employee 5")),
			key:     230,
			value:   "Employee 3",
			before:  101,
			wantErr: false,
			expect:  newFromPairs(kvp(987, "Employee 1"), kvp(443, "Employee 2"), kvp(230, "Employee 3"), kvp(101, "Employee 4"), kvp(814, "Employee 5")),
		},
		{
			name:    "errors if before key is not found in empty map",
			o:       New[int, string](),
			key:     1,
			value:   "3rd",
			before:  2,
			wantErr: true,
		},
		{
			name:    "errors if before key is not found in empty map (default key value)",
			o:       New[int, string](),
			key:     0,
			value:   "3rd",
			before:  2,
			wantErr: true,
		},
		{
			name:    "errors if before key is not found in populated map",
			o:       newFromPairs(kvp(987, "Employee 1"), kvp(443, "Employee 2"), kvp(101, "Employee 4"), kvp(814, "Employee 5")),
			key:     230,
			value:   "Employee 3",
			before:  11111,
			wantErr: true,
		},
		{
			name:    "errors if key and before are the same",
			o:       newFromPairs(kvp(987, "Employee 1"), kvp(443, "Employee 2"), kvp(230, "Employee 3"), kvp(101, "Employee 4"), kvp(814, "Employee 5")),
			key:     101,
			value:   "Employee 4",
			before:  101,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.o.InsertBefore(tt.key, tt.value, tt.before)
			if (err != nil) != tt.wantErr {
				t.Errorf("InsertBefore() error = %s, wantErr %v", err.Error(), tt.wantErr)
			}

			if err != nil {
				t.Logf("InsertBefore() error was: %s", err.Error())
			}

			if !tt.wantErr {
				if !Equal(tt.expect, tt.o) {
					diff := cmp.Diff(tt.expect.GoString(), tt.o.GoString())
					t.Errorf("InsertBefore() expected state mismatch:\n%s", diff)
				}
			}
		})
	}
}

func TestOrderedMap_Last(t *testing.T) {
	type testCase[K comparable, V any] struct {
		name string
		o    OrderedMap[K, V]
		want *pair[K, V]
	}
	tests := []testCase[string, int]{
		{
			name: "List is nil for empty map",
			o:    *New[string, int](),
			want: nil,
		},
		{
			name: "List is the element in single element map",
			o:    *newFromPairs[string, int](kvp("First", 1)),
			want: kvp("First", 1),
		},
		{
			name: "List is last element in multiple element map",
			o:    *newFromPairs[string, int](kvp("Z", 1), kvp("A", 2)),
			want: kvp("A", 2),
		},
		{
			name: "Last is last element in a manipulated map",
			o: func() OrderedMap[string, int] {
				original := newFromPairs[string, int](kvp("Z", 1), kvp("A", 2))
				_ = original.MoveBefore("A", "Z")
				return *original
			}(),
			want: kvp("Z", 1),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.o.Last()
			if !tt.want.Equals(got) {
				t.Errorf("Last() = %#v, want %#v", got, tt.want)
			}
		})
	}
}

func TestOrderedMap_MoveAfter(t *testing.T) {
	type testCase[K comparable, V any] struct {
		name    string
		o       *OrderedMap[K, V]
		key     K
		after   K
		wantErr bool
		expect  *OrderedMap[K, V]
	}
	tests := []testCase[string, string]{
		{
			name:    "MoveAfter moves the desired value at 'key' after the pair defined at 'after'",
			o:       newFromPairs(kvp("first", "1st"), kvp("second", "2nd"), kvp("fourth", "4th"), kvp("fifth", "5th")),
			key:     "first",
			after:   "second",
			wantErr: false,
			expect:  newFromPairs(kvp("second", "2nd"), kvp("first", "1st"), kvp("fourth", "4th"), kvp("fifth", "5th")),
		},
		{
			name:    "MoveAfter errors on empty map",
			o:       New[string, string](),
			key:     "first",
			after:   "second",
			wantErr: true,
		},
		{
			name:    "MoveAfter errors on populated map with missing key",
			o:       newFromPairs(kvp("second", "2nd"), kvp("fourth", "4th"), kvp("fifth", "5th")),
			key:     "first",
			after:   "second",
			wantErr: true,
		},
		{
			name:    "MoveAfter errors on populated map with missing 'after'",
			o:       newFromPairs(kvp("first", "1st"), kvp("fourth", "4th"), kvp("fifth", "5th")),
			key:     "first",
			after:   "second",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.o.MoveAfter(tt.key, tt.after)
			if (err != nil) != tt.wantErr {
				t.Errorf("MoveAfter() error = %s, wantErr %v", err.Error(), tt.wantErr)
			}

			if err != nil {
				t.Logf("MoveAfter() error was: %s", err.Error())
			}

			if !tt.wantErr {
				if !Equal(tt.expect, tt.o) {
					diff := cmp.Diff(tt.expect.GoString(), tt.o.GoString())
					t.Errorf("MoveAfter() expected state mismatch:\n%s", diff)
				}
			}
		})
	}
}

func TestOrderedMap_MoveBefore(t *testing.T) {
	type testCase[K comparable, V any] struct {
		name    string
		o       *OrderedMap[K, V]
		key     K
		before  K
		wantErr bool
		expect  *OrderedMap[K, V]
	}
	tests := []testCase[string, string]{
		{
			name:    "MoveBefore moves the desired value at 'key' after the pair defined at 'after'",
			o:       newFromPairs(kvp("first", "1st"), kvp("second", "2nd"), kvp("third", "3rd"), kvp("fourth", "4th"), kvp("fifth", "5th")),
			key:     "third",
			before:  "second",
			wantErr: false,
			expect:  newFromPairs(kvp("first", "1st"), kvp("third", "3rd"), kvp("second", "2nd"), kvp("fourth", "4th"), kvp("fifth", "5th")),
		},
		{
			name:    "MoveBefore errors on empty map",
			o:       New[string, string](),
			key:     "first",
			before:  "second",
			wantErr: true,
		},
		{
			name:    "MoveBefore errors on populated map with missing key",
			o:       newFromPairs(kvp("second", "2nd"), kvp("fourth", "4th"), kvp("fifth", "5th")),
			key:     "first",
			before:  "second",
			wantErr: true,
		},
		{
			name:    "MoveBefore errors on populated map with missing 'after'",
			o:       newFromPairs(kvp("first", "1st"), kvp("fourth", "4th"), kvp("fifth", "5th")),
			key:     "first",
			before:  "second",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.o.MoveBefore(tt.key, tt.before)
			if (err != nil) != tt.wantErr {
				t.Errorf("MoveBefore() error = %s, wantErr %v", err.Error(), tt.wantErr)
			}

			if err != nil {
				t.Logf("MoveBefore() error was: %s", err.Error())
			}

			if !tt.wantErr {
				if !Equal(tt.expect, tt.o) {
					diff := cmp.Diff(tt.expect.GoString(), tt.o.GoString())
					t.Errorf("MoveBefore() expected state mismatch:\n%s", diff)
				}
			}
		})
	}
}

func TestOrderedMap_MoveToBack(t *testing.T) {
	type testCase[K comparable, V any] struct {
		name    string
		key     K
		o       *OrderedMap[K, V]
		expect  *OrderedMap[K, V]
		wantErr bool
	}
	tests := []testCase[string, string]{
		{
			name:   "MoveToBack moves first element to back",
			o:      newFromPairs(kvp("first", "1st"), kvp("second", "2nd"), kvp("third", "3rd"), kvp("fourth", "4th"), kvp("fifth", "5th")),
			key:    "first",
			expect: newFromPairs(kvp("second", "2nd"), kvp("third", "3rd"), kvp("fourth", "4th"), kvp("fifth", "5th"), kvp("first", "1st")),
		},
		{
			name:   "MoveToBack moves Nth element to back",
			o:      newFromPairs(kvp("first", "1st"), kvp("second", "2nd"), kvp("third", "3rd"), kvp("fourth", "4th"), kvp("fifth", "5th")),
			key:    "third",
			expect: newFromPairs(kvp("first", "1st"), kvp("second", "2nd"), kvp("fourth", "4th"), kvp("fifth", "5th"), kvp("third", "3rd")),
		},
		{
			name:   "MoveToBack move last element is no-op",
			o:      newFromPairs(kvp("first", "1st"), kvp("second", "2nd"), kvp("third", "3rd"), kvp("fourth", "4th"), kvp("fifth", "5th")),
			key:    "fifth",
			expect: newFromPairs(kvp("first", "1st"), kvp("second", "2nd"), kvp("third", "3rd"), kvp("fourth", "4th"), kvp("fifth", "5th")),
		},
		{
			name:    "MoveToBack errors if the element is not found",
			o:       newFromPairs(kvp("first", "1st"), kvp("second", "2nd"), kvp("third", "3rd"), kvp("fourth", "4th"), kvp("fifth", "5th")),
			key:     "asdf",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.o.MoveToBack(tt.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("MoveToBack() error = %s, wantErr %v", err.Error(), tt.wantErr)
			}

			if err != nil {
				t.Logf("MoveToBack() error was: %s", err.Error())
			}

			if !tt.wantErr {
				if !Equal(tt.expect, tt.o) {
					diff := cmp.Diff(tt.expect.GoString(), tt.o.GoString())
					t.Errorf("MoveToBack() expected state mismatch:\n%s", diff)
				}
			}
		})
	}
}

func TestOrderedMap_MoveToFront(t *testing.T) {
	type testCase[K comparable, V any] struct {
		name    string
		key     K
		o       *OrderedMap[K, V]
		expect  *OrderedMap[K, V]
		wantErr bool
	}
	tests := []testCase[string, string]{
		{
			name:   "MoveToFront moves last element to front",
			o:      newFromPairs(kvp("first", "1st"), kvp("second", "2nd"), kvp("third", "3rd"), kvp("fourth", "4th"), kvp("fifth", "5th")),
			key:    "fifth",
			expect: newFromPairs(kvp("fifth", "5th"), kvp("first", "1st"), kvp("second", "2nd"), kvp("third", "3rd"), kvp("fourth", "4th")),
		},
		{
			name:   "MoveToFront moves Nth element to back",
			o:      newFromPairs(kvp("first", "1st"), kvp("second", "2nd"), kvp("third", "3rd"), kvp("fourth", "4th"), kvp("fifth", "5th")),
			key:    "third",
			expect: newFromPairs(kvp("third", "3rd"), kvp("first", "1st"), kvp("second", "2nd"), kvp("fourth", "4th"), kvp("fifth", "5th")),
		},
		{
			name:   "MoveToFront move first element is no-op",
			o:      newFromPairs(kvp("first", "1st"), kvp("second", "2nd"), kvp("third", "3rd"), kvp("fourth", "4th"), kvp("fifth", "5th")),
			key:    "first",
			expect: newFromPairs(kvp("first", "1st"), kvp("second", "2nd"), kvp("third", "3rd"), kvp("fourth", "4th"), kvp("fifth", "5th")),
		},
		{
			name:    "MoveToFront errors if the element is not found",
			o:       newFromPairs(kvp("first", "1st"), kvp("second", "2nd"), kvp("third", "3rd"), kvp("fourth", "4th"), kvp("fifth", "5th")),
			key:     "asdf",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.o.MoveToFront(tt.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("MoveToFront() error = %s, wantErr %v", err.Error(), tt.wantErr)
			}

			if err != nil {
				t.Logf("MoveToFront() error was: %s", err.Error())
			}

			if !tt.wantErr {
				if !Equal(tt.expect, tt.o) {
					diff := cmp.Diff(tt.expect.GoString(), tt.o.GoString())
					t.Errorf("MoveToFront() expected state mismatch:\n%s", diff)
				}
			}
		})
	}
}

func TestOrderedMap_Remove(t *testing.T) {
	type testCase[K comparable, V any] struct {
		name   string
		o      *OrderedMap[K, V]
		key    K
		want   *pair[K, V]
		wantOk bool
		expect *OrderedMap[K, V]
	}
	tests := []testCase[string, string]{
		{
			name:   "Removes an existing element",
			o:      newFromPairs(kvp("first", "1st"), kvp("second", "2nd"), kvp("fourth", "4th"), kvp("fifth", "5th")),
			key:    "second",
			want:   kvp("second", "2nd"),
			wantOk: true,
			expect: newFromPairs(kvp("first", "1st"), kvp("fourth", "4th"), kvp("fifth", "5th")),
		},
		{
			name:   "Removes does not modify map when element does not exist",
			o:      newFromPairs(kvp("first", "1st"), kvp("fourth", "4th"), kvp("fifth", "5th")),
			key:    "second",
			want:   kvp("second", "2nd"),
			wantOk: false,
			expect: newFromPairs(kvp("first", "1st"), kvp("fourth", "4th"), kvp("fifth", "5th")),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, ok := tt.o.Remove(tt.key)
			if ok {
				actual := new(pair[string, string])
				if got != nil {
					actual = kvp(got.Key, got.Value)
				}
				if !reflect.DeepEqual(actual, tt.want) {
					t.Errorf("Remove() got: %v, want %v", actual, tt.want)
				}
			}

			if ok != tt.wantOk {
				t.Errorf("Remove() ok %v, wantOk %v", ok, tt.wantOk)
			}

			if !Equal(tt.expect, tt.o) {
				diff := cmp.Diff(tt.expect.GoString(), tt.o.GoString())
				t.Errorf("Remove() expected state mismatch:\n%s", diff)
			}
		})
	}
}

func TestOrderedMap_Set(t *testing.T) {
	type testCase[K comparable, V any] struct {
		name   string
		o      *OrderedMap[K, V]
		key    K
		value  V
		expect *OrderedMap[K, V]
	}
	tests := []testCase[string, string]{
		{
			name:   "Set on new map",
			o:      New[string, string](),
			key:    "first",
			value:  "1st",
			expect: newFromPairs(kvp("first", "1st")),
		},
		{
			name:   "Set on single value map",
			o:      newFromPairs(kvp("first", "1st")),
			key:    "second",
			value:  "2nd",
			expect: newFromPairs(kvp("first", "1st"), kvp("second", "2nd")),
		},
		{
			name:   "Set on multiple value map",
			o:      newFromPairs(kvp("first", "1st"), kvp("second", "2nd")),
			key:    "third",
			value:  "3rd",
			expect: newFromPairs(kvp("first", "1st"), kvp("second", "2nd"), kvp("third", "3rd")),
		},
		{
			name:   "Updates existing value in map without changing order",
			o:      newFromPairs(kvp("first", "1st"), kvp("second", "2nd"), kvp("third", "3rd")),
			key:    "second",
			value:  ":(",
			expect: newFromPairs(kvp("first", "1st"), kvp("second", ":("), kvp("third", "3rd")),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.o.Set(tt.key, tt.value)
			if !Equal(tt.expect, tt.o) {
				diff := cmp.Diff(tt.expect.GoString(), tt.o.GoString())
				t.Errorf("Set() expected state mismatch:\n%s", diff)
			}
		})
	}
}
