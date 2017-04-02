package main

import (
	"math/rand"
	"runtime"
	"sync"
)

type chainMap map[Morph]*chain

func newChain() *chain {
	return &chain{c: make(chainMap)}
}

type chain struct {
	c     chainMap
	count int
}

func (c *chain) add(ms Morphs) {
	next, ok := c.c[*ms[0]]
	if !ok {
		next = newChain()
		c.c[*ms[0]] = next
	}
	if len(ms) == 2 {
		if _, ok := next.c[*ms[1]]; !ok {
			next.c[*ms[1]] = nil
		}
		return
	}
	next.add(ms[1:])
}

func (c *chain) inc() {
	c.count++
}

func (c *chain) findRand(ms Morphs) (*Morph, bool) {
	next, ok := c.c[*ms[0]]
	if !ok {
		return nil, false
	}
	if len(ms) == 1 {
		n := rand.Intn(len(c.c[*ms[0]].c))
		i := 0
		for k := range c.c[*ms[0]].c {
			if i >= n {
				return &k, true
			}
			i++
		}
		return nil, false
	}
	return next.findRand(ms[1:])
}

// MarkovParam defines Markov chain model's property.
type MarkovParam struct {
	n   int // ngram (n >= 2)
	lcs int // max length of chains
	lc  int // max length of each chain
	lms int // max length of generated Morphs
	try int // max trying count
}

// newMarkov makes new Markov instance
func newMarkov(param *MarkovParam) *markov {
	return &markov{
		param:    param,
		cs:       []*chain{newChain()},
		learning: newChain(),
		ready:    make(chan struct{}),
	}
}

// markov is a model of markov chain.
type markov struct {
	param    *MarkovParam
	cs       []*chain
	learning *chain
	once     sync.Once
	ready    chan struct{}
}

func (m *markov) shiftable() bool {
	return m.learning.count >= m.param.lc
}

func (m *markov) shift() {
	if len(m.cs) >= m.param.lcs {
		m.cs[0] = nil
		m.cs = m.cs[1:]
	}
	m.cs = append(m.cs, m.learning)
	m.learning = newChain()
	m.once.Do(func() { close(m.ready) })
}

func (m *markov) add(ms Morphs) {
	// make vector
	vec := make(Morphs, m.param.n)
	for i := range vec {
		vec[i] = &MorphBOS
	}

	// add vectors to markov
	for _, morph := range append(ms, &MorphEOS) {
		vec = append(vec[1:], morph)
		m.learning.add(vec)
	}
	m.learning.inc()

	// shift
	if m.shiftable() {
		m.shift()
	}
}

func (m *markov) generate() (Morphs, bool) {
gen:
	for i := 0; i < m.param.try; i++ {
		ms := make(Morphs, 0, m.param.lms)

		// make vector
		vec := make(Morphs, m.param.n)
		for i := range vec {
			vec[i] = &MorphBOS
		}

		// make phrase
		for i := 0; i < m.param.lms; i++ {
			cand := make(Morphs, 0, len(m.cs))

			// find
			for _, c := range m.cs {
				if morph, ok := c.findRand(vec[1:]); ok {
					cand = append(cand, morph)
				}
			}
			if len(cand) == 0 {
				continue
			}
			morph := cand[rand.Intn((len(cand)))]
			if *morph == MorphEOS {
				if len(ms) == 0 {
					continue gen
				}
				return ms, true
			}
			ms = append(ms, morph)
			vec = append(vec[1:], ms[len(ms)-1])
		}

		if len(ms) > 0 {
			return ms, true
		}
	}
	return nil, false
}

// MarkovServer starts server which can learn and generate phrases.
func MarkovServer(param *MarkovParam) (<-chan Morphs, chan<- Morphs) {
	out := make(chan Morphs)
	in := make(chan Morphs)

	go func() {
		wg := new(sync.WaitGroup)
		sema := make(chan struct{}, runtime.GOMAXPROCS(0))
		m := newMarkov(param)

		// server
		for {
			select {
			// add input
			case ms, ok := <-in:
				wg.Wait()
				if !ok {
					close(out)
					return
				}
				if len(ms) == 0 {
					continue
				}
				m.add(ms)

			// generate random Morphs
			default:
				sema <- struct{}{}
				wg.Add(1)
				go func() {
					defer wg.Done()
					if ms, ok := m.generate(); ok {
						out <- ms
					}
				}()
				<-sema
			}
		}
	}()
	return out, in
}