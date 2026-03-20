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

// ScanURL queries VirusTotal for URL analysis results.
func (c *URLClient) ScanURL(ctx context.Context, rawURL string) (*VTResult, error) {
	// base64url encode the URL (no padding) as the ID
	urlID := base64.RawURLEncoding.EncodeToString([]byte(rawURL))

	// Try to get existing report
	result, err := c.vtGetReport(ctx, urlID)
	if err == nil {
		return result, nil
	}

	// Submit URL for scanning and get analysis ID
	analysisID, err := c.vtSubmitURL(ctx, rawURL)
	if err != nil {
		return nil, fmt.Errorf("スキャン送信に失敗: %w", err)
	}

	// Poll analysis endpoint until completed
	result, err = c.vtWaitAnalysis(ctx, analysisID)
	if err != nil {
		return nil, fmt.Errorf("スキャン結果の取得に失敗: %w", err)
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
	defer func() { _ = resp.Body.Close() }()

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

func (c *URLClient) vtSubmitURL(ctx context.Context, rawURL string) (string, error) {
	body := "url=" + url.QueryEscape(rawURL)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost,
		"https://www.virustotal.com/api/v3/urls", strings.NewReader(body))
	if err != nil {
		return "", err
	}
	req.Header.Set("x-apikey", c.vtKey)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := c.http.Do(req)
	if err != nil {
		return "", err
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		respBody, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("VT submit error (status %d): %s", resp.StatusCode, string(respBody))
	}

	var result struct {
		Data struct {
			ID string `json:"id"`
		} `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", fmt.Errorf("failed to decode submit response: %w", err)
	}

	return result.Data.ID, nil
}

func (c *URLClient) vtWaitAnalysis(ctx context.Context, analysisID string) (*VTResult, error) {
	endpoint := "https://www.virustotal.com/api/v3/analyses/" + analysisID

	for {
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
		if err != nil {
			return nil, err
		}
		req.Header.Set("x-apikey", c.vtKey)

		resp, err := c.http.Do(req)
		if err != nil {
			return nil, err
		}

		var result struct {
			Data struct {
				Attributes struct {
					Status string `json:"status"`
					Stats  struct {
						Harmless   int `json:"harmless"`
						Malicious  int `json:"malicious"`
						Suspicious int `json:"suspicious"`
						Undetected int `json:"undetected"`
						Timeout    int `json:"timeout"`
					} `json:"stats"`
				} `json:"attributes"`
			} `json:"data"`
		}

		if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
			_ = resp.Body.Close()
			return nil, fmt.Errorf("failed to decode analysis response: %w", err)
		}
		_ = resp.Body.Close()

		if result.Data.Attributes.Status == "completed" {
			stats := result.Data.Attributes.Stats
			return &VTResult{
				Harmless:   stats.Harmless,
				Malicious:  stats.Malicious,
				Suspicious: stats.Suspicious,
				Undetected: stats.Undetected,
				Timeout:    stats.Timeout,
			}, nil
		}

		// Wait before retrying
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case <-time.After(3 * time.Second):
		}
	}
}
