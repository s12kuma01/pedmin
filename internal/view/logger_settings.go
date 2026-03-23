// opyright (c) 2025-2026 s12kuma01
// SPDX-License-Identifier: MPL-2.0

package view

import (
	"fmt"
	"strings"

	"github.com/disgoorg/disgo/discord"
	"github.com/s12kuma01/pedmin/internal/model"
)

// LoggerSettingsPanel builds the logger settings panel components.
func LoggerSettingsPanel(settings *model.LoggerSettings) []discord.LayoutComponent {
	// Build info text
	channelText := "未設定"
	if settings.ChannelID != 0 {
		channelText = fmt.Sprintf("<#%d>", settings.ChannelID)
	}

	var enabledNames []string
	for _, ev := range model.AllLogEvents {
		if settings.IsEventEnabled(ev.Key) {
			enabledNames = append(enabledNames, ev.Label)
		}
	}
	eventsText := "なし"
	if len(enabledNames) > 0 {
		eventsText = strings.Join(enabledNames, ", ")
	}

	infoDisplay := discord.NewTextDisplay(fmt.Sprintf(
		"**ログチャンネル:** %s\n**ログ対象:** %s",
		channelText, eventsText,
	))

	// Channel select menu
	channelSelect := discord.NewActionRow(
		discord.NewChannelSelectMenu(model.LoggerModuleID+":channel", "ログチャンネルを選択...").
			WithChannelTypes(discord.ChannelTypeGuildText),
	)

	// Event select menu
	var options []discord.StringSelectMenuOption
	for _, ev := range model.AllLogEvents {
		opt := discord.StringSelectMenuOption{
			Label: ev.Label,
			Value: ev.Key,
		}
		if settings.IsEventEnabled(ev.Key) {
			opt.Default = true
		}
		options = append(options, opt)
	}
	eventSelect := discord.NewActionRow(
		discord.NewStringSelectMenu(model.LoggerModuleID+":events", "ログ対象イベントを選択...", options...).
			WithMinValues(0).
			WithMaxValues(len(model.AllLogEvents)),
	)

	return []discord.LayoutComponent{
		infoDisplay,
		channelSelect,
		eventSelect,
	}
}
