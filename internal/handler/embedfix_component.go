// Copyright (c) 2025-2026 s12kuma01
// SPDX-License-Identifier: MPL-2.0

package handler

import (
	"context"
	"log/slog"
	"strings"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
	"github.com/Sumire-Labs/pedmin/internal/model"
	"github.com/Sumire-Labs/pedmin/internal/ui"
	"github.com/Sumire-Labs/pedmin/internal/view"
)

func (h *EmbedFixHandler) handleComponent(e *events.ComponentInteractionCreate) {
	customID := e.Data.CustomID()
	_, rest, _ := strings.Cut(customID, ":")
	action, rest, _ := strings.Cut(rest, ":")

	switch action {
	case "translate":
		h.handleTranslate(e, rest)
	case "platforms":
		h.handlePlatformSettings(e)
	}
}

func (h *EmbedFixHandler) handleTranslate(e *events.ComponentInteractionCreate, rest string) {
	_ = e.DeferUpdateMessage()

	if !h.service.IsTranslationAvailable() {
		h.respondTranslateError(e, "翻訳APIキーが設定されていないため、翻訳できません。")
		return
	}

	platform, params, _ := strings.Cut(rest, ":")
	ctx := context.Background()

	components, err := h.service.TranslateContent(ctx, platform, params)
	if err != nil {
		h.respondTranslateError(e, "翻訳に失敗しました。")
		return
	}

	_, _ = e.Client().Rest.UpdateInteractionResponse(e.ApplicationID(), e.Token(),
		discord.NewMessageUpdateV2(components))
}

func (h *EmbedFixHandler) handlePlatformSettings(e *events.ComponentInteractionCreate) {
	guildID := e.GuildID()
	if guildID == nil {
		return
	}

	data, ok := e.Data.(discord.StringSelectMenuInteractionData)
	if !ok {
		return
	}

	settings, err := h.service.LoadSettings(*guildID)
	if err != nil {
		h.logger.Error("failed to load embedfix settings", slog.Any("error", err))
		return
	}

	// Disable all, then enable selected
	for k := range settings.Platforms {
		settings.Platforms[k] = false
	}
	for _, v := range data.Values {
		settings.Platforms[model.Platform(v)] = true
	}

	if err := h.service.SaveSettings(*guildID, settings); err != nil {
		h.logger.Error("failed to save embedfix settings", slog.Any("error", err))
	}

	settingsUI := view.BuildEmbedFixSettingsPanel(settings)
	enabled := h.bot.IsModuleEnabled(*guildID, model.EmbedFixModuleID)
	_ = e.UpdateMessage(ui.BuildModulePanel(h.Info(), enabled, settingsUI))
}

func (h *EmbedFixHandler) respondTranslateError(e *events.ComponentInteractionCreate, msg string) {
	_, _ = e.Client().Rest.UpdateInteractionResponse(e.ApplicationID(), e.Token(),
		discord.NewMessageUpdateV2([]discord.LayoutComponent{
			discord.NewContainer(
				discord.NewTextDisplay(msg),
			),
		}))
}
