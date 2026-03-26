// Copyright (c) 2025-2026 s12kuma01
// SPDX-License-Identifier: MPL-2.0

package handler

import (
	"fmt"
	"log/slog"
	"strconv"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
	"github.com/Sumire-Labs/pedmin/internal/model"
	"github.com/Sumire-Labs/pedmin/internal/ui"
	"github.com/Sumire-Labs/pedmin/internal/view"
)

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
