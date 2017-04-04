package main

// Katakana is a table mapping katakana to Mora.
var Katakana = map[string]Mora{
	"ア": Mora{"", "a"}, "イ": Mora{"", "i"}, "ウ": Mora{"", "u"}, "エ": Mora{"", "e"}, "オ": Mora{"", "o"},
	"カ": Mora{"k", "a"}, "キ": Mora{"k", "i"}, "ク": Mora{"k", "u"}, "ケ": Mora{"k", "e"}, "コ": Mora{"k", "o"},
	"サ": Mora{"s", "a"}, "シ": Mora{"sh", "i"}, "ス": Mora{"s", "u"}, "セ": Mora{"s", "e"}, "ソ": Mora{"s", "o"},
	"タ": Mora{"t", "a"}, "チ": Mora{"ch", "i"}, "ツ": Mora{"ts", "u"}, "テ": Mora{"t", "e"}, "ト": Mora{"t", "o"},
	"ナ": Mora{"n", "a"}, "ニ": Mora{"n", "i"}, "ヌ": Mora{"n", "u"}, "ネ": Mora{"n", "e"}, "ノ": Mora{"n", "o"},
	"ハ": Mora{"h", "a"}, "ヒ": Mora{"h", "i"}, "フ": Mora{"f", "u"}, "ヘ": Mora{"h", "e"}, "ホ": Mora{"h", "o"},
	"マ": Mora{"m", "a"}, "ミ": Mora{"m", "i"}, "ム": Mora{"m", "u"}, "メ": Mora{"m", "e"}, "モ": Mora{"m", "o"},
	"ヤ": Mora{"y", "a"}, "ユ": Mora{"y", "u"}, "ヨ": Mora{"y", "o"},
	"ラ": Mora{"r", "a"}, "リ": Mora{"r", "i"}, "ル": Mora{"r", "u"}, "レ": Mora{"r", "e"}, "ロ": Mora{"r", "o"},
	"ワ": Mora{"w", "a"}, "ヲ": Mora{"", "o"}, "ン": Mora{"*n", "*n"},
	"ガ": Mora{"g", "a"}, "ギ": Mora{"g", "i"}, "グ": Mora{"g", "u"}, "ゲ": Mora{"g", "e"}, "ゴ": Mora{"g", "o"},
	"ザ": Mora{"z", "a"}, "ジ": Mora{"j", "i"}, "ズ": Mora{"z", "u"}, "ゼ": Mora{"z", "e"}, "ゾ": Mora{"z", "o"},
	"ダ": Mora{"d", "a"}, "ヂ": Mora{"j", "i"}, "ヅ": Mora{"z", "u"}, "デ": Mora{"d", "e"}, "ド": Mora{"d", "o"},
	"バ": Mora{"b", "a"}, "ビ": Mora{"b", "i"}, "ブ": Mora{"b", "u"}, "ベ": Mora{"b", "e"}, "ボ": Mora{"b", "o"},
	"パ": Mora{"p", "a"}, "ピ": Mora{"p", "i"}, "プ": Mora{"p", "u"}, "ペ": Mora{"p", "e"}, "ポ": Mora{"p", "o"},
	"キャ": Mora{"ky", "a"}, "キュ": Mora{"ky", "u"}, "キョ": Mora{"ky", "o"},
	"シャ": Mora{"sh", "a"}, "シュ": Mora{"sh", "u"}, "ショ": Mora{"sh", "o"},
	"チャ": Mora{"ch", "a"}, "チュ": Mora{"ch", "u"}, "チョ": Mora{"ch", "o"},
	"ニャ": Mora{"ny", "a"}, "ニュ": Mora{"ny", "u"}, "ニョ": Mora{"ny", "o"},
	"ヒャ": Mora{"hy", "a"}, "ヒュ": Mora{"hy", "u"}, "ヒョ": Mora{"hy", "o"},
	"ミャ": Mora{"my", "a"}, "ミュ": Mora{"my", "u"}, "ミョ": Mora{"my", "o"},
	"リャ": Mora{"ry", "a"}, "リュ": Mora{"ry", "u"}, "リョ": Mora{"ry", "o"},
	"ギャ": Mora{"gy", "a"}, "ギュ": Mora{"gy", "u"}, "ギョ": Mora{"gy", "o"},
	"ジャ": Mora{"j", "a"}, "ジュ": Mora{"j", "u"}, "ジョ": Mora{"j", "o"},
	"ビャ": Mora{"by", "a"}, "ビュ": Mora{"by", "u"}, "ビョ": Mora{"by", "o"},
	"ピャ": Mora{"py", "a"}, "ピュ": Mora{"py", "u"}, "ピョ": Mora{"py", "o"},
	"ファ": Mora{"f", "a"}, "フィ": Mora{"f", "i"}, "フェ": Mora{"f", "e"}, "フォ": Mora{"f", "o"},
	"フュ": Mora{"fy", "u"},
	"ウィ": Mora{"w", "i"}, "ウェ": Mora{"w", "e"}, "ウォ": Mora{"w", "o"},
	"ヴァ": Mora{"v", "a"}, "ヴィ": Mora{"v", "i"}, "ヴェ": Mora{"v", "e"}, "ヴォ": Mora{"v", "o"},
	"ツァ": Mora{"ts", "a"}, "ツィ": Mora{"ts", "i"}, "ツェ": Mora{"ts", "e"}, "ツォ": Mora{"ts", "o"},
	"チェ": Mora{"ch", "e"}, "シェ": Mora{"sh", "e"}, "ジェ": Mora{"j", "e"},
	"ティ": Mora{"t", "i"}, "ディ": Mora{"d", "i"},
	"デュ": Mora{"d", "u"}, "トゥ": Mora{"t", "u"},
	"ッ": Mora{"*xtu", "*xtu"},
}

