package embedfix

import "fmt"

// Application Emoji (registered in Developer Portal)
const (
	// X/Twitter (existing)
	emojiMessages = "<:messages:1484573490400067635>"
	emojiRepost   = "<:repost:1484573523929469079>"
	emojiLike     = "<:like:1484573534490595419>"
	emojiGraph    = "<:graph:1484573514571714712>"
	emojiX        = "<:x_:1484598038508081293>"

	// Reddit (register and replace IDs)
	emojiReddit = "<:reddit:0000000000000000000>"
	emojiUpvote = "<:upvote:0000000000000000000>"

	// TikTok (register and replace IDs)
	emojiTikTok = "<:tiktok:0000000000000000000>"
	emojiPlay   = "<:play:0000000000000000000>"
	emojiShare  = "<:share:0000000000000000000>"
)

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
