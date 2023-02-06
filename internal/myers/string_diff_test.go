package myers

import (
	"testing"
)

func TestDiff(t *testing.T) {
	type args struct {
		first  string
		second string
	}
	tests := []struct {
		name   string
		args   args
		want   string
		wantOk bool
	}{
		{
			name:   "example: ABCABBA -> CBABAC",
			args:   args{first: "ABCABBA", second: "CBABAC"},
			want:   "\033[31mA\033[0m\033[31mB\033[0mC\033[32mB\033[0mAB\033[31mB\033[0mA\033[32mC\033[0m",
			wantOk: true,
		},
		{
			name:   "equal strings are equal",
			args:   args{first: "anteater", second: "anteater"},
			want:   "",
			wantOk: false,
		},
		{
			name:   "equal strings are equal (multi-line)",
			args:   args{first: "anteaters\nare\nawesome", second: "anteaters\nare\nawesome"},
			want:   "",
			wantOk: false,
		},
		{
			name:   "unequal strings are not equal (single line)",
			args:   args{first: "anteater", second: "anteaters"},
			want:   "anteater\033[32ms\033[0m",
			wantOk: true,
		},
		{
			name:   "unequal strings are not equal (single line, reversed)",
			args:   args{first: "anteaters", second: "anteater"},
			want:   "anteater\033[31ms\033[0m",
			wantOk: true,
		},
		{
			name:   "unequal strings are not equal (multi-line)",
			args:   args{first: "anteaters\nare\nawesome", second: "anteaters\nare\nlame"},
			want:   "anteaters\nare\n\033[32ml\033[0ma\033[31mw\033[0m\033[31me\033[0m\033[31ms\033[0m\033[31mo\033[0mme",
			wantOk: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, ok := Diff(tt.args.first, tt.args.second)
			if got != tt.want {
				t.Errorf("Diff() got = %v, want %v", got, tt.want)
			}
			if ok {
				t.Logf("Diff: %s", got)
			}
			if ok != tt.wantOk {
				t.Errorf("Diff() ok = %v, want %v", ok, tt.wantOk)
			}
		})
	}
}
