package main

import (
	"reflect"
	"sync"
	"testing"
)

func TestMarkov_Add(t *testing.T) {
	tests := []struct {
		morphss  [][]*Morph
		params   *MarkovParams
		learning chain
		chains   []chain
	}{
		{
			[][]*Morph{
				[]*Morph{
					&Morph{"BOS", "", "", "", "", "", "", "", "", ""},
					&Morph{"おはよう", "感動詞", "*", "*", "*", "*", "*", "おはよう", "オハヨウ", "オハヨー"},
					&Morph{"ござい", "助動詞", "*", "*", "*", "五段・ラ行特殊", "連用形", "ござる", "ゴザイ", "ゴザイ"},
					&Morph{"ます", "助動詞", "*", "*", "*", "特殊・マス", "基本形", "ます", "マス", "マス"},
					&Morph{"EOS", "", "", "", "", "", "", "", "", ""},
				},
			},
			&MarkovParams{
				Ngram:          2,
				ChainNum:       2,
				ChainMorphsNum: 5,
			},
			chain{
				Morph{"BOS", "", "", "", "", "", "", "", "", ""}: chain{
					Morph{"おはよう", "感動詞", "*", "*", "*", "*", "*", "おはよう", "オハヨウ", "オハヨー"}: chain{},
				},
				Morph{"おはよう", "感動詞", "*", "*", "*", "*", "*", "おはよう", "オハヨウ", "オハヨー"}: chain{
					Morph{"ござい", "助動詞", "*", "*", "*", "五段・ラ行特殊", "連用形", "ござる", "ゴザイ", "ゴザイ"}: chain{},
				},
				Morph{"ござい", "助動詞", "*", "*", "*", "五段・ラ行特殊", "連用形", "ござる", "ゴザイ", "ゴザイ"}: chain{
					Morph{"ます", "助動詞", "*", "*", "*", "特殊・マス", "基本形", "ます", "マス", "マス"}: chain{},
				},
				Morph{"ます", "助動詞", "*", "*", "*", "特殊・マス", "基本形", "ます", "マス", "マス"}: chain{
					Morph{"EOS", "", "", "", "", "", "", "", "", ""}: chain{},
				},
			},
			nil,
		},
		{
			[][]*Morph{
				[]*Morph{
					&Morph{"BOS", "", "", "", "", "", "", "", "", ""},
					&Morph{"おはよう", "感動詞", "*", "*", "*", "*", "*", "おはよう", "オハヨウ", "オハヨー"},
					&Morph{"ござい", "助動詞", "*", "*", "*", "五段・ラ行特殊", "連用形", "ござる", "ゴザイ", "ゴザイ"},
					&Morph{"ます", "助動詞", "*", "*", "*", "特殊・マス", "基本形", "ます", "マス", "マス"},
					&Morph{"EOS", "", "", "", "", "", "", "", "", ""},
				},
			},
			&MarkovParams{
				Ngram:          3,
				ChainNum:       2,
				ChainMorphsNum: 5,
			},
			chain{
				Morph{"BOS", "", "", "", "", "", "", "", "", ""}: chain{
					Morph{"おはよう", "感動詞", "*", "*", "*", "*", "*", "おはよう", "オハヨウ", "オハヨー"}: chain{
						Morph{"ござい", "助動詞", "*", "*", "*", "五段・ラ行特殊", "連用形", "ござる", "ゴザイ", "ゴザイ"}: chain{},
					},
				},
				Morph{"おはよう", "感動詞", "*", "*", "*", "*", "*", "おはよう", "オハヨウ", "オハヨー"}: chain{
					Morph{"ござい", "助動詞", "*", "*", "*", "五段・ラ行特殊", "連用形", "ござる", "ゴザイ", "ゴザイ"}: chain{
						Morph{"ます", "助動詞", "*", "*", "*", "特殊・マス", "基本形", "ます", "マス", "マス"}: chain{},
					},
				},
				Morph{"ござい", "助動詞", "*", "*", "*", "五段・ラ行特殊", "連用形", "ござる", "ゴザイ", "ゴザイ"}: chain{
					Morph{"ます", "助動詞", "*", "*", "*", "特殊・マス", "基本形", "ます", "マス", "マス"}: chain{
						Morph{"EOS", "", "", "", "", "", "", "", "", ""}: chain{},
					},
				},
			},
			nil,
		},
		{
			[][]*Morph{
				[]*Morph{
					&Morph{"BOS", "", "", "", "", "", "", "", "", ""},
					&Morph{"おはよう", "感動詞", "*", "*", "*", "*", "*", "おはよう", "オハヨウ", "オハヨー"},
					&Morph{"ござい", "助動詞", "*", "*", "*", "五段・ラ行特殊", "連用形", "ござる", "ゴザイ", "ゴザイ"},
					&Morph{"ます", "助動詞", "*", "*", "*", "特殊・マス", "基本形", "ます", "マス", "マス"},
					&Morph{"EOS", "", "", "", "", "", "", "", "", ""},
				},
			},
			&MarkovParams{
				Ngram:          3,
				ChainNum:       2,
				ChainMorphsNum: 2,
			},
			chain{
				Morph{"ござい", "助動詞", "*", "*", "*", "五段・ラ行特殊", "連用形", "ござる", "ゴザイ", "ゴザイ"}: chain{
					Morph{"ます", "助動詞", "*", "*", "*", "特殊・マス", "基本形", "ます", "マス", "マス"}: chain{
						Morph{"EOS", "", "", "", "", "", "", "", "", ""}: chain{},
					},
				},
			},
			[]chain{
				chain{
					Morph{"BOS", "", "", "", "", "", "", "", "", ""}: chain{
						Morph{"おはよう", "感動詞", "*", "*", "*", "*", "*", "おはよう", "オハヨウ", "オハヨー"}: chain{
							Morph{"ござい", "助動詞", "*", "*", "*", "五段・ラ行特殊", "連用形", "ござる", "ゴザイ", "ゴザイ"}: chain{},
						},
					},
					Morph{"おはよう", "感動詞", "*", "*", "*", "*", "*", "おはよう", "オハヨウ", "オハヨー"}: chain{
						Morph{"ござい", "助動詞", "*", "*", "*", "五段・ラ行特殊", "連用形", "ござる", "ゴザイ", "ゴザイ"}: chain{
							Morph{"ます", "助動詞", "*", "*", "*", "特殊・マス", "基本形", "ます", "マス", "マス"}: chain{},
						},
					},
				},
			},
		},
		{
			[][]*Morph{
				[]*Morph{
					&Morph{"BOS", "", "", "", "", "", "", "", "", ""},
					&Morph{"おはよう", "感動詞", "*", "*", "*", "*", "*", "おはよう", "オハヨウ", "オハヨー"},
					&Morph{"ござい", "助動詞", "*", "*", "*", "五段・ラ行特殊", "連用形", "ござる", "ゴザイ", "ゴザイ"},
					&Morph{"ます", "助動詞", "*", "*", "*", "特殊・マス", "基本形", "ます", "マス", "マス"},
					&Morph{"EOS", "", "", "", "", "", "", "", "", ""},
				},
				[]*Morph{
					&Morph{"BOS", "", "", "", "", "", "", "", "", ""},
					&Morph{"おはよう", "感動詞", "*", "*", "*", "*", "*", "おはよう", "オハヨウ", "オハヨー"},
					&Morph{"さん", "名詞", "接尾", "人名", "*", "*", "*", "さん", "サン", "サン"},
					&Morph{"EOS", "", "", "", "", "", "", "", "", ""},
				},
			},
			&MarkovParams{
				Ngram:          2,
				ChainNum:       2,
				ChainMorphsNum: 2,
			},
			chain{
				Morph{"さん", "名詞", "接尾", "人名", "*", "*", "*", "さん", "サン", "サン"}: chain{
					Morph{"EOS", "", "", "", "", "", "", "", "", ""}: chain{},
				},
			},
			[]chain{
				chain{
					Morph{"ござい", "助動詞", "*", "*", "*", "五段・ラ行特殊", "連用形", "ござる", "ゴザイ", "ゴザイ"}: chain{
						Morph{"ます", "助動詞", "*", "*", "*", "特殊・マス", "基本形", "ます", "マス", "マス"}: chain{},
					},
					Morph{"ます", "助動詞", "*", "*", "*", "特殊・マス", "基本形", "ます", "マス", "マス"}: chain{
						Morph{"EOS", "", "", "", "", "", "", "", "", ""}: chain{},
					},
				},
				chain{
					Morph{"BOS", "", "", "", "", "", "", "", "", ""}: chain{
						Morph{"おはよう", "感動詞", "*", "*", "*", "*", "*", "おはよう", "オハヨウ", "オハヨー"}: chain{},
					},
					Morph{"おはよう", "感動詞", "*", "*", "*", "*", "*", "おはよう", "オハヨウ", "オハヨー"}: chain{
						Morph{"さん", "名詞", "接尾", "人名", "*", "*", "*", "さん", "サン", "サン"}: chain{},
					},
				},
			},
		},
	}

	for idx, test := range tests {
		m := NewMarkov(test.params)
		for _, morphs := range test.morphss {
			m.Add(morphs)
		}
		if !reflect.DeepEqual(test.learning, m.learning) {
			t.Errorf("[%d] learning: expected\n%v, but got\n%v", idx, test.learning, m.learning)
		}
		if !reflect.DeepEqual(test.chains, m.chains) {
			t.Errorf("[%d] chains: expected\n%v, but got\n%v", idx, test.chains, m.chains)
		}
	}
}

