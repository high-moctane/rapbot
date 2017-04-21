package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"regexp"
	"runtime"
	"syscall"
	"time"

	"github.com/dghubble/go-twitter/twitter"
)

var configFlag = flag.String("c", "./config.json", "path to config.json")

func main() {
	os.Exit(run())
}

func run() int {
	// pprof
	runtime.SetBlockProfileRate(1)
	go func() {
		log.Println(http.ListenAndServe("localhost:6060", nil))
	}()

	log.Print("launch!")
	flag.Parse()

	// load config
	config, err := newConfig(*configFlag)
	if err != nil {
		log.Print("invalid config file:", err)
		return 1
	}

	// make twitter client
	client := newClient(config)
	sendLog := func(s string) (*twitter.DirectMessage, *http.Response, error) {
		log.Printf("send DM to author: %q", s)
		dm, resp, err := client.DirectMessages.New(&twitter.DirectMessageNewParams{
			ScreenName: config.TwitterParam.LogToScreenName,
			Text:       s,
		})
		if err != nil {
			log.Print("logging DM send error:", err)
		}
		return dm, resp, err
	}

	// make channels
	sampleStr := make(chan string, 100)
	randMorphs := make(chan Morphs, 100)
	rhymes := make(chan []Morphs, 100)

	// connect twitter sample
	stream, err := client.Streams.Sample(&twitter.StreamSampleParams{
		StallWarnings: twitter.Bool(true),
		Language:      []string{"ja"},
	})
	if err != nil {
		log.Print("cannot connect twitter stream sample", err)
		return 1
	}
	defer stream.Stop()

	demuxStrm := twitter.NewSwitchDemux()
	demuxStrm.Tweet = func(tweet *twitter.Tweet) {
		// filter
		for _, level := range config.TwitterParam.Filter {
			if tweet.FilterLevel == level {
				return
			}
		}
		if tweet.RetweetedStatus != nil {
			return
		}
		if tweet.User.FavouritesCount < 10 {
			return
		}

		// remove some invalid sentence
		url := regexp.MustCompile(`(^|\p{Zs})(http|https|ttp|ttps)://.*?($|\p{Zs})`)
		mention := regexp.MustCompile(`(^|\p{Zs}|\.)@.*?($|\p{Zs})`)
		hashtag := regexp.MustCompile(`(^|\p{Zs})(♯|＃|#).*?($|\p{Zs})`)
		toWhite := regexp.MustCompile(
			`(「|」|\[|\]|（|）|\(|\)|。|、|\,|\.|，|．|【|】|『|』|〈|〉|［|］|《|》|？|！|\?|\!|…|〜)`,
		)

		text := tweet.Text

		text = url.ReplaceAllString(text, " ")
		text = mention.ReplaceAllString(text, " ")
		text = hashtag.ReplaceAllString(text, " ")
		text = toWhite.ReplaceAllString(text, " ")

		if text == "" {
			return
		}
		for _, t := range regexp.MustCompile(`\p{Zs}.+`).Split(text, -1) {
			sampleStr <- t
		}
	}
	demuxStrm.StreamDisconnect = func(dscn *twitter.StreamDisconnect) {
		log.Printf("sample stream disconnected: code: %v, stream_name: %q, reason: %q",
			dscn.Code, dscn.Reason, dscn.StreamName)
		sendLog("sample stream disconnected")
	}
	demuxStrm.Warning = func(warning *twitter.StallWarning) {
		log.Printf("sample stream stall warning: code: %q, message: %q, percent_full: %q",
			warning.Code, warning.Message, warning.PercentFull)
	}
	demuxStrm.FriendsList = func(_ *twitter.FriendsList) {
		log.Print("stream sample connected")
	}
	go demuxStrm.HandleChan(stream.Messages)

	// tokenize tweet text
	parsedMorphs, _ := NewMorphizer(sampleStr)

	// learn Markov and generate random Morphs
	for _, param := range config.MarkovParams {
		out := MarkovServer(param, parsedMorphs)
		go func() {
			for ms := range out {
				if len(ms) == 0 {
					continue
				}
				randMorphs <- ms
			}
		}()
	}

	// generate rhymes
	for _, param := range config.RhymerParams {
		r := NewRhymer(param, randMorphs)
		out := r.Server()
		go func() {
			for rhyme := range out {
				rhymes <- rhyme
			}
		}()
	}

	// buffer rhymes
	raps := make(chan []Morphs)
	NewStackServer(raps, rhymes, config.StackParam)

	// print
	go func() {
		for rhyme := range raps {
			for _, ms := range rhyme {
				p, _ := ms.Surface()
				fmt.Println(p)
			}
			fmt.Println("")
			time.Sleep(10 * time.Hour)
		}
	}()

	// ctrl+c
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT)
	<-sig
	log.Print("interrupted")

	return 0
}
