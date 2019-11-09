package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/dghubble/go-twitter/twitter"
)

// ChTweets is stream of tweets
var ChTweets = make(chan string, 10)

// ChTweetSentence is stream of sentences
var ChTweetSentence = make(chan Sentence, 10)

// ChRandomSentence is stream of random sentences.
var ChRandomSentence = make(chan Sentence, 10)

// ChLyric is a stream of lyrics
var ChLyric = make(chan Lyric, 5)

// markov is random sentence generator.
var markov *Markov

// rapper is main rapper.
var rapper *Rapper

// lyricStorage is global lyricStorage
var lyricStorage = NewLyricStorage(10000)

// TwiClient is Twitter client
var TwiClient *twitter.Client

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	// read .env
	if err := LoadEnv(); err != nil {
		return err
	}

	// initialize global variables
	markovParams, err := DefaultMarkovParams()
	if err != nil {
		return fmt.Errorf("cannot create markov: %w", err)
	}
	markov = NewMarkov(markovParams)
	rapper, err = DefaultRapper()
	if err != nil {
		return err
	}
	TwiClient = NewTwitterClient()

	// parse tweets
	go JapaneseParseServer(ChTweetSentence, ChTweets)

	// build markov chains
	go markov.AddServer(ChTweetSentence)

	// generate random sentence
	if err := markov.LaunchRandomSentenceServer(ChRandomSentence); err != nil {
		return err
	}

	// generate lyrics
	if err := rapper.LaunchRapServer(ChLyric, ChRandomSentence); err != nil {
		return err
	}

	// store lyrics
	go lyricStorage.PushServer(ChLyric)

	// Twitter stream
	streamSample, err := TwiClient.Streams.Sample(&twitter.StreamSampleParams{
		StallWarnings: twitter.Bool(true),
	})
	if err != nil {
		return fmt.Errorf("twitter sample error: %w", err)
	}
	defer streamSample.Stop()
	demuxSample := twitter.NewSwitchDemux()
	demuxSample.Tweet = ExtractTweet
	go demuxSample.HandleChan(streamSample.Messages)
	// TODO

	// regular tweet
	LaunchRegularTweetServer()

	// serve twitter reply
	streamReply, err := TwiClient.Streams.Filter(&twitter.StreamFilterParams{
		Track:         []string{os.Getenv("TWITTER_SCREENNAME")},
		StallWarnings: twitter.Bool(true),
	})
	if err != nil {
		return fmt.Errorf("twitter filter reply error: %w", err)
	}
	defer streamReply.Stop()
	demuxReply := twitter.NewSwitchDemux()
	demuxReply.Tweet = ServeReply
	go demuxReply.HandleChan(streamReply.Messages)

	// signal handling
	chSig := make(chan os.Signal, 1)
	signal.Notify(chSig, syscall.SIGINT, syscall.SIGTERM)
	log.Println(<-chSig)
	return nil
}
