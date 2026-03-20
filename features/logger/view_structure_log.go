package logger

import (
	"fmt"

	"github.com/disgoorg/disgo/discord"
)

func BuildRoleCreateLog(role discord.Role) discord.MessageCreate {
	return buildRoleLog("作成", role)
}

func BuildRoleUpdateLog(role discord.Role) discord.MessageCreate {
	return buildRoleLog("更新", role)
}

func BuildRoleDeleteLog(role discord.Role) discord.MessageCreate {
	return buildRoleLog("削除", role)
}

func buildRoleLog(action string, role discord.Role) discord.MessageCreate {
	colorText := "なし"
	if role.Color != 0 {
		colorText = fmt.Sprintf("#%06X", role.Color)
	}
	return logMessage(
		fmt.Sprintf("### 🏷️ ロール%s", action),
		fmt.Sprintf("**ロール:** %s\n**色:** %s",
			role.Name, colorText),
	)
}

func BuildChannelCreateLog(channel discord.GuildChannel) discord.MessageCreate {
	return buildChannelLog("作成", channel)
}

func BuildChannelUpdateLog(channel discord.GuildChannel) discord.MessageCreate {
	return buildChannelLog("更新", channel)
}

func BuildChannelDeleteLog(channel discord.GuildChannel) discord.MessageCreate {
	return buildChannelLog("削除", channel)
}

func buildChannelLog(action string, channel discord.GuildChannel) discord.MessageCreate {
	return logMessage(
		fmt.Sprintf("### 📁 チャンネル%s", action),
		fmt.Sprintf("**チャンネル:** %s\n**タイプ:** %s",
			channel.Name(), channelTypeName(channel.Type())),
	)
}

func channelTypeName(t discord.ChannelType) string {
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
