package main

import (
	"reflect"
	"testing"

	"github.com/ikawaha/kagome/tokenizer"
)

func TestNewMorae(t *testing.T) {
	tests := []struct {
		pronunciation string
		morae         Morae
		ok            bool
	}{
		{
			"キ",
			Morae{&Mora{"k", "i"}},
			true,
		},
		{
			"ビャイ",
			Morae{&Mora{"by", "a"}, &Mora{"", "i"}},
			true,
		},
		{
			"ゴファー",
			Morae{&Mora{"g", "o"}, &Mora{"f", "a"}, &Mora{"", "a"}},
			true,
		},
		{
			"　",
			nil,
			false,
		},
		{
			"",
			nil,
			false,
		},
	}

	for idx, test := range tests {
		morae, ok := NewMorae(test.pronunciation)
		if test.ok != ok {
			t.Errorf("[%d] expected %v, but got %v", idx, test.ok, ok)
		}
		if !reflect.DeepEqual(test.morae, morae) {
			t.Errorf("[%d] expected %v, but got %v", idx, test.morae, morae)
		}
	}
}

func TestAnalyzeText(t *testing.T) {
	tests := []struct {
		text     string
		sentence Sentence
	}{
		{
			"まじか",
			Sentence{
				&Morph{"BOS", "", "", "", "", "", "", "", "", ""},
				&Morph{"まじ", "名詞", "形容動詞語幹", "*", "*", "*", "*", "まじ", "マジ", "マジ"},
				&Morph{"か", "助詞", "副助詞／並立助詞／終助詞", "*", "*", "*", "*", "か", "カ", "カ"},
				&Morph{"EOS", "", "", "", "", "", "", "", "", ""},
			},
		},
	}

	tk := tokenizer.New()

	for idx, test := range tests {
		sentence := analyzeText(&tk, test.text)
		if !reflect.DeepEqual(test.sentence, sentence) {
			t.Errorf("[%d] expected %v, but got %v", idx, test.sentence, sentence)
		}
	}
}
