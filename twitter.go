package main

import (
	"time"

	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
)

// TwitterParam is keys of client.
type TwitterParam struct {
	ConsumerKey       string   `json:"consumer_key"`
	ConsumerSecret    string   `json:"consumer_secret"`
	AccessToken       string   `json:"access_token"`
	AccessTokenSecret string   `json:"access_token_secret"`
	LogToScreenName   string   `json:"log_to_screen_name"`
	Filter            []string `json:"filter"`
	Freq              int
	FreqSeconds       time.Duration `json:"freq_seconds"`
	RoutineMinutes    time.Duration `json:"routine_minutes"`
}

func newClient(c *config) *twitter.Client {
	config := oauth1.NewConfig(c.TwitterParam.ConsumerKey, c.TwitterParam.ConsumerSecret)
	token := oauth1.NewToken(c.TwitterParam.AccessToken, c.TwitterParam.AccessTokenSecret)
	httpClient := config.Client(oauth1.NoContext, token)
	return twitter.NewClient(httpClient)
}
