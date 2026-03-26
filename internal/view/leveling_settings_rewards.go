// Copyright (c) 2025-2026 s12kuma01
// SPDX-License-Identifier: MPL-2.0

package view

import (
	"fmt"

	"github.com/disgoorg/disgo/discord"
	"github.com/Sumire-Labs/pedmin/internal/model"
	"github.com/Sumire-Labs/pedmin/internal/ui"
)

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
