package service

import (
	"context"
	"net/url"

	"github.com/s12kuma01/pedmin/internal/client"
	"github.com/s12kuma01/pedmin/internal/model"
)

// URLService wraps the URL client with validation logic.
type URLService struct {
	client *client.URLClient
}

// NewURLService creates a new URLService.
func NewURLService(c *client.URLClient) *URLService {
	return &URLService{client: c}
}

// Shorten validates a URL and shortens it via x.gd.
func (s *URLService) Shorten(ctx context.Context, rawURL string) (string, error) {
	return s.client.ShortenURL(ctx, rawURL)
}

// Check validates a URL and scans it via VirusTotal.
func (s *URLService) Check(ctx context.Context, rawURL string) (*model.VTResult, error) {
	return s.client.ScanURL(ctx, rawURL)
}

// URLIsValid checks whether a raw URL string is a valid http or https URL.
func URLIsValid(rawURL string) bool {
	u, err := url.Parse(rawURL)
	if err != nil {
		return false
	}
	return u.Scheme == "http" || u.Scheme == "https"
}
