package ticket

import (
	"fmt"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/snowflake/v2"
)

func BuildCreateTicketModal() discord.ModalCreate {
	return discord.ModalCreate{
		CustomID: ModuleID + ":create_modal",
		Title:    "チケットを作成",
		Components: []discord.LayoutComponent{
			discord.NewLabel("件名",
				discord.NewShortTextInput(ModuleID+":subject").
					WithPlaceholder("チケットの件名を入力").
					WithRequired(true).
					WithMaxLength(100),
			),
			discord.NewLabel("説明",
				discord.NewParagraphTextInput(ModuleID+":description").
					WithPlaceholder("詳しい内容を入力してください").
					WithRequired(false).
					WithMaxLength(1000),
			),
		},
	}
}

func BuildTicketInfo(number int, userID snowflake.ID, subject, description string) discord.MessageCreate {
	body := fmt.Sprintf("**作成者:** <@%d>\n**件名:** %s", userID, subject)
	if description != "" {
		body += "\n\n" + description
	}

	return discord.NewMessageCreateV2(
		discord.NewContainer(
			discord.NewTextDisplay(fmt.Sprintf("## 🎫 チケット #%04d", number)),
			discord.NewSmallSeparator(),
			discord.NewTextDisplay(body),
			discord.NewLargeSeparator(),
			discord.NewActionRow(
				discord.NewDangerButton("チケットを閉じる", ModuleID+":close"),
			),
		),
	).WithAllowedMentions(&discord.AllowedMentions{})
}

func BuildArchiveInfo(number int, userID snowflake.ID, subject string) discord.MessageCreate {
	body := fmt.Sprintf("**作成者:** <@%d>\n**件名:** %s\n\nこのチケットはクローズされました。", userID, subject)

	return discord.NewMessageCreateV2(
		discord.NewContainer(
			discord.NewTextDisplay(fmt.Sprintf("## 🔒 チケット #%04d (アーカイブ)", number)),
			discord.NewSmallSeparator(),
			discord.NewTextDisplay(body),
			discord.NewLargeSeparator(),
			discord.NewActionRow(
				discord.NewDangerButton("チケットを削除", ModuleID+":delete"),
			),
		),
	).WithAllowedMentions(&discord.AllowedMentions{})
}
