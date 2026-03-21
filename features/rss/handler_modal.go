package rss

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
	"github.com/s12kuma01/pedmin/ui"
)

func (r *RSS) handleModal(e *events.ModalSubmitInteractionCreate) {
	customID := e.Data.CustomID
	_, action, _ := strings.Cut(customID, ":")

	if action != "add_modal" {
		return
	}

	feedURL := e.Data.Text(ModuleID + ":url")
	feedURL = strings.TrimSpace(feedURL)

	if feedURL == "" {
		_ = e.CreateMessage(ui.EphemeralV2(errorContainer("URLを入力してください。")))
		return
	}

	// Basic URL validation
	if !strings.HasPrefix(feedURL, "http://") && !strings.HasPrefix(feedURL, "https://") {
		feedURL = "https://" + feedURL
	}

	if _, err := url.ParseRequestURI(feedURL); err != nil {
		_ = e.CreateMessage(ui.EphemeralV2(errorContainer("無効なURLです。")))
		return
	}

	encodedURL := url.QueryEscape(feedURL)

	_ = e.CreateMessage(ui.EphemeralV2(
		discord.NewContainer(
			discord.NewTextDisplay(fmt.Sprintf("**フィードURL:** %s", feedURL)),
			discord.NewSmallSeparator(),
			discord.NewTextDisplay("配信先チャンネルを選択してください:"),
			discord.NewActionRow(
				discord.NewChannelSelectMenu(ModuleID+":add_channel:"+encodedURL, "チャンネルを選択...").
					WithChannelTypes(discord.ChannelTypeGuildText),
			),
		),
	))
}
