package main

import "log"

type stack struct {
	s [][]Morphs
}

func (s stack) len() int {
	return len(s.s)
}

func (s *stack) push(rap []Morphs) {
	s.s = append(s.s, rap)
}

func (s *stack) pop() []Morphs {
	if s.len() == 0 {
		return nil
	}
	ans := s.s[s.len()-1]
	s.s = s.s[:s.len()-1]
	return ans
}

func (s *stack) shift() {
	s.s = s.s[1:]
}

// StackParam has parameters of stack.
type StackParam struct {
	L int // stack length
}

// NewStackServer serves fresh rhymes.
func NewStackServer(out chan []Morphs, in chan []Morphs, conf *StackParam) {
	log.Println(conf.L)
	go func() {
		s := stack{}

		for {
			log.Println("stack len: ", s.len())
			switch {
			case s.len() == 0:
				s.push(<-in)
			case s.len() >= conf.L:
				out <- s.pop()
				s.shift() // refreshing
			default:
				next := s.pop()
				select {
				case rap := <-in:
					s.push(next)
					s.push(rap)
				case out <- next:
				}
			}
		}
	}()
}
