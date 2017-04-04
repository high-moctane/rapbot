package main

import (
	"testing"

	"reflect"

	"sync"

	"github.com/k0kubun/pp"
)

func TestChainAdd(t *testing.T) {
	tests := []struct {
		input []Morphs
		want  *chain
	}{
		{
			input: []Morphs{
				Morphs{
					{Surface: "a"},
					{Surface: "b"},
				},
			},
			want: &chain{
				c: chainMap{
					{Surface: "a"}: &chain{c: chainMap{{Surface: "b"}: nil}},
				},
			},
		},
		{
			input: []Morphs{
				Morphs{
					{Surface: "a"},
					{Surface: "b"},
				},
				Morphs{
					{Surface: "a"},
					{Surface: "b"},
				},
			},
			want: &chain{
				c: chainMap{
					{Surface: "a"}: &chain{c: chainMap{{Surface: "b"}: nil}},
				},
			},
		},
		{
			input: []Morphs{
				Morphs{
					{Surface: "a"},
					{Surface: "b"},
				},
				Morphs{
					{Surface: "a"},
					{Surface: "c"},
				},
			},
			want: &chain{
				c: chainMap{
					{Surface: "a"}: &chain{c: chainMap{
						{Surface: "b"}: nil,
						{Surface: "c"}: nil,
					}},
				},
			},
		},
		{
			input: []Morphs{
				Morphs{
					{Surface: "a"},
					{Surface: "b"},
				},
				Morphs{
					{Surface: "c"},
					{Surface: "b"},
				},
			},
			want: &chain{
				c: chainMap{
					{Surface: "a"}: &chain{c: chainMap{{Surface: "b"}: nil}},
					{Surface: "c"}: &chain{c: chainMap{{Surface: "b"}: nil}},
				},
			},
		},
		{
			input: []Morphs{
				Morphs{
					{Surface: "a"},
					{Surface: "b"},
					{Surface: "c"},
				},
			},
			want: &chain{c: chainMap{
				{Surface: "a"}: &chain{c: chainMap{
					{Surface: "b"}: &chain{c: chainMap{
						{Surface: "c"}: nil}},
				}},
			}},
		},
		{
			input: []Morphs{
				Morphs{
					{Surface: "A"},
					{Surface: "a"},
					{Surface: "0"},
				},
				Morphs{
					{Surface: "A"},
					{Surface: "a"},
					{Surface: "1"},
				},
				Morphs{
					{Surface: "A"},
					{Surface: "b"},
					{Surface: "0"},
				},
				Morphs{
					{Surface: "B"},
					{Surface: "a"},
					{Surface: "0"},
				},
			},
			want: &chain{
				c: chainMap{
					{Surface: "A"}: &chain{c: chainMap{
						{Surface: "a"}: &chain{c: chainMap{
							{Surface: "0"}: nil,
							{Surface: "1"}: nil,
						}},
						{Surface: "b"}: &chain{c: chainMap{
							{Surface: "0"}: nil,
						}},
					}},
					{Surface: "B"}: &chain{c: chainMap{
						{Surface: "a"}: &chain{c: chainMap{
							{Surface: "0"}: nil,
						}},
					}},
				},
			},
		},
	}

	for _, test := range tests {
		c := newChain()
		for _, ms := range test.input {
			c.add(ms)
		}
		if !reflect.DeepEqual(c, test.want) {
			t.Errorf("chain = %v, want %v",
				pp.Sprint(c), pp.Sprint(test.want))
		}
	}
}

func TestChainFindRand(t *testing.T) {
	type answer struct {
		m  *Morph
		ok bool
	}
	tests := []struct {
		input Morphs
		want  []answer
	}{
		{
			input: Morphs{{Surface: "A"}},
			want:  []answer{{m: &Morph{Surface: "a"}, ok: true}},
		},
		{
			input: Morphs{{Surface: "A"}, {Surface: "a"}},
			want:  []answer{{m: &Morph{Surface: "0"}, ok: true}},
		},
		{
			input: Morphs{{Surface: "B"}},
			want: []answer{
				{m: &Morph{Surface: "a"}, ok: true},
				{m: &Morph{Surface: "b"}, ok: true},
			},
		},
		{
			input: Morphs{{Surface: "C"}},
			want:  []answer{{m: nil, ok: false}},
		},
		{
			input: Morphs{{Surface: "A"}, {Surface: "c"}},
			want:  []answer{{m: nil, ok: false}},
		},
	}

	c := &chain{
		c: chainMap{
			{Surface: "A"}: &chain{c: chainMap{
				{Surface: "a"}: &chain{c: chainMap{
					{Surface: "0"}: nil,
				}},
			}},
			{Surface: "B"}: &chain{c: chainMap{
				{Surface: "a"}: &chain{c: chainMap{
					{Surface: "0"}: nil,
					{Surface: "1"}: nil,
				}},
				{Surface: "b"}: &chain{c: chainMap{
					{Surface: "0"}: nil,
					{Surface: "1"}: nil,
				}},
			}},
		},
	}

loop:
	for _, test := range tests {
		ans := answer{}
		ans.m, ans.ok = c.findRand(test.input)
		for _, a := range test.want {
			if reflect.DeepEqual(ans, a) {
				continue loop
			}
		}
		t.Errorf("(*chain).find(%v) = %v, want one of %v",
			test.input, pp.Sprint(ans), pp.Sprint(test.want))
	}
}

