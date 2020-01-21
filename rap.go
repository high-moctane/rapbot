package main

import (
	"container/list"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"sync"
)

// Weight is a similarity weight for a mora.
type Weight struct {
	consonant, vowel float64
}

// Rapper makes nice lyrics.
type Rapper struct {
	weights   []Weight
	maxWeight float64
	thresh    float64
	tryNum    int
}

// DefaultRapper returns default rapper
func DefaultRapper() (*Rapper, error) {
	weights, err := parseWeights()
	if err != nil {
		return nil, fmt.Errorf("cannot create rapper: weights: %w", err)
	}
	thresh, err := strconv.ParseFloat(os.Getenv("THRESH"), 64)
	if err != nil {
		return nil, fmt.Errorf("cannot create rapper: thresh: %w", err)
	}
	tryNum, err := strconv.Atoi(os.Getenv("TRY_NUM"))
	if err != nil {
		return nil, fmt.Errorf("cannot create rapper: tryNum: %w", err)
	}

	var maxWeight float64
	for _, weight := range weights {
		maxWeight += weight.consonant
		maxWeight += weight.vowel
	}

	return &Rapper{
		weights:   weights,
		maxWeight: maxWeight,
		thresh:    thresh,
		tryNum:    tryNum,
	}, nil
}

func parseWeights() (weights []Weight, err error) {
	for _, str := range strings.Split(os.Getenv("CONSONANT_WEIGHTS"), ",") {
		var val float64
		val, err = strconv.ParseFloat(str, 64)
		if err != nil {
			return
		}
		weights = append(weights, Weight{consonant: val})
	}
	for i, str := range strings.Split(os.Getenv("VOWEL_WEIGHTS"), ",") {
		var val float64
		val, err = strconv.ParseFloat(str, 64)
		if err != nil {
			return
		}
		if len(weights) < i+1 {
			err = errors.New("CONSONANT_WEIGHTS len and VOWEL_WEIGHTS are not equal")
			return
		}
		weights[i].vowel = val
	}
	return
}

// LaunchRapServer launches multi RapServer.
func (rap *Rapper) LaunchRapServer(chLyric chan<- Lyric, chSentence <-chan Sentence) error {
	for _, str := range strings.Split(os.Getenv("LYRIC_LINE_NUM"), ",") {
		val, err := strconv.Atoi(str)
		if err != nil {
			return fmt.Errorf("invalid LYRIC_LINE_NUM: %w", err)
		}

		go rap.RapServer(chLyric, chSentence, val)
	}

	return nil
}

// RapServer make lyrics forever.
func (rap *Rapper) RapServer(chLyric chan<- Lyric, chSentence <-chan Sentence, lineNum int) {
mainLoop:
	for {
		lyric := []Sentence{<-chSentence}
	lyricLoop:
		for len(lyric) < lineNum {
			for try := 0; try < rap.tryNum; try++ {
				// find pronounceable sentence
				var sentence Sentence
				for {
					sentence = <-chSentence
					if isValidRapSentence(sentence) {
						break
					}
				}

				// judge the lyric is valid
				if rap.IsAppendable(lyric, sentence) {
					lyric = append(lyric, sentence)
					continue lyricLoop
				}
			}
			continue mainLoop
		}
		chLyric <- lyric
	}
}

// isValidRapSentence returns whether the sentence is valid for lyric.
func isValidRapSentence(sentence Sentence) bool {
	return true && // for easy comment out
		sentence.IsPronounceable() &&
		sentence[len(sentence)-1].ConjugatedForm2 != "連用タ接続" && // 「なかっ」
		sentence[len(sentence)-1].ConjugatedForm2 != "連用形" && // 「（ありがとう）ござい」
		sentence[len(sentence)-1].ConjugatedForm2 != "未然形" && // 「い（ない）」
		sentence[len(sentence)-1].PartOfSpeech != "助詞" && // 「〇〇の」
		(sentence[len(sentence)-1].PartOfSpeechSection1 != "接尾" ||
			sentence[len(sentence)-1].PartOfSpeechSection2 != "人名") && // 「〇〇さん」
		true // for easy comment out
}

// IsAppendable returns if sentence is suitable for lyric
func (rap *Rapper) IsAppendable(lyric Lyric, sentence Sentence) bool {
	// rhyming
	if rap.Distance(lyric[len(lyric)-1], sentence) < rap.thresh {
		return false
	}

	// similarity
	// If sentences ends same character, return false.
	lastSentenceRune := []rune(lyric[len(lyric)-1].String())
	newSentenceRune := []rune(sentence.String())
	if lastSentenceRune[len(lastSentenceRune)-1] == newSentenceRune[len(newSentenceRune)-1] {
		return false
	}

	return true
}

// Distance is a similarity of two sentence.
func (rap *Rapper) Distance(sen1, sen2 Sentence) float64 {
	morae1, ok := sen1.Morae()
	if !ok {
		return 0.0
	}
	morae2, ok := sen2.Morae()
	if !ok {
		return 0.0
	}

	minLength := len(rap.weights)
	if len(morae1) < minLength {
		minLength = len(morae1)
	}
	if len(morae2) < minLength {
		minLength = len(morae2)
	}

	var sum float64
	for i := 0; i < minLength; i++ {
		weight := rap.weights[len(rap.weights)-1-i]
		mora1 := morae1[len(morae1)-1-i]
		mora2 := morae2[len(morae2)-1-i]

		if mora1.consonant == mora2.consonant {
			sum += weight.consonant
		}
		if mora1.vowel == mora2.vowel {
			sum += weight.vowel
		}
	}

	return sum / rap.maxWeight
}

// LyricStorage stores rhymes.
type LyricStorage struct {
	maxLen int
	length int
	mu     *sync.Mutex
	lyrics *list.List
}

// NewLyricStorage returns new LyricStorage.
func NewLyricStorage(maxLen int) *LyricStorage {
	return &LyricStorage{
		maxLen: maxLen,
		mu:     new(sync.Mutex),
		lyrics: list.New(),
	}
}

// PushServer receive lyrics forever.
func (ls *LyricStorage) PushServer(chLyric <-chan Lyric) {
	for lyric := range chLyric {
		ls.Push(lyric)
	}
}

// Push adds lyric.
func (ls *LyricStorage) Push(lyric Lyric) {
	ls.mu.Lock()
	defer ls.mu.Unlock()

	ls.lyrics.PushFront(lyric)
	ls.length++
	if ls.length > ls.maxLen {
		ls.lyrics.Remove(ls.lyrics.Back())
		ls.length--
	}
}

// Pop returns newest lyric
func (ls *LyricStorage) Pop() Lyric {
	ls.mu.Lock()
	defer ls.mu.Unlock()

	if ls.length == 0 {
		return nil
	}

	lyric := ls.lyrics.Front().Value.(Lyric)
	ls.lyrics.Remove(ls.lyrics.Front())
	ls.length--
	return lyric
}

// ContinueLyric returns most suitable lyric.
func (ls *LyricStorage) ContinueLyric(rapper *Rapper, sentence Sentence) Lyric {
	ls.mu.Lock()
	defer ls.mu.Unlock()

	if ls.length == 0 {
		return nil
	}

	most := ls.lyrics.Front()
	distance := rapper.Distance(sentence, most.Value.(Lyric)[0])
	for e := ls.lyrics.Front().Next(); e != nil; e = e.Next() {
		if d := rapper.Distance(sentence, e.Value.(Lyric)[0]); d > distance {
			most = e
			distance = d
		}
	}
	ls.lyrics.Remove(most)
	ls.length--
	return most.Value.(Lyric)
}
