package logger

import "github.com/disgoorg/disgo/events"

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
