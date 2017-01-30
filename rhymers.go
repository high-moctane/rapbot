package main

import "github.com/high-moctane/go-rhymer"

type Rhymers struct {
	Member []rhymer.Rhymer
}

func NewRhymers(d *Dict, w *rhymer.MoraWeight, sim float64, length []int) Rhymers {
	member := make([]rhymer.Rhymer, 0)
	for l := range length {
		member = append(member, rhymer.New(&d.Pro, w, sim, l))
	}
	return Rhymers{Member: member}
}
