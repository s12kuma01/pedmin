// Copyright (c) 2025-2026 s12kuma01
// SPDX-License-Identifier: MPL-2.0

package handler

import (
	"strings"

	"github.com/disgoorg/disgo/events"
	"github.com/Sumire-Labs/pedmin/internal/model"
	"github.com/Sumire-Labs/pedmin/internal/view"
)

func (h *TicketHandler) HandleComponent(e *events.ComponentInteractionCreate) {
	customID := e.Data.CustomID()
	_, rest, _ := strings.Cut(customID, ":")
	action, extra, _ := strings.Cut(rest, ":")

	guildID := e.GuildID()
	if guildID == nil {
		return
	}

	switch action {
	case "category":
		h.ticketHandleCategorySelect(e, *guildID)
	case "log_prompt":
		h.ticketHandleLogPrompt(e)
	case "log_channel":
		h.ticketHandleLogChannelSelect(e, *guildID)
	case "role_prompt":
		h.ticketHandleRolePrompt(e)
	case "role":
		h.ticketHandleRoleSelect(e, *guildID)
	case "deploy_prompt":
		h.ticketHandleDeployPrompt(e)
	case "deploy_channel":
		h.ticketHandleDeployChannelSelect(e)
	case "deploy_confirm":
		if extra == "" {
			return
		}
		h.ticketHandleDeployConfirm(e, extra)
	case "deploy_cancel":
		_ = e.DeferUpdateMessage()
	case "create":
		if !h.bot.IsModuleEnabled(*guildID, model.TicketModuleID) {
			return
		}
		_ = e.Modal(view.TicketCreateModal())
	case "close":
		h.ticketArchiveTicket(e, *guildID)
	case "delete":
		h.ticketDeleteTicket(e, *guildID)
	}
}
