package main

import (
	"fmt"
	"strings"

	"github.com/ikawaha/kagome/tokenizer"
)

// Morph has a morpheme and its properties.
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
	Pronunciation        string
}

var (
	// BOS means Begin Of String
	BOS = Morph{"BOS", "", "", "", "", "", "", "", "", ""}
	// EOS means End Of String
	EOS = Morph{"EOS", "", "", "", "", "", "", "", "", ""}
)

// NewMorph returns new Morph.
func NewMorph(token *tokenizer.Token) *Morph {
	morph := Morph{}
	morph.Surface = token.Surface

	if token.Class == tokenizer.DUMMY {
		return &morph
	}

	for i, feature := range token.Features() {
		switch i {
		case 0:
			morph.PartOfSpeech = feature
		case 1:
			morph.PartOfSpeechSection1 = feature
		case 2:
			morph.PartOfSpeechSection2 = feature
		case 3:
			morph.PartOfSpeechSection3 = feature
		case 4:
			morph.ConjugatedForm1 = feature
		case 5:
			morph.ConjugatedForm2 = feature
		case 6:
			morph.Inflection = feature
		case 7:
			morph.Reading = feature
		case 8:
			morph.Pronunciation = feature
		}
	}
	return &morph
}

// Morae returns morae. If has unreadable morph, ok will be false.
func (m *Morph) Morae() (morae Morae, ok bool) {
	morae, ok = NewMorae(m.Pronunciation)
	return
}

func (m *Morph) String() string {
	return "[" + fmt.Sprintf("%#v", m.Surface) + " " +
		strings.Join([]string{
			m.PartOfSpeech,
			m.PartOfSpeechSection1,
			m.PartOfSpeechSection2,
			m.PartOfSpeechSection3,
			m.ConjugatedForm1,
			m.ConjugatedForm2,
			m.Inflection,
			m.Reading,
			m.Pronunciation,
		}, " ") + "]"
}

var katakana = map[string]*Mora{
	"ア": &Mora{"", "a"}, "イ": &Mora{"", "i"}, "ウ": &Mora{"", "u"}, "エ": &Mora{"", "e"}, "オ": &Mora{"", "o"},
	"カ": &Mora{"k", "a"}, "キ": &Mora{"k", "i"}, "ク": &Mora{"k", "u"}, "ケ": &Mora{"k", "e"}, "コ": &Mora{"k", "o"},
	"サ": &Mora{"s", "a"}, "シ": &Mora{"sh", "i"}, "ス": &Mora{"s", "u"}, "セ": &Mora{"s", "e"}, "ソ": &Mora{"s", "o"},
	"タ": &Mora{"t", "a"}, "チ": &Mora{"ch", "i"}, "ツ": &Mora{"ts", "u"}, "テ": &Mora{"t", "e"}, "ト": &Mora{"t", "o"},
	"ナ": &Mora{"n", "a"}, "ニ": &Mora{"n", "i"}, "ヌ": &Mora{"n", "u"}, "ネ": &Mora{"n", "e"}, "ノ": &Mora{"n", "o"},
	"ハ": &Mora{"h", "a"}, "ヒ": &Mora{"h", "i"}, "フ": &Mora{"f", "u"}, "ヘ": &Mora{"h", "e"}, "ホ": &Mora{"h", "o"},
	"マ": &Mora{"m", "a"}, "ミ": &Mora{"m", "i"}, "ム": &Mora{"m", "u"}, "メ": &Mora{"m", "e"}, "モ": &Mora{"m", "o"},
	"ヤ": &Mora{"y", "a"}, "ユ": &Mora{"y", "u"}, "ヨ": &Mora{"y", "o"},
	"ラ": &Mora{"r", "a"}, "リ": &Mora{"r", "i"}, "ル": &Mora{"r", "u"}, "レ": &Mora{"r", "e"}, "ロ": &Mora{"r", "o"},
	"ワ": &Mora{"w", "a"}, "ヲ": &Mora{"", "o"}, "ン": &Mora{"*n", "*n"},
	"ガ": &Mora{"g", "a"}, "ギ": &Mora{"g", "i"}, "グ": &Mora{"g", "u"}, "ゲ": &Mora{"g", "e"}, "ゴ": &Mora{"g", "o"},
	"ザ": &Mora{"z", "a"}, "ジ": &Mora{"j", "i"}, "ズ": &Mora{"z", "u"}, "ゼ": &Mora{"z", "e"}, "ゾ": &Mora{"z", "o"},
	"ダ": &Mora{"d", "a"}, "ヂ": &Mora{"j", "i"}, "ヅ": &Mora{"z", "u"}, "デ": &Mora{"d", "e"}, "ド": &Mora{"d", "o"},
	"バ": &Mora{"b", "a"}, "ビ": &Mora{"b", "i"}, "ブ": &Mora{"b", "u"}, "ベ": &Mora{"b", "e"}, "ボ": &Mora{"b", "o"},
	"パ": &Mora{"p", "a"}, "ピ": &Mora{"p", "i"}, "プ": &Mora{"p", "u"}, "ペ": &Mora{"p", "e"}, "ポ": &Mora{"p", "o"},
	"キャ": &Mora{"ky", "a"}, "キュ": &Mora{"ky", "u"}, "キョ": &Mora{"ky", "o"},
	"シャ": &Mora{"sh", "a"}, "シュ": &Mora{"sh", "u"}, "ショ": &Mora{"sh", "o"},
	"チャ": &Mora{"ch", "a"}, "チュ": &Mora{"ch", "u"}, "チョ": &Mora{"ch", "o"},
	"ニャ": &Mora{"ny", "a"}, "ニュ": &Mora{"ny", "u"}, "ニョ": &Mora{"ny", "o"},
	"ヒャ": &Mora{"hy", "a"}, "ヒュ": &Mora{"hy", "u"}, "ヒョ": &Mora{"hy", "o"},
	"ミャ": &Mora{"my", "a"}, "ミュ": &Mora{"my", "u"}, "ミョ": &Mora{"my", "o"},
	"リャ": &Mora{"ry", "a"}, "リュ": &Mora{"ry", "u"}, "リョ": &Mora{"ry", "o"},
	"ギャ": &Mora{"gy", "a"}, "ギュ": &Mora{"gy", "u"}, "ギョ": &Mora{"gy", "o"},
	"ジャ": &Mora{"j", "a"}, "ジュ": &Mora{"j", "u"}, "ジョ": &Mora{"j", "o"},
	"ビャ": &Mora{"by", "a"}, "ビュ": &Mora{"by", "u"}, "ビョ": &Mora{"by", "o"},
	"ピャ": &Mora{"py", "a"}, "ピュ": &Mora{"py", "u"}, "ピョ": &Mora{"py", "o"},
	"ファ": &Mora{"f", "a"}, "フィ": &Mora{"f", "i"}, "フェ": &Mora{"f", "e"}, "フォ": &Mora{"f", "o"},
	"フュ": &Mora{"fy", "u"},
	"ウィ": &Mora{"w", "i"}, "ウェ": &Mora{"w", "e"}, "ウォ": &Mora{"w", "o"},
	"ヴァ": &Mora{"v", "a"}, "ヴィ": &Mora{"v", "i"}, "ヴェ": &Mora{"v", "e"}, "ヴォ": &Mora{"v", "o"},
	"ツァ": &Mora{"ts", "a"}, "ツィ": &Mora{"ts", "i"}, "ツェ": &Mora{"ts", "e"}, "ツォ": &Mora{"ts", "o"},
	"チェ": &Mora{"ch", "e"}, "シェ": &Mora{"sh", "e"}, "ジェ": &Mora{"j", "e"},
	"ティ": &Mora{"t", "i"}, "ディ": &Mora{"d", "i"},
	"デュ": &Mora{"d", "u"}, "トゥ": &Mora{"t", "u"},
	"ッ": &Mora{"*xtu", "*xtu"},
}

