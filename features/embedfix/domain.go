package embedfix

import "regexp"

// Platform identifies which SNS a detected URL belongs to.
type Platform string

const (
	PlatformTwitter   Platform = "twitter"
	PlatformReddit    Platform = "reddit"
	PlatformTikTok    Platform = "tiktok"
	PlatformInstagram Platform = "instagram"
)

// EmbedRef represents a detected SNS URL with platform-specific parameters.
type EmbedRef struct {
	Platform Platform
	Params   []string // platform-specific captured groups
}

var (
	tweetURLRegex     = regexp.MustCompile(`https?://(?:www\.)?(?:twitter\.com|x\.com)/(\w+)/status/(\d+)`)
	redditURLRegex    = regexp.MustCompile(`https?://(?:www\.|old\.|new\.)?reddit\.com/r/(\w+)/comments/(\w+)`)
	tiktokURLRegex = regexp.MustCompile(`https?://(?:www\.)?tiktok\.com/@([\w.]+)/video/(\d+)`)
)

// maxEmbedURLs limits total embed URLs processed per message.
const maxEmbedURLs = 4

type urlMatcher struct {
	regex    *regexp.Regexp
	platform Platform
}

var urlMatchers = []urlMatcher{
	{tweetURLRegex, PlatformTwitter},
	{redditURLRegex, PlatformReddit},
	{tiktokURLRegex, PlatformTikTok},
}

func extractEmbedURLs(content string) []EmbedRef {
	var refs []EmbedRef
	for _, m := range urlMatchers {
		matches := m.regex.FindAllStringSubmatch(content, maxEmbedURLs-len(refs))
		for _, match := range matches {
			refs = append(refs, EmbedRef{
				Platform: m.platform,
				Params:   match[1:], // captured groups only
			})
			if len(refs) >= maxEmbedURLs {
				return refs
			}
		}
	}
	return refs
}