func TestMarkovAdd(t *testing.T) {
	type Input struct {
		param MarkovParam
		mss   []Morphs
	}
	tests := []struct {
		input Input
		want  *markov
	}{
		{
			input: Input{
				param: MarkovParam{n: 2, lcs: 2, lc: 3},
				mss: []Morphs{
					{{Surface: "a"}},
					{{Surface: "b"}},
					{{Surface: "c"}},
					{{Surface: "d"}},
					{{Surface: "e"}},
					{{Surface: "f"}},
					{{Surface: "g"}},
					{{Surface: "h"}},
					{{Surface: "i"}},
					{{Surface: "j"}},
				},
			},
			want: &markov{
				param: &MarkovParam{n: 2, lcs: 2, lc: 3},
				cs: []*chain{
					&chain{count: 3, c: chainMap{
						MorphBOS: &chain{c: chainMap{
							{Surface: "d"}: nil,
							{Surface: "e"}: nil,
							{Surface: "f"}: nil,
						}},
						{Surface: "d"}: &chain{c: chainMap{MorphEOS: nil}},
						{Surface: "e"}: &chain{c: chainMap{MorphEOS: nil}},
						{Surface: "f"}: &chain{c: chainMap{MorphEOS: nil}},
					}},
					&chain{count: 3, c: chainMap{
						MorphBOS: &chain{c: chainMap{
							{Surface: "g"}: nil,
							{Surface: "h"}: nil,
							{Surface: "i"}: nil,
						}},
						{Surface: "g"}: &chain{c: chainMap{MorphEOS: nil}},
						{Surface: "h"}: &chain{c: chainMap{MorphEOS: nil}},
						{Surface: "i"}: &chain{c: chainMap{MorphEOS: nil}},
					}},
				},
				learning: &chain{count: 1, c: chainMap{
					MorphBOS: &chain{c: chainMap{
						{Surface: "j"}: nil,
					}},
					{Surface: "j"}: &chain{c: chainMap{MorphEOS: nil}},
				}},
				ready: nil,
			},
		},
	}

	for _, test := range tests {
		m := newMarkov(&test.input.param)
		for _, ms := range test.input.mss {
			m.add(ms)
		}
		<-m.ready
		m.once = sync.Once{}
		m.ready = nil
		if !reflect.DeepEqual(m, test.want) {
			t.Errorf("(*markov).add(%v) => m: %v, want %v",
				pp.Sprint(test.input), pp.Sprint(m), pp.Sprint(test.want))
		}
	}
}

func TestMarkovGenerate(t *testing.T) {
	mss := []Morphs{
		{{Surface: "a"}, {Surface: "a"}},
		{{Surface: "a"}, {Surface: "b"}},
		{{Surface: "a"}, {Surface: "c"}},
		{{Surface: "b"}, {Surface: "a"}},
		{{Surface: "b"}, {Surface: "b"}},
		{{Surface: "b"}, {Surface: "c"}},
		{{Surface: "c"}, {Surface: "a"}},
		{{Surface: "c"}, {Surface: "b"}},
		{{Surface: "c"}, {Surface: "c"}},
	}
	param := MarkovParam{n: 2, lcs: 3, lc: 3, lms: 2, try: 100}
	m := newMarkov(&param)
	for _, ms := range mss {
		m.add(ms)
	}
	ans := make([]Morphs, 0, 1000)
	for i := 0; i < 1000; i++ {
		if ms, ok := m.generate(); ok {
			ans = append(ans, ms)
		}
	}
	want := []Morphs{
		{{Surface: "a"}},
		{{Surface: "b"}},
		{{Surface: "c"}},
		{{Surface: "a"}, {Surface: "a"}},
		{{Surface: "a"}, {Surface: "b"}},
		{{Surface: "a"}, {Surface: "c"}},
		{{Surface: "b"}, {Surface: "a"}},
		{{Surface: "b"}, {Surface: "b"}},
		{{Surface: "b"}, {Surface: "c"}},
		{{Surface: "c"}, {Surface: "a"}},
		{{Surface: "c"}, {Surface: "b"}},
		{{Surface: "c"}, {Surface: "c"}},
	}

testLoop:
	for _, a := range ans {
		for _, w := range want {
			if reflect.DeepEqual(a, w) {
				continue testLoop
			}
		}
		t.Errorf("ans cannot contain %v, but it did.", pp.Sprint(a))
		break testLoop
	}
}

func TestMarkovServer(t *testing.T) {
	param := MarkovParam{n: 2, lcs: 3, lc: 3, lms: 2, try: 100}
	in := make(chan Morphs)
	out := MarkovServer(&param, in)

	input := []Morphs{
		{{Surface: "a"}, {Surface: "a"}},
		{{Surface: "a"}, {Surface: "b"}},
		{{Surface: "a"}, {Surface: "c"}},
		{{Surface: "b"}, {Surface: "a"}},
		{{Surface: "b"}, {Surface: "b"}},
		{{Surface: "b"}, {Surface: "c"}},
		{{Surface: "c"}, {Surface: "a"}},
		{{Surface: "c"}, {Surface: "b"}},
		{{Surface: "c"}, {Surface: "c"}},
		{{Surface: "a"}, {Surface: "a"}},
		{{Surface: "a"}, {Surface: "b"}},
		{{Surface: "a"}, {Surface: "c"}},
		{{Surface: "b"}, {Surface: "a"}},
		{{Surface: "b"}, {Surface: "b"}},
		{{Surface: "b"}, {Surface: "c"}},
		{{Surface: "c"}, {Surface: "a"}},
		{{Surface: "c"}, {Surface: "b"}},
		{{Surface: "c"}, {Surface: "c"}},
	}

	go func() {
		for _, ms := range input {
			in <- ms
		}
		close(in)
	}()

	ans := []Morphs{}
	for ms := range out {
		ans = append(ans, ms)
		if len(ans) > 100 {
			return
		}
	}
}
