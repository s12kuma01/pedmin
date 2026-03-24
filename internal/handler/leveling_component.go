// Copyright (c) 2025-2026 s12kuma01
// SPDX-License-Identifier: MPL-2.0

package handler

import (
	"strings"

	"github.com/disgoorg/disgo/events"
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
	case "settings_tab":
		h.handleSettingsTab(e, extra)
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
	case "lb_page":
		h.handleLeaderboardPage(e, extra)
	}
}
