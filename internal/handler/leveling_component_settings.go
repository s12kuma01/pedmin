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

func (h *LevelingHandler) handleSettingsTab(e *events.ComponentInteractionCreate, tab string) {
	switch tab {
	case "general":
		h.refreshGeneralTab(e, *e.GuildID())
	case "rewards":
		rewards, err := h.service.GetRoleRewards(*e.GuildID())
		if err != nil {
			h.logger.Error("failed to get role rewards", slog.Any("error", err))
		}
		_ = e.UpdateMessage(view.LevelingRewardsTab(rewards))
	case "multipliers":
		settings, err := h.service.LoadSettings(*e.GuildID())
		if err != nil {
			h.logger.Error("failed to load settings", slog.Any("error", err))
			settings = model.DefaultLevelingSettings()
		}
		_ = e.UpdateMessage(view.LevelingMultipliersTab(settings))
	}
}

func (h *LevelingHandler) refreshGeneralTab(e *events.ComponentInteractionCreate, guildID snowflake.ID) {
	settings, err := h.service.LoadSettings(guildID)
	if err != nil {
		h.logger.Error("failed to load settings", slog.Any("error", err))
		settings = model.DefaultLevelingSettings()
	}
	rewardCount, _ := h.service.CountRoleRewards(guildID)
	settingsUI := view.LevelingSettingsPanel(settings, rewardCount)
	enabled := h.bot.IsModuleEnabled(guildID, model.LevelingModuleID)
	_ = e.UpdateMessage(ui.BuildModulePanel(h.Info(), enabled, settingsUI))
}

func (h *LevelingHandler) handleNotificationMode(e *events.ComponentInteractionCreate) {
	data, ok := e.Data.(discord.StringSelectMenuInteractionData)
	if !ok || len(data.Values) == 0 {
		return
	}

	settings, err := h.service.LoadSettings(*e.GuildID())
	if err != nil {
		h.logger.Error("failed to load settings", slog.Any("error", err))
		return
	}

	settings.NotificationMode = data.Values[0]
	if err := h.service.SaveSettings(*e.GuildID(), settings); err != nil {
		h.logger.Error("failed to save settings", slog.Any("error", err))
	}

	h.refreshGeneralTab(e, *e.GuildID())
}

func (h *LevelingHandler) handleNotificationChannel(e *events.ComponentInteractionCreate) {
	data, ok := e.Data.(discord.ChannelSelectMenuInteractionData)
	if !ok || len(data.Values) == 0 {
		return
	}

	settings, err := h.service.LoadSettings(*e.GuildID())
	if err != nil {
		h.logger.Error("failed to load settings", slog.Any("error", err))
		return
	}

	settings.NotificationChID = data.Values[0]
	if err := h.service.SaveSettings(*e.GuildID(), settings); err != nil {
		h.logger.Error("failed to save settings", slog.Any("error", err))
	}

	h.refreshGeneralTab(e, *e.GuildID())
}

func (h *LevelingHandler) handleXPRangePrompt(e *events.ComponentInteractionCreate) {
	_ = e.Modal(discord.ModalCreate{
		CustomID: model.LevelingModuleID + ":xp_range_modal",
		Title:    "XP範囲設定",
		Components: []discord.LayoutComponent{
			discord.NewLabel("最小XP",
				discord.NewShortTextInput(model.LevelingModuleID+":min_xp").
					WithRequired(true).WithPlaceholder("15"),
			),
			discord.NewLabel("最大XP",
				discord.NewShortTextInput(model.LevelingModuleID+":max_xp").
					WithRequired(true).WithPlaceholder("25"),
			),
		},
	})
}

func (h *LevelingHandler) handleCooldownPrompt(e *events.ComponentInteractionCreate) {
	_ = e.Modal(discord.ModalCreate{
		CustomID: model.LevelingModuleID + ":cooldown_modal",
		Title:    "クールダウン設定",
		Components: []discord.LayoutComponent{
			discord.NewLabel("クールダウン（秒）",
				discord.NewShortTextInput(model.LevelingModuleID+":cooldown").
					WithRequired(true).WithPlaceholder("60"),
			),
		},
	})
}

func (h *LevelingHandler) handleVoiceXPPrompt(e *events.ComponentInteractionCreate) {
	_ = e.Modal(discord.ModalCreate{
		CustomID: model.LevelingModuleID + ":voice_xp_modal",
		Title:    "ボイスXP設定",
		Components: []discord.LayoutComponent{
			discord.NewLabel("XP/分",
				discord.NewShortTextInput(model.LevelingModuleID+":voice_xp_val").
					WithRequired(true).WithPlaceholder("5"),
			),
		},
	})
}
