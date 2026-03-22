package logger

import (
	"github.com/disgoorg/snowflake/v2"
	"github.com/s12kuma01/pedmin/store"
)

type LoggerSettings struct {
	ChannelID snowflake.ID    `json:"channel_id"`
	Events    map[string]bool `json:"events"`
}

const (
	EventMessageEdit   = "message_edit"
	EventMessageDelete = "message_delete"
	EventMemberJoin    = "member_join"
	EventMemberLeave   = "member_leave"
	EventBanAdd        = "ban_add"
	EventBanRemove     = "ban_remove"
	EventRoleChange    = "role_change"
	EventChannelChange = "channel_change"
)

var AllEvents = []struct {
	Key   string
	Label string
}{
	{EventMessageEdit, "メッセージ編集"},
	{EventMessageDelete, "メッセージ削除"},
	{EventMemberJoin, "メンバー参加"},
	{EventMemberLeave, "メンバー退出"},
	{EventBanAdd, "BAN"},
	{EventBanRemove, "BAN解除"},
	{EventRoleChange, "ロール変更"},
	{EventChannelChange, "チャンネル変更"},
}

func LoadSettings(guildStore store.GuildStore, guildID snowflake.ID) (*LoggerSettings, error) {
	s, err := store.LoadModuleSettings(guildStore, guildID, ModuleID, func() *LoggerSettings {
		return &LoggerSettings{Events: make(map[string]bool)}
	})
	if err != nil {
		return nil, err
	}
	if s.Events == nil {
		s.Events = make(map[string]bool)
	}
	return s, nil
}

func SaveSettings(guildStore store.GuildStore, guildID snowflake.ID, settings *LoggerSettings) error {
	return store.SaveModuleSettings(guildStore, guildID, ModuleID, settings)
}

func (s *LoggerSettings) IsEventEnabled(event string) bool {
	return s.Events[event]
}
