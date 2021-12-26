package strs

import (
	"reflect"
	"testing"
)

func TestAllBetweenPattern(t *testing.T) {
	type args struct {
		s       string
		pattern string
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			name: "no match",
			args: args{
				s:       "abcsome text to extractabc",
				pattern: "abcd",
			},
			want: []string{},
		},
		{
			name: "multi string pattern single match",
			args: args{
				s:       "abcsome text to extractabc",
				pattern: "abc",
			},
			want: []string{"some text to extract"},
		},
		{
			name: "multi string pattern multi match",
			args: args{
				s:       "abcsome text to extractabcabcsome other textabc",
				pattern: "abc",
			},
			want: []string{"some text to extract", "some other text"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := AllBetweenPattern(tt.args.s, tt.args.pattern)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("AllBetweenPattern() = %v, want %v", got, tt.want)
			}
		})
	}
}
