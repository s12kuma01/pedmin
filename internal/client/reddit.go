package client

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/s12kuma01/pedmin/internal/model"
)

// RedditClient fetches post data from the Reddit public JSON API.
type RedditClient struct {
	http *http.Client
}

func NewRedditClient(timeout time.Duration) *RedditClient {
	return &RedditClient{
		http: &http.Client{Timeout: timeout},
	}
}

type redditListing struct {
	Data struct {
		Children []struct {
			Data redditPostData `json:"data"`
		} `json:"children"`
	} `json:"data"`
}

type redditPostData struct {
	Title       string  `json:"title"`
	Selftext    string  `json:"selftext"`
	Author      string  `json:"author"`
	Subreddit   string  `json:"subreddit"`
	Score       int     `json:"score"`
	NumComments int     `json:"num_comments"`
	URL         string  `json:"url"`
	Thumbnail   string  `json:"thumbnail"`
	IsVideo     bool    `json:"is_video"`
	CreatedUTC  float64 `json:"created_utc"`
	PostHint    string  `json:"post_hint"`
	Preview     *struct {
		Images []struct {
			Source struct {
				URL string `json:"url"`
			} `json:"source"`
		} `json:"images"`
	} `json:"preview"`
	SRDetail *struct {
		IconImg       string `json:"icon_img"`
		CommunityIcon string `json:"community_icon"`
	} `json:"sr_detail"`
}

func (c *RedditClient) GetPost(ctx context.Context, subreddit, postID string) (*model.RedditPost, error) {
	endpoint := fmt.Sprintf("https://www.reddit.com/r/%s/comments/%s.json?sr_detail=1", subreddit, postID)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "pedmin-bot/1.0")

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
		return nil, fmt.Errorf("reddit API returned %d: %s", resp.StatusCode, string(body))
	}

	// Reddit returns an array: [listing (post), listing (comments)]
	var listings []redditListing
	if err := json.Unmarshal(body, &listings); err != nil {
		return nil, fmt.Errorf("invalid reddit response: %w", err)
	}

	if len(listings) == 0 || len(listings[0].Data.Children) == 0 {
		return nil, fmt.Errorf("reddit post not found")
	}

	data := listings[0].Data.Children[0].Data

	post := &model.RedditPost{
		Title:       data.Title,
		Selftext:    data.Selftext,
		Author:      data.Author,
		Subreddit:   data.Subreddit,
		Score:       data.Score,
		NumComments: data.NumComments,
		URL:         data.URL,
		Thumbnail:   data.Thumbnail,
		IsVideo:     data.IsVideo,
		CreatedUTC:  time.Unix(int64(data.CreatedUTC), 0),
		PostHint:    data.PostHint,
	}

	// Extract subreddit icon
	if data.SRDetail != nil {
		icon := data.SRDetail.CommunityIcon
		if icon == "" {
			icon = data.SRDetail.IconImg
		}
		if icon != "" {
			post.SubredditIcon = decodeHTMLEntities(icon)
		}
	}

	// Extract preview images
	if data.Preview != nil {
		for _, img := range data.Preview.Images {
			if img.Source.URL != "" {
				post.Preview = append(post.Preview, decodeHTMLEntities(img.Source.URL))
			}
		}
	}

	// Truncate selftext for display
	if len(post.Selftext) > 300 {
		post.Selftext = post.Selftext[:300] + "..."
	}

	return post, nil
}

// decodeHTMLEntities replaces &amp; with & in Reddit preview URLs.
func decodeHTMLEntities(s string) string {
	result := s
	for {
		i := 0
		found := false
		for i < len(result) {
			if i+4 < len(result) && result[i:i+5] == "&amp;" {
				result = result[:i] + "&" + result[i+5:]
				found = true
			}
			i++
		}
		if !found {
			break
		}
	}
	return result
}
