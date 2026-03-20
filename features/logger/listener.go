package logger

import (
	"log/slog"

	disgobot "github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/snowflake/v2"
)

func SetupListeners(client *disgobot.Client, l *Logger) {
	client.AddEventListeners(
		disgobot.NewListenerFunc(l.onMessageUpdate),
		disgobot.NewListenerFunc(l.onMessageDelete),
		disgobot.NewListenerFunc(l.onMemberJoin),
		disgobot.NewListenerFunc(l.onMemberLeave),
		disgobot.NewListenerFunc(l.onBan),
		disgobot.NewListenerFunc(l.onUnban),
		disgobot.NewListenerFunc(l.onRoleCreate),
		disgobot.NewListenerFunc(l.onRoleUpdate),
		disgobot.NewListenerFunc(l.onRoleDelete),
		disgobot.NewListenerFunc(l.onChannelCreate),
		disgobot.NewListenerFunc(l.onChannelUpdate),
		disgobot.NewListenerFunc(l.onChannelDelete),
	)
}

func (l *Logger) sendLog(guildID snowflake.ID, event string, msg discord.MessageCreate) {
	if !l.bot.IsModuleEnabled(guildID, ModuleID) {
		return
	}

	settings, err := LoadSettings(l.store, guildID)
	if err != nil {
		l.logger.Error("failed to load logger settings", slog.Any("error", err))
		return
	}

	if settings.ChannelID == 0 || !settings.IsEventEnabled(event) {
		return
	}

	if _, err := l.client.Rest.CreateMessage(settings.ChannelID, msg); err != nil {
		l.logger.Error("failed to send log message",
			slog.String("event", event),
			slog.Any("error", err),
		)
	}
}
