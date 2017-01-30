package main

import (
	"github.com/high-moctane/go-markov_chain_Japanese"
	"github.com/high-moctane/go-mecab_slice"
)

const Order = 2

type Dict struct {
	Trainee, Pro markov.Markov
	mecabs       *mecabs.MeCabS
	Count        int
	Ready        bool
}

func NewDict(m *mecabs.MeCabS) Dict {
	dmT := markov.NewDataMap(Order)
	dmP := markov.NewDataMap(Order)
	return Dict{
		Trainee: markov.New(m, &dmT),
		Pro:     markov.New(m, &dmP),
		mecabs:  m,
		Ready:   false,
	}
}

func (d *Dict) Shift() {
	dm := markov.NewDataMap(Order)
	d.Pro.Data = d.Trainee.Data
	d.Trainee = markov.New(d.mecabs, &dm)
	d.Ready = true
}

func (d *Dict) TryShift() {
	d.Count++
	if d.Count > LearnMax {
		d.Shift()
		d.Count = 0
	}
}
