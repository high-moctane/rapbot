package main

import (
	"encoding/json"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"regexp"
	"time"

	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
	"github.com/high-moctane/go-mecab_slice"
	"github.com/high-moctane/go-rhymer"
)

const LearnMax = 10000

var TokenPath = filepath.Join(os.Getenv("HOME"), ".config", "go-rapbot", "token.json")
var MeCabParam = map[string]string{}

var MoraWeight = rhymer.NewMoraWeight([]rhymer.MoraWeightCell{{5.0, 10.0}, {5.0, 15.0}, {10.0, 20.0}, {20.0, 50.0}})

const Similarity = 0.8

var MoraLen = 6

const TweetDuration = 5 * time.Second

type Token struct {
	ConsumerKey       string
	ConsumerSecret    string
	AccessToken       string
	AccessTokenSecret string
}

func LoadToken(path string) (Token, error) {
	f, err := os.Open(path)
	if err != nil {
		return Token{}, err
	}
	defer f.Close()

	var token Token
	err = json.NewDecoder(f).Decode(&token)
	return token, nil
}

func NewClient(t Token) *twitter.Client {
	config := oauth1.NewConfig(t.ConsumerKey, t.ConsumerSecret)
	token := oauth1.NewToken(t.AccessToken, t.AccessTokenSecret)
	httpClient := config.Client(oauth1.NoContext, token)
	return twitter.NewClient(httpClient)
}

func isValidTweet(tweet *twitter.Tweet) bool {
	if tweet.Lang != "ja" {
		return false
	}
	if tweet.RetweetCount > 0 {
		return false
	}
	return true
}

func BodyText(tweet *twitter.Tweet) string {
	mention := regexp.MustCompile(`(^|\s|\.)@.*?($|\s)`)
	url := regexp.MustCompile(`(^|\s)(http|https)://.*?($|\s)`)
	hashtag := regexp.MustCompile(`(^|\s)#.*?($|\s)`)

	str := tweet.Text

	str = mention.ReplaceAllString(str, " ")
	str = url.ReplaceAllString(str, " ")
	str = hashtag.ReplaceAllString(str, " ")

	return str
}

func Learn(s *twitter.Stream, d *Dict) {
	demux := twitter.NewSwitchDemux()
	demux.Tweet = func(tweet *twitter.Tweet) {
		if !isValidTweet(tweet) {
			return
		}
		d.Trainee.Add(BodyText(tweet))
		d.TryShift()
	}
	demux.HandleChan(s.Messages)
}

func ContinuousLearn(s *twitter.Stream, d *Dict) {
	go func() {
		for {
			Learn(s, d)
			time.Sleep(1 * time.Minute)
		}
	}()
}

func PeriodicTweet(rhymes <-chan []mecabs.Phrase) {
	go func() {
		for rhyme := range rhymes {
			var str string
			for _, p := range rhyme {
				str += p.OriginalForm() + "\n"
			}
			log.Println(str)

			time.Sleep(TweetDuration)
		}
	}()
}

func main() {
	rand.Seed(time.Now().UnixNano())

	t, err := LoadToken(TokenPath)
	if err != nil {
		log.Fatal(err)
	}

	client := NewClient(t)

	mecabs, err := mecabs.New(MeCabParam)
	if err != nil {
		log.Fatal(err)
	}
	defer mecabs.Destroy()

	dict := NewDict(&mecabs)

	stream, err := client.Streams.Sample(
		&twitter.StreamSampleParams{StallWarnings: twitter.Bool(true)},
	)
	if err != nil {
		log.Fatal(err)
	}
	ContinuousLearn(stream, &dict)

	for !dict.Ready {
		time.Sleep(1 * time.Second)
	}

	rhymers := NewRhymers(&dict, &MoraWeight, Similarity, []int{5, 6, 6, 6, 6, 7, 7})
	rhymes := rhymers.Stream([]int{2, 3, 3, 4, 4}, 100)

	PeriodicTweet(rhymes)

	time.Sleep(24 * time.Hour)
}
