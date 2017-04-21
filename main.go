package main

import (
	"flag"
	"log"
	"math/rand"
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
	myTwitterStatus, _, err := client.Accounts.VerifyCredentials(nil)
	if err != nil {
		log.Println("twitter verify error:", err)
		return 1
	}

	// make channels
	sampleStr := make(chan string)
	randMorphs := make(chan Morphs)
	rhymes := make(chan []Morphs)

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
		for _, t := range regexp.MustCompile(`(\p{Zs}|\n).+`).Split(text, -1) {
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

	// Time Line
	tl, err := client.Streams.User(&twitter.StreamUserParams{
		StallWarnings: twitter.Bool(true),
	})
	if err != nil {
		log.Print("cannot connect twitter user stream", err)
		return 1
	}
	defer tl.Stop()

	repFreq := map[int64][]time.Time{}
	demuxTL := twitter.NewSwitchDemux()
	demuxTL.Tweet = func(tweet *twitter.Tweet) {
		myScreenName := "@" + myTwitterStatus.ScreenName
		if !regexp.MustCompile(`(\A|\A\.)` + myScreenName).Match([]byte(tweet.Text)) {
			return
		}

		if freqQueue, ok := repFreq[tweet.User.ID]; ok {
			now := time.Now()
			for i, t := range repFreq[tweet.User.ID] {
				if now.After(t.Add(config.TwitterParam.FreqSeconds * time.Second)) {
					freqQueue = append(freqQueue[:i], freqQueue[i+1:]...)
				}
			}
			repFreq[tweet.User.ID] = freqQueue
			if len(freqQueue) >= config.TwitterParam.Freq {
				log.Println("received too freq reply from:", tweet.User.ScreenName)
				return
			}
		}

		message := "@" + tweet.User.ScreenName + " "
		select {
		case rap := <-raps:
			for i, r := range rap {
				if i != 0 {
					message += "\n"
				}
				s, _ := r.Surface()
				message += s
			}
		default:
			message += "ネタ切れ御免。。。"
		}
		_, _, err := client.Statuses.Update(message, &twitter.StatusUpdateParams{
			InReplyToStatusID: tweet.User.ID,
		})
		if err != nil {
			log.Printf("tweet error: %v, message: %q", err, message)
			return
		}
		log.Printf("tweet: %q", message)
		repFreq[tweet.User.ID] = append(repFreq[tweet.User.ID], time.Now())
	}
	demuxTL.StreamDisconnect = func(dscn *twitter.StreamDisconnect) {
		log.Printf("user stream disconnected: code: %v, stream_name: %q, reason: %q",
			dscn.Code, dscn.Reason, dscn.StreamName)
		sendLog("user stream disconnected")
	}
	demuxTL.Warning = func(warning *twitter.StallWarning) {
		log.Printf("user stream stall warning: code: %q, message: %q, percent_full: %q",
			warning.Code, warning.Message, warning.PercentFull)
	}
	demuxTL.FriendsList = func(_ *twitter.FriendsList) {
		log.Print("user stream connected")
	}
	go demuxTL.HandleChan(tl.Messages)

	// routine tweet
	go func() {
		for rap := range raps {
			var message string
			for i, r := range rap {
				if i != 0 {
					message += "\n"
				}
				s, _ := r.Surface()
				message += s
			}
			_, _, err := client.Statuses.Update(message, nil)
			if err != nil {
				log.Printf("routine tweet error: %v, message: %q", err, message)
				continue
			}
			log.Printf("routine tweet: %q", message)
			time.Sleep(config.TwitterParam.RoutineMinutes*time.Minute +
				time.Duration(float64(10*time.Minute)*rand.Float64()))
			//todo
		}
	}()

	// ctrl+c
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT)
	<-sig
	log.Print("interrupted")

	return 0
}