// Mora consists of consonant and vowel.
type Mora struct {
	consonant, vowel string
}

// NewMora returns a mora corresponding to kana.
// If there is no suitable mora, ok will be false and mora will be non-nil
// empty value.
func NewMora(kana string) (mora *Mora, ok bool) {
	mora, ok = katakana[kana]
	if !ok {
		mora = &Mora{}
	}
	return
}

func (m *Mora) String() string {
	return "[" + m.consonant + " " + m.vowel + "]"
}

// Morae is a slice of Mora
type Morae []*Mora

// NewMorae returns new morae from katakana pronunciation.
// If cannot build morae completely, ok will be false.
func NewMorae(pronunciation string) (morae Morae, ok bool) {
	runes := append([]rune(pronunciation), '*') // "*" is dummy rune

	for i := 0; i < len(runes)-1; i++ {
		if mora, ok2 := NewMora(string(runes[i : i+2])); ok2 {
			// 拗音
			morae = append(morae, mora)
			i++
		} else if mora, ok2 := NewMora(string(runes[i])); ok2 {
			morae = append(morae, mora)
		} else if len(morae) > 0 && runes[i] == 'ー' {
			mora := &Mora{"", morae[len(morae)-1].vowel}
			morae = append(morae, mora)
		} else {
			return
		}
	}
	ok = len(morae) > 0
	return
}

// Sentence is a sentence which consist of Morph arrays.
type Sentence []*Morph

// IsPronounceable returns the sentence can pronounce.
func (se Sentence) IsPronounceable() bool {
	for _, morph := range se {
		if _, ok := morph.Morae(); !ok {
			return false
		}
	}
	return true
}

// Morae returns sentence's morae.
func (se Sentence) Morae() (morae Morae, ok bool) {
	for _, morph := range se {
		m, ok2 := morph.Morae()
		if !ok2 {
			return
		}
		morae = append(morae, m...)
	}
	ok = len(morae) > 0
	return
}

// String joins all surfaces.
func (se Sentence) String() string {
	builder := new(strings.Builder)
	for _, morph := range se {
		builder.WriteString(morph.Surface)
	}
	return builder.String()
}

// Lyric is Lyric
type Lyric []Sentence

func (l Lyric) String() string {
	strs := []string{}
	for _, line := range l {
		strs = append(strs, line.String())
	}
	return strings.Join(strs, "\n")
}

// JapaneseParseServer parse chTweets text and send it to ChMorphs.
func JapaneseParseServer(chSentence chan<- Sentence, chString <-chan string) {
	t := tokenizer.New()

	for text := range chString {
		chSentence <- analyzeText(&t, text)
	}
}

// analyzeText analyzes text into Sentence.
func analyzeText(t *tokenizer.Tokenizer, text string) Sentence {
	tokens := t.Tokenize(text)
	sentence := make(Sentence, 0, len(tokens))
	for _, token := range tokens {
		morph := NewMorph(&token)
		sentence = append(sentence, morph)
	}
	return sentence
}
