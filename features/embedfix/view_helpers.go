package embedfix

import (
	"fmt"
	"regexp"
)

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

// Application Emoji (registered in Developer Portal)
const (
	emojiMessages = "<:messages:1484573490400067635>"
	emojiRepost   = "<:repost:1484573523929469079>"
	emojiLike     = "<:like:1484573534490595419>"
	emojiGraph    = "<:graph:1484573514571714712>"
	emojiX        = "<:x_:1484598038508081293>"
)

var (
	tweetURLRegex     = regexp.MustCompile(`https?://(?:www\.)?(?:twitter\.com|x\.com)/(\w+)/status/(\d+)`)
	redditURLRegex    = regexp.MustCompile(`https?://(?:www\.|old\.|new\.)?reddit\.com/r/(\w+)/comments/(\w+)`)
	tiktokURLRegex    = regexp.MustCompile(`https?://(?:www\.)?tiktok\.com/@([\w.]+)/video/(\d+)`)
	instagramURLRegex = regexp.MustCompile(`https?://(?:www\.)?instagram\.com/(?:p|reel|tv)/([\w-]+)`)
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
	{instagramURLRegex, PlatformInstagram},
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

func formatCount(n int) string {
	switch {
	case n >= 100_000_000:
		return fmt.Sprintf("%.1f億", float64(n)/100_000_000)
	case n >= 10_000:
		return fmt.Sprintf("%.1f万", float64(n)/10_000)
	case n >= 1_000:
		return fmt.Sprintf("%.1fK", float64(n)/1_000)
	default:
		return fmt.Sprintf("%d", n)
	}
}

var langNames = map[string]string{
	"en": "英語",
	"es": "スペイン語",
	"fr": "フランス語",
	"de": "ドイツ語",
	"it": "イタリア語",
	"pt": "ポルトガル語",
	"ru": "ロシア語",
	"ko": "韓国語",
	"zh": "中国語",
	"ar": "アラビア語",
	"hi": "ヒンディー語",
	"th": "タイ語",
	"vi": "ベトナム語",
	"id": "インドネシア語",
	"tr": "トルコ語",
	"nl": "オランダ語",
	"pl": "ポーランド語",
	"sv": "スウェーデン語",
	"da": "デンマーク語",
	"fi": "フィンランド語",
	"no": "ノルウェー語",
	"uk": "ウクライナ語",
}

func langName(code string) string {
	if name, ok := langNames[code]; ok {
		return name
	}
	return code
}
