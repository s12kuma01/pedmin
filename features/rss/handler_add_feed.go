package rss

import (
	"context"
	"fmt"
	"log/slog"
	"net/url"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
	"github.com/s12kuma01/pedmin/ui"
)

func (r *RSS) handleAddPrompt(e *events.ComponentInteractionCreate) {
	_ = e.Modal(discord.ModalCreate{
		CustomID: ModuleID + ":add_modal",
		Title:    "RSSフィード追加",
		Components: []discord.LayoutComponent{
			discord.NewLabel("フィードURL",
				discord.NewShortTextInput(ModuleID+":url").
					WithRequired(true).
					WithPlaceholder("https://example.com/feed.xml"),
			),
		},
	})
}

func (r *RSS) handleAddChannel(e *events.ComponentInteractionCreate, encodedURL string) {
	feedURL, err := url.QueryUnescape(encodedURL)
	if err != nil {
		_ = e.CreateMessage(ui.EphemeralV2(errorContainer("無効なURLです。")))
		return
	}

	data, ok := e.Data.(discord.ChannelSelectMenuInteractionData)
	if !ok || len(data.Values) == 0 {
		return
	}
	channelID := data.Values[0]

	_ = e.DeferCreateMessage(true)

	ctx, cancel := context.WithTimeout(context.Background(), r.feedTimeout)
	defer cancel()

	feed, err := r.AddFeed(ctx, *e.GuildID(), channelID, feedURL)
	if err != nil {
		r.logger.Error("failed to add feed", slog.Any("error", err))
		_, _ = e.Client().Rest.UpdateInteractionResponse(e.ApplicationID(), e.Token(), discord.NewMessageUpdateV2([]discord.LayoutComponent{
			errorContainer(fmt.Sprintf("フィード追加に失敗しました:\n%s", err.Error())),
		}))
		return
	}

	_, _ = e.Client().Rest.UpdateInteractionResponse(e.ApplicationID(), e.Token(), discord.NewMessageUpdateV2([]discord.LayoutComponent{
		discord.NewContainer(
			discord.NewTextDisplay(fmt.Sprintf("**%s** を <#%d> に追加しました。", feed.Title, feed.ChannelID)),
		),
	}))
}