// MoraDummy is a dummy mora.
var MoraDummy = Mora{"*", "*"}

// Mora is a property of mora.
type Mora struct{ Consonant, Vowel string }

// NewMorae converts s to []Mora.
func NewMorae(s string) ([]Mora, bool) {
	runes := []rune(s + "*")
	ans := make([]Mora, 0, len(runes)-1)
	for i := 0; i < len(runes)-1; i++ {
		if m, ok := Katakana[string(runes[i:i+2])]; ok {
			ans = append(ans, m)
			i++
		} else if m, ok := Katakana[string(runes[i])]; ok {
			ans = append(ans, m)
		} else if runes[i] == 'ー' && len(ans) > 0 {
			m := Mora{Consonant: "", Vowel: ans[len(ans)-1].Vowel}
			ans = append(ans, m)
		} else {
			return nil, false
		}
	}
	if len(ans) == 0 {
		return nil, false
	}
	return ans, true
}

// MoraWeight is a part of MoraeWeight.
type MoraWeight struct{ Consonant, Vowel float64 }

// NewMoraeWeight returns new MoraWeight instance from mw.
func NewMoraeWeight(mw []MoraWeight) *MoraeWeight {
	sum := 0.0
	for _, m := range mw {
		sum += m.Consonant
		sum += m.Vowel
	}
	return &MoraeWeight{MW: mw, sum: sum}
}

// MoraeWeight defines a weight of rhyming.
type MoraeWeight struct {
	MW  []MoraWeight
	sum float64
}

// Similarity returns coolness of rhyming on m0 + m1.
func (mw *MoraeWeight) Similarity(m0, m1 []Mora) float64 {
	minLen := len(mw.MW)
	for _, l := range []int{len(m0), len(m1)} {
		if minLen > l {
			minLen = l
		}
	}
	if minLen == 0 {
		return 0.0
	}

	weight := mw.MW[len(mw.MW)-minLen:]
	myM0 := m0[len(m0)-minLen:]
	myM1 := m1[len(m1)-minLen:]

	var sum float64
	for i := 0; i < minLen; i++ {
		if myM0[i] == MoraDummy || myM1[i] == MoraDummy {
			continue
		}
		if myM0[i].Consonant == myM1[i].Consonant {
			sum += weight[i].Consonant
		}
		if myM0[i].Vowel == myM1[i].Vowel {
			sum += weight[i].Vowel
		}
	}
	return sum / mw.sum
}

// SimMorphs returns similarity between ms0 and ms1.
func (mw *MoraeWeight) SimMorphs(ms0, ms1 Morphs) float64 {
	toMora := func(ms Morphs) []Mora {
		ans := make([]Mora, 0, len(ms))
		for _, m := range ms {
			morae, ok := NewMorae(m.Pronounciation)
			if !ok && m.Pronounciation == "" {
				morae = []Mora{MoraDummy}
			}
			ans = append(ans, morae...)
		}
		return ans
	}
	return mw.Similarity(toMora(ms0), toMora(ms1))
}
