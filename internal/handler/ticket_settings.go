// Copyright (c) 2025-2026 s12kuma01
// SPDX-License-Identifier: MPL-2.0

package handler

import (
	"log/slog"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
	"github.com/disgoorg/snowflake/v2"
	"github.com/Sumire-Labs/pedmin/internal/model"
	"github.com/Sumire-Labs/pedmin/internal/ui"
	"github.com/Sumire-Labs/pedmin/internal/view"
)

func (h *TicketHandler) ticketHandleCategorySelect(e *events.ComponentInteractionCreate, guildID snowflake.ID) {
	data, ok := e.Data.(discord.ChannelSelectMenuInteractionData)
	if !ok || len(data.Values) == 0 {
		return
	}
	if err := h.service.UpdateCategory(guildID, data.Values[0]); err != nil {
		h.logger.Error("failed to update category", slog.Any("error", err))
	}
	h.ticketRefreshSettingsPanel(e, guildID)
}

func (h *TicketHandler) ticketHandleLogPrompt(e *events.ComponentInteractionCreate) {
	_ = e.CreateMessage(discord.NewMessageCreateV2(
		discord.NewContainer(
			discord.NewTextDisplay("ログチャンネルを選択してください:"),
			discord.NewActionRow(
				discord.NewChannelSelectMenu(model.TicketModuleID+":log_channel", "ログチャンネルを選択...").
					WithChannelTypes(discord.ChannelTypeGuildText),
			),
		),
	).WithEphemeral(true))
}

func (h *TicketHandler) ticketHandleLogChannelSelect(e *events.ComponentInteractionCreate, guildID snowflake.ID) {
	data, ok := e.Data.(discord.ChannelSelectMenuInteractionData)
	if !ok || len(data.Values) == 0 {
		return
	}
	if err := h.service.UpdateLogChannel(guildID, data.Values[0]); err != nil {
		h.logger.Error("failed to update log channel", slog.Any("error", err))
	}
	h.ticketRefreshSettingsPanel(e, guildID)
}

func (h *TicketHandler) ticketHandleRolePrompt(e *events.ComponentInteractionCreate) {
	_ = e.CreateMessage(discord.NewMessageCreateV2(
		discord.NewContainer(
			discord.NewTextDisplay("サポートロールを選択してください:"),
			discord.NewActionRow(
				discord.NewRoleSelectMenu(model.TicketModuleID+":role", "サポートロールを選択..."),
			),
		),
	).WithEphemeral(true))
}

func (h *TicketHandler) ticketHandleRoleSelect(e *events.ComponentInteractionCreate, guildID snowflake.ID) {
	data, ok := e.Data.(discord.RoleSelectMenuInteractionData)
	if !ok || len(data.Values) == 0 {
		return
	}
	if err := h.service.UpdateSupportRole(guildID, data.Values[0]); err != nil {
		h.logger.Error("failed to update support role", slog.Any("error", err))
	}
	h.ticketRefreshSettingsPanel(e, guildID)
}

func (h *TicketHandler) ticketRefreshSettingsPanel(e *events.ComponentInteractionCreate, guildID snowflake.ID) {
	settings, err := h.service.LoadSettings(guildID)
	if err != nil {
		h.logger.Error("failed to load ticket settings for refresh", slog.Any("error", err))
		_ = e.DeferUpdateMessage()
		return
	}
	settingsUI := view.TicketSettingsPanel(settings)
	enabled := h.bot.IsModuleEnabled(guildID, model.TicketModuleID)
	_ = e.UpdateMessage(ui.BuildModulePanel(h.Info(), enabled, settingsUI))
}

func (h *TicketHandler) ticketArchiveTicket(e *events.ComponentInteractionCreate, guildID snowflake.ID) {
	channelID := e.Channel().ID()

	_ = e.DeferUpdateMessage()

	ticket, err := h.service.ArchiveTicket(guildID, channelID, e.User().ID)
	if err != nil {
		h.logger.Error("failed to archive ticket", slog.Any("error", err))
		return
	}

	archiveUI := view.TicketArchiveInfo(ticket.Number, ticket.UserID, ticket.Subject)
	_, _ = e.Client().Rest.UpdateInteractionResponse(e.ApplicationID(), e.Token(), discord.NewMessageUpdateV2(archiveUI.Components))
}

func (h *TicketHandler) ticketDeleteTicket(e *events.ComponentInteractionCreate, guildID snowflake.ID) {
	channelID := e.Channel().ID()

	_ = e.DeferUpdateMessage()

	if _, err := h.service.DeleteTicket(guildID, channelID); err != nil {
		h.logger.Error("failed to delete ticket", slog.Any("error", err))
	}
}
