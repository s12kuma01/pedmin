// Copyright (c) 2025-2026 s12kuma01
// SPDX-License-Identifier: MPL-2.0

package handler

import (
	"fmt"
	"log/slog"
	"strings"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
	"github.com/Sumire-Labs/pedmin/internal/model"
	"github.com/Sumire-Labs/pedmin/internal/ui"
)

func (h *CounterHandler) HandleModal(e *events.ModalSubmitInteractionCreate) {
	customID := e.Data.CustomID
	_, rest, _ := strings.Cut(customID, ":")
	action, extra, _ := strings.Cut(rest, ":")

	if action != "add_modal" {
		return
	}

	matchType := model.MatchType(extra)

	word := e.Data.Text(model.CounterModuleID + ":word")
	word = strings.TrimSpace(word)

	if word == "" {
		_ = e.CreateMessage(ui.EphemeralError("ワードを入力してください。"))
		return
	}

	guildID := e.GuildID()
	if guildID == nil {
		return
	}

	counter, err := h.service.AddCounter(*guildID, word, matchType)
	if err != nil {
		h.logger.Error("failed to add counter", slog.Any("error", err))
		_ = e.CreateMessage(ui.EphemeralError(fmt.Sprintf("カウンターの追加に失敗しました:\n%s", err.Error())))
		return
	}

	_ = e.CreateMessage(ui.EphemeralV2(
		discord.NewContainer(
			discord.NewTextDisplay(fmt.Sprintf(
				"**%s** (%s) を追加しました。",
				counter.Word, model.MatchTypeLabel(counter.MatchType),
			)),
		),
	))
}
