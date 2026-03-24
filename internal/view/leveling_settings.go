// Copyright (c) 2025-2026 s12kuma01
// SPDX-License-Identifier: MPL-2.0

package view

import (
	"fmt"

	"github.com/disgoorg/disgo/discord"
	"github.com/s12kuma01/pedmin/internal/model"
	"github.com/s12kuma01/pedmin/internal/ui"
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

	// Notification mode select
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

	// Show channel select if mode is "channel"
	if settings.NotificationMode == "channel" {
		chRow := discord.NewActionRow(
			discord.NewChannelSelectMenu(model.LevelingModuleID+":notification_ch", "通知チャンネルを選択...").
				WithChannelTypes(discord.ChannelTypeGuildText),
		)
		components = append(components, chRow)
	}

	// Tab buttons
	tabRow := discord.NewActionRow(
		discord.NewPrimaryButton("一般", model.LevelingModuleID+":settings_tab:general").AsDisabled(),
		discord.NewSecondaryButton("報酬", model.LevelingModuleID+":settings_tab:rewards"),
		discord.NewSecondaryButton("倍率", model.LevelingModuleID+":settings_tab:multipliers"),
	)
	components = append(components, tabRow)

	return components
}

// LevelingRewardsTab builds the role rewards management tab.
func LevelingRewardsTab(rewards []model.LevelRoleReward) discord.MessageUpdate {
	components := []discord.ContainerSubComponent{
		discord.NewTextDisplay("### ロール報酬"),
		discord.NewSmallSeparator(),
	}

	if len(rewards) == 0 {
		components = append(components, discord.NewTextDisplay("ロール報酬はまだ設定されていません。"))
	} else {
		var text string
		for _, r := range rewards {
			text += fmt.Sprintf("Lv.%d → <@&%d>\n", r.Level, r.RoleID)
		}
		components = append(components, discord.NewTextDisplay(text))
	}

	addBtn := discord.NewPrimaryButton("報酬追加", model.LevelingModuleID+":reward_add")
	if len(rewards) >= model.MaxRoleRewardsPerGuild {
		addBtn = addBtn.AsDisabled()
	}
	manageBtn := discord.NewSecondaryButton("報酬管理", model.LevelingModuleID+":reward_manage")
	if len(rewards) == 0 {
		manageBtn = manageBtn.AsDisabled()
	}

	components = append(components,
		discord.NewLargeSeparator(),
		discord.NewActionRow(addBtn, manageBtn),
		discord.NewActionRow(
			discord.NewSecondaryButton("一般", model.LevelingModuleID+":settings_tab:general"),
			discord.NewPrimaryButton("報酬", model.LevelingModuleID+":settings_tab:rewards").AsDisabled(),
			discord.NewSecondaryButton("倍率", model.LevelingModuleID+":settings_tab:multipliers"),
		),
	)

	return discord.NewMessageUpdateV2([]discord.LayoutComponent{
		discord.NewContainer(components...),
	})
}

