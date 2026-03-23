// opyright (c) 2025-2026 s12kuma01
// SPDX-License-Identifier: MPL-2.0

package client

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

type URLClient struct {
	http   *http.Client
	xgdKey string
	vtKey  string
}

func NewURLClient(xgdKey, vtKey string, timeout time.Duration) *URLClient {
	return &URLClient{
		http:   &http.Client{Timeout: timeout},
		xgdKey: xgdKey,
		vtKey:  vtKey,
	}
}

// ShortenURL calls the x.gd API to shorten a URL.
func (c *URLClient) ShortenURL(ctx context.Context, rawURL string) (string, error) {
	endpoint := fmt.Sprintf("https://xgd.io/V1/shorten?url=%s&key=%s",
		url.QueryEscape(rawURL), url.QueryEscape(c.xgdKey))

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return "", err
	}

	resp, err := c.http.Do(req)
	if err != nil {
		return "", err
	}
	defer func() { _ = resp.Body.Close() }()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var result struct {
		ShortURL string `json:"shorturl"`
		Error    string `json:"error"`
	}
	if err := json.Unmarshal(body, &result); err != nil {
		return "", fmt.Errorf("invalid response: %s", string(body))
	}

	if result.Error != "" {
		return "", fmt.Errorf("%s", result.Error)
	}
	if result.ShortURL == "" {
		return "", fmt.Errorf("empty short URL returned")
	}

	return result.ShortURL, nil
}
