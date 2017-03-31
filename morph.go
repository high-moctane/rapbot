package main

import (
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
	morphs := make(Morphs, len(tokens))
	for i, token := range tokens {
		morphs[i] = NewMorph(token)
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
