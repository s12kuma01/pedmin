// Copyright (c) 2025-2026 s12kuma01
// SPDX-License-Identifier: MPL-2.0

package handler

import (
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
	"github.com/disgoorg/snowflake/v2"
	"github.com/Sumire-Labs/pedmin/internal/model"
)

func (h *PlayerHandler) handleAddModal(e *events.ComponentInteractionCreate) {
	_ = e.Modal(discord.ModalCreate{
		CustomID: model.PlayerModuleID + ":add_modal",
		Title:    "キューに追加",
		Components: []discord.LayoutComponent{
			discord.NewLabel("検索キーワードまたはURL",
				discord.NewShortTextInput(model.PlayerModuleID+":query").
					WithPlaceholder("曲名またはYouTube/SpotifyのURL").
					WithRequired(true),
			),
		},
	})
}

func (h *PlayerHandler) handleShowQueue(e *events.ComponentInteractionCreate, guildID snowflake.ID) {
	queueUI := h.service.BuildQueueUI(guildID)
	_ = e.UpdateMessage(discord.NewMessageUpdateV2([]discord.LayoutComponent{queueUI}))
}

func (h *PlayerHandler) handleBack(e *events.ComponentInteractionCreate, guildID snowflake.ID) {
	playerUI := h.service.BuildPlayerUI(guildID)
	_ = e.UpdateMessage(discord.NewMessageUpdateV2([]discord.LayoutComponent{playerUI}))
}

func (h *PlayerHandler) handleClearQueue(e *events.ComponentInteractionCreate, guildID snowflake.ID) {
	h.service.ClearQueue(guildID)

	queueUI := h.service.BuildQueueUI(guildID)
	_ = e.UpdateMessage(discord.NewMessageUpdateV2([]discord.LayoutComponent{queueUI}))
}
