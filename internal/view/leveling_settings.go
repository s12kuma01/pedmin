// Copyright (c) 2025-2026 s12kuma01
// SPDX-License-Identifier: MPL-2.0

package view

import (
	"fmt"

	"github.com/disgoorg/disgo/discord"
	"github.com/s12kuma01/pedmin/internal/model"
)

// LevelingSettingsPanel builds the general settings tab (default view in settings panel).
func LevelingSettingsPanel(settings *model.LevelingSettings, rewardCount int) []discord.LayoutComponent {
	notifLabel := model.NotificationModeLabel(settings.NotificationMode)
	if settings.NotificationMode == "channel" && settings.NotificationChID != 0 {
		notifLabel = fmt.Sprintf("<#%d>", settings.NotificationChID)
	}

	infoDisplay := discord.NewTextDisplay(fmt.Sprintf(
		"**XP範囲:** %d〜%d　**クールダウン:** %ds\n**通知:** %s　**ボイスXP:** %d/分\n**ロール報酬:** %d/%d件",
		settings.MinXP, settings.MaxXP, settings.CooldownSeconds,
		notifLabel, settings.VoiceXPPerMinute,
		rewardCount, model.MaxRoleRewardsPerGuild,
	))

	editRow := discord.NewActionRow(
		discord.NewSecondaryButton("XP範囲", model.LevelingModuleID+":xp_range"),
		discord.NewSecondaryButton("クールダウン", model.LevelingModuleID+":cooldown"),
		discord.NewSecondaryButton("ボイスXP", model.LevelingModuleID+":voice_xp"),
	)

	notifOptions := []discord.StringSelectMenuOption{
		{Label: "同じチャンネル", Value: "same"},
		{Label: "指定チャンネル", Value: "channel"},
		{Label: "オフ", Value: "off"},
	}
	for i, opt := range notifOptions {
		if opt.Value == settings.NotificationMode {
			notifOptions[i] = opt.WithDefault(true)
		}
	}
	notifRow := discord.NewActionRow(
		discord.NewStringSelectMenu(model.LevelingModuleID+":notification_mode", "通知モード", notifOptions...),
	)

	components := []discord.LayoutComponent{infoDisplay, editRow, notifRow}

	if settings.NotificationMode == "channel" {
		chRow := discord.NewActionRow(
			discord.NewChannelSelectMenu(model.LevelingModuleID+":notification_ch", "通知チャンネルを選択...").
				WithChannelTypes(discord.ChannelTypeGuildText),
		)
		components = append(components, chRow)
	}

	tabRow := discord.NewActionRow(
		discord.NewPrimaryButton("一般", model.LevelingModuleID+":settings_tab:general").AsDisabled(),
		discord.NewSecondaryButton("報酬", model.LevelingModuleID+":settings_tab:rewards"),
		discord.NewSecondaryButton("倍率", model.LevelingModuleID+":settings_tab:multipliers"),
	)
	components = append(components, tabRow)

	return components
}
