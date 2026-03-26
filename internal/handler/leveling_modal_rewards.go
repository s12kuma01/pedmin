// Copyright (c) 2025-2026 s12kuma01
// SPDX-License-Identifier: MPL-2.0

package handler

import (
	"fmt"
	"log/slog"
	"strconv"
	"strings"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
	"github.com/disgoorg/snowflake/v2"
	"github.com/Sumire-Labs/pedmin/internal/model"
	"github.com/Sumire-Labs/pedmin/internal/ui"
)

func (h *LevelingHandler) handleRewardLevelModal(e *events.ModalSubmitInteractionCreate) {
	levelStr := strings.TrimSpace(e.Data.Text(model.LevelingModuleID + ":reward_level"))

	level, err := strconv.Atoi(levelStr)
	if err != nil || level < 1 || level > model.MaxLevel {
		_ = e.CreateMessage(ui.EphemeralError(fmt.Sprintf("レベルは1〜%dの整数を入力してください。", model.MaxLevel)))
		return
	}

	_ = e.CreateMessage(ui.EphemeralV2(
		discord.NewContainer(
			discord.NewTextDisplay(fmt.Sprintf("**Lv.%d** に付与するロールを選択してください:", level)),
			discord.NewActionRow(
				discord.NewRoleSelectMenu(
					fmt.Sprintf("%s:reward_add_role:%d", model.LevelingModuleID, level),
					"ロールを選択...",
				),
			),
		),
	))
}

func (h *LevelingHandler) handleRoleMultModal(e *events.ModalSubmitInteractionCreate, roleIDStr string) {
	multStr := strings.TrimSpace(e.Data.Text(model.LevelingModuleID + ":mult_value"))

	mult, err := strconv.ParseFloat(multStr, 64)
	if err != nil || mult < 0.1 || mult > 10.0 {
		_ = e.CreateMessage(ui.EphemeralError("倍率は0.1〜10.0の数値を入力してください。"))
		return
	}

	roleID, err := strconv.ParseUint(roleIDStr, 10, 64)
	if err != nil {
		return
	}

	settings, err := h.service.LoadSettings(*e.GuildID())
	if err != nil {
		h.logger.Error("failed to load settings", slog.Any("error", err))
		return
	}

	totalMults := len(settings.RoleMultipliers) + len(settings.ChannelMultipliers)
	if totalMults >= model.MaxMultipliersPerGuild {
		_ = e.CreateMessage(ui.EphemeralError(fmt.Sprintf("倍率数が上限(%d)に達しています。", model.MaxMultipliersPerGuild)))
		return
	}

	settings.RoleMultipliers = append(settings.RoleMultipliers, model.RoleMultiplier{
		RoleID:     snowflake.ID(roleID),
		Multiplier: mult,
	})
	if err := h.service.SaveSettings(*e.GuildID(), settings); err != nil {
		h.logger.Error("failed to save settings", slog.Any("error", err))
		_ = e.CreateMessage(ui.EphemeralError("設定の保存に失敗しました。"))
		return
	}

	_ = e.CreateMessage(ui.EphemeralV2(
		discord.NewContainer(
			discord.NewTextDisplay(fmt.Sprintf("<@&%d> に **%.1fx** 倍率を設定しました。", roleID, mult)),
		),
	))
}

func (h *LevelingHandler) handleChMultModal(e *events.ModalSubmitInteractionCreate, chIDStr string) {
	multStr := strings.TrimSpace(e.Data.Text(model.LevelingModuleID + ":mult_value"))

	mult, err := strconv.ParseFloat(multStr, 64)
	if err != nil || mult < 0.1 || mult > 10.0 {
		_ = e.CreateMessage(ui.EphemeralError("倍率は0.1〜10.0の数値を入力してください。"))
		return
	}

	chID, err := strconv.ParseUint(chIDStr, 10, 64)
	if err != nil {
		return
	}

	settings, err := h.service.LoadSettings(*e.GuildID())
	if err != nil {
		h.logger.Error("failed to load settings", slog.Any("error", err))
		return
	}

	totalMults := len(settings.RoleMultipliers) + len(settings.ChannelMultipliers)
	if totalMults >= model.MaxMultipliersPerGuild {
		_ = e.CreateMessage(ui.EphemeralError(fmt.Sprintf("倍率数が上限(%d)に達しています。", model.MaxMultipliersPerGuild)))
		return
	}

	settings.ChannelMultipliers = append(settings.ChannelMultipliers, model.ChannelMultiplier{
		ChannelID:  snowflake.ID(chID),
		Multiplier: mult,
	})
	if err := h.service.SaveSettings(*e.GuildID(), settings); err != nil {
		h.logger.Error("failed to save settings", slog.Any("error", err))
		_ = e.CreateMessage(ui.EphemeralError("設定の保存に失敗しました。"))
		return
	}

	_ = e.CreateMessage(ui.EphemeralV2(
		discord.NewContainer(
			discord.NewTextDisplay(fmt.Sprintf("<#%d> に **%.1fx** 倍率を設定しました。", chID, mult)),
		),
	))
}
