package main

import (
	"fmt"
	"html"
	"os"
	"strconv"
	"time"

	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
	"github.com/ikawaha/kagome/tokenizer"
)

// NewTwitterClient returns new twitter client.
func NewTwitterClient() *twitter.Client {
	config := oauth1.NewConfig(os.Getenv("CONSUMER_KEY"), os.Getenv("CONSUMER_SECRET"))
	token := oauth1.NewToken(os.Getenv("ACCESS_TOKEN"), os.Getenv("ACCESS_TOKEN_SECRET"))
	httpClient := config.Client(oauth1.NoContext, token)
	return twitter.NewClient(httpClient)
}

// ExtractTweet extract valid text and send if to chTweets.
func ExtractTweet(tweet *twitter.Tweet) {
	if !isLearnableTweet(tweet) {
		return
	}

	text := html.UnescapeString(tweet.Text)

	// avoid blocking
	select {
	case ChTweets <- text:
	default:
	}
}

// isLearnableTweet returns the tweet is valid.
// Filter spams by this function.
func isLearnableTweet(tweet *twitter.Tweet) bool {
	return true && // dummy for easy comment out
		tweet.Lang == "ja" && // Japanese lang tweet
		// tweet.RetweetedStatus == nil && // not retweet
		// tweet.QuotedStatus == nil && // not retweet
		tweet.InReplyToScreenName == "" && // not reply
		len(tweet.Entities.Hashtags) == 0 && // no hashtags
		len(tweet.Entities.Media) == 0 && // no media
		len(tweet.Entities.Urls) == 0 && // no urls
		len(tweet.Entities.UserMentions) == 0 && // no mentions
		tweet.User.FriendsCount > 10 && // has some friends
		tweet.User.FollowersCount > 10 && // has some followers
		true // dummy for easy comment out
}

// ServeReply serves a reply.
func ServeReply(tweet *twitter.Tweet) {
	t := tokenizer.New()
	sentence := analyzeText(&t, tweet.Text)
	lyric := lyricStorage.ContinueLyric(rapper, sentence)

	header := "@" + tweet.User.ScreenName

	var body string
	if lyric == nil {
		body = "準備中です(｀･ω･´)"
	} else {
		body = lyric.String()
	}

	TwiClient.Statuses.Update(
		header+"\n"+body,
		&twitter.StatusUpdateParams{
			InReplyToStatusID: tweet.ID,
		},
	)
}

// LaunchRegularTweetServer post tweet.
func LaunchRegularTweetServer() error {
	duration, err := strconv.ParseInt(os.Getenv("REGULAR_TWEET_MINUTES"), 10, 64)
	if err != nil {
		return fmt.Errorf("invalid REGULAR_TWEET_MINUTES: %w", err)
	}

	go func() {
		ticker := time.NewTicker(time.Duration(duration) * time.Minute)
		<-markov.Ready
		for {
			<-ticker.C
			lyric := lyricStorage.Pop()
			if lyric == nil {
				continue
			}

			TwiClient.Statuses.Update(lyric.String(), nil)
		}
	}()
	return nil
}
