package logger

import (
	"log/slog"

	disgobot "github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
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

func (l *Logger) onMessageUpdate(e *events.GuildMessageUpdate) {
	if e.Message.Author.Bot {
		return
	}
	oldContent := e.OldMessage.Content
	newContent := e.Message.Content
	if oldContent == newContent && AttachmentsEqual(e.OldMessage.Attachments, e.Message.Attachments) {
		return
	}
	l.sendLog(e.GuildID, EventMessageEdit,
		BuildMessageEditLog(e.Message.Author, e.ChannelID, oldContent, newContent, e.OldMessage.Attachments, e.Message.Attachments),
	)
}

func (l *Logger) onMessageDelete(e *events.GuildMessageDelete) {
	var user *discord.User
	content := e.Message.Content
	if e.Message.Author.ID != 0 {
		user = &e.Message.Author
	}
	if user != nil && user.Bot {
		return
	}
	l.sendLog(e.GuildID, EventMessageDelete,
		BuildMessageDeleteLog(user, e.ChannelID, content, e.Message.Attachments),
	)
}

func (l *Logger) onMemberJoin(e *events.GuildMemberJoin) {
	l.sendLog(e.GuildID, EventMemberJoin,
		BuildMemberJoinLog(e.Member),
	)
}

func (l *Logger) onMemberLeave(e *events.GuildMemberLeave) {
	l.sendLog(e.GuildID, EventMemberLeave,
		BuildMemberLeaveLog(e.User),
	)
}

func (l *Logger) onBan(e *events.GuildBan) {
	l.sendLog(e.GuildID, EventBanAdd,
		BuildBanLog(e.User),
	)
}

func (l *Logger) onUnban(e *events.GuildUnban) {
	l.sendLog(e.GuildID, EventBanRemove,
		BuildUnbanLog(e.User),
	)
}

func (l *Logger) onRoleCreate(e *events.RoleCreate) {
	l.sendLog(e.GuildID, EventRoleChange,
		BuildRoleCreateLog(e.Role),
	)
}

func (l *Logger) onRoleUpdate(e *events.RoleUpdate) {
	l.sendLog(e.GuildID, EventRoleChange,
		BuildRoleUpdateLog(e.Role),
	)
}

func (l *Logger) onRoleDelete(e *events.RoleDelete) {
	l.sendLog(e.GuildID, EventRoleChange,
		BuildRoleDeleteLog(e.Role),
	)
}

func (l *Logger) onChannelCreate(e *events.GuildChannelCreate) {
	l.sendLog(e.GuildID, EventChannelChange,
		BuildChannelCreateLog(e.Channel),
	)
}

func (l *Logger) onChannelUpdate(e *events.GuildChannelUpdate) {
	l.sendLog(e.GuildID, EventChannelChange,
		BuildChannelUpdateLog(e.Channel),
	)
}

func (l *Logger) onChannelDelete(e *events.GuildChannelDelete) {
	l.sendLog(e.GuildID, EventChannelChange,
		BuildChannelDeleteLog(e.Channel),
	)
}
