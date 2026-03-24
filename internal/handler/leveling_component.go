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
	"github.com/s12kuma01/pedmin/internal/view"
)

func (h *LevelingHandler) HandleComponent(e *events.ComponentInteractionCreate) {
	customID := e.Data.CustomID()
	_, rest, _ := strings.Cut(customID, ":")
	action, extra, _ := strings.Cut(rest, ":")

	guildID := e.GuildID()
	if guildID == nil {
		return
	}

	switch action {
	// --- Settings tabs ---
	case "settings_tab":
		h.handleSettingsTab(e, extra)

	// --- General settings ---
	case "notification_mode":
		h.handleNotificationMode(e)
	case "notification_ch":
		h.handleNotificationChannel(e)
	case "xp_range":
		h.handleXPRangePrompt(e)
	case "cooldown":
		h.handleCooldownPrompt(e)
	case "voice_xp":
		h.handleVoiceXPPrompt(e)

	// --- Rewards ---
	case "reward_add":
		h.handleRewardAddPrompt(e)
	case "reward_add_role":
		h.handleRewardAddRole(e, extra)
	case "reward_manage":
		h.handleRewardManage(e)
	case "reward_manage_select":
		h.handleRewardManageSelect(e)
	case "reward_remove":
		h.handleRewardRemove(e, extra)

	// --- Multipliers ---
	case "role_mult_add":
		h.handleRoleMultAddPrompt(e)
	case "role_mult_add_role":
		h.handleRoleMultAddRole(e)
	case "role_mult_remove":
		h.handleRoleMultRemove(e, extra)
	case "ch_mult_add":
		h.handleChMultAddPrompt(e)
	case "ch_mult_add_ch":
		h.handleChMultAddChannel(e)
	case "ch_mult_remove":
		h.handleChMultRemove(e, extra)

	// --- Leaderboard ---
	case "lb_page":
		h.handleLeaderboardPage(e, extra)
	}
}

// --- Settings Tabs ---

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

// --- General Settings ---

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

// --- Rewards ---

func (h *LevelingHandler) handleRewardAddPrompt(e *events.ComponentInteractionCreate) {
	_ = e.Modal(discord.ModalCreate{
		CustomID: model.LevelingModuleID + ":reward_level_modal",
		Title:    "ロール報酬追加",
		Components: []discord.LayoutComponent{
			discord.NewLabel("レベル (1-100)",
				discord.NewShortTextInput(model.LevelingModuleID+":reward_level").
					WithRequired(true).WithPlaceholder("10"),
			),
		},
	})
}

func (h *LevelingHandler) handleRewardAddRole(e *events.ComponentInteractionCreate, levelStr string) {
	data, ok := e.Data.(discord.RoleSelectMenuInteractionData)
	if !ok || len(data.Values) == 0 {
		return
	}

	level, err := strconv.Atoi(levelStr)
	if err != nil {
		return
	}

	roleID := data.Values[0]
	if err := h.service.AddRoleReward(*e.GuildID(), level, roleID); err != nil {
		h.logger.Error("failed to add role reward", slog.Any("error", err))
		_ = e.CreateMessage(ui.EphemeralError("報酬の追加に失敗しました: " + err.Error()))
		return
	}

	_ = e.CreateMessage(ui.EphemeralV2(
		discord.NewContainer(
			discord.NewTextDisplay(fmt.Sprintf("Lv.%d → <@&%d> を追加しました。", level, roleID)),
		),
	))
}

func (h *LevelingHandler) handleRewardManage(e *events.ComponentInteractionCreate) {
	rewards, err := h.service.GetRoleRewards(*e.GuildID())
	if err != nil {
		h.logger.Error("failed to get rewards", slog.Any("error", err))
		_ = e.CreateMessage(ui.EphemeralError("報酬一覧の取得に失敗しました。"))
		return
	}
	_ = e.CreateMessage(view.LevelingRewardManagePanel(rewards))
}

func (h *LevelingHandler) handleRewardManageSelect(e *events.ComponentInteractionCreate) {
	data, ok := e.Data.(discord.StringSelectMenuInteractionData)
	if !ok || len(data.Values) == 0 {
		return
	}

	rewardID, err := strconv.ParseInt(data.Values[0], 10, 64)
	if err != nil {
		return
	}

	rewards, err := h.service.GetRoleRewards(*e.GuildID())
	if err != nil {
		return
	}

	for _, r := range rewards {
		if r.ID == rewardID {
			_ = e.UpdateMessage(view.LevelingRewardDetail(r))
			return
		}
	}
}

func (h *LevelingHandler) handleRewardRemove(e *events.ComponentInteractionCreate, idStr string) {
	rewardID, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return
	}

	if err := h.service.RemoveRoleReward(rewardID, *e.GuildID()); err != nil {
		h.logger.Error("failed to remove reward", slog.Any("error", err))
		_ = e.CreateMessage(ui.EphemeralError("報酬の削除に失敗しました。"))
		return
	}

	_ = e.UpdateMessage(discord.NewMessageUpdateV2([]discord.LayoutComponent{
		discord.NewContainer(discord.NewTextDisplay("報酬を削除しました。")),
	}))
}

