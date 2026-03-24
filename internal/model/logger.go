// Copyright (c) 2025-2026 s12kuma01
// SPDX-License-Identifier: MPL-2.0

package model

import "github.com/disgoorg/snowflake/v2"

// LoggerSettings holds per-guild logging configuration.
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

// AllLogEvents lists all supported log event types with display labels.
var AllLogEvents = []struct {
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

// IsEventEnabled checks if a specific event type is enabled.
func (s *LoggerSettings) IsEventEnabled(event string) bool {
	return s.Events[event]
}
