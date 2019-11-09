package main

import (
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"sync"
)

// MarkovParams is a parameter of a markov chain.
type MarkovParams struct {
	Ngram          int // ngram (n >= 2).
	ChainNum       int // max number of markov chains.
	ChainMorphsNum int // max number of morphemes which each chain has.
}

// DefaultMarkovParams uses .env values.
func DefaultMarkovParams() (*MarkovParams, error) {
	ngram, err := envAtoiErr("NGRAM")
	if err != nil {
		return nil, err
	}
	chainNum, err := envAtoiErr("CHAIN_NUM")
	if err != nil {
		return nil, err
	}
	chainMorphsNum, err := envAtoiErr("CHAIN_MORPHS_NUM")
	if err != nil {
		return nil, err
	}

	return &MarkovParams{
		Ngram:          ngram,
		ChainNum:       chainNum,
		ChainMorphsNum: chainMorphsNum,
	}, nil
}

func envAtoiErr(name string) (int, error) {
	val, err := strconv.Atoi(os.Getenv(name))
	if err != nil {
		return 0, fmt.Errorf("invalid %v: %w", name, err)
	}
	return val, nil
}

// Markov has Markov chains. It can generate random sentences.
type Markov struct {
	once     *sync.Once
	Ready    chan struct{} // close ready when learning completed.
	params   *MarkovParams
	learning chain // under learning chain
	mu       *sync.RWMutex
	chains   []chain // Markov chains
}

// NewMarkov returns new Markov.
func NewMarkov(params *MarkovParams) *Markov {
	return &Markov{
		once:     new(sync.Once),
		Ready:    make(chan struct{}),
		params:   params,
		learning: make(chain),
		mu:       new(sync.RWMutex),
	}
}

// AddServer build Markov chains from ch.
func (m *Markov) AddServer(ch <-chan Sentence) {
	for se := range ch {
		m.Add(se)
	}
}

// Add adds sentence to Markov learning chain. This function cannot be called
// concurrently.
func (m *Markov) Add(sentence Sentence) {
	for i := 0; i < len(sentence)-m.params.Ngram+1; i++ {
		morphs := sentence[i : i+m.params.Ngram]
		m.learning.Add(morphs)

		if len(m.learning) >= m.params.ChainMorphsNum {
			m.shiftChain()
		}
	}
}

// shiftChain shift Markov chains and initialize learning.
func (m *Markov) shiftChain() {
	m.mu.Lock()
	if len(m.chains) >= m.params.ChainNum {
		m.chains = m.chains[1:]
	}
	m.chains = append(m.chains, m.learning)
	if len(m.chains) >= m.params.ChainNum {
		m.once.Do(func() { close(m.Ready) })
	}
	m.mu.Unlock()

	m.learning = make(chain)
}

// RandomSentenceServer generate random sentence forever.
func (m *Markov) RandomSentenceServer(chSentence chan<- Sentence, morphLen int) {
	for {
		sentence, ok := m.RandomSentence(morphLen)
		if !ok {
			continue
		}
		chSentence <- sentence
	}
}

// LaunchRandomSentenceServer launch multi RandomSentenceServer
func (m *Markov) LaunchRandomSentenceServer(chSentence chan<- Sentence) error {
	var morphLens []int
	for _, str := range strings.Split(os.Getenv("RANDOM_MORPH_LEN"), ",") {
		val, err := strconv.Atoi(str)
		if err != nil {
			return fmt.Errorf("invalid RANDOM_MORPH_LEN")
		}
		morphLens = append(morphLens, val)
	}

	for _, morphLen := range morphLens {
		go func(morphLen int) {
			<-markov.Ready
			markov.RandomSentenceServer(chSentence, morphLen)
		}(morphLen)
	}
	return nil
}

// RandomSentence generates random sentence
func (m *Markov) RandomSentence(morphLen int) (sentence Sentence, ok bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	// generate head of sentence
	sentence = []*Morph{&BOS}
	for i := 0; i < m.params.Ngram-2; i++ {
		var morph *Morph
		morph, ok = m.RandomMorph(sentence)
		if !ok {
			return
		}
		sentence = append(sentence, morph)
	}

	for len(sentence) < morphLen+1 {
		var morph *Morph
		morph, ok = m.RandomMorph(sentence[len(sentence)-m.params.Ngram+1:])
		if !ok {
			return
		}
		if *morph == EOS {
			break
		}
		sentence = append(sentence, morph)
	}

	return sentence[1:], true
}

// RandomMorph find random morph from all chains.
func (m *Markov) RandomMorph(morphs []*Morph) (morph *Morph, ok bool) {
	for _, idx := range randomIndice(len(m.chains)) {
		chain := m.chains[idx]
		morph, ok = chain.RandomMorph(morphs)
		if !ok {
			continue
		}
		return
	}

	return
}

// randomIndice generate random indice.
func randomIndice(num int) []int {
	indice := make([]int, num)
	for i := 0; i < num; i++ {
		indice[i] = i
	}
	for i := num - 1; i > 0; i-- {
		j := rand.Intn(i + 1)
		indice[i], indice[j] = indice[j], indice[i]
	}
	return indice
}

// chain is a markov chain map.
type chain map[Morph]chain

// Add adds morphs recursively.
func (c *chain) Add(morphs []*Morph) {
	if len(morphs) == 0 {
		return
	}

	next, ok := (*c)[*morphs[0]]
	if !ok {
		next = make(chain)
		(*c)[*morphs[0]] = next
	}

	(&next).Add(morphs[1:])
}

// RandomMorph returns random Morph from morphs.
func (c chain) RandomMorph(morphs []*Morph) (morph *Morph, ok bool) {
	if len(morphs) == 0 {
		for m := range c {
			morph = &m
			ok = true
			return
		}

		// cannot reach this code
		panic("chain has no entry")
	}

	var next chain
	next, ok = c[*morphs[0]]
	if !ok {
		return
	}
	return next.RandomMorph(morphs[1:])
}
