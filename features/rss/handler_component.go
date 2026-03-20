package rss

import (
	"context"
	"fmt"
	"log/slog"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
)

func (r *RSS) handleComponent(e *events.ComponentInteractionCreate) {
	customID := e.Data.CustomID()
	_, rest, _ := strings.Cut(customID, ":")
	action, extra, _ := strings.Cut(rest, ":")

	guildID := e.GuildID()
	if guildID == nil {
		return
	}

	switch action {
	case "add_prompt":
		r.handleAddPrompt(e)

	case "add_channel":
		r.handleAddChannel(e, extra)

	case "manage":
		r.handleManage(e)

	case "manage_select":
		r.handleManageSelect(e)

	case "delete":
		r.handleDelete(e, extra)
	}
}

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
		_ = e.CreateMessage(ephemeralV2(errorContainer("無効なURLです。")))
		return
	}

	data, ok := e.Data.(discord.ChannelSelectMenuInteractionData)
	if !ok || len(data.Values) == 0 {
		return
	}
	channelID := data.Values[0]

	_ = e.DeferCreateMessage(true)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
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

func (r *RSS) handleManage(e *events.ComponentInteractionCreate) {
	feeds, err := r.GetFeeds(*e.GuildID())
	if err != nil {
		r.logger.Error("failed to get feeds", slog.Any("error", err))
		_ = e.CreateMessage(ephemeralV2(errorContainer("フィード一覧の取得に失敗しました。")))
		return
	}

	_ = e.CreateMessage(BuildManagePanel(feeds))
}

func (r *RSS) handleManageSelect(e *events.ComponentInteractionCreate) {
	data, ok := e.Data.(discord.StringSelectMenuInteractionData)
	if !ok || len(data.Values) == 0 {
		return
	}

	feedID, err := strconv.ParseInt(data.Values[0], 10, 64)
	if err != nil {
		return
	}

	feed, err := r.GetFeed(*e.GuildID(), feedID)
	if err != nil {
		r.logger.Error("failed to get feed", slog.Any("error", err))
		_ = e.UpdateMessage(discord.NewMessageUpdateV2([]discord.LayoutComponent{
			errorContainer("フィードが見つかりませんでした。"),
		}))
		return
	}

	_ = e.UpdateMessage(BuildFeedDetail(*feed))
}

func (r *RSS) handleDelete(e *events.ComponentInteractionCreate, feedIDStr string) {
	feedID, err := strconv.ParseInt(feedIDStr, 10, 64)
	if err != nil {
		return
	}

	feeds, err := r.DeleteFeedAndList(feedID, *e.GuildID())
	if err != nil {
		r.logger.Error("failed to delete feed", slog.Any("error", err))
		_ = e.CreateMessage(ephemeralV2(errorContainer("フィードの削除に失敗しました。")))
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

	msg := BuildManagePanel(feeds)
	_ = e.UpdateMessage(discord.NewMessageUpdateV2(msg.Components))
}
