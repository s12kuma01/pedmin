// opyright (c) 2025-2026 s12kuma01
// SPDX-License-Identifier: MPL-2.0

package handler

import (
	"log/slog"
	"strconv"
	"strings"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
	"github.com/s12kuma01/pedmin/internal/ui"
	"github.com/s12kuma01/pedmin/internal/view"
)

func (h *RSSHandler) HandleComponent(e *events.ComponentInteractionCreate) {
	customID := e.Data.CustomID()
	_, rest, _ := strings.Cut(customID, ":")
	action, extra, _ := strings.Cut(rest, ":")

	guildID := e.GuildID()
	if guildID == nil {
		return
	}

	switch action {
	case "add_prompt":
		h.rssHandleAddPrompt(e)

	case "add_channel":
		h.rssHandleAddChannel(e, extra)

	case "manage":
		h.rssHandleManage(e)

	case "manage_select":
		h.rssHandleManageSelect(e)

	case "delete":
		h.rssHandleDelete(e, extra)
	}
}

func (h *RSSHandler) rssHandleManage(e *events.ComponentInteractionCreate) {
	feeds, err := h.service.GetFeeds(*e.GuildID())
	if err != nil {
		h.logger.Error("failed to get feeds", slog.Any("error", err))
		_ = e.CreateMessage(ui.EphemeralV2(view.RSSErrorContainer("フィード一覧の取得に失敗しました。")))
		return
	}

	_ = e.CreateMessage(view.RSSManagePanel(feeds))
}

func (h *RSSHandler) rssHandleManageSelect(e *events.ComponentInteractionCreate) {
	data, ok := e.Data.(discord.StringSelectMenuInteractionData)
	if !ok || len(data.Values) == 0 {
		return
	}

	feedID, err := strconv.ParseInt(data.Values[0], 10, 64)
	if err != nil {
		return
	}

	feed, err := h.service.GetFeed(*e.GuildID(), feedID)
	if err != nil {
		h.logger.Error("failed to get feed", slog.Any("error", err))
		_ = e.UpdateMessage(discord.NewMessageUpdateV2([]discord.LayoutComponent{
			view.RSSErrorContainer("フィードが見つかりませんでした。"),
		}))
		return
	}

	_ = e.UpdateMessage(view.RSSFeedDetail(*feed))
}

func (h *RSSHandler) rssHandleDelete(e *events.ComponentInteractionCreate, feedIDStr string) {
	feedID, err := strconv.ParseInt(feedIDStr, 10, 64)
	if err != nil {
		return
	}

	feeds, err := h.service.DeleteFeedAndList(feedID, *e.GuildID())
	if err != nil {
		h.logger.Error("failed to delete feed", slog.Any("error", err))
		_ = e.CreateMessage(ui.EphemeralV2(view.RSSErrorContainer("フィードの削除に失敗しました。")))
		return
	}

	if len(feeds) == 0 {
		_ = e.UpdateMessage(discord.NewMessageUpdateV2([]discord.LayoutComponent{
			discord.NewContainer(
				discord.NewTextDisplay("登録されているフィードはありません。"),
			),
		}))
		return
	}

	msg := view.RSSManagePanel(feeds)
	_ = e.UpdateMessage(discord.NewMessageUpdateV2(msg.Components))
}
