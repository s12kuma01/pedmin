// Copyright (c) 2025-2026 s12kuma01
// SPDX-License-Identifier: MPL-2.0

package client

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/Sumire-Labs/pedmin/internal/model"
)

// FxTwitterClient fetches tweet data from the fxtwitter API.
type FxTwitterClient struct {
	http *http.Client
}

func NewFxTwitterClient(timeout time.Duration) *FxTwitterClient {
	return &FxTwitterClient{
		http: &http.Client{Timeout: timeout},
	}
}

type fxResponse struct {
	Code  int     `json:"code"`
	Tweet fxTweet `json:"tweet"`
}

type fxTweet struct {
	Text            string   `json:"text"`
	Author          fxAuthor `json:"author"`
	Media           *fxMedia `json:"media"`
	Replies         int      `json:"replies"`
	Retweets        int      `json:"retweets"`
	Likes           int      `json:"likes"`
	Views           int      `json:"views"`
	CreatedAt       string   `json:"created_at"`
	CreatedTimestamp int64   `json:"created_timestamp"`
	Lang            string   `json:"lang"`
}

type fxAuthor struct {
	Name       string `json:"name"`
	ScreenName string `json:"screen_name"`
	AvatarURL  string `json:"avatar_url"`
}

type fxMedia struct {
	Photos []fxPhoto `json:"photos"`
	Videos []fxVideo `json:"videos"`
}

type fxPhoto struct {
	URL string `json:"url"`
}

type fxVideo struct {
	URL          string `json:"url"`
	ThumbnailURL string `json:"thumbnail_url"`
}

// twitterTimeFormat matches the date format returned by the fxtwitter API
// (e.g. "Sat Jun 14 01:05:09 +0000 2025").
const twitterTimeFormat = "Mon Jan 02 15:04:05 -0700 2006"

func (c *FxTwitterClient) GetTweet(ctx context.Context, screenName, tweetID string) (*model.Tweet, error) {
	endpoint := fmt.Sprintf("https://api.fxtwitter.com/%s/status/%s", screenName, tweetID)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("fxtwitter API returned %d: %s", resp.StatusCode, string(body))
	}

	var fxResp fxResponse
	if err := json.Unmarshal(body, &fxResp); err != nil {
		return nil, fmt.Errorf("invalid fxtwitter response: %w", err)
	}

	createdAt := parseTwitterTime(fxResp.Tweet.CreatedAt, fxResp.Tweet.CreatedTimestamp)

	tweet := &model.Tweet{
		Text: fxResp.Tweet.Text,
		Author: model.TweetAuthor{
			Name:       fxResp.Tweet.Author.Name,
			ScreenName: fxResp.Tweet.Author.ScreenName,
			AvatarURL:  fxResp.Tweet.Author.AvatarURL,
		},
		Replies:   fxResp.Tweet.Replies,
		Retweets:  fxResp.Tweet.Retweets,
		Likes:     fxResp.Tweet.Likes,
		Views:     fxResp.Tweet.Views,
		CreatedAt: createdAt,
		Lang:      fxResp.Tweet.Lang,
	}

	if fxResp.Tweet.Media != nil {
		for _, p := range fxResp.Tweet.Media.Photos {
			tweet.Media = append(tweet.Media, model.TweetMedia{
				Type: "photo",
				URL:  p.URL,
			})
		}
		for _, v := range fxResp.Tweet.Media.Videos {
			tweet.Media = append(tweet.Media, model.TweetMedia{
				Type:         "video",
				URL:          v.URL,
				ThumbnailURL: v.ThumbnailURL,
			})
		}
	}

	return tweet, nil
}

// parseTwitterTime tries the Unix timestamp first, then the Twitter date format string.
func parseTwitterTime(dateStr string, timestamp int64) time.Time {
	if timestamp > 0 {
		return time.Unix(timestamp, 0)
	}
	if t, err := time.Parse(twitterTimeFormat, dateStr); err == nil {
		return t
	}
	return time.Now()
}