// --- Multipliers ---

func (h *LevelingHandler) handleRoleMultAddPrompt(e *events.ComponentInteractionCreate) {
	_ = e.CreateMessage(ui.EphemeralV2(
		discord.NewContainer(
			discord.NewTextDisplay("倍率を適用するロールを選択してください:"),
			discord.NewActionRow(
				discord.NewRoleSelectMenu(model.LevelingModuleID+":role_mult_add_role", "ロールを選択..."),
			),
		),
	))
}

func (h *LevelingHandler) handleRoleMultAddRole(e *events.ComponentInteractionCreate) {
	data, ok := e.Data.(discord.RoleSelectMenuInteractionData)
	if !ok || len(data.Values) == 0 {
		return
	}

	roleID := data.Values[0]
	_ = e.Modal(discord.ModalCreate{
		CustomID: fmt.Sprintf("%s:role_mult_modal:%d", model.LevelingModuleID, roleID),
		Title:    "ロールXP倍率",
		Components: []discord.LayoutComponent{
			discord.NewLabel("倍率 (例: 1.5)",
				discord.NewShortTextInput(model.LevelingModuleID+":mult_value").
					WithRequired(true).WithPlaceholder("1.5"),
			),
		},
	})
}

func (h *LevelingHandler) handleRoleMultRemove(e *events.ComponentInteractionCreate, indexStr string) {
	index, err := strconv.Atoi(indexStr)
	if err != nil {
		return
	}

	settings, err := h.service.LoadSettings(*e.GuildID())
	if err != nil {
		return
	}

	if index < 0 || index >= len(settings.RoleMultipliers) {
		return
	}

	settings.RoleMultipliers = append(settings.RoleMultipliers[:index], settings.RoleMultipliers[index+1:]...)
	if err := h.service.SaveSettings(*e.GuildID(), settings); err != nil {
		h.logger.Error("failed to save settings", slog.Any("error", err))
	}

	_ = e.UpdateMessage(view.LevelingMultipliersTab(settings))
}

func (h *LevelingHandler) handleChMultAddPrompt(e *events.ComponentInteractionCreate) {
	_ = e.CreateMessage(ui.EphemeralV2(
		discord.NewContainer(
			discord.NewTextDisplay("倍率を適用するチャンネルを選択してください:"),
			discord.NewActionRow(
				discord.NewChannelSelectMenu(model.LevelingModuleID+":ch_mult_add_ch", "チャンネルを選択...").
					WithChannelTypes(discord.ChannelTypeGuildText),
			),
		),
	))
}

func (h *LevelingHandler) handleChMultAddChannel(e *events.ComponentInteractionCreate) {
	data, ok := e.Data.(discord.ChannelSelectMenuInteractionData)
	if !ok || len(data.Values) == 0 {
		return
	}

	channelID := data.Values[0]
	_ = e.Modal(discord.ModalCreate{
		CustomID: fmt.Sprintf("%s:ch_mult_modal:%d", model.LevelingModuleID, channelID),
		Title:    "チャンネルXP倍率",
		Components: []discord.LayoutComponent{
			discord.NewLabel("倍率 (例: 2.0)",
				discord.NewShortTextInput(model.LevelingModuleID+":mult_value").
					WithRequired(true).WithPlaceholder("2.0"),
			),
		},
	})
}

func (h *LevelingHandler) handleChMultRemove(e *events.ComponentInteractionCreate, indexStr string) {
	index, err := strconv.Atoi(indexStr)
	if err != nil {
		return
	}

	settings, err := h.service.LoadSettings(*e.GuildID())
	if err != nil {
		return
	}

	if index < 0 || index >= len(settings.ChannelMultipliers) {
		return
	}

	settings.ChannelMultipliers = append(settings.ChannelMultipliers[:index], settings.ChannelMultipliers[index+1:]...)
	if err := h.service.SaveSettings(*e.GuildID(), settings); err != nil {
		h.logger.Error("failed to save settings", slog.Any("error", err))
	}

	_ = e.UpdateMessage(view.LevelingMultipliersTab(settings))
}

// --- Leaderboard ---

func (h *LevelingHandler) handleLeaderboardPage(e *events.ComponentInteractionCreate, pageStr string) {
	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 0 {
		return
	}

	entries, err := h.service.GetLeaderboard(*e.GuildID(), 10, page*10)
	if err != nil {
		h.logger.Error("failed to get leaderboard", slog.Any("error", err))
		return
	}

	totalPages := page + 1
	if len(entries) == 10 {
		totalPages = page + 2
	}

	_ = e.UpdateMessage(view.LevelingLeaderboardUpdate(entries, page, totalPages))
}
