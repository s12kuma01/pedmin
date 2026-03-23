// opyright (c) 2025-2026 s12kuma01
// SPDX-License-Identifier: MPL-2.0

package handler

import (
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
	"github.com/s12kuma01/pedmin/internal/model"
	"github.com/s12kuma01/pedmin/internal/ui"
)

func (h *PlayerHandler) handleModal(e *events.ModalSubmitInteractionCreate) {
	customID := e.Data.CustomID
	if customID != model.PlayerModuleID+":add_modal" {
		return
	}

	guildID := e.GuildID()
	if guildID == nil {
		return
	}

	ti, ok := e.Data.TextInput(model.PlayerModuleID + ":query")
	query := ""
	if ok {
		query = ti.Value
	}

	if query == "" {
		_ = e.CreateMessage(ui.EphemeralError("検索キーワードまたはURLを入力してください。"))
		return
	}

	_ = e.DeferCreateMessage(true)

	sendFollowup := func(text string) {
		_, _ = e.Client().Rest.UpdateInteractionResponse(
			e.ApplicationID(), e.Token(),
			discord.NewMessageUpdateV2([]discord.LayoutComponent{
				discord.NewContainer(discord.NewTextDisplay(text)),
			}),
		)
	}

	h.service.LoadAndPlay(*guildID, e.Member().User.ID, query, sendFollowup)
}
