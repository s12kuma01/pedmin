package embedfix

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
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

type TweetAuthor struct {
	Name       string
	ScreenName string
	AvatarURL  string
}

type TweetMedia struct {
	Type         string // "photo" or "video"
	URL          string
	ThumbnailURL string
}

type fxResponse struct {
	Code  int     `json:"code"`
	Tweet fxTweet `json:"tweet"`
}

type fxTweet struct {
	Text      string   `json:"text"`
	Author    fxAuthor `json:"author"`
	Media     *fxMedia `json:"media"`
	Replies   int      `json:"replies"`
	Retweets  int      `json:"retweets"`
	Likes     int      `json:"likes"`
	Views     int      `json:"views"`
	CreatedAt string   `json:"created_at"`
	Lang      string   `json:"lang"`
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

func (c *FxTwitterClient) GetTweet(ctx context.Context, screenName, tweetID string) (*Tweet, error) {
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

	createdAt, _ := time.Parse(time.RFC1123, fxResp.Tweet.CreatedAt)

	tweet := &Tweet{
		Text: fxResp.Tweet.Text,
		Author: TweetAuthor{
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
			tweet.Media = append(tweet.Media, TweetMedia{
				Type: "photo",
				URL:  p.URL,
			})
		}
		for _, v := range fxResp.Tweet.Media.Videos {
			tweet.Media = append(tweet.Media, TweetMedia{
				Type:         "video",
				URL:          v.URL,
				ThumbnailURL: v.ThumbnailURL,
			})
		}
	}

	return tweet, nil
}

// TranslateClient calls the DeepL API Free.
type TranslateClient struct {
	http   *http.Client
	apiKey string
}

func NewTranslateClient(apiKey string, timeout time.Duration) *TranslateClient {
	return &TranslateClient{
		http:   &http.Client{Timeout: timeout},
		apiKey: apiKey,
	}
}

type TranslateResult struct {
	TranslatedText   string
	DetectedLanguage string
}

type deeplResponse struct {
	Translations []struct {
		Text               string `json:"text"`
		DetectedSourceLang string `json:"detected_source_language"`
	} `json:"translations"`
}

func (c *TranslateClient) Translate(ctx context.Context, text, targetLang string) (*TranslateResult, error) {
	// DeepL Free API uses api-free.deepl.com
	endpoint := "https://api-free.deepl.com/v2/translate"

	form := url.Values{}
	form.Set("text", text)
	form.Set("target_lang", strings.ToUpper(targetLang))

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, strings.NewReader(form.Encode()))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Authorization", "DeepL-Auth-Key "+c.apiKey)

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
		return nil, fmt.Errorf("DeepL API returned %d: %s", resp.StatusCode, string(body))
	}

	var dResp deeplResponse
	if err := json.Unmarshal(body, &dResp); err != nil {
		return nil, fmt.Errorf("invalid DeepL response: %w", err)
	}

	if len(dResp.Translations) == 0 {
		return nil, fmt.Errorf("no translations returned")
	}

	t := dResp.Translations[0]
	return &TranslateResult{
		TranslatedText:   t.Text,
		DetectedLanguage: strings.ToLower(t.DetectedSourceLang),
	}, nil
}