func TestMarkov_RandomSentence(t *testing.T) {
	tests := []struct {
		morphLen int
		markov   Markov
		sentence Sentence
	}{
		{
			3,
			Markov{
				params: &MarkovParams{
					Ngram: 2,
				},
				mu: new(sync.RWMutex),
				chains: []chain{
					chain{
						Morph{"BOS", "", "", "", "", "", "", "", "", ""}: chain{
							Morph{"あ", "", "", "", "", "", "", "", "", ""}: chain{},
						},
						Morph{"あ", "", "", "", "", "", "", "", "", ""}: chain{
							Morph{"い", "", "", "", "", "", "", "", "", ""}: chain{},
						},
						Morph{"い", "", "", "", "", "", "", "", "", ""}: chain{
							Morph{"う", "", "", "", "", "", "", "", "", ""}: chain{},
						},
						Morph{"う", "", "", "", "", "", "", "", "", ""}: chain{
							Morph{"え", "", "", "", "", "", "", "", "", ""}: chain{},
						},
					},
				},
			},
			Sentence{
				&Morph{"あ", "", "", "", "", "", "", "", "", ""},
				&Morph{"い", "", "", "", "", "", "", "", "", ""},
				&Morph{"う", "", "", "", "", "", "", "", "", ""},
			},
		},
		{
			3,
			Markov{
				params: &MarkovParams{
					Ngram: 2,
				},
				mu: new(sync.RWMutex),
				chains: []chain{
					chain{
						Morph{"BOS", "", "", "", "", "", "", "", "", ""}: chain{
							Morph{"あ", "", "", "", "", "", "", "", "", ""}: chain{},
						},
						Morph{"あ", "", "", "", "", "", "", "", "", ""}: chain{
							Morph{"い", "", "", "", "", "", "", "", "", ""}: chain{},
						},
						Morph{"い", "", "", "", "", "", "", "", "", ""}: chain{
							EOS: chain{},
						},
						Morph{"う", "", "", "", "", "", "", "", "", ""}: chain{
							Morph{"え", "", "", "", "", "", "", "", "", ""}: chain{},
						},
					},
				},
			},
			Sentence{
				&Morph{"あ", "", "", "", "", "", "", "", "", ""},
				&Morph{"い", "", "", "", "", "", "", "", "", ""},
			},
		},
	}

	for idx, test := range tests {
		sentence, _ := test.markov.RandomSentence(test.morphLen)
		if !reflect.DeepEqual(test.sentence, sentence) {
			t.Errorf("[%d] expected %v, but got %v", idx, test.sentence, sentence)
		}
	}
}

