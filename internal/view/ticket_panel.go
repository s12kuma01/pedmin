// Copyright (c) 2025-2026 s12kuma01
// SPDX-License-Identifier: MPL-2.0

package view

import (
	"github.com/disgoorg/disgo/discord"
	"github.com/Sumire-Labs/pedmin/internal/model"
)

// TicketPanel builds the ticket creation panel message.
func TicketPanel() discord.MessageCreate {
	return discord.NewMessageCreateV2(
		discord.NewContainer(
			discord.NewTextDisplay("## 🎫 チケットサポート"),
			discord.NewSmallSeparator(),
			discord.NewTextDisplay("サポートが必要な場合は、下のボタンからチケットを作成してください。"),
			discord.NewActionRow(
				discord.NewPrimaryButton("チケットを作成", model.TicketModuleID+":create"),
			),
		),
	).WithAllowedMentions(&discord.AllowedMentions{})
}
