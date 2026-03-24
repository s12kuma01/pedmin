// Copyright (c) 2025-2026 s12kuma01
// SPDX-License-Identifier: MPL-2.0

package handler

import (
	"log/slog"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
	"github.com/disgoorg/snowflake/v2"
	"github.com/s12kuma01/pedmin/internal/model"
	"github.com/s12kuma01/pedmin/internal/view"
)

func (h *TicketHandler) ticketHandleDeployPrompt(e *events.ComponentInteractionCreate) {
	_ = e.CreateMessage(discord.NewMessageCreateV2(
		discord.NewContainer(
			discord.NewTextDisplay("パネルを設置するチャンネルを選択してください:"),
			discord.NewActionRow(
				discord.NewChannelSelectMenu(model.TicketModuleID+":deploy_channel", "チャンネルを選択...").
					WithChannelTypes(discord.ChannelTypeGuildText),
			),
		),
	).WithEphemeral(true))
}

func (h *TicketHandler) ticketHandleDeployChannelSelect(e *events.ComponentInteractionCreate) {
	data, ok := e.Data.(discord.ChannelSelectMenuInteractionData)
	if !ok || len(data.Values) == 0 {
		return
	}
	channelID := data.Values[0]
	_ = e.UpdateMessage(view.TicketDeployConfirm(channelID))
}

func (h *TicketHandler) ticketHandleDeployConfirm(e *events.ComponentInteractionCreate, channelIDStr string) {
	_ = e.DeferUpdateMessage()

	channelID, err := snowflake.Parse(channelIDStr)
	if err != nil {
		h.logger.Error("failed to parse channel ID", slog.Any("error", err))
		return
	}

	if err := h.service.DeployPanel(channelID); err != nil {
		h.logger.Error("failed to deploy ticket panel", slog.Any("error", err))
	}
}
