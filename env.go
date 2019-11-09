package main

import (
	"errors"
	"log"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

// LoadEnv read .env and verify it.
func LoadEnv() error {
	if err := godotenv.Load(); err != nil {
		log.Fatal("cannot read .env")
	}
	if err := verifyEnv(os.Environ()); err != nil {
		log.Fatal("invalid .env: ", err)
	}
	return nil
}

func verifyEnv(envs []string) error {
	entryies := []string{
		"CONSUMER_KEY",
		"CONSUMER_SECRET",
		"ACCESS_TOKEN",
		"ACCESS_TOKEN_SECRET",
		"TWITTER_SCREENNAME",
		"REGULAR_TWEET_MINUTES",
		"NGRAM",
		"CHAIN_NUM",
		"CHAIN_MORPHS_NUM",
		"RANDOM_MORPH_LEN",
		"TRY_NUM",
		"THRESH",
		"CONSONANT_WEIGHTS",
		"VOWEL_WEIGHTS",
		"LYRIC_LINE_NUM",
	}

	// return error when ent is not found in envs.
findLoop:
	for _, ent := range entryies {
		for _, env := range envs {
			if strings.HasPrefix(env, ent+"=") {
				continue findLoop
			}
		}
		return errors.New("not found " + ent)
	}

	return nil
}
