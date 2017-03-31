package main

import (
	"testing"

	"github.com/k0kubun/pp"
)

func TestNewMorphizer(t *testing.T) {
	tests := []struct {
		input []string
		want  []Morphs
	}{
		{
			input: []string{"a"},
			want: []Morphs{
				Morphs{
					&Morph{"a", "名詞", "固有名詞", "組織", "*", "*", "*", "*", "", ""},
				},
			},
		},
		{
			input: []string{"a", "b", "c"},
			want: []Morphs{
				Morphs{
					&Morph{"a", "名詞", "固有名詞", "組織", "*", "*", "*", "*", "", ""},
				},
				Morphs{
					&Morph{"b", "名詞", "固有名詞", "組織", "*", "*", "*", "*", "", ""},
				},
				Morphs{
					&Morph{"c", "名詞", "固有名詞", "組織", "*", "*", "*", "*", "", ""},
				},
			},
		},
	}

	for _, test := range tests {
		out, in, _ := NewMorphizer()

		// input
		go func() {
			for _, str := range test.input {
				in <- str
			}
			close(in)
		}()

		// output
		ans := make([]Morphs, 0, len(test.input))
		for mok := range out {
			ans = append(ans, mok)
		}
		got := ans[:]

		// verify
		for l := len(ans); l > 0 && len(ans) == l; l-- {
			for _, want := range test.want {
				for i, ms := range ans {
					if want.IsEqual(ms) {
						ans = append(ans[:i], ans[i+1:]...)
					}
				}
			}
		}
		if len(ans) > 0 {
			t.Errorf("in %v, out %v, want %v", test.input, pp.Sprint(got), pp.Sprint(test.want))
		}
	}
}
