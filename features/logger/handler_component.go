package logger

import (
	"log/slog"
	"strings"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
)

func (l *Logger) handleComponent(e *events.ComponentInteractionCreate) {
	customID := e.Data.CustomID()
	_, action, _ := strings.Cut(customID, ":")

	guildID := e.GuildID()
	if guildID == nil {
		return
	}

	settings, err := LoadSettings(l.store, *guildID)
	if err != nil {
		l.logger.Error("failed to load logger settings", slog.Any("error", err))
		return
	}

	switch action {
	case "channel":
		data, ok := e.Data.(discord.ChannelSelectMenuInteractionData)
		if !ok {
			return
		}
		if len(data.Values) > 0 {
			settings.ChannelID = data.Values[0]
		}

	case "events":
		data, ok := e.Data.(discord.StringSelectMenuInteractionData)
		if !ok {
			return
		}
		for k := range settings.Events {
			settings.Events[k] = false
		}
		for _, v := range data.Values {
			settings.Events[v] = true
		}

	default:
		return
	}

	if err := SaveSettings(l.store, *guildID, settings); err != nil {
		l.logger.Error("failed to save logger settings", slog.Any("error", err))
	}

	_ = e.DeferUpdateMessage()
}
