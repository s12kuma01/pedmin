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
	emojiReddit = "<:reddit:1485288295884914688>"
	emojiUpvote = "<:trend:1485288337026846770>"

	// TikTok (register and replace IDs)
	emojiTikTok = "<:tiktok:1485288359289946223>"
	emojiPlay   = "<:play:1485288373047263424>"
	emojiShare  = "<:send:1485288382509875231>"
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
