package main

import (
	"encoding/json"
	"os"

	"github.com/pkg/errors"
)

type config struct {
	MarkovParams []*MarkovParam `json:"markov_params"`
	RhymerParams []*RhymerParam `json:"rhymer_params"`
	TwitterParam *TwitterParam  `json:"twitter_param"`
}

func newConfig(p string) (*config, error) {
	f, err := os.Open(p)
	if err != nil {
		return nil, errors.Wrapf(err, "cannot open %s", p)
	}
	defer f.Close()

	ans := new(config)
	dec := json.NewDecoder(f)
	dec.Decode(ans)
	return ans, nil
}
