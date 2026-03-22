package embedfix

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

// InstagramClient fetches post data from the Meta oEmbed API.
type InstagramClient struct {
	http        *http.Client
	accessToken string
}

func NewInstagramClient(accessToken string, timeout time.Duration) *InstagramClient {
	return &InstagramClient{
		http:        &http.Client{Timeout: timeout},
		accessToken: accessToken,
	}
}

// IsAvailable returns true if the Instagram client has valid credentials.
func (c *InstagramClient) IsAvailable() bool {
	return c.accessToken != ""
}

type InstagramPost struct {
	AuthorName   string
	Title        string
	ThumbnailURL string
	CreatedAt    time.Time
}

type oembedResponse struct {
	AuthorName   string `json:"author_name"`
	AuthorURL    string `json:"author_url"`
	Title        string `json:"title"`
	ThumbnailURL string `json:"thumbnail_url"`
	HTML         string `json:"html"`
}

func (c *InstagramClient) GetPost(ctx context.Context, shortcode string) (*InstagramPost, error) {
	postURL := fmt.Sprintf("https://www.instagram.com/p/%s/", shortcode)
	endpoint := fmt.Sprintf("https://graph.facebook.com/v22.0/instagram_oembed?url=%s&access_token=%s",
		url.QueryEscape(postURL), url.QueryEscape(c.accessToken))

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
		return nil, fmt.Errorf("instagram oEmbed API returned %d: %s", resp.StatusCode, string(body))
	}

	var oembed oembedResponse
	if err := json.Unmarshal(body, &oembed); err != nil {
		return nil, fmt.Errorf("invalid instagram oEmbed response: %w", err)
	}

	return &InstagramPost{
		AuthorName:   oembed.AuthorName,
		Title:        oembed.Title,
		ThumbnailURL: oembed.ThumbnailURL,
		CreatedAt:    time.Now(), // oEmbed does not provide created_at
	}, nil
}
