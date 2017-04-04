package main

import (
	"runtime"
	"sync"

	"github.com/ikawaha/kagome/tokenizer"
)

// NewMorphizer provides new Morphizer.
func NewMorphizer(in <-chan string) (<-chan Morphs, func()) {
	n := runtime.GOMAXPROCS(0)
	wg := new(sync.WaitGroup)

	out := make(chan Morphs)

	sign := make(chan struct{})
	stop := func() {
		close(sign)
		for range out {
		}
	}

	// kagome is NOT goroutine safe
	for i := 0; i < n; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			kagome := tokenizer.New()
			for {
				select {
				// interrupt
				case <-sign:
					return

				// parse
				case str, ok := <-in:
					if !ok {
						return
					}
					tokens := kagome.Tokenize(str)
					out <- NewMorphs(tokens)
				}
			}
		}()
	}

	// tell finishing
	go func() {
		wg.Wait()
		close(out)
	}()

	return out, stop
}
