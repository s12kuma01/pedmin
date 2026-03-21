package logger

import (
	"fmt"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/snowflake/v2"
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
		discord.NewContainer(components...),
	).WithAllowedMentions(&discord.AllowedMentions{})
}

func BuildMessageDeleteLog(user *discord.User, channelID snowflake.ID, content string, attachments []discord.Attachment, forwarded bool) discord.MessageCreate {
	userText := "*不明*"
	if user != nil {
		userText = fmt.Sprintf("<@%d>", user.ID)
	}

	title := "### 🗑️ メッセージ削除"
	hasContent := content != ""
	hasAttachments := len(attachments) > 0

	var body string
	switch {
	case hasContent && forwarded:
		body = fmt.Sprintf("**ユーザー:** %s\n**チャンネル:** <#%d>\n**転送メッセージの内容:**\n> %s",
			userText, channelID, content)
	case hasContent:
		body = fmt.Sprintf("**ユーザー:** %s\n**チャンネル:** <#%d>\n**内容:**\n> %s",
			userText, channelID, content)
	case !hasContent && !hasAttachments && user == nil:
		body = fmt.Sprintf("**ユーザー:** %s\n**チャンネル:** <#%d>\n**内容:**\n> *内容を取得できませんでした*",
			userText, channelID)
	default:
		body = fmt.Sprintf("**ユーザー:** %s\n**チャンネル:** <#%d>",
			userText, channelID)
	}

	components := []discord.ContainerSubComponent{
		discord.NewTextDisplay(title),
		discord.NewSmallSeparator(),
		discord.NewTextDisplay(body),
	}

	if hasAttachments {
		label := "**添付ファイル:**"
		if forwarded {
			label = "**転送メッセージの添付ファイル:**"
		}
		components = append(components, discord.NewSmallSeparator())
		components = append(components, discord.NewTextDisplay(label))
		components = append(components, buildAttachmentComponents(attachments)...)
	}

	return discord.NewMessageCreateV2(
		discord.NewContainer(components...),
	).WithAllowedMentions(&discord.AllowedMentions{})
}