func TestMarkov_ChoiceMorph(t *testing.T) {
	tests := []struct {
		morphs []*Morph
		markov Markov
		morph  *Morph
	}{
		{
			[]*Morph{
				&Morph{"あ", "", "", "", "", "", "", "", "", ""},
			},
			Markov{
				chains: []chain{
					chain{
						Morph{"あ", "", "", "", "", "", "", "", "", ""}: chain{
							Morph{"い", "", "", "", "", "", "", "", "", ""}: chain{},
						},
					},
					chain{
						Morph{"う", "", "", "", "", "", "", "", "", ""}: chain{
							Morph{"え", "", "", "", "", "", "", "", "", ""}: chain{},
						},
					},
				},
			},
			&Morph{"い", "", "", "", "", "", "", "", "", ""},
		},
		{
			[]*Morph{
				&Morph{"あ", "", "", "", "", "", "", "", "", ""},
				&Morph{"い", "", "", "", "", "", "", "", "", ""},
			},
			Markov{
				chains: []chain{
					chain{
						Morph{"あ", "", "", "", "", "", "", "", "", ""}: chain{
							Morph{"い", "", "", "", "", "", "", "", "", ""}: chain{
								Morph{"う", "", "", "", "", "", "", "", "", ""}: chain{},
							},
							Morph{"え", "", "", "", "", "", "", "", "", ""}: chain{
								Morph{"お", "", "", "", "", "", "", "", "", ""}: chain{},
								Morph{"か", "", "", "", "", "", "", "", "", ""}: chain{},
							},
						},
					},
					chain{
						Morph{"う", "", "", "", "", "", "", "", "", ""}: chain{
							Morph{"え", "", "", "", "", "", "", "", "", ""}: chain{
								Morph{"お", "", "", "", "", "", "", "", "", ""}: chain{},
								Morph{"か", "", "", "", "", "", "", "", "", ""}: chain{},
							},
						},
					},
				},
			},
			&Morph{"う", "", "", "", "", "", "", "", "", ""},
		},
	}

	for idx, test := range tests {
		morph, _ := test.markov.RandomMorph(test.morphs)
		if *test.morph != *morph {
			t.Errorf("[%d] expected %v, but got %v", idx, *test.morph, *morph)
		}
	}
}

