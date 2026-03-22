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
