// Copyright (c) 2025-2026 s12kuma01
// SPDX-License-Identifier: MPL-2.0

package handler

import (
	"log/slog"
	"strings"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
	"github.com/disgoorg/snowflake/v2"
	"github.com/Sumire-Labs/pedmin/internal/model"
	"github.com/Sumire-Labs/pedmin/internal/ui"
	"github.com/Sumire-Labs/pedmin/internal/view"
)

func (h *LoggerHandler) HandleComponent(e *events.ComponentInteractionCreate) {
	customID := e.Data.CustomID()
	_, action, _ := strings.Cut(customID, ":")

	guildID := e.GuildID()
	if guildID == nil {
		return
	}

	settings, err := h.service.LoadSettings(*guildID)
	if err != nil {
		h.logger.Error("failed to load logger settings", slog.Any("error", err))
		return
	}

	switch action {
	case "channel":
		data, ok := e.Data.(discord.ChannelSelectMenuInteractionData)
		if !ok {
			return
		}
		if len(data.Values) > 0 {
			settings.ChannelID = data.Values[0]
		}

	case "events":
		data, ok := e.Data.(discord.StringSelectMenuInteractionData)
		if !ok {
			return
		}
		for k := range settings.Events {
			settings.Events[k] = false
		}
		for _, v := range data.Values {
			settings.Events[v] = true
		}

	default:
		return
	}

	if err := h.service.SaveSettings(*guildID, settings); err != nil {
		h.logger.Error("failed to save logger settings", slog.Any("error", err))
	}

	h.loggerRefreshSettingsPanel(e, *guildID, settings)
}

func (h *LoggerHandler) loggerRefreshSettingsPanel(e *events.ComponentInteractionCreate, guildID snowflake.ID, settings *model.LoggerSettings) {
	settingsUI := view.LoggerSettingsPanel(settings)
	enabled := h.bot.IsModuleEnabled(guildID, model.LoggerModuleID)
	_ = e.UpdateMessage(ui.BuildModulePanel(h.Info(), enabled, settingsUI))
}
