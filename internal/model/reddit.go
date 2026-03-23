package model

import "time"

// RedditPost represents a processed Reddit post for embed display.
type RedditPost struct {
	Title         string
	Selftext      string
	Author        string
	Subreddit     string
	SubredditIcon string
	Score         int
	NumComments   int
	URL           string
	Thumbnail     string
	IsVideo       bool
	CreatedUTC    time.Time
	PostHint      string // "self", "link", "image", "hosted:video", "rich:video"
	Preview       []string
}
