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
	"github.com/s12kuma01/pedmin/internal/model"
	"github.com/s12kuma01/pedmin/internal/ui"
)

func (h *LevelingHandler) HandleModal(e *events.ModalSubmitInteractionCreate) {
	customID := e.Data.CustomID
	_, rest, _ := strings.Cut(customID, ":")
	action, extra, _ := strings.Cut(rest, ":")

	guildID := e.GuildID()
	if guildID == nil {
		return
	}

	switch action {
	case "xp_range_modal":
		h.handleXPRangeModal(e)
	case "cooldown_modal":
		h.handleCooldownModal(e)
	case "voice_xp_modal":
		h.handleVoiceXPModal(e)
	case "reward_level_modal":
		h.handleRewardLevelModal(e)
	case "role_mult_modal":
		h.handleRoleMultModal(e, extra)
	case "ch_mult_modal":
		h.handleChMultModal(e, extra)
	}
}

func (h *LevelingHandler) handleXPRangeModal(e *events.ModalSubmitInteractionCreate) {
	minStr := strings.TrimSpace(e.Data.Text(model.LevelingModuleID + ":min_xp"))
	maxStr := strings.TrimSpace(e.Data.Text(model.LevelingModuleID + ":max_xp"))

	minXP, err := strconv.Atoi(minStr)
	if err != nil || minXP < 1 || minXP > 1000 {
		_ = e.CreateMessage(ui.EphemeralError("最小XPは1〜1000の整数を入力してください。"))
		return
	}

	maxXP, err := strconv.Atoi(maxStr)
	if err != nil || maxXP < 1 || maxXP > 1000 {
		_ = e.CreateMessage(ui.EphemeralError("最大XPは1〜1000の整数を入力してください。"))
		return
	}

	if minXP > maxXP {
		_ = e.CreateMessage(ui.EphemeralError("最小XPは最大XP以下にしてください。"))
		return
	}

	settings, err := h.service.LoadSettings(*e.GuildID())
	if err != nil {
		h.logger.Error("failed to load settings", slog.Any("error", err))
		return
	}

	settings.MinXP = minXP
	settings.MaxXP = maxXP
	if err := h.service.SaveSettings(*e.GuildID(), settings); err != nil {
		h.logger.Error("failed to save settings", slog.Any("error", err))
		_ = e.CreateMessage(ui.EphemeralError("設定の保存に失敗しました。"))
		return
	}

	_ = e.CreateMessage(ui.EphemeralV2(
		discord.NewContainer(
			discord.NewTextDisplay(fmt.Sprintf("XP範囲を **%d〜%d** に設定しました。", minXP, maxXP)),
		),
	))
}

func (h *LevelingHandler) handleCooldownModal(e *events.ModalSubmitInteractionCreate) {
	cdStr := strings.TrimSpace(e.Data.Text(model.LevelingModuleID + ":cooldown"))

	cd, err := strconv.Atoi(cdStr)
	if err != nil || cd < 1 || cd > 3600 {
		_ = e.CreateMessage(ui.EphemeralError("クールダウンは1〜3600秒の整数を入力してください。"))
		return
	}

	settings, err := h.service.LoadSettings(*e.GuildID())
	if err != nil {
		h.logger.Error("failed to load settings", slog.Any("error", err))
		return
	}

	settings.CooldownSeconds = cd
	if err := h.service.SaveSettings(*e.GuildID(), settings); err != nil {
		h.logger.Error("failed to save settings", slog.Any("error", err))
		_ = e.CreateMessage(ui.EphemeralError("設定の保存に失敗しました。"))
		return
	}

	_ = e.CreateMessage(ui.EphemeralV2(
		discord.NewContainer(
			discord.NewTextDisplay(fmt.Sprintf("クールダウンを **%d秒** に設定しました。", cd)),
		),
	))
}

func (h *LevelingHandler) handleVoiceXPModal(e *events.ModalSubmitInteractionCreate) {
	xpStr := strings.TrimSpace(e.Data.Text(model.LevelingModuleID + ":voice_xp_val"))

	xp, err := strconv.Atoi(xpStr)
	if err != nil || xp < 0 || xp > 100 {
		_ = e.CreateMessage(ui.EphemeralError("ボイスXPは0〜100の整数を入力してください。"))
		return
	}

	settings, err := h.service.LoadSettings(*e.GuildID())
	if err != nil {
		h.logger.Error("failed to load settings", slog.Any("error", err))
		return
	}

	settings.VoiceXPPerMinute = xp
	if err := h.service.SaveSettings(*e.GuildID(), settings); err != nil {
		h.logger.Error("failed to save settings", slog.Any("error", err))
		_ = e.CreateMessage(ui.EphemeralError("設定の保存に失敗しました。"))
		return
	}

	_ = e.CreateMessage(ui.EphemeralV2(
		discord.NewContainer(
			discord.NewTextDisplay(fmt.Sprintf("ボイスXPを **%d/分** に設定しました。", xp)),
		),
	))
}

func (h *LevelingHandler) handleRewardLevelModal(e *events.ModalSubmitInteractionCreate) {
	levelStr := strings.TrimSpace(e.Data.Text(model.LevelingModuleID + ":reward_level"))

	level, err := strconv.Atoi(levelStr)
	if err != nil || level < 1 || level > model.MaxLevel {
		_ = e.CreateMessage(ui.EphemeralError(fmt.Sprintf("レベルは1〜%dの整数を入力してください。", model.MaxLevel)))
		return
	}

	// Show role select menu
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
