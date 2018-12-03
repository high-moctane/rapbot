package main

import (
	"log"
	"math"
	"reflect"
	"testing"

	"github.com/ikawaha/kagome/tokenizer"
)

func TestNewMorae(t *testing.T) {
	type answer struct {
		m  []Mora
		ok bool
	}
	tests := []struct {
		input string
		want  answer
	}{
		{
			input: "ム",
			want: answer{
				m:  []Mora{{"m", "u"}},
				ok: true,
			},
		},
		{
			input: "タピオカ",
			want: answer{
				m:  []Mora{{"t", "a"}, {"p", "i"}, {"", "o"}, {"k", "a"}},
				ok: true,
			},
		},
		{
			input: "ケーキ",
			want: answer{
				m:  []Mora{{"k", "e"}, {"", "e"}, {"k", "i"}},
				ok: true,
			},
		},
		{
			input: "クッキー",
			want: answer{
				m:  []Mora{{"k", "u"}, {"*xtu", "*xtu"}, {"k", "i"}, {"", "i"}},
				ok: true,
			},
		},
		{
			input: "",
			want:  answer{nil, false},
		},
		{
			input: "くれーぷ",
			want:  answer{nil, false},
		},
	}

	for _, test := range tests {
		var ans answer
		ans.m, ans.ok = NewMorae(test.input)
		if !reflect.DeepEqual(ans, test.want) {
			t.Errorf("NewMorae(%q) = %v, %v, want %v, %v",
				test.input, ans.m, ans.ok, test.want.m, test.want.ok)
		}
	}
}

func TestNewMoraeWeight(t *testing.T) {
	tests := []struct {
		input []MoraWeight
		want  *MoraeWeight
	}{
		{
			input: []MoraWeight{},
			want:  &MoraeWeight{MW: []MoraWeight{}, Sum: 0.0},
		},
		{
			input: []MoraWeight{{1.0, 2.0}},
			want:  &MoraeWeight{MW: []MoraWeight{{1.0, 2.0}}, Sum: 3.0},
		},
		{
			input: []MoraWeight{{1.0, 2.0}, {4.0, 8.0}},
			want: &MoraeWeight{MW: []MoraWeight{
				{1.0, 2.0}, {4.0, 8.0}}, Sum: 15.0,
			},
		},
	}

	for _, test := range tests {
		ans := NewMoraeWeight(test.input)
		if !reflect.DeepEqual(ans, test.want) {
			t.Errorf("NewMoraeWeight(%v) = %v, want %v",
				test.input, ans, test.want)
		}
	}
}

func TestMoraeWeightSimilarity(t *testing.T) {
	mw := NewMoraeWeight([]MoraWeight{{1.0, 2.0}, {4.0, 8.0}})

	tests := []struct {
		input [][]Mora
		want  float64
	}{
		{
			input: [][]Mora{{}, {{"", "a"}}},
			want:  0.0,
		},
		{
			input: [][]Mora{{{"", "a"}}, {{"", "a"}}},
			want:  0.8,
		},
		{
			input: [][]Mora{
				{{"", "a"}, {"", "i"}},
				{{"", "i"}, {"k", "a"}},
			},
			want: 0.0667,
		},
		{
			input: [][]Mora{
				{{"k", "e"}, {"", "e"}, {"k", "i"}},
				{{"k", "u"}, {"*xtu", "*xtu"}, {"k", "i"}, {"", "i"}},
			},
			want: 0.533,
		},
		{
			input: [][]Mora{
				{{"p", "a"}, {"s", "e"}, {"r", "i"}},
				{{"", "i"}},
			},
			want: 0.533,
		},
	}

	for _, test := range tests {
		ans := mw.Similarity(test.input[0], test.input[1])
		if math.Abs(ans-test.want) > 0.01 {
			t.Errorf("(*MoraeWeight).Similarity(%v, %v) ~= %.2f, want %.2f",
				test.input[0], test.input[1], ans, test.want)
		}
	}
}

func TestMoraeWeightSimMorphs(t *testing.T) {
	mw := NewMoraeWeight([]MoraWeight{{1.0, 2.0}, {4.0, 8.0}})
	kagome := tokenizer.New()

	tests := []struct {
		input [2]string
		want  float64
	}{
		{
			input: [2]string{"美味しいケーキ", "まずいケーキ"},
			want:  1.0,
		},
	}

	for _, test := range tests {
		ms0 := NewMorphs(kagome.Tokenize(test.input[0]))
		ms1 := NewMorphs(kagome.Tokenize(test.input[1]))
		log.Print(mw.SimMorphs(ms0, ms1))
		if math.Abs(mw.SimMorphs(ms0, ms1)-test.want) > 0.001 {
			t.Errorf("(*MoraeWeight).SimMorphs(%v, %v) = %.3f, want %.3f",
				ms0, ms1, mw.SimMorphs(ms0, ms1), test.want)
		}
	}
}
