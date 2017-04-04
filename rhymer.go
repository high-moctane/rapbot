package main

import (
	"runtime"
	"sync"
)

// RhymerParam is a parameter for generating rhymes.
type RhymerParam struct {
	Try    int
	Line   int
	MW     *MoraeWeight
	Thresh float64
}

// NewRhymer returns new Rhymer instance.
func NewRhymer(param *RhymerParam, in <-chan Morphs) *Rhymer {
	return &Rhymer{param: param, ms: in}
}

// Rhymer instance generate rhymes.
type Rhymer struct {
	param *RhymerParam
	ms    <-chan Morphs
}

// Generate generates rhyming morphs.
func (r *Rhymer) Generate(seed Morphs) ([]Morphs, bool) {
	ans := make([]Morphs, 1, r.param.Line)
	ans[0] = seed
loop:
	for len(ans) < r.param.Line {
		var ms Morphs
		for i := 0; i < r.param.Try; i++ {
			var ok bool
			ms, ok = <-r.ms
			if !ok {
				return nil, false
			}
			if r.param.MW.SimMorphs(ans[len(ans)-1], ms) >= r.param.Thresh {
				ans = append(ans, ms)
				continue loop
			}
		}
		if len(ans) >= 2 {
			return ans, true
		}
		return nil, false
	}
	return ans, true
}

// Server runs new Server which makes rhymes.
func (r *Rhymer) Server() <-chan []Morphs {
	wg := new(sync.WaitGroup)
	out := make(chan []Morphs)

	for i := 0; i < runtime.GOMAXPROCS(0); i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for {
				select {
				case seed, ok := <-r.ms:
					if !ok {
						return
					}
					if rhyme, ok := r.Generate(seed); ok {
						out <- rhyme
					}
				}
			}
		}()
	}

	go func() {
		wg.Wait()
		close(out)
	}()

	return out
}
