// opyright (c) 2025-2026 s12kuma01
// SPDX-License-Identifier: MPL-2.0

package view

import (
	"fmt"

	"github.com/disgoorg/disgo/discord"
)

// LoggerRoleCreateLog builds the role create log entry.
func LoggerRoleCreateLog(role discord.Role) discord.MessageCreate {
	return loggerBuildRoleLog("作成", role)
}

// LoggerRoleUpdateLog builds the role update log entry.
func LoggerRoleUpdateLog(role discord.Role) discord.MessageCreate {
	return loggerBuildRoleLog("更新", role)
}

// LoggerRoleDeleteLog builds the role delete log entry.
func LoggerRoleDeleteLog(role discord.Role) discord.MessageCreate {
	return loggerBuildRoleLog("削除", role)
}

func loggerBuildRoleLog(action string, role discord.Role) discord.MessageCreate {
	colorText := "なし"
	if role.Color != 0 {
		colorText = fmt.Sprintf("#%06X", role.Color)
	}
	return loggerMessage(
		fmt.Sprintf("### 🏷️ ロール%s", action),
		fmt.Sprintf("**ロール:** %s\n**色:** %s",
			role.Name, colorText),
	)
}

// LoggerChannelCreateLog builds the channel create log entry.
func LoggerChannelCreateLog(channel discord.GuildChannel) discord.MessageCreate {
	return loggerBuildChannelLog("作成", channel)
}

// LoggerChannelUpdateLog builds the channel update log entry.
func LoggerChannelUpdateLog(channel discord.GuildChannel) discord.MessageCreate {
	return loggerBuildChannelLog("更新", channel)
}

// LoggerChannelDeleteLog builds the channel delete log entry.
func LoggerChannelDeleteLog(channel discord.GuildChannel) discord.MessageCreate {
	return loggerBuildChannelLog("削除", channel)
}

func loggerBuildChannelLog(action string, channel discord.GuildChannel) discord.MessageCreate {
	return loggerMessage(
		fmt.Sprintf("### 📁 チャンネル%s", action),
		fmt.Sprintf("**チャンネル:** %s\n**タイプ:** %s",
			channel.Name(), loggerChannelTypeName(channel.Type())),
	)
}

func loggerChannelTypeName(t discord.ChannelType) string {
	switch t {
	case discord.ChannelTypeGuildText:
		return "テキスト"
	case discord.ChannelTypeGuildVoice:
		return "ボイス"
	case discord.ChannelTypeGuildCategory:
		return "カテゴリ"
	case discord.ChannelTypeGuildNews:
		return "ニュース"
	case discord.ChannelTypeGuildStageVoice:
		return "ステージ"
	case discord.ChannelTypeGuildForum:
		return "フォーラム"
	default:
		return "その他"
	}
}
