package url

import (
	"context"
	"net/url"
)

// Shorten validates a URL and shortens it via x.gd.
func (u *URL) Shorten(ctx context.Context, rawURL string) (string, error) {
	return u.client.ShortenURL(ctx, rawURL)
}

// Check validates a URL and scans it via VirusTotal.
func (u *URL) Check(ctx context.Context, rawURL string) (*VTResult, error) {
	return u.client.ScanURL(ctx, rawURL)
}

func isValidURL(rawURL string) bool {
	u, err := url.Parse(rawURL)
	if err != nil {
		return false
	}
	return u.Scheme == "http" || u.Scheme == "https"
}
