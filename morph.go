package main

import (
	"fmt"
	"strings"

	"github.com/ikawaha/kagome/tokenizer"
)

// NewMorph makes new Morph from token.
func NewMorph(token tokenizer.Token) *Morph {
	// BOS or EOS
	if token.Class == tokenizer.DUMMY {
		if token.Surface == "BOS" {
			return &MorphBOS
		}
		return &MorphEOS
	}

	// ordinary morph
	m := &Morph{Surface: token.Surface}
	for i, f := range token.Features() {
		switch i {
		case 0:
			m.PartOfSpeech = f
		case 1:
			m.PartOfSpeechSection1 = f
		case 2:
			m.PartOfSpeechSection2 = f
		case 3:
			m.PartOfSpeechSection3 = f
		case 4:
			m.ConjugatedForm1 = f
		case 5:
			m.ConjugatedForm2 = f
		case 6:
			m.Inflection = f
		case 7:
			m.Reading = f
		case 8:
			m.Pronounciation = f
		}
	}
	return m
}

// Morph stores morpheme's properties.
type Morph struct {
	Surface              string
	PartOfSpeech         string
	PartOfSpeechSection1 string
	PartOfSpeechSection2 string
	PartOfSpeechSection3 string
	ConjugatedForm1      string
	ConjugatedForm2      string
	Inflection           string
	Reading              string
	Pronounciation       string
}

var (
	// MorphBOS is BOS Morph.
	MorphBOS = Morph{PartOfSpeech: "BOS"}

	// MorphEOS is EOS Morph.
	MorphEOS = Morph{PartOfSpeech: "EOS"}
)

func (m *Morph) String() string {
	return m.Surface + "\t" + strings.Join([]string{
		m.PartOfSpeech,
		m.PartOfSpeechSection1,
		m.PartOfSpeechSection2,
		m.PartOfSpeechSection3,
		m.ConjugatedForm1,
		m.ConjugatedForm2,
		m.Inflection,
		m.Reading,
		m.Pronounciation,
	}, ",")
}

// NewMorphs makes new Morphs from tokens.
func NewMorphs(tokens []tokenizer.Token) Morphs {
	morphs := make(Morphs, 0, len(tokens))
	for _, token := range tokens {
		if m := NewMorph(token); *m != MorphBOS && *m != MorphEOS {
			morphs = append(morphs, m)
		}
	}
	return morphs
}

// Morphs is a slice of Morph.
type Morphs []*Morph

// IsEqual reports whether ms and a have same Morph(s).
func (ms Morphs) IsEqual(a Morphs) bool {
	if len(ms) != len(a) {
		return false
	}
	for i := range ms {
		if *ms[i] != *a[i] {
			return false
		}
	}
	return true
}

func (ms Morphs) String() string {
	ans := []string{}
	for _, m := range ms {
		ans = append(ans, m.String())
	}
	return fmt.Sprint(ans)
}

// Surface returns ms's surface.
// When ms contains empty surface, it returns false.
func (ms Morphs) Surface() (string, bool) {
	ans := make([]string, 0, len(ms))
	ok := true
	for _, m := range ms {
		if m.Surface == "" {
			ok = false
		}
		ans = append(ans, m.Surface)
	}
	return strings.Join(ans, ""), ok
}

// Pronounciation returns ms's pronounciation.
// When ms contains empty pronounciation, it returns false.
func (ms Morphs) Pronounciation() (string, bool) {
	ans := make([]string, 0, len(ms))
	ok := true
	for _, m := range ms {
		if m.Pronounciation == "" {
			ok = false
		}
		ans = append(ans, m.Pronounciation)
	}
	return strings.Join(ans, ""), ok
}
