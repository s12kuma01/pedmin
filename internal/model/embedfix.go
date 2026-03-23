// opyright (c) 2025-2026 s12kuma01
// SPDX-License-Identifier: MPL-2.0

package model

import "regexp"

// Platform identifies which SNS a detected URL belongs to.
type Platform string

const (
	PlatformTwitter Platform = "twitter"
	PlatformReddit  Platform = "reddit"
	PlatformTikTok  Platform = "tiktok"
)

// EmbedRef represents a detected SNS URL with platform-specific parameters.
type EmbedRef struct {
	Platform Platform
	Params   []string // platform-specific captured groups
}

// EmbedFixSettings holds per-guild embedfix configuration.
type EmbedFixSettings struct {
	Platforms map[Platform]bool `json:"platforms"`
}

// AllPlatforms lists all supported platforms with display labels.
var AllPlatforms = []struct {
	Key   Platform
	Label string
}{
	{PlatformTwitter, "X / Twitter"},
	{PlatformReddit, "Reddit"},
	{PlatformTikTok, "TikTok"},
}

// DefaultEmbedFixSettings returns settings with all platforms enabled.
func DefaultEmbedFixSettings() *EmbedFixSettings {
	platforms := make(map[Platform]bool, len(AllPlatforms))
	for _, p := range AllPlatforms {
		platforms[p.Key] = true
	}
	return &EmbedFixSettings{Platforms: platforms}
}

// IsPlatformEnabled checks if a platform is enabled in the settings.
func (s *EmbedFixSettings) IsPlatformEnabled(p Platform) bool {
	enabled, ok := s.Platforms[p]
	if !ok {
		return true // enabled by default
	}
	return enabled
}

// URL regex matchers for embed URL extraction.
var (
	TweetURLRegex  = regexp.MustCompile(`https?://(?:www\.)?(?:twitter\.com|x\.com)/(\w+)/status/(\d+)`)
	RedditURLRegex = regexp.MustCompile(`https?://(?:www\.|old\.|new\.)?reddit\.com/r/(\w+)/comments/(\w+)`)
	TikTokURLRegex = regexp.MustCompile(`https?://(?:www\.)?tiktok\.com/@([\w.]+)/video/(\d+)`)
)

// MaxEmbedURLs limits total embed URLs processed per message.
const MaxEmbedURLs = 4

// URLMatcher pairs a regex with a platform for URL extraction.
type URLMatcher struct {
	Regex    *regexp.Regexp
	Platform Platform
}

// URLMatchers is the ordered list of URL matchers for embed extraction.
var URLMatchers = []URLMatcher{
	{TweetURLRegex, PlatformTwitter},
	{RedditURLRegex, PlatformReddit},
	{TikTokURLRegex, PlatformTikTok},
}

// ExtractEmbedURLs extracts SNS URLs from message content.
func ExtractEmbedURLs(content string) []EmbedRef {
	var refs []EmbedRef
	for _, m := range URLMatchers {
		matches := m.Regex.FindAllStringSubmatch(content, MaxEmbedURLs-len(refs))
		for _, match := range matches {
			refs = append(refs, EmbedRef{
				Platform: m.Platform,
				Params:   match[1:], // captured groups only
			})
			if len(refs) >= MaxEmbedURLs {
				return refs
			}
		}
	}
	return refs
}
