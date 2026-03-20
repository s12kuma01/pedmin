package url

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type URLClient struct {
	http  *http.Client
	xgdKey string
	vtKey  string
}

func NewURLClient(xgdKey, vtKey string) *URLClient {
	return &URLClient{
		http:   &http.Client{Timeout: 10 * time.Second},
		xgdKey: xgdKey,
		vtKey:  vtKey,
	}
}

type VTResult struct {
	Harmless   int
	Malicious  int
	Suspicious int
	Undetected int
	Timeout    int
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
	defer resp.Body.Close()

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

// ScanURL queries VirusTotal for URL analysis results.
func (c *URLClient) ScanURL(ctx context.Context, rawURL string) (*VTResult, error) {
	// base64url encode the URL (no padding) as the ID
	urlID := base64.RawURLEncoding.EncodeToString([]byte(rawURL))

	// Try to get existing report
	result, err := c.vtGetReport(ctx, urlID)
	if err == nil {
		return result, nil
	}

	// Submit URL for scanning
	if err := c.vtSubmitURL(ctx, rawURL); err != nil {
		return nil, fmt.Errorf("スキャン送信に失敗: %w", err)
	}

	// Re-fetch report after submission
	result, err = c.vtGetReport(ctx, urlID)
	if err != nil {
		return nil, fmt.Errorf("レポート取得に失敗: %w", err)
	}

	return result, nil
}

func (c *URLClient) vtGetReport(ctx context.Context, urlID string) (*VTResult, error) {
	endpoint := "https://www.virustotal.com/api/v3/urls/" + urlID

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("x-apikey", c.vtKey)

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, fmt.Errorf("not found")
	}
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("VT API error (status %d): %s", resp.StatusCode, string(body))
	}

	var report struct {
		Data struct {
			Attributes struct {
				LastAnalysisStats struct {
					Harmless   int `json:"harmless"`
					Malicious  int `json:"malicious"`
					Suspicious int `json:"suspicious"`
					Undetected int `json:"undetected"`
					Timeout    int `json:"timeout"`
				} `json:"last_analysis_stats"`
			} `json:"attributes"`
		} `json:"data"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&report); err != nil {
		return nil, fmt.Errorf("failed to decode VT response: %w", err)
	}

	stats := report.Data.Attributes.LastAnalysisStats
	return &VTResult{
		Harmless:   stats.Harmless,
		Malicious:  stats.Malicious,
		Suspicious: stats.Suspicious,
		Undetected: stats.Undetected,
		Timeout:    stats.Timeout,
	}, nil
}

func (c *URLClient) vtSubmitURL(ctx context.Context, rawURL string) error {
	body := "url=" + url.QueryEscape(rawURL)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost,
		"https://www.virustotal.com/api/v3/urls", strings.NewReader(body))
	if err != nil {
		return err
	}
	req.Header.Set("x-apikey", c.vtKey)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := c.http.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		respBody, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("VT submit error (status %d): %s", resp.StatusCode, string(respBody))
	}

	return nil
}
