package main

import (
	"github.com/high-moctane/go-mecab_slice"
	"github.com/high-moctane/go-rhymer"
)

type Rhymers struct {
	Member []rhymer.Rhymer
}

func NewRhymers(d *Dict, w *rhymer.MoraWeight, sim float64, length []int) Rhymers {
	member := make([]rhymer.Rhymer, 0, len(length))
	for l := range length {
		member = append(member, rhymer.New(&d.Pro, w, sim, l))
	}
	return Rhymers{Member: member}
}

func (r *Rhymers) Stream(ls []int, buflen int) <-chan []mecabs.Phrase {
	ans := make(chan []mecabs.Phrase, buflen)
	for _, m := range r.Member {
		for _, l := range ls {
			stream, _ := m.Stream(l)
			go func() {
				for p := range stream {
					ans <- p
				}
			}()
		}
	}
	return (<-chan []mecabs.Phrase)(ans)
}
