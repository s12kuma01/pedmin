package rss

import (
	"fmt"

	"github.com/disgoorg/disgo/discord"
	"github.com/s12kuma01/pedmin/store"
)

func BuildManagePanel(feeds []store.RSSFeed) discord.MessageCreate {
	if len(feeds) == 0 {
		return ephemeralV2(
			discord.NewContainer(
				discord.NewTextDisplay("登録されているフィードはありません。"),
			),
		)
	}

	var options []discord.StringSelectMenuOption
	for _, f := range feeds {
		label := f.Title
		if label == "" {
			label = f.URL
		}
		options = append(options, discord.StringSelectMenuOption{
			Label:       truncate(label, 100),
			Value:       fmt.Sprintf("%d", f.ID),
			Description: fmt.Sprintf("→ #%d", f.ChannelID),
		})
	}

	return ephemeralV2(
		discord.NewContainer(
			discord.NewTextDisplay("### RSSフィード管理"),
			discord.NewSmallSeparator(),
			discord.NewTextDisplay(fmt.Sprintf("登録フィード: %d/%d", len(feeds), MaxFeedsPerGuild)),
			discord.NewActionRow(
				discord.NewStringSelectMenu(ModuleID+":manage_select", "フィードを選択...", options...),
			),
		),
	)
}

func BuildFeedDetail(feed store.RSSFeed) discord.MessageUpdate {
	title := feed.Title
	if title == "" {
		title = feed.URL
	}

	return discord.NewMessageUpdateV2([]discord.LayoutComponent{
		discord.NewContainer(
			discord.NewTextDisplay(fmt.Sprintf("### %s", title)),
			discord.NewSmallSeparator(),
			discord.NewTextDisplay(fmt.Sprintf(
				"**URL:** %s\n**配信先:** <#%d>\n**追加日:** %s",
				feed.URL, feed.ChannelID, feed.AddedAt.Format("2006-01-02"),
			)),
			discord.NewLargeSeparator(),
			discord.NewActionRow(
				discord.NewDangerButton("削除", fmt.Sprintf("%s:delete:%d", ModuleID, feed.ID)),
			),
		),
	})
}
