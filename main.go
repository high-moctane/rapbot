package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"os"
	"os/signal"
	"path/filepath"
	"regexp"
	"syscall"
	"time"
	"unicode/utf8"

	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
	"github.com/high-moctane/go-mecab_slice"
	"github.com/high-moctane/go-rhymer"
)

// Config
const LearnMax = 10000

var TokenPath = filepath.Join(os.Getenv("HOME"), ".config", "go-rapbot", "token.json")
var MeCabParam = map[string]string{}

var MoraWeight = rhymer.NewMoraWeight([]rhymer.MoraWeightCell{{5.0, 10.0}, {5.0, 15.0}, {10.0, 20.0}, {20.0, 50.0}})

const Similarity = 0.8

const TweetDuration = 1*time.Hour + 30*time.Minute

const ReplyDuration = 1 * time.Minute

const ReplyCount = 5

// Global variable

var MyClient *twitter.Client
var MyUser *twitter.User

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
	toWhite := regexp.MustCompile(
		`(「|」|\[|\]|（|）|\(|\)|。|、|\,|\.|，|．|【|】|『|』|〈|〉|［|］|《|》|？|！|\?|\!|…|〜)`,
	)

	str := tweet.Text

	str = mention.ReplaceAllString(str, " ")
	str = url.ReplaceAllString(str, " ")
	str = hashtag.ReplaceAllString(str, " ")
	str = toWhite.ReplaceAllString(str, " ")

	return str
}

func Learn(d *Dict) func(*twitter.Tweet) {
	white := regexp.MustCompile(`^.\s$`)
	return func(tweet *twitter.Tweet) {
		if !isValidTweet(tweet) {
			return
		}
		txt := BodyText(tweet)
		if white.Match([]byte(txt)) {
			return
		}
		d.Trainee.Add(txt)
		d.TryShift()
	}
}

func PeriodicTweet(rhymes <-chan []mecabs.Phrase) {
	for rhyme := range rhymes {
		str := RapToStr(rhyme)
		MyClient.Statuses.Update(str, nil)

		time.Sleep(TweetDuration)
	}
}

func Reply(repHist *ReplyHistory, rhymes <-chan []mecabs.Phrase) func(*twitter.Tweet) {
	return func(tweet *twitter.Tweet) {
		if !ReplyForMe(tweet) {
			return
		}
		if repHist.isTooFreq(tweet) {
			return
		}
		var rap []mecabs.Phrase
		var mes string
		for {
			rap = <-rhymes
			mes = MakeReplyStr(tweet, RapToStr(rap))
			if utf8.RuneCountInString(mes) <= 140 {
				break
			}
		}
		MyClient.Statuses.Update(
			mes,
			&twitter.StatusUpdateParams{
				InReplyToStatusID: tweet.ID,
			},
		)
		repHist.Add(tweet)
	}
}

func ReplyForMe(tweet *twitter.Tweet) bool {
	if tweet.User.ID == MyUser.ID {
		return false
	}
	if tweet.InReplyToUserID != MyUser.ID {
		return false
	}
	return true
}

func RapToStr(rap []mecabs.Phrase) string {
	ans := ""
	for i, r := range rap {
		ans += r.OriginalForm()
		if i < len(rap)-1 {
			ans += "\n"
		}
	}
	return ans
}

func MakeReplyStr(tweet *twitter.Tweet, str string) string {
	return fmt.Sprint("@", tweet.User.ScreenName, " ", str)
}

func StreamSample(dict *Dict) {
	demux := twitter.NewSwitchDemux()
	demux.Tweet = func(tweet *twitter.Tweet) {
		Learn(dict)(tweet)
	}
	stream, err := MyClient.Streams.Sample(
		&twitter.StreamSampleParams{StallWarnings: twitter.Bool(true)},
	)
	if err != nil {
		log.Fatal(err)
	}
	demux.HandleChan(stream.Messages)
}

func StreamUser(rhyme <-chan []mecabs.Phrase, dict *Dict) {
	// TODO: ここを後で外部の変数にしておく
	repHist := ReplyHistory{
		Duration: ReplyDuration,
		Count:    ReplyCount,
		List:     map[int64][]time.Time{},
	}
	demux := twitter.NewSwitchDemux()
	demux.Tweet = func(tweet *twitter.Tweet) {
		if !dict.Ready {
			return
		}
		Reply(&repHist, rhyme)(tweet)
	}
	stream, err := MyClient.Streams.User(
		&twitter.StreamUserParams{StallWarnings: twitter.Bool(true)},
	)
	if err != nil {
		log.Fatal(err)
	}
	demux.HandleChan(stream.Messages)
}

func main() {
	rand.Seed(time.Now().UnixNano())

	t, err := LoadToken(TokenPath)
	if err != nil {
		log.Fatal(err)
	}

	MyClient = NewClient(t)
	MyUser, _, err = MyClient.Accounts.VerifyCredentials(
		&twitter.AccountVerifyParams{
			SkipStatus:   twitter.Bool(true),
			IncludeEmail: twitter.Bool(true),
		},
	)
	if err != nil {
		log.Fatal(err)
	}

	mecabs, err := mecabs.New(MeCabParam)
	if err != nil {
		log.Fatal(err)
	}
	defer mecabs.Destroy()

	dict := NewDict(&mecabs)

	rhymers := NewRhymers(&dict, &MoraWeight, Similarity, []int{5, 6, 6, 6, 6, 7, 7})
	rhymes := rhymers.Stream([]int{2, 3, 3, 4, 4}, 100)

	go StreamSample(&dict)
	go StreamUser(rhymes, &dict)
	go PeriodicTweet(rhymes)

	log.Println("start(｀･ω･´)")

	// Wait for SIGINT and SIGTERM (HIT CTRL-C)
	ch := make(chan os.Signal)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	log.Println(<-ch)
	log.Println("see you (｀･ω･´)")
}
