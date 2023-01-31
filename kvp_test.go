package orderedmap

import (
	"testing"
)

type kvpValue struct {
	Name string
	X    int
	Y    int
}

func TestKeyValuePair_String(t *testing.T) {
	type testCase[K comparable, V any] struct {
		name string
		k    KeyValuePair[K, V]
		want string
	}
	tests := []testCase[string, any]{
		{
			name: "String prints key=value format without unexported fields (struct)",
			k:    KeyValuePair[string, any]{Key: "MyValue", Value: kvpValue{Name: "Value", X: 2, Y: 6}},
			want: "MyValue={Name:Value X:2 Y:6}",
		},
		{
			name: "String prints key=value format (string)",
			k:    KeyValuePair[string, any]{Key: "MyValue", Value: "asdfjkl;qwertyuiop"},
			want: "MyValue=asdfjkl;qwertyuiop",
		},
		{
			name: "String prints key=value format (slice)",
			k:    KeyValuePair[string, any]{Key: "MyValue", Value: []int{1, 2, 3}},
			want: "MyValue=[1 2 3]",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.k.String(); got != tt.want {
				t.Errorf("String() = %v, want %v", got, tt.want)
			}
		})
	}
}
