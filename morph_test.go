package main

import (
	"testing"

	"github.com/ikawaha/kagome/tokenizer"
)

func TestMorphString(t *testing.T) {
	tests := []struct {
		input Morph
		want  string
	}{
		{
			input: MorphBOS,
			want:  "\tBOS,,,,,,,,",
		},
		{
			input: Morph{"0", "1", "2", "3", "4", "5", "6", "7", "8", "9"},
			want:  "0\t1,2,3,4,5,6,7,8,9",
		},
	}

	for _, test := range tests {
		ans := test.input.String()
		if ans != test.want {
			t.Errorf("(*Morph).String(&%+v) = %q, want %q\n",
				test.input, ans, test.want)
		}
	}
}

func TestMorphsIsEqual(t *testing.T) {
	tests := []struct {
		input [2]Morphs
		want  bool
	}{
		{
			input: [2]Morphs{{}, {}},
			want:  true,
		},
		{
			input: [2]Morphs{
				{{Surface: "a"}, {Surface: "b"}},
				{{Surface: "a"}, {Surface: "b"}},
			},
			want: true,
		},
		{
			input: [2]Morphs{
				{{Surface: "a"}},
				{{Surface: "a"}, {Surface: "b"}},
			},
			want: false,
		},
		{
			input: [2]Morphs{
				{{Surface: "a"}, {Surface: "b"}},
				{{Surface: "a"}, {Surface: "c"}},
			},
			want: false,
		},
	}

	for _, test := range tests {
		if test.input[0].IsEqual(test.input[1]) != test.want {
			t.Errorf("(%v).IsEqual(%v) != %v\n",
				test.input[0], test.input[1], test.want)
		}
	}
}

func TestNewMorphs(t *testing.T) {
	type answer struct {
		Morphs
		ok bool
	}
	tests := []struct {
		input string
		want  answer
	}{
		{
			input: "",
			want:  answer{Morphs: Morphs{&MorphBOS, &MorphEOS}, ok: true},
		},
		{
			input: "愛",
			want: answer{
				Morphs: Morphs{
					&MorphBOS,
					&Morph{"愛", "名詞", "一般", "*", "*", "*", "*", "愛", "アイ", "アイ"},
					&MorphEOS,
				},
				ok: true,
			},
		},
	}

	kagome := tokenizer.New()

	for _, test := range tests {
		tokens := kagome.Tokenize(test.input)
		ms, ok := NewMorphs(tokens)
		if !ms.IsEqual(test.want.Morphs) || ok != test.want.ok {
			t.Errorf("NewMorphs(%q) = %v, %v, want %v, %v\n",
				test.input, ms, ok, test.want.Morphs, test.want.ok)
		}
	}
}
