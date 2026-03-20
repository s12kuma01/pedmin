package logger

import (
	"fmt"
	"strings"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/snowflake/v2"
)

const (
	colorMessageEdit   = 0xF39C12
	colorMessageDelete = 0xE74C3C
	colorMemberJoin    = 0x2ECC71
	colorMemberLeave   = 0xE67E22
	colorBan           = 0xE74C3C
	colorUnban         = 0x3498DB
	colorRole          = 0x9B59B6
	colorChannel       = 0x1ABC9C
)

func BuildMessageEditLog(user discord.User, channelID snowflake.ID, oldContent, newContent string, oldAttachments, newAttachments []discord.Attachment) discord.MessageCreate {
	title := "### ✏️ メッセージ編集"
	body := fmt.Sprintf("**ユーザー:** <@%d>\n**チャンネル:** <#%d>\n**変更前:**\n> %s\n**変更後:**\n> %s",
		user.ID, channelID, oldContent, newContent)

	components := []discord.ContainerSubComponent{
		discord.NewTextDisplay(title),
		discord.NewSmallSeparator(),
		discord.NewTextDisplay(body),
	}

	removed, added := diffAttachments(oldAttachments, newAttachments)
	if len(removed) > 0 {
		components = append(components, discord.NewSmallSeparator())
		components = append(components, discord.NewTextDisplay("**削除された添付ファイル:**"))
		components = append(components, buildAttachmentComponents(removed)...)
	}
	if len(added) > 0 {
		components = append(components, discord.NewSmallSeparator())
		components = append(components, discord.NewTextDisplay("**追加された添付ファイル:**"))
		components = append(components, buildAttachmentComponents(added)...)
	}

	return discord.NewMessageCreateV2(
		discord.NewContainer(components...).WithAccentColor(colorMessageEdit),
	).WithAllowedMentions(&discord.AllowedMentions{})
}

func BuildMessageDeleteLog(user *discord.User, channelID snowflake.ID, content string, attachments []discord.Attachment) discord.MessageCreate {
	userText := "*不明*"
	if user != nil {
		userText = fmt.Sprintf("<@%d>", user.ID)
	}
	contentText := content
	if contentText == "" {
		contentText = "*内容を取得できませんでした*"
	}

	title := "### 🗑️ メッセージ削除"
	body := fmt.Sprintf("**ユーザー:** %s\n**チャンネル:** <#%d>\n**内容:**\n> %s",
		userText, channelID, contentText)

	components := []discord.ContainerSubComponent{
		discord.NewTextDisplay(title),
		discord.NewSmallSeparator(),
		discord.NewTextDisplay(body),
	}

	if len(attachments) > 0 {
		components = append(components, discord.NewSmallSeparator())
		components = append(components, discord.NewTextDisplay("**添付ファイル:**"))
		components = append(components, buildAttachmentComponents(attachments)...)
	}

	return discord.NewMessageCreateV2(
		discord.NewContainer(components...).WithAccentColor(colorMessageDelete),
	).WithAllowedMentions(&discord.AllowedMentions{})
}

func BuildMemberJoinLog(member discord.Member) discord.MessageCreate {
	createdAt := member.User.CreatedAt().Format("2006-01-02")
	return logMessage(colorMemberJoin,
		"### 📥 メンバー参加",
		fmt.Sprintf("**ユーザー:** <@%d> (%s)\n**アカウント作成:** %s",
			member.User.ID, member.User.Username, createdAt),
	)
}

func BuildMemberLeaveLog(user discord.User) discord.MessageCreate {
	return logMessage(colorMemberLeave,
		"### 📤 メンバー退出",
		fmt.Sprintf("**ユーザー:** <@%d> (%s)",
			user.ID, user.Username),
	)
}

func BuildBanLog(user discord.User) discord.MessageCreate {
	return logMessage(colorBan,
		"### 🔨 BAN",
		fmt.Sprintf("**ユーザー:** <@%d> (ID: %d)",
			user.ID, user.ID),
	)
}

func BuildUnbanLog(user discord.User) discord.MessageCreate {
	return logMessage(colorUnban,
		"### 🔓 BAN解除",
		fmt.Sprintf("**ユーザー:** <@%d> (ID: %d)",
			user.ID, user.ID),
	)
}

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
	return logMessage(colorRole,
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
	return logMessage(colorChannel,
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

func buildAttachmentComponents(attachments []discord.Attachment) []discord.ContainerSubComponent {
	var images []discord.MediaGalleryItem
	var files []string

	for _, a := range attachments {
		if a.ContentType != nil && strings.HasPrefix(*a.ContentType, "image/") {
			images = append(images, discord.MediaGalleryItem{
				Media: discord.UnfurledMediaItem{URL: a.URL},
			})
		} else {
			size := formatFileSize(a.Size)
			files = append(files, fmt.Sprintf("📎 %s (%s)", a.Filename, size))
		}
	}

	var components []discord.ContainerSubComponent
	if len(images) > 0 {
		components = append(components, discord.NewMediaGallery(images...))
	}
	if len(files) > 0 {
		components = append(components, discord.NewTextDisplay(strings.Join(files, "\n")))
	}
	return components
}

func formatFileSize(bytes int) string {
	switch {
	case bytes >= 1024*1024:
		return fmt.Sprintf("%.1f MB", float64(bytes)/(1024*1024))
	case bytes >= 1024:
		return fmt.Sprintf("%.1f KB", float64(bytes)/1024)
	default:
		return fmt.Sprintf("%d B", bytes)
	}
}

func diffAttachments(old, new []discord.Attachment) (removed, added []discord.Attachment) {
	oldIDs := make(map[snowflake.ID]discord.Attachment, len(old))
	for _, a := range old {
		oldIDs[a.ID] = a
	}
	newIDs := make(map[snowflake.ID]struct{}, len(new))
	for _, a := range new {
		newIDs[a.ID] = struct{}{}
		if _, exists := oldIDs[a.ID]; !exists {
			added = append(added, a)
		}
	}
	for _, a := range old {
		if _, exists := newIDs[a.ID]; !exists {
			removed = append(removed, a)
		}
	}
	return
}

func AttachmentsEqual(old, new []discord.Attachment) bool {
	if len(old) != len(new) {
		return false
	}
	for i := range old {
		if old[i].ID != new[i].ID {
			return false
		}
	}
	return true
}

func logMessage(color int, title, body string) discord.MessageCreate {
	return discord.NewMessageCreateV2(
		discord.NewContainer(
			discord.NewTextDisplay(title),
			discord.NewSmallSeparator(),
			discord.NewTextDisplay(body),
		).WithAccentColor(color),
	).WithAllowedMentions(&discord.AllowedMentions{})
}
