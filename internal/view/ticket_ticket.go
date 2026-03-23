// opyright (c) 2025-2026 s12kuma01
// SPDX-License-Identifier: MPL-2.0

package view

import (
	"fmt"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/snowflake/v2"
	"github.com/s12kuma01/pedmin/internal/model"
)

// TicketCreateModal builds the ticket creation modal.
func TicketCreateModal() discord.ModalCreate {
	return discord.ModalCreate{
		CustomID: model.TicketModuleID + ":create_modal",
		Title:    "チケットを作成",
		Components: []discord.LayoutComponent{
			discord.NewLabel("件名",
				discord.NewShortTextInput(model.TicketModuleID+":subject").
					WithPlaceholder("チケットの件名を入力").
					WithRequired(true).
					WithMaxLength(100),
			),
			discord.NewLabel("説明",
				discord.NewParagraphTextInput(model.TicketModuleID+":description").
					WithPlaceholder("詳しい内容を入力してください").
					WithRequired(false).
					WithMaxLength(1000),
			),
		},
	}
}

// TicketInfo builds the ticket channel info message.
func TicketInfo(number int, userID snowflake.ID, subject, description string) discord.MessageCreate {
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
				discord.NewDangerButton("チケットを閉じる", model.TicketModuleID+":close"),
			),
		),
	).WithAllowedMentions(&discord.AllowedMentions{})
}

// TicketArchiveInfo builds the archived ticket info message.
func TicketArchiveInfo(number int, userID snowflake.ID, subject string) discord.MessageCreate {
	body := fmt.Sprintf("**作成者:** <@%d>\n**件名:** %s\n\nこのチケットはクローズされました。", userID, subject)

	return discord.NewMessageCreateV2(
		discord.NewContainer(
			discord.NewTextDisplay(fmt.Sprintf("## 🔒 チケット #%04d (アーカイブ)", number)),
			discord.NewSmallSeparator(),
			discord.NewTextDisplay(body),
			discord.NewLargeSeparator(),
			discord.NewActionRow(
				discord.NewDangerButton("チケットを削除", model.TicketModuleID+":delete"),
			),
		),
	).WithAllowedMentions(&discord.AllowedMentions{})
}
