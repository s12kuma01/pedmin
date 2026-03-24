// Copyright (c) 2025-2026 s12kuma01
// SPDX-License-Identifier: MPL-2.0

package handler

import (
	"github.com/disgoorg/disgo/events"
	"github.com/disgoorg/snowflake/v2"
	"github.com/s12kuma01/pedmin/internal/model"
)

func (h *LevelingHandler) onMessageCreate(e *events.GuildMessageCreate) {
	if e.Message.Author.Bot {
		return
	}
	if !h.bot.IsModuleEnabled(e.GuildID, model.LevelingModuleID) {
		return
	}

	var roles []snowflake.ID
	if e.Message.Member != nil {
		roles = e.Message.Member.RoleIDs
	}

	h.service.ProcessMessage(e.GuildID, e.Message.Author.ID, e.ChannelID, roles)
}
