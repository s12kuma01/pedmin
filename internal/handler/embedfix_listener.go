// Copyright (c) 2025-2026 s12kuma01
// SPDX-License-Identifier: MPL-2.0

package handler

import (
	"context"

	"github.com/disgoorg/disgo/events"
	"github.com/Sumire-Labs/pedmin/internal/model"
)

func (h *EmbedFixHandler) onMessageCreate(e *events.GuildMessageCreate) {
	if e.Message.Author.Bot {
		return
	}
	if !h.bot.IsModuleEnabled(e.GuildID, model.EmbedFixModuleID) {
		return
	}

	h.service.ProcessMessageURLs(context.Background(), e.GuildID, e.ChannelID, e.MessageID, e.Message.Content)
}