// LevelingMultipliersTab builds the multipliers management tab.
func LevelingMultipliersTab(settings *model.LevelingSettings) discord.MessageUpdate {
	components := []discord.ContainerSubComponent{
		discord.NewTextDisplay("### XP倍率"),
		discord.NewSmallSeparator(),
	}

	// Role multipliers
	roleText := "**ロール倍率:**\n"
	if len(settings.RoleMultipliers) == 0 {
		roleText += "なし\n"
	} else {
		for i, rm := range settings.RoleMultipliers {
			roleText += fmt.Sprintf("<@&%d> → %.1fx", rm.RoleID, rm.Multiplier)
			roleText += fmt.Sprintf("　[`削除`](%s:role_mult_remove:%d)\n", model.LevelingModuleID, i)
		}
	}

	// Channel multipliers
	chText := "**チャンネル倍率:**\n"
	if len(settings.ChannelMultipliers) == 0 {
		chText += "なし\n"
	} else {
		for i, cm := range settings.ChannelMultipliers {
			chText += fmt.Sprintf("<#%d> → %.1fx", cm.ChannelID, cm.Multiplier)
			chText += fmt.Sprintf("　[`削除`](%s:ch_mult_remove:%d)\n", model.LevelingModuleID, i)
		}
	}

	components = append(components, discord.NewTextDisplay(roleText+"\n"+chText))

	totalMults := len(settings.RoleMultipliers) + len(settings.ChannelMultipliers)
	addRoleBtn := discord.NewSecondaryButton("ロール倍率追加", model.LevelingModuleID+":role_mult_add")
	addChBtn := discord.NewSecondaryButton("チャンネル倍率追加", model.LevelingModuleID+":ch_mult_add")
	if totalMults >= model.MaxMultipliersPerGuild {
		addRoleBtn = addRoleBtn.AsDisabled()
		addChBtn = addChBtn.AsDisabled()
	}

	// Remove buttons for role multipliers
	var removeButtons []discord.InteractiveComponent
	for i := range settings.RoleMultipliers {
		removeButtons = append(removeButtons, discord.NewDangerButton(
			fmt.Sprintf("ロール%d削除", i+1),
			fmt.Sprintf("%s:role_mult_remove:%d", model.LevelingModuleID, i),
		))
	}
	for i := range settings.ChannelMultipliers {
		removeButtons = append(removeButtons, discord.NewDangerButton(
			fmt.Sprintf("CH%d削除", i+1),
			fmt.Sprintf("%s:ch_mult_remove:%d", model.LevelingModuleID, i),
		))
	}

	components = append(components,
		discord.NewLargeSeparator(),
		discord.NewActionRow(addRoleBtn, addChBtn),
	)

	if len(removeButtons) > 0 {
		// Discord allows max 5 buttons per action row
		for i := 0; i < len(removeButtons); i += 5 {
			end := i + 5
			if end > len(removeButtons) {
				end = len(removeButtons)
			}
			components = append(components, discord.NewActionRow(removeButtons[i:end]...))
		}
	}

	components = append(components,
		discord.NewActionRow(
			discord.NewSecondaryButton("一般", model.LevelingModuleID+":settings_tab:general"),
			discord.NewSecondaryButton("報酬", model.LevelingModuleID+":settings_tab:rewards"),
			discord.NewPrimaryButton("倍率", model.LevelingModuleID+":settings_tab:multipliers").AsDisabled(),
		),
	)

	return discord.NewMessageUpdateV2([]discord.LayoutComponent{
		discord.NewContainer(components...),
	})
}

// LevelingRewardManagePanel builds the reward list for selection.
func LevelingRewardManagePanel(rewards []model.LevelRoleReward) discord.MessageCreate {
	var options []discord.StringSelectMenuOption
	for _, r := range rewards {
		options = append(options, discord.StringSelectMenuOption{
			Label:       fmt.Sprintf("Lv.%d", r.Level),
			Value:       fmt.Sprintf("%d", r.ID),
			Description: fmt.Sprintf("ロール ID: %d", r.RoleID),
		})
	}

	return ui.EphemeralV2(
		discord.NewContainer(
			discord.NewTextDisplay("### ロール報酬管理"),
			discord.NewSmallSeparator(),
			discord.NewActionRow(
				discord.NewStringSelectMenu(model.LevelingModuleID+":reward_manage_select", "報酬を選択...", options...),
			),
		),
	)
}

// LevelingRewardDetail builds the reward detail view with delete button.
func LevelingRewardDetail(reward model.LevelRoleReward) discord.MessageUpdate {
	return discord.NewMessageUpdateV2([]discord.LayoutComponent{
		discord.NewContainer(
			discord.NewTextDisplay(fmt.Sprintf("### Lv.%d → <@&%d>", reward.Level, reward.RoleID)),
			discord.NewLargeSeparator(),
			discord.NewActionRow(
				discord.NewDangerButton("削除", fmt.Sprintf("%s:reward_remove:%d", model.LevelingModuleID, reward.ID)),
			),
		),
	})
}
