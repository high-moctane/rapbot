package main

import (
	"runtime"
	"sync"

	"github.com/ikawaha/kagome/tokenizer"
)

// NewMorphizer provides new Morphizer.
func NewMorphizer() (<-chan Morphs, chan<- string, func()) {
	n := runtime.GOMAXPROCS(0)
	wg := new(sync.WaitGroup)

	out := make(chan Morphs)
	in := make(chan string)

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
			kagome := tokenizer.New()
			for {
				select {
				// interrupt
				case <-sign:
					wg.Done()
					return

				// parse
				case str, ok := <-in:
					if !ok {
						wg.Done()
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

	return out, in, stop
}
