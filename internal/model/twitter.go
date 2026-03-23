package model

import "time"

// Tweet represents a processed tweet for embed display.
type Tweet struct {
	Text      string
	Author    TweetAuthor
	Media     []TweetMedia
	Replies   int
	Retweets  int
	Likes     int
	Views     int
	CreatedAt time.Time
	Lang      string
}

// TweetAuthor holds tweet author information.
type TweetAuthor struct {
	Name       string
	ScreenName string
	AvatarURL  string
}

// TweetMedia holds tweet media information.
type TweetMedia struct {
	Type         string // "photo" or "video"
	URL          string
	ThumbnailURL string
}
