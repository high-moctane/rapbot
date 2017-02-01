package main

import (
	"time"

	"github.com/dghubble/go-twitter/twitter"
)

type ReplyHistory struct {
	Duration time.Duration
	Count    int
	List     map[int64][]time.Time
}

func (r *ReplyHistory) queue(tweet *twitter.Tweet) []time.Time {
	id := tweet.User.ID
	q, ok := r.List[id]
	if !ok {
		q = make([]time.Time, r.Count)
		r.List[id] = q
	}
	return q
}

func (r *ReplyHistory) Add(tweet *twitter.Tweet) {
	q := r.queue(tweet)
	copy(q, q[1:])
	q[len(q)-1] = time.Now()
	r.List[tweet.User.ID] = q
}

func (r *ReplyHistory) isTooFreq(tweet *twitter.Tweet) bool {
	thresh := time.Now().Add(-r.Duration)
	q := r.queue(tweet)
	for _, t := range q {
		if t.Before(thresh) {
			return false
		}
	}
	return true
}
