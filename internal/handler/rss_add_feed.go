// opyright (c) 2025-2026 s12kuma01
// SPDX-License-Identifier: MPL-2.0

package handler

import (
	"context"
	"fmt"
	"log/slog"
	"net/url"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
	"github.com/s12kuma01/pedmin/internal/model"
	"github.com/s12kuma01/pedmin/internal/ui"
	"github.com/s12kuma01/pedmin/internal/view"
)

func (h *RSSHandler) rssHandleAddPrompt(e *events.ComponentInteractionCreate) {
	_ = e.Modal(discord.ModalCreate{
		CustomID: model.RSSModuleID + ":add_modal",
		Title:    "RSSフィード追加",
		Components: []discord.LayoutComponent{
			discord.NewLabel("フィードURL",
				discord.NewShortTextInput(model.RSSModuleID+":url").
					WithRequired(true).
					WithPlaceholder("https://example.com/feed.xml"),
			),
		},
	})
}

func (h *RSSHandler) rssHandleAddChannel(e *events.ComponentInteractionCreate, encodedURL string) {
	feedURL, err := url.QueryUnescape(encodedURL)
	if err != nil {
		_ = e.CreateMessage(ui.EphemeralV2(view.RSSErrorContainer("無効なURLです。")))
		return
	}

	data, ok := e.Data.(discord.ChannelSelectMenuInteractionData)
	if !ok || len(data.Values) == 0 {
		return
	}
	channelID := data.Values[0]

	_ = e.DeferCreateMessage(true)

	ctx, cancel := context.WithTimeout(context.Background(), h.feedTimeout_)
	defer cancel()

	feed, err := h.service.AddFeed(ctx, *e.GuildID(), channelID, feedURL)
	if err != nil {
		h.logger.Error("failed to add feed", slog.Any("error", err))
		_, _ = e.Client().Rest.UpdateInteractionResponse(e.ApplicationID(), e.Token(), discord.NewMessageUpdateV2([]discord.LayoutComponent{
			view.RSSErrorContainer(fmt.Sprintf("フィード追加に失敗しました:\n%s", err.Error())),
		}))
		return
	}

	_, _ = e.Client().Rest.UpdateInteractionResponse(e.ApplicationID(), e.Token(), discord.NewMessageUpdateV2([]discord.LayoutComponent{
		discord.NewContainer(
			discord.NewTextDisplay(fmt.Sprintf("**%s** を <#%d> に追加しました。", feed.Title, feed.ChannelID)),
		),
	}))
}
