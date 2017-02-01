package main

import (
	"time"

	"github.com/high-moctane/go-mecab_slice"
	"github.com/high-moctane/go-rhymer"
)

type Rhymers struct {
	Member []rhymer.Rhymer
	Dict   *Dict
}

func NewRhymers(d *Dict, w *rhymer.MoraWeight, sim float64, length []int) Rhymers {
	member := make([]rhymer.Rhymer, 0, len(length))
	for l := range length {
		member = append(member, rhymer.New(&d.Pro, w, sim, l))
	}
	return Rhymers{Member: member, Dict: d}
}

func (r *Rhymers) Stream(ls []int, buflen int) <-chan []mecabs.Phrase {
	ans := make(chan []mecabs.Phrase, buflen)
	go func() {
		for !r.Dict.Ready {
			time.Sleep(1 * time.Second)
		}
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
	}()
	go func() {
		for {
			<-ans
			time.Sleep(10 * time.Minute)
		}
	}()
	return (<-chan []mecabs.Phrase)(ans)
}
