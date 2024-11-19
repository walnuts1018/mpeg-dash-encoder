package sliceequal

import (
	"strings"
	"testing"
)

func TestEqual(t *testing.T) {
	type args struct {
		s1  []string
		s2  []string
		cmp func(a, b string) int
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "deep equal",
			args: args{
				s1: []string{"a", "b", "c"},
				s2: []string{"a", "b", "c"},
				cmp: func(a, b string) int {
					return strings.Compare(a, b)
				},
			},
			want: true,
		},
		{
			name: "order changed",
			args: args{
				s1: []string{"a", "b", "c"},
				s2: []string{"c", "b", "a"},
				cmp: func(a, b string) int {
					return strings.Compare(a, b)
				},
			},
			want: true,
		},
		{
			name: "not equal",
			args: args{
				s1: []string{"a", "b", "c"},
				s2: []string{"a", "b", "d"},
				cmp: func(a, b string) int {
					return strings.Compare(a, b)
				},
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Equal(tt.args.s1, tt.args.s2, tt.args.cmp); got != tt.want {
				t.Errorf("Equal() = %v, want %v", got, tt.want)
			}
		})
	}
}
