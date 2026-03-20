package rss

import (
	"fmt"

	"github.com/disgoorg/disgo/discord"
	"github.com/mmcdole/gofeed"
)

func BuildFeedAnnouncement(feedTitle string, item *gofeed.Item) discord.MessageCreate {
	desc := truncate(stripHTML(item.Description), 300)

	body := fmt.Sprintf("**%s**", item.Title)
	if desc != "" {
		body += "\n" + desc
	}
	if item.Link != "" {
		body += fmt.Sprintf("\n\n[続きを読む](%s)", item.Link)
	}

	return discord.NewMessageCreateV2(
		discord.NewContainer(
			discord.NewTextDisplay(fmt.Sprintf("### %s", feedTitle)),
			discord.NewSmallSeparator(),
			discord.NewTextDisplay(body),
		),
	).WithAllowedMentions(&discord.AllowedMentions{})
}
