package models

import (
	"time"
)

type Tweet struct {
	Username string `json:"username" validate:"required"`
	Message  string `json:"message" validate:"required,max=150,HasSwearWords"`
	Created  time.Time
}

var tweets []Tweet

func init() {
	// tweets = append(tweets, Tweet{Username: "amul", Message: "Hello, first tweet"})
}

func AddTweet(tweet Tweet) {
	tweets = append(tweets, tweet)
}

func FetchTweets() []Tweet {
	return tweets
}