func TestRandomIndice(t *testing.T) {
	// t.Fail()
	t.SkipNow()

	t.Log(randomIndice(10))
	t.Log(randomIndice(10))
	t.Log(randomIndice(10))
	t.Log(randomIndice(10))
}

func TestChain_Add(t *testing.T) {
	tests := []struct {
		morphss [][]*Morph
		c       chain
	}{
		{
			[][]*Morph{
				[]*Morph{
					&Morph{"ぽ", "", "", "", "", "", "", "", "", ""},
					&Morph{"わ", "", "", "", "", "", "", "", "", ""},
				},
			},
			chain{
				Morph{"ぽ", "", "", "", "", "", "", "", "", ""}: chain{
					Morph{"わ", "", "", "", "", "", "", "", "", ""}: chain{},
				},
			},
		},
		{
			[][]*Morph{
				[]*Morph{
					&Morph{"ぽ", "", "", "", "", "", "", "", "", ""},
					&Morph{"わ", "", "", "", "", "", "", "", "", ""},
				},
				[]*Morph{
					&Morph{"め", "", "", "", "", "", "", "", "", ""},
					&Morph{"う", "", "", "", "", "", "", "", "", ""},
				},
			},
			chain{
				Morph{"ぽ", "", "", "", "", "", "", "", "", ""}: chain{
					Morph{"わ", "", "", "", "", "", "", "", "", ""}: chain{},
				},
				Morph{"め", "", "", "", "", "", "", "", "", ""}: chain{
					Morph{"う", "", "", "", "", "", "", "", "", ""}: chain{},
				},
			},
		},
		{
			[][]*Morph{
				[]*Morph{
					&Morph{"ぽ", "", "", "", "", "", "", "", "", ""},
					&Morph{"わ", "", "", "", "", "", "", "", "", ""},
				},
				[]*Morph{
					&Morph{"ぽ", "", "", "", "", "", "", "", "", ""},
					&Morph{"い", "", "", "", "", "", "", "", "", ""},
				},
				[]*Morph{
					&Morph{"め", "", "", "", "", "", "", "", "", ""},
					&Morph{"う", "", "", "", "", "", "", "", "", ""},
				},
			},
			chain{
				Morph{"ぽ", "", "", "", "", "", "", "", "", ""}: chain{
					Morph{"わ", "", "", "", "", "", "", "", "", ""}: chain{},
					Morph{"い", "", "", "", "", "", "", "", "", ""}: chain{},
				},
				Morph{"め", "", "", "", "", "", "", "", "", ""}: chain{
					Morph{"う", "", "", "", "", "", "", "", "", ""}: chain{},
				},
			},
		},
	}

	for idx, test := range tests {
		c := make(chain)

		for _, morphs := range test.morphss {
			c.Add(morphs)
		}
		if !reflect.DeepEqual(test.c, c) {
			t.Errorf("[%d] expected\n%v, but got\n%v", idx, test.c, c)
		}
	}
}

func TestChain_Choice(t *testing.T) {
	tests := []struct {
		morphs []*Morph
		chain  chain
		morph  *Morph
		ok     bool
	}{
		{
			[]*Morph{&Morph{"あ", "", "", "", "", "", "", "", "", ""}},
			chain{
				Morph{"あ", "", "", "", "", "", "", "", "", ""}: chain{
					Morph{"い", "", "", "", "", "", "", "", "", ""}: chain{},
				},
			},
			&Morph{"い", "", "", "", "", "", "", "", "", ""},
			true,
		},
		{
			[]*Morph{
				&Morph{"あ", "", "", "", "", "", "", "", "", ""},
				&Morph{"い", "", "", "", "", "", "", "", "", ""},
			},
			chain{
				Morph{"あ", "", "", "", "", "", "", "", "", ""}: chain{
					Morph{"い", "", "", "", "", "", "", "", "", ""}: chain{
						Morph{"う", "", "", "", "", "", "", "", "", ""}: chain{},
					},
					Morph{"え", "", "", "", "", "", "", "", "", ""}: chain{
						Morph{"お", "", "", "", "", "", "", "", "", ""}: chain{},
						Morph{"か", "", "", "", "", "", "", "", "", ""}: chain{},
					},
				},
			},
			&Morph{"う", "", "", "", "", "", "", "", "", ""},
			true,
		},
	}

	for idx, test := range tests {
		morph, ok := test.chain.RandomMorph(test.morphs)
		if test.ok != ok {
			t.Errorf("[%d] ok: expected %v, but got %v", idx, test.ok, ok)
		}
		if *test.morph != *morph {
			t.Errorf("[%d] morph: expected %v, but got %v", idx, *test.morph, *morph)
		}
	}
}
