// Copyright (c) 2025-2026 s12kuma01
// SPDX-License-Identifier: MPL-2.0

package view

import (
	"fmt"
	"math"

	"github.com/disgoorg/disgo/discord"
	"github.com/Sumire-Labs/pedmin/internal/model"
)

// TicketLog builds the ticket log message for the log channel.
func TicketLog(ticket *model.Ticket) discord.MessageCreate {
	createdAt := ticket.CreatedAt.Format("2006-01-02 15:04")

	closedAtText := "-"
	closedByText := "-"
	durationText := "-"

	if ticket.ClosedAt != nil {
		closedAtText = ticket.ClosedAt.Format("2006-01-02 15:04")

		dur := ticket.ClosedAt.Sub(ticket.CreatedAt)
		minutes := int(math.Round(dur.Minutes()))
		if minutes < 60 {
			durationText = fmt.Sprintf("%d分", minutes)
		} else {
			h := minutes / 60
			m := minutes % 60
			if m == 0 {
				durationText = fmt.Sprintf("%d時間", h)
			} else {
				durationText = fmt.Sprintf("%d時間%d分", h, m)
			}
		}
	}
	if ticket.ClosedBy != nil {
		closedByText = fmt.Sprintf("<@%d>", *ticket.ClosedBy)
	}

	body := fmt.Sprintf(
		"**チケット:** #%04d\n**作成者:** <@%d>\n**件名:** %s\n**作成日時:** %s\n**クローズ日時:** %s\n**対応時間:** %s\n**クローズ者:** %s",
		ticket.Number, ticket.UserID, ticket.Subject,
		createdAt, closedAtText, durationText, closedByText,
	)

	return discord.NewMessageCreateV2(
		discord.NewContainer(
			discord.NewTextDisplay("### 📋 チケットログ"),
			discord.NewSmallSeparator(),
			discord.NewTextDisplay(body),
		),
	).WithAllowedMentions(&discord.AllowedMentions{})
}
